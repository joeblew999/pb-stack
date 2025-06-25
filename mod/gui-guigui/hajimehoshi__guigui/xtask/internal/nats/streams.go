package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

// CreateStreams creates the necessary JetStream streams for xtask coordination
func CreateStreams(js nats.JetStreamContext) error {
	// Create xtask commands stream
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "XTASK_COMMANDS",
		Subjects: []string{"xtask.commands.>"},
		Storage:  nats.FileStorage,
		MaxAge:   24 * 60 * 60 * 1000000000, // 24 hours in nanoseconds
		MaxMsgs:  10000,
	})
	if err != nil && !isStreamExistsError(err) {
		log.Printf("Failed to create XTASK_COMMANDS stream: %v", err)
		return err
	}

	// Create xtask results stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "XTASK_RESULTS",
		Subjects: []string{"xtask.results.>"},
		Storage:  nats.FileStorage,
		MaxAge:   24 * 60 * 60 * 1000000000, // 24 hours in nanoseconds
		MaxMsgs:  10000,
	})
	if err != nil && !isStreamExistsError(err) {
		log.Printf("Failed to create XTASK_RESULTS stream: %v", err)
		return err
	}

	// Create xtask nodes stream for node discovery
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "XTASK_NODES",
		Subjects: []string{"xtask.nodes.>"},
		Storage:  nats.FileStorage,
		MaxAge:   60 * 60 * 1000000000, // 1 hour in nanoseconds
		MaxMsgs:  1000,
	})
	if err != nil && !isStreamExistsError(err) {
		log.Printf("Failed to create XTASK_NODES stream: %v", err)
		return err
	}

	// Create xtask events stream for GUI SSE
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "XTASK_EVENTS",
		Subjects: []string{"xtask.events.>"},
		Storage:  nats.FileStorage,
		MaxAge:   60 * 60 * 1000000000, // 1 hour in nanoseconds
		MaxMsgs:  5000,
	})
	if err != nil && !isStreamExistsError(err) {
		log.Printf("Failed to create XTASK_EVENTS stream: %v", err)
		return err
	}

	log.Println("✅ JetStream streams created successfully")
	return nil
}

// isStreamExistsError checks if the error is due to stream already existing
func isStreamExistsError(err error) bool {
	return err.Error() == "nats: stream name already in use" ||
		err.Error() == "nats: subjects overlap with an existing stream"
}

// PublishCommand publishes a command to the NATS stream
func PublishCommand(js nats.JetStreamContext, nodeID, command string, args []string) error {
	subject := "xtask." + nodeID + "." + command

	// TODO: Implement proper message serialization
	message := command
	for _, arg := range args {
		message += " " + arg
	}

	_, err := js.Publish(subject, []byte(message))
	return err
}

// SubscribeToCommands subscribes to commands for a specific node
func SubscribeToCommands(js nats.JetStreamContext, nodeID string, handler func(msg *nats.Msg)) error {
	subject := "xtask." + nodeID + ".>"

	_, err := js.Subscribe(subject, handler, nats.Durable("xtask-"+nodeID))
	return err
}
