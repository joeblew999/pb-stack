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
	"sort"
)

func Execute(assets embed.FS, bootFlag bool, debootFlag bool, packageName string, appDebugMode bool) {
	var scriptBaseName string

	if bootFlag {
		scriptBaseName = "boot"
	} else if debootFlag {
		scriptBaseName = "deboot"
	} else {
		log.Println("CLI mode selected. Please specify an action:") // Changed to log
		programName := "your-program"                               // Default
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
	log.Printf("%s...\n", actionMessage) // Changed to log

	var cmd *exec.Cmd
	var err error

	if packageName != "" { // If a specific package is named, handle it directly (always local now)
		cmd, err = handleLocalPackageOperation(scriptBaseName, packageName)
		if err != nil {
			log.Fatalf("Failed to prepare local package operation: %v", err)
		}
		if cmd != nil { // If a command was actually prepared (e.g., OS supported)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if runErr := cmd.Run(); runErr != nil {
				log.Fatalf("Failed to execute local package operation '%s': %v", cmd.Path, runErr)
			}
		}
	} else {
		err = executeScriptOperations(assets, scriptBaseName, "", packageName, appDebugMode) // Pass empty string for targetHost
		if err != nil {
			log.Fatalf("Script execution failed: %v", err)
		}
	}
	log.Printf("%s complete.\n", actionMessage) // Changed to log
}

func handleLocalPackageOperation(scriptBaseName, packageName string) (*exec.Cmd, error) {
	log.Println("Performing local package operation...") // Changed to log
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
		log.Printf("Using Homebrew: %s %s %s\n", pkgManagerCmd, opCmd, packageName) // Changed to log
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
			log.Println("sudo command not found, attempting to run package manager without it.") // Changed to log
			sudo = ""                                                                            // sudo not found, try without
		}

		log.Printf("Using %s: %s %s %s -y %s\n", pkgManagerCmd, sudo, pkgManagerCmd, opCmd, packageName) // Changed to log
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
		log.Printf("Using Winget: %s %s --source winget --exact --id %s --accept-package-agreements --accept-source-agreements\n", pkgManagerCmd, opCmd, packageName) // Changed to log
		cmd = exec.Command(pkgManagerCmd, opCmd, "--source", "winget", "--exact", "--id", packageName, "--accept-package-agreements", "--accept-source-agreements")
	default:
		return nil, fmt.Errorf("unsupported OS for local package operation: %s", runtime.GOOS)
	}
	return cmd, nil
}

func executeScriptOperations(assets embed.FS, scriptBaseName, _ /* targetHost - no longer used */ string, packageName string, appDebugMode bool) error {
	var tempExtensionsPath string
	// Prepare extensions.txt first, as its path will be an argument to scripts
	extensionsFilePathInAssets := "migrations/extensions.txt"
	extensionsBytes, err := assets.ReadFile(extensionsFilePathInAssets)
	// Defer cleanup of extensions.txt if created
	defer func() { // Ensure tempExtensionsPath is captured by the closure if it gets set
		if tempExtensionsPath != "" {
			defer os.Remove(tempExtensionsPath)
		}
	}()

	if appDebugMode {
		fmt.Println("Listing embedded migration assets (debug mode):")
		fs.WalkDir(assets, "migrations", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Warning: Error walking embedded assets at %s: %v", path, err)
				if d != nil && d.IsDir() {
					return fs.SkipDir // Skip this directory if error occurs on the directory itself
				}
				return nil // Continue walking other files/dirs even if one file errors
			}
			log.Printf("[DEBUG] Embedded asset: %s", path) // Changed to log
			return nil
		})
	}

	if err == nil { // extensions.txt File exists and was read
		extTempFile, errCreate := os.CreateTemp(os.TempDir(), "pb-stack-extensions-*.txt")
		if errCreate != nil {
			return fmt.Errorf("failed to create temp file for extensions.txt: %w", errCreate)
		}
		tempExtensionsPath = extTempFile.Name() // Assign here
		// Now that tempExtensionsPath is assigned, the defer above will work if this function exits.

		if _, errWrite := extTempFile.Write(extensionsBytes); errWrite != nil {
			extTempFile.Close()
			// os.Remove(tempExtensionsPath) // defer will handle this if path was set
			return fmt.Errorf("failed to write to temp file for extensions.txt: %w", errWrite)
		}
		if errClose := extTempFile.Close(); errClose != nil {
			// os.Remove(tempExtensionsPath) // defer will handle this
			return fmt.Errorf("failed to close temp file for extensions.txt: %w", errClose)
		}
		log.Printf("Found and prepared extensions.txt at temporary path: %s\n", tempExtensionsPath) // Changed to log
	} else if !os.IsNotExist(err) { // An error other than "not found"
		log.Printf("Warning: Error reading embedded extensions.txt: %v. Proceeding without it.", err)
	} // If os.IsNotExist(err), we just proceed without it silently.

	var scriptArgs []string
	// targetHost argument removed

	// Note: packageName is handled differently now for single vs multi-script
	if tempExtensionsPath != "" {
		scriptArgs = append(scriptArgs, tempExtensionsPath) // Common for all scripts if extensions.txt exists
	}

	isUnixLike := runtime.GOOS == "darwin" || runtime.GOOS == "linux"
	scriptSuffix := ".sh"
	if !isUnixLike {
		if runtime.GOOS == "windows" {
			scriptSuffix = ".ps1"
		} else {
			return fmt.Errorf("unsupported OS for script execution: %s", runtime.GOOS)
		}
	}

	if packageName != "" {
		// Single script execution for a specific package
		log.Printf("Performing operation for package '%s' via script...\n", packageName) // Changed to log
		singleScriptArgs := append([]string{}, scriptArgs...)                            // Copy base args
		singleScriptArgs = append(singleScriptArgs, packageName)                         // Add packageName specifically for this script
		// The script name for package operations is still the generic boot/deboot
		scriptPathInAssets := filepath.Join("migrations", scriptBaseName+scriptSuffix)
		return runSingleScript(assets, scriptPathInAssets, isUnixLike, singleScriptArgs)
	} else {
		// Migration-style: multiple ordered scripts
		scriptPrefix := "setup_"
		if scriptBaseName == "deboot" {
			scriptPrefix = "teardown_"
		}
		log.Printf("Performing migration-style operation with prefix '%s'...\n", scriptPrefix) // Changed to log

		entries, err := assets.ReadDir("migrations")
		if err != nil {
			return fmt.Errorf("failed to list embedded migration scripts: %w", err)
		}

		var scriptsToRun []string
		for _, entry := range entries {
			if !entry.IsDir() && filepath.HasPrefix(entry.Name(), scriptPrefix) && filepath.Ext(entry.Name()) == scriptSuffix {
				scriptsToRun = append(scriptsToRun, entry.Name())
			}
		}
		sort.Strings(scriptsToRun) // Ensure alphabetical order

		if len(scriptsToRun) == 0 {
			log.Printf("No scripts found matching pattern '%s*%s'. Nothing to do.\n", scriptPrefix, scriptSuffix) // Changed to log
			return nil
		}

		log.Printf("Found migration scripts to run: %v\n", scriptsToRun) // Changed to log
		for _, scriptFilename := range scriptsToRun {
			scriptPathInAssets := filepath.Join("migrations", scriptFilename)
			log.Printf("Executing script: %s\n", scriptPathInAssets) // Changed to log
			// scriptArgs for migration scripts do not include packageName
			err := runSingleScript(assets, scriptPathInAssets, isUnixLike, scriptArgs)
			if err != nil {
				return fmt.Errorf("error executing script %s: %w", scriptFilename, err)
			}
		}
	}
	return nil
}

// runSingleScript prepares and executes one script.
func runSingleScript(assets embed.FS, scriptPathInAssets string, isUnixLike bool, scriptArgs []string) error {
	cmd, tempScriptPath, err := prepareScriptCmd(assets, scriptPathInAssets, isUnixLike, scriptArgs)
	if err != nil {
		return err
	}
	defer os.Remove(tempScriptPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = os.TempDir()
	return cmd.Run()
}

func prepareScriptCmd(assets embed.FS, scriptPathInAssets string, isUnixLike bool, scriptArgs []string) (*exec.Cmd, string, error) {
	scriptBytes, err := assets.ReadFile(scriptPathInAssets)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read embedded script %s: %w", scriptPathInAssets, err)
	}

	// Determine suffix from scriptPathInAssets for temp file
	scriptSuffix := filepath.Ext(scriptPathInAssets)
	tempFile, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("pb-stack-script-*%s", scriptSuffix))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp file for %s: %w", scriptPathInAssets, err)
	}
	tempFilePath := tempFile.Name()

	if _, err := tempFile.Write(scriptBytes); err != nil {
		tempFile.Close()
		os.Remove(tempFilePath)
		return nil, "", fmt.Errorf("failed to write to temp file for %s: %w", scriptPathInAssets, err)
	}

	if isUnixLike {
		if err := tempFile.Chmod(0755); err != nil {
			tempFile.Close()
			os.Remove(tempFilePath)
			return nil, "", fmt.Errorf("failed to set executable permission for temp file %s: %w", scriptPathInAssets, err)
		}
	}

	if err := tempFile.Close(); err != nil {
		os.Remove(tempFilePath)
		return nil, "", fmt.Errorf("failed to close temp file for %s: %w", scriptPathInAssets, err)
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
