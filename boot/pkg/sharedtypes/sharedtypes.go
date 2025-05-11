package sharedtypes

import (
	"log"
	"os"
	"os/exec"
)

// --- Structs for config.json ---
type BrewConfig struct {
	PackageName string `json:"packageName"`
	IsCask      bool   `json:"isCask,omitempty"`
}

type WingetConfig struct {
	PackageID string `json:"packageId"`
}

type PackageInfo struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Brew         BrewConfig   `json:"brew,omitempty"`
	Winget       WingetConfig `json:"winget,omitempty"`
	CheckCommand string       `json:"checkCommand,omitempty"` // Command to check if installed
}

type AppConfig struct {
	Packages             []PackageInfo `json:"packages"`
	VSCodeExtensionsFile string        `json:"vscodeExtensionsFile"`
}

// --- End Structs ---

func RunCommand(cmd *exec.Cmd) { // Renamed to RunCommand for export
	if cmd == nil {
		log.Println("No command to run.")
		return
	}
	log.Printf("Executing: %s %v", cmd.Path, cmd.Args[1:]) // Args[0] is the command itself
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing command '%s': %v", cmd.Path, err)
	} else {
		log.Printf("Command '%s' completed successfully.", cmd.Path)
	}
}
