package bridge

import (
	"image/color"
	"strconv"
	"strings"
)

// Presentation represents a complete presentation
type Presentation struct {
	Slides []PresentationSlide
}

// PresentationSlide represents a single slide in the presentation
type PresentationSlide struct {
	Background string
	Foreground string
	Elements   []PresentationElement
}

// PresentationElement represents an element on a slide
type PresentationElement struct {
	Type       string            // text, image, rect, circle, etc.
	Attributes map[string]string // XML attributes
	Content    string            // Text content
}

// GetFloat gets a float attribute with a default value
func (e *PresentationElement) GetFloat(key string, defaultValue float64) float64 {
	if val, ok := e.Attributes[key]; ok {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

// GetString gets a string attribute with a default value
func (e *PresentationElement) GetString(key string, defaultValue string) string {
	if val, ok := e.Attributes[key]; ok {
		return val
	}
	return defaultValue
}

// GetColor parses a color string and returns a color.NRGBA
func (e *PresentationElement) GetColor(key string, defaultColor color.NRGBA) color.NRGBA {
	colorStr := e.GetString(key, "")
	if colorStr == "" {
		return defaultColor
	}
	
	return ParseColor(colorStr, defaultColor)
}

// ParseColor parses various color formats and returns color.NRGBA
func ParseColor(colorStr string, defaultColor color.NRGBA) color.NRGBA {
	colorStr = strings.TrimSpace(colorStr)
	
	// Handle hex colors
	if strings.HasPrefix(colorStr, "#") {
		return parseHexColor(colorStr, defaultColor)
	}
	
	// Handle rgb() colors
	if strings.HasPrefix(colorStr, "rgb(") {
		return parseRGBColor(colorStr, defaultColor)
	}
	
	// Handle named colors
	if namedColor, ok := namedColors[strings.ToLower(colorStr)]; ok {
		return namedColor
	}
	
	return defaultColor
}

// parseHexColor parses hex color strings like #ff0000
func parseHexColor(hex string, defaultColor color.NRGBA) color.NRGBA {
	if len(hex) < 7 {
		return defaultColor
	}
	
	hex = hex[1:] // Remove #
	
	var r, g, b uint8 = 0, 0, 0
	var a uint8 = 255
	
	if len(hex) >= 6 {
		if val, err := strconv.ParseUint(hex[0:2], 16, 8); err == nil {
			r = uint8(val)
		}
		if val, err := strconv.ParseUint(hex[2:4], 16, 8); err == nil {
			g = uint8(val)
		}
		if val, err := strconv.ParseUint(hex[4:6], 16, 8); err == nil {
			b = uint8(val)
		}
	}
	
	if len(hex) >= 8 {
		if val, err := strconv.ParseUint(hex[6:8], 16, 8); err == nil {
			a = uint8(val)
		}
	}
	
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// parseRGBColor parses rgb(r,g,b) color strings
func parseRGBColor(rgb string, defaultColor color.NRGBA) color.NRGBA {
	rgb = strings.TrimPrefix(rgb, "rgb(")
	rgb = strings.TrimSuffix(rgb, ")")
	
	parts := strings.Split(rgb, ",")
	if len(parts) < 3 {
		return defaultColor
	}
	
	var r, g, b uint8 = 0, 0, 0
	var a uint8 = 255
	
	if val, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 8); err == nil {
		r = uint8(val)
	}
	if val, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 8); err == nil {
		g = uint8(val)
	}
	if val, err := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 8); err == nil {
		b = uint8(val)
	}
	
	if len(parts) >= 4 {
		if val, err := strconv.ParseUint(strings.TrimSpace(parts[3]), 10, 8); err == nil {
			a = uint8(val)
		}
	}
	
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

// namedColors maps color names to NRGBA values
var namedColors = map[string]color.NRGBA{
	"black":   {0, 0, 0, 255},
	"white":   {255, 255, 255, 255},
	"red":     {255, 0, 0, 255},
	"green":   {0, 255, 0, 255},
	"blue":    {0, 0, 255, 255},
	"yellow":  {255, 255, 0, 255},
	"cyan":    {0, 255, 255, 255},
	"magenta": {255, 0, 255, 255},
	"gray":    {128, 128, 128, 255},
	"grey":    {128, 128, 128, 255},
	"orange":  {255, 165, 0, 255},
	"purple":  {128, 0, 128, 255},
	"brown":   {165, 42, 42, 255},
	"pink":    {255, 192, 203, 255},
}
