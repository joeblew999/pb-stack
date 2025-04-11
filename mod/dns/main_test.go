// main_test.go
package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

// --- Mocking DNS Lookup ---

// mockLookupNS is our mock implementation of net.LookupNS
func mockLookupNS(domain string) ([]*net.NS, error) {
	// Predefined responses for specific domains used in tests
	switch domain {
	case "www.dhaba517.com": // Example domain using weebly
		return []*net.NS{
			{Host: "ns1.weebly.com."},
			{Host: "ns2.weebly.com."},
		}, nil
	case "www.google.com": // Example domain *not* using weebly
		return []*net.NS{
			{Host: "ns1.google.com."},
			{Host: "ns2.google.com."},
		}, nil
	case "nonexistentdomain.xyz": // Example domain that causes a lookup error
		return nil, &net.DNSError{Err: "no such host", Name: domain, IsNotFound: true}
	default:
		return nil, fmt.Errorf("unexpected domain in mockLookupNS: %s", domain)
	}
}

// --- Test Cases for checkDomainProvider ---

func TestCheckDomainProvider(t *testing.T) {
	// --- Setup Mock ---
	// Store original lookup function
	originalLookupNS := lookupNS
	// Replace with mock
	lookupNS = mockLookupNS
	// Ensure original is restored after test
	t.Cleanup(func() {
		lookupNS = originalLookupNS
	})
	// --- End Setup Mock ---

	testCases := []struct {
		name          string
		domain        string
		providerKey   string
		wantMatch     bool
		wantErr       bool
		wantErrSubstr string // Optional: check for specific error text
	}{
		{
			name:        "Matching Domain (Weebly)",
			domain:      "www.dhaba517.com",
			providerKey: "weebly",
			wantMatch:   true,
			wantErr:     false,
		},
		{
			name:        "Matching Domain with https prefix (Weebly)",
			domain:      "https://www.dhaba517.com",
			providerKey: "weebly",
			wantMatch:   true,
			wantErr:     false,
		},
		{
			name:        "Matching Domain case insensitive provider",
			domain:      "www.dhaba517.com",
			providerKey: "Weebly", // Mixed case
			wantMatch:   true,
			wantErr:     false,
		},
		{
			name:        "Non-Matching Domain (Google vs Weebly)",
			domain:      "www.google.com",
			providerKey: "weebly",
			wantMatch:   false,
			wantErr:     false,
		},
		{
			name:          "Unknown Provider",
			domain:        "www.dhaba517.com",
			providerKey:   "unknownprovider",
			wantMatch:     false,
			wantErr:       true,
			wantErrSubstr: "unknown provider: unknownprovider",
		},
		{
			name:          "DNS Lookup Error (NXDOMAIN)",
			domain:        "nonexistentdomain.xyz",
			providerKey:   "weebly",
			wantMatch:     false,
			wantErr:       true,
			wantErrSubstr: "failed to lookup NS records", // Check general failure
		},
		{
			name:          "DNS Lookup Error (NXDOMAIN) - Check specific error",
			domain:        "nonexistentdomain.xyz",
			providerKey:   "weebly",
			wantMatch:     false,
			wantErr:       true,
			wantErrSubstr: "no such host", // Check underlying DNS error text
		},
		{
			name:          "Empty Domain Input", // Assuming checkDomainProvider handles this gracefully or main prevents it
			domain:        "",
			providerKey:   "weebly",
			wantMatch:     false,                         // Expect no match or error depending on implementation
			wantErr:       true,                          // Likely an error from lookupNS or prior validation
			wantErrSubstr: "failed to lookup NS records", // Or specific validation error
		},
		// Add more test cases for other providers if defined in providerNameservers
		// Add tests for edge cases in domain normalization if needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotMatch, err := checkDomainProvider(tc.domain, tc.providerKey)

			// Check error expectation
			if tc.wantErr {
				if err == nil {
					t.Errorf("checkDomainProvider() error = nil, wantErr %v", tc.wantErr)
					return // Avoid further checks if error was expected but missing
				}
				if tc.wantErrSubstr != "" && !strings.Contains(err.Error(), tc.wantErrSubstr) {
					t.Errorf("checkDomainProvider() error = %q, want error containing %q", err, tc.wantErrSubstr)
				}
			} else { // No error expected
				if err != nil {
					t.Errorf("checkDomainProvider() unexpected error = %v", err)
				}
			}

			// Check match expectation (only if no error was expected or if error handling allows returning a match value)
			// Adjust this logic based on how checkDomainProvider behaves on error
			if !tc.wantErr && gotMatch != tc.wantMatch {
				t.Errorf("checkDomainProvider() gotMatch = %v, want %v", gotMatch, tc.wantMatch)
			}
			// If an error is expected, you might not care about the match value,
			// or you might expect it to be false.
			if tc.wantErr && gotMatch != false {
				t.Logf("checkDomainProvider() returned match=%v on expected error, ensuring it's false", gotMatch)
				if gotMatch { // Explicitly fail if match is true when error is expected
					t.Errorf("checkDomainProvider() gotMatch = true, want false when error is expected")
				}
			}
		})
	}
}

// Note: Testing the `main` function itself (argument parsing, output) often involves
// more complex setup using os/exec or dedicated CLI testing libraries.
// The tests above focus on the core logic function `checkDomainProvider`.
