package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

const (
	defaultNatsURL            = "nats://localhost:4222"
	defaultPocketBaseServeURL = "http://localhost:8090"
	defaultPocketBaseAdmin    = "admin" // Informational, as PB handles admin creation
	defaultPocketBaseDataPath = "./data"
	jetStreamName             = "POCKETBASE_EVENTS"
)

// Helper to get environment variables or return a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Environment variable %s not set, using default: %s", key, fallback)
	return fallback
}

func main() {
	// --- Configuration from Environment Variables ---
	natsURL := getEnv("NATS_SERVER_URL", defaultNatsURL)
	pbServeURL := getEnv("POCKETBASE_SERVER_URL", defaultPocketBaseServeURL)
	pbAdminUser := getEnv("POCKETBASE_ADMIN_USER", defaultPocketBaseAdmin)
	pbDataPath := getEnv("POCKETBASE_DATA_PATH", defaultPocketBaseDataPath)

	log.Printf("--- Configuration ---")
	log.Printf("NATS Server URL: %s", natsURL)
	log.Printf("PocketBase Serve URL: %s", pbServeURL)
	log.Printf("PocketBase Admin User (info): %s", pbAdminUser)
	log.Printf("PocketBase Data Path: %s", pbDataPath)
	log.Printf("--------------------")

	// --- Initialize PocketBase ---
	// PocketBase typically reads its own env vars (PB_ENCRYPTION_KEY, etc.)
	// or uses command-line flags. We'll set flags programmatically.
	// Ensure "serve" is the first argument if not already present for flag parsing.
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "serve")
	}

	// Remove http:// or https:// prefix for PocketBase's --http flag
	pbHostPort := strings.TrimPrefix(pbServeURL, "http://")
	pbHostPort = strings.TrimPrefix(pbHostPort, "https://")

	os.Args = append(os.Args, "--dir="+pbDataPath)
	os.Args = append(os.Args, "--http="+pbHostPort)

	app := pocketbase.New()

	// --- NATS JetStream Connection ---
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS server at %s: %v", natsURL, err)
	}
	defer nc.Close()
	log.Printf("Successfully connected to NATS server: %s", natsURL)

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error getting NATS JetStream context: %v", err)
	}

	// Ensure JetStream stream exists
	streamInfo, err := js.StreamInfo(jetStreamName)
	if err != nil {
		if err == nats.ErrStreamNotFound {
			log.Printf("JetStream stream '%s' not found, creating it...", jetStreamName)
			_, err = js.AddStream(&nats.StreamConfig{
				Name:      jetStreamName,
				Subjects:  []string{fmt.Sprintf("%s.>", jetStreamName)}, // e.g., POCKETBASE_EVENTS.users.created
				Storage:   nats.FileStorage,                             // Or nats.MemoryStorage
				Retention: nats.WorkQueuePolicy,                         // Or LimitsPolicy, InterestPolicy
			})
			if err != nil {
				log.Fatalf("Error creating JetStream stream '%s': %v", jetStreamName, err)
			}
			log.Printf("JetStream stream '%s' created.", jetStreamName)
		} else {
			log.Fatalf("Error getting stream info for '%s': %v", jetStreamName, err)
		}
	} else {
		log.Printf("Using existing JetStream stream: %s", streamInfo.Config.Name)
	}

	// --- PocketBase Hooks to Publish to NATS JetStream ---
	publishEventToJetStream := func(collectionName, action, recordId string, recordData map[string]any) {
		subject := fmt.Sprintf("%s.%s.%s", jetStreamName, collectionName, action)
		eventPayload := map[string]any{
			"timestamp":  time.Now().UTC().Format(time.RFC3339Nano),
			"action":     action,
			"collection": collectionName,
			"recordId":   recordId,
			"data":       recordData, // This is the full record data
		}

		jsonData, err := json.Marshal(eventPayload)
		if err != nil {
			log.Printf("Error marshalling event data for NATS: %v. Subject: %s", err, subject)
			return
		}

		ack, err := js.Publish(subject, jsonData)
		if err != nil {
			log.Printf("Error publishing to JetStream subject '%s': %v", subject, err)
		} else {
			log.Printf("Published to JetStream: Subject='%s', RecordID='%s', StreamSeq=%d", subject, recordId, ack.Sequence)
		}
	}

	app.OnRecordAfterCreateRequest().Add(func(e *core.RecordCreateEvent) error {
		if e.Record != nil {
			publishEventToJetStream(e.Record.Collection().Name, "create", e.Record.GetId(), e.Record.PublicExport())
		}
		return nil
	})

	app.OnRecordAfterUpdateRequest().Add(func(e *core.RecordUpdateEvent) error {
		if e.Record != nil {
			publishEventToJetStream(e.Record.Collection().Name, "update", e.Record.GetId(), e.Record.PublicExport())
		}
		return nil
	})

	app.OnRecordAfterDeleteRequest().Add(func(e *core.RecordDeleteEvent) error {
		if e.Record != nil {
			// For delete, PublicExport() gives the state of the record *before* deletion.
			publishEventToJetStream(e.Record.Collection().Name, "delete", e.Record.GetId(), e.Record.PublicExport())
		}
		return nil
	})

	// --- Optional: Setup for other custom routes if needed ---
	// If you had other non-Datastar routes, you could set them up here
	// using e.Router.GET, e.Router.POST, etc., or by mounting another
	// chi router instance if you prefer.
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Example: e.Router.GET("/custom-api/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong") })
		log.Println("Custom route setup (if any) completed in OnBeforeServe.")
		return nil
	})

	// --- Start PocketBase ---
	// This will also start the web server with the mounted Chi router.
	log.Println("Starting PocketBase server...")
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start PocketBase: %v", err)
	}

	log.Println("pb-nats application finished.") // Should not be reached if Start() is blocking
}
