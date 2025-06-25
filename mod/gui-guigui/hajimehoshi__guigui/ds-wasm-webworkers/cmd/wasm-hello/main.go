//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syscall/js"
	"time"

	"github.com/hack-pad/safejs"
	"github.com/magodo/go-wasmww"
)

// HelloWorldState represents the state managed by this worker
type HelloWorldState struct {
	Message     string    `json:"message"`
	Counter     int       `json:"counter"`
	LastUpdate  string    `json:"lastUpdate"`
	WorkerName  string    `json:"workerName"`
	Environment []string  `json:"environment"`
	Arguments   []string  `json:"arguments"`
}

var state = &HelloWorldState{
	Message:    "Hello from go-wasmww Worker!",
	Counter:    0,
	LastUpdate: "",
}

func main() {
	fmt.Println("🌟 Hello World WASM Worker - Starting with go-wasmww...")

	// Create self connection using go-wasmww
	self, err := wasmww.NewSelfConn()
	if err != nil {
		log.Fatalf("Failed to create self connection: %v", err)
	}

	// Get worker name
	name, err := self.Name()
	if err != nil {
		log.Fatalf("Failed to get worker name: %v", err)
	}

	// Initialize state with worker info
	state.WorkerName = name
	state.Arguments = os.Args
	state.Environment = os.Environ()
	state.LastUpdate = time.Now().Format("15:04:05")

	fmt.Printf("✅ Worker '%s' initialized\n", name)
	fmt.Printf("📋 Args: %v\n", os.Args)
	fmt.Printf("🌍 Env: %v\n", os.Environ())

	// Setup connection and get event channel
	ch, err := self.SetupConn()
	if err != nil {
		log.Fatalf("Failed to setup connection: %v", err)
	}

	fmt.Println("🔗 Connection established, ready to receive messages")

	// Send initial state to controller
	sendStateUpdate(self)

	// Start periodic counter updates
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				state.Counter++
				state.LastUpdate = time.Now().Format("15:04:05")
				sendStateUpdate(self)
			}
		}
	}()

	// Handle incoming messages
	for event := range ch {
		handleMessage(self, event)
	}

	fmt.Println("👋 Worker shutting down...")
}

// handleMessage processes incoming messages from the controller
func handleMessage(self *wasmww.SelfConn, event interface{}) {
	// Extract message data using type assertion
	// The event should be from go-webworkers types.MessageEventMessage
	if msgEvent, ok := event.(interface{ Data() (js.Value, error) }); ok {
		data, err := msgEvent.Data()
		if err != nil {
			fmt.Printf("❌ Error extracting event data: %v\n", err)
			return
		}

		str := data.String()
		fmt.Printf("📨 Received message: %s\n", str)

		// Handle different message types
		switch str {
		case "ping":
			// Respond to ping
			response := fmt.Sprintf("pong from %s at %s", state.WorkerName, time.Now().Format("15:04:05"))
			sendMessage(self, response)

		case "increment":
			// Increment counter
			state.Counter++
			state.LastUpdate = time.Now().Format("15:04:05")
			sendStateUpdate(self)

		case "reset":
			// Reset counter
			state.Counter = 0
			state.LastUpdate = time.Now().Format("15:04:05")
			sendStateUpdate(self)

		case "status":
			// Send current status
			sendStateUpdate(self)

		case "close":
			// Close worker
			fmt.Println("🛑 Received close command")
			if err := self.Close(); err != nil {
				fmt.Printf("❌ Error closing worker: %v\n", err)
			}

		default:
			// Echo back unknown messages
			response := fmt.Sprintf("Echo: %s (from %s)", str, state.WorkerName)
			sendMessage(self, response)
		}
	} else {
		fmt.Printf("❌ Unknown event type: %T\n", event)
	}
}

// sendStateUpdate sends the current state to the controller
func sendStateUpdate(self *wasmww.SelfConn) {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("❌ Error marshaling state: %v\n", err)
		return
	}

	// Send as JSON with special prefix to identify as state update
	message := fmt.Sprintf("STATE_UPDATE:%s", string(stateJSON))
	sendMessage(self, message)
}

// sendMessage sends a message to the controller
func sendMessage(self *wasmww.SelfConn, message string) {
	if err := self.PostMessage(safejs.Safe(js.ValueOf(message)), nil); err != nil {
		fmt.Printf("❌ Error sending message: %v\n", err)
	}
}
