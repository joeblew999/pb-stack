package bridge

import (
	"image/color"
	"strings"
	"testing"
)

func TestProcessDeckshString(t *testing.T) {
	// Simple decksh content
	deckshContent := `
deck
    slide "white" "black"
        ctext "Hello, World!" 50 50 5
    eslide
edeck
`

	// This test will only work if decksh is installed
	xmlData, err := ProcessDeckshString(deckshContent)
	if err != nil {
		if strings.Contains(err.Error(), "decksh command not found") {
			t.Skip("decksh command not found, skipping test")
		}
		t.Fatalf("ProcessDeckshString failed: %v", err)
	}

	if len(xmlData) == 0 {
		t.Fatal("Expected XML data, got empty result")
	}

	// Check that we got XML
	xmlStr := string(xmlData)
	if !strings.Contains(xmlStr, "<deck>") {
		t.Errorf("Expected XML to contain <deck>, got: %s", xmlStr)
	}
}

func TestParseDeckXML(t *testing.T) {
	// Sample deck XML
	xmlData := []byte(`
<deck>
    <slide bg="white" fg="black">
        <text align="c" xp="50" yp="50" sp="5">Hello, World!</text>
    </slide>
</deck>
`)

	presentation, err := ParseDeckXML(xmlData)
	if err != nil {
		t.Fatalf("ParseDeckXML failed: %v", err)
	}

	if len(presentation.Slides) != 1 {
		t.Errorf("Expected 1 slide, got %d", len(presentation.Slides))
	}

	slide := presentation.Slides[0]
	if slide.Background != "white" {
		t.Errorf("Expected background 'white', got '%s'", slide.Background)
	}

	if slide.Foreground != "black" {
		t.Errorf("Expected foreground 'black', got '%s'", slide.Foreground)
	}

	if len(slide.Elements) != 1 {
		t.Errorf("Expected 1 element, got %d", len(slide.Elements))
	}

	element := slide.Elements[0]
	if element.Type != "text" {
		t.Errorf("Expected element type 'text', got '%s'", element.Type)
	}

	if element.Content != "Hello, World!" {
		t.Errorf("Expected content 'Hello, World!', got '%s'", element.Content)
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		input    string
		expected color.NRGBA
	}{
		{"red", color.NRGBA{255, 0, 0, 255}},
		{"blue", color.NRGBA{0, 0, 255, 255}},
		{"#ff0000", color.NRGBA{255, 0, 0, 255}},
		{"#00ff00", color.NRGBA{0, 255, 0, 255}},
		{"rgb(255,0,0)", color.NRGBA{255, 0, 0, 255}},
		{"rgb(0, 255, 0)", color.NRGBA{0, 255, 0, 255}},
	}

	defaultColor := color.NRGBA{128, 128, 128, 255}

	for _, test := range tests {
		result := ParseColor(test.input, defaultColor)
		if result != test.expected {
			t.Errorf("ParseColor(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestPresentationElementGetters(t *testing.T) {
	element := PresentationElement{
		Type: "text",
		Attributes: map[string]string{
			"xp":    "50.5",
			"yp":    "25.0",
			"sp":    "3.5",
			"color": "red",
		},
		Content: "Test text",
	}

	// Test GetFloat
	if x := element.GetFloat("xp", 0); x != 50.5 {
		t.Errorf("GetFloat('xp') = %f, expected 50.5", x)
	}

	if missing := element.GetFloat("missing", 99.9); missing != 99.9 {
		t.Errorf("GetFloat('missing') = %f, expected 99.9", missing)
	}

	// Test GetString
	if color := element.GetString("color", ""); color != "red" {
		t.Errorf("GetString('color') = %s, expected 'red'", color)
	}

	if missing := element.GetString("missing", "default"); missing != "default" {
		t.Errorf("GetString('missing') = %s, expected 'default'", missing)
	}

	// Test GetColor
	expectedRed := color.NRGBA{255, 0, 0, 255}
	if c := element.GetColor("color", color.NRGBA{}); c != expectedRed {
		t.Errorf("GetColor('color') = %v, expected %v", c, expectedRed)
	}
}
