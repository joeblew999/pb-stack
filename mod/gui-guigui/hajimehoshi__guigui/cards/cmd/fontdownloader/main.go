package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

// Font represents a font file to download
type Font struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	URL         string `json:"url"`
	Family      string `json:"family"`
	Weight      string `json:"weight"`
	Style       string `json:"style"`
}

// FontConfig represents the JSON configuration file structure
type FontConfig struct {
	Fonts []Font `json:"fonts"`
}

// GoogleFont represents a font from the Google Fonts API
type GoogleFont struct {
	Family   string            `json:"family"`
	Variants []string          `json:"variants"`
	Subsets  []string          `json:"subsets"`
	Category string            `json:"category"`
	Files    map[string]string `json:"files"`
}

// GoogleFontsResponse represents the response from Google Fonts API
type GoogleFontsResponse struct {
	Items []GoogleFont `json:"items"`
}

// FontDownloader handles downloading fonts with progress tracking
type FontDownloader struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// NewFontDownloader creates a new font downloader with reasonable timeouts
func NewFontDownloader() *FontDownloader {
	return &FontDownloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://www.googleapis.com/webfonts/v1/webfonts",
	}
}

func main() {
	var (
		search      = flag.String("search", "", "Search for fonts by family name")
		list        = flag.Bool("list", false, "List available fonts from Google Fonts")
		category    = flag.String("category", "", "Filter fonts by category (serif, sans-serif, display, handwriting, monospace)")
		interactive = flag.Bool("interactive", false, "Interactive font selection mode")
		add         = flag.String("add", "", "Add a specific font family to fonts.json")
		update      = flag.Bool("update", false, "Update fonts.json with latest font URLs")
		apiKey      = flag.String("api-key", "", "Google Fonts API key (optional, uses public endpoint if not provided)")
		help        = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Get output directory from remaining args or use default
	args := flag.Args()
	outputDir := "./.assets/fonts"
	if len(args) > 0 {
		outputDir = args[0]
	}

	// Get config file path from remaining args or use default
	configPath := "fonts.json"
	if len(args) > 1 {
		configPath = args[1]
	}

	downloader := NewFontDownloader()
	if *apiKey != "" {
		downloader.apiKey = *apiKey
	}

	// Handle different modes
	switch {
	case *search != "":
		handleSearchMode(downloader, *search, *category)
	case *list:
		handleListMode(downloader, *category)
	case *interactive:
		handleInteractiveMode(downloader, configPath)
	case *add != "":
		handleAddMode(downloader, *add, configPath)
	case *update:
		handleUpdateMode(downloader, configPath)
	default:
		// Default behavior: download fonts from config
		handleDownloadMode(downloader, outputDir, configPath)
	}
}

func showHelp() {
	fmt.Println("🔤 Enhanced Font Downloader for Card Design")
	fmt.Println("==========================================")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Printf("  %s [flags] [output_dir] [config_file]\n", os.Args[0])
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  --search <query>      Search Google Fonts for fonts matching query")
	fmt.Println("  --list               List all available Google Fonts")
	fmt.Println("  --category <cat>     Filter by category (serif, sans-serif, display, handwriting, monospace)")
	fmt.Println("  --add <family>       Add a font family to fonts.json")
	fmt.Println("  --update             Update existing fonts in fonts.json with latest URLs")
	fmt.Println("  --interactive        Interactive font selection (coming soon)")
	fmt.Println("  --api-key <key>      Google Fonts API key (optional)")
	fmt.Println("  --help               Show this help")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Download fonts from fonts.json (default mode)")
	fmt.Printf("  %s\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Search for Roboto fonts")
	fmt.Printf("  %s --search Roboto\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # List all serif fonts")
	fmt.Printf("  %s --list --category serif\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Add Open Sans to your font config")
	fmt.Printf("  %s --add \"Open Sans\"\n", os.Args[0])
	fmt.Println()
	fmt.Println("  # Update all fonts in config with latest URLs")
	fmt.Printf("  %s --update\n", os.Args[0])
}

// handleSearchMode searches for fonts and displays results
func handleSearchMode(downloader *FontDownloader, query, category string) {
	fmt.Printf("🔍 Searching for fonts matching: %s\n", query)
	if category != "" {
		fmt.Printf("📂 Category filter: %s\n", category)
	}
	fmt.Println()

	fonts, err := downloader.searchFonts(query, category)
	if err != nil {
		fmt.Printf("❌ Search failed: %v\n", err)
		os.Exit(1)
	}

	if len(fonts) == 0 {
		fmt.Println("❌ No fonts found matching your criteria")
		return
	}

	fmt.Printf("✅ Found %d fonts:\n\n", len(fonts))
	for i, font := range fonts {
		fmt.Printf("%d. %s\n", i+1, font.Family)
		fmt.Printf("   Category: %s\n", font.Category)
		fmt.Printf("   Variants: %s\n", strings.Join(font.Variants, ", "))
		if len(font.Subsets) > 0 {
			fmt.Printf("   Subsets: %s\n", strings.Join(font.Subsets, ", "))
		}
		fmt.Println()
	}

	fmt.Println("💡 To add a font to your project:")
	fmt.Printf("   %s --add \"<FONT_FAMILY_NAME>\"\n", os.Args[0])
}

// handleListMode lists all available fonts
func handleListMode(downloader *FontDownloader, category string) {
	fmt.Println("📋 Listing Google Fonts...")
	if category != "" {
		fmt.Printf("📂 Category filter: %s\n", category)
	}
	fmt.Println()

	fonts, err := downloader.listAllFonts(category)
	if err != nil {
		fmt.Printf("❌ Failed to list fonts: %v\n", err)
		os.Exit(1)
	}

	// Group by category
	categories := make(map[string][]GoogleFont)
	for _, font := range fonts {
		categories[font.Category] = append(categories[font.Category], font)
	}

	// Sort categories
	var sortedCategories []string
	for cat := range categories {
		sortedCategories = append(sortedCategories, cat)
	}
	sort.Strings(sortedCategories)

	for _, cat := range sortedCategories {
		fontList := categories[cat]
		fmt.Printf("## %s (%d fonts)\n", strings.Title(cat), len(fontList))

		// Sort fonts within category
		sort.Slice(fontList, func(i, j int) bool {
			return fontList[i].Family < fontList[j].Family
		})

		for _, font := range fontList {
			fmt.Printf("  • %s (%d variants)\n", font.Family, len(font.Variants))
		}
		fmt.Println()
	}

	fmt.Printf("Total: %d fonts available\n", len(fonts))
}

// handleInteractiveMode provides an interactive font selection interface
func handleInteractiveMode(downloader *FontDownloader, configPath string) {
	fmt.Println("🎛️  Interactive Font Selection")
	fmt.Println("=============================")
	fmt.Println("This feature will be implemented in the next version.")
	fmt.Println("For now, use:")
	fmt.Printf("  %s --search <query>\n", os.Args[0])
	fmt.Printf("  %s --add \"<font-family>\"\n", os.Args[0])
}

// handleAddMode adds a font to the fonts.json configuration
func handleAddMode(downloader *FontDownloader, family, configPath string) {
	fmt.Printf("➕ Adding font family: %s\n", family)

	// Search for the specific font
	fonts, err := downloader.searchFonts(family, "")
	if err != nil {
		fmt.Printf("❌ Search failed: %v\n", err)
		os.Exit(1)
	}

	// Find exact match
	var targetFont *GoogleFont
	for _, font := range fonts {
		if strings.EqualFold(font.Family, family) {
			targetFont = &font
			break
		}
	}

	if targetFont == nil {
		fmt.Printf("❌ Font family '%s' not found\n", family)
		fmt.Println("\n💡 Available similar fonts:")
		for _, font := range fonts {
			fmt.Printf("   • %s\n", font.Family)
		}
		return
	}

	// Convert to our Font format and add to config
	err = downloader.addFontToConfig(*targetFont, configPath)
	if err != nil {
		fmt.Printf("❌ Failed to add font: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Added %s to %s\n", family, configPath)
}

// handleUpdateMode updates existing fonts in fonts.json with latest URLs
func handleUpdateMode(downloader *FontDownloader, configPath string) {
	fmt.Printf("🔄 Updating font URLs in %s\n", configPath)

	// Load existing config
	fonts, err := loadFontsFromJSON(configPath)
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 Found %d fonts to update\n", len(fonts))

	// For each font, try to find updated URL
	updated := 0
	for i, font := range fonts {
		fmt.Printf("🔍 Checking %s...\n", font.Family)

		// Search for current font family
		googleFonts, err := downloader.searchFonts(font.Family, "")
		if err != nil {
			fmt.Printf("   ⚠️  Search failed: %v\n", err)
			continue
		}

		// Find matching Google font
		for _, gf := range googleFonts {
			if strings.EqualFold(gf.Family, font.Family) {
				// Update URL if different
				newURL := downloader.getDownloadURL(gf, font.Weight, font.Style)
				if newURL != "" && newURL != font.URL {
					fonts[i].URL = newURL
					fmt.Printf("   ✅ Updated URL\n")
					updated++
				} else {
					fmt.Printf("   ➡️  No changes needed\n")
				}
				break
			}
		}
	}

	if updated > 0 {
		// Save updated config
		err = saveFontsToJSON(fonts, configPath)
		if err != nil {
			fmt.Printf("❌ Failed to save config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Updated %d fonts in %s\n", updated, configPath)
	} else {
		fmt.Println("✅ All fonts are up to date")
	}
}

// handleDownloadMode downloads fonts from the configuration (default behavior)
func handleDownloadMode(downloader *FontDownloader, outputDir, configPath string) {
	// Load font configuration from JSON
	fonts, err := loadFontsFromJSON(configPath)
	if err != nil {
		fmt.Printf("❌ Failed to load font configuration: %v\n", err)
		os.Exit(1)
	}

	// Ensure fonts directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create fonts directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🔤 Professional Font Downloader for Card Design")
	fmt.Println("===============================================")
	fmt.Printf("📦 Downloading %d fonts to %s/...\n\n", len(fonts), outputDir)

	successCount := 0
	for i, font := range fonts {
		fmt.Printf("[%d/%d] 📥 %s\n", i+1, len(fonts), font.DisplayName)

		// Full path to font file
		fontPath := fmt.Sprintf("%s/%s", outputDir, font.Name)

		// Check if font already exists
		if _, err := os.Stat(fontPath); err == nil {
			fmt.Printf("    ⏭️  Already exists, skipping\n")
			successCount++
			continue
		}

		err := downloader.downloadFont(font, fontPath)
		if err != nil {
			fmt.Printf("    ❌ Error: %v\n", err)
			continue
		}

		// Verify the download
		info, err := os.Stat(fontPath)
		if err != nil {
			fmt.Printf("    ❌ Verification failed: %v\n", err)
			continue
		}

		fmt.Printf("    ✅ Success (%s)\n", formatFileSize(info.Size()))
		successCount++
	}

	fmt.Printf("\n🎉 Font download complete: %d/%d successful\n", successCount, len(fonts))

	// Generate font manifest for documentation
	if successCount > 0 {
		generateFontManifest(fonts, successCount, outputDir)
	}

	// Create a backup font list using alternative sources
	if successCount < len(fonts) {
		fmt.Println("\n💡 Alternative: Using system fonts or font CDN services")
		fmt.Println("   Consider downloading fonts manually from:")
		fmt.Println("   - https://fonts.google.com")
		fmt.Println("   - https://fontsource.org")
		fmt.Println("   - Font releases from GitHub repositories")
	}
}

// searchFonts searches Google Fonts API for fonts matching the query
func (fd *FontDownloader) searchFonts(query, category string) ([]GoogleFont, error) {
	// Build URL with query parameters
	params := url.Values{}
	if fd.apiKey != "" {
		params.Set("key", fd.apiKey)
	}
	params.Set("sort", "popularity")

	searchURL := fd.baseURL + "?" + params.Encode()

	// Make request
	resp, err := fd.client.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fonts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %s", resp.Status)
	}

	var response GoogleFontsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Filter results
	var filtered []GoogleFont
	queryLower := strings.ToLower(query)
	categoryLower := strings.ToLower(category)

	for _, font := range response.Items {
		// Filter by query (family name)
		if query != "" && !strings.Contains(strings.ToLower(font.Family), queryLower) {
			continue
		}

		// Filter by category
		if category != "" && strings.ToLower(font.Category) != categoryLower {
			continue
		}

		filtered = append(filtered, font)
	}

	return filtered, nil
}

// listAllFonts lists all fonts, optionally filtered by category
func (fd *FontDownloader) listAllFonts(category string) ([]GoogleFont, error) {
	return fd.searchFonts("", category)
}

// getDownloadURL constructs a download URL for a specific font variant
func (fd *FontDownloader) getDownloadURL(font GoogleFont, weight, style string) string {
	// Try to find the exact variant
	variant := weight
	if style == "italic" {
		if weight == "400" {
			variant = "italic"
		} else {
			variant = weight + "italic"
		}
	}

	// Check if variant exists in files
	if url, exists := font.Files[variant]; exists {
		return url
	}

	// Fallback to regular if exact variant not found
	if url, exists := font.Files["regular"]; exists {
		return url
	}

	// Fallback to 400 if regular not found
	if url, exists := font.Files["400"]; exists {
		return url
	}

	// Return any available variant as last resort
	for _, url := range font.Files {
		return url
	}

	return ""
}

// addFontToConfig adds a Google Font to the fonts.json configuration
func (fd *FontDownloader) addFontToConfig(googleFont GoogleFont, configPath string) error {
	// Load existing config or create new one
	var config FontConfig
	if _, err := os.Stat(configPath); err == nil {
		fonts, err := loadFontsFromJSON(configPath)
		if err != nil {
			return fmt.Errorf("failed to load existing config: %w", err)
		}
		config.Fonts = fonts
	}

	// Convert Google Font to our Font format
	// Add regular variant by default
	regularURL := fd.getDownloadURL(googleFont, "400", "normal")
	if regularURL == "" {
		return fmt.Errorf("no usable font files found for %s", googleFont.Family)
	}

	// Extract filename from URL
	parts := strings.Split(regularURL, "/")
	filename := parts[len(parts)-1]

	// Handle Google Fonts URLs that might not have file extensions
	if !strings.Contains(filename, ".") {
		// Use a reasonable default based on font family
		safeName := strings.ReplaceAll(googleFont.Family, " ", "")
		filename = safeName + "-Regular.ttf"
	}

	newFont := Font{
		Name:        filename,
		DisplayName: googleFont.Family + " Regular",
		URL:         regularURL,
		Family:      googleFont.Family,
		Weight:      "400",
		Style:       "normal",
	}

	// Check if font already exists
	for _, existingFont := range config.Fonts {
		if existingFont.Family == newFont.Family && existingFont.Weight == newFont.Weight && existingFont.Style == newFont.Style {
			return fmt.Errorf("font %s (weight: %s, style: %s) already exists in config", newFont.Family, newFont.Weight, newFont.Style)
		}
	}

	// Add new font
	config.Fonts = append(config.Fonts, newFont)

	// Save config
	return saveFontsToJSON(config.Fonts, configPath)
}

// loadFontsFromJSON loads font configuration from a JSON file
func loadFontsFromJSON(configPath string) ([]Font, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", configPath, err)
	}
	defer file.Close()

	var config FontConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %w", err)
	}

	if len(config.Fonts) == 0 {
		return nil, fmt.Errorf("no fonts found in config file")
	}

	return config.Fonts, nil
}

// saveFontsToJSON saves fonts configuration to a JSON file
func saveFontsToJSON(fonts []Font, configPath string) error {
	config := FontConfig{Fonts: fonts}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (fd *FontDownloader) downloadFont(font Font, fontPath string) error {
	// Create the file
	out, err := os.Create(fontPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Make the request
	req, err := http.NewRequest("GET", font.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent to be nice to font servers
	req.Header.Set("User-Agent", "pb-stack-card-generator/1.0")

	resp, err := fd.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d: %s", resp.StatusCode, resp.Status)
	}

	// Verify content type (accept various font types)
	contentType := resp.Header.Get("Content-Type")
	acceptedTypes := []string{"font", "octet-stream", "application", "woff"}
	isValidFont := false
	for _, t := range acceptedTypes {
		if strings.Contains(contentType, t) {
			isValidFont = true
			break
		}
	}

	if !isValidFont {
		return fmt.Errorf("unexpected content type: %s", contentType)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func generateFontManifest(fonts []Font, successCount int, outputDir string) {
	manifestPath := fmt.Sprintf("%s/FONTS.md", outputDir)
	file, err := os.Create(manifestPath)
	if err != nil {
		fmt.Printf("⚠️  Could not generate font manifest: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "# Font Manifest\n\n")
	fmt.Fprintf(file, "Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "Total fonts: %d\n\n", successCount)

	fmt.Fprintf(file, "## Available Fonts\n\n")
	fmt.Fprintf(file, "| Font File | Display Name | Family | Weight | Style |\n")
	fmt.Fprintf(file, "|-----------|--------------|--------|--------|\n")

	for _, font := range fonts {
		fontPath := fmt.Sprintf("%s/%s", outputDir, font.Name)
		if _, err := os.Stat(fontPath); err == nil {
			fmt.Fprintf(file, "| `%s` | %s | %s | %s | %s |\n",
				font.Name, font.DisplayName, font.Family, font.Weight, font.Style)
		}
	}

	fmt.Fprintf(file, "\n## Usage in decksh\n\n")
	fmt.Fprintf(file, "Set the font directory:\n")
	fmt.Fprintf(file, "```bash\n")
	fmt.Fprintf(file, "export DECKFONTS=\"$(pwd)/fonts\"\n")
	fmt.Fprintf(file, "# or use -fontdir flag\n")
	fmt.Fprintf(file, "pngdeck -fontdir ./fonts input.xml output.png\n")
	fmt.Fprintf(file, "```\n")

	fmt.Fprintf(file, "\n## Font Notes\n\n")
	fmt.Fprintf(file, "- Some fonts are downloaded as WOFF2 format (web fonts)\n")
	fmt.Fprintf(file, "- decksh works best with TTF/OTF fonts\n")
	fmt.Fprintf(file, "- Consider converting WOFF2 to TTF if needed\n")
	fmt.Fprintf(file, "- Fira Sans and Fira Mono are reliably available as TTF\n")

	fmt.Printf("📄 Generated font manifest: %s\n", manifestPath)
}
