package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"process-compose-huma/api"
	"process-compose-huma/wrapper"
)

// waitForProcessComposeAPI waits for the process-compose API to be available
func waitForProcessComposeAPI() error {
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get("http://localhost:8080/live")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		log.Printf("⏳ Waiting for Process Compose API... (%d/%d)", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("process-compose API not available after %d attempts", maxRetries)
}

func main() {
	log.Println("🚀 Starting Enhanced Process Compose API...")

	// Wait for existing process-compose API to be available
	log.Println("🔍 Waiting for Process Compose API at http://localhost:8080...")
	if err := waitForProcessComposeAPI(); err != nil {
		log.Fatalf("Process Compose API not available: %v", err)
	}
	log.Println("✅ Process Compose API is ready")

	// Initialize API client (no process management)
	controller := wrapper.NewProcessComposeClient()

	// Create Chi router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))

	// Create Huma API configuration
	config := huma.DefaultConfig("Process Compose API", "1.0.0")
	config.Info.Description = "Enhanced Process Compose API with real-time SSE support powered by Huma v2"
	config.Servers = []*huma.Server{
		{URL: "http://localhost:8888", Description: "Development server"},
	}

	// Create Huma API instance
	humaAPI := humachi.New(router, config)

	// Initialize API handlers
	apiHandlers := api.NewAPIHandlers(controller)
	sseHandlers := api.NewSSEHandlers(controller)

	// Register REST API routes
	apiHandlers.RegisterRoutes(humaAPI)

	// Register SSE routes
	sseHandlers.RegisterSSERoutes(humaAPI)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8888",
		Handler: router,
	}

	// Graceful shutdown handling
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("🛑 Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	// Start the server
	log.Printf("🌐 Enhanced Process Compose API starting on http://localhost:8888")
	log.Printf("📚 API Documentation: http://localhost:8888/docs")
	log.Printf("🔄 SSE Process Events: http://localhost:8888/events/processes")
	log.Printf("📝 SSE Log Events: http://localhost:8888/events/logs")
	log.Printf("🔧 SSE System Events: http://localhost:8888/events/system")
	log.Printf("🌐 SSE All Events: http://localhost:8888/events/all")
	log.Printf("💚 Health Check: http://localhost:8888/health")
	log.Printf("📊 Original Process Compose API: http://localhost:8080")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
