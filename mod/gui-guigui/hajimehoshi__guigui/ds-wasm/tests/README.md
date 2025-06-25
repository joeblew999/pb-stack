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
