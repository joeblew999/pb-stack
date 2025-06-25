package main

import (
	"fmt"
	"log"
	"time"

	"gui-robot/pkg/robot"
)

// guiControlMain runs the GUI control demonstration
func guiControlMain() {
	fmt.Println("🎯 GUI Robot - Real Application Control Demo")
	fmt.Println("============================================")

	// Create AI controller
	controller := robot.NewAIController()

	// Get screen info first
	screenInfo, err := controller.GetScreenInfo()
	if err != nil {
		log.Fatalf("Failed to get screen info: %v", err)
	}

	fmt.Printf("Screen size: %dx%d\n", screenInfo.Width, screenInfo.Height)
	centerX := float64(screenInfo.Width / 2)
	centerY := float64(screenInfo.Height / 2)

	fmt.Println("\n🚀 Starting GUI Control Demo...")
	fmt.Println("This will demonstrate controlling a GUI application")
	fmt.Println("Make sure you have a text editor or browser open!")

	// Wait for user to get ready
	fmt.Println("\nPress Enter when ready...")
	fmt.Scanln()

	// Demo 1: Open Spotlight Search (macOS)
	fmt.Println("\n1️⃣ Opening Spotlight Search...")
	_, err = controller.Execute("key_press", map[string]interface{}{
		"key": []string{"cmd", "space"},
	})
	if err != nil {
		log.Printf("Failed to open Spotlight: %v", err)
	} else {
		fmt.Println("✓ Spotlight opened")
	}

	time.Sleep(1 * time.Second)

	// Demo 2: Type application name
	fmt.Println("\n2️⃣ Typing 'TextEdit'...")
	_, err = controller.Execute("type", map[string]interface{}{
		"text": "TextEdit",
	})
	if err != nil {
		log.Printf("Failed to type: %v", err)
	} else {
		fmt.Println("✓ Typed 'TextEdit'")
	}

	time.Sleep(1 * time.Second)

	// Demo 3: Press Enter to launch
	fmt.Println("\n3️⃣ Pressing Enter to launch...")
	_, err = controller.Execute("key_press", map[string]interface{}{
		"key": "enter",
	})
	if err != nil {
		log.Printf("Failed to press Enter: %v", err)
	} else {
		fmt.Println("✓ Pressed Enter")
	}

	time.Sleep(3 * time.Second) // Wait for app to launch

	// Demo 4: Type some content
	fmt.Println("\n4️⃣ Typing content in the application...")
	content := `Hello! This text is being typed by GUI Robot!

This demonstrates:
- Opening applications via Spotlight
- Typing text content
- Controlling GUI applications programmatically

Current time: ` + time.Now().Format("2006-01-02 15:04:05")

	_, err = controller.Execute("type", map[string]interface{}{
		"text": content,
	})
	if err != nil {
		log.Printf("Failed to type content: %v", err)
	} else {
		fmt.Println("✓ Content typed successfully")
	}

	time.Sleep(2 * time.Second)

	// Demo 5: Select all text
	fmt.Println("\n5️⃣ Selecting all text (Cmd+A)...")
	_, err = controller.Execute("key_press", map[string]interface{}{
		"key": []string{"cmd", "a"},
	})
	if err != nil {
		log.Printf("Failed to select all: %v", err)
	} else {
		fmt.Println("✓ All text selected")
	}

	time.Sleep(1 * time.Second)

	// Demo 6: Copy text
	fmt.Println("\n6️⃣ Copying text (Cmd+C)...")
	_, err = controller.Execute("key_press", map[string]interface{}{
		"key": []string{"cmd", "c"},
	})
	if err != nil {
		log.Printf("Failed to copy: %v", err)
	} else {
		fmt.Println("✓ Text copied to clipboard")
	}

	time.Sleep(1 * time.Second)

	// Demo 7: Move cursor to end and add more text
	fmt.Println("\n7️⃣ Moving to end and adding more text...")
	_, err = controller.Execute("key_press", map[string]interface{}{
		"key": []string{"cmd", "right"}, // Move to end
	})
	if err != nil {
		log.Printf("Failed to move cursor: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	_, err = controller.Execute("type", map[string]interface{}{
		"text": "\n\n--- Added by GUI Robot ---\nThis line was added after copying the text above!",
	})
	if err != nil {
		log.Printf("Failed to add text: %v", err)
	} else {
		fmt.Println("✓ Additional text added")
	}

	// Demo 8: Mouse control - click somewhere
	fmt.Println("\n8️⃣ Demonstrating mouse control...")
	fmt.Printf("Clicking at center of screen (%.0f, %.0f)...\n", centerX, centerY)

	_, err = controller.Execute("click", map[string]interface{}{
		"x":      centerX,
		"y":      centerY,
		"button": "left",
	})
	if err != nil {
		log.Printf("Failed to click: %v", err)
	} else {
		fmt.Println("✓ Mouse click performed")
	}

	time.Sleep(1 * time.Second)

	// Demo 9: Save the document
	fmt.Println("\n9️⃣ Saving document (Cmd+S)...")
	_, err = controller.Execute("key_press", map[string]interface{}{
		"key": []string{"cmd", "s"},
	})
	if err != nil {
		log.Printf("Failed to save: %v", err)
	} else {
		fmt.Println("✓ Save dialog opened")
	}

	time.Sleep(1 * time.Second)

	// Type filename
	_, err = controller.Execute("type", map[string]interface{}{
		"text": "GUI_Robot_Demo_" + time.Now().Format("20060102_150405"),
	})
	if err != nil {
		log.Printf("Failed to type filename: %v", err)
	} else {
		fmt.Println("✓ Filename entered")
	}

	time.Sleep(1 * time.Second)

	// Press Enter to save
	_, err = controller.Execute("key_press", map[string]interface{}{
		"key": "enter",
	})
	if err != nil {
		log.Printf("Failed to confirm save: %v", err)
	} else {
		fmt.Println("✓ Document saved")
	}

	fmt.Println("\n🎉 GUI Control Demo Complete!")
	fmt.Println("\nWhat was demonstrated:")
	fmt.Println("✅ Opening applications via keyboard shortcuts")
	fmt.Println("✅ Typing text content")
	fmt.Println("✅ Text selection and clipboard operations")
	fmt.Println("✅ Cursor navigation")
	fmt.Println("✅ Mouse clicking")
	fmt.Println("✅ File saving operations")
	fmt.Println("\n💡 This shows how AI can fully control GUI applications!")
}
