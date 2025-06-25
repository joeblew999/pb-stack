# 🤖 GUI Robot Examples

This directory contains comprehensive examples demonstrating GUI Robot's capabilities.

## 📋 Quick Start

### Interactive Examples Menu
```bash
task examples
```
This launches an interactive menu where you can choose which example to run.

### Direct Examples
```bash
# Basic functionality demonstration
task example-basic

# Real GUI application control
task gui-demo

# Basic automation demo
task demo
```

## 📁 File Structure

### `main.go` - Interactive Examples Menu
- **Purpose**: Unified entry point for all examples
- **Features**: Interactive menu system
- **Usage**: `go run examples/main.go` or `task examples`

### `basic_usage.go` - Basic Functionality
- **Purpose**: Demonstrates core robot functionality
- **Features**: 
  - Basic robot API usage
  - AI controller examples
  - Screenshot and mouse control
  - GUI testing scenarios
- **Usage**: `go run examples/basic_usage.go` or `task example-basic`

### `gui_control_demo.go` - Real Application Control
- **Purpose**: Controls actual GUI applications
- **Features**:
  - Opens applications via Spotlight
  - Types content into TextEdit
  - Demonstrates file operations
  - Shows complete automation workflow
- **Usage**: `go run examples/gui_control_demo.go` or `task gui-demo`

## 🎯 Example Categories

### 1️⃣ Basic Usage Examples
- **Robot API**: Direct function calls
- **AI Controller**: High-level command interface
- **Screenshots**: Screen capture functionality
- **Mouse Control**: Movement and clicking
- **Keyboard Input**: Text typing and key combinations

### 2️⃣ GUI Control Demonstrations
- **Application Launching**: Opening apps via system shortcuts
- **Text Manipulation**: Typing, selecting, copying content
- **File Operations**: Saving documents with custom names
- **Workflow Automation**: Complete task sequences

### 3️⃣ API Examples
- **Low-level Functions**: Direct robot interface usage
- **Error Handling**: Proper error management
- **Configuration**: Robot settings and options

### 4️⃣ Batch Commands
- **Command Sequences**: Multiple operations in order
- **Automation Workflows**: Complex task automation
- **Session Management**: Tracking automation sessions

## 🚀 Running Examples

### Prerequisites
```bash
# Build the project first
task build

# Or build the macOS app for better permissions
task build-macos-app
```

### Interactive Mode
```bash
# Start the examples menu
task examples

# Choose from:
# 1 - Basic Usage Example
# 2 - GUI Control Demo  
# 3 - API Examples
# 4 - Batch Commands
# 0 - Exit
```

### Direct Execution
```bash
# Run specific examples directly
go run examples/basic_usage.go
go run examples/gui_control_demo.go
go run examples/main.go
```

## ⚠️ Important Notes

### macOS Permissions
For full functionality (especially screenshots), you may need to:
1. Grant **Accessibility** permissions
2. Grant **Screen Recording** permissions
3. Use the macOS app bundle: `task build-macos-app`

### GUI Control Demo
The `gui_control_demo.go` will actually control your GUI:
- ✅ Opens real applications
- ✅ Types actual content
- ✅ Saves real files
- ⚠️ Make sure you're ready before running!

## 🎮 What Each Example Shows

### Basic Usage (`basic_usage.go`)
- How to create robot instances
- Taking screenshots and saving them
- Getting mouse position and moving cursor
- Using the AI controller interface
- Session management and command execution

### GUI Control (`gui_control_demo.go`)
- Real-world application automation
- Opening apps via Spotlight search
- Typing content into text editors
- File saving operations
- Complete automation workflows

### Interactive Menu (`main.go`)
- User-friendly example selection
- Organized demonstration categories
- Easy navigation between examples

## 💡 Tips for Development

1. **Start with Basic Usage** to understand the API
2. **Use the Interactive Menu** to explore different features
3. **Check the GUI Demo** to see real automation in action
4. **Build the macOS App** for better permission handling
5. **Read the source code** to understand implementation details

## 🔧 Customization

You can modify these examples to:
- Add your own automation scenarios
- Test specific GUI applications
- Create custom command sequences
- Develop application-specific workflows

Happy automating! 🚀
