package tests

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
)

// TestHelper provides common testing utilities for Rod-based tests
type TestHelper struct {
	page *rod.Page
	t    *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(page *rod.Page, t *testing.T) *TestHelper {
	return &TestHelper{
		page: page,
		t:    t,
	}
}

// WaitForWASM waits for WASM module to load successfully
func (h *TestHelper) WaitForWASM(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return rod.Try(func() {
		h.page.Context(ctx).MustElement(".success").MustWaitVisible()
	})
}

// WaitForElement waits for an element to be visible
func (h *TestHelper) WaitForElement(selector string, timeout time.Duration) (*rod.Element, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var element *rod.Element
	err := rod.Try(func() {
		element = h.page.Context(ctx).MustElement(selector).MustWaitVisible()
	})

	return element, err
}

// ClickAndWait clicks an element and waits for a response
func (h *TestHelper) ClickAndWait(selector string, waitTime time.Duration) error {
	element, err := h.WaitForElement(selector, 5*time.Second)
	if err != nil {
		return fmt.Errorf("element not found: %s", selector)
	}

	element.MustClick()
	time.Sleep(waitTime)
	return nil
}

// GetElementText gets text from an element with error handling
func (h *TestHelper) GetElementText(selector string) (string, error) {
	element, err := h.page.Element(selector)
	if err != nil {
		return "", fmt.Errorf("element not found: %s", selector)
	}

	text, err := element.Text()
	if err != nil {
		return "", fmt.Errorf("failed to get text from %s: %v", selector, err)
	}

	return text, nil
}

// Screenshot takes a screenshot and saves it to the tests directory
func (h *TestHelper) Screenshot(filename string) error {
	// Take screenshot
	screenshot, err := h.page.Screenshot(true, nil)
	if err != nil {
		return fmt.Errorf("failed to take screenshot: %v", err)
	}

	// Save screenshot
	screenshotPath := filepath.Join("screenshots", filename+".png")
	err = utils.OutputFile(screenshotPath, screenshot)
	if err != nil {
		return fmt.Errorf("failed to save screenshot: %v", err)
	}

	h.t.Logf("📸 Screenshot saved: %s", screenshotPath)
	return nil
}

// CheckConsoleErrors checks for JavaScript console errors (simplified)
func (h *TestHelper) CheckConsoleErrors() []string {
	// Simplified version - just return empty for now
	// In a real implementation, you'd set up event listeners
	return []string{}
}

// WaitForDataStarUpdate waits for DataStar to update the UI
func (h *TestHelper) WaitForDataStarUpdate(timeout time.Duration) error {
	// Simplified - just wait for a short time
	time.Sleep(500 * time.Millisecond)
	return nil
}

// TestDataStarReactivity tests DataStar reactive features
func (h *TestHelper) TestDataStarReactivity(storeUpdate string, expectedSelector string) error {
	// Execute JavaScript to update DataStar store
	_, err := h.page.Eval(fmt.Sprintf(`
		if (window.datastar) {
			%s
		}
	`, storeUpdate))

	if err != nil {
		return fmt.Errorf("failed to update DataStar store: %v", err)
	}

	// Wait for UI to update
	time.Sleep(500 * time.Millisecond)

	// Check if expected element is visible
	_, err = h.WaitForElement(expectedSelector, 5*time.Second)
	return err
}

// SimulateUserInteraction simulates realistic user interactions
func (h *TestHelper) SimulateUserInteraction(actions []UserAction) error {
	for _, action := range actions {
		switch action.Type {
		case "click":
			err := h.ClickAndWait(action.Selector, action.Wait)
			if err != nil {
				return err
			}
		case "input":
			element, err := h.WaitForElement(action.Selector, 5*time.Second)
			if err != nil {
				return err
			}
			element.MustInput(action.Value)
			time.Sleep(action.Wait)
		case "wait":
			time.Sleep(action.Wait)
		case "screenshot":
			err := h.Screenshot(action.Value)
			if err != nil {
				h.t.Logf("Warning: Screenshot failed: %v", err)
			}
		}
	}
	return nil
}

// UserAction represents a user interaction
type UserAction struct {
	Type     string        // "click", "input", "wait", "screenshot"
	Selector string        // CSS selector for the element
	Value    string        // Value for input or screenshot filename
	Wait     time.Duration // Time to wait after action
}

// ValidateWASMLoading validates that WASM loaded correctly
func (h *TestHelper) ValidateWASMLoading() error {
	// Check for WASM success indicator
	err := h.WaitForWASM(10 * time.Second)
	if err != nil {
		// Check for error message
		errorText, _ := h.GetElementText(".error")
		if errorText != "" {
			return fmt.Errorf("WASM loading failed: %s", errorText)
		}
		return fmt.Errorf("WASM did not load within timeout")
	}

	// Check console for WASM-related messages
	errors := h.CheckConsoleErrors()
	if len(errors) > 0 {
		h.t.Logf("Console errors detected: %v", errors)
	}

	return nil
}

// TestFormSubmission tests form submission with DataStar
func (h *TestHelper) TestFormSubmission(formSelector, inputSelector, submitSelector string, inputValue string) error {
	// Fill form
	input, err := h.WaitForElement(inputSelector, 5*time.Second)
	if err != nil {
		return fmt.Errorf("input field not found: %s", inputSelector)
	}

	input.MustInput(inputValue)

	// Submit form
	err = h.ClickAndWait(submitSelector, 1*time.Second)
	if err != nil {
		return fmt.Errorf("submit button not found: %s", submitSelector)
	}

	return nil
}

// CompareElements compares text content of two elements
func (h *TestHelper) CompareElements(selector1, selector2 string) (bool, error) {
	text1, err := h.GetElementText(selector1)
	if err != nil {
		return false, err
	}

	text2, err := h.GetElementText(selector2)
	if err != nil {
		return false, err
	}

	return text1 == text2, nil
}

// WaitForNetworkIdle waits for network activity to settle
func (h *TestHelper) WaitForNetworkIdle(timeout time.Duration) error {
	// Simplified - just wait
	time.Sleep(1 * time.Second)
	return nil
}

// GetPagePerformance gets basic page performance metrics
func (h *TestHelper) GetPagePerformance() (map[string]interface{}, error) {
	result, err := h.page.Eval(`() => {
		const timing = performance.timing;
		return {
			loadTime: timing.loadEventEnd - timing.navigationStart,
			domReady: timing.domContentLoadedEventEnd - timing.navigationStart,
			firstPaint: performance.getEntriesByType('paint')[0]?.startTime || 0
		};
	}`)

	if err != nil {
		return nil, err
	}

	// Convert to map[string]interface{}
	resultMap := make(map[string]interface{})
	for k, v := range result.Value.Map() {
		resultMap[k] = v
	}
	return resultMap, nil
}
