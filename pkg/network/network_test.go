package network

import (
	"testing"
)

func TestGetAvailableIPs(t *testing.T) {
	// Test GetAvailableIPs function
	ips := GetAvailableIPs()

	// The function may return nil or empty on machines without network
	// Just verify it doesn't panic and returns a valid type
	if ips == nil {
		// No IPs available, which is acceptable in some environments
		return
	}

	// If IPs are returned, verify they're valid strings
	for _, ip := range ips {
		if ip == "" {
			t.Errorf("GetAvailableIPs() returned empty IP string")
		}
	}
}
