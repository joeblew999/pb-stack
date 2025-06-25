package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
)

// BrowserConfig configures browser behavior for testing
type BrowserConfig struct {
	Headless       bool
	Screenshots    bool
	ScreenshotDir  string
	Timeout        time.Duration
	SlowMotion     time.Duration
	DevTools       bool
	UserAgent      string
}

// DefaultConfig returns sensible defaults for testing
func DefaultConfig() *BrowserConfig {
	return &BrowserConfig{
		Headless:      true,
		Screenshots:   true,
		ScreenshotDir: "screenshots",
		Timeout:       30 * time.Second,
		SlowMotion:    0,
		DevTools:      false,
		UserAgent:     "",
	}
}

// BrowserManager handles browser lifecycle and provides testing utilities
type BrowserManager struct {
	config   *BrowserConfig
	browser  *rod.Browser
	launcher *launcher.Launcher
}

// NewBrowserManager creates a new browser manager with the given configuration
func NewBrowserManager(config *BrowserConfig) *BrowserManager {
	if config == nil {
		config = DefaultConfig()
	}
	
	return &BrowserManager{
		config: config,
	}
}

// Start initializes and starts the browser
func (bm *BrowserManager) Start() error {
	// Create launcher
	bm.launcher = launcher.New().
		Headless(bm.config.Headless).
		Devtools(bm.config.DevTools)
	
	if bm.config.UserAgent != "" {
		bm.launcher = bm.launcher.Set("user-agent", bm.config.UserAgent)
	}
	
	// Launch browser
	url := bm.launcher.MustLaunch()
	bm.browser = rod.New().ControlURL(url).MustConnect()
	
	// Configure browser
	if bm.config.SlowMotion > 0 {
		bm.browser = bm.browser.SlowMotion(bm.config.SlowMotion)
	}
	
	// Create screenshot directory if needed
	if bm.config.Screenshots {
		if err := os.MkdirAll(bm.config.ScreenshotDir, 0755); err != nil {
			return fmt.Errorf("failed to create screenshot directory: %w", err)
		}
	}
	
	return nil
}

// Stop closes the browser and cleans up resources
func (bm *BrowserManager) Stop() error {
	if bm.browser != nil {
		bm.browser.MustClose()
	}
	if bm.launcher != nil {
		bm.launcher.Cleanup()
	}
	return nil
}

// NewPage creates a new page with testing utilities
func (bm *BrowserManager) NewPage(url string) (*PageHelper, error) {
	if bm.browser == nil {
		return nil, fmt.Errorf("browser not started")
	}
	
	page := bm.browser.MustPage()
	
	// Set timeout
	page = page.Timeout(bm.config.Timeout)
	
	// Navigate to URL if provided
	if url != "" {
		page.MustNavigate(url).MustWaitLoad()
	}
	
	return &PageHelper{
		page:    page,
		config:  bm.config,
		browser: bm,
	}, nil
}

// PageHelper provides testing utilities for a specific page
type PageHelper struct {
	page    *rod.Page
	config  *BrowserConfig
	browser *BrowserManager
}

// Page returns the underlying Rod page
func (ph *PageHelper) Page() *rod.Page {
	return ph.page
}

// Navigate navigates to a URL and waits for load
func (ph *PageHelper) Navigate(url string) *PageHelper {
	ph.page.MustNavigate(url).MustWaitLoad()
	return ph
}

// WaitForElement waits for an element to be present and returns it
func (ph *PageHelper) WaitForElement(selector string, timeout time.Duration) (*rod.Element, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	element, err := ph.page.Context(ctx).Element(selector)
	if err != nil {
		return nil, fmt.Errorf("element %s not found: %w", selector, err)
	}
	
	return element, nil
}

// WaitForElementVisible waits for an element to be visible
func (ph *PageHelper) WaitForElementVisible(selector string, timeout time.Duration) (*rod.Element, error) {
	element, err := ph.WaitForElement(selector, timeout)
	if err != nil {
		return nil, err
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	err = element.Context(ctx).WaitVisible()
	if err != nil {
		return nil, fmt.Errorf("element %s not visible: %w", selector, err)
	}
	
	return element, nil
}

// Screenshot captures a screenshot with optional name
func (ph *PageHelper) Screenshot(name string) error {
	if !ph.config.Screenshots {
		return nil
	}
	
	if name == "" {
		name = fmt.Sprintf("screenshot_%d", time.Now().Unix())
	}
	
	filename := filepath.Join(ph.config.ScreenshotDir, name+".png")
	
	data, err := ph.page.Screenshot(true, nil)
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}
	
	err = utils.OutputFile(filename, data)
	if err != nil {
		return fmt.Errorf("failed to save screenshot: %w", err)
	}
	
	return nil
}

// ClickAndWait clicks an element and waits for a response
func (ph *PageHelper) ClickAndWait(selector string, waitTime time.Duration) error {
	element, err := ph.WaitForElementVisible(selector, 5*time.Second)
	if err != nil {
		return err
	}
	
	element.MustClick()
	time.Sleep(waitTime)
	
	return nil
}

// FillForm fills a form with the provided data
func (ph *PageHelper) FillForm(formSelector string, data map[string]string) error {
	for selector, value := range data {
		element, err := ph.WaitForElement(selector, 5*time.Second)
		if err != nil {
			return fmt.Errorf("failed to find form field %s: %w", selector, err)
		}
		
		element.MustSelectAllText().MustInput(value)
	}
	
	return nil
}

// Close closes the page
func (ph *PageHelper) Close() {
	if ph.page != nil {
		ph.page.MustClose()
	}
}
