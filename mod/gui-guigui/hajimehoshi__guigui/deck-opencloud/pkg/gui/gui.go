package gui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
)

// App represents the main GUI application
type App struct {
	guigui.DefaultWidget
	background basicwidget.Background
	title      basicwidget.Text
	searchBox  basicwidget.TextInput
	searchBtn  basicwidget.Button
	statusText basicwidget.Text
}

// Build implements the guigui.Widget interface
func (a *App) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	// Set up the background
	appender.AppendChildWidgetWithBounds(&a.background, context.Bounds(a))

	// Create a simple grid layout for the main content
	bounds := context.Bounds(a).Inset(20)

	// Calculate positions manually since VBoxLayout doesn't exist
	titleHeight := 40
	inputHeight := 30
	buttonHeight := 30
	spacing := 10

	currentY := bounds.Min.Y

	// Title
	a.title.SetValue("OpenCloud - Search Interface")
	a.title.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	titleBounds := image.Rect(bounds.Min.X, currentY, bounds.Max.X, currentY+titleHeight)
	appender.AppendChildWidgetWithBounds(&a.title, titleBounds)
	currentY += titleHeight + spacing

	// Search input
	a.searchBox.SetEditable(true)
	inputBounds := image.Rect(bounds.Min.X, currentY, bounds.Max.X, currentY+inputHeight)
	appender.AppendChildWidgetWithBounds(&a.searchBox, inputBounds)
	currentY += inputHeight + spacing

	// Search button
	a.searchBtn.SetText("Search")
	a.searchBtn.SetOnDown(func() {
		query := a.searchBox.Value()
		if query != "" {
			a.statusText.SetValue("Searching for: " + query)
		} else {
			a.statusText.SetValue("Please enter a search query")
		}
	})
	buttonBounds := image.Rect(bounds.Min.X, currentY, bounds.Max.X, currentY+buttonHeight)
	appender.AppendChildWidgetWithBounds(&a.searchBtn, buttonBounds)
	currentY += buttonHeight + spacing

	// Status text
	if a.statusText.Value() == "" {
		a.statusText.SetValue("Ready to search...")
	}
	statusBounds := image.Rect(bounds.Min.X, currentY, bounds.Max.X, currentY+30)
	appender.AppendChildWidgetWithBounds(&a.statusText, statusBounds)

	return nil
}

// Start initializes and starts the GUI
func Start() error {
	app := &App{}

	options := &guigui.RunOptions{
		Title:      "OpenCloud GUI",
		WindowSize: image.Pt(800, 600),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}

	return guigui.Run(app, options)
}
