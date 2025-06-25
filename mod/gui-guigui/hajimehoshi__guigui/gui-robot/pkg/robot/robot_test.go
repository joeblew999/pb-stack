package robot

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.MouseSpeed != 100 {
		t.Errorf("Expected MouseSpeed to be 100, got %d", config.MouseSpeed)
	}
	
	if config.ClickDelay != 50*time.Millisecond {
		t.Errorf("Expected ClickDelay to be 50ms, got %v", config.ClickDelay)
	}
	
	if config.ImageFormat != "png" {
		t.Errorf("Expected ImageFormat to be 'png', got %s", config.ImageFormat)
	}
	
	if !config.EnableAI {
		t.Error("Expected EnableAI to be true")
	}
}

func TestNewRobot(t *testing.T) {
	robot := New()
	if robot == nil {
		t.Error("Expected robot to be created, got nil")
	}
	
	// Test that it implements the Robot interface
	var _ Robot = robot
}

func TestNewAIController(t *testing.T) {
	controller := NewAIController()
	if controller == nil {
		t.Error("Expected controller to be created, got nil")
	}
	
	// Test that it implements the AIController interface
	var _ AIController = controller
}

func TestAIControllerCapabilities(t *testing.T) {
	controller := NewAIController()
	capabilities := controller.GetCapabilities()
	
	expectedCaps := []string{
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
	
	if len(capabilities) != len(expectedCaps) {
		t.Errorf("Expected %d capabilities, got %d", len(expectedCaps), len(capabilities))
	}
	
	capMap := make(map[string]bool)
	for _, cap := range capabilities {
		capMap[cap] = true
	}
	
	for _, expected := range expectedCaps {
		if !capMap[expected] {
			t.Errorf("Expected capability '%s' not found", expected)
		}
	}
}

func TestKeyToString(t *testing.T) {
	robot := &DefaultRobot{config: DefaultConfig()}
	
	testCases := []struct {
		key      Key
		expected string
	}{
		{KeyEnter, "enter"},
		{KeyEscape, "esc"},
		{KeySpace, "space"},
		{KeyTab, "tab"},
		{KeyUp, "up"},
		{KeyDown, "down"},
		{KeyLeft, "left"},
		{KeyRight, "right"},
		{KeyShift, "shift"},
		{KeyCtrl, "ctrl"},
		{KeyAlt, "alt"},
		{KeyCmd, "cmd"},
		{KeyF1, "f1"},
		{KeyF12, "f12"},
	}
	
	for _, tc := range testCases {
		result := robot.keyToString(tc.key)
		if result != tc.expected {
			t.Errorf("keyToString(%v) = %s, expected %s", tc.key, result, tc.expected)
		}
	}
	
	// Test unknown key
	unknownKey := Key(9999)
	result := robot.keyToString(unknownKey)
	if result != "" {
		t.Errorf("keyToString(unknown) = %s, expected empty string", result)
	}
}

func TestStringToKey(t *testing.T) {
	controller := &DefaultAIController{config: DefaultConfig()}
	
	testCases := []struct {
		keyStr   string
		expected Key
	}{
		{"enter", KeyEnter},
		{"escape", KeyEscape},
		{"esc", KeyEscape},
		{"space", KeySpace},
		{"tab", KeyTab},
		{"up", KeyUp},
		{"down", KeyDown},
		{"left", KeyLeft},
		{"right", KeyRight},
		{"shift", KeyShift},
		{"ctrl", KeyCtrl},
		{"alt", KeyAlt},
		{"cmd", KeyCmd},
		{"f1", KeyF1},
		{"f12", KeyF12},
	}
	
	for _, tc := range testCases {
		result := controller.stringToKey(tc.keyStr)
		if result != tc.expected {
			t.Errorf("stringToKey(%s) = %v, expected %v", tc.keyStr, result, tc.expected)
		}
	}
	
	// Test unknown key
	result := controller.stringToKey("unknown")
	if result != -1 {
		t.Errorf("stringToKey(unknown) = %v, expected -1", result)
	}
}

func TestSessionManagement(t *testing.T) {
	controller := NewAIController()
	
	// Start a session
	session, err := controller.StartSession()
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}
	
	if session.ID == "" {
		t.Error("Expected session ID to be non-empty")
	}
	
	if session.Robot == nil {
		t.Error("Expected session robot to be non-nil")
	}
	
	if session.History == nil {
		t.Error("Expected session history to be initialized")
	}
	
	// End the session
	err = controller.EndSession(session.ID)
	if err != nil {
		t.Errorf("Failed to end session: %v", err)
	}
}

func TestExecuteWaitCommand(t *testing.T) {
	controller := NewAIController()
	
	start := time.Now()
	result, err := controller.Execute("wait", map[string]interface{}{
		"duration": 0.1, // 100ms
	})
	duration := time.Since(start)
	
	if err != nil {
		t.Fatalf("Wait command failed: %v", err)
	}
	
	if !result.Success {
		t.Error("Expected wait command to succeed")
	}
	
	// Check that it actually waited (with some tolerance)
	if duration < 90*time.Millisecond || duration > 200*time.Millisecond {
		t.Errorf("Expected wait duration around 100ms, got %v", duration)
	}
}

func TestExecuteUnknownCommand(t *testing.T) {
	controller := NewAIController()
	
	result, err := controller.Execute("unknown_command", map[string]interface{}{})
	
	if err == nil {
		t.Error("Expected error for unknown command")
	}
	
	if result.Success {
		t.Error("Expected unknown command to fail")
	}
	
	if result.Error == nil {
		t.Error("Expected result.Error to be set")
	}
}

func TestBatchExecution(t *testing.T) {
	controller := NewAIController()
	
	commands := []Command{
		{
			Name:   "wait",
			Params: map[string]interface{}{"duration": 0.01},
		},
		{
			Name:   "get_screen_info",
			Params: map[string]interface{}{},
		},
	}
	
	results, err := controller.ExecuteBatch(commands)
	if err != nil {
		t.Fatalf("Batch execution failed: %v", err)
	}
	
	if len(results) != len(commands) {
		t.Errorf("Expected %d results, got %d", len(commands), len(results))
	}
	
	for i, result := range results {
		if !result.Success {
			t.Errorf("Command %d failed: %s", i, result.Message)
		}
	}
}

func TestBatchExecutionWithError(t *testing.T) {
	controller := NewAIController()
	
	commands := []Command{
		{
			Name:   "wait",
			Params: map[string]interface{}{"duration": 0.01},
		},
		{
			Name:   "unknown_command",
			Params: map[string]interface{}{},
		},
		{
			Name:   "wait",
			Params: map[string]interface{}{"duration": 0.01},
		},
	}
	
	results, err := controller.ExecuteBatch(commands)
	
	// Should fail on the second command
	if err == nil {
		t.Error("Expected batch execution to fail")
	}
	
	// Should have results for the first two commands (including the failed one)
	if len(results) != 2 {
		t.Errorf("Expected 2 results (up to the failed command), got %d", len(results))
	}
	
	// First command should succeed
	if !results[0].Success {
		t.Error("Expected first command to succeed")
	}
	
	// Second command should fail
	if results[1].Success {
		t.Error("Expected second command to fail")
	}
}
