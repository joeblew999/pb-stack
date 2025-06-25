package bridge

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// ProcessDeckshFile processes a .dsh file and returns the generated deck XML
func ProcessDeckshFile(filename string) ([]byte, error) {
	// Check if decksh command is available
	deckshPath, err := exec.LookPath("decksh")
	if err != nil {
		return nil, fmt.Errorf("decksh command not found in PATH. Please install decksh: go install github.com/ajstarks/decksh/cmd/decksh@latest")
	}

	// Read the input file
	input, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	// Execute decksh command
	cmd := exec.Command(deckshPath)
	cmd.Stdin = input

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("decksh execution failed: %w\nStderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// ProcessDeckshString processes decksh content from a string and returns deck XML
func ProcessDeckshString(content string) ([]byte, error) {
	deckshPath, err := exec.LookPath("decksh")
	if err != nil {
		return nil, fmt.Errorf("decksh command not found in PATH")
	}

	cmd := exec.Command(deckshPath)
	cmd.Stdin = strings.NewReader(content)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("decksh execution failed: %w\nStderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// Deck represents the root deck structure
type Deck struct {
	XMLName xml.Name `xml:"deck"`
	Slides  []Slide  `xml:"slide"`
}

// Slide represents a single slide
type Slide struct {
	XMLName    xml.Name `xml:"slide"`
	Background string   `xml:"bg,attr,omitempty"`
	Foreground string   `xml:"fg,attr,omitempty"`
	Elements   []Element
}

// Element represents any slide element (text, image, shape, etc.)
type Element struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",chardata"`
}

// UnmarshalXML custom unmarshaling for slide elements
func (s *Slide) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Set background and foreground from attributes
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "bg":
			s.Background = attr.Value
		case "fg":
			s.Foreground = attr.Value
		}
	}

	// Parse child elements
	for {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			var elem Element
			elem.XMLName = t.Name
			elem.Attrs = t.Attr

			// Read the content
			if err := d.DecodeElement(&elem, &t); err != nil {
				return err
			}

			s.Elements = append(s.Elements, elem)
		case xml.EndElement:
			if t.Name == start.Name {
				return nil
			}
		}
	}

	return nil
}

// ParseDeckXML parses deck XML into a structured presentation
func ParseDeckXML(xmlData []byte) (*Presentation, error) {
	var deck Deck
	err := xml.Unmarshal(xmlData, &deck)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deck XML: %w", err)
	}

	presentation := &Presentation{
		Slides: make([]PresentationSlide, len(deck.Slides)),
	}

	for i, slide := range deck.Slides {
		pSlide := PresentationSlide{
			Background: slide.Background,
			Foreground: slide.Foreground,
			Elements:   make([]PresentationElement, len(slide.Elements)),
		}

		for j, elem := range slide.Elements {
			pElem := PresentationElement{
				Type:       elem.XMLName.Local,
				Attributes: make(map[string]string),
				Content:    strings.TrimSpace(elem.Content),
			}

			// Convert XML attributes to map
			for _, attr := range elem.Attrs {
				pElem.Attributes[attr.Name.Local] = attr.Value
			}

			pSlide.Elements[j] = pElem
		}

		presentation.Slides[i] = pSlide
	}

	return presentation, nil
}
