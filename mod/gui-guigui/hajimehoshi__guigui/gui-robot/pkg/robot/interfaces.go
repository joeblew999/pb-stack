package robot

import (
	"image"
	"time"
)

// Robot defines the core interface for GUI automation
type Robot interface {
	// Screen capture
	Screenshot() (image.Image, error)
	ScreenshotArea(x, y, width, height int) (image.Image, error)
	ScreenshotWindow(windowID string) (image.Image, error)

	// Mouse control
	MoveMouse(x, y int) error
	GetMousePosition() (int, int, error)
	Click(button MouseButton) error
	DoubleClick(button MouseButton) error
	Drag(fromX, fromY, toX, toY int) error
	Scroll(x, y int, direction ScrollDirection, amount int) error

	// Keyboard control
	KeyPress(key Key) error
	KeyRelease(key Key) error
	KeyCombo(keys ...Key) error
	TypeText(text string) error

	// Window management
	FindWindow(title string) (*Window, error)
	FindWindowByClass(className string) (*Window, error)
	GetActiveWindow() (*Window, error)
	FocusWindow(window *Window) error
	ListWindows() ([]*Window, error)

	// Utility
	Sleep(duration time.Duration)
	GetScreenSize() (int, int, error)
}

// AIController provides high-level commands optimized for AI interaction
type AIController interface {
	// Execute a command and return result with optional screenshot
	Execute(command string, params map[string]interface{}) (*CommandResult, error)
	
	// Session management
	StartSession() (*Session, error)
	EndSession(sessionID string) error
	
	// Capabilities
	GetCapabilities() []string
	GetScreenInfo() (*ScreenInfo, error)
	
	// Visual analysis
	FindElement(elementType string, criteria map[string]interface{}) (*Element, error)
	CompareImages(img1, img2 image.Image) (*ImageComparison, error)
	
	// Batch operations
	ExecuteBatch(commands []Command) ([]*CommandResult, error)
}

// MouseButton represents mouse button types
type MouseButton int

const (
	LeftButton MouseButton = iota
	RightButton
	MiddleButton
)

// ScrollDirection represents scroll directions
type ScrollDirection int

const (
	ScrollUp ScrollDirection = iota
	ScrollDown
	ScrollLeft
	ScrollRight
)

// Key represents keyboard keys
type Key int

const (
	// Common keys
	KeyEnter Key = iota
	KeyEscape
	KeySpace
	KeyTab
	KeyBackspace
	KeyDelete
	
	// Arrow keys
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	
	// Modifier keys
	KeyShift
	KeyCtrl
	KeyAlt
	KeyCmd  // Command key on macOS
	KeyWin  // Windows key
	
	// Function keys
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

// Window represents a GUI window
type Window struct {
	ID       string
	Title    string
	Class    string
	X, Y     int
	Width    int
	Height   int
	IsActive bool
	PID      int
}

// Element represents a UI element found on screen
type Element struct {
	Type     string
	X, Y     int
	Width    int
	Height   int
	Text     string
	Metadata map[string]interface{}
}

// CommandResult represents the result of an executed command
type CommandResult struct {
	Success    bool
	Message    string
	Screenshot image.Image
	Data       map[string]interface{}
	Duration   time.Duration
	Error      error
}

// Command represents a command to be executed
type Command struct {
	Name   string
	Params map[string]interface{}
}

// Session represents an automation session
type Session struct {
	ID        string
	StartTime time.Time
	Robot     Robot
	History   []*CommandResult
}

// ScreenInfo provides information about the screen
type ScreenInfo struct {
	Width       int
	Height      int
	ColorDepth  int
	RefreshRate int
	Displays    []DisplayInfo
}

// DisplayInfo represents information about a display
type DisplayInfo struct {
	ID       int
	X, Y     int
	Width    int
	Height   int
	Primary  bool
	Name     string
}

// ImageComparison represents the result of comparing two images
type ImageComparison struct {
	Similar    bool
	Difference float64
	DiffImage  image.Image
	Regions    []DifferenceRegion
}

// DifferenceRegion represents an area where images differ
type DifferenceRegion struct {
	X, Y   int
	Width  int
	Height int
	Score  float64
}

// Config holds configuration for the robot
type Config struct {
	// Mouse settings
	MouseSpeed    int
	ClickDelay    time.Duration
	DoubleClickDelay time.Duration
	
	// Keyboard settings
	KeyDelay      time.Duration
	TypeDelay     time.Duration
	
	// Screenshot settings
	ImageFormat   string
	ImageQuality  int
	
	// AI settings
	EnableAI      bool
	AITimeout     time.Duration
	
	// Debug settings
	Debug         bool
	LogLevel      string
}
