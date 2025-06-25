package wrapper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// ProcessComposeWrapper wraps process-compose-main without modifying it
type ProcessComposeWrapper struct {
	configPath  string
	mainProcess *exec.Cmd
	isRunning   bool
	mutex       sync.RWMutex

	// SSE channels
	processEvents chan ProcessStatusEvent
	logEvents     chan LogEvent
	systemEvents  chan SystemEvent

	// State tracking for change detection
	lastStates  map[string]ProcessState
	statesMutex sync.RWMutex

	// Monitoring
	stopMonitoring chan bool
}

// NewProcessComposeWrapper creates a new wrapper instance
func NewProcessComposeWrapper(configPath string) *ProcessComposeWrapper {
	return &ProcessComposeWrapper{
		configPath:     configPath,
		processEvents:  make(chan ProcessStatusEvent, 100),
		logEvents:      make(chan LogEvent, 100),
		systemEvents:   make(chan SystemEvent, 100),
		lastStates:     make(map[string]ProcessState),
		stopMonitoring: make(chan bool, 1),
	}
}

// NewProcessComposeClient creates a client that connects to existing process-compose API
func NewProcessComposeClient() *ProcessComposeWrapper {
	client := &ProcessComposeWrapper{
		configPath:     "", // Not needed for client-only mode
		processEvents:  make(chan ProcessStatusEvent, 100),
		logEvents:      make(chan LogEvent, 100),
		systemEvents:   make(chan SystemEvent, 100),
		lastStates:     make(map[string]ProcessState),
		stopMonitoring: make(chan bool, 1),
		isRunning:      true, // Mark as running since we're connecting to existing API
	}

	// Start monitoring for changes (but don't manage process lifecycle)
	go client.monitorChanges()

	// Emit system event
	client.systemEvents <- SystemEvent{
		Type:      "client_connected",
		Message:   "Connected to existing Process Compose API",
		Timestamp: time.Now(),
	}

	log.Println("✅ Connected to existing Process Compose API")
	return client
}

// getProcessComposeBinaryPath returns the path to the process-compose binary
func (w *ProcessComposeWrapper) getProcessComposeBinaryPath() string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: Could not get current directory: %v", err)
		return "process-compose" // fallback to PATH
	}

	log.Printf("🔍 Current working directory: %s", cwd)

	// Try multiple possible locations for the binary
	possiblePaths := []string{
		// If running from .bin directory (task run scenario)
		filepath.Join(".", "process-compose"),
		// If running from process-compose-huma directory
		filepath.Join("..", ".bin", "process-compose"),
		// If running from mod-process root directory
		filepath.Join(".bin", "process-compose"),
		// If running from parent directory
		filepath.Join(filepath.Dir(cwd), ".bin", "process-compose"),
	}

	for _, binPath := range possiblePaths {
		log.Printf("🔍 Checking for binary at: %s", binPath)
		if _, err := os.Stat(binPath); err == nil {
			// Convert to absolute path for exec.Command
			absPath, err := filepath.Abs(binPath)
			if err != nil {
				log.Printf("Warning: Could not get absolute path for %s: %v", binPath, err)
				return binPath // fallback to relative path
			}
			log.Printf("✅ Found binary at: %s (absolute: %s)", binPath, absPath)
			return absPath
		}
	}

	// Fallback to PATH
	log.Printf("❌ Binary not found in any expected location, using PATH")
	return "process-compose"
}

// Start launches process-compose-main without TUI
func (w *ProcessComposeWrapper) Start() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.isRunning {
		return fmt.Errorf("process-compose is already running")
	}

	// Start process-compose with API enabled, TUI disabled
	// Use just the filename since config is copied to the same directory
	configFile := "process-compose.yaml"
	w.mainProcess = exec.Command(w.getProcessComposeBinaryPath(),
		"-f", configFile,
		"-t=false",   // Disable TUI
		"-p", "8080", // API port
	)

	// Set up process output
	w.mainProcess.Stdout = os.Stdout
	w.mainProcess.Stderr = os.Stderr

	if err := w.mainProcess.Start(); err != nil {
		return fmt.Errorf("failed to start process-compose: %w", err)
	}

	w.isRunning = true

	// Wait for API to be ready
	if err := w.waitForAPI(); err != nil {
		// Don't call Stop() here as it would cause deadlock
		// Just clean up the process directly
		if w.mainProcess != nil {
			w.mainProcess.Process.Kill()
		}
		w.isRunning = false
		return fmt.Errorf("process-compose API not ready: %w", err)
	}

	// Start monitoring for changes
	go w.monitorChanges()

	// Emit system event
	w.systemEvents <- SystemEvent{
		Type:      "system_started",
		Message:   "Process Compose started successfully",
		Timestamp: time.Now(),
	}

	log.Println("✅ Process Compose started successfully")
	return nil
}

// waitForAPI waits for process-compose API to be ready
func (w *ProcessComposeWrapper) waitForAPI() error {
	for i := 0; i < 30; i++ { // 30 second timeout
		resp, err := http.Get("http://localhost:8080/live")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			log.Printf("🔗 Process Compose API ready after %d seconds", i+1)
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("API not ready after 30 seconds")
}

// monitorChanges polls for state changes and emits SSE events
func (w *ProcessComposeWrapper) monitorChanges() {
	ticker := time.NewTicker(2 * time.Second) // Poll every 2 seconds
	defer ticker.Stop()

	for {
		select {
		case <-w.stopMonitoring:
			return
		case <-ticker.C:
			if !w.IsRunning() {
				return
			}

			// Get current state
			currentStates, err := w.getProcessesFromAPI()
			if err != nil {
				log.Printf("Error getting process states: %v", err)
				continue
			}

			// Compare with last known state and emit events
			w.detectAndEmitChanges(currentStates)

			// Update last known state
			w.statesMutex.Lock()
			w.lastStates = currentStates
			w.statesMutex.Unlock()
		}
	}
}

// getProcessesFromAPI fetches current process states from process-compose API
func (w *ProcessComposeWrapper) getProcessesFromAPI() (map[string]ProcessState, error) {
	resp, err := http.Get("http://localhost:8080/processes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var processesResp ProcessesResponse
	if err := json.NewDecoder(resp.Body).Decode(&processesResp); err != nil {
		return nil, err
	}

	states := make(map[string]ProcessState)
	for _, state := range processesResp.Data {
		states[state.Name] = state
	}

	return states, nil
}

// detectAndEmitChanges compares states and emits SSE events for changes
func (w *ProcessComposeWrapper) detectAndEmitChanges(currentStates map[string]ProcessState) {
	w.statesMutex.RLock()
	lastStates := w.lastStates
	w.statesMutex.RUnlock()

	for name, currentState := range currentStates {
		lastState, existed := lastStates[name]

		// New process
		if !existed {
			w.processEvents <- ProcessStatusEvent{
				Type:      "process_added",
				Process:   currentState,
				Timestamp: time.Now(),
			}
			continue
		}

		// Status change
		if currentState.Status != lastState.Status {
			w.processEvents <- ProcessStatusEvent{
				Type:      "process_status_changed",
				Process:   currentState,
				Timestamp: time.Now(),
			}
		}

		// PID change (indicates restart)
		if currentState.PID != lastState.PID && currentState.PID > 0 {
			w.processEvents <- ProcessStatusEvent{
				Type:      "process_restarted",
				Process:   currentState,
				Timestamp: time.Now(),
			}
		}

		// Restart count change
		if currentState.Restarts != lastState.Restarts {
			w.processEvents <- ProcessStatusEvent{
				Type:      "process_restart_count_changed",
				Process:   currentState,
				Timestamp: time.Now(),
			}
		}
	}

	// Detect removed processes
	for name, lastState := range lastStates {
		if _, exists := currentStates[name]; !exists {
			w.processEvents <- ProcessStatusEvent{
				Type:      "process_removed",
				Process:   lastState,
				Timestamp: time.Now(),
			}
		}
	}
}

// IsRunning returns whether process-compose is currently running
func (w *ProcessComposeWrapper) IsRunning() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.isRunning
}

// Stop shuts down process-compose (or disconnects client)
func (w *ProcessComposeWrapper) Stop() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.isRunning {
		return nil
	}

	// Stop monitoring
	select {
	case w.stopMonitoring <- true:
	default:
	}

	w.isRunning = false

	// Only manage process lifecycle if we started it (not in client mode)
	if w.mainProcess != nil {
		// Try graceful shutdown first
		if err := w.mainProcess.Process.Signal(os.Interrupt); err != nil {
			// Force kill if graceful shutdown fails
			w.mainProcess.Process.Kill()
		}
		w.mainProcess.Wait()

		// Emit system event for process shutdown
		w.systemEvents <- SystemEvent{
			Type:      "system_stopped",
			Message:   "Process Compose stopped",
			Timestamp: time.Now(),
		}
		log.Println("🛑 Process Compose stopped")
	} else {
		// Client mode - just disconnect
		w.systemEvents <- SystemEvent{
			Type:      "client_disconnected",
			Message:   "Disconnected from Process Compose API",
			Timestamp: time.Now(),
		}
		log.Println("🔌 Disconnected from Process Compose API")
	}

	return nil
}

// API Methods - Get data from process-compose API
func (w *ProcessComposeWrapper) GetProcesses() (*ProcessesResponse, error) {
	resp, err := http.Get("http://localhost:8080/processes")
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}
	defer resp.Body.Close()

	var processesResp ProcessesResponse
	if err := json.NewDecoder(resp.Body).Decode(&processesResp); err != nil {
		return nil, fmt.Errorf("failed to decode processes response: %w", err)
	}

	return &processesResp, nil
}

func (w *ProcessComposeWrapper) GetProcess(name string) (*ProcessState, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/process/%s", name))
	if err != nil {
		return nil, fmt.Errorf("failed to get process %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("process %s not found", name)
	}

	var process ProcessState
	if err := json.NewDecoder(resp.Body).Decode(&process); err != nil {
		return nil, fmt.Errorf("failed to decode process response: %w", err)
	}

	return &process, nil
}

func (w *ProcessComposeWrapper) GetProjectState() (*ProjectState, error) {
	resp, err := http.Get("http://localhost:8080/project")
	if err != nil {
		return nil, fmt.Errorf("failed to get project state: %w", err)
	}
	defer resp.Body.Close()

	var project ProjectState
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode project response: %w", err)
	}

	return &project, nil
}

// Control Methods - Use process-compose CLI commands
func (w *ProcessComposeWrapper) StartProcess(name string) error {
	cmd := exec.Command(w.getProcessComposeBinaryPath(),
		"start", name, "-f", "process-compose.yaml")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start process %s: %w", name, err)
	}

	// Emit event (will be confirmed by monitoring)
	w.processEvents <- ProcessStatusEvent{
		Type:      "process_start_requested",
		Process:   ProcessState{Name: name, Status: "Starting"},
		Timestamp: time.Now(),
	}

	return nil
}

func (w *ProcessComposeWrapper) StopProcess(name string) error {
	cmd := exec.Command(w.getProcessComposeBinaryPath(),
		"stop", name, "-f", "process-compose.yaml")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop process %s: %w", name, err)
	}

	// Emit event (will be confirmed by monitoring)
	w.processEvents <- ProcessStatusEvent{
		Type:      "process_stop_requested",
		Process:   ProcessState{Name: name, Status: "Stopping"},
		Timestamp: time.Now(),
	}

	return nil
}

func (w *ProcessComposeWrapper) RestartProcess(name string) error {
	cmd := exec.Command(w.getProcessComposeBinaryPath(),
		"restart", name, "-f", "process-compose.yaml")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart process %s: %w", name, err)
	}

	// Emit event (will be confirmed by monitoring)
	w.processEvents <- ProcessStatusEvent{
		Type:      "process_restart_requested",
		Process:   ProcessState{Name: name, Status: "Restarting"},
		Timestamp: time.Now(),
	}

	return nil
}

// Event channel getters
func (w *ProcessComposeWrapper) ProcessEvents() <-chan ProcessStatusEvent {
	return w.processEvents
}

func (w *ProcessComposeWrapper) LogEvents() <-chan LogEvent {
	return w.logEvents
}

func (w *ProcessComposeWrapper) SystemEvents() <-chan SystemEvent {
	return w.systemEvents
}
