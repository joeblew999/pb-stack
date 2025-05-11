package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	// Assuming your module path allows these imports.
	// If your go.mod is in the 'boot' directory, these might be "myprogram/cmd/cli"
	"main/cmd/cli"
	"main/cmd/gui" // Import the gui package
)

//go:embed all:migrations
var embeddedAssets embed.FS // Embed migrations folder

func main() {
	debootFlag := flag.Bool("deboot", false, "Run deboot scripts instead of boot scripts")
	bootFlag := flag.Bool("boot", false, "Run boot scripts (used with -cli)")
	cliModeFlag := flag.Bool("cli", false, "Run command-line boot/deboot process instead of GUI")
	targetHostFlag := flag.String("target", "", "Target host/IP for boot/deboot operations (used with -cli)")
	packageNameFlag := flag.String("package", "", "Specific package name for boot/deboot (e.g., Winget ID or Homebrew formula)")
	logFileFlag := flag.String("logFile", "", "Path to log file. If empty, logs to stderr only.")
	debugFlag := flag.Bool("debug", false, "Enable debug logging.")
	flag.Parse()

	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Standard flags with file/line number

	if *logFileFlag != "" {
		logFile, err := os.OpenFile(*logFileFlag, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatalf("Failed to open log file %s: %v", *logFileFlag, err)
		}
		defer logFile.Close()
		// Log to both stderr and the file
		multiWriter := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(multiWriter)
		log.Printf("Logging to console and file: %s", *logFileFlag)
	} else {
		log.Printf("Logging to console only.")
	}

	if *cliModeFlag {
		// --- CLI Mode ---
		//var action string
		//var scriptBaseName string

		if *bootFlag && *debootFlag {
			fmt.Fprintln(os.Stderr, "Error: -boot and -deboot flags are mutually exclusive when using -cli.")
			os.Exit(1)
		}
		// The cli.Execute function will handle its specific logic,
		// including the case where neither -boot nor -deboot is specified.
		cli.Execute(embeddedAssets, *bootFlag, *debootFlag, *targetHostFlag, *packageNameFlag, *debugFlag)
	} else {
		// --- GUI Mode (Default) ---
		// If -boot or -deboot flags were passed without -cli, they are effectively ignored here,
		// which matches the original behavior where these flags were only checked within the cliModeFlag block.

		log.Println("Launching GUI application...")
		gui.Launch() // Call the Launch function from the gui package
		log.Println("GUI application closed.")
	}
}
