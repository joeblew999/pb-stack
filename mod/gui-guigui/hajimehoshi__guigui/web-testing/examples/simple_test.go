package examples

import (
	"testing"
	"time"

	"web-testing/core"
)

// TestSimpleBrowser demonstrates basic browser functionality
func TestSimpleBrowser(t *testing.T) {
	// Create browser manager with default config
	config := core.DefaultConfig()
	config.Headless = true
	config.Screenshots = false // Disable screenshots for simple test
	config.Timeout = 10 * time.Second

	browser := core.NewBrowserManager(config)

	// Start browser
	if err := browser.Start(); err != nil {
		t.Fatalf("Failed to start browser: %v", err)
	}
	defer browser.Stop()

	// Create a page
	page, err := browser.NewPage("")
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}
	defer page.Close()

	// Navigate to a simple data URL
	simpleHTML := "data:text/html,<h1 id='test'>Hello Rod Testing Framework!</h1>"
	page.Navigate(simpleHTML)

	// Test element finding
	element, err := page.WaitForElement("#test", 5*time.Second)
	if err != nil {
		t.Fatalf("Element not found: %v", err)
	}

	// Get text without using Must methods to avoid panics
	text, err := element.Text()
	if err != nil {
		t.Fatalf("Failed to get text: %v", err)
	}

	if text != "Hello Rod Testing Framework!" {
		t.Errorf("Expected 'Hello Rod Testing Framework!', got '%s'", text)
	}

	t.Log("✅ Simple browser test passed!")
}

// TestBrowserCapabilities tests basic browser capabilities
func TestBrowserCapabilities(t *testing.T) {
	config := core.DefaultConfig()
	config.Headless = true
	config.Screenshots = false

	browser := core.NewBrowserManager(config)
	if err := browser.Start(); err != nil {
		t.Fatalf("Failed to start browser: %v", err)
	}
	defer browser.Stop()

	page, err := browser.NewPage("")
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}
	defer page.Close()

	// Test JavaScript evaluation
	result, err := page.Page().Eval(`() => {
		return {
			userAgent: navigator.userAgent.includes('Chrome'),
			title: document.title || 'No Title',
			url: window.location.href
		};
	}`)

	if err != nil {
		t.Fatalf("JavaScript evaluation failed: %v", err)
	}

	// Verify we got a result
	resultMap := result.Value.Map()
	if len(resultMap) == 0 {
		t.Fatalf("No result from JavaScript evaluation")
	}

	t.Log("✅ Browser capabilities test passed!")
	t.Logf("📋 Browser info: %v", resultMap)
}

// TestPageNavigation tests basic page navigation
func TestPageNavigation(t *testing.T) {
	config := core.DefaultConfig()
	config.Headless = true
	config.Screenshots = false

	browser := core.NewBrowserManager(config)
	if err := browser.Start(); err != nil {
		t.Fatalf("Failed to start browser: %v", err)
	}
	defer browser.Stop()

	page, err := browser.NewPage("")
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}
	defer page.Close()

	// Test navigation to different pages
	pages := []struct {
		url      string
		expected string
	}{
		{"data:text/html,<title>Page 1</title><h1>First Page</h1>", "Page 1"},
		{"data:text/html,<title>Page 2</title><h1>Second Page</h1>", "Page 2"},
		{"data:text/html,<title>Page 3</title><h1>Third Page</h1>", "Page 3"},
	}

	for i, testPage := range pages {
		page.Navigate(testPage.url)

		// Get page title
		title, err := page.Page().Info()
		if err != nil {
			t.Fatalf("Failed to get page info for page %d: %v", i+1, err)
		}

		if title.Title != testPage.expected {
			t.Errorf("Page %d: expected title '%s', got '%s'", i+1, testPage.expected, title.Title)
		}
	}

	t.Log("✅ Page navigation test passed!")
}
