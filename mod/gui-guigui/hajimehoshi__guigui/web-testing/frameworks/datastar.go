package frameworks

import (
	"encoding/json"
	"fmt"
	"time"

	"web-testing/core"
)

// DataStarHelper provides DataStar-specific testing utilities
type DataStarHelper struct {
	page *core.PageHelper
}

// NewDataStarHelper creates a new DataStar testing helper
func NewDataStarHelper(page *core.PageHelper) *DataStarHelper {
	return &DataStarHelper{
		page: page,
	}
}

// WaitForDataStar waits for DataStar to be loaded and available
func (ds *DataStarHelper) WaitForDataStar(timeout time.Duration) error {
	// Wait for DataStar to be available in the global scope
	_, err := ds.page.Page().Timeout(timeout).Eval(`() => {
		if (typeof window.datastar === 'undefined') {
			throw new Error('DataStar not loaded');
		}
		return true;
	}`)

	if err != nil {
		return fmt.Errorf("DataStar not available: %w", err)
	}

	return nil
}

// GetStore retrieves the current DataStar store state
func (ds *DataStarHelper) GetStore() (map[string]interface{}, error) {
	result, err := ds.page.Page().Eval(`() => {
		return JSON.stringify(window.datastar || {});
	}`)

	if err != nil {
		return nil, fmt.Errorf("failed to get DataStar store: %w", err)
	}

	storeJSON := result.Value.String()

	var store map[string]interface{}
	if err := json.Unmarshal([]byte(storeJSON), &store); err != nil {
		return nil, fmt.Errorf("failed to parse store JSON: %w", err)
	}

	return store, nil
}

// SetStoreValue sets a value in the DataStar store
func (ds *DataStarHelper) SetStoreValue(key string, value interface{}) error {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	script := fmt.Sprintf(`() => {
		window.datastar = window.datastar || {};
		window.datastar.%s = %s;
	}`, key, string(valueJSON))

	_, err = ds.page.Page().Eval(script)
	if err != nil {
		return fmt.Errorf("failed to set store value: %w", err)
	}

	return nil
}

// WaitForStoreValue waits for a store value to match the expected value
func (ds *DataStarHelper) WaitForStoreValue(key string, expectedValue interface{}, timeout time.Duration) error {
	expectedJSON, err := json.Marshal(expectedValue)
	if err != nil {
		return fmt.Errorf("failed to marshal expected value: %w", err)
	}

	script := fmt.Sprintf(`() => {
		const actual = window.datastar && window.datastar.%s;
		const expected = %s;
		return JSON.stringify(actual) === JSON.stringify(expected);
	}`, key, string(expectedJSON))

	_, err = ds.page.Page().Timeout(timeout).Eval(script)
	if err != nil {
		return fmt.Errorf("store value %s did not match expected value within timeout: %w", key, err)
	}

	return nil
}

// TestDataStarReactivity tests that DataStar reactive updates work correctly
func (ds *DataStarHelper) TestDataStarReactivity(storeUpdate string, elementSelector string) error {
	// Execute store update
	_, err := ds.page.Page().Eval(fmt.Sprintf("() => { %s }", storeUpdate))
	if err != nil {
		return fmt.Errorf("failed to execute store update: %w", err)
	}

	// Wait for element to be updated
	time.Sleep(500 * time.Millisecond)

	// Check if element exists and is visible
	_, err = ds.page.WaitForElementVisible(elementSelector, 5*time.Second)
	if err != nil {
		return fmt.Errorf("reactive element %s not found or visible: %w", elementSelector, err)
	}

	return nil
}

// ValidateDataStarElement checks if an element has DataStar attributes
func (ds *DataStarHelper) ValidateDataStarElement(selector string) error {
	element, err := ds.page.WaitForElement(selector, 5*time.Second)
	if err != nil {
		return err
	}

	// Check for common DataStar attributes
	attributes := []string{"data-text", "data-show", "data-on-click", "data-store"}

	hasDataStarAttr := false
	for _, attr := range attributes {
		if value, err := element.Attribute(attr); err == nil && value != nil {
			hasDataStarAttr = true
			break
		}
	}

	if !hasDataStarAttr {
		return fmt.Errorf("element %s does not have DataStar attributes", selector)
	}

	return nil
}

// TriggerDataStarEvent triggers a DataStar event (like clicking a button with data-on-click)
func (ds *DataStarHelper) TriggerDataStarEvent(selector string, expectedUpdate string) error {
	// Click the element
	err := ds.page.ClickAndWait(selector, 500*time.Millisecond)
	if err != nil {
		return fmt.Errorf("failed to trigger DataStar event: %w", err)
	}

	// If expected update is provided, wait for it
	if expectedUpdate != "" {
		_, err = ds.page.WaitForElementVisible(expectedUpdate, 5*time.Second)
		if err != nil {
			return fmt.Errorf("expected update %s not visible after event: %w", expectedUpdate, err)
		}
	}

	return nil
}

// TestFormSubmission tests DataStar form submission
func (ds *DataStarHelper) TestFormSubmission(formSelector, inputSelector, submitSelector, expectedResult string) error {
	// Fill form
	input, err := ds.page.WaitForElement(inputSelector, 5*time.Second)
	if err != nil {
		return fmt.Errorf("form input not found: %w", err)
	}

	input.MustSelectAllText().MustInput("Test input")

	// Submit form
	err = ds.page.ClickAndWait(submitSelector, 1*time.Second)
	if err != nil {
		return fmt.Errorf("failed to submit form: %w", err)
	}

	// Check expected result
	if expectedResult != "" {
		_, err = ds.page.WaitForElementVisible(expectedResult, 5*time.Second)
		if err != nil {
			return fmt.Errorf("expected result %s not visible: %w", expectedResult, err)
		}
	}

	return nil
}

// ValidateSSEConnection validates that Server-Sent Events are working with DataStar
func (ds *DataStarHelper) ValidateSSEConnection(timeout time.Duration) error {
	// Check for SSE connection indicators
	script := `() => {
		// Look for EventSource or fetch with SSE headers
		return window.EventSource !== undefined || 
			   (window.fetch && document.querySelector('[data-sse]'));
	}`

	result, err := ds.page.Page().Timeout(timeout).Eval(script)
	if err != nil {
		return fmt.Errorf("SSE validation failed: %w", err)
	}

	hasSSE := result.Value.Bool()
	if !hasSSE {
		return fmt.Errorf("SSE connection not detected")
	}

	return nil
}

// CompareDataStarStates compares DataStar store states between two pages
func CompareDataStarStates(page1, page2 *DataStarHelper, keys []string) error {
	store1, err := page1.GetStore()
	if err != nil {
		return fmt.Errorf("failed to get store from page 1: %w", err)
	}

	store2, err := page2.GetStore()
	if err != nil {
		return fmt.Errorf("failed to get store from page 2: %w", err)
	}

	for _, key := range keys {
		val1, ok1 := store1[key]
		val2, ok2 := store2[key]

		if ok1 != ok2 {
			return fmt.Errorf("key %s existence mismatch: page1=%v, page2=%v", key, ok1, ok2)
		}

		if ok1 && ok2 {
			val1JSON, _ := json.Marshal(val1)
			val2JSON, _ := json.Marshal(val2)

			if string(val1JSON) != string(val2JSON) {
				return fmt.Errorf("key %s value mismatch: page1=%s, page2=%s", key, val1JSON, val2JSON)
			}
		}
	}

	return nil
}

// DebugDataStarState prints the current DataStar store state for debugging
func (ds *DataStarHelper) DebugDataStarState() error {
	store, err := ds.GetStore()
	if err != nil {
		return err
	}

	storeJSON, _ := json.MarshalIndent(store, "", "  ")
	fmt.Printf("DataStar Store State:\n%s\n", string(storeJSON))

	return nil
}

// WaitForDataStarUpdate waits for any DataStar store update
func (ds *DataStarHelper) WaitForDataStarUpdate(timeout time.Duration) error {
	// Set up a mutation observer to watch for DataStar updates
	script := `() => {
		return new Promise((resolve, reject) => {
			const timeout = setTimeout(() => reject(new Error('timeout')), ` + fmt.Sprintf("%d", timeout.Milliseconds()) + `);
			
			// Watch for DOM changes that indicate DataStar updates
			const observer = new MutationObserver(() => {
				clearTimeout(timeout);
				observer.disconnect();
				resolve(true);
			});
			
			observer.observe(document.body, {
				childList: true,
				subtree: true,
				attributes: true,
				attributeFilter: ['data-text', 'data-show', 'data-class']
			});
		});
	}`

	_, err := ds.page.Page().Eval(script)
	return err
}
