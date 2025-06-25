package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	url     = flag.String("url", "", "URL to request (required)")
	method  = flag.String("method", "GET", "HTTP method")
	timeout = flag.Duration("timeout", 10*time.Second, "Request timeout")
	pretty  = flag.Bool("pretty", false, "Pretty print JSON responses")
	timing  = flag.Bool("timing", false, "Show timing information")
	headers = flag.Bool("headers", false, "Show response headers")
)

func main() {
	flag.Parse()

	if *url == "" {
		fmt.Println("❌ URL is required")
		flag.Usage()
		os.Exit(1)
	}

	start := time.Now()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: *timeout,
	}

	// Create request
	req, err := http.NewRequest(*method, *url, nil)
	if err != nil {
		fmt.Printf("❌ Failed to create request: %v\n", err)
		os.Exit(1)
	}

	// Set user agent
	req.Header.Set("User-Agent", "OpenCloud-HTTP-Test/1.0")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	// Show timing if requested
	if *timing {
		fmt.Printf("⏱️  Request took: %v\n", elapsed)
	}

	// Show status
	statusIcon := "✅"
	if resp.StatusCode >= 400 {
		statusIcon = "❌"
	} else if resp.StatusCode >= 300 {
		statusIcon = "⚠️"
	}

	fmt.Printf("%s %s %s\n", statusIcon, resp.Status, *url)

	// Show headers if requested
	if *headers {
		fmt.Println("\n📋 Response Headers:")
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
		fmt.Println()
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read response: %v\n", err)
		os.Exit(1)
	}

	// Try to pretty print JSON if requested and content type is JSON
	contentType := resp.Header.Get("Content-Type")
	if *pretty && (contentType == "application/json" || contentType == "application/json; charset=utf-8") {
		var jsonData interface{}
		if err := json.Unmarshal(body, &jsonData); err == nil {
			prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
			if err == nil {
				fmt.Println(string(prettyJSON))
				return
			}
		}
	}

	// Print raw response
	fmt.Print(string(body))
	
	// Add newline if response doesn't end with one
	if len(body) > 0 && body[len(body)-1] != '\n' {
		fmt.Println()
	}
}
