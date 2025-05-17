package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var tpl *template.Template

func main() {
	// Get the directory of the executable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Construct the path to the templates directory relative to the executable
	// Assuming templates are in ../../web/templates from cmd/server/server(.exe)
	// Adjust this path if your structure is different or if you run `go run` from project root.
	// For `go run ./cmd/server/main.go` from project root, templatesDir would be "web/templates"
	templatesDir := filepath.Join(filepath.Dir(filepath.Dir(exeDir)), "web", "templates")
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		// Fallback for `go run` from project root
		wd, _ := os.Getwd()
		templatesDir = filepath.Join(wd, "web", "templates")
	}

	tpl, err = template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		log.Fatalf("Failed to parse templates from %s: %v", templatesDir, err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/api/git/placeholder", gitApiPlaceholderHandler)

	log.Println("Datastar backend server starting on :8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func gitApiPlaceholderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Write([]byte("Git API placeholder received POST request"))
		return
	}
	w.Write([]byte("This is a Git API placeholder. Try POSTing."))
}
