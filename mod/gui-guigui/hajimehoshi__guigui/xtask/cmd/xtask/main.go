package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"xtask/internal/client"
	"xtask/internal/server"
	"xtask/internal/tools"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "xtask",
		Short: "Cross-platform development task runner with embedded NATS coordination",
		Long: `xtask is a unified Go binary that embeds Task + all cross-platform development tools.
It can run as both a CLI tool and a server with embedded NATS JetStream for team coordination.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		Run:     runCLI,
	}

	// Global flags
	rootCmd.PersistentFlags().Bool("server", false, "Run as server with embedded NATS JetStream")
	rootCmd.PersistentFlags().String("server-url", "", "Connect to specific server URL")
	rootCmd.PersistentFlags().String("cluster-url", "", "Connect to NATS cluster")
	rootCmd.PersistentFlags().Bool("local", false, "Force local execution")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging")

	// Add tool subcommands
	rootCmd.AddCommand(createToolCommands()...)

	// Add server command
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Start xtask server with embedded NATS JetStream",
		Run:   runServer,
	}
	serverCmd.Flags().Int("port", 8080, "HTTP API server port")
	serverCmd.Flags().Int("nats-port", 4222, "NATS server port")
	serverCmd.Flags().String("data-dir", "./.data", "Data directory for JetStream")
	rootCmd.AddCommand(serverCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runCLI(cmd *cobra.Command, args []string) {
	// Check if server flag is set
	if server, _ := cmd.Flags().GetBool("server"); server {
		runServer(cmd, args)
		return
	}

	// Smart execution routing
	if local, _ := cmd.Flags().GetBool("local"); !local {
		if client := tryConnectToServer(cmd); client != nil {
			if err := client.ExecuteCommand(args); err != nil {
				log.Printf("Server execution failed: %v, falling back to local", err)
			} else {
				return
			}
		}
	}

	// Fallback to local execution
	executeLocally(args)
}

func runServer(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt("port")
	natsPort, _ := cmd.Flags().GetInt("nats-port")
	dataDir, _ := cmd.Flags().GetString("data-dir")
	verbose, _ := cmd.Flags().GetBool("verbose")

	log.Println("🚀 Starting xtask server...")

	// Create server instance
	s, err := server.New(&server.Config{
		Port:     port,
		NATSPort: natsPort,
		DataDir:  dataDir,
		Verbose:  verbose,
	})
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Printf("✅ xtask server ready")
	log.Printf("🌐 API: http://localhost:%d", port)
	log.Printf("📡 NATS: nats://localhost:%d", natsPort)
	log.Printf("🌍 Web UI: http://localhost:%d/web", port)

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Shutting down xtask server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	log.Println("✅ xtask server stopped")
}

func tryConnectToServer(cmd *cobra.Command) *client.Client {
	// Try specific server URL first
	if serverURL, _ := cmd.Flags().GetString("server-url"); serverURL != "" {
		if client := client.New(serverURL); client.IsAvailable() {
			return client
		}
	}

	// Try local server
	if client := client.New("http://localhost:8080"); client.IsAvailable() {
		return client
	}

	// Try cluster coordination
	if clusterURL, _ := cmd.Flags().GetString("cluster-url"); clusterURL != "" {
		// TODO: Implement cluster client
		log.Printf("Cluster coordination not yet implemented: %s", clusterURL)
	}

	return nil
}

func executeLocally(args []string) {
	if len(args) == 0 {
		// Default to task execution
		tools.ExecuteTask([]string{})
		return
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "which":
		tools.ExecuteWhich(commandArgs)
	case "got":
		tools.ExecuteGot(commandArgs)
	case "silent":
		tools.ExecuteSilent(commandArgs)
	case "kill-port":
		tools.ExecuteKillPort(commandArgs)
	case "wait-for-port":
		tools.ExecuteWaitForPort(commandArgs)
	case "tree":
		tools.ExecuteTree(commandArgs)
	case "health-check":
		tools.ExecuteHealthCheck(commandArgs)
	default:
		// Default to task execution
		tools.ExecuteTask(args)
	}
}

func createToolCommands() []*cobra.Command {
	var commands []*cobra.Command

	// which command
	whichCmd := &cobra.Command{
		Use:   "which [binary]",
		Short: "Find binary location (cross-platform)",
		Run: func(cmd *cobra.Command, args []string) {
			tools.ExecuteWhich(args)
		},
	}
	commands = append(commands, whichCmd)

	// got command
	gotCmd := &cobra.Command{
		Use:   "got [url] [options]",
		Short: "Download files (cross-platform)",
		Run: func(cmd *cobra.Command, args []string) {
			tools.ExecuteGot(args)
		},
	}
	commands = append(commands, gotCmd)

	// silent command
	silentCmd := &cobra.Command{
		Use:   "silent [command] [args...]",
		Short: "Execute command silently (cross-platform 2>/dev/null)",
		Run: func(cmd *cobra.Command, args []string) {
			tools.ExecuteSilent(args)
		},
	}
	commands = append(commands, silentCmd)

	// kill-port command
	killPortCmd := &cobra.Command{
		Use:   "kill-port [port]",
		Short: "Kill process on port (cross-platform)",
		Run: func(cmd *cobra.Command, args []string) {
			tools.ExecuteKillPort(args)
		},
	}
	commands = append(commands, killPortCmd)

	// wait-for-port command
	waitPortCmd := &cobra.Command{
		Use:   "wait-for-port [port] [timeout]",
		Short: "Wait for port to be available",
		Run: func(cmd *cobra.Command, args []string) {
			tools.ExecuteWaitForPort(args)
		},
	}
	commands = append(commands, waitPortCmd)

	// tree command
	treeCmd := &cobra.Command{
		Use:   "tree [path]",
		Short: "Display directory tree (cross-platform)",
		Run: func(cmd *cobra.Command, args []string) {
			tools.ExecuteTree(args)
		},
	}
	commands = append(commands, treeCmd)

	// health-check command
	healthCmd := &cobra.Command{
		Use:   "health-check [url]",
		Short: "HTTP health check",
		Run: func(cmd *cobra.Command, args []string) {
			tools.ExecuteHealthCheck(args)
		},
	}
	commands = append(commands, healthCmd)

	return commands
}
