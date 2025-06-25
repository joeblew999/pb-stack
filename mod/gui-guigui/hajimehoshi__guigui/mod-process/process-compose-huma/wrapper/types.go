package wrapper

import "time"

// ProcessState represents the state of a process (matching process-compose types)
type ProcessState struct {
	Name             string  `json:"name" doc:"Process name"`
	Namespace        string  `json:"namespace" doc:"Process namespace"`
	Status           string  `json:"status" doc:"Current status (Running, Completed, Failed, etc.)"`
	SystemTime       string  `json:"system_time" doc:"System time information"`
	Age              int64   `json:"age" doc:"Age in nanoseconds"`
	IsReady          string  `json:"is_ready" doc:"Readiness status"`
	HasReadyProbe    bool    `json:"has_ready_probe" doc:"Whether process has readiness probe"`
	Restarts         int     `json:"restarts" doc:"Number of restarts"`
	ExitCode         int     `json:"exit_code" doc:"Exit code of the process"`
	PID              int     `json:"pid" doc:"Process ID"`
	IsElevated       bool    `json:"is_elevated" doc:"Whether process runs with elevated privileges"`
	PasswordProvided bool    `json:"password_provided" doc:"Whether password was provided"`
	Mem              int64   `json:"mem" doc:"Memory usage in bytes"`
	CPU              float64 `json:"cpu" doc:"CPU usage percentage"`
	IsRunning        bool    `json:"is_running" doc:"Whether process is currently running"`
}

// ProcessesResponse represents the response for getting all processes
type ProcessesResponse struct {
	Data []ProcessState `json:"data" doc:"List of all processes"`
}

// ProjectState represents the overall project state
type ProjectState struct {
	HostName          string `json:"hostName" doc:"Host name"`
	UserName          string `json:"userName" doc:"User name"`
	Version           string `json:"version" doc:"Process-compose version"`
	StartTime         string `json:"startTime" doc:"Project start time"`
	UpTime            int64  `json:"upTime" doc:"Uptime in seconds"`
	ProcessNum        int    `json:"processNum" doc:"Total number of processes"`
	RunningProcessNum int    `json:"runningProcessNum" doc:"Number of running processes"`
}

// SSE Event types for real-time updates
type ProcessStatusEvent struct {
	Type      string       `json:"type" doc:"Event type"`
	Process   ProcessState `json:"process" doc:"Process that changed"`
	Timestamp time.Time    `json:"timestamp" doc:"When the event occurred"`
}

type LogEvent struct {
	ProcessName string    `json:"process_name" doc:"Name of the process"`
	Line        string    `json:"line" doc:"Log line content"`
	Timestamp   time.Time `json:"timestamp" doc:"When the log was generated"`
}

type SystemEvent struct {
	Type      string    `json:"type" doc:"System event type"`
	Message   string    `json:"message" doc:"Event message"`
	Timestamp time.Time `json:"timestamp" doc:"When the event occurred"`
}

// ProcessController interface for controlling process-compose
type ProcessController interface {
	// Lifecycle
	Start() error
	Stop() error
	IsRunning() bool

	// Process Management
	GetProcesses() (*ProcessesResponse, error)
	GetProcess(name string) (*ProcessState, error)
	GetProjectState() (*ProjectState, error)
	StartProcess(name string) error
	StopProcess(name string) error
	RestartProcess(name string) error

	// Event Streams
	ProcessEvents() <-chan ProcessStatusEvent
	LogEvents() <-chan LogEvent
	SystemEvents() <-chan SystemEvent
}
