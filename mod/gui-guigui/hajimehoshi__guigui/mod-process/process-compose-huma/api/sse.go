package api

import (
	"context"
	"log"
	"time"

	"process-compose-huma/wrapper"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
)

// SSEHandlers contains all the SSE endpoint handlers
type SSEHandlers struct {
	controller wrapper.ProcessController
}

// NewSSEHandlers creates a new SSE handlers instance
func NewSSEHandlers(controller wrapper.ProcessController) *SSEHandlers {
	return &SSEHandlers{
		controller: controller,
	}
}

// ProcessEventsHandler handles SSE connections for process events
func (s *SSEHandlers) ProcessEventsHandler(ctx context.Context, input *struct{}, send sse.Sender) {
	log.Println("🔄 New SSE client connected to process events")

	// Send initial state
	processes, err := s.controller.GetProcesses()
	if err != nil {
		log.Printf("Error getting initial processes: %v", err)
		return
	}

	for _, process := range processes.Data {
		send.Data(wrapper.ProcessStatusEvent{
			Type:      "initial_state",
			Process:   process,
			Timestamp: time.Now(),
		})
	}

	// Listen for real-time events
	for {
		select {
		case <-ctx.Done():
			log.Println("🔌 SSE client disconnected from process events")
			return
		case event := <-s.controller.ProcessEvents():
			send.Data(event)
		}
	}
}

// LogEventsHandler handles SSE connections for log events
func (s *SSEHandlers) LogEventsHandler(ctx context.Context, input *struct{}, send sse.Sender) {
	log.Println("📝 New SSE client connected to log events")

	for {
		select {
		case <-ctx.Done():
			log.Println("🔌 SSE client disconnected from log events")
			return
		case event := <-s.controller.LogEvents():
			send.Data(event)
		}
	}
}

// SystemEventsHandler handles SSE connections for system events
func (s *SSEHandlers) SystemEventsHandler(ctx context.Context, input *struct{}, send sse.Sender) {
	log.Println("🔧 New SSE client connected to system events")

	for {
		select {
		case <-ctx.Done():
			log.Println("🔌 SSE client disconnected from system events")
			return
		case event := <-s.controller.SystemEvents():
			send.Data(event)
		}
	}
}

// AllEventsHandler handles SSE connections for all events combined
func (s *SSEHandlers) AllEventsHandler(ctx context.Context, input *struct{}, send sse.Sender) {
	log.Println("🌐 New SSE client connected to all events")

	// Send initial process states
	processes, err := s.controller.GetProcesses()
	if err != nil {
		log.Printf("Error getting initial processes: %v", err)
		return
	}

	for _, process := range processes.Data {
		send.Data(wrapper.ProcessStatusEvent{
			Type:      "initial_state",
			Process:   process,
			Timestamp: time.Now(),
		})
	}

	// Listen for all event types
	for {
		select {
		case <-ctx.Done():
			log.Println("🔌 SSE client disconnected from all events")
			return
		case event := <-s.controller.ProcessEvents():
			send.Data(event)
		case event := <-s.controller.LogEvents():
			send.Data(event)
		case event := <-s.controller.SystemEvents():
			send.Data(event)
		}
	}
}

// RegisterSSERoutes registers all SSE routes with Huma
func (s *SSEHandlers) RegisterSSERoutes(api huma.API) {
	// Process events SSE endpoint
	sse.Register(api, huma.Operation{
		OperationID: "subscribeToProcessEvents",
		Method:      "GET",
		Path:        "/events/processes",
		Summary:     "Subscribe to process events",
		Description: "Real-time stream of process status changes via Server-Sent Events",
		Tags:        []string{"SSE", "Real-time"},
	}, map[string]any{
		"process_status_changed":        wrapper.ProcessStatusEvent{},
		"process_added":                 wrapper.ProcessStatusEvent{},
		"process_removed":               wrapper.ProcessStatusEvent{},
		"process_restarted":             wrapper.ProcessStatusEvent{},
		"process_start_requested":       wrapper.ProcessStatusEvent{},
		"process_stop_requested":        wrapper.ProcessStatusEvent{},
		"process_restart_requested":     wrapper.ProcessStatusEvent{},
		"process_restart_count_changed": wrapper.ProcessStatusEvent{},
		"initial_state":                 wrapper.ProcessStatusEvent{},
	}, s.ProcessEventsHandler)

	// Log events SSE endpoint
	sse.Register(api, huma.Operation{
		OperationID: "subscribeToLogEvents",
		Method:      "GET",
		Path:        "/events/logs",
		Summary:     "Subscribe to log events",
		Description: "Real-time stream of process log output via Server-Sent Events",
		Tags:        []string{"SSE", "Real-time", "Logs"},
	}, map[string]any{
		"log_line": wrapper.LogEvent{},
	}, s.LogEventsHandler)

	// System events SSE endpoint
	sse.Register(api, huma.Operation{
		OperationID: "subscribeToSystemEvents",
		Method:      "GET",
		Path:        "/events/system",
		Summary:     "Subscribe to system events",
		Description: "Real-time stream of system-level events via Server-Sent Events",
		Tags:        []string{"SSE", "Real-time", "System"},
	}, map[string]any{
		"system_started": wrapper.SystemEvent{},
		"system_stopped": wrapper.SystemEvent{},
	}, s.SystemEventsHandler)

	// All events combined SSE endpoint
	sse.Register(api, huma.Operation{
		OperationID: "subscribeToAllEvents",
		Method:      "GET",
		Path:        "/events/all",
		Summary:     "Subscribe to all events",
		Description: "Real-time stream of all system events via Server-Sent Events",
		Tags:        []string{"SSE", "Real-time"},
	}, map[string]any{
		"process_status_changed":        wrapper.ProcessStatusEvent{},
		"process_added":                 wrapper.ProcessStatusEvent{},
		"process_removed":               wrapper.ProcessStatusEvent{},
		"process_restarted":             wrapper.ProcessStatusEvent{},
		"process_start_requested":       wrapper.ProcessStatusEvent{},
		"process_stop_requested":        wrapper.ProcessStatusEvent{},
		"process_restart_requested":     wrapper.ProcessStatusEvent{},
		"process_restart_count_changed": wrapper.ProcessStatusEvent{},
		"initial_state":                 wrapper.ProcessStatusEvent{},
		"log_line":                      wrapper.LogEvent{},
		"system_started":                wrapper.SystemEvent{},
		"system_stopped":                wrapper.SystemEvent{},
	}, s.AllEventsHandler)
}
