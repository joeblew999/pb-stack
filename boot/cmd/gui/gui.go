package gui

import (
	"image" // Added for placeholder view backgrounds
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"         // Import Ebiten
	"github.com/hajimehoshi/ebiten/v2/text/v2" // For font handling
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

// HomeViewWidget will contain the original UI elements of the application
type HomeViewWidget struct {
	guigui.DefaultWidget
	background          basicwidget.Background
	packageLabel        basicwidget.Text
	packageInput        basicwidget.TextInput
	migrationSetLabel   basicwidget.Text
	migrationSetInput   basicwidget.TextInput
	bootButton          basicwidget.TextButton
	initialMigrationSet string // To store the value passed from main.go
	debootButton        basicwidget.TextButton
	clearButton         basicwidget.TextButton
	statusText          basicwidget.Text
}

// Root is the main container widget for the application, holding sidebar and views
type Root struct {
	guigui.DefaultWidget

	background basicwidget.Background
	sidebar    SidebarWidget // From sidebar.go
	model      *Model        // From model.go

	// Views
	homeView     HomeViewWidget     // The main view for setup/teardown
	packagesView PackagesViewWidget // View for listing packages
	settingsView SettingsViewWidget // Placeholder
}

// Build for HomeViewWidget (contains the original Root.Build logic)
func (h *HomeViewWidget) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	// h.background.SetColor(color.NRGBA{R: 0x1e, G: 0x20, B: 0x26, A: 0xff}) // Dark background
	appender.AppendChildWidgetWithBounds(&h.background, context.Bounds(h))

	h.packageLabel.SetValue("Package Name (Winget/Homebrew):")

	h.migrationSetLabel.SetValue("Migration Set (e.g., main, test):")
	if h.initialMigrationSet != "" && h.migrationSetInput.Value() == "" {
		h.migrationSetInput.SetValue(h.initialMigrationSet)
	} else if h.migrationSetInput.Value() == "" {
		h.migrationSetInput.SetValue("main")
	}

	h.statusText.SetSelectable(true)
	h.statusText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter) // Use basicwidget package
	h.statusText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)     // Use basicwidget package
	h.statusText.SetScale(1)
	if h.statusText.Value() == "" {
		h.statusText.SetValue("Select an action.")
	}

	h.bootButton.SetText("Setup")
	h.bootButton.SetOnUp(func() {
		go runCLIProcess("Setting up", "-setup", h.packageInput.Value(), h.migrationSetInput.Value(), &h.statusText)
	})

	h.debootButton.SetText("Teardown")
	h.debootButton.SetOnUp(func() {
		go runCLIProcess("Tearing down", "-teardown", h.packageInput.Value(), h.migrationSetInput.Value(), &h.statusText)
	})

	h.clearButton.SetText("Clear")
	h.clearButton.SetOnUp(func() {
		h.packageInput.SetValue("")
		h.migrationSetInput.SetValue("main")
		h.statusText.SetValue("Select an action.")
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(h).Inset(u),
		Heights: []layout.Size{
			layout.FixedSize(u),
			layout.FixedSize(u * 2),
			layout.FixedSize(u),
			layout.FixedSize(u * 2),
			layout.FixedSize(u * 2),
			layout.FixedSize(u * 2),
			layout.FixedSize(u * 2),
			layout.FlexibleSize(1),
		},
		RowGap: u,
	}
	appender.AppendChildWidgetWithBounds(&h.packageLabel, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&h.packageInput, gl.CellBounds(0, 1))
	appender.AppendChildWidgetWithBounds(&h.migrationSetLabel, gl.CellBounds(0, 2))
	appender.AppendChildWidgetWithBounds(&h.migrationSetInput, gl.CellBounds(0, 3))
	appender.AppendChildWidgetWithBounds(&h.bootButton, gl.CellBounds(0, 4))
	appender.AppendChildWidgetWithBounds(&h.debootButton, gl.CellBounds(0, 5))
	appender.AppendChildWidgetWithBounds(&h.clearButton, gl.CellBounds(0, 6))
	appender.AppendChildWidgetWithBounds(&h.statusText, gl.CellBounds(0, 7))

	return nil
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	r.sidebar.SetModel(r.model)

	unitSize := basicwidget.UnitSize(context)
	if unitSize == 0 { // Should not happen if fonts are loaded
		log.Println("Warning: UnitSize is 0 in Root.Build. Sidebar might not render correctly.")
		unitSize = 16 // Fallback, though font loading should prevent this
	}

	gl := layout.GridLayout{
		Bounds: context.Bounds(r),
		Widths: []layout.Size{
			layout.FixedSize(8 * unitSize), // Sidebar width
			layout.FlexibleSize(1),         // Content area
		},
		Heights: []layout.Size{layout.FlexibleSize(1)}, // Ensure full height
	}
	appender.AppendChildWidgetWithBounds(&r.sidebar, gl.CellBounds(0, 0))

	contentBounds := gl.CellBounds(1, 0)
	switch r.model.Mode() {
	case "home":
		appender.AppendChildWidgetWithBounds(&r.homeView, contentBounds)
	case "packages":
		appender.AppendChildWidgetWithBounds(&r.packagesView, contentBounds)
	case "settings":
		appender.AppendChildWidgetWithBounds(&r.settingsView, contentBounds)
	default:
		// Fallback to home view if mode is unknown
		appender.AppendChildWidgetWithBounds(&r.homeView, contentBounds)
	}

	return nil
}

// PackagesViewWidget displays information about packages, loaded from the CLI.
type PackagesViewWidget struct {
	guigui.DefaultWidget
	statusText    basicwidget.Text // To display CLI output (package list or errors)
	bg            basicwidget.Background
	loadAttempted bool // To ensure CLI is called only once when view becomes active
}

func (pv *PackagesViewWidget) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	appender.AppendChildWidgetWithBounds(&pv.bg, context.Bounds(pv))

	if !pv.loadAttempted {
		pv.loadAttempted = true // Mark that we are attempting to load
		// Set initial status
		pv.statusText.SetValue("Loading package information...")
		// Use the -inspect-config flag from the CLI to get package information.
		// The migrationSet defaults to "main" in runCLIProcess if not specified.
		go runCLIProcess("Inspecting packages", "-inspect-config", "", "", &pv.statusText)
	}

	// Configure statusText properties
	pv.statusText.SetSelectable(true)
	pv.statusText.SetHorizontalAlign(basicwidget.HorizontalAlignStart) // Use basicwidget package (Start usually means Left)
	pv.statusText.SetVerticalAlign(basicwidget.VerticalAlignTop)       // Use basicwidget package
	pv.statusText.SetScale(1)                                          // Default text scale

	// Layout the statusText to fill most of the view, with some padding
	unitSize := basicwidget.UnitSize(context)
	contentBounds := context.Bounds(pv).Inset(unitSize) // Add padding
	appender.AppendChildWidgetWithBounds(&pv.statusText, contentBounds)

	return nil
}

// SettingsViewWidget - Placeholder
type SettingsViewWidget struct {
	guigui.DefaultWidget
	label basicwidget.Text
	bg    basicwidget.Background
}

func (sv *SettingsViewWidget) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	// The SetColor method is not available on basicwidget.Background.
	// This will now use the default theme background color.
	appender.AppendChildWidgetWithBounds(&sv.bg, context.Bounds(sv))
	sv.label.SetValue("Application Settings View.")
	labelSize := sv.label.DefaultSize(context)
	bounds := context.Bounds(sv)
	pos := image.Pt(bounds.Min.X+(bounds.Dx()-labelSize.X)/2, bounds.Min.Y+(bounds.Dy()-labelSize.Y)/2)
	appender.AppendChildWidgetWithPosition(&sv.label, pos)
	return nil
}

// Launch is the public entry point to start the GUI.
func Launch(initialMigrationSet string) {
	// Set up the default font face sources for guigui widgets
	// This is similar to how the gallery example does it.
	faceSources := []*text.GoTextFaceSource{
		basicwidget.DefaultFaceSource(),
		// You can add other font sources here if needed, e.g., for CJK characters
		// basicwidget.cjkfont.FaceSourceJP(), // Example from gallery
	}
	basicwidget.SetFaceSources(faceSources)

	// Get current screen dimensions to set the initial window size
	// Use ScreenSizeInFullscreen() for broader compatibility
	screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()

	appModel := &Model{} // Create the application model

	rootWidget := &Root{
		model: appModel,
		homeView: HomeViewWidget{ // Initialize HomeViewWidget with its specific needs
			initialMigrationSet: initialMigrationSet,
		},
		packagesView: PackagesViewWidget{
			// loadAttempted will be false by default.
			// statusText (basicwidget.Text) will be zero-initialized and ready to use.
		},
	}

	op := &guigui.RunOptions{
		Title:         "PB-Stack Environment Setup",
		WindowSize:    image.Pt(screenWidth, screenHeight), // Set initial size to screen dimensions
		WindowMinSize: image.Pt(320, 240),                  // Still allow it to be shrunk
		// WindowMaxSize: image.Pt(640, 480), // Remove or comment out to allow full maximization
	}
	// Allow the window to be resized by the user. This should be called before Run.
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := guigui.Run(rootWidget, op); err != nil {
		log.Printf("Error running GUI: %v", err) // Changed to log
		os.Exit(1)
	}
}
