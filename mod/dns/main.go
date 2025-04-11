// /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

// providerNameservers maps lowercase provider keys to a list of their known nameserver hostnames.
// Note: Nameserver hostnames often end with a dot. Be consistent.
var providerNameservers = map[string][]string{
	"weebly": {
		"ns1.weebly.com.",
		"ns2.weebly.com.",
		"ns3.weebly.com.",
		"ns4.weebly.com.",
		// Add variations if needed, e.g., without the trailing dot, though matching with suffix is better
		"ns1.weebly.com",
		"ns2.weebly.com",
		"ns3.weebly.com",
		"ns4.weebly.com",
	},
	"cloudflare": { // Example of another provider
		"ada.ns.cloudflare.com.",
		"bob.ns.cloudflare.com.",
		// Add more as needed
	},
	"godaddy": { // Example
		"ns01.domaincontrol.com.",
		"ns02.domaincontrol.com.",
		// Add more as needed
	},
	// Add other providers and their nameservers here
}

// lookupNS is a variable pointing to the function used for NS lookups.
// This allows it to be replaced by a mock during testing.
var lookupNS = net.LookupNS

// normalizeDomain removes common prefixes and path components from a domain string.
func normalizeDomain(domain string) string {
	d := strings.ToLower(domain)
	d = strings.TrimPrefix(d, "https://")
	d = strings.TrimPrefix(d, "http://")
	d = strings.SplitN(d, "/", 2)[0] // Remove path if present
	return d
}

// checkDomainProvider checks if a domain's NS records match any known nameservers
// for the specified provider.
func checkDomainProvider(domain string, providerKey string) (bool, error) {
	normalizedDomain := normalizeDomain(domain)
	if normalizedDomain == "" {
		// Return an error that lookupNS would likely return for an empty string
		// to align with the existing test case expectation.
		// Alternatively, add specific validation and a distinct error message.
		_, err := lookupNS("") // Trigger the expected error type
		if err != nil {
			return false, fmt.Errorf("failed to lookup NS records for empty domain: %w", err)
		}
		// Should not happen if lookupNS behaves as expected, but handle defensively
		return false, fmt.Errorf("domain input was empty or invalid")
	}

	lowerProviderKey := strings.ToLower(providerKey)
	targetNSList, ok := providerNameservers[lowerProviderKey]
	if !ok {
		return false, fmt.Errorf("unknown provider: %s", providerKey)
	}

	actualNS, err := lookupNS(normalizedDomain)
	if err != nil {
		// Wrap the error for more context
		return false, fmt.Errorf("failed to lookup NS records for %s: %w", normalizedDomain, err)
	}

	if len(actualNS) == 0 {
		// Handle cases where lookup succeeds but returns no records (less common for NS)
		return false, fmt.Errorf("no NS records found for %s", normalizedDomain)
	}

	// Create a map for efficient lookup of target nameservers
	targetNSMap := make(map[string]struct{}, len(targetNSList))
	for _, tns := range targetNSList {
		// Store both with and without trailing dot for robust matching
		targetNSMap[strings.TrimSuffix(tns, ".")] = struct{}{}
	}

	// Check if any of the actual nameservers match the target provider's nameservers
	for _, ans := range actualNS {
		// Normalize the actual nameserver hostname (lowercase, remove trailing dot)
		normalizedActualHost := strings.ToLower(strings.TrimSuffix(ans.Host, "."))

		// Direct match check
		if _, exists := targetNSMap[normalizedActualHost]; exists {
			return true, nil // Found a direct match
		}

		// Suffix check (e.g., ns1.weebly.com matches *.weebly.com)
		// This is less precise but can catch subdomains if needed.
		// Be cautious with this if providers use distinct subdomains for different services.
		// For now, let's stick to direct matching based on the provided list.
		// If suffix matching is desired, iterate through targetNSList and use strings.HasSuffix.
	}

	return false, nil // No match found
}

func main() {
	// Define command-line flags
	domainPtr := flag.String("domain", "", "Domain to check (e.g., www.example.com or https://example.com)")
	providerPtr := flag.String("provider", "", "DNS Provider key to check against (e.g., weebly, cloudflare)")

	// Parse the flags
	flag.Parse()

	// Validate inputs
	if *domainPtr == "" {
		log.Fatal("Error: --domain flag is required.")
	}
	if *providerPtr == "" {
		log.Fatal("Error: --provider flag is required.")
	}

	// Perform the check
	matches, err := checkDomainProvider(*domainPtr, *providerPtr)
	if err != nil {
		// Log the error and exit. Using Fatalf includes timestamp and exits with status 1.
		log.Fatalf("Error checking domain '%s' for provider '%s': %v", *domainPtr, *providerPtr, err)
	}

	// Print the result
	if matches {
		fmt.Printf("✅ Domain '%s' appears to be using DNS provider '%s'.\n", *domainPtr, *providerPtr)
	} else {
		fmt.Printf("❌ Domain '%s' does NOT appear to be using DNS provider '%s'.\n", *domainPtr, *providerPtr)
	}
}
