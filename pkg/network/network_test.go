package network

import (
	"testing"
)

func TestGetLocalIPs(t *testing.T) {
	// Test GetLocalIPs function
	ips, err := GetLocalIPs()
	if err != nil {
		t.Errorf("GetLocalIPs() error = %v, want nil", err)
	}

	// Check that at least one IP is returned
	if len(ips) == 0 {
		t.Errorf("GetLocalIPs() returned no IPs")
	}
}
