//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/hack-pad/safejs"
	"github.com/magodo/go-wasmww"
)

// ControllerState manages the overall application state
type ControllerState struct {
	IsWorkerLoaded   bool                   `json:"isWorkerLoaded"`
	WorkerStatus     string                 `json:"workerStatus"`
	LastMessage      string                 `json:"lastMessage"`
	MessageCount     int                    `json:"messageCount"`
	WorkerState      map[string]interface{} `json:"workerState"`
	AvailableWorkers []string               `json:"availableWorkers"`
}

var (
	controllerState = &ControllerState{
		IsWorkerLoaded:   false,
		WorkerStatus:     "Not loaded",
		LastMessage:      "",
		MessageCount:     0,
		WorkerState:      make(map[string]interface{}),
		AvailableWorkers: []string{"hello-worker.wasm", "todo-worker.wasm", "sse-worker.wasm"},
	}

	currentWorker *wasmww.WasmWebWorkerConn
)

func main() {
	fmt.Println("🎮 go-wasmww Controller - Starting...")

	// Register controller functions for DataStar to call
	js.Global().Set("loadWorker", js.FuncOf(loadWorker))
	js.Global().Set("unloadWorker", js.FuncOf(unloadWorker))
	js.Global().Set("sendToWorker", js.FuncOf(sendToWorker))
	js.Global().Set("pingWorker", js.FuncOf(pingWorker))
	js.Global().Set("incrementWorker", js.FuncOf(incrementWorker))
	js.Global().Set("resetWorker", js.FuncOf(resetWorker))
	js.Global().Set("getControllerState", js.FuncOf(getControllerState))

	fmt.Println("✅ Controller functions registered:")
	fmt.Println("  - loadWorker(workerPath)")
	fmt.Println("  - unloadWorker()")
	fmt.Println("  - sendToWorker(message)")
	fmt.Println("  - pingWorker()")
	fmt.Println("  - incrementWorker()")
	fmt.Println("  - resetWorker()")
	fmt.Println("  - getControllerState()")

	// Initialize DataStar with controller state
	updateDataStarStore()

	fmt.Println("🚀 Controller ready for DataStar integration!")

	// Keep the program running
	select {}
}

// loadWorker dynamically loads a WASM worker
func loadWorker(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		fmt.Println("❌ loadWorker requires worker path argument")
		return nil
	}

	workerPath := args[0].String()
	fmt.Printf("📦 Loading worker: %s\n", workerPath)

	// Unload existing worker if any
	if currentWorker != nil {
		fmt.Println("🔄 Unloading existing worker...")
		currentWorker.Terminate()
		currentWorker = nil
	}

	// Create new worker connection
	conn := &wasmww.WasmWebWorkerConn{
		Name: "hello-world-worker",
		Path: workerPath,
		Args: []string{"hello-worker", "--mode=production"},
		Env:  []string{"WORKER_ENV=production", "DEBUG=true"},
	}

	// Start the worker
	if err := conn.Start(); err != nil {
		fmt.Printf("❌ Failed to start worker: %v\n", err)
		controllerState.WorkerStatus = fmt.Sprintf("Failed to load: %v", err)
		updateDataStarStore()
		return nil
	}

	currentWorker = conn
	controllerState.IsWorkerLoaded = true
	controllerState.WorkerStatus = "Loaded and running"
	controllerState.MessageCount = 0

	fmt.Println("✅ Worker loaded successfully")

	// Start handling worker events
	go handleWorkerEvents()

	updateDataStarStore()
	return nil
}

// unloadWorker terminates the current worker
func unloadWorker(this js.Value, args []js.Value) interface{} {
	if currentWorker == nil {
		fmt.Println("⚠️ No worker to unload")
		return nil
	}

	fmt.Println("🛑 Unloading worker...")
	currentWorker.Terminate()
	currentWorker = nil

	controllerState.IsWorkerLoaded = false
	controllerState.WorkerStatus = "Unloaded"
	controllerState.WorkerState = make(map[string]interface{})

	updateDataStarStore()
	return nil
}

// sendToWorker sends a message to the current worker
func sendToWorker(this js.Value, args []js.Value) interface{} {
	if currentWorker == nil {
		fmt.Println("❌ No worker loaded")
		return nil
	}

	if len(args) < 1 {
		fmt.Println("❌ sendToWorker requires message argument")
		return nil
	}

	message := args[0].String()
	fmt.Printf("📤 Sending to worker: %s\n", message)

	if err := currentWorker.PostMessage(safejs.Safe(js.ValueOf(message)), nil); err != nil {
		fmt.Printf("❌ Failed to send message: %v\n", err)
		return nil
	}

	controllerState.MessageCount++
	updateDataStarStore()
	return nil
}

// pingWorker sends a ping to the worker
func pingWorker(this js.Value, args []js.Value) interface{} {
	return sendToWorker(this, []js.Value{js.ValueOf("ping")})
}

// incrementWorker tells worker to increment its counter
func incrementWorker(this js.Value, args []js.Value) interface{} {
	return sendToWorker(this, []js.Value{js.ValueOf("increment")})
}

// resetWorker tells worker to reset its counter
func resetWorker(this js.Value, args []js.Value) interface{} {
	return sendToWorker(this, []js.Value{js.ValueOf("reset")})
}

// getControllerState returns current controller state as JSON
func getControllerState(this js.Value, args []js.Value) interface{} {
	data, err := json.Marshal(controllerState)
	if err != nil {
		fmt.Printf("❌ Error marshaling controller state: %v\n", err)
		return "{}"
	}
	return string(data)
}

// handleWorkerEvents processes events from the worker
func handleWorkerEvents() {
	if currentWorker == nil {
		return
	}

	fmt.Println("👂 Starting to listen for worker events...")

	for event := range currentWorker.EventChannel() {
		data, err := event.Data()
		if err != nil {
			fmt.Printf("❌ Error getting event data: %v\n", err)
			continue
		}

		str, err := data.String()
		if err != nil {
			fmt.Printf("❌ Error converting event data to string: %v\n", err)
			continue
		}
		fmt.Printf("📥 Received from worker: %s\n", str)

		// Handle different message types
		if strings.HasPrefix(str, "STATE_UPDATE:") {
			// Parse worker state update
			stateJSON := str[len("STATE_UPDATE:"):]
			var workerState map[string]interface{}
			if err := json.Unmarshal([]byte(stateJSON), &workerState); err != nil {
				fmt.Printf("❌ Error parsing worker state: %v\n", err)
			} else {
				controllerState.WorkerState = workerState
			}
		} else {
			// Regular message
			controllerState.LastMessage = str
		}

		controllerState.MessageCount++
		updateDataStarStore()
	}

	fmt.Println("👋 Worker event handler finished")
	controllerState.IsWorkerLoaded = false
	controllerState.WorkerStatus = "Disconnected"
	updateDataStarStore()
}

// updateDataStarStore sends current state to DataStar
func updateDataStarStore() {
	data, err := json.Marshal(controllerState)
	if err != nil {
		fmt.Printf("❌ Error marshaling controller state: %v\n", err)
		return
	}

	// Call DataStar update function via JavaScript
	js.Global().Call("updateDataStarStore", string(data))
}
