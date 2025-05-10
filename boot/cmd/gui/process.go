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
// and the statusText widget to provide feedback.
func runCLIProcess(actionName string, cliActionFlag string, statusText *basicwidget.Text) {
	exePath, err := os.Executable()
	if err != nil {
		errMsg := fmt.Sprintf("GUI Error: Failed to get executable path: %v", err)
		log.Println(errMsg)
		statusText.SetValue("Error: Cannot find executable.") // User-friendly message
		return
	}

	statusText.SetValue(fmt.Sprintf("%s... See console.", actionName))
	log.Printf("GUI: %s system...", actionName)

	args := []string{"-cli", cliActionFlag}
	cmd := exec.Command(exePath, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("GUI: Running command: %s %s %s", exePath, args[0], args[1])
	err = cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("GUI Error: Command '%s %s' failed: %v", args[0], args[1], err)
		log.Println(errMsg)
		statusText.SetValue(fmt.Sprintf("%s failed. See console.", actionName))
	} else {
		log.Printf("GUI: Command '%s %s' finished successfully.", args[0], args[1])
		statusText.SetValue(fmt.Sprintf("%s successful.", actionName))
	}
}
