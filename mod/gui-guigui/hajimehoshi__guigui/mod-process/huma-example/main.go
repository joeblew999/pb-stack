package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/sse" // Huma's SSE package
	"github.com/go-chi/chi/v5"
)

// Message defines the structure of a regular event message.
type Message struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

// HeartbeatEvent defines the structure of a heartbeat event.
type HeartbeatEvent struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
}

// UserConnectedEvent defines the structure for a new user connection event.
type UserConnectedEvent struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"userId"`
}

func main() {
	// Create a new Chi router
	router := chi.NewRouter()

	// Create a new Huma API with the router
	api := humachi.New(router, huma.DefaultConfig("SSE Example API", "1.0.0"))

	// --- SSE Endpoint Definition ---
	// Register the SSE endpoint using sse.Register instead of huma.Register.
	// The third argument is a map that defines the different event types (Go structs)
	// that this SSE stream can send. Huma uses this to generate the OpenAPI spec.
	sse.Register(api, huma.Operation{
		OperationID: "subscribeToEvents",
		Method:      http.MethodGet,
		Path:        "/events",
		Summary:     "Subscribe to real-time events via Server-Sent Events (SSE)",
		Description: "Streams various real-time updates including messages, heartbeats, and user connection notifications.",
		Tags:        []string{"SSE"},
	}, map[string]any{
		"message":       Message{},            // Represents regular messages
		"heartbeat":     HeartbeatEvent{},     // Represents heartbeat pings
		"userConnected": UserConnectedEvent{}, // Represents new user connections
	}, func(ctx context.Context, input *struct{}, send sse.Sender) {
		// This function is executed when a client connects to the SSE endpoint.
		log.Println("New SSE client connected.")

		// Simulate a new user connection event immediately
		send.Data(UserConnectedEvent{
			Event:     "userConnected",
			Timestamp: time.Now(),
			UserID:    fmt.Sprintf("user-%d", time.Now().UnixNano()),
		})

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop() // Ensure the ticker is stopped when the goroutine exits

		messageCounter := 0

		for {
			select {
			case <-ctx.Done():
				// Context cancelled (client disconnected or server shutting down)
				log.Println("SSE client disconnected.")
				return
			case t := <-ticker.C:
				// Send a regular message event
				messageCounter++
				msg := Message{
					ID:        fmt.Sprintf("msg-%d", messageCounter),
					Timestamp: t,
					Content:   fmt.Sprintf("This is a streamed message number %d.", messageCounter),
				}
				send.Data(msg)
				log.Printf("Sent message event: %s", msg.ID)

				// Every 5th message, send a heartbeat
				if messageCounter%5 == 0 {
					heartbeat := HeartbeatEvent{
						Event:     "heartbeat",
						Timestamp: time.Now(),
					}
					send.Data(heartbeat)
					log.Println("Sent heartbeat event.")
				}
			}
		}
	})

	// --- Serve the OpenAPI Specification ---
	// Huma automatically serves the OpenAPI spec at /openapi.json and /openapi.yaml.
	// You can view it in your browser or use tools like Swagger UI.
	// Example: go to http://localhost:8888/openapi.yaml or http://localhost:8888/openapi.json

	// Start the HTTP server.
	fmt.Println("Huma API server started on http://localhost:8888")
	fmt.Println("SSE Endpoint: http://localhost:8888/events")
	fmt.Println("OpenAPI Spec (JSON): http://localhost:8888/openapi.json")
	log.Fatal(http.ListenAndServe(":8888", router))
}
