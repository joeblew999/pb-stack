package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	fmt.Println("🎉 Hello from xtask example!")
	fmt.Println("============================")
	fmt.Println()

	// Show system information
	fmt.Printf("📍 Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("🕐 Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("📁 Working Directory: %s\n", getWorkingDir())
	fmt.Printf("🔧 Go Version: %s\n", runtime.Version())
	fmt.Println()

	// Demonstrate cross-platform compatibility
	fmt.Println("✅ This application was built and run using xtask!")
	fmt.Println("🌍 xtask provides cross-platform development tools:")
	fmt.Println("   • Binary detection (which)")
	fmt.Println("   • File downloads (got)")
	fmt.Println("   • Silent execution")
	fmt.Println("   • Port management")
	fmt.Println("   • Health checks")
	fmt.Println("   • Directory trees")
	fmt.Println()

	fmt.Println("🚀 Try running 'task server-start' to test the server features!")
	fmt.Println("🌐 The server provides HTTP API and NATS coordination")
}

func getWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
