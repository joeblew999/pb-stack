package examples

import (
	"testing"
	"time"

	"web-testing/core"
)

// TestBasicBrowserAutomation demonstrates basic browser automation
func TestBasicBrowserAutomation(t *testing.T) {
	// Create browser manager with default config
	config := core.DefaultConfig()
	config.Headless = true
	config.Screenshots = true

	browser := core.NewBrowserManager(config)

	// Start browser
	if err := browser.Start(); err != nil {
		t.Fatalf("Failed to start browser: %v", err)
	}
	defer browser.Stop()

	// Create a page with simple HTML
	page, err := browser.NewPage("")
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}
	defer page.Close()

	// Navigate to a data URL with test content
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head><title>Rod Testing Example</title></head>
	<body>
		<h1 id="title">Hello Rod Testing!</h1>
		<button id="test-btn">Click Me</button>
		<div id="result" style="display:none;">Button Clicked!</div>
		<script>
			document.getElementById('test-btn').onclick = function() {
				document.getElementById('result').style.display = 'block';
			};
		</script>
	</body>
	</html>
	`

	dataURL := "data:text/html;charset=utf-8," + htmlContent
	page.Navigate(dataURL)

	// Test basic element interaction
	title, err := page.WaitForElement("#title", 5*time.Second)
	if err != nil {
		t.Fatalf("Title element not found: %v", err)
	}

	titleText := title.MustText()
	if titleText != "Hello Rod Testing!" {
		t.Errorf("Expected 'Hello Rod Testing!', got '%s'", titleText)
	}

	// Test button click
	err = page.ClickAndWait("#test-btn", 500*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to click button: %v", err)
	}

	// Verify result is visible
	result, err := page.WaitForElementVisible("#result", 5*time.Second)
	if err != nil {
		t.Fatalf("Result element not visible: %v", err)
	}

	resultText := result.MustText()
	if resultText != "Button Clicked!" {
		t.Errorf("Expected 'Button Clicked!', got '%s'", resultText)
	}

	// Take a screenshot for debugging
	if err := page.Screenshot("basic_test_success"); err != nil {
		t.Logf("Failed to take screenshot: %v", err)
	}

	t.Log("✅ Basic browser automation test passed!")
}

// TestFormInteraction demonstrates form testing
func TestFormInteraction(t *testing.T) {
	config := core.DefaultConfig()
	config.Headless = true

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

	// Create a form test page
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head><title>Form Test</title></head>
	<body>
		<form id="test-form">
			<input type="text" id="name" placeholder="Name" required>
			<input type="email" id="email" placeholder="Email" required>
			<button type="submit" id="submit-btn">Submit</button>
		</form>
		<div id="form-result" style="display:none;">Form Submitted!</div>
		<script>
			document.getElementById('test-form').onsubmit = function(e) {
				e.preventDefault();
				document.getElementById('form-result').style.display = 'block';
			};
		</script>
	</body>
	</html>
	`

	dataURL := "data:text/html;charset=utf-8," + htmlContent
	page.Navigate(dataURL)

	// Test form filling
	formData := map[string]string{
		"#name":  "John Doe",
		"#email": "john@example.com",
	}

	err = page.FillForm("#test-form", formData)
	if err != nil {
		t.Fatalf("Failed to fill form: %v", err)
	}

	// Submit form
	err = page.ClickAndWait("#submit-btn", 500*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to submit form: %v", err)
	}

	// Verify form submission
	_, err = page.WaitForElementVisible("#form-result", 5*time.Second)
	if err != nil {
		t.Fatalf("Form result not visible: %v", err)
	}

	t.Log("✅ Form interaction test passed!")
}

// TestMultiplePages demonstrates managing multiple pages
func TestMultiplePages(t *testing.T) {
	config := core.DefaultConfig()
	config.Headless = true

	browser := core.NewBrowserManager(config)
	if err := browser.Start(); err != nil {
		t.Fatalf("Failed to start browser: %v", err)
	}
	defer browser.Stop()

	// Create first page
	page1, err := browser.NewPage("")
	if err != nil {
		t.Fatalf("Failed to create page 1: %v", err)
	}
	defer page1.Close()

	// Create second page
	page2, err := browser.NewPage("")
	if err != nil {
		t.Fatalf("Failed to create page 2: %v", err)
	}
	defer page2.Close()

	// Navigate pages to different content
	html1 := "data:text/html,<h1 id='page'>Page 1</h1>"
	html2 := "data:text/html,<h1 id='page'>Page 2</h1>"

	page1.Navigate(html1)
	page2.Navigate(html2)

	// Verify each page has correct content
	title1, err := page1.WaitForElement("#page", 5*time.Second)
	if err != nil {
		t.Fatalf("Page 1 title not found: %v", err)
	}

	title2, err := page2.WaitForElement("#page", 5*time.Second)
	if err != nil {
		t.Fatalf("Page 2 title not found: %v", err)
	}

	text1 := title1.MustText()
	text2 := title2.MustText()

	if text1 != "Page 1" {
		t.Errorf("Page 1 expected 'Page 1', got '%s'", text1)
	}

	if text2 != "Page 2" {
		t.Errorf("Page 2 expected 'Page 2', got '%s'", text2)
	}

	t.Log("✅ Multiple pages test passed!")
}

// TestScreenshotCapture demonstrates screenshot functionality
func TestScreenshotCapture(t *testing.T) {
	config := core.DefaultConfig()
	config.Headless = true
	config.Screenshots = true
	config.ScreenshotDir = "test_screenshots"

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

	// Create a colorful test page
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head><title>Screenshot Test</title></head>
	<body style="background: linear-gradient(45deg, #ff6b6b, #4ecdc4); padding: 50px;">
		<h1 style="color: white; text-align: center;">Screenshot Test Page</h1>
		<div style="background: white; padding: 20px; border-radius: 10px; margin: 20px;">
			<p>This page is designed to test screenshot functionality.</p>
			<button id="color-btn" style="background: #ff6b6b; color: white; padding: 10px;">Colorful Button</button>
		</div>
	</body>
	</html>
	`

	dataURL := "data:text/html;charset=utf-8," + htmlContent
	page.Navigate(dataURL)

	// Take initial screenshot
	if err := page.Screenshot("initial_state"); err != nil {
		t.Fatalf("Failed to take initial screenshot: %v", err)
	}

	// Click button and take another screenshot
	err = page.ClickAndWait("#color-btn", 500*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to click button: %v", err)
	}

	if err := page.Screenshot("after_click"); err != nil {
		t.Fatalf("Failed to take after-click screenshot: %v", err)
	}

	t.Log("✅ Screenshot capture test passed!")
	t.Log("📸 Screenshots saved to test_screenshots/ directory")
}
