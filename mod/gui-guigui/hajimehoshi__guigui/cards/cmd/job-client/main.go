package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type JobRequest struct {
	ID       string                 `json:"id"`
	Workflow string                 `json:"workflow"`
	Data     map[string]interface{} `json:"data"`
}

type JobResult struct {
	ID       string `json:"id"`
	Success  bool   `json:"success"`
	Output   string `json:"output"`
	Error    string `json:"error,omitempty"`
	Duration string `json:"duration"`
}

func main() {
	var (
		natsURL  = flag.String("nats", "nats://localhost:4222", "NATS server URL")
		workflow = flag.String("workflow", "test-command", "Workflow to execute")
		jobID    = flag.String("id", "", "Job ID (auto-generated if empty)")
	)
	flag.Parse()

	// Connect to NATS
	nc, err := nats.Connect(*natsURL)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("Failed to create JetStream:", err)
	}

	// Generate job ID if not provided
	if *jobID == "" {
		*jobID = fmt.Sprintf("job-%d", time.Now().Unix())
	}

	// Create job request
	job := JobRequest{
		ID:       *jobID,
		Workflow: *workflow,
		Data: map[string]interface{}{
			"name":       "World",
			"deck_id":    "test-deck",
			"output_dir": "/tmp/cards",
			"card_name":  "ace-of-spades",
			"theme":      "classic",
			"style":      "minimal",
		},
	}

	// Subscribe to results
	resultSubject := fmt.Sprintf("jobs.results.%s", *jobID)
	resultSub, err := nc.Subscribe(resultSubject, func(msg *nats.Msg) {
		var result JobResult
		if err := json.Unmarshal(msg.Data, &result); err != nil {
			log.Printf("Failed to unmarshal result: %v", err)
			return
		}

		fmt.Printf("\n=== Job Result ===\n")
		fmt.Printf("ID: %s\n", result.ID)
		fmt.Printf("Success: %v\n", result.Success)
		fmt.Printf("Duration: %s\n", result.Duration)
		if result.Success {
			fmt.Printf("Output: %s\n", result.Output)
		} else {
			fmt.Printf("Error: %s\n", result.Error)
		}
		fmt.Printf("==================\n")
	})
	if err != nil {
		log.Fatal("Failed to subscribe to results:", err)
	}
	defer resultSub.Unsubscribe()

	// Publish job
	jobData, _ := json.Marshal(job)
	subject := fmt.Sprintf("jobs.%s", *workflow)

	ctx := context.Background()
	if _, err := js.Publish(ctx, subject, jobData); err != nil {
		log.Fatal("Failed to publish job:", err)
	}

	fmt.Printf("Published job %s to workflow %s\n", *jobID, *workflow)
	fmt.Printf("Waiting for result on %s...\n", resultSubject)

	// Wait for results
	time.Sleep(10 * time.Second)
	fmt.Println("Done waiting for results")
}
