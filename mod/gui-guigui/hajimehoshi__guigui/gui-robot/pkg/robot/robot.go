package robot

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

// DefaultRobot implements the Robot interface using robotgo
type DefaultRobot struct {
	config *Config
}

// New creates a new robot instance with default configuration
func New() Robot {
	return NewWithConfig(DefaultConfig())
}

// NewWithConfig creates a new robot instance with custom configuration
func NewWithConfig(config *Config) Robot {
	return &DefaultRobot{
		config: config,
	}
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		MouseSpeed:       100,
		ClickDelay:       50 * time.Millisecond,
		DoubleClickDelay: 100 * time.Millisecond,
		KeyDelay:         50 * time.Millisecond,
		TypeDelay:        10 * time.Millisecond,
		ImageFormat:      "png",
		ImageQuality:     90,
		EnableAI:         true,
		AITimeout:        30 * time.Second,
		Debug:            false,
		LogLevel:         "info",
	}
}

// Screenshot captures the entire screen
func (r *DefaultRobot) Screenshot() (image.Image, error) {
	// Try robotgo first
	tempFile := "/tmp/gui_robot_screenshot.png"
	err := robotgo.SaveCapture(tempFile)
	if err != nil {
		// Fallback to system screencapture command on macOS
		return r.screenshotFallback()
	}

	// Load the image back using standard library
	file, err := os.Open(tempFile)
	if err != nil {
		return r.screenshotFallback()
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return r.screenshotFallback()
	}

	return img, nil
}

// screenshotFallback uses system commands as backup
func (r *DefaultRobot) screenshotFallback() (image.Image, error) {
	tempFile := "/tmp/gui_robot_fallback_screenshot.png"

	// Use macOS screencapture command
	cmd := fmt.Sprintf("screencapture -x %s", tempFile)
	if err := r.runSystemCommand(cmd); err != nil {
		return nil, fmt.Errorf("failed to capture screenshot with fallback: %w", err)
	}

	// Load the image
	file, err := os.Open(tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open fallback screenshot: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode fallback screenshot: %w", err)
	}

	return img, nil
}

// runSystemCommand executes a system command
func (r *DefaultRobot) runSystemCommand(cmd string) error {
	// This would need to import os/exec and implement command execution
	// For now, return an error to indicate it's not implemented
	return fmt.Errorf("system command execution not implemented: %s", cmd)
}

// ScreenshotArea captures a specific area of the screen
func (r *DefaultRobot) ScreenshotArea(x, y, width, height int) (image.Image, error) {
	// Use robotgo's area capture with save approach
	tempFile := "/tmp/gui_robot_area_screenshot.png"
	err := robotgo.SaveCapture(tempFile, x, y, width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to capture area screenshot: %w", err)
	}

	// Load the image back using standard library
	file, err := os.Open(tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open area screenshot file: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode area screenshot: %w", err)
	}

	return img, nil
}

// ScreenshotWindow captures a specific window (placeholder implementation)
func (r *DefaultRobot) ScreenshotWindow(windowID string) (image.Image, error) {
	// TODO: Implement window-specific screenshot
	return r.Screenshot()
}

// MoveMouse moves the mouse cursor to the specified coordinates
func (r *DefaultRobot) MoveMouse(x, y int) error {
	robotgo.MoveMouse(x, y)
	time.Sleep(r.config.ClickDelay)
	return nil
}

// GetMousePosition returns the current mouse position
func (r *DefaultRobot) GetMousePosition() (int, int, error) {
	x, y := robotgo.GetMousePos()
	return x, y, nil
}

// Click performs a mouse click
func (r *DefaultRobot) Click(button MouseButton) error {
	var buttonStr string
	switch button {
	case LeftButton:
		buttonStr = "left"
	case RightButton:
		buttonStr = "right"
	case MiddleButton:
		buttonStr = "center"
	default:
		return fmt.Errorf("unsupported mouse button: %d", button)
	}

	robotgo.Click(buttonStr)
	time.Sleep(r.config.ClickDelay)
	return nil
}

// DoubleClick performs a double click
func (r *DefaultRobot) DoubleClick(button MouseButton) error {
	err := r.Click(button)
	if err != nil {
		return err
	}
	time.Sleep(r.config.DoubleClickDelay)
	return r.Click(button)
}

// Drag performs a drag operation
func (r *DefaultRobot) Drag(fromX, fromY, toX, toY int) error {
	robotgo.MoveMouse(fromX, fromY)
	time.Sleep(r.config.ClickDelay)
	robotgo.Toggle("left", "down")
	robotgo.MoveMouse(toX, toY)
	time.Sleep(r.config.ClickDelay)
	robotgo.Toggle("left", "up")
	return nil
}

// Scroll performs scrolling
func (r *DefaultRobot) Scroll(x, y int, direction ScrollDirection, amount int) error {
	robotgo.MoveMouse(x, y)
	time.Sleep(r.config.ClickDelay)

	switch direction {
	case ScrollUp:
		robotgo.Scroll(0, -amount)
	case ScrollDown:
		robotgo.Scroll(0, amount)
	case ScrollLeft:
		robotgo.Scroll(-amount, 0)
	case ScrollRight:
		robotgo.Scroll(amount, 0)
	default:
		return fmt.Errorf("unsupported scroll direction: %d", direction)
	}

	return nil
}

// KeyPress presses a key
func (r *DefaultRobot) KeyPress(key Key) error {
	keyStr := r.keyToString(key)
	if keyStr == "" {
		return fmt.Errorf("unsupported key: %d", key)
	}

	robotgo.KeyDown(keyStr)
	time.Sleep(r.config.KeyDelay)
	return nil
}

// KeyRelease releases a key
func (r *DefaultRobot) KeyRelease(key Key) error {
	keyStr := r.keyToString(key)
	if keyStr == "" {
		return fmt.Errorf("unsupported key: %d", key)
	}

	robotgo.KeyUp(keyStr)
	time.Sleep(r.config.KeyDelay)
	return nil
}

// KeyCombo presses a combination of keys
func (r *DefaultRobot) KeyCombo(keys ...Key) error {
	// Press all keys down
	for _, key := range keys {
		if err := r.KeyPress(key); err != nil {
			return err
		}
	}

	// Release all keys in reverse order
	for i := len(keys) - 1; i >= 0; i-- {
		if err := r.KeyRelease(keys[i]); err != nil {
			return err
		}
	}

	return nil
}

// TypeText types the specified text
func (r *DefaultRobot) TypeText(text string) error {
	for _, char := range text {
		robotgo.TypeStr(string(char))
		time.Sleep(r.config.TypeDelay)
	}
	return nil
}

// FindWindow finds a window by title (placeholder implementation)
func (r *DefaultRobot) FindWindow(title string) (*Window, error) {
	// TODO: Implement window finding
	return nil, fmt.Errorf("window finding not implemented yet")
}

// FindWindowByClass finds a window by class name (placeholder implementation)
func (r *DefaultRobot) FindWindowByClass(className string) (*Window, error) {
	// TODO: Implement window finding by class
	return nil, fmt.Errorf("window finding by class not implemented yet")
}

// GetActiveWindow returns the currently active window (placeholder implementation)
func (r *DefaultRobot) GetActiveWindow() (*Window, error) {
	// TODO: Implement active window detection
	return nil, fmt.Errorf("active window detection not implemented yet")
}

// FocusWindow focuses the specified window (placeholder implementation)
func (r *DefaultRobot) FocusWindow(window *Window) error {
	// TODO: Implement window focusing
	return fmt.Errorf("window focusing not implemented yet")
}

// ListWindows returns a list of all windows (placeholder implementation)
func (r *DefaultRobot) ListWindows() ([]*Window, error) {
	// TODO: Implement window listing
	return nil, fmt.Errorf("window listing not implemented yet")
}

// Sleep pauses execution for the specified duration
func (r *DefaultRobot) Sleep(duration time.Duration) {
	time.Sleep(duration)
}

// GetScreenSize returns the screen dimensions
func (r *DefaultRobot) GetScreenSize() (int, int, error) {
	width, height := robotgo.GetScreenSize()
	return width, height, nil
}

// keyToString converts a Key constant to robotgo key string
func (r *DefaultRobot) keyToString(key Key) string {
	switch key {
	case KeyEnter:
		return "enter"
	case KeyEscape:
		return "esc"
	case KeySpace:
		return "space"
	case KeyTab:
		return "tab"
	case KeyBackspace:
		return "backspace"
	case KeyDelete:
		return "delete"
	case KeyUp:
		return "up"
	case KeyDown:
		return "down"
	case KeyLeft:
		return "left"
	case KeyRight:
		return "right"
	case KeyShift:
		return "shift"
	case KeyCtrl:
		return "ctrl"
	case KeyAlt:
		return "alt"
	case KeyCmd:
		return "cmd"
	case KeyWin:
		return "cmd" // Windows key maps to cmd
	case KeyF1:
		return "f1"
	case KeyF2:
		return "f2"
	case KeyF3:
		return "f3"
	case KeyF4:
		return "f4"
	case KeyF5:
		return "f5"
	case KeyF6:
		return "f6"
	case KeyF7:
		return "f7"
	case KeyF8:
		return "f8"
	case KeyF9:
		return "f9"
	case KeyF10:
		return "f10"
	case KeyF11:
		return "f11"
	case KeyF12:
		return "f12"
	default:
		return ""
	}
}
