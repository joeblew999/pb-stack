package main

import (
	"embed"
	"flag"
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
	log.Println("Application starting...")

	teardownFlag := flag.Bool("teardown", false, "Run teardown scripts/operations (used with -cli)")
	setupFlag := flag.Bool("setup", false, "Run setup scripts/operations (used with -cli)")
	cliModeFlag := flag.Bool("cli", false, "Run command-line boot/deboot process instead of GUI")
	packageNameFlag := flag.String("package", "", "Specific package name for boot/deboot (e.g., Winget ID or Homebrew formula)")
	logFileFlag := flag.String("logFile", "", "Path to log file. If empty, logs to stderr only.")
	migrationSetFlag := flag.String("migrationSet", "main", "The set of migrations to use (e.g., 'main', 'test') from the migrations folder.")
	inspectConfigFlag := flag.Bool("inspect-config", false, "Inspect the config.json for the selected migration set and exit.")
	debugFlag := flag.Bool("debug", false, "Enable debug logging.")

	// Flags for the GitHub asset finder command (config file mode)
	assetConfigFlag := flag.String("asset-config", "", "Path to a YAML file with asset search configurations.")
	assetTokenFlag := flag.String("asset-token", os.Getenv("GITHUB_TOKEN"), "GitHub API token (for -find-asset, defaults to GITHUB_TOKEN env var).")

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

	if *debugFlag {
		log.Println("Debug mode enabled.")
		log.Printf("Parsed flags: -cli=%t, -setup=%t, -teardown=%t, -package='%s', -logFile='%s', -migrationSet='%s', -inspect-config=%t, -debug=%t",
			*cliModeFlag, *setupFlag, *teardownFlag, *packageNameFlag, *logFileFlag, *migrationSetFlag, *inspectConfigFlag, *debugFlag)
		if *assetConfigFlag != "" {
			log.Printf("Asset Finder flags: -asset-config='%s', -asset-token set: %t", *assetConfigFlag, *assetTokenFlag != "")
		}
	}

	if *assetConfigFlag != "" {
		// --- GitHub Asset Finder Mode (Config File) ---
		log.Println("Entering GitHub Asset Finder mode (using config file).")
		if *cliModeFlag || *setupFlag || *teardownFlag {
			log.Fatalln("Error: -asset-config mode is mutually exclusive with -cli, -setup, or -teardown flags.")
		}
		cli.FindAssetsFromConfig(*assetConfigFlag, *assetTokenFlag)
		log.Println("GitHub Asset Finder (config file) mode complete.")
	} else if *cliModeFlag {

		// --- CLI Mode ---
		log.Println("Entering CLI mode.")

		if *setupFlag && *teardownFlag {
			// Log this error as well, before exiting
			errMsg := "Error: -setup and -teardown flags are mutually exclusive when using -cli."
			log.Fatalln(errMsg) // Fatalln will print to log and exit
		}
		// The cli.Execute function will handle its specific logic,
		// including the case where neither -boot nor -deboot is specified.
		log.Println("Calling cli.Execute...")
		cli.Execute(embeddedAssets, *setupFlag, *teardownFlag, *packageNameFlag, *migrationSetFlag, *inspectConfigFlag, *debugFlag)
		log.Println("CLI mode complete.")
	} else {
		// --- GUI Mode (Default) ---
		// If -boot or -deboot flags were passed without -cli, they are effectively ignored here,
		// which matches the original behavior where these flags were only checked within the cliModeFlag block.
		log.Println("Entering GUI mode.")

		log.Println("Calling gui.Launch...")
		gui.Launch(*migrationSetFlag) // Pass the migrationSet to the GUI
		log.Println("GUI application closed.")
	}
}
