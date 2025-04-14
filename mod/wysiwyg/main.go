// /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/mod/wysiwyg/main.go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	// Ensure you have your migrations directory structure setup
	// (e.g., 'migrations/1_init.up.sql') if using migrations.
	// Adjust the import path if your migrations are elsewhere relative to main.go
	// _ "github.com/joeblew999/pb-stack/mod/wysiwyg/migrations" // Example if migrations are in a subpackage
)

func main() {
	// --- Optional: Read Environment Variables ---
	// PocketBase CLI flags usually override these, but they can be useful
	// for custom logic within your Go code if needed.
	// Note: The Makefile sets these, but PocketBase's standard 'serve'
	// command doesn't directly use POCKETBASE_ADMIN_EMAIL/PASSWORD env vars
	// for initial admin creation (that's usually interactive or via flags).
	// POCKETBASE_URL is also usually determined by the --http flag or default.
	pbURL := os.Getenv("POCKETBASE_URL")
	adminEmail := os.Getenv("POCKETBASE_ADMIN_EMAIL")
	adminPassword := os.Getenv("POCKETBASE_ADMIN_PASSWORD") // Be careful with passwords

	log.Printf("Info: POCKETBASE_URL (from env): %s", pbURL)
	log.Printf("Info: POCKETBASE_ADMIN_EMAIL (from env): %s", adminEmail)
	if adminPassword != "" {
		log.Println("Info: POCKETBASE_ADMIN_PASSWORD (from env): [set]")
	} else {
		log.Println("Info: POCKETBASE_ADMIN_PASSWORD (from env): [not set]")
	}
	// --- End Optional Environment Variable Reading ---

	// --- Determine Data Directory ---
	// This allows `go run .` to work correctly by finding the executable's dir
	// You might want to customize this further.
	var dataDir string
	if goRunPath := os.Getenv("GORUN_WD"); goRunPath != "" {
		// If running via 'go run .', use the original working directory
		dataDir = filepath.Join(goRunPath, "pb_data")
	} else {
		// Otherwise, assume рядом с бинарником (next to the built binary)
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to get executable path: %v", err)
		}
		dataDir = filepath.Join(filepath.Dir(exePath), "pb_data")
	}
	log.Printf("Info: Using data directory: %s", dataDir)
	// --- End Determine Data Directory ---

	// --- Initialize PocketBase ---
	// Pass the data directory explicitly if needed, otherwise PocketBase uses default flags
	// app := pocketbase.NewWithConfig(pocketbase.Config{
	// 	DefaultDataDir: dataDir,
	// })
	// Or simply:
	app := pocketbase.New()

	// --- Migrations ---
	// Uncomment the next line if you have migrations defined in Go files
	// (requires the migrations directory and import above)
	// migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
	// 	Automigrate: true, // Auto-run migrations on startup
	// })

	// --- Custom WYSIWYG Hooks/Logic ---
	// This is where you'd add your specific module's functionality.
	// For example, adding a hook before a record is created/updated
	// in a collection that uses your WYSIWYG editor.
	app.OnRecordBeforeUpdateRequest("your_collection_name").Add(func(e *core.RecordUpdateEvent) error {
		log.Printf("WYSIWYG Hook: Processing update for record %s in collection %s", e.Record.Id, e.Collection.Name)

		// Example: Sanitize HTML content from a field named 'wysiwyg_content'
		// content := e.Record.GetString("wysiwyg_content")
		// sanitizedContent := yourSanitizeFunction(content) // Replace with actual sanitization
		// e.Record.Set("wysiwyg_content", sanitizedContent)

		// Example: Check for something specific in the content
		// if strings.Contains(content, "<script>") {
		// 	return apis.NewBadRequestError("Script tags are not allowed", nil)
		// }

		return nil // Return nil to allow the update, or an error to prevent it
	})

	app.OnRecordBeforeCreateRequest("your_collection_name").Add(func(e *core.RecordCreateEvent) error {
		log.Printf("WYSIWYG Hook: Processing create for collection %s", e.Collection.Name)
		// Add similar logic as above for creation if needed
		// content := e.Record.GetString("wysiwyg_content")
		// ...
		return nil
	})

	// Example: Add a custom API route if needed
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Example: e.Router.AddRoute(echo.Route{...})
		log.Println("WYSIWYG Module: Adding custom routes or middleware if necessary...")
		return nil
	})
	// --- End Custom WYSIWYG Hooks/Logic ---

	// --- Start PocketBase ---
	// This will parse command-line flags (like --http, --dir) and start the server
	// or execute other commands like 'migrate'.
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// Placeholder for a sanitization function (replace with a real one like bluemonday)
// func yourSanitizeFunction(htmlContent string) string {
// 	// p := bluemonday.UGCPolicy()
// 	// return p.Sanitize(htmlContent)
// 	log.Println("Warning: HTML sanitization not implemented!")
// 	return htmlContent
// }
