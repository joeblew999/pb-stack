package main

import (
	"fmt"
	"image/color"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	cardWidth    = 120
	cardHeight   = 168
)

// Card represents a playing card
type Card struct {
	Suit     string
	Value    string
	Image    *ebiten.Image
	X, Y     float64
	Selected bool
}

// Game represents the main game state
type Game struct {
	cards       []*Card
	cardImages  map[string]*ebiten.Image
	initialized bool
}

func NewGame() *Game {
	return &Game{
		cards:      make([]*Card, 0),
		cardImages: make(map[string]*ebiten.Image),
	}
}

// LoadCardImages loads all card images from the assets directory
func (g *Game) LoadCardImages() error {
	suits := []string{"hearts", "diamonds", "clubs", "spades"}
	values := []string{"a", "2", "3", "4", "5", "6", "7", "8", "9", "10", "j", "q", "k"}

	for _, suit := range suits {
		for _, value := range values {
			filename := fmt.Sprintf("%s_%s.png", suit, value)
			imagePath := filepath.Join(".assets", "png", filename)

			img, _, err := ebitenutil.NewImageFromFile(imagePath)
			if err != nil {
				// If image doesn't exist, create a placeholder
				img = ebiten.NewImage(cardWidth, cardHeight)
				img.Fill(color.White)
				log.Printf("Warning: Could not load %s, using placeholder", imagePath)
			}

			g.cardImages[fmt.Sprintf("%s_%s", suit, value)] = img
		}
	}

	return nil
}

// InitializeGame sets up the initial game state
func (g *Game) InitializeGame() error {
	if g.initialized {
		return nil
	}

	// Load card images
	if err := g.LoadCardImages(); err != nil {
		return err
	}

	// Create a sample hand of cards for demonstration
	sampleCards := []struct {
		suit  string
		value string
		x, y  float64
	}{
		{"hearts", "a", 50, 200},
		{"diamonds", "k", 180, 200},
		{"clubs", "q", 310, 200},
		{"spades", "j", 440, 200},
		{"hearts", "10", 570, 200},
	}

	for _, cardData := range sampleCards {
		key := fmt.Sprintf("%s_%s", cardData.suit, cardData.value)
		if img, exists := g.cardImages[key]; exists {
			card := &Card{
				Suit:  cardData.suit,
				Value: cardData.value,
				Image: img,
				X:     cardData.x,
				Y:     cardData.y,
			}
			g.cards = append(g.cards, card)
		}
	}

	g.initialized = true
	return nil
}

// Update handles game logic updates
func (g *Game) Update() error {
	if !g.initialized {
		if err := g.InitializeGame(); err != nil {
			return err
		}
	}

	// Handle mouse input for card selection
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		for _, card := range g.cards {
			// Check if mouse is over this card
			if float64(x) >= card.X && float64(x) <= card.X+cardWidth &&
				float64(y) >= card.Y && float64(y) <= card.Y+cardHeight {
				card.Selected = !card.Selected
			}
		}
	}

	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with green background (like a card table)
	screen.Fill(color.RGBA{0, 100, 0, 255})

	// Draw title
	ebitenutil.DebugPrint(screen, "Ebiten Card Game Demo\nClick cards to select them")

	// Draw cards
	for _, card := range g.cards {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(card.X, card.Y)

		// If selected, add a yellow tint
		if card.Selected {
			op.ColorM.Scale(1.2, 1.2, 0.8, 1.0)
		}

		screen.DrawImage(card.Image, op)

		// Draw card info
		info := fmt.Sprintf("%s of %s", card.Value, card.Suit)
		ebitenutil.DebugPrintAt(screen, info, int(card.X), int(card.Y+cardHeight+5))
	}

	// Draw instructions
	ebitenutil.DebugPrintAt(screen, "Cards loaded from: ./assets/cards/", 10, screenHeight-40)
	ebitenutil.DebugPrintAt(screen, "Run 'go run ../cardgen' to generate card images", 10, screenHeight-20)
}

// Layout returns the screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten Card Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
