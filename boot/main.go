package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec" // Import the filepath package
	"runtime" // Import the runtime package
)

//go:embed all:migrations
var embeddedAssets embed.FS // Embed migrations folder

func main() {
	debootFlag := flag.Bool("deboot", false, "Run deboot scripts instead of boot scripts")
	flag.Parse()

	action := "Booting"
	if *debootFlag {
		action = "Debooting"
	}
	fmt.Printf("%s up...\n", action)

	// Example: List files in the embedded migrations directory (using the new embeddedAssets)
	fs.WalkDir(embeddedAssets, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(path)
		return nil
	})

	fmt.Println("Checking OS ...")

	fmt.Printf("Detected OS: %s\n", runtime.GOOS) // Use runtime.GOOS

	scriptBaseName := "boot"
	if *debootFlag {
		scriptBaseName = "deboot"
	}
	switch runtime.GOOS { // Use runtime.GOOS
	case "darwin", "linux":
		{
			scriptName := fmt.Sprintf("migrations/%s.sh", scriptBaseName) // Path within the embedded FS
			scriptBytes, err := embeddedAssets.ReadFile(scriptName)
			if err != nil {
				log.Fatalf("Failed to read embedded script %s: %v", scriptName, err)
			}

			tempFile, err := os.CreateTemp(os.TempDir(), "pb-stack-boot-*.sh")
			if err != nil {
				log.Fatalf("Failed to create temp file for %s: %v", scriptName, err)
			}
			defer os.Remove(tempFile.Name()) // Clean up the temp file

			if _, err := tempFile.Write(scriptBytes); err != nil {
				tempFile.Close() // Close before attempting remove on error
				log.Fatalf("Failed to write to temp file for %s: %v", scriptName, err)
			}
			// Make the script executable
			if err := tempFile.Chmod(0755); err != nil {
				tempFile.Close()
				log.Fatalf("Failed to set executable permission for temp file %s: %v", scriptName, err)
			}

			tempFilePath := tempFile.Name()
			if err := tempFile.Close(); err != nil { // Close the file before executing
				log.Fatalf("Failed to close temp file for %s: %v", scriptName, err)
			}

			cmd := exec.Command(tempFilePath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			// Set working directory for the script.
			// Using os.TempDir() to avoid locking the executable's directory,
			// which can cause issues with `go run` cleanup.
			// If the script needs to be in its own directory (e.g., to access co-located files
			// that are also extracted), then cmd.Dir should be filepath.Dir(tempFilePath).
			cmd.Dir = os.TempDir()
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to execute %s: %v", scriptName, err)
			}
		}
	case "windows":
		{
			scriptName := fmt.Sprintf("migrations/%s.ps1", scriptBaseName) // Path within the embedded FS
			scriptBytes, err := embeddedAssets.ReadFile(scriptName)
			if err != nil {
				log.Fatalf("Failed to read embedded script %s: %v", scriptName, err)
			}

			tempFile, err := os.CreateTemp(os.TempDir(), "pb-stack-boot-*.ps1")
			if err != nil {
				log.Fatalf("Failed to create temp file for %s: %v", scriptName, err)
			}
			defer os.Remove(tempFile.Name()) // Clean up the temp file

			if _, err := tempFile.Write(scriptBytes); err != nil {
				tempFile.Close() // Close before attempting remove on error
				log.Fatalf("Failed to write to temp file for %s: %v", scriptName, err)
			}

			tempFilePath := tempFile.Name()
			if err := tempFile.Close(); err != nil { // Close the file before passing its name to PowerShell
				log.Fatalf("Failed to close temp file for %s: %v", scriptName, err)
			}

			cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempFilePath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			// Set working directory for the script.
			// Using os.TempDir() to avoid locking the executable's directory,
			// which can cause issues with `go run` cleanup on Windows.
			// If the script needs to be in its own directory (e.g., to access co-located files
			// that are also extracted), then cmd.Dir should be filepath.Dir(tempFilePath).
			cmd.Dir = os.TempDir()
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to execute %s: %v", scriptName, err)
			}
		}
	}

	fmt.Printf("%s up complete.\n", action)
}
