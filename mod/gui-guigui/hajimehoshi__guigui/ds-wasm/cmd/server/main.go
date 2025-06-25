package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/starfederation/datastar/sdk/go"
)

var port = flag.Int("port", 8081, "Port to run server on")

func main() {
	flag.Parse()

	fmt.Printf("🌟 DataStar Hello World - Server Mode\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("🚀 Starting server on http://localhost:%d\n", *port)
	fmt.Printf("📋 Available endpoints:\n")
	fmt.Printf("  GET  /           - Hello World page\n")
	fmt.Printf("  GET  /health     - Health check\n")
	fmt.Printf("  POST /click      - Handle button click\n")
	fmt.Printf("  GET  /time       - Current time\n")
	fmt.Printf("\n")

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/click", clickHandler)
	http.HandleFunc("/time", timeHandler)

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>DataStar Hello World - Server Mode</title>
    <script type="module" defer src="https://cdn.jsdelivr.net/npm/@starfederation/datastar@latest"></script>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .container { background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0; }
        button { background: #007cba; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; margin: 5px; }
        button:hover { background: #005a87; }
        .result { background: white; padding: 10px; border-radius: 4px; margin: 10px 0; }
        .mode-badge { background: #28a745; color: white; padding: 5px 10px; border-radius: 4px; font-size: 12px; }
    </style>
</head>
<body>
    <h1>🌟 DataStar Hello World</h1>
    <div class="mode-badge">SERVER MODE</div>
    
    <div class="container">
        <h2>Basic Example</h2>
        <div data-store='{"count": 0, "message": "Hello from server!"}'>
            <p data-text="$message"></p>
            <p>Count: <span data-text="$count"></span></p>
            <button data-on-click="$count++">Increment</button>
            <button data-on-click="$count--">Decrement</button>
            <button data-on-click="$count = 0">Reset</button>
        </div>
    </div>

    <div class="container">
        <h2>Server Interaction</h2>
        <div id="server-result" class="result">
            Click the button to interact with the server
        </div>
        <button 
            data-post="/click"
            data-target="#server-result">
            Click Me (Server)
        </button>
    </div>

    <div class="container">
        <h2>Live Time</h2>
        <div id="time-display" class="result">
            <span data-text="$currentTime">Loading...</span>
        </div>
        <button 
            data-get="/time"
            data-target="#time-display">
            Get Current Time
        </button>
    </div>

    <div class="container">
        <h2>Mode Comparison</h2>
        <p>This is the <strong>server mode</strong> of the DataStar hello-world example.</p>
        <p>🔗 <a href="http://localhost:8082" target="_blank">Open WASM Mode</a> (if running)</p>
        <p>The same DataStar application running in two different modes:</p>
        <ul>
            <li><strong>Server Mode:</strong> Traditional Go web server</li>
            <li><strong>WASM Mode:</strong> Go compiled to WebAssembly running in browser</li>
        </ul>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "mode": "server", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

func clickHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use DataStar to send back a fragment
	sse := datastar.NewSSE(w, r)

	clickTime := time.Now().Format("15:04:05")
	fragment := fmt.Sprintf(`
		<div class="result">
			✅ Server responded at %s<br>
			🖥️ Running in server mode<br>
			🔄 Request processed by Go HTTP server
		</div>
	`, clickTime)

	sse.MergeFragments(fragment, datastar.WithSelector("#server-result"))
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Send back data to update the store
	sse.MarshalAndMergeSignals(map[string]interface{}{
		"currentTime": currentTime,
	})
}
