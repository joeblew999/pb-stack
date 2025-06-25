package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"opencloud/pkg/server"
)

var (
	indexDir = flag.String("index", "./index", "Index directory for the server")
	port     = flag.String("port", "8080", "Server port")
	dataDir  = flag.String("data", "./data", "Data directory for documents")
	host     = flag.String("host", "localhost", "Server host")
	debug    = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	flag.Parse()

	fmt.Println("🌩️  OpenCloud - Collaboration Server")
	fmt.Println("====================================")
	fmt.Printf("🚀 Starting server on %s:%s\n", *host, *port)
	fmt.Printf("📁 Index directory: %s\n", *indexDir)
	fmt.Printf("📂 Data directory: %s\n", *dataDir)

	if *debug {
		fmt.Println("🐛 Debug mode enabled")
	}

	// Create server configuration
	config := &server.Config{
		Host:     *host,
		Port:     *port,
		IndexDir: *indexDir,
		DataDir:  *dataDir,
		Debug:    *debug,
	}

	// Create and start the server
	srv, err := server.New(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("🌐 Server listening on http://%s:%s\n", *host, *port)
		fmt.Println("📋 Available endpoints:")
		fmt.Println("  GET  /health           - Health check")
		fmt.Println("  GET  /api/search       - Search documents")
		fmt.Println("  POST /api/index        - Index a document")
		fmt.Println("  GET  /api/documents    - List documents")
		fmt.Println("  GET  /                 - Web interface")
		fmt.Println("")
		fmt.Println("Press Ctrl+C to stop the server")

		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")
	if err := srv.Shutdown(); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	fmt.Println("✅ Server stopped")
}
