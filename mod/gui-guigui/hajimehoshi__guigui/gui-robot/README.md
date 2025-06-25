# GUI Robot

**AI-Powered GUI Automation and Control System**

GUI Robot is a Go-based system that enables AI assistants and automated systems to visually interact with graphical user interfaces. It provides screen capture, mouse control, keyboard input, and visual analysis capabilities.

## 🎯 Vision

Enable AI assistants to:
- **See** what's on the screen through automated screenshots
- **Control** applications via mouse and keyboard automation  
- **Test** GUI applications with visual feedback
- **Debug** interface issues in real-time
- **Document** applications with automated screenshots

## 📊 Current Status

**✅ WORKING**: Mouse control, keyboard input, screen info, AI integration
**⚠️ PARTIAL**: Screen capture (requires macOS permissions setup)
**🔄 PLANNED**: Window management, visual analysis, remote control

## 🚀 Features

### Core Capabilities
- ✅ **Mouse Control** - Move cursor, click, drag, scroll
- ✅ **Keyboard Input** - Send key presses, key combinations, text input
- ✅ **Screen Information** - Get screen dimensions and display info
- ✅ **AI Integration** - High-level commands designed for AI assistants
- ⚠️ **Screen Capture** - Available but requires macOS permissions setup
- 🔄 **Window Management** - Find, focus, resize application windows (planned)
- 🔄 **Visual Analysis** - Compare images, detect UI elements (planned)
- ✅ **Cross-Platform** - Works on Windows, macOS, and Linux

### Advanced Features
- 🔄 **Real-time Streaming** - Live screen sharing via WebSocket
- 🎮 **Remote Control** - Control GUI applications over network
- 🤖 **AI Integration** - Tools designed for AI assistant interaction
- 📊 **Test Automation** - Automated GUI testing framework
- 📸 **Smart Screenshots** - Automatic element detection and cropping

## 🛠 Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   AI Assistant  │◄──►│   GUI Robot     │◄──►│  Target GUI App │
│                 │    │                 │    │                 │
│ • Commands      │    │ • Screen Capture│    │ • Presentation  │
│ • Analysis      │    │ • Input Control │    │ • Game          │
│ • Feedback      │    │ • Image Analysis│    │ • Any GUI App   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📦 Installation

### Prerequisites
- Go 1.21 or higher
- Platform-specific dependencies:
  - **macOS**: Accessibility permissions for screen capture
  - **Linux**: X11 development libraries
  - **Windows**: No additional dependencies

### Install
```bash
go install github.com/your-org/gui-robot/cmd/gui-robot@latest
```

### Build from Source
```bash
git clone https://github.com/your-org/gui-robot.git
cd gui-robot
go build -o bin/gui-robot ./cmd/gui-robot
```

## 🎮 Quick Start

### Basic Screen Capture
```go
package main

import (
    "github.com/your-org/gui-robot/pkg/robot"
    "image/png"
    "os"
)

func main() {
    // Create robot instance
    r := robot.New()
    
    // Take screenshot
    img, err := r.Screenshot()
    if err != nil {
        panic(err)
    }
    
    // Save to file
    f, _ := os.Create("screenshot.png")
    defer f.Close()
    png.Encode(f, img)
}
```

### Mouse and Keyboard Control
```go
// Move mouse and click
r.MoveMouse(100, 200)
r.Click(robot.LeftButton)

// Type text
r.TypeText("Hello, World!")

// Press key combinations
r.KeyCombo(robot.KeyCmd, robot.KeyC) // Cmd+C on macOS
```

### AI Assistant Integration
```go
// Create AI-friendly controller
controller := robot.NewAIController()

// Execute high-level commands
result, err := controller.Execute("click_button", map[string]interface{}{
    "text": "Next",
    "screenshot": true,
})
```

## 🔧 API Reference

### Core Robot Interface
```go
type Robot interface {
    // Screen capture
    Screenshot() (image.Image, error)
    ScreenshotArea(x, y, width, height int) (image.Image, error)
    
    // Mouse control
    MoveMouse(x, y int) error
    Click(button MouseButton) error
    Drag(fromX, fromY, toX, toY int) error
    Scroll(x, y int, direction ScrollDirection, amount int) error
    
    // Keyboard control
    KeyPress(key Key) error
    KeyRelease(key Key) error
    KeyCombo(keys ...Key) error
    TypeText(text string) error
    
    // Window management
    FindWindow(title string) (*Window, error)
    FocusWindow(window *Window) error
}
```

### AI Controller Interface
```go
type AIController interface {
    Execute(command string, params map[string]interface{}) (*Result, error)
    StartSession() (*Session, error)
    GetCapabilities() []string
}
```

## 🎯 Use Cases

### 1. GUI Testing Automation
```go
// Test a presentation app
controller.Execute("open_app", map[string]interface{}{
    "name": "presentation-viewer",
})
controller.Execute("click_button", map[string]interface{}{
    "text": "Next Slide",
})
controller.Execute("verify_text", map[string]interface{}{
    "expected": "Slide 2 of 10",
})
```

### 2. AI-Driven QA
```go
// AI assistant can now visually test applications
screenshot := controller.Screenshot()
// AI analyzes screenshot and decides next action
nextAction := ai.AnalyzeAndDecide(screenshot)
controller.Execute(nextAction.Command, nextAction.Params)
```

### 3. Remote GUI Control
```go
// Start remote control server
server := robot.NewRemoteServer(":8080")
server.EnableScreenStreaming()
server.EnableInputControl()
server.Start()
```

## 🔒 Security & Permissions

### macOS
- Requires **Accessibility** permissions in System Preferences
- Requires **Screen Recording** permissions for screenshots

### Linux
- Requires X11 server access
- May need to run with appropriate user permissions

### Windows
- May require administrator privileges for some operations
- Windows Defender might flag automation tools

## 🛣 Roadmap

High priority features is to test on Windows. TO do that we will need to move to a windows virtual machine and then code there. It will be important to use golang compile tags to conditionally compile the code ?


### Phase 1: Core Functionality ✅
- [x] Basic screen capture
- [x] Mouse and keyboard control
- [x] Cross-platform support

### Phase 2: AI Integration 🔄
- [ ] AI-friendly command interface
- [ ] Visual element detection
- [ ] Smart screenshot analysis
- [ ] Command result validation

### Phase 3: Advanced Features 📋
- [ ] Real-time screen streaming
- [ ] Remote control capabilities
- [ ] GUI testing framework
- [ ] Performance optimization

### Phase 4: Ecosystem 🌟
- [ ] Plugin system
- [ ] Integration with popular GUI frameworks
- [ ] Cloud-based control
- [ ] Mobile device support

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup
```bash
git clone https://github.com/your-org/gui-robot.git
cd gui-robot
go mod tidy
go test ./...
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [robotgo](https://github.com/go-vgo/robotgo) - Cross-platform automation
- [screenshot](https://github.com/kbinani/screenshot) - Screen capture utilities
- [guigui](https://github.com/hajimehoshi/guigui) - GUI framework inspiration

---

**Built with ❤️ for AI-powered automation**
