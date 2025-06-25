package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ProcessState represents a process from the API
type ProcessState struct {
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	PID       int     `json:"pid"`
	IsRunning bool    `json:"is_running"`
	Restarts  int     `json:"restarts"`
	ExitCode  int     `json:"exit_code"`
	Mem       int64   `json:"mem"`
	CPU       float64 `json:"cpu"`
}

// ProcessesResponse from the API
type ProcessesResponse struct {
	Data []ProcessState `json:"data"`
}

// ProjectState from the API
type ProjectState struct {
	HostName          string `json:"hostName"`
	Version           string `json:"version"`
	StartTime         string `json:"startTime"`
	ProcessNum        int    `json:"processNum"`
	RunningProcessNum int    `json:"runningProcessNum"`
}

const processComposeAPIURL = "http://localhost:8888"

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Serve static files (CSS, JS)
	router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Main dashboard page
	router.Get("/", handleDashboard)

	// API endpoints for Datastar
	router.Get("/api/processes", handleProcesses)
	router.Get("/api/project", handleProject)
	router.Post("/api/process/{name}/start", handleStartProcess)
	router.Post("/api/process/{name}/stop", handleStopProcess)
	router.Post("/api/process/{name}/restart", handleRestartProcess)

	// SSE endpoint for real-time updates
	router.Get("/events", handleSSE)

	log.Println("🌟 Process Compose Datastar GUI starting on http://localhost:3000")
	log.Println("📊 Connecting to Process Compose API at", processComposeAPIURL)
	log.Fatal(http.ListenAndServe(":3000", router))
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Process Compose Dashboard</title>
    <script type="module" src="https://cdn.jsdelivr.net/npm/@sudodevnull/datastar@latest/dist/datastar.js"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        .status-running { @apply bg-green-100 text-green-800; }
        .status-completed { @apply bg-blue-100 text-blue-800; }
        .status-failed { @apply bg-red-100 text-red-800; }
        .status-pending { @apply bg-yellow-100 text-yellow-800; }
    </style>
</head>
<body class="bg-gray-50 min-h-screen">
    <div class="container mx-auto px-4 py-8" data-store='{"processes": [], "project": {}, "events": []}'>
        <!-- Header -->
        <div class="mb-8">
            <h1 class="text-3xl font-bold text-gray-900 mb-2">Process Compose Dashboard</h1>
            <div class="flex items-center space-x-4 text-sm text-gray-600">
                <span data-text="project.hostName"></span>
                <span data-text="project.version"></span>
                <span>Processes: <span data-text="project.processNum"></span></span>
                <span>Running: <span data-text="project.runningProcessNum"></span></span>
            </div>
        </div>

        <!-- Auto-refresh controls -->
        <div class="mb-6 flex justify-between items-center">
            <button 
                class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                data-on-click="$get('/api/processes')"
            >
                🔄 Refresh
            </button>
            <div class="text-sm text-gray-500">
                Last updated: <span id="lastUpdate"></span>
            </div>
        </div>

        <!-- Process List -->
        <div class="bg-white shadow rounded-lg overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-200">
                <h2 class="text-lg font-semibold text-gray-900">Processes</h2>
            </div>
            <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                    <thead class="bg-gray-50">
                        <tr>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">PID</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Restarts</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Memory</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">CPU</th>
                            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="bg-white divide-y divide-gray-200" data-for="process in processes">
                        <tr>
                            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900" data-text="process.name"></td>
                            <td class="px-6 py-4 whitespace-nowrap">
                                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full" 
                                      data-class="'status-' + process.status.toLowerCase()"
                                      data-text="process.status">
                                </span>
                            </td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500" data-text="process.pid"></td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500" data-text="process.restarts"></td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500" data-text="(process.mem / 1024 / 1024).toFixed(1) + ' MB'"></td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500" data-text="process.cpu.toFixed(2) + '%'"></td>
                            <td class="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                                <button
                                    class="bg-green-500 hover:bg-green-700 text-white font-bold py-1 px-2 rounded text-xs"
                                    data-on-click="$post('/api/process/' + process.name + '/start')"
                                >
                                    ▶️ Start
                                </button>
                                <button
                                    class="bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-2 rounded text-xs"
                                    data-on-click="$post('/api/process/' + process.name + '/stop')"
                                >
                                    ⏹️ Stop
                                </button>
                                <button 
                                    class="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-1 px-2 rounded text-xs"
                                    data-on-click="$post('/api/process/' + process.name + '/restart')"
                                >
                                    🔄 Restart
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Events Log -->
        <div class="mt-8 bg-white shadow rounded-lg overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-200">
                <h2 class="text-lg font-semibold text-gray-900">Live Events</h2>
            </div>
            <div class="p-4 h-64 overflow-y-auto bg-gray-900 text-green-400 font-mono text-sm" id="events-log">
                <div data-for="event in events.slice(-20)">
                    <div data-text="event.timestamp + ' - ' + event.type + ': ' + event.message"></div>
                </div>
            </div>
        </div>
    </div>

    <script>
        // Auto-refresh every 5 seconds
        setInterval(() => {
            document.querySelector('[data-on-click*="get"]').click();
            document.getElementById('lastUpdate').textContent = new Date().toLocaleTimeString();
        }, 5000);

        // Initial load
        document.addEventListener('DOMContentLoaded', () => {
            document.querySelector('[data-on-click*="get"]').click();
            document.getElementById('lastUpdate').textContent = new Date().toLocaleTimeString();
        });

        // Connect to SSE for real-time events
        const eventSource = new EventSource('/events');
        eventSource.onmessage = function(event) {
            const data = JSON.parse(event.data);
            const store = window.datastar.store;
            store.events = store.events || [];
            store.events.push({
                timestamp: new Date().toLocaleTimeString(),
                type: data.type || 'info',
                message: data.message || JSON.stringify(data)
            });
            // Keep only last 50 events
            if (store.events.length > 50) {
                store.events = store.events.slice(-50);
            }
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func handleProcesses(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(processComposeAPIURL + "/processes")
	if err != nil {
		http.Error(w, "Failed to fetch processes", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var processes ProcessesResponse
	if err := json.NewDecoder(resp.Body).Decode(&processes); err != nil {
		http.Error(w, "Failed to decode processes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("datastar-merge-store", "true")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"processes": processes.Data,
	})
}

func handleProject(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(processComposeAPIURL + "/project")
	if err != nil {
		http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var project ProjectState
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		http.Error(w, "Failed to decode project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("datastar-merge-store", "true")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"project": project,
	})
}

func handleStartProcess(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	resp, err := http.Post(processComposeAPIURL+"/process/"+name+"/start", "application/json", nil)
	if err != nil {
		http.Error(w, "Failed to start process", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("datastar-merge-store", "true")
	w.WriteHeader(http.StatusOK)
}

func handleStopProcess(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	resp, err := http.Post(processComposeAPIURL+"/process/"+name+"/stop", "application/json", nil)
	if err != nil {
		http.Error(w, "Failed to stop process", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("datastar-merge-store", "true")
	w.WriteHeader(http.StatusOK)
}

func handleRestartProcess(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	resp, err := http.Post(processComposeAPIURL+"/process/"+name+"/restart", "application/json", nil)
	if err != nil {
		http.Error(w, "Failed to restart process", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("datastar-merge-store", "true")
	w.WriteHeader(http.StatusOK)
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Connect to Process Compose SSE
	resp, err := http.Get(processComposeAPIURL + "/events/all")
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\": \"Failed to connect to Process Compose events\"}\n\n")
		return
	}
	defer resp.Body.Close()

	// Stream events from Process Compose to client
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading SSE stream: %v", err)
			}
			break
		}
		w.Write(buffer[:n])
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}
