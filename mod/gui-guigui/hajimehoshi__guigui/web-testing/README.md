# Web Testing Framework

A **production-ready testing framework** built on [Rod](https://github.com/go-rod/rod) for comprehensive browser automation and web application testing.

## 🎯 **What This Framework Provides**

This framework offers **professional-grade testing utilities** for:
- **Browser automation** with robust error handling
- **Framework-agnostic testing** (DataStar, React, Vue, etc.)
- **WASM application testing** with specialized support
- **Multi-service coordination** for integration testing
- **Visual debugging** with screenshot capture
- **Cross-platform compatibility** for CI/CD pipelines

## 🌟 **Key Features**

### **🤖 Smart Browser Management**
- Automatic browser lifecycle management
- Headless and visible modes
- Screenshot capture on failures
- Robust timeout and retry handling

### **🎯 Framework Support**
- **DataStar** - Reactive UI testing with store validation
- **WASM** - Module loading and worker testing
- **Generic DOM** - Works with any web framework
- **Extensible** - Easy to add new framework support

### **🔧 Developer Experience**
- **Fluent API** for readable test code
- **Built-in debugging** with visual feedback
- **Service management** for complex scenarios
- **Comprehensive examples** and documentation

## 🚀 **Quick Start**

### **Installation**

Add to your Go project:
```bash
# In your go.work or go.mod
go get web-testing
```

### **Basic Usage**

```go
package main

import (
    "testing"
    "time"

    "web-testing/core"
    "web-testing/frameworks"
)

func TestMyApp(t *testing.T) {
    // Create browser
    browser := core.NewBrowserManager(core.DefaultConfig())
    browser.Start()
    defer browser.Stop()

    // Create page
    page, _ := browser.NewPage("http://localhost:8080")
    defer page.Close()

    // Test DataStar app
    datastar := frameworks.NewDataStarHelper(page)
    datastar.WaitForDataStar(10 * time.Second)

    // Test interactions
    datastar.TriggerDataStarEvent("#submit-btn", "#success-msg")

    // Take screenshot
    page.Screenshot("test_result")
}
```

## 📁 **Project Structure**

```
web-testing/
├── core/                     # Core browser management
│   └── browser.go           # Browser lifecycle and page utilities
├── frameworks/              # Framework-specific helpers
│   ├── datastar.go         # DataStar reactive UI testing
│   └── wasm.go             # WASM application testing
├── services/               # Service management
│   └── manager.go          # Multi-service coordination
├── examples/               # Usage examples
│   ├── basic_test.go       # Basic browser automation
│   ├── datastar_test.go    # DataStar testing examples
│   └── simple_test.go      # Simple functionality tests
├── Taskfile.yml           # Build automation
├── go.mod                 # Dependencies
└── README.md              # This file
```

## 🎯 **Core Features**

### **🤖 Browser Management**

```go
// Configure browser
config := core.DefaultConfig()
config.Headless = false        // Visible browser for debugging
config.Screenshots = true      // Enable screenshot capture
config.ScreenshotDir = "debug" // Custom screenshot directory
config.Timeout = 30 * time.Second

// Create and start browser
browser := core.NewBrowserManager(config)
browser.Start()
defer browser.Stop()

// Create pages
page1, _ := browser.NewPage("http://localhost:8080")
page2, _ := browser.NewPage("http://localhost:8081")
```

### **📸 Screenshot & Debugging**

```go
// Automatic screenshots on test steps
page.Screenshot("initial_state")

// Click and capture
page.ClickAndWait("#submit-btn", 1*time.Second)
page.Screenshot("after_submit")

// Form interaction with debugging
page.FillForm("#contact-form", map[string]string{
    "#name":  "John Doe",
    "#email": "john@example.com",
})
page.Screenshot("form_filled")
```

### **🎯 DataStar Testing**

```go
// Create DataStar helper
datastar := frameworks.NewDataStarHelper(page)

// Wait for DataStar to load
datastar.WaitForDataStar(10 * time.Second)

// Test store manipulation
datastar.SetStoreValue("count", 5)
datastar.WaitForStoreValue("count", 5, 5*time.Second)

// Test reactivity
datastar.TestDataStarReactivity(
    "window.datastar.count = 5",
    "[data-text='$count']"
)

// Test form submission
datastar.TestFormSubmission(
    "#todo-form",
    "#todo-input",
    "#submit-btn",
    "#success-message"
)

// Debug store state
datastar.DebugDataStarState()
```

### **🌐 WASM Testing**

```go
// Create WASM helper
wasm := frameworks.NewWASMHelper(page)

// Check WASM support
wasm.WaitForWASMSupport()

// Wait for WASM module to load
wasm.WaitForWASMLoad(10 * time.Second)

// Validate specific WASM module
wasm.ValidateWASMModule("main.wasm")

// Test WASM function calls
result, _ := wasm.CallWASMFunction("calculateSum", 5, 10)

// Test WASM worker
wasm.ValidateWASMWorker("worker.wasm", 10*time.Second)

// Performance testing
duration, _ := wasm.CheckWASMPerformance("heavyFunction", 100)

// Debug WASM state
wasm.DebugWASMState()
```

### **🔧 Service Management**

```go
// Create service manager
services := services.NewServiceManager()

// Register services
services.RegisterService("api", "http://localhost:8080", 12345)
services.RegisterService("wasm", "http://localhost:8081", 12346)

// Wait for all services to be healthy
services.WaitForHealth(30 * time.Second)

// Check individual service
services.CheckHealth("api")

// Auto-discover services
services.AutoDiscoverServices([]int{8080, 8081, 8082, 8083})

// Print status
services.PrintStatus()
```

## 🧪 **Testing Patterns**

### **Cross-Application Testing**

```go
// Test consistency between server and WASM modes
func TestServerVsWASM(t *testing.T) {
    // Create two pages
    serverPage, _ := browser.NewPage("http://localhost:8080")
    wasmPage, _ := browser.NewPage("http://localhost:8081")

    // Create DataStar helpers
    serverDS := frameworks.NewDataStarHelper(serverPage)
    wasmDS := frameworks.NewDataStarHelper(wasmPage)

    // Compare store states
    frameworks.CompareDataStarStates(serverDS, wasmDS, []string{"count", "items"})
}
```

### **Multi-Service Integration**

```go
func TestFullStack(t *testing.T) {
    // Setup services
    services := services.NewServiceManager()
    services.RegisterService("api", "http://localhost:8080", 0)
    services.RegisterService("frontend", "http://localhost:3000", 0)
    services.WaitForHealth(30 * time.Second)

    // Test frontend
    page, _ := browser.NewPage("http://localhost:3000")
    datastar := frameworks.NewDataStarHelper(page)

    // Test API integration
    datastar.TriggerDataStarEvent("#load-data", "#data-loaded")
}
```

## 🎮 **Available Tasks**

```bash
# Development
task init             # Initialize dependencies
task test             # Run all tests
task test-simple      # Run simple tests
task test-datastar    # Run DataStar tests

# Examples
task example-simple   # Simple browser example
task example-datastar # DataStar example

# Utilities
task clean            # Clean test artifacts
task lint             # Run linting
task format           # Format code
task demo             # Show framework demo
```

## 🌟 **Why Use This Framework?**

### **Production Ready**
- **Robust error handling** with retries and timeouts
- **Cross-platform compatibility** for CI/CD
- **Memory management** and resource cleanup
- **Professional logging** and debugging

### **Framework Agnostic**
- **DataStar support** with reactive UI testing
- **WASM testing** for modern web applications
- **Generic DOM testing** for any framework
- **Extensible architecture** for new frameworks

### **Developer Focused**
- **Fluent API** for readable test code
- **Built-in debugging** with screenshots
- **Comprehensive examples** and documentation
- **Service coordination** for integration testing

## 🔗 **Integration Examples**

### **CI/CD Pipeline**

```yaml
# .github/workflows/test.yml
- name: Run Rod Tests
  run: |
    task init
    task test-simple
    task test-datastar
```

### **Docker Testing**

```dockerfile
FROM golang:1.24-alpine
RUN apk add --no-cache chromium
COPY . /app
WORKDIR /app
RUN task init && task test
```

## 📚 **Documentation**

- **Core API**: See `core/browser.go` for browser management
- **DataStar**: See `frameworks/datastar.go` for reactive UI testing
- **WASM**: See `frameworks/wasm.go` for WASM application testing
- **Services**: See `services/manager.go` for service coordination
- **Examples**: See `examples/` directory for usage patterns

## 🎯 **Perfect For**

1. **DataStar Applications** - Reactive UI testing with store validation
2. **WASM Projects** - Module loading and worker testing
3. **Multi-Service Apps** - Integration testing across services
4. **CI/CD Pipelines** - Automated browser testing
5. **Development Debugging** - Visual testing with screenshots

**This framework transforms Rod from a basic browser automation tool into a comprehensive web testing platform for modern applications!** 🚀


# DataStar WASM Rod Testing

This directory contains Rod-based browser automation tests for the DataStar WASM project.

## 🤖 What is Rod?

[Rod](https://github.com/go-rod/rod) is a high-level driver directly based on DevTools Protocol. It's designed for web automation and scraping, making it perfect for testing DataStar applications in real browsers.

## 🎯 Why Rod for DataStar Testing?

Rod is ideal for testing DataStar applications because:

- **Real Browser Testing**: Tests run in actual browsers (Chrome, Firefox, etc.)
- **DataStar Validation**: Can verify reactive UI updates work correctly
- **WASM Testing**: Ensures WASM modules load and function properly
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Screenshot Capture**: Visual debugging and test documentation
- **Headless/Visible**: Run tests with or without browser UI

## 📁 Test Files

| File | Purpose |
|------|---------|
| `rod_test.go` | Main Rod tests for DataStar applications |
| `helpers.go` | Common testing utilities and helper functions |
| `example_test.go` | Simple examples demonstrating Rod capabilities |
| `run_tests.sh` | Script to run tests with proper service setup |

## 🚀 Running Tests

### Quick Start

```bash
# Run all Rod tests (headless mode)
task test-rod-headless

# Run Rod tests with visible browser
task test-rod

# Run all automated tests
task test-all
```

### Manual Testing

```bash
# Run specific test
go test -v ./tests -run TestExample

# Run all tests in package
go test -v ./tests

# Run with visible browser (for debugging)
go test -v ./tests -args -headless=false
```

### Using the Test Script

```bash
# Run with automatic service management
./tests/run_tests.sh

# Run with visible browser
./tests/run_tests.sh --visible
```

## 🧪 Test Categories

### 1. Example Tests (`example_test.go`)
- **TestExample**: Demonstrates basic Rod usage with JavaScript
- **TestRodCapabilities**: Shows Rod's testing capabilities

### 2. DataStar Application Tests (`rod_test.go`)
- **TestServerMode**: Tests traditional Go HTTP server
- **TestWASMServiceWorker**: Tests WASM service worker functionality
- **TestTodoWASM**: Tests Todo WASM application
- **TestCrossApplication**: Compares server vs WASM consistency

### 3. Helper Functions (`helpers.go`)
- **TestHelper**: Common utilities for Rod testing
- **Screenshot capture**: Visual debugging support
- **Element waiting**: Robust element interaction
- **DataStar integration**: Reactive UI testing helpers

## 📋 Test Requirements

### Services Must Be Running

Before running tests, ensure these services are available:

- **Server Mode**: http://localhost:8081
- **WASM Service Worker**: http://localhost:8082  
- **Todo WASM**: http://localhost:8083

### Starting Services

```bash
# Start all services
task dev

# Or start individually
task server        # Port 8081
task wasm          # Port 8082
task wasm-todo     # Port 8083
```

## 🔧 Test Configuration

### Browser Setup

Rod automatically downloads and manages Chrome/Chromium:
- First run downloads browser (~150MB)
- Subsequent runs use cached browser
- Supports both headless and visible modes

### Timeouts

Default timeouts are configured for reliable testing:
- **Test Timeout**: 30 seconds
- **Element Wait**: 5 seconds
- **WASM Load**: 10 seconds

## 📸 Screenshots

Tests can capture screenshots for debugging:
- Saved to `tests/screenshots/`
- Automatically named with test context
- Useful for visual debugging and documentation

## 🐛 Debugging Tests

### Visible Browser Mode

Run tests with visible browser for debugging:

```bash
go test -v ./tests -args -headless=false
```

### Console Logs

Rod can capture browser console logs:
- JavaScript errors
- DataStar debug messages
- WASM loading status

### Screenshots

Capture screenshots at any point:

```go
helper := NewTestHelper(page, t)
helper.Screenshot("test-state")
```

## 🎯 Testing DataStar Features

### Reactive UI Updates

```go
// Test DataStar reactivity
helper.TestDataStarReactivity(
    "window.datastar.count = 5",
    "[data-text='$count']"
)
```

### Form Interactions

```go
// Test form submission
helper.TestFormSubmission(
    "#todo-form",
    "#todo-input", 
    "#submit-btn",
    "New todo item"
)
```

### WASM Loading

```go
// Validate WASM loaded correctly
err := helper.ValidateWASMLoading()
if err != nil {
    t.Fatalf("WASM failed to load: %v", err)
}
```

## 🌟 Best Practices

1. **Wait for Elements**: Always wait for elements before interacting
2. **Use Helpers**: Leverage helper functions for common operations
3. **Screenshot on Failure**: Capture screenshots when tests fail
4. **Test Real Scenarios**: Test actual user workflows
5. **Verify State**: Check both UI and underlying state changes

## 🔄 Continuous Integration

Rod tests can run in CI environments:
- Headless mode for automated testing
- Screenshot capture for failure analysis
- Cross-platform browser testing
- Parallel test execution support

## 📚 Resources

- [Rod Documentation](https://go-rod.github.io/)
- [DataStar Documentation](https://github.com/starfederation/datastar)
- [DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)
- [Browser Automation Best Practices](https://go-rod.github.io/#/best-practices)

## 🎉 Success Metrics

A successful Rod test suite should verify:
- ✅ All DataStar applications load correctly
- ✅ WASM modules initialize without errors
- ✅ Reactive UI updates work as expected
- ✅ User interactions trigger correct responses
- ✅ Cross-application consistency maintained
- ✅ Performance within acceptable limits

**This framework transforms Rod from a basic browser automation tool into a comprehensive web testing platform for modern applications!** 🚀

## WASM Browser Testing

https://github.com/agnivade/wasmbrowsertest can be integrated with this.