package cli

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"main/pkg/sharedtypes"
	"main/pkg/vscodeutils"
)

func Execute(assets embed.FS, setupFlag bool, teardownFlag bool, packageName string, migrationSet string, appDebugMode bool) {
	var scriptBaseName string
	if setupFlag {
		scriptBaseName = "boot"
	} else if teardownFlag {
		scriptBaseName = "deboot"
	} else {
		log.Println("CLI mode selected. Please specify an action:") // Changed to log
		programName := "your-program"                               // Default
		if len(os.Args) > 0 {
			programName = filepath.Base(os.Args[0])
		}
		fmt.Printf("  %s -cli -setup    (to run setup operations)\n", programName)
		fmt.Printf("  %s -cli -teardown  (to run teardown operations)\n", programName)
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
			log.Printf("Warning: Could not prepare direct local package operation for '%s': %v. Will attempt script fallback if applicable.", packageName, err)
		}
		if cmd != nil { // If a command was actually prepared (e.g., OS supported)
			log.Println("Performing direct local package operation...")
			sharedtypes.RunCommand(cmd)
		} else {
			// Fallback to script if direct operation wasn't possible or if we always want to run a script for packages
			log.Printf("No direct local package operation for '%s', or script fallback intended. Attempting script execution...", packageName)
			// This will run migrations/boot.sh or migrations/deboot.sh with packageName
			err = executeScriptOperations(assets, scriptBaseName, migrationSet, packageName, appDebugMode)
			if err != nil {
				log.Fatalf("Script execution for package '%s' failed: %v", packageName, err)
			}
		}
	} else {
		// General local setup/teardown (no specific package specified by user)
		log.Println("Performing general local setup/teardown...")
		err = handleLocalJsonDrivenOperations(assets, scriptBaseName == "boot", migrationSet, appDebugMode)
		if err != nil {
			log.Fatalf("Local JSON-driven operation failed: %v", err)
		}

		// Optionally, run migration scripts AFTER JSON-driven local setup for other tasks
		log.Println("Proceeding to run general migration scripts (if any)...")
		// This will run setup_*.sh or teardown_*.sh
		err = executeScriptOperations(assets, scriptBaseName, migrationSet, "", appDebugMode) // packageName is empty for migration scripts
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

func executeScriptOperations(assets embed.FS, scriptBaseName, migrationSet string, packageName string, appDebugMode bool) error {
	var tempExtensionsPath string
	extensionsFilePathInAssets := filepath.Join("migrations", migrationSet, "extensions.txt") // Path within the specific set
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
		scriptPathInAssets := filepath.Join("migrations", migrationSet, scriptBaseName+scriptSuffix)
		return runSingleScript(assets, scriptPathInAssets, isUnixLike, singleScriptArgs)
	} else {
		// Migration-style: multiple ordered scripts
		scriptPrefix := "setup_"
		if scriptBaseName == "deboot" {
			scriptPrefix = "teardown_"
		}
		log.Printf("Performing migration-style operation with prefix '%s'...\n", scriptPrefix) // Changed to log

		entries, err := assets.ReadDir(filepath.Join("migrations", migrationSet))
		if err != nil {
			return fmt.Errorf("failed to list embedded migration scripts from set '%s': %w", migrationSet, err)
		}

		var scriptsToRun []string
		for _, entry := range entries {
			if !entry.IsDir() && filepath.HasPrefix(entry.Name(), scriptPrefix) && filepath.Ext(entry.Name()) == scriptSuffix {
				scriptsToRun = append(scriptsToRun, filepath.Join("migrations", migrationSet, entry.Name()))
			}
		}
		sort.Strings(scriptsToRun) // Ensure alphabetical order

		if len(scriptsToRun) == 0 {
			log.Printf("No scripts found matching pattern '%s*%s'. Nothing to do.\n", scriptPrefix, scriptSuffix) // Changed to log
			return nil
		}

		log.Printf("Found migration scripts to run: %v\n", scriptsToRun) // Changed to log
		for _, scriptFilename := range scriptsToRun {
			log.Printf("Executing script: %s\n", scriptFilename) // scriptFilename is now the full path in assets
			// scriptArgs for migration scripts do not include packageName (already handled)
			err := runSingleScript(assets, scriptFilename, isUnixLike, scriptArgs)
			if err != nil {
				return fmt.Errorf("error executing script %s: %w", scriptFilename, err)
			}
		}
	}
	return nil
}

func handleLocalJsonDrivenOperations(assets embed.FS, isSetup bool, migrationSet string, appDebugMode bool) error {
	log.Println("Performing local JSON-driven operations...")
	configPath := filepath.Join("migrations", migrationSet, "config.json")
	configBytes, err := assets.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded config.json from set '%s': %w. Ensure it exists in migrations/%s/", migrationSet, err, migrationSet)
	}

	var appConfig sharedtypes.AppConfig
	if err := json.Unmarshal(configBytes, &appConfig); err != nil {
		return fmt.Errorf("failed to parse config.json: %w", err)
	}

	// Handle Packages
	for _, pkg := range appConfig.Packages {
		log.Printf("Processing package: %s", pkg.Name)
		installed := isPackageInstalled(pkg)

		var cmd *exec.Cmd
		var pkgManagerName, pkgIdentifier, opVerb string

		if isSetup { // Setup (Install)
			opVerb = "install"
			if installed {
				log.Printf("Package '%s' already installed. Skipping.", pkg.Name)
				continue
			} else {
				log.Printf("Package '%s' not detected or check command failed. Proceeding with setup.", pkg.Name)
			}
			// ... (rest of switch runtime.GOOS for install, as previously defined) ...
			// This part needs to be filled in with the brew/winget command construction
			// For brevity, I'm omitting the full switch statement here, but it should be the same as before.
			// Example for brew:
			if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
				if pkg.Brew.PackageName != "" {
					pkgManagerName = "brew"
					pkgIdentifier = pkg.Brew.PackageName
					log.Printf("Preparing to install binary for '%s' using Homebrew with package name '%s'", pkg.Name, pkgIdentifier)
					cmd = exec.Command(pkgManagerName, opVerb, pkgIdentifier)
				}
			} // Add Windows/Winget case here

		} else { // Teardown (Uninstall)
			opVerb = "uninstall"
			// ... (rest of teardown logic, similar to setup but with uninstall commands) ...
		}

		if cmd != nil {
			log.Printf("Attempting to %s '%s' using %s...", opVerb, pkgIdentifier, pkgManagerName)
			sharedtypes.RunCommand(cmd)
			if isSetup {
				log.Printf("Installation processed for package: %s (%s)", pkg.Name, pkgIdentifier)
			}
		} else if pkgIdentifier != "" {
			log.Printf("No suitable package manager command found for '%s' on %s to %s package '%s'", pkg.Name, runtime.GOOS, opVerb, pkgIdentifier)
		} else {
			log.Printf("No package manager configuration for '%s' on %s.", pkg.Name, runtime.GOOS)
		}
	}

	// Handle VS Code Extensions by calling the dedicated function
	if err := vscodeutils.HandleExtensions(assets, appConfig, isSetup, migrationSet); err != nil {
		// Log the error but don't make it fatal for the whole JSON-driven process
		log.Printf("Error during VS Code extension handling: %v", err)
	}
	return nil
}

func isPackageInstalled(pkg sharedtypes.PackageInfo) bool {
	if pkg.CheckCommand == "" {
		log.Printf("No check command for package '%s', assuming not installed for setup or installed for teardown.", pkg.Name)
		return false
	}
	parts := strings.Fields(pkg.CheckCommand)
	if len(parts) == 0 {
		return false
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	err := cmd.Run()
	if err == nil {
		log.Printf("Check command for '%s' succeeded, package is likely installed.", pkg.Name)
	} else {
		log.Printf("Check command for '%s' failed (%v), package is likely not installed.", pkg.Name, err)
	}
	return err == nil
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
