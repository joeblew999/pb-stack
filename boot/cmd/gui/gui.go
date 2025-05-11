package gui

import (
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.DefaultWidget

	background   basicwidget.Background
	packageLabel basicwidget.Text
	packageInput basicwidget.TextInput
	bootButton   basicwidget.TextButton
	debootButton basicwidget.TextButton
	clearButton  basicwidget.TextButton
	statusText   basicwidget.Text
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	// r.background.SetColor(color.NRGBA{R: 0x1e, G: 0x20, B: 0x26, A: 0xff}) // Dark background - Method not found in current guigui version
	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	r.packageLabel.SetValue("Package Name (Winget/Homebrew):")
	// r.packageInput.SetPlaceholder("e.g., Git.Git or htop") // Placeholder not available

	r.statusText.SetSelectable(true) // Make status text selectable for copy-pasting
	r.statusText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	r.statusText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	r.statusText.SetScale(1)        // Normal scale for status text
	if r.statusText.Value() == "" { // Set initial value only if not already set
		r.statusText.SetValue("Select an action.")
	}

	r.bootButton.SetText("Setup")
	r.bootButton.SetOnUp(func() {
		// Run in a goroutine to avoid blocking the UI thread
		go runCLIProcess("Setting up", "-boot", r.packageInput.Value(), &r.statusText)
	})

	r.debootButton.SetText("Teardown")
	r.debootButton.SetOnUp(func() {
		// Run in a goroutine to avoid blocking the UI thread
		go runCLIProcess("Tearing down", "-deboot", r.packageInput.Value(), &r.statusText)
	})

	r.clearButton.SetText("Clear")
	r.clearButton.SetOnUp(func() {
		r.packageInput.SetValue("")
		r.statusText.SetValue("Select an action.")
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r).Inset(u),
		Heights: []layout.Size{
			layout.FixedSize(u),     // 0: packageLabel
			layout.FixedSize(u * 2), // 1: packageInput
			layout.FixedSize(u * 2), // 2: bootButton
			layout.FixedSize(u * 2), // 3: debootButton
			layout.FixedSize(u * 2), // 4: clearButton
			layout.FlexibleSize(1),  // 5: statusText
		},
		RowGap: u,
	}
	// Add package label and input
	appender.AppendChildWidgetWithBounds(&r.packageLabel, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&r.packageInput, gl.CellBounds(0, 1))
	// Add boot button
	appender.AppendChildWidgetWithBounds(&r.bootButton, gl.CellBounds(0, 2))
	// Add deboot button
	appender.AppendChildWidgetWithBounds(&r.debootButton, gl.CellBounds(0, 3))
	// Add clear button
	appender.AppendChildWidgetWithBounds(&r.clearButton, gl.CellBounds(0, 4))
	// Add status text
	appender.AppendChildWidgetWithBounds(&r.statusText, gl.CellBounds(0, 5))

	return nil
}

// Launch is the public entry point to start the GUI.
func Launch() {
	op := &guigui.RunOptions{
		Title:         "PB-Stack Bootstrapper",
		WindowMinSize: image.Pt(320, 240), // Adjusted size
		// WindowMaxSize: image.Pt(640, 480), // Remove or comment out to allow full maximization
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		log.Printf("Error running GUI: %v", err) // Changed to log
		os.Exit(1)
	}
}
