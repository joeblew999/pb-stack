//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"
)

// ServiceWorkerState manages the WASM service worker state
type ServiceWorkerState struct {
	ServerURL    string                 `json:"serverUrl"`
	IsConnected  bool                   `json:"isConnected"`
	LastUpdate   string                 `json:"lastUpdate"`
	MessageCount int                    `json:"messageCount"`
	Routes       map[string]interface{} `json:"routes"`
}

var state = &ServiceWorkerState{
	ServerURL:    "http://localhost:8081",
	IsConnected:  false,
	LastUpdate:   "",
	MessageCount: 0,
	Routes: map[string]interface{}{
		"/":       "Home - DataStar WASM Service Worker",
		"/hello":  "Hello World endpoint",
		"/status": "Status and health check",
	},
}

func main() {
	fmt.Println("🌐 DataStar WASM Service Worker - Starting...")
	fmt.Println("============================================")

	// Register WASM functions for DataStar integration
	js.Global().Set("connectToServer", js.FuncOf(connectToServer))
	js.Global().Set("disconnectFromServer", js.FuncOf(disconnectFromServer))
	js.Global().Set("handleRoute", js.FuncOf(handleRoute))
	js.Global().Set("getServiceWorkerState", js.FuncOf(getServiceWorkerState))
	js.Global().Set("sendSSERequest", js.FuncOf(sendSSERequest))

	fmt.Println("✅ WASM Service Worker functions registered:")
	fmt.Println("  - connectToServer")
	fmt.Println("  - disconnectFromServer")
	fmt.Println("  - handleRoute")
	fmt.Println("  - getServiceWorkerState")
	fmt.Println("  - sendSSERequest")

	// Initialize the UI
	updateUI()

	fmt.Println("🚀 DataStar WASM Service Worker ready!")
	fmt.Println("📡 Ready to receive SSE from server")

	// Keep the program running
	select {}
}

// connectToServer establishes SSE connection to the server
func connectToServer(this js.Value, args []js.Value) interface{} {
	fmt.Println("� Connecting to server...")

	state.IsConnected = true
	state.LastUpdate = time.Now().Format("15:04:05")

	// Start SSE connection simulation
	go func() {
		// Simulate receiving SSE messages
		for state.IsConnected {
			time.Sleep(2 * time.Second)
			if state.IsConnected {
				state.MessageCount++
				state.LastUpdate = time.Now().Format("15:04:05")
				updateUI()
			}
		}
	}()

	updateUI()
	return nil
}

// disconnectFromServer closes SSE connection
func disconnectFromServer(this js.Value, args []js.Value) interface{} {
	fmt.Println("🔌 Disconnecting from server...")

	state.IsConnected = false
	state.LastUpdate = time.Now().Format("15:04:05")

	updateUI()
	return nil
}

// handleRoute processes different URL routes
func handleRoute(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}

	route := args[0].String()
	fmt.Printf("🛣️ Handling route: %s\n", route)

	var response string
	switch route {
	case "/":
		response = "🏠 Home - DataStar WASM Service Worker is running!"
	case "/hello":
		response = fmt.Sprintf("👋 Hello World! Time: %s", time.Now().Format("15:04:05"))
	case "/status":
		response = fmt.Sprintf("✅ Status: Connected=%t, Messages=%d", state.IsConnected, state.MessageCount)
	default:
		response = fmt.Sprintf("❓ Unknown route: %s", route)
	}

	// Update the response in the UI
	document := js.Global().Get("document")
	target := document.Call("querySelector", "#route-response")
	if !target.IsNull() {
		target.Set("innerHTML", fmt.Sprintf("<div class='route-result'>%s</div>", response))
	}

	return nil
}

// getServiceWorkerState returns current state as JSON
func getServiceWorkerState(this js.Value, args []js.Value) interface{} {
	data, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("Error marshaling state: %v\n", err)
		return "{}"
	}
	return string(data)
}

// sendSSERequest simulates sending SSE request to server
func sendSSERequest(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}

	endpoint := args[0].String()
	fmt.Printf("📡 Sending SSE request to: %s\n", endpoint)

	// Simulate SSE response
	response := fmt.Sprintf("SSE response from %s at %s", endpoint, time.Now().Format("15:04:05"))

	// Update UI with SSE response
	document := js.Global().Get("document")
	target := document.Call("querySelector", "#sse-response")
	if !target.IsNull() {
		target.Set("innerHTML", fmt.Sprintf("<div class='sse-result'>%s</div>", response))
	}

	return nil
}

// updateUI sends current state to DataStar for UI updates
func updateUI() {
	data, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("Error updating UI: %v\n", err)
		return
	}

	// Call DataStar update function via JavaScript
	js.Global().Call("updateDataStarStore", string(data))
}
