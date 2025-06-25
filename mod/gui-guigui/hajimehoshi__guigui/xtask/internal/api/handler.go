package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"xtask/internal/tools"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	natsgo "github.com/nats-io/nats.go"
)

// Handler handles HTTP API requests
type Handler struct {
	natsConn *natsgo.Conn
}

// CommandRequest represents a command execution request
type CommandRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// CommandResponse represents a command execution response
type CommandResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// WhichResponse represents a which tool response
type WhichResponse struct {
	Found bool   `json:"found"`
	Path  string `json:"path"`
	Error string `json:"error,omitempty"`
}

// DownloadRequest represents a download request
type DownloadRequest struct {
	URL    string `json:"url"`
	Output string `json:"output"`
}

// DownloadResponse represents a download response
type DownloadResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// New creates a new API handler
func New(natsConn *natsgo.Conn) *Handler {
	return &Handler{
		natsConn: natsConn,
	}
}

// SetupRoutes sets up the HTTP routes
func (h *Handler) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	// CORS middleware for web UI
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Health check
	r.Get("/health", h.handleHealth)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/tasks", h.handleExecuteCommand)

		r.Route("/tools", func(r chi.Router) {
			r.Get("/which/{binary}", h.handleWhich)
			r.Post("/got", h.handleDownload)
		})
	})

	// Web UI (serve static files)
	r.Get("/web", h.handleWebUI)
	r.Get("/web/*", h.handleWebUI)

	return r
}

// handleHealth handles health check requests
func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "xtask-server",
	})
}

// handleExecuteCommand handles command execution requests
func (h *Handler) handleExecuteCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute command locally for now
	// TODO: Implement NATS-based distributed execution
	output, err := h.executeCommandLocally(req.Command, req.Args)

	response := CommandResponse{
		Success: err == nil,
		Output:  output,
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleWhich handles which tool requests
func (h *Handler) handleWhich(w http.ResponseWriter, r *http.Request) {
	binary := chi.URLParam(r, "binary")
	if binary == "" {
		http.Error(w, "Binary name required", http.StatusBadRequest)
		return
	}

	path, err := h.findBinary(binary)
	response := WhichResponse{
		Found: err == nil,
		Path:  path,
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleDownload handles download requests
func (h *Handler) handleDownload(w http.ResponseWriter, r *http.Request) {
	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use tools package for download
	err := h.downloadFile(req.URL, req.Output)

	response := DownloadResponse{
		Success: err == nil,
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleWebUI serves the web UI
func (h *Handler) handleWebUI(w http.ResponseWriter, r *http.Request) {
	// Serve the DataStar-powered web UI from web/index.html
	indexPath := filepath.Join("web", "index.html")

	// Check if the file exists
	if _, err := os.Stat(indexPath); err != nil {
		// Fallback to basic HTML if file doesn't exist
		h.serveBasicHTML(w, r)
		return
	}

	// Read and serve the index.html file
	htmlContent, err := os.ReadFile(indexPath)
	if err != nil {
		log.Printf("Error reading web/index.html: %v", err)
		h.serveBasicHTML(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(htmlContent)
}

// serveBasicHTML serves a basic fallback HTML when the main index.html is not available
func (h *Handler) serveBasicHTML(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>xtask Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { color: #333; }
        .status { color: green; }
        .section { margin: 20px 0; }
    </style>
</head>
<body>
    <h1 class="header">🚀 xtask Dashboard</h1>
    <div class="section">
        <h2>Status</h2>
        <p class="status">✅ Server is running</p>
        <p>📡 NATS JetStream is active</p>
    </div>
    <div class="section">
        <h2>API Endpoints</h2>
        <ul>
            <li><code>GET /health</code> - Health check</li>
            <li><code>POST /api/v1/tasks</code> - Execute commands</li>
            <li><code>GET /api/v1/tools/which/{binary}</code> - Find binary</li>
            <li><code>POST /api/v1/tools/got</code> - Download files</li>
        </ul>
    </div>
    <div class="section">
        <p><em>Note: Advanced DataStar UI not found. Using fallback interface.</em></p>
    </div>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// executeCommandLocally executes a command locally
func (h *Handler) executeCommandLocally(command string, args []string) (string, error) {
	switch command {
	case "which":
		return h.executeWhich(args)
	case "got":
		return h.executeGot(args)
	case "silent":
		return h.executeSilent(args)
	case "kill-port":
		return h.executeKillPort(args)
	case "wait-for-port":
		return h.executeWaitForPort(args)
	case "tree":
		return h.executeTree(args)
	case "health-check":
		return h.executeHealthCheck(args)
	default:
		// Default to task execution
		return h.executeTask(append([]string{command}, args...))
	}
}

// executeWhich implements the which command
func (h *Handler) executeWhich(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("binary name required")
	}

	path, err := h.findBinary(args[0])
	if err != nil {
		return "", err
	}

	return path, nil
}

// executeGot implements the got (download) command
func (h *Handler) executeGot(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("URL and output path required")
	}

	err := h.downloadFile(args[0], args[1])
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Downloaded %s to %s", args[0], args[1]), nil
}

// executeSilent implements the silent command
func (h *Handler) executeSilent(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("command required")
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, _ := cmd.CombinedOutput() // Ignore errors for silent execution
	return string(output), nil
}

// executeKillPort implements the kill-port command
func (h *Handler) executeKillPort(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("port number required")
	}

	// Use tools package if available, otherwise implement basic version
	tools.ExecuteKillPort(args)
	return fmt.Sprintf("Attempted to kill processes on port %s", args[0]), nil
}

// executeWaitForPort implements the wait-for-port command
func (h *Handler) executeWaitForPort(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("port number required")
	}

	// Use tools package if available, otherwise implement basic version
	tools.ExecuteWaitForPort(args)
	return fmt.Sprintf("Waited for port %s", args[0]), nil
}

// executeTree implements the tree command
func (h *Handler) executeTree(args []string) (string, error) {
	// Use tools package if available, otherwise implement basic version
	tools.ExecuteTree(args)
	return "Tree command executed", nil
}

// executeHealthCheck implements the health-check command
func (h *Handler) executeHealthCheck(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("URL required")
	}

	// Use tools package if available, otherwise implement basic version
	tools.ExecuteHealthCheck(args)
	return fmt.Sprintf("Health check for %s completed", args[0]), nil
}

// executeTask implements task execution
func (h *Handler) executeTask(args []string) (string, error) {
	// Use tools package for task execution
	tools.ExecuteTask(args)
	return "Task executed", nil
}

// findBinary finds the path to a binary executable
func (h *Handler) findBinary(binary string) (string, error) {
	// Add common executable extensions on Windows
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(binary, ".exe") && !strings.HasSuffix(binary, ".cmd") && !strings.HasSuffix(binary, ".bat") {
			// Try with .exe first
			if path, err := exec.LookPath(binary + ".exe"); err == nil {
				return path, nil
			}
			// Try with .cmd
			if path, err := exec.LookPath(binary + ".cmd"); err == nil {
				return path, nil
			}
			// Try with .bat
			if path, err := exec.LookPath(binary + ".bat"); err == nil {
				return path, nil
			}
		}
	}

	// Try the binary as-is
	path, err := exec.LookPath(binary)
	if err != nil {
		return "", fmt.Errorf("binary '%s' not found in PATH", binary)
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path, nil // Return relative path if absolute fails
	}

	return absPath, nil
}

// downloadFile downloads a file from URL to the specified output path
func (h *Handler) downloadFile(url, output string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Use tools package for download if available
	// For now, implement a basic version using curl or wget
	var cmd *exec.Cmd

	// Try curl first (more common)
	if _, err := exec.LookPath("curl"); err == nil {
		cmd = exec.Command("curl", "-L", "-o", output, url)
	} else if _, err := exec.LookPath("wget"); err == nil {
		cmd = exec.Command("wget", "-O", output, url)
	} else {
		return fmt.Errorf("neither curl nor wget found for downloading")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Verify file was created
	if _, err := os.Stat(output); err != nil {
		return fmt.Errorf("download completed but file not found: %w", err)
	}

	return nil
}
