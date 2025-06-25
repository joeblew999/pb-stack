package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := ":8085"
	webDir := "./web"

	// Serve static files from web directory
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	fmt.Printf("🎮 DataStar WASM WebWorkers Demo Server\n")
	fmt.Printf("🌐 Serving %s on http://localhost%s\n", webDir, port)
	fmt.Printf("📦 WASM files available at /wasm/\n")
	fmt.Printf("🚀 Open http://localhost%s in your browser\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
