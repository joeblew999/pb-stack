package cli

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func Execute(assets embed.FS, bootFlag bool, debootFlag bool, targetHost string) {
	var action string
	var scriptBaseName string

	if bootFlag {
		action = "Booting (CLI)"
		scriptBaseName = "boot"
	} else if debootFlag {
		action = "Debooting (CLI)"
		scriptBaseName = "deboot"
	} else {
		fmt.Println("CLI mode selected. Please specify an action:")
		programName := "your-program" // Default
		if len(os.Args) > 0 {
			programName = filepath.Base(os.Args[0])
		}
		fmt.Printf("  %s -cli -boot    (to run boot scripts)\n", programName)
		fmt.Printf("  %s -cli -deboot  (to run deboot scripts)\n", programName)
		os.Exit(0) // Exit cleanly after showing options
	}
	actionMessage := fmt.Sprintf("%s up", action)
	if targetHost != "" {
		actionMessage = fmt.Sprintf("%s on target '%s'", actionMessage, targetHost)
	}
	fmt.Printf("%s...\n", actionMessage)

	// Example: List files in the embedded migrations directory
	fs.WalkDir(assets, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(path)
		return nil
	})

	fmt.Println("Checking OS ...")
	fmt.Printf("Detected OS: %s\n", runtime.GOOS)

	var scriptName string
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		{
			scriptName = fmt.Sprintf("migrations/%s.sh", scriptBaseName) // Path within the embedded FS
			scriptBytes, err := assets.ReadFile(scriptName)
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
			if err := tempFile.Chmod(0755); err != nil {
				tempFile.Close()
				log.Fatalf("Failed to set executable permission for temp file %s: %v", scriptName, err)
			}

			tempFilePath := tempFile.Name()
			if err := tempFile.Close(); err != nil { // Close the file before executing
				log.Fatalf("Failed to close temp file for %s: %v", scriptName, err)
			}
			if targetHost != "" {
				cmd = exec.Command(tempFilePath, targetHost)
			} else {
				cmd = exec.Command(tempFilePath)
			}
		}
	case "windows":
		{
			scriptName = fmt.Sprintf("migrations/%s.ps1", scriptBaseName) // Path within the embedded FS
			scriptBytes, err := assets.ReadFile(scriptName)
			if err != nil {
				log.Fatalf("Failed to read embedded script %s: %v", scriptName, err)
			}

			tempFile, err := os.CreateTemp(os.TempDir(), "pb-stack-boot-*.ps1")
			if err != nil {
				log.Fatalf("Failed to create temp file for %s: %v", scriptName, err)
			}
			defer os.Remove(tempFile.Name())

			if _, err := tempFile.Write(scriptBytes); err != nil {
				tempFile.Close()
				log.Fatalf("Failed to write to temp file for %s: %v", scriptName, err)
			}
			tempFilePath := tempFile.Name()
			if err := tempFile.Close(); err != nil {
				log.Fatalf("Failed to close temp file for %s: %v", scriptName, err)
			}
			if targetHost != "" {
				cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempFilePath, targetHost)
			} else {
				cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempFilePath)
			}
		}
	default:
		log.Fatalf("Unsupported OS for CLI mode: %s", runtime.GOOS)
	}

	if cmd != nil {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = os.TempDir() // Set working directory for the script if necessary
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to execute %s: %v", scriptName, err)
		}
	}
	fmt.Printf("%s complete.\n", actionMessage)
}
