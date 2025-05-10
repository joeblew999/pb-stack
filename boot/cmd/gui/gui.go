package gui

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.DefaultWidget

	background   basicwidget.Background
	targetLabel  basicwidget.Text
	targetInput  basicwidget.TextInput
	bootButton   basicwidget.TextButton
	debootButton basicwidget.TextButton
	statusText   basicwidget.Text
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	// r.background.SetColor(color.NRGBA{R: 0x1e, G: 0x20, B: 0x26, A: 0xff}) // Dark background - Method not found in current guigui version
	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	r.targetLabel.SetValue("Target Host/IP:")
	// r.targetInput.SetPlaceholder("e.g., 192.168.1.100 or server.example.com") // Method not available
	if r.targetInput.Value() == "" {
		// You could set a default target if desired
		// r.targetInput.SetValue("localhost")
	}

	r.statusText.SetSelectable(false) // Status text usually isn't selectable
	r.statusText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	r.statusText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	r.statusText.SetScale(1)        // Normal scale for status text
	if r.statusText.Value() == "" { // Set initial value only if not already set
		r.statusText.SetValue("Select an action.")
	}

	r.bootButton.SetText("Boot System")
	r.bootButton.SetOnUp(func() {
		// Run in a goroutine to avoid blocking the UI thread
		go runCLIProcess("Booting", "-boot", r.targetInput.Value(), &r.statusText)
	})

	r.debootButton.SetText("Deboot System")
	r.debootButton.SetOnUp(func() {
		// Run in a goroutine to avoid blocking the UI thread
		go runCLIProcess("Debooting", "-deboot", r.targetInput.Value(), &r.statusText)
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r).Inset(u),
		Heights: []layout.Size{
			layout.FixedSize(u),     // For targetLabel
			layout.FixedSize(u * 2), // For targetInput
			layout.FixedSize(u * 2), // For bootButton
			layout.FixedSize(u * 2), // For debootButton
			layout.FlexibleSize(1),  // For statusText
		},
		RowGap: u,
	}

	// Add boot button to the first row
	appender.AppendChildWidgetWithBounds(&r.bootButton, gl.CellBounds(0, 0))

	// Add target label
	appender.AppendChildWidgetWithBounds(&r.targetLabel, gl.CellBounds(0, 0))
	// Add target input field
	appender.AppendChildWidgetWithBounds(&r.targetInput, gl.CellBounds(0, 1))
	// Add boot button
	appender.AppendChildWidgetWithBounds(&r.bootButton, gl.CellBounds(0, 2))
	// Add deboot button
	appender.AppendChildWidgetWithBounds(&r.debootButton, gl.CellBounds(0, 3))
	// Add status text
	appender.AppendChildWidgetWithBounds(&r.statusText, gl.CellBounds(0, 4))

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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
