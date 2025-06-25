package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ajstarks/decksh"
)

// Card represents a playing card
type Card struct {
	Suit  string // hearts, diamonds, clubs, spades
	Value string // A, 2-10, J, Q, K
}

// DeckshCardGenerator uses decksh to generate card markup
type DeckshCardGenerator struct {
	outputDir  string
	cardWidth  int
	cardHeight int
}

func NewDeckshCardGenerator(outputDir string, width, height int) *DeckshCardGenerator {
	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
	}

	return &DeckshCardGenerator{
		outputDir:  outputDir,
		cardWidth:  width,
		cardHeight: height,
	}
}

// GenerateCardScript creates a decksh script for a single card
func (g *DeckshCardGenerator) GenerateCardScript(card Card) string {
	suitSymbol := getSuitSymbol(card.Suit)
	suitColor := getSuitColor(card.Suit)

	script := fmt.Sprintf(`
// %s %s
deck
	canvas %d %d
	slide "white" "%s"
		// Top-left value and suit
		text "%s" 10 85 8 "sans" "%s"
		text "%s" 10 75 12 "sans" "%s"
		
		// Center suit symbol (large)
		ctext "%s" 50 50 24 "sans" "%s"
		
		// Bottom-right value and suit (rotated)
		// Note: decksh coordinate system has origin at bottom-left
		rtext "%s" 90 15 180 8 "sans" "%s"
		rtext "%s" 90 25 180 12 "sans" "%s"
		
		// Card border
		rect 50 50 95 98 "transparent" 100
		line 2.5 2 97.5 2 0.5 "black"
		line 2.5 98 97.5 98 0.5 "black"
		line 2.5 2 2.5 98 0.5 "black"
		line 97.5 2 97.5 98 0.5 "black"
		
		// Rounded corners (approximated with small arcs)
		arc 8 8 8 8 180 270 0.5 "black"
		arc 92 8 8 8 270 360 0.5 "black"
		arc 8 92 8 8 90 180 0.5 "black"
		arc 92 92 8 8 0 90 0.5 "black"
	eslide
edeck
`, card.Suit, card.Value, g.cardWidth, g.cardHeight,
		suitColor, card.Value, suitColor, suitSymbol, suitColor,
		suitSymbol, suitColor,
		card.Value, suitColor, suitSymbol, suitColor)

	return script
}

// GenerateCard creates a single card using decksh
func (g *DeckshCardGenerator) GenerateCard(card Card) error {
	// Generate decksh script
	script := g.GenerateCardScript(card)

	// Process with decksh to get deck XML
	var xmlOutput strings.Builder
	err := decksh.Process(&xmlOutput, strings.NewReader(script))
	if err != nil {
		return fmt.Errorf("failed to process decksh script: %w", err)
	}

	// Save the XML markup
	filename := fmt.Sprintf("%s_%s.xml", strings.ToLower(card.Suit), strings.ToLower(card.Value))
	xmlPath := filepath.Join(g.outputDir, filename)

	err = os.WriteFile(xmlPath, []byte(xmlOutput.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write XML file: %w", err)
	}

	fmt.Printf("Generated XML: %s\n", xmlPath)
	return nil
}

// GenerateAllCards creates a complete deck
func (g *DeckshCardGenerator) GenerateAllCards() error {
	suits := []string{"hearts", "diamonds", "clubs", "spades"}
	values := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, suit := range suits {
		for _, value := range values {
			card := Card{Suit: suit, Value: value}
			if err := g.GenerateCard(card); err != nil {
				return fmt.Errorf("failed to generate card %s %s: %w", suit, value, err)
			}
		}
	}

	fmt.Printf("Generated %d card XML files in %s\n", len(suits)*len(values), g.outputDir)
	return nil
}

// ConvertToImages uses deck clients to convert XML to images
func (g *DeckshCardGenerator) ConvertToImages(format string) error {
	pattern := filepath.Join(g.outputDir, "*.xml")
	xmlFiles, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to find XML files: %w", err)
	}

	fmt.Printf("Converting %d XML files to %s format...\n", len(xmlFiles), format)

	for _, xmlFile := range xmlFiles {
		baseName := strings.TrimSuffix(filepath.Base(xmlFile), ".xml")

		var cmd string
		var outputFile string

		switch format {
		case "svg":
			cmd = "svgdeck"
			outputFile = filepath.Join(g.outputDir, baseName+".svg")
		case "png":
			cmd = "pngdeck"
			outputFile = filepath.Join(g.outputDir, baseName+".png")
		case "pdf":
			cmd = "pdfdeck"
			outputFile = filepath.Join(g.outputDir, baseName+".pdf")
		default:
			return fmt.Errorf("unsupported format: %s", format)
		}

		fmt.Printf("Converting %s -> %s\n", xmlFile, outputFile)
		// Note: In a real implementation, you would exec.Command() here
		// For now, we just show what commands would be run
		fmt.Printf("  Command: %s %s\n", cmd, xmlFile)
	}

	return nil
}

func getSuitSymbol(suit string) string {
	switch suit {
	case "hearts":
		return "♥"
	case "diamonds":
		return "♦"
	case "clubs":
		return "♣"
	case "spades":
		return "♠"
	default:
		return "?"
	}
}

func getSuitColor(suit string) string {
	if suit == "hearts" || suit == "diamonds" {
		return "red"
	}
	return "black"
}

func main() {
	// Get output directory from command line argument or use default
	outputDir := "./.assets/xml"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	// Create generator
	generator := NewDeckshCardGenerator(outputDir, 360, 504) // 3:4.2 ratio, larger for quality

	fmt.Println("=== Generating cards using decksh ===")
	fmt.Printf("Output directory: %s\n", outputDir)

	// Generate all cards as XML
	if err := generator.GenerateAllCards(); err != nil {
		fmt.Printf("Error generating cards: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== Card XML Generation Complete ===")
	fmt.Println("To convert to images, use the Taskfile:")
	fmt.Println("  task cards:decksh      # Generate XML + SVG")
	fmt.Println("  task cards:decksh-png  # Generate XML + PNG")
	fmt.Println()
	fmt.Println("Or install tools locally and convert manually:")
	fmt.Printf("  ./.bin/svgdeck %s/hearts_a.xml\n", outputDir)
	fmt.Printf("  ./.bin/pngdeck %s/hearts_a.xml\n", outputDir)
}
