package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"decksh-ebdeck/pkg/bridge"

	"github.com/hajimehoshi/ebiten/v2"
)

// Canvas wraps ebiten.Image with percentage-based coordinate system
type Canvas struct {
	Image  *ebiten.Image
	Width  int
	Height int
}

// NewCanvas creates a new canvas with the specified dimensions
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		Image:  ebiten.NewImage(width, height),
		Width:  width,
		Height: height,
	}
}

// Clear fills the canvas with the specified color
func (c *Canvas) Clear(color color.NRGBA) {
	c.Image.Fill(color)
}

// pctToPixel converts percentage coordinates to pixel coordinates
func (c *Canvas) pctToPixel(x, y float64) (int, int) {
	px := int(x * float64(c.Width) / 100.0)
	py := int((100.0 - y) * float64(c.Height) / 100.0) // Flip Y axis
	return px, py
}

// pctToPixelSize converts percentage size to pixel size
func (c *Canvas) pctToPixelSize(size float64) int {
	return int(size * float64(c.Width) / 100.0)
}

// RenderSlide renders a single slide to the canvas
func (c *Canvas) RenderSlide(slide bridge.PresentationSlide) error {
	// Set background color
	bgColor := bridge.ParseColor(slide.Background, color.NRGBA{255, 255, 255, 255})
	c.Clear(bgColor)

	// Render each element
	for _, element := range slide.Elements {
		if err := c.renderElement(element, slide.Foreground); err != nil {
			return fmt.Errorf("failed to render element %s: %w", element.Type, err)
		}
	}

	return nil
}

// renderElement renders a single element on the canvas
func (c *Canvas) renderElement(element bridge.PresentationElement, defaultColor string) error {
	switch element.Type {
	case "text":
		return c.renderText(element, defaultColor)
	case "rect":
		return c.renderRect(element, defaultColor)
	case "circle":
		return c.renderCircle(element, defaultColor)
	case "line":
		return c.renderLine(element, defaultColor)
	case "image":
		return c.renderImage(element)
	default:
		// For unsupported elements, just log and continue
		fmt.Printf("Unsupported element type: %s\n", element.Type)
		return nil
	}
}

// renderText renders text elements
func (c *Canvas) renderText(element bridge.PresentationElement, defaultColor string) error {
	x := element.GetFloat("xp", 50)
	y := element.GetFloat("yp", 50)
	size := element.GetFloat("sp", 3)

	// Get text color
	textColor := element.GetColor("color", bridge.ParseColor(defaultColor, color.NRGBA{0, 0, 0, 255}))

	// Convert to pixel coordinates
	px, py := c.pctToPixel(x, y)
	fontSize := c.pctToPixelSize(size)

	// For now, we'll use a simple text rendering approach
	// In a full implementation, you'd use a proper font rendering library
	fmt.Printf("Rendering text at (%d,%d) size %d color %v: %s\n",
		px, py, fontSize, textColor, element.Content)

	// TODO: Implement actual text rendering with ebiten
	// This would require loading fonts and using ebiten's text rendering

	return nil
}

// renderRect renders rectangle elements
func (c *Canvas) renderRect(element bridge.PresentationElement, defaultColor string) error {
	x := element.GetFloat("xp", 50)
	y := element.GetFloat("yp", 50)
	w := element.GetFloat("wp", 10)
	h := element.GetFloat("hp", 10)

	fillColor := element.GetColor("color", bridge.ParseColor(defaultColor, color.NRGBA{128, 128, 128, 255}))

	// Convert to pixel coordinates
	px, py := c.pctToPixel(x, y)
	pw := c.pctToPixelSize(w)
	ph := c.pctToPixelSize(h)

	// Create a rectangle image and draw it
	rectImg := ebiten.NewImage(pw, ph)
	rectImg.Fill(fillColor)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(px-pw/2), float64(py-ph/2))
	c.Image.DrawImage(rectImg, opts)

	return nil
}

// renderCircle renders circle elements
func (c *Canvas) renderCircle(element bridge.PresentationElement, defaultColor string) error {
	x := element.GetFloat("xp", 50)
	y := element.GetFloat("yp", 50)
	r := element.GetFloat("rp", 5)

	fillColor := element.GetColor("color", bridge.ParseColor(defaultColor, color.NRGBA{128, 128, 128, 255}))

	// Convert to pixel coordinates
	px, py := c.pctToPixel(x, y)
	pr := c.pctToPixelSize(r)

	// Create a circle image (simplified as a square for now)
	circleImg := ebiten.NewImage(pr*2, pr*2)
	circleImg.Fill(fillColor)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(px-pr), float64(py-pr))
	c.Image.DrawImage(circleImg, opts)

	return nil
}

// renderLine renders line elements
func (c *Canvas) renderLine(element bridge.PresentationElement, defaultColor string) error {
	x1 := element.GetFloat("x1p", 0)
	y1 := element.GetFloat("y1p", 0)
	x2 := element.GetFloat("x2p", 100)
	y2 := element.GetFloat("y2p", 100)

	strokeColor := element.GetColor("color", bridge.ParseColor(defaultColor, color.NRGBA{0, 0, 0, 255}))

	// Convert to pixel coordinates
	px1, py1 := c.pctToPixel(x1, y1)
	px2, py2 := c.pctToPixel(x2, y2)

	// Simple line rendering (would need proper line drawing in full implementation)
	fmt.Printf("Rendering line from (%d,%d) to (%d,%d) color %v\n",
		px1, py1, px2, py2, strokeColor)

	// TODO: Implement actual line rendering

	return nil
}

// renderImage renders image elements
func (c *Canvas) renderImage(element bridge.PresentationElement) error {
	x := element.GetFloat("xp", 50)
	y := element.GetFloat("yp", 50)
	scale := element.GetFloat("scale", 100)

	filename := element.GetString("href", "")
	if filename == "" {
		return fmt.Errorf("image element missing href attribute")
	}

	// Convert to pixel coordinates
	px, py := c.pctToPixel(x, y)

	fmt.Printf("Rendering image %s at (%d,%d) scale %.1f%%\n",
		filename, px, py, scale)

	// TODO: Load and render actual image

	return nil
}

// RenderToFile renders a presentation to a PNG file
func RenderToFile(presentation *bridge.Presentation, filename string, width, height int) error {
	if len(presentation.Slides) == 0 {
		return fmt.Errorf("presentation has no slides")
	}

	// For now, just render the first slide
	canvas := NewCanvas(width, height)

	err := canvas.RenderSlide(presentation.Slides[0])
	if err != nil {
		return fmt.Errorf("failed to render slide: %w", err)
	}

	// Convert ebiten.Image to standard image.Image
	bounds := canvas.Image.Bounds()
	img := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := canvas.Image.At(x, y)
			img.Set(x, y, c)
		}
	}

	// Save to file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}
