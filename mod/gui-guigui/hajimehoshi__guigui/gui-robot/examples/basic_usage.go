package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"gui-robot/pkg/robot"
)

// basicUsageMain runs the basic usage examples
func basicUsageMain() {
	// Example 1: Basic Robot Usage
	fmt.Println("=== Basic Robot Usage ===")
	basicRobotExample()

	// Example 2: AI Controller Usage
	fmt.Println("\n=== AI Controller Usage ===")
	aiControllerExample()

	// Example 3: GUI Testing Example
	fmt.Println("\n=== GUI Testing Example ===")
	guiTestingExample()
}

func basicRobotExample() {
	// Create a robot instance
	r := robot.New()

	// Get screen size
	width, height, err := r.GetScreenSize()
	if err != nil {
		log.Printf("Failed to get screen size: %v", err)
		return
	}
	fmt.Printf("Screen size: %dx%d\n", width, height)

	// Take a screenshot
	img, err := r.Screenshot()
	if err != nil {
		log.Printf("Failed to take screenshot: %v", err)
		return
	}

	// Save screenshot
	file, err := os.Create("basic_screenshot.png")
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		log.Printf("Failed to save screenshot: %v", err)
		return
	}
	fmt.Println("Screenshot saved to: basic_screenshot.png")

	// Get current mouse position
	x, y, err := r.GetMousePosition()
	if err != nil {
		log.Printf("Failed to get mouse position: %v", err)
		return
	}
	fmt.Printf("Current mouse position: (%d, %d)\n", x, y)

	// Move mouse to center of screen
	centerX, centerY := width/2, height/2
	if err := r.MoveMouse(centerX, centerY); err != nil {
		log.Printf("Failed to move mouse: %v", err)
		return
	}
	fmt.Printf("Mouse moved to center: (%d, %d)\n", centerX, centerY)
}

func aiControllerExample() {
	// Create AI controller
	controller := robot.NewAIController()

	// Show capabilities
	caps := controller.GetCapabilities()
	fmt.Println("AI Controller capabilities:")
	for _, cap := range caps {
		fmt.Printf("  - %s\n", cap)
	}

	// Take screenshot using AI controller
	result, err := controller.Execute("screenshot", map[string]interface{}{})
	if err != nil {
		log.Printf("AI screenshot failed: %v", err)
		return
	}
	fmt.Printf("AI Screenshot result: %s (took %v)\n", result.Message, result.Duration)

	// Get screen info
	result, err = controller.Execute("get_screen_info", map[string]interface{}{})
	if err != nil {
		log.Printf("Get screen info failed: %v", err)
		return
	}
	fmt.Printf("Screen info result: %s\n", result.Message)
	if screenInfo, ok := result.Data["screen_info"].(*robot.ScreenInfo); ok {
		fmt.Printf("  Screen dimensions: %dx%d\n", screenInfo.Width, screenInfo.Height)
	}

	// Move mouse using AI controller
	result, err = controller.Execute("move_mouse", map[string]interface{}{
		"x": 100.0,
		"y": 100.0,
	})
	if err != nil {
		log.Printf("Move mouse failed: %v", err)
		return
	}
	fmt.Printf("Move mouse result: %s\n", result.Message)
}

func guiTestingExample() {
	// This example shows how to test a GUI application
	controller := robot.NewAIController()

	// Start a session
	session, err := controller.StartSession()
	if err != nil {
		log.Printf("Failed to start session: %v", err)
		return
	}
	fmt.Printf("Started session: %s\n", session.ID)

	// Test scenario: Take screenshot, click somewhere, take another screenshot
	fmt.Println("Executing GUI test scenario...")

	// Step 1: Initial screenshot
	_, err = controller.Execute("screenshot", map[string]interface{}{})
	if err != nil {
		log.Printf("Initial screenshot failed: %v", err)
		return
	}
	fmt.Println("✓ Initial screenshot taken")

	// Step 2: Click at a safe location (center of screen)
	screenInfo, _ := controller.GetScreenInfo()
	centerX := float64(screenInfo.Width / 2)
	centerY := float64(screenInfo.Height / 2)

	result2, err := controller.Execute("click", map[string]interface{}{
		"x":          centerX,
		"y":          centerY,
		"screenshot": true, // Take screenshot after click
	})
	if err != nil {
		log.Printf("Click failed: %v", err)
		return
	}
	fmt.Printf("✓ Click performed at (%.0f, %.0f)\n", centerX, centerY)

	// Save the post-click screenshot
	if result2.Screenshot != nil {
		file, err := os.Create("post_click_screenshot.png")
		if err == nil {
			defer file.Close()
			png.Encode(file, result2.Screenshot)
			fmt.Println("✓ Post-click screenshot saved to: post_click_screenshot.png")
		}
	}

	// Step 3: Type some text (this will type wherever the cursor is)
	_, err = controller.Execute("type", map[string]interface{}{
		"text": "Hello from GUI Robot!",
	})
	if err != nil {
		log.Printf("Type failed: %v", err)
		return
	}
	fmt.Println("✓ Text typed successfully")

	// Step 4: Press Escape key
	_, err = controller.Execute("key_press", map[string]interface{}{
		"key": "escape",
	})
	if err != nil {
		log.Printf("Key press failed: %v", err)
		return
	}
	fmt.Println("✓ Escape key pressed")

	// Step 5: Wait a moment
	_, err = controller.Execute("wait", map[string]interface{}{
		"duration": 1.0, // 1 second
	})
	if err != nil {
		log.Printf("Wait failed: %v", err)
		return
	}
	fmt.Println("✓ Waited 1 second")

	// Step 6: Final screenshot
	result6, err := controller.Execute("screenshot", map[string]interface{}{})
	if err != nil {
		log.Printf("Final screenshot failed: %v", err)
		return
	}
	fmt.Println("✓ Final screenshot taken")

	// Save final screenshot
	if result6.Screenshot != nil {
		file, err := os.Create("final_screenshot.png")
		if err == nil {
			defer file.Close()
			png.Encode(file, result6.Screenshot)
			fmt.Println("✓ Final screenshot saved to: final_screenshot.png")
		}
	}

	// End session
	controller.EndSession(session.ID)
	fmt.Printf("✓ Session %s ended\n", session.ID)

	fmt.Println("\nGUI test scenario completed successfully!")
	fmt.Println("Check the generated screenshot files to see the results.")
}

// Example of batch command execution
func batchCommandExample() {
	controller := robot.NewAIController()

	// Define a sequence of commands
	commands := []robot.Command{
		{
			Name:   "screenshot",
			Params: map[string]interface{}{},
		},
		{
			Name: "move_mouse",
			Params: map[string]interface{}{
				"x": 200.0,
				"y": 200.0,
			},
		},
		{
			Name: "click",
			Params: map[string]interface{}{
				"button": "left",
			},
		},
		{
			Name: "wait",
			Params: map[string]interface{}{
				"duration": 0.5,
			},
		},
		{
			Name:   "screenshot",
			Params: map[string]interface{}{},
		},
	}

	// Execute batch
	results, err := controller.ExecuteBatch(commands)
	if err != nil {
		log.Printf("Batch execution failed: %v", err)
		return
	}

	fmt.Printf("Batch execution completed. %d commands executed.\n", len(results))
	for i, result := range results {
		fmt.Printf("Command %d: %s (success: %t)\n", i+1, result.Message, result.Success)
	}
}
