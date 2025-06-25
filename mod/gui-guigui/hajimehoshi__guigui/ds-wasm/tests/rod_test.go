package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const (
	ServerURL   = "http://localhost:8081"
	WASMUrl     = "http://localhost:8082"
	TodoURL     = "http://localhost:8083"
	TestTimeout = 30 * time.Second
)

// TestMain sets up and tears down the browser for all tests
func TestMain(m *testing.M) {
	// Setup browser launcher
	launcher := launcher.New().
		Headless(true). // Set to false for debugging
		Devtools(false)

	// Launch browser
	url := launcher.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()

	// Store browser in global variable for tests
	testBrowser = browser

	// Run tests
	code := m.Run()

	// Cleanup
	browser.MustClose()
	launcher.Cleanup()

	os.Exit(code)
}

var testBrowser *rod.Browser

// Helper function to create a new page with timeout
func newPage(t *testing.T) *rod.Page {
	page := testBrowser.MustPage()
	page = page.Timeout(TestTimeout)

	// Add cleanup
	t.Cleanup(func() {
		page.MustClose()
	})

	return page
}

// TestServerMode tests the traditional Go HTTP server
func TestServerMode(t *testing.T) {
	page := newPage(t)

	t.Run("LoadServerPage", func(t *testing.T) {
		page.MustNavigate(ServerURL)
		page.MustWaitLoad()

		// Check title
		title := page.MustInfo().Title
		if title == "" {
			t.Error("Expected page title, got empty string")
		}

		// Check for DataStar content
		page.MustElement("body").MustWaitVisible()

		t.Logf("✅ Server page loaded successfully: %s", title)
	})

	t.Run("TestServerHealth", func(t *testing.T) {
		// Navigate to health endpoint
		page.MustNavigate(ServerURL + "/health")
		page.MustWaitLoad()

		// Should get a health response
		body := page.MustElement("body").MustText()
		if body == "" {
			t.Error("Expected health response, got empty body")
		}

		t.Logf("✅ Server health check passed: %s", body)
	})
}

// TestWASMServiceWorker tests the WASM service worker functionality
func TestWASMServiceWorker(t *testing.T) {
	page := newPage(t)

	t.Run("LoadWASMPage", func(t *testing.T) {
		page.MustNavigate(WASMUrl)
		page.MustWaitLoad()

		// Check title
		title := page.MustInfo().Title
		if title != "DataStar WASM Service Worker" {
			t.Errorf("Expected 'DataStar WASM Service Worker', got '%s'", title)
		}

		// Wait for WASM to load
		page.MustElement("#wasm-status").MustWaitVisible()

		// Wait for WASM success message (with timeout)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := rod.Try(func() {
			page.Context(ctx).MustElement(".success").MustWaitVisible()
		})

		if err != nil {
			// Check if there's an error message instead
			errorEl, _ := page.Element(".error")
			if errorEl != nil {
				errorText, _ := errorEl.Text()
				t.Fatalf("WASM failed to load: %s", errorText)
			}
			t.Fatal("WASM did not load within timeout")
		}

		t.Log("✅ WASM Service Worker loaded successfully")
	})

	t.Run("TestServiceWorkerRoutes", func(t *testing.T) {
		page.MustNavigate(WASMUrl)
		page.MustWaitLoad()

		// Wait for WASM to be ready
		page.MustElement(".success").MustWaitVisible()

		// Test home route
		homeBtn := page.MustElement(`button[data-on-click="handleRoute('/')"]`)
		homeBtn.MustClick()

		// Wait for response
		time.Sleep(500 * time.Millisecond)
		responseEl := page.MustElement("#route-response")
		response := responseEl.MustText()

		if response == "" {
			t.Error("Expected route response, got empty")
		}

		t.Logf("✅ Home route test passed: %s", response)

		// Test hello route
		helloBtn := page.MustElement(`button[data-on-click="handleRoute('/hello')"]`)
		helloBtn.MustClick()

		time.Sleep(500 * time.Millisecond)
		response = responseEl.MustText()

		if response == "" {
			t.Error("Expected hello route response, got empty")
		}

		t.Logf("✅ Hello route test passed: %s", response)
	})

	t.Run("TestSSECommunication", func(t *testing.T) {
		page.MustNavigate(WASMUrl)
		page.MustWaitLoad()

		// Wait for WASM to be ready
		page.MustElement(".success").MustWaitVisible()

		// Test SSE request
		sseBtn := page.MustElement(`button[data-on-click="sendSSERequest('/events')"]`)
		sseBtn.MustClick()

		// Wait for SSE response
		time.Sleep(500 * time.Millisecond)
		responseEl := page.MustElement("#sse-response")
		response := responseEl.MustText()

		if response == "" {
			t.Error("Expected SSE response, got empty")
		}

		t.Logf("✅ SSE communication test passed: %s", response)
	})

	t.Run("TestConnectionToggle", func(t *testing.T) {
		page.MustNavigate(WASMUrl)
		page.MustWaitLoad()

		// Wait for WASM to be ready
		page.MustElement(".success").MustWaitVisible()

		// Initially should be disconnected
		connectBtn := page.MustElement(`button[data-on-click="connectToServer()"]`)
		if !connectBtn.MustVisible() {
			t.Error("Expected connect button to be visible initially")
		}

		// Click connect
		connectBtn.MustClick()

		// Wait for state change
		time.Sleep(1 * time.Second)

		// Should now show disconnect button
		disconnectBtn := page.MustElement(`button[data-on-click="disconnectFromServer()"]`)
		if !disconnectBtn.MustVisible() {
			t.Error("Expected disconnect button to be visible after connecting")
		}

		t.Log("✅ Connection toggle test passed")
	})
}

// TestTodoWASM tests the Todo WASM application
func TestTodoWASM(t *testing.T) {
	page := newPage(t)

	t.Run("LoadTodoPage", func(t *testing.T) {
		page.MustNavigate(TodoURL)
		page.MustWaitLoad()

		// Check title
		title := page.MustInfo().Title
		if title != "DataStar Todo - WASM Mode" {
			t.Errorf("Expected 'DataStar Todo - WASM Mode', got '%s'", title)
		}

		// Wait for WASM to load
		page.MustElement("#todoContainer").MustWaitVisible()

		t.Log("✅ Todo WASM page loaded successfully")
	})

	t.Run("TestAddTodo", func(t *testing.T) {
		page.MustNavigate(TodoURL)
		page.MustWaitLoad()

		// Wait for WASM to be ready
		time.Sleep(2 * time.Second)

		// Add a new todo
		input := page.MustElement("#newTodoInput")
		input.MustInput("Test todo from Rod")

		addBtn := page.MustElement(`button[data-on-click="addTodo()"]`)
		addBtn.MustClick()

		// Wait for todo to appear
		time.Sleep(1 * time.Second)

		// Check if todo was added
		todoContainer := page.MustElement("#todoContainer")
		todoText := todoContainer.MustText()

		if todoText == "" {
			t.Error("Expected todo to be added, but container is empty")
		}

		t.Logf("✅ Add todo test passed: %s", todoText)
	})

	t.Run("TestTodoFilters", func(t *testing.T) {
		page.MustNavigate(TodoURL)
		page.MustWaitLoad()

		// Wait for WASM to be ready
		time.Sleep(2 * time.Second)

		// Test filter buttons
		allBtn := page.MustElement(`button[data-on-click="setFilter('all')"]`)
		activeBtn := page.MustElement(`button[data-on-click="setFilter('active')"]`)
		completedBtn := page.MustElement(`button[data-on-click="setFilter('completed')"]`)

		// Click different filters
		activeBtn.MustClick()
		time.Sleep(500 * time.Millisecond)

		completedBtn.MustClick()
		time.Sleep(500 * time.Millisecond)

		allBtn.MustClick()
		time.Sleep(500 * time.Millisecond)

		t.Log("✅ Todo filters test passed")
	})
}

// TestCrossApplication tests interactions between different modes
func TestCrossApplication(t *testing.T) {
	t.Run("CompareServerAndWASM", func(t *testing.T) {
		// This test could compare responses between server and WASM modes
		// to ensure they behave consistently

		serverPage := newPage(t)
		wasmPage := newPage(t)

		// Load both pages
		serverPage.MustNavigate(ServerURL)
		wasmPage.MustNavigate(WASMUrl)

		serverPage.MustWaitLoad()
		wasmPage.MustWaitLoad()

		// Both should load successfully
		serverTitle := serverPage.MustInfo().Title
		wasmTitle := wasmPage.MustInfo().Title

		if serverTitle == "" || wasmTitle == "" {
			t.Error("Both server and WASM pages should load successfully")
		}

		t.Logf("✅ Cross-application test passed - Server: %s, WASM: %s", serverTitle, wasmTitle)
	})
}

// BenchmarkPageLoad benchmarks page loading times
func BenchmarkPageLoad(b *testing.B) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	b.Run("ServerMode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			page := browser.MustPage(ServerURL)
			page.MustWaitLoad()
			page.MustClose()
		}
	})

	b.Run("WASMMode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			page := browser.MustPage(WASMUrl)
			page.MustWaitLoad()
			page.MustClose()
		}
	})
}
