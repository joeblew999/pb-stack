package tools

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// ExecuteTask runs the embedded Task functionality
func ExecuteTask(args []string) {
	// For now, delegate to system task command
	// TODO: Embed actual go-task/task functionality
	cmd := exec.Command("task", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		log.Printf("Task execution failed: %v", err)
		os.Exit(1)
	}
}

// ExecuteWhich finds binary location (cross-platform)
func ExecuteWhich(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: xtask which <binary>")
		os.Exit(1)
	}

	binary := args[0]
	path, err := exec.LookPath(binary)
	if err != nil {
		fmt.Printf("❌ %s not found\n", binary)
		os.Exit(1)
	}

	fmt.Printf("✅ %s found at: %s\n", binary, path)
}

// ExecuteGot downloads files (cross-platform)
func ExecuteGot(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: xtask got <url> [-o output]")
		os.Exit(1)
	}

	url := args[0]
	output := ""

	// Parse output flag
	for i, arg := range args {
		if arg == "-o" && i+1 < len(args) {
			output = args[i+1]
			break
		}
	}

	if output == "" {
		output = filepath.Base(url)
	}

	fmt.Printf("📥 Downloading %s to %s...\n", url, output)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Download failed: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Download failed: HTTP %d", resp.StatusCode)
		os.Exit(1)
	}

	file, err := os.Create(output)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	// Copy response body to file
	// TODO: Add progress bar and better error handling
	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Downloaded %s\n", output)
}

// ExecuteSilent executes command silently (cross-platform 2>/dev/null)
func ExecuteSilent(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: xtask silent <command> [args...]")
		os.Exit(1)
	}

	cmd := exec.Command(args[0], args[1:]...)

	// Redirect stderr to null (cross-platform)
	if runtime.GOOS == "windows" {
		cmd.Stderr = nil
	} else {
		devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err == nil {
			defer devNull.Close()
			cmd.Stderr = devNull
		}
	}

	// Run command and ignore errors
	_ = cmd.Run()
}

// ExecuteKillPort kills process on port (cross-platform)
func ExecuteKillPort(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: xtask kill-port <port>")
		os.Exit(1)
	}

	portStr := args[0]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid port: %s", portStr)
		os.Exit(1)
	}

	fmt.Printf("🔍 Looking for processes on port %d...\n", port)

	// Get all processes
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Failed to get processes: %v", err)
		os.Exit(1)
	}

	killed := false
	for _, p := range processes {
		// Get process connections
		connections, err := p.Connections()
		if err != nil {
			continue
		}

		for _, conn := range connections {
			if conn.Laddr.Port == uint32(port) {
				name, _ := p.Name()
				pid := p.Pid

				fmt.Printf("🎯 Found process: %s (PID: %d)\n", name, pid)

				if err := p.Kill(); err != nil {
					log.Printf("Failed to kill process %d: %v", pid, err)
				} else {
					fmt.Printf("✅ Killed process %d\n", pid)
					killed = true
				}
			}
		}
	}

	if !killed {
		fmt.Printf("❌ No processes found on port %d\n", port)
	}
}

// ExecuteWaitForPort waits for port to be available
func ExecuteWaitForPort(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: xtask wait-for-port <port> [timeout]")
		os.Exit(1)
	}

	portStr := args[0]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid port: %s", portStr)
		os.Exit(1)
	}

	timeout := 30 * time.Second
	if len(args) > 1 {
		if t, err := time.ParseDuration(args[1]); err == nil {
			timeout = t
		}
	}

	fmt.Printf("⏳ Waiting for port %d (timeout: %v)...\n", port, timeout)

	start := time.Now()
	for time.Since(start) < timeout {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Second)
		if err == nil {
			conn.Close()
			fmt.Printf("✅ Port %d is ready\n", port)
			return
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("❌ Timeout waiting for port %d\n", port)
	os.Exit(1)
}

// ExecuteTree displays directory tree (cross-platform)
func ExecuteTree(args []string) {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	fmt.Printf("📁 Directory tree for: %s\n", path)
	printTree(path, "", true)
}

func printTree(path, prefix string, isLast bool) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	// Print current item
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	fmt.Printf("%s%s%s\n", prefix, connector, info.Name())

	// If it's a directory, print contents
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return
		}

		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		for i, entry := range entries {
			isLastEntry := i == len(entries)-1
			entryPath := filepath.Join(path, entry.Name())
			printTree(entryPath, newPrefix, isLastEntry)
		}
	}
}

// ExecuteHealthCheck performs HTTP health check
func ExecuteHealthCheck(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: xtask health-check <url>")
		os.Exit(1)
	}

	url := args[0]
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	fmt.Printf("🏥 Health checking: %s\n", url)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("✅ Health check passed: HTTP %d\n", resp.StatusCode)
	} else {
		fmt.Printf("⚠️  Health check warning: HTTP %d\n", resp.StatusCode)
		os.Exit(1)
	}
}
