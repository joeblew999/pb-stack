package main

import (
	"embed"
	"flag"
	"fmt"
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
	flag.Parse()

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
		cli.Execute(embeddedAssets, *bootFlag, *debootFlag, *targetHostFlag, *packageNameFlag)
	} else {
		// --- GUI Mode (Default) ---
		// If -boot or -deboot flags were passed without -cli, they are effectively ignored here,
		// which matches the original behavior where these flags were only checked within the cliModeFlag block.

		fmt.Println("Launching GUI application...") // gui.Launch() will handle its own logging if needed
		gui.Launch()                                // Call the Launch function from the gui package
		fmt.Println("GUI application closed.")
	}
}
