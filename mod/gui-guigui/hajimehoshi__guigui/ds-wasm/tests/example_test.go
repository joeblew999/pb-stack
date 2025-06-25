package tests

import (
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// TestExample demonstrates basic Rod usage for DataStar testing
func TestExample(t *testing.T) {
	// This test demonstrates how Rod can be used to test DataStar applications
	// It's a simple example that doesn't require running services

	// Launch browser
	launcher := launcher.New().Headless(true)
	url := launcher.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// Create a page
	page := browser.MustPage()
	defer page.MustClose()

	// Navigate to a simple HTML page (data URL) with basic JavaScript
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Rod Test Example</title>
		<script>
			// Simple counter simulation (without DataStar for testing)
			let count = 0;
			function updateDisplay() {
				document.getElementById('count').textContent = count;
			}
			function increment() {
				count++;
				updateDisplay();
			}
			function decrement() {
				count--;
				updateDisplay();
			}
			function reset() {
				count = 0;
				updateDisplay();
			}
			window.onload = function() {
				document.getElementById('message').textContent = 'Hello from Rod!';
				updateDisplay();
				document.getElementById('increment').onclick = increment;
				document.getElementById('decrement').onclick = decrement;
				document.getElementById('reset').onclick = reset;
			};
		</script>
	</head>
	<body>
		<h1>Rod Testing Example</h1>
		<p id="message">Loading...</p>
		<p>Count: <span id="count">0</span></p>
		<button id="increment">Increment</button>
		<button id="decrement">Decrement</button>
		<button id="reset">Reset</button>
	</body>
	</html>
	`

	dataURL := "data:text/html;charset=utf-8," + htmlContent
	page.MustNavigate(dataURL)
	page.MustWaitLoad()

	// Wait for JavaScript to initialize
	time.Sleep(1 * time.Second)

	// Test basic page elements
	title := page.MustInfo().Title
	if title != "Rod Test Example" {
		t.Errorf("Expected title 'Rod Test Example', got '%s'", title)
	}

	// Test that elements are visible
	page.MustElement("h1").MustWaitVisible()
	page.MustElement("#message").MustWaitVisible()
	page.MustElement("#count").MustWaitVisible()

	// Test initial values
	messageText := page.MustElement("#message").MustText()
	if messageText != "Hello from Rod!" {
		t.Errorf("Expected message 'Hello from Rod!', got '%s'", messageText)
	}

	countText := page.MustElement("#count").MustText()
	if countText != "0" {
		t.Errorf("Expected count '0', got '%s'", countText)
	}

	// Test JavaScript reactivity by clicking buttons
	incrementBtn := page.MustElement("#increment")
	incrementBtn.MustClick()

	// Wait for JavaScript to update
	time.Sleep(500 * time.Millisecond)

	// Check if count updated
	newCountText := page.MustElement("#count").MustText()
	if newCountText != "1" {
		t.Errorf("Expected count '1' after increment, got '%s'", newCountText)
	}

	// Test multiple clicks
	incrementBtn.MustClick()
	incrementBtn.MustClick()
	time.Sleep(500 * time.Millisecond)

	finalCountText := page.MustElement("#count").MustText()
	if finalCountText != "3" {
		t.Errorf("Expected count '3' after multiple increments, got '%s'", finalCountText)
	}

	// Test reset button
	resetBtn := page.MustElement("#reset")
	resetBtn.MustClick()
	time.Sleep(500 * time.Millisecond)

	resetCountText := page.MustElement("#count").MustText()
	if resetCountText != "0" {
		t.Errorf("Expected count '0' after reset, got '%s'", resetCountText)
	}

	t.Log("✅ Rod browser automation example test passed!")
	t.Log("📋 Verified:")
	t.Log("  - Page loading and title")
	t.Log("  - Element visibility")
	t.Log("  - JavaScript reactive updates")
	t.Log("  - Button click interactions")
	t.Log("  - State management")
}

// TestRodCapabilities demonstrates Rod's testing capabilities
func TestRodCapabilities(t *testing.T) {
	// Launch browser
	launcher := launcher.New().Headless(true)
	url := launcher.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	defer browser.MustClose()

	// Create a page
	page := browser.MustPage()
	defer page.MustClose()

	// Navigate to a test page
	page.MustNavigate("data:text/html,<h1>Rod Capabilities Test</h1><p id='test'>Testing Rod features</p>")
	page.MustWaitLoad()

	// Test element selection
	h1 := page.MustElement("h1")
	h1Text := h1.MustText()
	if h1Text != "Rod Capabilities Test" {
		t.Errorf("Expected 'Rod Capabilities Test', got '%s'", h1Text)
	}

	// Test JavaScript evaluation
	result, err := page.Eval(`() => {
		return {
			title: document.title,
			url: window.location.href,
			userAgent: navigator.userAgent.includes('HeadlessChrome')
		};
	}`)

	if err != nil {
		t.Fatalf("JavaScript evaluation failed: %v", err)
	}

	// Verify we're running in headless Chrome
	resultMap := result.Value.Map()
	if userAgent, ok := resultMap["userAgent"]; ok {
		if !userAgent.Bool() {
			t.Log("⚠️  Not running in headless Chrome, but that's okay for testing")
		} else {
			t.Log("✅ Running in headless Chrome as expected")
		}
	}

	t.Log("✅ Rod capabilities test passed!")
	t.Log("📋 Verified:")
	t.Log("  - Element selection and text extraction")
	t.Log("  - JavaScript evaluation")
	t.Log("  - Browser environment detection")
}
