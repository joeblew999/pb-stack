package vscodeutils

import (
	"embed"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"main/pkg/sharedtypes" // Import shared types
)

// HandleExtensions manages the setup or teardown of VS Code extensions.
func HandleExtensions(assets embed.FS, appConfig sharedtypes.AppConfig, isSetup bool, migrationSet string) error {
	if appConfig.VSCodeExtensionsFile == "" {
		log.Println("No VS Code extensions file specified in config.json. Skipping VS Code extension processing.")
		return nil
	}

	log.Printf("Processing VS Code extensions using file: %s", appConfig.VSCodeExtensionsFile)
	extensionsFilePathInAssets := filepath.Join("migrations", migrationSet, appConfig.VSCodeExtensionsFile)
	extensionsBytes, err := assets.ReadFile(extensionsFilePathInAssets)
	if err != nil {
		log.Printf("Warning: Could not read embedded extensions file '%s': %v. Skipping VS Code extension processing.", extensionsFilePathInAssets, err)
		return nil // Not a fatal error for the whole process, just skip this part
	}

	desiredExtensions := strings.Split(strings.ReplaceAll(string(extensionsBytes), "\r\n", "\n"), "\n")
	var validDesiredExtensions []string
	for _, ext := range desiredExtensions {
		trimmedExt := strings.TrimSpace(ext)
		if trimmedExt != "" && !strings.HasPrefix(trimmedExt, "#") {
			validDesiredExtensions = append(validDesiredExtensions, trimmedExt)
		}
	}

	if _, err := exec.LookPath("code"); err != nil {
		log.Println("VS Code 'code' command not found in PATH. Skipping extension management.")
		return nil // Not fatal, just can't manage extensions.
	}

	if isSetup {
		log.Println("Setting up VS Code extensions...")
		for _, ext := range validDesiredExtensions {
			log.Printf("Checking/Installing VS Code extension: %s", ext)
			cmd := exec.Command("code", "--install-extension", ext)
			sharedtypes.RunCommand(cmd) // RunCommand logs errors but doesn't make them fatal for other extensions
		}
	} else { // Teardown
		log.Println("Tearing down VS Code extensions (uninstalling those not in the desired list)...")
		installedExtBytes, err := exec.Command("code", "--list-extensions").Output()
		if err != nil {
			log.Printf("Warning: Could not list installed VS Code extensions: %v", err)
		} else {
			installedExtensions := strings.Split(strings.ReplaceAll(string(installedExtBytes), "\r\n", "\n"), "\n")
			desiredMap := make(map[string]bool)
			for _, de := range validDesiredExtensions {
				desiredMap[de] = true
			}

			for _, instExt := range installedExtensions {
				trimmedInstExt := strings.TrimSpace(instExt)
				if trimmedInstExt != "" && !desiredMap[trimmedInstExt] {
					log.Printf("Uninstalling VS Code extension not in desired list: %s", trimmedInstExt)
					cmd := exec.Command("code", "--uninstall-extension", trimmedInstExt)
					sharedtypes.RunCommand(cmd) // RunCommand logs errors
				}
			}
		}
	}
	return nil
}
