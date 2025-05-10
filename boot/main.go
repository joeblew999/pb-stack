package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec" // Import the filepath package
	"runtime" // Import the runtime package
)

//go:embed all:migrations
var embeddedAssets embed.FS // Embed migrations folder

func main() {
	debootFlag := flag.Bool("deboot", false, "Run deboot scripts instead of boot scripts")
	bootFlag := flag.Bool("boot", false, "Run boot scripts (used with -cli)")
	cliModeFlag := flag.Bool("cli", false, "Run command-line boot/deboot process instead of GUI")
	flag.Parse()

	if *cliModeFlag {
		// --- CLI Mode ---
		var action string
		var scriptBaseName string

		if *bootFlag && *debootFlag {
			fmt.Fprintln(os.Stderr, "Error: -boot and -deboot flags are mutually exclusive when using -cli.")
			os.Exit(1)
		}

		if *bootFlag {
			action = "Booting (CLI)"
			scriptBaseName = "boot"
		} else if *debootFlag {
			action = "Debooting (CLI)"
			scriptBaseName = "deboot"
		} else {
			fmt.Println("CLI mode selected. Please specify an action:")
			fmt.Println("  go run ./main.go -cli -boot    (to run boot scripts)")
			fmt.Println("  go run ./main.go -cli -deboot  (to run deboot scripts)")
			os.Exit(0) // Exit cleanly after showing options
		}
		fmt.Printf("%s up...\n", action)

		// Example: List files in the embedded migrations directory
		fs.WalkDir(embeddedAssets, "migrations", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(path)
			return nil
		})

		fmt.Println("Checking OS ...")
		fmt.Printf("Detected OS: %s\n", runtime.GOOS)

		var scriptName string
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "darwin", "linux":
			{
				scriptName = fmt.Sprintf("migrations/%s.sh", scriptBaseName) // Path within the embedded FS
				scriptBytes, err := embeddedAssets.ReadFile(scriptName)
				if err != nil {
					log.Fatalf("Failed to read embedded script %s: %v", scriptName, err)
				}

				tempFile, err := os.CreateTemp(os.TempDir(), "pb-stack-boot-*.sh")
				if err != nil {
					log.Fatalf("Failed to create temp file for %s: %v", scriptName, err)
				}
				defer os.Remove(tempFile.Name()) // Clean up the temp file

				if _, err := tempFile.Write(scriptBytes); err != nil {
					tempFile.Close() // Close before attempting remove on error
					log.Fatalf("Failed to write to temp file for %s: %v", scriptName, err)
				}
				// Make the script executable
				if err := tempFile.Chmod(0755); err != nil {
					tempFile.Close()
					log.Fatalf("Failed to set executable permission for temp file %s: %v", scriptName, err)
				}

				tempFilePath := tempFile.Name()
				if err := tempFile.Close(); err != nil { // Close the file before executing
					log.Fatalf("Failed to close temp file for %s: %v", scriptName, err)
				}

				cmd = exec.Command(tempFilePath)
			}
		case "windows":
			{
				scriptName = fmt.Sprintf("migrations/%s.ps1", scriptBaseName) // Path within the embedded FS
				scriptBytes, err := embeddedAssets.ReadFile(scriptName)
				if err != nil {
					log.Fatalf("Failed to read embedded script %s: %v", scriptName, err)
				}

				tempFile, err := os.CreateTemp(os.TempDir(), "pb-stack-boot-*.ps1")
				if err != nil {
					log.Fatalf("Failed to create temp file for %s: %v", scriptName, err)
				}
				defer os.Remove(tempFile.Name()) // Clean up the temp file

				if _, err := tempFile.Write(scriptBytes); err != nil {
					tempFile.Close() // Close before attempting remove on error
					log.Fatalf("Failed to write to temp file for %s: %v", scriptName, err)
				}

				tempFilePath := tempFile.Name()
				if err := tempFile.Close(); err != nil { // Close the file before passing its name to PowerShell
					log.Fatalf("Failed to close temp file for %s: %v", scriptName, err)
				}

				cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempFilePath)
			}
		default:
			log.Fatalf("Unsupported OS for CLI mode: %s", runtime.GOOS)
		}

		if cmd != nil {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = os.TempDir()
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to execute %s: %v", scriptName, err)
			}
		}
		fmt.Printf("%s up complete.\n", action)
	} else {
		// --- GUI Mode (Default) ---
		fmt.Println("Launching GUI application...")
		launchEbitenGUI() // Call your Ebiten launch function here
		fmt.Println("GUI application closed.")
	}
}

// launchEbitenGUI is a placeholder for your Ebiten application launch.
func launchEbitenGUI() {
	fmt.Println("*************************************")
	fmt.Println("*    Placeholder for Ebiten GUI     *")
	fmt.Println("* This would launch your GUI app.   *")
	fmt.Println("*************************************")
	// Example Ebiten launch (you'll need to import "github.com/hajimehoshi/ebiten/v2"):
	//
	// type Game struct{}
	// func (g *Game) Update() error { return nil }
	// func (g *Game) Draw(screen *ebiten.Image) { screen.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff}) }
	// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) { return 320, 240 }
	//
	// ebiten.SetWindowSize(640, 480)
	// ebiten.SetWindowTitle("My Ebiten App")
	// if err := ebiten.RunGame(&Game{}); err != nil {
	// 	log.Fatal(err)
	// }
}
