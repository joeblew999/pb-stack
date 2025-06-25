package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"decksh-ebdeck/pkg/bridge"
	"decksh-ebdeck/pkg/gui"
	"decksh-ebdeck/pkg/renderer"
)

var (
	inputFile  = flag.String("input", "", "Input decksh file (.dsh)")
	outputFile = flag.String("output", "", "Output file (optional, for saving rendered content)")
	guiMode    = flag.Bool("gui", false, "Launch GUI presentation viewer")
	width      = flag.Int("width", 1200, "Window width")
	height     = flag.Int("height", 800, "Window height")
	verbose    = flag.Bool("verbose", false, "Verbose output")
)

func main() {
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -input <file.dsh> [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate input file
	if !strings.HasSuffix(*inputFile, ".dsh") {
		log.Fatalf("Input file must have .dsh extension: %s", *inputFile)
	}

	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		log.Fatalf("Input file does not exist: %s", *inputFile)
	}

	// Process decksh file
	if *verbose {
		fmt.Printf("Processing decksh file: %s\n", *inputFile)
	}

	deckXML, err := bridge.ProcessDeckshFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to process decksh file: %v", err)
	}

	if *verbose {
		fmt.Printf("Generated deck XML (%d bytes)\n", len(deckXML))
	}

	// Parse deck XML into presentation data
	presentation, err := bridge.ParseDeckXML(deckXML)
	if err != nil {
		log.Fatalf("Failed to parse deck XML: %v", err)
	}

	if *verbose {
		fmt.Printf("Parsed presentation with %d slides\n", len(presentation.Slides))
	}

	// Launch GUI or render to file
	if *guiMode {
		if *verbose {
			fmt.Println("Launching GUI presentation viewer...")
		}

		err = gui.LaunchPresentation(presentation, *width, *height)
		if err != nil {
			log.Fatalf("Failed to launch GUI: %v", err)
		}
	} else {
		// Render to file or display info
		if *outputFile != "" {
			err = renderer.RenderToFile(presentation, *outputFile, *width, *height)
			if err != nil {
				log.Fatalf("Failed to render to file: %v", err)
			}
			fmt.Printf("Rendered presentation to: %s\n", *outputFile)
		} else {
			// Just display presentation info
			fmt.Printf("Presentation: %s\n", filepath.Base(*inputFile))
			fmt.Printf("Slides: %d\n", len(presentation.Slides))
			fmt.Println("Use -gui flag to view presentation or -output to render to file")
		}
	}
}
