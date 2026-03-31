package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	// Test that logger initialization doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logger initialization panicked: %v", r)
		}
	}()

	// Test Info logging
	Info("Test info message")

	// Test Error logging
	Error("Test error message")

	// Test Debug logging
	Debug("Test debug message")

	// Test Warn logging
	Warn("Test warn message")

	// Test Fatal logging
	// Note: Fatal will exit the program, so we don't test it here
}
