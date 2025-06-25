package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

// CreateStreams creates the necessary JetStream streams for xtask
func CreateStreams(js nats.JetStreamContext) error {
	// Create events stream for real-time updates
	eventsStream := &nats.StreamConfig{
		Name:     "XTASK_EVENTS",
		Subjects: []string{"xtask.events.>"},
		Storage:  nats.FileStorage,
		MaxAge:   24 * 60 * 60 * 1000000000, // 24 hours in nanoseconds
		MaxMsgs:  10000,
	}

	_, err := js.AddStream(eventsStream)
	if err != nil {
		log.Printf("Failed to create events stream: %v", err)
		return err
	}

	// Create tasks stream for task coordination
	tasksStream := &nats.StreamConfig{
		Name:     "XTASK_TASKS",
		Subjects: []string{"xtask.tasks.>"},
		Storage:  nats.FileStorage,
		MaxAge:   7 * 24 * 60 * 60 * 1000000000, // 7 days in nanoseconds
		MaxMsgs:  50000,
	}

	_, err = js.AddStream(tasksStream)
	if err != nil {
		log.Printf("Failed to create tasks stream: %v", err)
		return err
	}

	// Create nodes stream for node discovery
	nodesStream := &nats.StreamConfig{
		Name:     "XTASK_NODES",
		Subjects: []string{"xtask.nodes.>"},
		Storage:  nats.FileStorage,
		MaxAge:   60 * 60 * 1000000000, // 1 hour in nanoseconds
		MaxMsgs:  1000,
	}

	_, err = js.AddStream(nodesStream)
	if err != nil {
		log.Printf("Failed to create nodes stream: %v", err)
		return err
	}

	log.Println("✅ JetStream streams created successfully")
	return nil
}
