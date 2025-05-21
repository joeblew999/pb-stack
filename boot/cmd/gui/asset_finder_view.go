package gui

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"

	"main/cmd/cli" // Import the cli package to use the asset finder
)

// AssetFinderViewWidget is the GUI view for the GitHub Asset Finder feature.
type AssetFinderViewWidget struct {
	guigui.DefaultWidget

	background basicwidget.Background

	configPathLabel basicwidget.Text
	configPathInput basicwidget.TextInput
	searchButton    basicwidget.TextButton
	resultsText     basicwidget.Text // Widget to display the search results
}

// Build constructs the UI for the AssetFinderViewWidget.
func (v *AssetFinderViewWidget) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	// v.background.SetColor(...) // Optional: Set a background color
	appender.AppendChildWidgetWithBounds(&v.background, context.Bounds(v))

	v.configPathLabel.SetValue("Asset Search Config File (YAML):")
	// Set a default value, maybe relative to the executable or a known location
	if v.configPathInput.Value() == "" {
		// A reasonable default path might be relative to the executable or a standard config location
		// For now, let's use a placeholder or a known test file path
		v.configPathInput.SetValue("./asset-searches.yml") // Assuming it's in the boot directory
	}

	v.searchButton.SetText("Run Searches")
	v.searchButton.SetOnUp(func() {
		// Run the asset search logic in a goroutine to avoid blocking the GUI
		go v.runAssetSearches(context)
	})

	// Configure the results display area
	v.resultsText.SetSelectable(true)
	v.resultsText.SetHorizontalAlign(basicwidget.HorizontalAlignStart)
	v.resultsText.SetVerticalAlign(basicwidget.VerticalAlignTop)
	v.resultsText.SetMultiline(true)
	v.resultsText.SetAutoWrap(true)
	v.resultsText.SetValue("Click 'Run Searches' to load configurations and find assets.")

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(v).Inset(u), // Add padding
		Heights: []layout.Size{
			layout.FixedSize(u),     // Label
			layout.FixedSize(u * 2), // Input
			layout.FixedSize(u * 2), // Button
			layout.FlexibleSize(1),  // Results text area
		},
		RowGap: u,
	}

	appender.AppendChildWidgetWithBounds(&v.configPathLabel, gl.CellBounds(0, 0))
	appender.AppendChildWidgetWithBounds(&v.configPathInput, gl.CellBounds(0, 1))
	appender.AppendChildWidgetWithBounds(&v.searchButton, gl.CellBounds(0, 2))
	appender.AppendChildWidgetWithBounds(&v.resultsText, gl.CellBounds(0, 3))

	return nil
}

// runAssetSearches is called in a goroutine to execute the asset finding logic.
func (v *AssetFinderViewWidget) runAssetSearches(context *guigui.Context) {
	configPath := v.configPathInput.Value()
	apiToken := os.Getenv("GITHUB_TOKEN") // Get token from env var for now

	// Update status text in the GUI (must be on the main GUI thread)
	v.resultsText.SetValue("Loading and running searches...")
	context.RequestUpdate() // Request an update to redraw the widget

	results, err := cli.FindAssetsFromConfig(configPath, apiToken)

	// Update results display in the GUI (must be on the main GUI thread)
	v.updateResultsDisplay(context, results, err)
}

// updateResultsDisplay updates the resultsText widget on the main GUI thread.
func (v *AssetFinderViewWidget) updateResultsDisplay(context *guigui.Context, results []cli.SearchResult, err error) {
	if err != nil {
		v.resultsText.SetValue(fmt.Sprintf("Error loading config: %v", err))
		log.Printf("GUI: Error loading asset config: %v", err)
		context.RequestUpdate()
		return
	}

	var output strings.Builder
	if len(results) == 0 {
		output.WriteString("No search results returned.")
	} else {
		for _, res := range results {
			output.WriteString(fmt.Sprintf("--- Search: %s ---\n", res.SearchName))
			if res.SearchDescription != "" {
				output.WriteString(fmt.Sprintf("Description: %s\n", res.SearchDescription))
			}
			if res.Error != nil {
				output.WriteString(fmt.Sprintf("Error: %v\n", res.Error))
			} else {
				if len(res.SelectedOutput) > 0 {
					output.WriteString("Output:\n")
					for _, line := range res.SelectedOutput {
						output.WriteString(fmt.Sprintf("  %s\n", line))
					}
				}
				if res.DownloadStatus != "" {
					output.WriteString(fmt.Sprintf("Download Status:\n%s\n", res.DownloadStatus))
				}
				if len(res.SelectedOutput) == 0 && res.DownloadStatus == "" {
					output.WriteString("No output or download status.\n")
				}
			}
			output.WriteString("\n") // Separator between searches
		}
	}
	v.resultsText.SetValue(output.String())
	log.Println("GUI: Asset search results updated.")
	context.RequestUpdate() // Request an update to redraw the widget
}
