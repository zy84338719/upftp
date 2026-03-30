package cli

import (
	"fmt"
	"strings"
)

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func truncatePath(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return "..." + s[len(s)-maxLen+3:]
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func maskPassword(password string) string {
	if len(password) == 0 {
		return ""
	}
	if len(password) <= 2 {
		return "****"
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
