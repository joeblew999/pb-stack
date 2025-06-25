package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("🤖 GUI Robot - Examples & Demonstrations")
	fmt.Println("=========================================")
	fmt.Println("")

	for {
		showMenu()
		choice := getUserChoice()

		switch choice {
		case 1:
			fmt.Println("\n" + strings.Repeat("=", 50))
			runBasicUsageExample()
		case 2:
			fmt.Println("\n" + strings.Repeat("=", 50))
			runGUIControlDemo()
		case 3:
			fmt.Println("\n" + strings.Repeat("=", 50))
			runAPIExample()
		case 4:
			fmt.Println("\n" + strings.Repeat("=", 50))
			runBatchExample()
		case 0:
			fmt.Println("\n👋 Goodbye!")
			return
		default:
			fmt.Println("\n❌ Invalid choice. Please try again.")
		}

		fmt.Println("\n" + strings.Repeat("-", 50))
		fmt.Print("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func showMenu() {
	fmt.Println("📋 Available Examples:")
	fmt.Println("")
	fmt.Println("  1️⃣  Basic Usage Example")
	fmt.Println("      • Basic robot functionality")
	fmt.Println("      • AI controller usage")
	fmt.Println("      • Screenshot and mouse control")
	fmt.Println("")
	fmt.Println("  2️⃣  GUI Control Demo")
	fmt.Println("      • Real application control")
	fmt.Println("      • Opens TextEdit and types content")
	fmt.Println("      • Demonstrates full workflow")
	fmt.Println("")
	fmt.Println("  3️⃣  API Examples")
	fmt.Println("      • Low-level robot API")
	fmt.Println("      • Direct function calls")
	fmt.Println("      • Error handling examples")
	fmt.Println("")
	fmt.Println("  4️⃣  Batch Commands")
	fmt.Println("      • Execute multiple commands")
	fmt.Println("      • Command sequences")
	fmt.Println("      • Automation workflows")
	fmt.Println("")
	fmt.Println("  0️⃣  Exit")
	fmt.Println("")
}

func getUserChoice() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("🎯 Choose an example (0-4): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return -1
	}

	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}

	return choice
}

func runBasicUsageExample() {
	fmt.Println("🚀 Running Basic Usage Example...")
	fmt.Println("This example shows basic robot functionality and AI controller usage.")
	fmt.Println("")
	fmt.Println("💡 For now, run: task example-basic")
	fmt.Println("📝 This will execute the full basic usage example")
}

func runGUIControlDemo() {
	fmt.Println("🎯 Running GUI Control Demo...")
	fmt.Println("⚠️  This will actually control GUI applications!")
	fmt.Println("📝 Make sure you have TextEdit or similar app available")
	fmt.Println("")
	fmt.Println("💡 For now, run: task gui-demo")
	fmt.Println("📝 This will execute the full GUI control demonstration")
}

func runAPIExample() {
	fmt.Println("🔧 Running API Examples...")
	fmt.Println("This shows direct usage of the robot API.")
	fmt.Println("")

	// Run API examples
	apiExampleMain()
}

func runBatchExample() {
	fmt.Println("📦 Running Batch Command Example...")
	fmt.Println("This demonstrates executing multiple commands in sequence.")
	fmt.Println("")

	// Run batch examples
	batchExampleMain()
}

// apiExampleMain demonstrates direct API usage
func apiExampleMain() {
	fmt.Println("🔧 Direct API Usage Examples")
	fmt.Println("============================")

	// This would contain direct robot API calls
	fmt.Println("💡 This example shows direct usage of robot functions")
	fmt.Println("📝 Check basic_usage.go for detailed API examples")
}

// batchExampleMain demonstrates batch command execution
func batchExampleMain() {
	fmt.Println("📦 Batch Command Execution")
	fmt.Println("==========================")

	// This would contain batch execution examples
	fmt.Println("💡 This example shows batch command execution")
	fmt.Println("📝 Check basic_usage.go batchCommandExample() for implementation")
}
