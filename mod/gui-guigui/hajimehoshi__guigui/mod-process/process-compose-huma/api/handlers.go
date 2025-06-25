package api

import (
	"context"
	"fmt"
	"time"

	"process-compose-huma/wrapper"

	"github.com/danielgtaylor/huma/v2"
)

// Input/Output types for Huma API
type GetProcessInput struct {
	Name string `path:"name" doc:"Process name"`
}

type GetProcessOutput struct {
	Body wrapper.ProcessState `json:"process"`
}

type GetProcessesOutput struct {
	Body wrapper.ProcessesResponse `json:"processes"`
}

type GetProjectOutput struct {
	Body wrapper.ProjectState `json:"project"`
}

type StartProcessInput struct {
	Name string `path:"name" doc:"Process name"`
}

type StopProcessInput struct {
	Name string `path:"name" doc:"Process name"`
}

type RestartProcessInput struct {
	Name string `path:"name" doc:"Process name"`
}

type ProcessActionOutput struct {
	Body struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Process string `json:"process"`
	} `json:"result"`
}

type HealthOutput struct {
	Body struct {
		Status    string    `json:"status"`
		Timestamp time.Time `json:"timestamp"`
		Version   string    `json:"version"`
	} `json:"health"`
}

// APIHandlers contains all the Huma API handlers
type APIHandlers struct {
	controller wrapper.ProcessController
}

// NewAPIHandlers creates a new API handlers instance
func NewAPIHandlers(controller wrapper.ProcessController) *APIHandlers {
	return &APIHandlers{
		controller: controller,
	}
}

// GetProcesses returns all processes
func (h *APIHandlers) GetProcesses(ctx context.Context, input *struct{}) (*GetProcessesOutput, error) {
	processes, err := h.controller.GetProcesses()
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get processes", err)
	}

	return &GetProcessesOutput{Body: *processes}, nil
}

// GetProcess returns a specific process by name
func (h *APIHandlers) GetProcess(ctx context.Context, input *GetProcessInput) (*GetProcessOutput, error) {
	process, err := h.controller.GetProcess(input.Name)
	if err != nil {
		return nil, huma.Error404NotFound("Process not found", err)
	}

	return &GetProcessOutput{Body: *process}, nil
}

// GetProject returns the project state
func (h *APIHandlers) GetProject(ctx context.Context, input *struct{}) (*GetProjectOutput, error) {
	project, err := h.controller.GetProjectState()
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get project state", err)
	}

	return &GetProjectOutput{Body: *project}, nil
}

// StartProcess starts a specific process
func (h *APIHandlers) StartProcess(ctx context.Context, input *StartProcessInput) (*ProcessActionOutput, error) {
	err := h.controller.StartProcess(input.Name)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to start process", err)
	}

	return &ProcessActionOutput{
		Body: struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
			Process string `json:"process"`
		}{
			Success: true,
			Message: fmt.Sprintf("Process %s start requested", input.Name),
			Process: input.Name,
		},
	}, nil
}

// StopProcess stops a specific process
func (h *APIHandlers) StopProcess(ctx context.Context, input *StopProcessInput) (*ProcessActionOutput, error) {
	err := h.controller.StopProcess(input.Name)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to stop process", err)
	}

	return &ProcessActionOutput{
		Body: struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
			Process string `json:"process"`
		}{
			Success: true,
			Message: fmt.Sprintf("Process %s stop requested", input.Name),
			Process: input.Name,
		},
	}, nil
}

// RestartProcess restarts a specific process
func (h *APIHandlers) RestartProcess(ctx context.Context, input *RestartProcessInput) (*ProcessActionOutput, error) {
	err := h.controller.RestartProcess(input.Name)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to restart process", err)
	}

	return &ProcessActionOutput{
		Body: struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
			Process string `json:"process"`
		}{
			Success: true,
			Message: fmt.Sprintf("Process %s restart requested", input.Name),
			Process: input.Name,
		},
	}, nil
}

// HealthCheck returns the health status
func (h *APIHandlers) HealthCheck(ctx context.Context, input *struct{}) (*HealthOutput, error) {
	return &HealthOutput{
		Body: struct {
			Status    string    `json:"status"`
			Timestamp time.Time `json:"timestamp"`
			Version   string    `json:"version"`
		}{
			Status:    "healthy",
			Timestamp: time.Now(),
			Version:   "1.0.0-huma",
		},
	}, nil
}

// RegisterRoutes registers all API routes with Huma
func (h *APIHandlers) RegisterRoutes(api huma.API) {
	// Process management endpoints
	huma.Register(api, huma.Operation{
		OperationID: "getProcesses",
		Method:      "GET",
		Path:        "/processes",
		Summary:     "Get all processes",
		Description: "Retrieve the current state of all managed processes",
		Tags:        []string{"Processes"},
	}, h.GetProcesses)

	huma.Register(api, huma.Operation{
		OperationID: "getProcess",
		Method:      "GET",
		Path:        "/process/{name}",
		Summary:     "Get specific process",
		Description: "Retrieve the current state of a specific process by name",
		Tags:        []string{"Processes"},
	}, h.GetProcess)

	huma.Register(api, huma.Operation{
		OperationID: "getProject",
		Method:      "GET",
		Path:        "/project",
		Summary:     "Get project state",
		Description: "Retrieve the overall project state and statistics",
		Tags:        []string{"Project"},
	}, h.GetProject)

	// Process control endpoints
	huma.Register(api, huma.Operation{
		OperationID: "startProcess",
		Method:      "POST",
		Path:        "/process/{name}/start",
		Summary:     "Start process",
		Description: "Start a specific process by name",
		Tags:        []string{"Process Control"},
	}, h.StartProcess)

	huma.Register(api, huma.Operation{
		OperationID: "stopProcess",
		Method:      "POST",
		Path:        "/process/{name}/stop",
		Summary:     "Stop process",
		Description: "Stop a specific process by name",
		Tags:        []string{"Process Control"},
	}, h.StopProcess)

	huma.Register(api, huma.Operation{
		OperationID: "restartProcess",
		Method:      "POST",
		Path:        "/process/{name}/restart",
		Summary:     "Restart process",
		Description: "Restart a specific process by name",
		Tags:        []string{"Process Control"},
	}, h.RestartProcess)

	// Health check endpoint
	huma.Register(api, huma.Operation{
		OperationID: "healthCheck",
		Method:      "GET",
		Path:        "/health",
		Summary:     "Health check",
		Description: "Check if the API is healthy and responsive",
		Tags:        []string{"Health"},
	}, h.HealthCheck)
}
