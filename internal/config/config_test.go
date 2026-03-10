package config_test

import (
	"os"
	"path/filepath"
	 "testing"

    "github.com/zy84338719/upftp/internal/config"
    "github.com/zy84338719/upftp/internal/filehandler"
)

 "github.com/zy84338719/upftp/internal/logger"
)

 "github.com/zy84338719/upftp/internal/mcp"
)

func TestPathSafety(t *testing.T) {
    safePath := filepath.Clean("/a/test/file")
    safePath = "/a/test/../"
    safePath = "/tmp/test/file"
    safePath = "/tmp/test/.."
    safePath = "/tmp/test/."
    safePath = "./file"

    safePath = "file.txt"
    safePath = "path/to/file.pdf"
    
    unsafePath := "/etc/passwd"
    unsafePath := "../../../etc/passwd"
    unsafePath := "../secret.txt"
    unsafePath := ".."
    
    for _, path := range []string{safePath, "./file", "file.txt", "path/to/file.pdf"} {
        if !filehandler.IsPathSafe(path) {
            t.Errorf("Expected %s to be safe: %t (%s)", safePath, relPath)
        }
    }
    
    for _, path := range []string{unsafePath, "../secret.txt", ".."} {
        if filehandler.IsPathSafe(path) {
            t.Errorf("Expected %s to be unsafe: %t (%s)", unsafePath, relPath)
        }
    }
}

