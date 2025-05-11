package gui

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hajimehoshi/guigui/basicwidget"
)

// runCLIProcess is a helper function to execute the main program with CLI flags.
// It runs the command in a new goroutine to avoid blocking the GUI.
// It now accepts an action string (e.g., "Booting", "Debooting"), the CLI flag,
// the packageName string, the migrationSet string, and the statusText widget to provide feedback.
func runCLIProcess(actionName string, cliActionFlag string, packageName string, migrationSet string, statusText *basicwidget.Text) {
	exePath, err := os.Executable()
	if err != nil {
		errMsg := fmt.Sprintf("GUI Error: Failed to get executable path: %v", err)
		log.Println(errMsg)
		statusText.SetValue("Error: Cannot find executable.") // User-friendly message
		return
	}

	actionLog := fmt.Sprintf("%s system", actionName)
	if packageName != "" {
		actionLog = fmt.Sprintf("%s package '%s'", actionName, packageName)
	}
	if migrationSet != "" && migrationSet != "main" { // Only add if not default or empty
		actionLog = fmt.Sprintf("%s (using set: %s)", actionLog, migrationSet)
	}

	statusText.SetValue(fmt.Sprintf("%s... See console.", actionLog))
	log.Printf("GUI: %s...", actionLog)

	args := []string{"-cli", cliActionFlag}
	if packageName != "" {
		args = append(args, "-package", packageName)
	}
	// Use the provided migrationSet, defaulting to "main" if empty (though GUI now defaults to "main")
	currentMigrationSet := migrationSet
	if currentMigrationSet == "" {
		currentMigrationSet = "main"
	}
	args = append(args, "-migrationSet", currentMigrationSet)
	cmd := exec.Command(exePath, args...)

	// We will capture the output instead of sending it directly to os.Stdout/os.Stderr
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	log.Printf("GUI: Running command: %s %s %s", exePath, args[0], args[1])

	// Capture combined output (stdout and stderr)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		errMsg := fmt.Sprintf("GUI: Command '%s %s' failed: %v", args[0], args[1], err)
		log.Println(errMsg)
		statusText.SetValue(fmt.Sprintf("%s failed.\nError: %v\nOutput:\n%s", actionLog, err, outputStr))
	} else {
		log.Printf("GUI: Command '%s %s' finished successfully. Output below:\n%s", args[0], args[1], outputStr)
		statusText.SetValue(fmt.Sprintf("%s successful.\nOutput:\n%s", actionLog, outputStr))
	}
}
