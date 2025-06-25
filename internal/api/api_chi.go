package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/go-chi/chi/v5"
	natsgo "github.com/nats-io/nats.go"
	"github.com/starfederation/datastar/sdk/go/datastar"

	"xtask/internal/tools"
)

type Handler struct {
	nats *natsgo.Conn
}

func New(natsConn *natsgo.Conn) *Handler {
	return &Handler{
		nats: natsConn,
	}
}

// SetupRoutes configures Chi router with all xtask endpoints
func (h *Handler) SetupRoutes() chi.Router {
	r := chi.NewRouter()

	// Health check
	r.Get("/health", h.HandleHealth)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Tool endpoints
		r.Route("/tools", func(r chi.Router) {
			r.Get("/which/{binary}", h.HandleWhich)
			r.Post("/got", h.HandleGot)
			r.Post("/silent", h.HandleSilent)
			r.Post("/kill-port", h.HandleKillPort)
			r.Post("/wait-for-port", h.HandleWaitForPort)
			r.Get("/tree", h.HandleTree)
			r.Post("/health-check", h.HandleHealthCheck)
		})

		// Node management
		r.Get("/nodes", h.HandleNodes)
		r.Post("/tasks", h.HandleTasks)
	})

	// Server-Sent Events for DataStar
	r.Get("/events", h.HandleSSE)

	// Static files for web interface
	r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("./web/"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusMovedPermanently)
	})

	return r
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "dev",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleWhich(w http.ResponseWriter, r *http.Request) {
	binary := chi.URLParam(r, "binary")

	path, err := exec.LookPath(binary)
	response := map[string]interface{}{}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response["found"] = false
		response["error"] = err.Error()
	} else {
		response["found"] = true
		response["path"] = path
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleGot(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL    string `json:"url"`
		Output string `json:"output"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   "URL is required",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Execute download using tools package
	args := []string{req.URL}
	if req.Output != "" {
		args = append(args, "-o", req.Output)
	}

	// This is a simplified version - in reality we'd capture output
	tools.ExecuteGot(args)

	response := map[string]interface{}{
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleSilent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Command == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   "Command is required",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	args := append([]string{req.Command}, req.Args...)
	tools.ExecuteSilent(args)

	response := map[string]interface{}{
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleKillPort(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Port string `json:"port"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Port == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   "Port is required",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	tools.ExecuteKillPort([]string{req.Port})

	response := map[string]interface{}{
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleWaitForPort(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Port    string `json:"port"`
		Timeout string `json:"timeout"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Port == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   "Port is required",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	args := []string{req.Port}
	if req.Timeout != "" {
		args = append(args, req.Timeout)
	}

	tools.ExecuteWaitForPort(args)

	response := map[string]interface{}{
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleTree(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}

	tools.ExecuteTree([]string{path})

	response := map[string]interface{}{
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   "URL is required",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	tools.ExecuteHealthCheck([]string{req.URL})

	response := map[string]interface{}{
		"success": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleNodes(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement node discovery via NATS
	response := map[string]interface{}{
		"nodes": []map[string]interface{}{
			{
				"id":       "local",
				"hostname": "localhost",
				"status":   "online",
				"capabilities": []string{
					"task", "which", "got", "silent", "kill-port",
				},
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Command == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"success": false,
			"error":   "Command is required",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Execute command based on type
	switch req.Command {
	case "which":
		tools.ExecuteWhich(req.Args)
	case "got":
		tools.ExecuteGot(req.Args)
	case "silent":
		tools.ExecuteSilent(req.Args)
	case "kill-port":
		tools.ExecuteKillPort(req.Args)
	case "wait-for-port":
		tools.ExecuteWaitForPort(req.Args)
	case "tree":
		tools.ExecuteTree(req.Args)
	case "health-check":
		tools.ExecuteHealthCheck(req.Args)
	default:
		args := append([]string{req.Command}, req.Args...)
		tools.ExecuteTask(args)
	}

	response := map[string]interface{}{
		"success": true,
		"output":  "Command executed successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleSSE provides Server-Sent Events driven by NATS for DataStar
func (h *Handler) HandleSSE(w http.ResponseWriter, r *http.Request) {
	// Use DataStar SSE helper
	sse := datastar.NewSSE(w, r)

	// Get context and project scope
	project := r.URL.Query().Get("project")
	if project == "" {
		project = "default"
	}

	// Subscribe to NATS events for this project
	subject := fmt.Sprintf("xtask.events.%s.>", project)
	sub, err := h.nats.Subscribe(subject, func(msg *natsgo.Msg) {
		// Convert NATS message to DataStar event
		event := map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"subject":   msg.Subject,
			"data":      string(msg.Data),
		}

		// Send as DataStar fragment
		sse.MergeFragments(fmt.Sprintf(`
			<div id="event-%d" class="event-item">
				<span class="timestamp">%s</span>
				<span class="subject">%s</span>
				<span class="data">%s</span>
			</div>
		`, event["timestamp"], time.Now().Format("15:04:05"), msg.Subject, string(msg.Data)))
	})

	if err != nil {
		log.Printf("Failed to subscribe to NATS events: %v", err)
		return
	}
	defer sub.Unsubscribe()

	// Keep connection alive and handle client disconnect
	<-r.Context().Done()
}

// PublishEvent publishes an event to NATS for SSE streaming
func (h *Handler) PublishEvent(project, eventType string, data interface{}) error {
	subject := fmt.Sprintf("xtask.events.%s.%s", project, eventType)

	eventData, err := json.Marshal(map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().Unix(),
		"data":      data,
	})
	if err != nil {
		return err
	}

	return h.nats.Publish(subject, eventData)
}
