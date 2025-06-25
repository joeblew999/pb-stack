package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// WorkflowConfig defines a configurable workflow step
type WorkflowConfig struct {
	Name       string            `json:"name"`
	Binary     string            `json:"binary"`
	Args       []string          `json:"args"`
	Env        map[string]string `json:"env"`
	WorkingDir string            `json:"working_dir"`
	Timeout    string            `json:"timeout"`
	OnSuccess  *NextStep         `json:"on_success,omitempty"`
	OnError    *NextStep         `json:"on_error,omitempty"`
}

type NextStep struct {
	Subject string      `json:"subject"`
	Data    interface{} `json:"data"`
}

// JobRequest represents an incoming job
type JobRequest struct {
	ID       string                 `json:"id"`
	Workflow string                 `json:"workflow"`
	Data     map[string]interface{} `json:"data"`
}

// JobResult represents job completion
type JobResult struct {
	ID       string `json:"id"`
	Success  bool   `json:"success"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	Duration string `json:"duration"`
}

// Orchestrator manages workflow execution
type Orchestrator struct {
	nc        *nats.Conn
	js        jetstream.JetStream
	workflows map[string]WorkflowConfig
	binDir    string
}

func main() {
	var (
		natsURL    = flag.String("nats", "nats://localhost:4222", "NATS server URL")
		configFile = flag.String("config", ".config/orchestrator.yaml", "Workflow config file")
		binDir     = flag.String("bin-dir", ".bin", "Binary directory")
		subject    = flag.String("subject", "jobs.>", "NATS subject pattern to listen on")
	)
	flag.Parse()

	orchestrator, err := NewOrchestrator(*natsURL, *configFile, *binDir)
	if err != nil {
		log.Fatal("Failed to create orchestrator:", err)
	}
	defer orchestrator.Close()

	log.Printf("Starting orchestrator, listening on %s", *subject)

	// Subscribe to job requests
	_, err = orchestrator.js.Subscribe(*subject, orchestrator.handleJob, jetstream.Durable("orchestrator"))
	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Shutting down orchestrator...")
}

func NewOrchestrator(natsURL, configFile, binDir string) (*Orchestrator, error) {
	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream: %w", err)
	}

	// Load workflow configurations
	workflows, err := loadWorkflows(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load workflows: %w", err)
	}

	return &Orchestrator{
		nc:        nc,
		js:        js,
		workflows: workflows,
		binDir:    binDir,
	}, nil
}

func (o *Orchestrator) Close() {
	if o.nc != nil {
		o.nc.Close()
	}
}

func (o *Orchestrator) handleJob(msg jetstream.Msg) {
	var job JobRequest
	if err := json.Unmarshal(msg.Data(), &job); err != nil {
		log.Printf("Failed to unmarshal job: %v", err)
		msg.Nak()
		return
	}

	log.Printf("Processing job %s with workflow %s", job.ID, job.Workflow)

	// Get workflow config
	workflow, exists := o.workflows[job.Workflow]
	if !exists {
		log.Printf("Unknown workflow: %s", job.Workflow)
		msg.Nak()
		return
	}

	// Execute workflow
	result := o.executeWorkflow(job, workflow)

	// Publish result
	resultData, _ := json.Marshal(result)
	o.nc.Publish(fmt.Sprintf("jobs.results.%s", job.ID), resultData)

	// Handle next steps
	if result.Success && workflow.OnSuccess != nil {
		o.publishNextStep(job, workflow.OnSuccess)
	} else if !result.Success && workflow.OnError != nil {
		o.publishNextStep(job, workflow.OnError)
	}

	msg.Ack()
}

func (o *Orchestrator) executeWorkflow(job JobRequest, workflow WorkflowConfig) JobResult {
	start := time.Now()
	result := JobResult{
		ID:      job.ID,
		Success: false,
	}

	// Build command
	binaryPath := filepath.Join(o.binDir, workflow.Binary)

	// Substitute variables in args
	args := make([]string, len(workflow.Args))
	for i, arg := range workflow.Args {
		args[i] = substituteVariables(arg, job.Data)
	}

	cmd := exec.Command(binaryPath, args...)

	// Set environment
	cmd.Env = os.Environ()
	for key, value := range workflow.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, substituteVariables(value, job.Data)))
	}

	// Set working directory
	if workflow.WorkingDir != "" {
		cmd.Dir = workflow.WorkingDir
	}

	// Parse timeout
	timeout := 30 * time.Second
	if workflow.Timeout != "" {
		if t, err := time.ParseDuration(workflow.Timeout); err == nil {
			timeout = t
		}
	}

	// Execute with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	output, err := cmd.Output()
	result.Duration = time.Since(start).String()

	if err != nil {
		result.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.Error = fmt.Sprintf("%s: %s", err.Error(), string(exitErr.Stderr))
		}
		log.Printf("Job %s failed: %s", job.ID, result.Error)
	} else {
		result.Success = true
		result.Output = string(output)
		log.Printf("Job %s completed successfully in %s", job.ID, result.Duration)
	}

	return result
}

func (o *Orchestrator) publishNextStep(job JobRequest, next *NextStep) {
	data := next.Data
	if data == nil {
		data = job.Data
	}

	nextData, _ := json.Marshal(data)
	subject := substituteVariables(next.Subject, job.Data)

	if err := o.nc.Publish(subject, nextData); err != nil {
		log.Printf("Failed to publish next step to %s: %v", subject, err)
	} else {
		log.Printf("Published next step to %s", subject)
	}
}

func loadWorkflows(configFile string) (map[string]WorkflowConfig, error) {
	// For now, return hardcoded workflows
	// TODO: Load from YAML file
	return map[string]WorkflowConfig{
		"generate-cards": {
			Name:    "Generate Card Images",
			Binary:  "cardgen",
			Args:    []string{"-deck", "${deck_id}", "-output", "${output_dir}"},
			Timeout: "30s",
			OnSuccess: &NextStep{
				Subject: "jobs.svg-deck",
				Data: map[string]interface{}{
					"deck_id":    "${deck_id}",
					"card_dir":   "${output_dir}",
					"output_svg": "${output_dir}/deck.svg",
				},
			},
		},
		"svg-deck": {
			Name:    "Generate SVG Deck",
			Binary:  "svgdeck",
			Args:    []string{"-input", "${card_dir}", "-output", "${output_svg}"},
			Timeout: "15s",
			OnSuccess: &NextStep{
				Subject: "jobs.png-deck",
				Data: map[string]interface{}{
					"input_svg":  "${output_svg}",
					"output_png": "${output_dir}/deck.png",
				},
			},
		},
		"png-deck": {
			Name:    "Generate PNG Deck",
			Binary:  "pngdeck",
			Args:    []string{"-input", "${input_svg}", "-output", "${output_png}"},
			Timeout: "20s",
		},
	}, nil
}

func substituteVariables(template string, data map[string]interface{}) string {
	// Simple variable substitution - replace ${var} with data["var"]
	// TODO: Use a proper template engine
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("${%s}", key)
		if str, ok := value.(string); ok {
			result = fmt.Sprintf("%s", str) // This is simplified
		}
	}
	return result
}
