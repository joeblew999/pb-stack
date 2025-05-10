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

// debugMode can be set via a flag or environment variable in a real application
const debugMode = false // Set to true to enable debug logging like fs.WalkDir

func Execute(assets embed.FS, bootFlag bool, debootFlag bool, targetHost string, packageName string) {
	var scriptBaseName string

	if bootFlag {
		scriptBaseName = "boot"
	} else if debootFlag {
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

	// Adjust actionMessage construction slightly for new terms
	baseActionVerb := "Setting up"
	if scriptBaseName == "deboot" {
		baseActionVerb = "Tearing down"
	}
	actionMessage := baseActionVerb // e.g., "Setting up" or "Tearing down"
	if packageName != "" {
		actionMessage = fmt.Sprintf("%s for package '%s'", actionMessage, packageName)
	}
	if targetHost != "" {
		actionMessage = fmt.Sprintf("%s on target '%s'", actionMessage, targetHost)
	}
	fmt.Printf("%s...\n", actionMessage)

	var cmd *exec.Cmd
	var err error
	var tempScriptPath string // To hold the path of a temporary script, if one is created

	if packageName != "" && targetHost == "" {
		cmd, err = handleLocalPackageOperation(scriptBaseName, packageName)
		// No temporary script path for local package operations
	} else {
		cmd, tempScriptPath, err = handleScriptOperation(assets, scriptBaseName, targetHost, packageName)
		if tempScriptPath != "" {
			defer os.Remove(tempScriptPath) // Defer removal after Execute finishes
		}
	}

	if err != nil {
		log.Fatalf("Failed to prepare command: %v", err)
	}

	if cmd != nil {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = os.TempDir() // Set working directory for the script if necessary
		if runErr := cmd.Run(); runErr != nil {
			log.Fatalf("Failed to execute command '%s': %v", cmd.Path, runErr)
		}
	}
	fmt.Printf("%s complete.\n", actionMessage)
}

func handleLocalPackageOperation(scriptBaseName, packageName string) (*exec.Cmd, error) {
	fmt.Println("Performing local package operation...")
	var cmd *exec.Cmd
	var pkgManagerCmd string

	switch runtime.GOOS {
	case "darwin": // macOS
		pkgManagerCmd = "brew"
		if _, err := exec.LookPath(pkgManagerCmd); err != nil {
			return nil, fmt.Errorf("%s command not found in PATH: %w. Please install Homebrew", pkgManagerCmd, err)
		}
		opCmd := "install"
		if scriptBaseName == "deboot" {
			opCmd = "uninstall"
		}
		fmt.Printf("Using Homebrew: %s %s %s\n", pkgManagerCmd, opCmd, packageName)
		cmd = exec.Command(pkgManagerCmd, opCmd, packageName)
	case "linux":
		// Basic example for apt. A real-world scenario might need more robust detection.
		pkgManagerCmd = "apt" // Could also be apt-get
		if _, err := exec.LookPath(pkgManagerCmd); err != nil {
			// Try apt-get if apt is not found directly
			pkgManagerCmd = "apt-get"
			if _, errGet := exec.LookPath(pkgManagerCmd); errGet != nil {
				return nil, fmt.Errorf("neither apt nor apt-get command found in PATH: %w. Please ensure a Debian/Ubuntu based package manager is available", err)
			}
		}

		opCmd := "install"
		if scriptBaseName == "deboot" {
			opCmd = "remove"
		}

		sudo := "sudo"
		if _, err := exec.LookPath(sudo); err != nil {
			fmt.Println("sudo command not found, attempting to run package manager without it.")
			sudo = "" // sudo not found, try without
		}

		fmt.Printf("Using %s: %s %s %s -y %s\n", pkgManagerCmd, sudo, pkgManagerCmd, opCmd, packageName)
		if sudo != "" {
			cmd = exec.Command(sudo, pkgManagerCmd, opCmd, "-y", packageName)
		} else {
			cmd = exec.Command(pkgManagerCmd, opCmd, "-y", packageName)
		}
	case "windows":
		pkgManagerCmd = "winget"
		if _, err := exec.LookPath(pkgManagerCmd); err != nil {
			return nil, fmt.Errorf("%s command not found in PATH: %w. Please install Winget", pkgManagerCmd, err)
		}
		opCmd := "install"
		if scriptBaseName == "deboot" {
			opCmd = "uninstall"
		}
		fmt.Printf("Using Winget: %s %s --source winget --exact --id %s --accept-package-agreements --accept-source-agreements\n", pkgManagerCmd, opCmd, packageName)
		cmd = exec.Command(pkgManagerCmd, opCmd, "--source", "winget", "--exact", "--id", packageName, "--accept-package-agreements", "--accept-source-agreements")
	default:
		return nil, fmt.Errorf("unsupported OS for local package operation: %s", runtime.GOOS)
	}
	return cmd, nil
}

func handleScriptOperation(assets embed.FS, scriptBaseName, targetHost, packageName string) (*exec.Cmd, string, error) {
	if packageName != "" && targetHost != "" {
		fmt.Println("Performing remote package operation via script...")
	} else {
		fmt.Println("Performing generic script operation...")
	}

	if debugMode {
		fmt.Println("Listing embedded migration assets (debug mode):")
		fs.WalkDir(assets, "migrations", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Warning: Error walking embedded assets at %s: %v", path, err)
				if d != nil && d.IsDir() {
					return fs.SkipDir // Skip this directory if error occurs on the directory itself
				}
				return nil // Continue walking other files/dirs even if one file errors
			}
			fmt.Println(path)
			return nil
		})
	}

	fmt.Println("Determining script type for execution...")

	var scriptArgs []string
	if targetHost != "" {
		scriptArgs = append(scriptArgs, targetHost)
	}
	if packageName != "" {
		scriptArgs = append(scriptArgs, packageName)
	}

	var cmd *exec.Cmd
	var tempFilePath string
	var err error

	// The script type (.sh or .ps1) is chosen based on the OS running this Go program.
	// The script itself must handle remoting if targetHost is specified.
	switch runtime.GOOS {
	case "darwin", "linux":
		cmd, tempFilePath, err = prepareScriptCmd(assets, scriptBaseName, "sh", true, scriptArgs)
		if err != nil {
			return nil, "", fmt.Errorf("error preparing shell script: %w", err)
		}
	case "windows":
		cmd, tempFilePath, err = prepareScriptCmd(assets, scriptBaseName, "ps1", false, scriptArgs)
		if err != nil {
			return nil, "", fmt.Errorf("error preparing PowerShell script: %w", err)
		}
	default:
		return nil, "", fmt.Errorf("unsupported OS for script execution: %s", runtime.GOOS)
	}
	return cmd, tempFilePath, nil
}

func prepareScriptCmd(assets embed.FS, scriptBaseName, scriptSuffix string, isUnixLike bool, scriptArgs []string) (*exec.Cmd, string, error) {
	scriptName := fmt.Sprintf("migrations/%s.%s", scriptBaseName, scriptSuffix)
	scriptBytes, err := assets.ReadFile(scriptName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read embedded script %s: %w", scriptName, err)
	}

	tempFile, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("pb-stack-boot-*.%s", scriptSuffix))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp file for %s: %w", scriptName, err)
	}
	tempFilePath := tempFile.Name()

	if _, err := tempFile.Write(scriptBytes); err != nil {
		tempFile.Close()
		os.Remove(tempFilePath)
		return nil, "", fmt.Errorf("failed to write to temp file for %s: %w", scriptName, err)
	}

	if isUnixLike {
		if err := tempFile.Chmod(0755); err != nil {
			tempFile.Close()
			os.Remove(tempFilePath)
			return nil, "", fmt.Errorf("failed to set executable permission for temp file %s: %w", scriptName, err)
		}
	}

	if err := tempFile.Close(); err != nil {
		os.Remove(tempFilePath)
		return nil, "", fmt.Errorf("failed to close temp file for %s: %w", scriptName, err)
	}

	var cmd *exec.Cmd
	if isUnixLike {
		cmd = exec.Command(tempFilePath, scriptArgs...)
	} else { // PowerShell
		psArgs := []string{"-ExecutionPolicy", "Bypass", "-File", tempFilePath}
		psArgs = append(psArgs, scriptArgs...)
		cmd = exec.Command("powershell", psArgs...)
	}
	return cmd, tempFilePath, nil
}
