package gui

import (
	"fmt"
	"image"
	"image/color"

	"decksh-ebdeck/pkg/bridge"
	"decksh-ebdeck/pkg/renderer"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

// PresentationViewer is the main GUI application for viewing presentations
type PresentationViewer struct {
	guigui.DefaultWidget

	presentation *bridge.Presentation
	currentSlide int
	canvas       *renderer.Canvas
	slideImage   *ebiten.Image

	// UI components
	background   basicwidget.Background
	slideDisplay basicwidget.Image
	prevButton   basicwidget.Button
	nextButton   basicwidget.Button
	slideCounter basicwidget.Text
	titleText    basicwidget.Text
	statusText   basicwidget.Text
}

// NewPresentationViewer creates a new presentation viewer
func NewPresentationViewer(presentation *bridge.Presentation, width, height int) *PresentationViewer {
	viewer := &PresentationViewer{
		presentation: presentation,
		currentSlide: 0,
		canvas:       renderer.NewCanvas(width-100, height-150), // Leave space for UI
	}

	viewer.renderCurrentSlide()
	return viewer
}

// Build implements the guigui.Widget interface
func (pv *PresentationViewer) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	bounds := context.Bounds(pv)

	// Background
	appender.AppendChildWidgetWithBounds(&pv.background, bounds)

	// Title
	pv.titleText.SetValue("Decksh Presentation Viewer")
	pv.titleText.SetBold(true)
	pv.titleText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	pv.titleText.SetScale(2)

	// Slide counter
	slideInfo := fmt.Sprintf("Slide %d of %d", pv.currentSlide+1, len(pv.presentation.Slides))
	pv.slideCounter.SetValue(slideInfo)
	pv.slideCounter.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)

	// Navigation buttons
	pv.prevButton.SetText("Previous")
	pv.prevButton.SetOnUp(func() {
		pv.previousSlide()
	})
	context.SetEnabled(&pv.prevButton, pv.currentSlide > 0)

	pv.nextButton.SetText("Next")
	pv.nextButton.SetOnUp(func() {
		pv.nextSlide()
	})
	context.SetEnabled(&pv.nextButton, pv.currentSlide < len(pv.presentation.Slides)-1)

	// Status text
	pv.statusText.SetValue("Use arrow keys or buttons to navigate. Press ESC to exit.")
	pv.statusText.SetScale(0.8)

	// Layout
	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: bounds.Inset(u),
		Heights: []layout.Size{
			layout.FixedSize(u * 2), // Title
			layout.FlexibleSize(1),  // Slide display
			layout.FixedSize(u),     // Slide counter
			layout.FixedSize(u),     // Navigation buttons
			layout.FixedSize(u / 2), // Status
		},
		RowGap: u / 2,
	}

	// Title
	appender.AppendChildWidgetWithBounds(&pv.titleText, gl.CellBounds(0, 0))

	// Slide display area
	slideArea := gl.CellBounds(0, 1)
	if pv.slideImage != nil {
		// Always update the image in case it changed
		pv.slideDisplay.SetImage(pv.slideImage)
	}
	appender.AppendChildWidgetWithBounds(&pv.slideDisplay, slideArea)

	// Slide counter
	appender.AppendChildWidgetWithBounds(&pv.slideCounter, gl.CellBounds(0, 2))

	// Navigation buttons
	buttonArea := gl.CellBounds(0, 3)
	buttonLayout := layout.GridLayout{
		Bounds: buttonArea,
		Widths: []layout.Size{
			layout.FixedSize(u * 4), // Previous button
			layout.FlexibleSize(1),  // Spacer
			layout.FixedSize(u * 4), // Next button
		},
		ColumnGap: u / 2,
	}
	appender.AppendChildWidgetWithBounds(&pv.prevButton, buttonLayout.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&pv.nextButton, buttonLayout.CellBounds(2, 0))

	// Status
	appender.AppendChildWidgetWithBounds(&pv.statusText, gl.CellBounds(0, 4))

	return nil
}

// HandleButtonInput handles keyboard input
func (pv *PresentationViewer) HandleButtonInput(context *guigui.Context) guigui.HandleInputResult {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		pv.previousSlide()
		return guigui.HandleInputByWidget(pv)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyPageDown) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		pv.nextSlide()
		return guigui.HandleInputByWidget(pv)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		pv.firstSlide()
		return guigui.HandleInputByWidget(pv)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		pv.lastSlide()
		return guigui.HandleInputByWidget(pv)
	}

	return guigui.HandleInputResult{}
}

// previousSlide navigates to the previous slide
func (pv *PresentationViewer) previousSlide() {
	if pv.currentSlide > 0 {
		pv.currentSlide--
		pv.renderCurrentSlide()
		pv.updateUI()
	}
}

// nextSlide navigates to the next slide
func (pv *PresentationViewer) nextSlide() {
	if pv.currentSlide < len(pv.presentation.Slides)-1 {
		pv.currentSlide++
		pv.renderCurrentSlide()
		pv.updateUI()
	}
}

// firstSlide navigates to the first slide
func (pv *PresentationViewer) firstSlide() {
	pv.currentSlide = 0
	pv.renderCurrentSlide()
	pv.updateUI()
}

// lastSlide navigates to the last slide
func (pv *PresentationViewer) lastSlide() {
	pv.currentSlide = len(pv.presentation.Slides) - 1
	pv.renderCurrentSlide()
	pv.updateUI()
}

// renderCurrentSlide renders the current slide to the slide image
func (pv *PresentationViewer) renderCurrentSlide() {
	if pv.currentSlide >= 0 && pv.currentSlide < len(pv.presentation.Slides) {
		slide := pv.presentation.Slides[pv.currentSlide]

		// Clear and render the slide
		pv.canvas.Clear(color.NRGBA{255, 255, 255, 255})
		err := pv.canvas.RenderSlide(slide)
		if err != nil {
			fmt.Printf("Error rendering slide %d: %v\n", pv.currentSlide, err)
		}

		// Update the slide image (Build method will call SetImage)
		pv.slideImage = pv.canvas.Image
	}
}

// updateUI updates the slide counter and button states
func (pv *PresentationViewer) updateUI() {
	// Update slide counter
	slideInfo := fmt.Sprintf("Slide %d of %d", pv.currentSlide+1, len(pv.presentation.Slides))
	pv.slideCounter.SetValue(slideInfo)
}

// LaunchPresentation launches the GUI presentation viewer
func LaunchPresentation(presentation *bridge.Presentation, width, height int) error {
	if len(presentation.Slides) == 0 {
		return fmt.Errorf("presentation has no slides")
	}

	viewer := NewPresentationViewer(presentation, width, height)

	options := &guigui.RunOptions{
		Title:         "Decksh Presentation Viewer",
		WindowMinSize: image.Pt(800, 600),
		WindowSize:    image.Pt(width, height),
	}

	return guigui.Run(viewer, options)
}
