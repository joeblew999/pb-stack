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

func Execute(assets embed.FS, bootFlag bool, debootFlag bool, targetHost string, packageName string) {
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
	if packageName != "" {
		actionMessage = fmt.Sprintf("%s for package '%s'", actionMessage, packageName)
	}
	if targetHost != "" {
		actionMessage = fmt.Sprintf("%s on target '%s'", actionMessage, targetHost)
	}
	fmt.Printf("%s...\n", actionMessage)

	var scriptName string
	var cmd *exec.Cmd

	if packageName != "" && targetHost == "" {
		// Local package operation - direct command execution
		fmt.Println("Performing local package operation...")
		switch runtime.GOOS {
		case "darwin": // macOS
			pkgCmd := "install"
			if scriptBaseName == "deboot" {
				pkgCmd = "uninstall"
			}
			fmt.Printf("Using Homebrew: brew %s %s\n", pkgCmd, packageName)
			cmd = exec.Command("brew", pkgCmd, packageName)
		case "linux":
			// Assuming apt for Linux, could be yum, dnf, etc.
			// This part might need more sophisticated OS/distro detection or configuration
			pkgCmd := "install"
			sudo := "sudo" // Most package managers need sudo
			if _, err := exec.LookPath("sudo"); err != nil {
				sudo = "" // sudo not found, try without
			}
			if scriptBaseName == "deboot" {
				pkgCmd = "remove"
			}
			// Example for apt
			fmt.Printf("Using apt: %s apt %s -y %s\n", sudo, pkgCmd, packageName)
			if sudo != "" {
				cmd = exec.Command(sudo, "apt", pkgCmd, "-y", packageName)
			} else {
				cmd = exec.Command("apt", pkgCmd, "-y", packageName)
			}
		case "windows":
			pkgCmd := "install"
			if scriptBaseName == "deboot" {
				pkgCmd = "uninstall"
			}
			// Winget often needs specific IDs. The packageName should be the ID.
			// Adding common flags for non-interactive use.
			fmt.Printf("Using Winget: winget %s --source winget --exact --id %s --accept-package-agreements --accept-source-agreements\n", pkgCmd, packageName)
			cmd = exec.Command("winget", pkgCmd, "--source", "winget", "--exact", "--id", packageName, "--accept-package-agreements", "--accept-source-agreements")
		default:
			log.Fatalf("Unsupported OS for local package operation: %s", runtime.GOOS)
		}
	} else {
		// Generic script operation OR remote package operation (handled by the script)
		if packageName != "" && targetHost != "" {
			fmt.Println("Performing remote package operation via script...")
		} else {
			fmt.Println("Performing generic script operation...")
		}

		// List files in the embedded migrations directory (can be removed if not needed for debugging)
		fs.WalkDir(assets, "migrations", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Warning: Error walking embedded assets: %v", err) // Log as warning
				return fs.SkipDir                                             // Skip this directory if error occurs
			}
			// fmt.Println(path) // Can be verbose
			return nil
		})

		fmt.Println("Checking OS for script execution context (local or remote)...")
		// Note: runtime.GOOS here is the OS of the machine running *this Go program*.
		// If targetHost is set, the script runs on the target, but this Go app still needs to pick the right script type (.sh or .ps1).
		// For simplicity, we assume if a target is Linux/macOS, we use .sh, if Windows, .ps1.
		// A more robust solution might involve probing the target or having user specify target OS type.

		var scriptArgs []string
		if targetHost != "" {
			scriptArgs = append(scriptArgs, targetHost) // Script receives target host as $1 (or %1)
		}
		if packageName != "" { // If it's a package op (and thus targetHost is also set for this branch)
			scriptArgs = append(scriptArgs, packageName) // Script receives package name as $2 (or %2)
		}

		// Determine script type based on typical target OS if targetHost is set, else local OS.
		// This is a heuristic. If targetHost is "linux-server" but this Go app runs on Windows, we still want to use .sh.
		// For now, we'll rely on the Go app's OS to select the script type, assuming the admin knows which script to make.
		// This means if you run this on Windows to target a Linux machine, you'd ideally want to execute a .sh script via SSH.
		// The current model just extracts and runs locally, or passes to a local powershell/bash.
		// For true remote execution, this part needs to change significantly (e.g. invoke ssh or WinRM).
		// Given the current structure, we assume the script is run *locally* and if it's for a target, the script handles remoting.

		switch runtime.GOOS { // This is the OS of the machine running the Go program
		case "darwin", "linux":
			scriptName = fmt.Sprintf("migrations/%s.sh", scriptBaseName)
			scriptBytes, err := assets.ReadFile(scriptName)
			if err != nil {
				log.Fatalf("Failed to read embedded script %s: %v", scriptName, err)
			}
			tempFile, err := os.CreateTemp(os.TempDir(), "pb-stack-boot-*.sh")
			if err != nil {
				log.Fatalf("Failed to create temp file for %s: %v", scriptName, err)
			}
			defer os.Remove(tempFile.Name())
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
			cmd = exec.Command(tempFilePath, scriptArgs...)

		case "windows":
			scriptName = fmt.Sprintf("migrations/%s.ps1", scriptBaseName)
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
				tempFile.Close() // Close before attempting remove on error
				log.Fatalf("Failed to write to temp file for %s: %v", scriptName, err)
			}
			tempFilePath := tempFile.Name()
			if err := tempFile.Close(); err != nil { // Close the file before executing
				log.Fatalf("Failed to close temp file for %s: %v", scriptName, err)
			}
			psArgs := []string{"-ExecutionPolicy", "Bypass", "-File", tempFilePath}
			psArgs = append(psArgs, scriptArgs...)
			cmd = exec.Command("powershell", psArgs...)
		default:
			log.Fatalf("Unsupported OS for script execution: %s", runtime.GOOS)
		}
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
