package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar"
)

//go:embed index.html
var helloWorldHTML []byte

const serverPort = ":8080"

func main() {
	r := chi.NewRouter()

	const message = "Hello, world!"
	type Store struct {
		// Assuming 'delay' in JSON is a number representing milliseconds.
		Delay int64 `json:"delay"`
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(helloWorldHTML)
	})

	r.Get("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		store := &Store{Delay: 100} // Default delay if not provided by signals
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sse := datastar.NewSSE(w, r)

		for i := 0; i < len(message); i++ {
			sse.MergeFragments(`<div id="message">` + message[:i+1] + `</div>`)
			time.Sleep(time.Duration(store.Delay) * time.Millisecond)
		}
	})

	log.Printf("Starting server on port %s\n", serverPort)
	err := http.ListenAndServe(serverPort, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
