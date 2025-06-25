package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gui-robot/pkg/robot"
)

func main() {
	var (
		command    = flag.String("command", "screenshot", "Command to execute")
		params     = flag.String("params", "{}", "Command parameters as JSON")
		output     = flag.String("output", "", "Output file for screenshots")
		interactive = flag.Bool("interactive", false, "Start interactive mode")
		version    = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *version {
		fmt.Println("GUI Robot v1.0.0")
		fmt.Println("AI-Powered GUI Automation and Control System")
		return
	}

	if *interactive {
		runInteractiveMode()
		return
	}

	// Parse parameters
	var paramMap map[string]interface{}
	if err := json.Unmarshal([]byte(*params), &paramMap); err != nil {
		log.Fatalf("Failed to parse parameters: %v", err)
	}

	// Create AI controller
	controller := robot.NewAIController()

	// Execute command
	result, err := controller.Execute(*command, paramMap)
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	// Print result
	fmt.Printf("Success: %t\n", result.Success)
	fmt.Printf("Message: %s\n", result.Message)
	fmt.Printf("Duration: %v\n", result.Duration)

	// Save screenshot if available and output specified
	if result.Screenshot != nil && *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer file.Close()

		if err := png.Encode(file, result.Screenshot); err != nil {
			log.Fatalf("Failed to save screenshot: %v", err)
		}
		fmt.Printf("Screenshot saved to: %s\n", *output)
	}

	// Print additional data
	if len(result.Data) > 0 {
		fmt.Println("Additional data:")
		for key, value := range result.Data {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

func runInteractiveMode() {
	fmt.Println("GUI Robot Interactive Mode")
	fmt.Println("Type 'help' for available commands, 'quit' to exit")
	fmt.Println()

	controller := robot.NewAIController()

	for {
		fmt.Print("gui-robot> ")
		
		var input string
		fmt.Scanln(&input)
		
		if input == "quit" || input == "exit" {
			break
		}
		
		if input == "help" {
			showHelp()
			continue
		}
		
		if input == "capabilities" {
			caps := controller.GetCapabilities()
			fmt.Println("Available commands:")
			for _, cap := range caps {
				fmt.Printf("  - %s\n", cap)
			}
			continue
		}
		
		if input == "screenshot" {
			takeScreenshot(controller)
			continue
		}
		
		if strings.HasPrefix(input, "click ") {
			parts := strings.Fields(input)
			if len(parts) >= 3 {
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])
				performClick(controller, x, y)
			} else {
				fmt.Println("Usage: click <x> <y>")
			}
			continue
		}
		
		if strings.HasPrefix(input, "type ") {
			text := strings.TrimPrefix(input, "type ")
			performType(controller, text)
			continue
		}
		
		if strings.HasPrefix(input, "key ") {
			key := strings.TrimPrefix(input, "key ")
			performKeyPress(controller, key)
			continue
		}
		
		if input == "screen_info" {
			showScreenInfo(controller)
			continue
		}
		
		fmt.Printf("Unknown command: %s\n", input)
		fmt.Println("Type 'help' for available commands")
	}
	
	fmt.Println("Goodbye!")
}

func showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                 - Show this help")
	fmt.Println("  capabilities         - Show available AI commands")
	fmt.Println("  screenshot           - Take a screenshot")
	fmt.Println("  click <x> <y>        - Click at coordinates")
	fmt.Println("  type <text>          - Type text")
	fmt.Println("  key <key>            - Press a key")
	fmt.Println("  screen_info          - Show screen information")
	fmt.Println("  quit/exit            - Exit interactive mode")
}

func takeScreenshot(controller robot.AIController) {
	result, err := controller.Execute("screenshot", map[string]interface{}{})
	if err != nil {
		fmt.Printf("Screenshot failed: %v\n", err)
		return
	}
	
	if result.Screenshot != nil {
		filename := fmt.Sprintf("screenshot_%d.png", time.Now().Unix())
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Failed to create file: %v\n", err)
			return
		}
		defer file.Close()
		
		if err := png.Encode(file, result.Screenshot); err != nil {
			fmt.Printf("Failed to save screenshot: %v\n", err)
			return
		}
		
		fmt.Printf("Screenshot saved to: %s\n", filename)
	}
}

func performClick(controller robot.AIController, x, y int) {
	params := map[string]interface{}{
		"x": float64(x),
		"y": float64(y),
	}
	
	result, err := controller.Execute("click", params)
	if err != nil {
		fmt.Printf("Click failed: %v\n", err)
		return
	}
	
	fmt.Printf("Click result: %s\n", result.Message)
}

func performType(controller robot.AIController, text string) {
	params := map[string]interface{}{
		"text": text,
	}
	
	result, err := controller.Execute("type", params)
	if err != nil {
		fmt.Printf("Type failed: %v\n", err)
		return
	}
	
	fmt.Printf("Type result: %s\n", result.Message)
}

func performKeyPress(controller robot.AIController, key string) {
	params := map[string]interface{}{
		"key": key,
	}
	
	result, err := controller.Execute("key_press", params)
	if err != nil {
		fmt.Printf("Key press failed: %v\n", err)
		return
	}
	
	fmt.Printf("Key press result: %s\n", result.Message)
}

func showScreenInfo(controller robot.AIController) {
	screenInfo, err := controller.GetScreenInfo()
	if err != nil {
		fmt.Printf("Failed to get screen info: %v\n", err)
		return
	}
	
	fmt.Printf("Screen size: %dx%d\n", screenInfo.Width, screenInfo.Height)
}
