package examples

import (
	"testing"
	"time"

	"web-testing/core"
	"web-testing/frameworks"
)

// TestDataStarBasic demonstrates basic DataStar testing
func TestDataStarBasic(t *testing.T) {
	config := core.DefaultConfig()
	config.Headless = true
	config.Screenshots = true

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

	// Create DataStar test page
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>DataStar Test</title>
		<script type="module" defer src="https://cdn.jsdelivr.net/npm/@starfederation/datastar@latest"></script>
	</head>
	<body>
		<div data-store='{"count": 0, "message": "Hello DataStar!"}'>
			<h1 data-text="$message"></h1>
			<p>Count: <span id="counter" data-text="$count"></span></p>
			<button id="increment" data-on-click="$count++">Increment</button>
			<button id="reset" data-on-click="$count = 0">Reset</button>
			<div id="status" data-show="$count > 5" style="color: red;">
				Count is high!
			</div>
		</div>
	</body>
	</html>
	`

	dataURL := "data:text/html;charset=utf-8," + htmlContent
	page.Navigate(dataURL)

	// Create DataStar helper
	datastar := frameworks.NewDataStarHelper(page)

	// Wait for DataStar to load
	if err := datastar.WaitForDataStar(10 * time.Second); err != nil {
		t.Fatalf("DataStar not loaded: %v", err)
	}

	// Test initial state
	store, err := datastar.GetStore()
	if err != nil {
		t.Fatalf("Failed to get store: %v", err)
	}

	if store["count"].(float64) != 0 {
		t.Errorf("Expected count 0, got %v", store["count"])
	}

	if store["message"].(string) != "Hello DataStar!" {
		t.Errorf("Expected 'Hello DataStar!', got %v", store["message"])
	}

	// Test increment functionality
	err = datastar.TriggerDataStarEvent("#increment", "#counter")
	if err != nil {
		t.Fatalf("Failed to trigger increment: %v", err)
	}

	// Verify count updated
	err = datastar.WaitForStoreValue("count", float64(1), 5*time.Second)
	if err != nil {
		t.Fatalf("Count not updated: %v", err)
	}

	// Test multiple increments
	for i := 0; i < 5; i++ {
		err = datastar.TriggerDataStarEvent("#increment", "")
		if err != nil {
			t.Fatalf("Failed to increment (iteration %d): %v", i, err)
		}
	}

	// Verify high count status shows
	_, err = page.WaitForElementVisible("#status", 5*time.Second)
	if err != nil {
		t.Fatalf("Status element not visible when count > 5: %v", err)
	}

	// Test reset
	err = datastar.TriggerDataStarEvent("#reset", "")
	if err != nil {
		t.Fatalf("Failed to reset: %v", err)
	}

	err = datastar.WaitForStoreValue("count", float64(0), 5*time.Second)
	if err != nil {
		t.Fatalf("Count not reset: %v", err)
	}

	// Take screenshot of final state
	if err := page.Screenshot("datastar_test_final"); err != nil {
		t.Logf("Failed to take screenshot: %v", err)
	}

	t.Log("✅ DataStar basic test passed!")
}

// TestDataStarReactivity demonstrates DataStar reactivity testing
func TestDataStarReactivity(t *testing.T) {
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

	// Create reactive DataStar page
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>DataStar Reactivity Test</title>
		<script type="module" defer src="https://cdn.jsdelivr.net/npm/@starfederation/datastar@latest"></script>
	</head>
	<body>
		<div data-store='{"name": "", "email": "", "isValid": false}'>
			<form>
				<input id="name-input" type="text" data-model="$name" placeholder="Name">
				<input id="email-input" type="email" data-model="$email" placeholder="Email">
				<button id="validate-btn" type="button" 
					data-on-click="$isValid = $name.length > 0 && $email.includes('@')">
					Validate
				</button>
			</form>
			
			<div id="validation-result" data-show="$isValid" style="color: green;">
				Form is valid!
			</div>
			
			<div id="name-display" data-text="'Hello, ' + $name"></div>
			<div id="email-display" data-text="'Email: ' + $email"></div>
		</div>
	</body>
	</html>
	`

	dataURL := "data:text/html;charset=utf-8," + htmlContent
	page.Navigate(dataURL)

	datastar := frameworks.NewDataStarHelper(page)

	// Wait for DataStar
	if err := datastar.WaitForDataStar(10 * time.Second); err != nil {
		t.Fatalf("DataStar not loaded: %v", err)
	}

	// Test form reactivity
	formData := map[string]string{
		"#name-input":  "John Doe",
		"#email-input": "john@example.com",
	}

	err = page.FillForm("form", formData)
	if err != nil {
		t.Fatalf("Failed to fill form: %v", err)
	}

	// Trigger validation
	err = datastar.TriggerDataStarEvent("#validate-btn", "#validation-result")
	if err != nil {
		t.Fatalf("Failed to trigger validation: %v", err)
	}

	// Verify reactive updates
	nameDisplay, err := page.WaitForElement("#name-display", 5*time.Second)
	if err != nil {
		t.Fatalf("Name display not found: %v", err)
	}

	nameText := nameDisplay.MustText()
	if nameText != "Hello, John Doe" {
		t.Errorf("Expected 'Hello, John Doe', got '%s'", nameText)
	}

	emailDisplay, err := page.WaitForElement("#email-display", 5*time.Second)
	if err != nil {
		t.Fatalf("Email display not found: %v", err)
	}

	emailText := emailDisplay.MustText()
	if emailText != "Email: john@example.com" {
		t.Errorf("Expected 'Email: john@example.com', got '%s'", emailText)
	}

	t.Log("✅ DataStar reactivity test passed!")
}

// TestDataStarStoreManipulation demonstrates direct store manipulation
func TestDataStarStoreManipulation(t *testing.T) {
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

	// Simple DataStar page
	htmlContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>DataStar Store Test</title>
		<script type="module" defer src="https://cdn.jsdelivr.net/npm/@starfederation/datastar@latest"></script>
	</head>
	<body>
		<div data-store='{"items": [], "total": 0}'>
			<div id="item-count" data-text="$items.length + ' items'"></div>
			<div id="total-display" data-text="'Total: $' + $total"></div>
			<ul id="item-list" data-for="item in $items">
				<li data-text="item.name + ' - $' + item.price"></li>
			</ul>
		</div>
	</body>
	</html>
	`

	dataURL := "data:text/html;charset=utf-8," + htmlContent
	page.Navigate(dataURL)

	datastar := frameworks.NewDataStarHelper(page)

	// Wait for DataStar
	if err := datastar.WaitForDataStar(10 * time.Second); err != nil {
		t.Fatalf("DataStar not loaded: %v", err)
	}

	// Test direct store manipulation
	items := []map[string]interface{}{
		{"name": "Apple", "price": 1.50},
		{"name": "Banana", "price": 0.75},
		{"name": "Orange", "price": 2.00},
	}

	err = datastar.SetStoreValue("items", items)
	if err != nil {
		t.Fatalf("Failed to set items: %v", err)
	}

	err = datastar.SetStoreValue("total", 4.25)
	if err != nil {
		t.Fatalf("Failed to set total: %v", err)
	}

	// Wait for updates to propagate
	time.Sleep(1 * time.Second)

	// Verify UI updates
	itemCount, err := page.WaitForElement("#item-count", 5*time.Second)
	if err != nil {
		t.Fatalf("Item count not found: %v", err)
	}

	countText := itemCount.MustText()
	if countText != "3 items" {
		t.Errorf("Expected '3 items', got '%s'", countText)
	}

	totalDisplay, err := page.WaitForElement("#total-display", 5*time.Second)
	if err != nil {
		t.Fatalf("Total display not found: %v", err)
	}

	totalText := totalDisplay.MustText()
	if totalText != "Total: $4.25" {
		t.Errorf("Expected 'Total: $4.25', got '%s'", totalText)
	}

	// Debug store state
	if err := datastar.DebugDataStarState(); err != nil {
		t.Logf("Failed to debug store state: %v", err)
	}

	t.Log("✅ DataStar store manipulation test passed!")
}
