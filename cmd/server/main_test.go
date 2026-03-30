package main

import (
	"testing"
)

func TestVersionInfo(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
	if LastCommit == "" {
		t.Error("LastCommit should not be empty")
	}
}
