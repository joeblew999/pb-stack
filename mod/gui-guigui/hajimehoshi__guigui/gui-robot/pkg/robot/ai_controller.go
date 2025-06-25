package robot

import (
	"fmt"
	"image"
	"time"
	"crypto/rand"
	"encoding/hex"
)

// DefaultAIController implements the AIController interface
type DefaultAIController struct {
	robot    Robot
	sessions map[string]*Session
	config   *Config
}

// NewAIController creates a new AI controller
func NewAIController() AIController {
	return NewAIControllerWithRobot(New())
}

// NewAIControllerWithRobot creates a new AI controller with a specific robot
func NewAIControllerWithRobot(robot Robot) AIController {
	return &DefaultAIController{
		robot:    robot,
		sessions: make(map[string]*Session),
		config:   DefaultConfig(),
	}
}

// Execute executes a high-level command
func (c *DefaultAIController) Execute(command string, params map[string]interface{}) (*CommandResult, error) {
	startTime := time.Now()
	result := &CommandResult{
		Success:  false,
		Data:     make(map[string]interface{}),
		Duration: 0,
	}

	defer func() {
		result.Duration = time.Since(startTime)
	}()

	switch command {
	case "screenshot":
		return c.executeScreenshot(params)
	case "click":
		return c.executeClick(params)
	case "type":
		return c.executeType(params)
	case "key_press":
		return c.executeKeyPress(params)
	case "move_mouse":
		return c.executeMoveMouse(params)
	case "scroll":
		return c.executeScroll(params)
	case "find_window":
		return c.executeFindWindow(params)
	case "get_screen_info":
		return c.executeGetScreenInfo(params)
	case "wait":
		return c.executeWait(params)
	default:
		result.Error = fmt.Errorf("unknown command: %s", command)
		result.Message = fmt.Sprintf("Command '%s' is not supported", command)
		return result, result.Error
	}
}

// executeScreenshot takes a screenshot
func (c *DefaultAIController) executeScreenshot(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	// Check if area screenshot is requested
	if x, hasX := params["x"]; hasX {
		if y, hasY := params["y"]; hasY {
			if width, hasWidth := params["width"]; hasWidth {
				if height, hasHeight := params["height"]; hasHeight {
					img, err := c.robot.ScreenshotArea(
						int(x.(float64)), int(y.(float64)),
						int(width.(float64)), int(height.(float64)))
					if err != nil {
						result.Error = err
						result.Message = "Failed to capture area screenshot"
						return result, err
					}
					result.Screenshot = img
					result.Success = true
					result.Message = "Area screenshot captured successfully"
					return result, nil
				}
			}
		}
	}
	
	// Full screenshot
	img, err := c.robot.Screenshot()
	if err != nil {
		result.Error = err
		result.Message = "Failed to capture screenshot"
		return result, err
	}
	
	result.Screenshot = img
	result.Success = true
	result.Message = "Screenshot captured successfully"
	return result, nil
}

// executeClick performs a mouse click
func (c *DefaultAIController) executeClick(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	// Move to position if specified
	if x, hasX := params["x"]; hasX {
		if y, hasY := params["y"]; hasY {
			err := c.robot.MoveMouse(int(x.(float64)), int(y.(float64)))
			if err != nil {
				result.Error = err
				result.Message = "Failed to move mouse"
				return result, err
			}
		}
	}
	
	// Determine button
	button := LeftButton
	if buttonStr, hasButton := params["button"]; hasButton {
		switch buttonStr.(string) {
		case "left":
			button = LeftButton
		case "right":
			button = RightButton
		case "middle":
			button = MiddleButton
		}
	}
	
	// Perform click
	err := c.robot.Click(button)
	if err != nil {
		result.Error = err
		result.Message = "Failed to click"
		return result, err
	}
	
	result.Success = true
	result.Message = "Click performed successfully"
	
	// Take screenshot if requested
	if takeScreenshot, ok := params["screenshot"].(bool); ok && takeScreenshot {
		img, _ := c.robot.Screenshot()
		result.Screenshot = img
	}
	
	return result, nil
}

// executeType types text
func (c *DefaultAIController) executeType(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	text, hasText := params["text"]
	if !hasText {
		result.Error = fmt.Errorf("text parameter is required")
		result.Message = "Missing text parameter"
		return result, result.Error
	}
	
	err := c.robot.TypeText(text.(string))
	if err != nil {
		result.Error = err
		result.Message = "Failed to type text"
		return result, err
	}
	
	result.Success = true
	result.Message = fmt.Sprintf("Typed text: %s", text.(string))
	
	// Take screenshot if requested
	if takeScreenshot, ok := params["screenshot"].(bool); ok && takeScreenshot {
		img, _ := c.robot.Screenshot()
		result.Screenshot = img
	}
	
	return result, nil
}

// executeKeyPress presses a key or key combination
func (c *DefaultAIController) executeKeyPress(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	keyParam, hasKey := params["key"]
	if !hasKey {
		result.Error = fmt.Errorf("key parameter is required")
		result.Message = "Missing key parameter"
		return result, result.Error
	}
	
	// Handle single key or key combination
	switch v := keyParam.(type) {
	case string:
		key := c.stringToKey(v)
		if key == -1 {
			result.Error = fmt.Errorf("unknown key: %s", v)
			result.Message = fmt.Sprintf("Unknown key: %s", v)
			return result, result.Error
		}
		err := c.robot.KeyPress(key)
		if err == nil {
			err = c.robot.KeyRelease(key)
		}
		if err != nil {
			result.Error = err
			result.Message = "Failed to press key"
			return result, err
		}
	case []interface{}:
		keys := make([]Key, len(v))
		for i, keyStr := range v {
			key := c.stringToKey(keyStr.(string))
			if key == -1 {
				result.Error = fmt.Errorf("unknown key: %s", keyStr.(string))
				result.Message = fmt.Sprintf("Unknown key: %s", keyStr.(string))
				return result, result.Error
			}
			keys[i] = key
		}
		err := c.robot.KeyCombo(keys...)
		if err != nil {
			result.Error = err
			result.Message = "Failed to press key combination"
			return result, err
		}
	}
	
	result.Success = true
	result.Message = "Key press performed successfully"
	
	// Take screenshot if requested
	if takeScreenshot, ok := params["screenshot"].(bool); ok && takeScreenshot {
		img, _ := c.robot.Screenshot()
		result.Screenshot = img
	}
	
	return result, nil
}

// executeMoveMouse moves the mouse cursor
func (c *DefaultAIController) executeMoveMouse(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	x, hasX := params["x"]
	y, hasY := params["y"]
	
	if !hasX || !hasY {
		result.Error = fmt.Errorf("x and y parameters are required")
		result.Message = "Missing x or y parameter"
		return result, result.Error
	}
	
	err := c.robot.MoveMouse(int(x.(float64)), int(y.(float64)))
	if err != nil {
		result.Error = err
		result.Message = "Failed to move mouse"
		return result, err
	}
	
	result.Success = true
	result.Message = fmt.Sprintf("Mouse moved to (%d, %d)", int(x.(float64)), int(y.(float64)))
	return result, nil
}

// executeScroll performs scrolling
func (c *DefaultAIController) executeScroll(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	direction := ScrollDown
	if dirStr, hasDir := params["direction"]; hasDir {
		switch dirStr.(string) {
		case "up":
			direction = ScrollUp
		case "down":
			direction = ScrollDown
		case "left":
			direction = ScrollLeft
		case "right":
			direction = ScrollRight
		}
	}
	
	amount := 3
	if amountParam, hasAmount := params["amount"]; hasAmount {
		amount = int(amountParam.(float64))
	}
	
	x, y := 0, 0
	if xParam, hasX := params["x"]; hasX {
		x = int(xParam.(float64))
	}
	if yParam, hasY := params["y"]; hasY {
		y = int(yParam.(float64))
	}
	
	err := c.robot.Scroll(x, y, direction, amount)
	if err != nil {
		result.Error = err
		result.Message = "Failed to scroll"
		return result, err
	}
	
	result.Success = true
	result.Message = "Scroll performed successfully"
	return result, nil
}

// executeFindWindow finds a window
func (c *DefaultAIController) executeFindWindow(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	title, hasTitle := params["title"]
	if !hasTitle {
		result.Error = fmt.Errorf("title parameter is required")
		result.Message = "Missing title parameter"
		return result, result.Error
	}
	
	window, err := c.robot.FindWindow(title.(string))
	if err != nil {
		result.Error = err
		result.Message = "Failed to find window"
		return result, err
	}
	
	result.Success = true
	result.Message = "Window found successfully"
	result.Data["window"] = window
	return result, nil
}

// executeGetScreenInfo gets screen information
func (c *DefaultAIController) executeGetScreenInfo(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	width, height, err := c.robot.GetScreenSize()
	if err != nil {
		result.Error = err
		result.Message = "Failed to get screen info"
		return result, err
	}
	
	screenInfo := &ScreenInfo{
		Width:  width,
		Height: height,
	}
	
	result.Success = true
	result.Message = "Screen info retrieved successfully"
	result.Data["screen_info"] = screenInfo
	return result, nil
}

// executeWait waits for a specified duration
func (c *DefaultAIController) executeWait(params map[string]interface{}) (*CommandResult, error) {
	result := &CommandResult{Data: make(map[string]interface{})}
	
	duration := 1.0 // Default 1 second
	if durationParam, hasDuration := params["duration"]; hasDuration {
		duration = durationParam.(float64)
	}
	
	c.robot.Sleep(time.Duration(duration * float64(time.Second)))
	
	result.Success = true
	result.Message = fmt.Sprintf("Waited for %.2f seconds", duration)
	return result, nil
}

// StartSession starts a new automation session
func (c *DefaultAIController) StartSession() (*Session, error) {
	sessionID := c.generateSessionID()
	session := &Session{
		ID:        sessionID,
		StartTime: time.Now(),
		Robot:     c.robot,
		History:   make([]*CommandResult, 0),
	}
	
	c.sessions[sessionID] = session
	return session, nil
}

// EndSession ends an automation session
func (c *DefaultAIController) EndSession(sessionID string) error {
	delete(c.sessions, sessionID)
	return nil
}

// GetCapabilities returns the list of supported commands
func (c *DefaultAIController) GetCapabilities() []string {
	return []string{
		"screenshot",
		"click",
		"type",
		"key_press",
		"move_mouse",
		"scroll",
		"find_window",
		"get_screen_info",
		"wait",
	}
}

// GetScreenInfo returns screen information
func (c *DefaultAIController) GetScreenInfo() (*ScreenInfo, error) {
	width, height, err := c.robot.GetScreenSize()
	if err != nil {
		return nil, err
	}
	
	return &ScreenInfo{
		Width:  width,
		Height: height,
	}, nil
}

// FindElement finds a UI element (placeholder implementation)
func (c *DefaultAIController) FindElement(elementType string, criteria map[string]interface{}) (*Element, error) {
	// TODO: Implement element finding using image recognition
	return nil, fmt.Errorf("element finding not implemented yet")
}

// CompareImages compares two images (placeholder implementation)
func (c *DefaultAIController) CompareImages(img1, img2 image.Image) (*ImageComparison, error) {
	// TODO: Implement image comparison
	return nil, fmt.Errorf("image comparison not implemented yet")
}

// ExecuteBatch executes multiple commands in sequence
func (c *DefaultAIController) ExecuteBatch(commands []Command) ([]*CommandResult, error) {
	results := make([]*CommandResult, len(commands))
	
	for i, cmd := range commands {
		result, err := c.Execute(cmd.Name, cmd.Params)
		results[i] = result
		if err != nil {
			// Stop on first error
			return results[:i+1], err
		}
	}
	
	return results, nil
}

// Helper methods

// stringToKey converts a string to a Key constant
func (c *DefaultAIController) stringToKey(keyStr string) Key {
	switch keyStr {
	case "enter":
		return KeyEnter
	case "escape", "esc":
		return KeyEscape
	case "space":
		return KeySpace
	case "tab":
		return KeyTab
	case "backspace":
		return KeyBackspace
	case "delete":
		return KeyDelete
	case "up":
		return KeyUp
	case "down":
		return KeyDown
	case "left":
		return KeyLeft
	case "right":
		return KeyRight
	case "shift":
		return KeyShift
	case "ctrl":
		return KeyCtrl
	case "alt":
		return KeyAlt
	case "cmd":
		return KeyCmd
	case "win":
		return KeyWin
	case "f1":
		return KeyF1
	case "f2":
		return KeyF2
	case "f3":
		return KeyF3
	case "f4":
		return KeyF4
	case "f5":
		return KeyF5
	case "f6":
		return KeyF6
	case "f7":
		return KeyF7
	case "f8":
		return KeyF8
	case "f9":
		return KeyF9
	case "f10":
		return KeyF10
	case "f11":
		return KeyF11
	case "f12":
		return KeyF12
	default:
		return -1
	}
}

// generateSessionID generates a unique session ID
func (c *DefaultAIController) generateSessionID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
