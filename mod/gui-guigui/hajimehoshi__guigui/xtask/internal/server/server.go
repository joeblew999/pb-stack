package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"xtask/internal/api"
	"xtask/internal/nats"

	"github.com/nats-io/nats-server/v2/server"
	natsgo "github.com/nats-io/nats.go"
)

type Config struct {
	Port     int
	NATSPort int
	DataDir  string
	Verbose  bool
}

type Server struct {
	config     *Config
	natsServer *server.Server
	natsConn   *natsgo.Conn
	httpServer *http.Server
	apiHandler *api.Handler
}

func New(config *Config) (*Server, error) {
	if config == nil {
		config = &Config{
			Port:     8080,
			NATSPort: 4222,
			DataDir:  "./.data",
			Verbose:  false,
		}
	}

	return &Server{
		config: config,
	}, nil
}

func (s *Server) Start() error {
	// Create data directory
	if err := os.MkdirAll(s.config.DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Start embedded NATS server
	if err := s.startNATS(); err != nil {
		return fmt.Errorf("failed to start NATS: %w", err)
	}

	// Connect to NATS
	if err := s.connectNATS(); err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Start HTTP API server
	if err := s.startHTTP(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	var lastErr error

	// Shutdown HTTP server
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			lastErr = err
			log.Printf("Error shutting down HTTP server: %v", err)
		}
	}

	// Close NATS connection
	if s.natsConn != nil {
		s.natsConn.Close()
	}

	// Shutdown NATS server
	if s.natsServer != nil {
		s.natsServer.Shutdown()
		s.natsServer.WaitForShutdown()
	}

	return lastErr
}

func (s *Server) startNATS() error {
	// Ensure absolute path for NATS data directory
	absDataDir, err := filepath.Abs(s.config.DataDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for data directory: %w", err)
	}

	// Use organized data structure: .data/nats/jetstream
	natsDataDir := filepath.Join(absDataDir, "nats")
	jetStreamDir := filepath.Join(natsDataDir, "jetstream")
	if err := os.MkdirAll(jetStreamDir, 0755); err != nil {
		return fmt.Errorf("failed to create NATS JetStream directory: %w", err)
	}

	opts := &server.Options{
		Host:      "127.0.0.1",
		Port:      s.config.NATSPort,
		HTTPHost:  "127.0.0.1",
		HTTPPort:  8222,
		JetStream: true,
		StoreDir:  jetStreamDir,
		Debug:     s.config.Verbose,
		Trace:     s.config.Verbose,
		// Force JetStream to use local directory
		JetStreamMaxMemory: 64 * 1024 * 1024,   // 64MB
		JetStreamMaxStore:  1024 * 1024 * 1024, // 1GB
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return fmt.Errorf("failed to create NATS server: %w", err)
	}

	s.natsServer = ns

	// Start server in background
	go ns.Start()

	// Wait for server to be ready
	if !ns.ReadyForConnections(10 * time.Second) {
		return fmt.Errorf("NATS server failed to start within timeout")
	}

	log.Printf("📡 NATS server started on port %d", s.config.NATSPort)
	log.Printf("📊 NATS monitoring on http://localhost:8222")

	return nil
}

func (s *Server) connectNATS() error {
	natsURL := fmt.Sprintf("nats://localhost:%d", s.config.NATSPort)

	// Wait for NATS to be ready
	var nc *natsgo.Conn
	var err error

	for i := 0; i < 10; i++ {
		nc, err = natsgo.Connect(natsURL)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	s.natsConn = nc
	log.Printf("✅ Connected to NATS at %s", natsURL)

	// Initialize JetStream
	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("failed to get JetStream context: %w", err)
	}

	// Create streams
	if err := nats.CreateStreams(js); err != nil {
		return fmt.Errorf("failed to create JetStream streams: %w", err)
	}

	return nil
}

func (s *Server) startHTTP() error {
	// Create API handler
	s.apiHandler = api.New(s.natsConn)

	// Setup Chi router
	router := s.apiHandler.SetupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: router,
	}

	// Start server in background
	go func() {
		log.Printf("🌐 HTTP server starting on port %d", s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return nil
}
