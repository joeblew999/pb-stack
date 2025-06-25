package frameworks

import (
	"fmt"
	"strings"
	"time"

	"web-testing/core"
)

// WASMHelper provides WASM-specific testing utilities
type WASMHelper struct {
	page *core.PageHelper
}

// NewWASMHelper creates a new WASM testing helper
func NewWASMHelper(page *core.PageHelper) *WASMHelper {
	return &WASMHelper{
		page: page,
	}
}

// WaitForWASMSupport checks if the browser supports WASM
func (w *WASMHelper) WaitForWASMSupport() error {
	result, err := w.page.Page().Eval(`() => {
		return typeof WebAssembly !== 'undefined' && 
			   typeof WebAssembly.instantiate === 'function';
	}`)

	if err != nil {
		return fmt.Errorf("failed to check WASM support: %w", err)
	}

	supported := result.Value.Bool()
	if !supported {
		return fmt.Errorf("WASM not supported in this browser")
	}

	return nil
}

// WaitForWASMLoad waits for a WASM module to load successfully
func (w *WASMHelper) WaitForWASMLoad(timeout time.Duration) error {
	// Wait for Go WASM runtime to be available
	_, err := w.page.Page().Timeout(timeout).Eval(`() => {
		if (typeof Go === 'undefined') {
			throw new Error('Go WASM runtime not loaded');
		}
		return true;
	}`)

	if err != nil {
		return fmt.Errorf("Go WASM runtime not available: %w", err)
	}

	return nil
}

// ValidateWASMModule checks if a specific WASM module is loaded
func (w *WASMHelper) ValidateWASMModule(moduleName string) error {
	script := fmt.Sprintf(`() => {
		// Check if the WASM module has been instantiated
		return typeof WebAssembly !== 'undefined' && 
			   window.wasmModule !== undefined;
	}`)

	result, err := w.page.Page().Eval(script)
	if err != nil {
		return fmt.Errorf("failed to validate WASM module: %w", err)
	}

	loaded := result.Value.Bool()
	if !loaded {
		return fmt.Errorf("WASM module %s not loaded", moduleName)
	}

	return nil
}

// WaitForWASMFunction waits for a specific WASM function to be available
func (w *WASMHelper) WaitForWASMFunction(functionName string, timeout time.Duration) error {
	script := fmt.Sprintf(`() => {
		if (typeof %s === 'undefined') {
			throw new Error('Function %s not available');
		}
		return true;
	}`, functionName, functionName)

	_, err := w.page.Page().Timeout(timeout).Eval(script)
	if err != nil {
		return fmt.Errorf("WASM function %s not available: %w", functionName, err)
	}

	return nil
}

// CallWASMFunction calls a WASM function and returns the result
func (w *WASMHelper) CallWASMFunction(functionName string, args ...interface{}) (interface{}, error) {
	// Convert args to JavaScript format
	jsArgs := make([]string, len(args))
	for i, arg := range args {
		jsArgs[i] = fmt.Sprintf("%v", arg)
	}

	script := fmt.Sprintf(`() => {
		return %s(%s);
	}`, functionName, strings.Join(jsArgs, ", "))

	result, err := w.page.Page().Eval(script)
	if err != nil {
		return nil, fmt.Errorf("failed to call WASM function %s: %w", functionName, err)
	}

	return result.Value, nil
}

// ValidateWASMWorker checks if a Web Worker with WASM is functioning
func (w *WASMHelper) ValidateWASMWorker(workerPath string, timeout time.Duration) error {
	script := fmt.Sprintf(`() => {
		return new Promise((resolve, reject) => {
			const timeout = setTimeout(() => reject(new Error('Worker timeout')), %d);
			
			const worker = new Worker('%s');
			
			worker.onmessage = function(e) {
				clearTimeout(timeout);
				worker.terminate();
				resolve(e.data);
			};
			
			worker.onerror = function(error) {
				clearTimeout(timeout);
				worker.terminate();
				reject(error);
			};
			
			// Send test message
			worker.postMessage('test');
		});
	}`, timeout.Milliseconds(), workerPath)

	_, err := w.page.Page().Eval(script)
	if err != nil {
		return fmt.Errorf("WASM worker validation failed: %w", err)
	}

	return nil
}

// CheckWASMPerformance measures WASM execution performance
func (w *WASMHelper) CheckWASMPerformance(functionName string, iterations int) (time.Duration, error) {
	script := fmt.Sprintf(`() => {
		const start = performance.now();
		
		for (let i = 0; i < %d; i++) {
			%s();
		}
		
		const end = performance.now();
		return end - start;
	}`, iterations, functionName)

	result, err := w.page.Page().Eval(script)
	if err != nil {
		return 0, fmt.Errorf("failed to measure WASM performance: %w", err)
	}

	ms := result.Value.Num()

	return time.Duration(ms) * time.Millisecond, nil
}

// ValidateWASMMemory checks WASM memory usage
func (w *WASMHelper) ValidateWASMMemory() (map[string]interface{}, error) {
	script := `() => {
		const memInfo = {};
		
		// Check WebAssembly memory if available
		if (typeof WebAssembly !== 'undefined' && window.wasmMemory) {
			memInfo.wasmMemory = {
				buffer: window.wasmMemory.buffer.byteLength,
				pages: window.wasmMemory.buffer.byteLength / 65536
			};
		}
		
		// Check JavaScript memory if available
		if (performance.memory) {
			memInfo.jsMemory = {
				used: performance.memory.usedJSHeapSize,
				total: performance.memory.totalJSHeapSize,
				limit: performance.memory.jsHeapSizeLimit
			};
		}
		
		return memInfo;
	}`

	result, err := w.page.Page().Eval(script)
	if err != nil {
		return nil, fmt.Errorf("failed to get WASM memory info: %w", err)
	}

	memInfo := make(map[string]interface{})
	for k, v := range result.Value.Map() {
		memInfo[k] = v
	}
	return memInfo, nil
}

// WaitForWASMError waits for and captures WASM errors
func (w *WASMHelper) WaitForWASMError(timeout time.Duration) (string, error) {
	script := fmt.Sprintf(`() => {
		return new Promise((resolve, reject) => {
			const timeout = setTimeout(() => reject(new Error('No WASM error within timeout')), %d);
			
			const originalError = window.onerror;
			window.onerror = function(message, source, lineno, colno, error) {
				clearTimeout(timeout);
				window.onerror = originalError;
				resolve(message);
			};
			
			// Also listen for unhandled promise rejections
			const originalRejection = window.onunhandledrejection;
			window.onunhandledrejection = function(event) {
				clearTimeout(timeout);
				window.onunhandledrejection = originalRejection;
				resolve(event.reason.toString());
			};
		});
	}`, timeout.Milliseconds())

	result, err := w.page.Page().Eval(script)
	if err != nil {
		return "", fmt.Errorf("failed to wait for WASM error: %w", err)
	}

	errorMsg := result.Value.String()

	return errorMsg, nil
}

// TestWASMServiceWorker tests WASM functionality in a Service Worker
func (w *WASMHelper) TestWASMServiceWorker(swPath string, timeout time.Duration) error {
	script := fmt.Sprintf(`() => {
		return new Promise((resolve, reject) => {
			const timeout = setTimeout(() => reject(new Error('Service Worker timeout')), %d);
			
			if ('serviceWorker' in navigator) {
				navigator.serviceWorker.register('%s')
					.then(registration => {
						clearTimeout(timeout);
						resolve(registration);
					})
					.catch(error => {
						clearTimeout(timeout);
						reject(error);
					});
			} else {
				clearTimeout(timeout);
				reject(new Error('Service Workers not supported'));
			}
		});
	}`, timeout.Milliseconds(), swPath)

	_, err := w.page.Page().Eval(script)
	if err != nil {
		return fmt.Errorf("WASM Service Worker test failed: %w", err)
	}

	return nil
}

// DebugWASMState prints current WASM state for debugging
func (w *WASMHelper) DebugWASMState() error {
	script := `() => {
		const state = {
			wasmSupported: typeof WebAssembly !== 'undefined',
			goRuntime: typeof Go !== 'undefined',
			wasmModule: typeof window.wasmModule !== 'undefined',
			serviceWorker: 'serviceWorker' in navigator,
			webWorkers: typeof Worker !== 'undefined'
		};
		
		return JSON.stringify(state, null, 2);
	}`

	result, err := w.page.Page().Eval(script)
	if err != nil {
		return fmt.Errorf("failed to get WASM debug state: %w", err)
	}

	state := result.Value.String()

	fmt.Printf("WASM Debug State:\n%s\n", state)
	return nil
}
