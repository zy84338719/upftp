package dal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathResolver_IsPathSafe(t *testing.T) {
	root := t.TempDir()
	resolver := NewPathResolver(root)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"root path", "/", true},
		{"empty path", "", true},
		{"normal file", "test.txt", true},
		{"nested file", "sub/dir/test.txt", true},
		{"parent traversal", "..", false},
		{"parent traversal prefix", "../etc/passwd", false},
		{"nested parent traversal", "sub/../../etc", false},
		{"dot", ".", true},
		{"dot in path", "./test.txt", true},
		{"leading slash normal", "/sub/file.txt", true},
		{"double dot in name", "file..txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolver.IsPathSafe(tt.path); got != tt.expected {
				t.Errorf("IsPathSafe(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestPathResolver_ResolvePath(t *testing.T) {
	root := t.TempDir()
	resolver := NewPathResolver(root)

	resolved, err := resolver.ResolvePath("test.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Join(root, "test.txt")
	if resolved != expected {
		t.Errorf("ResolvePath(%q) = %q, want %q", "test.txt", resolved, expected)
	}

	_, err = resolver.ResolvePath("../../etc/passwd")
	if err == nil {
		t.Error("expected error for path traversal, got nil")
	}

	_, err = resolver.ResolvePath("..")
	if err == nil {
		t.Error("expected error for parent traversal, got nil")
	}
}

func TestPathResolver_SafeJoin(t *testing.T) {
	root := t.TempDir()
	resolver := NewPathResolver(root)

	result := resolver.SafeJoin("/sub/file.txt")
	expected := filepath.Join(root, "sub/file.txt")
	if result != expected {
		t.Errorf("SafeJoin(%q) = %q, want %q", "/sub/file.txt", result, expected)
	}
}

func TestPathResolver_TraversalVectors(t *testing.T) {
	root := t.TempDir()
	resolver := NewPathResolver(root)

	vectors := []string{
		"..",
		"../",
		"../../etc/passwd",
		"sub/../../../etc",
	}

	for _, v := range vectors {
		_, err := resolver.ResolvePath(v)
		if err == nil {
			t.Errorf("expected error for traversal vector %q, got nil", v)
		}
	}
}

func TestLocalFileStore_ReadWrite(t *testing.T) {
	store := NewLocalFileStore()
	dir := t.TempDir()

	data := []byte("hello world")
	path := filepath.Join(dir, "test.txt")

	if err := store.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	got, err := store.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("ReadFile = %q, want %q", got, data)
	}
}

func TestLocalFileStore_DeleteFile(t *testing.T) {
	store := NewLocalFileStore()
	dir := t.TempDir()

	path := filepath.Join(dir, "del.txt")
	store.WriteFile(path, []byte("x"), 0644)

	if err := store.DeleteFile(path); err != nil {
		t.Fatalf("DeleteFile error: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestLocalFileStore_ListDir(t *testing.T) {
	store := NewLocalFileStore()
	dir := t.TempDir()

	store.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	store.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)

	entries, err := store.ListDir(dir)
	if err != nil {
		t.Fatalf("ListDir error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("ListDir returned %d entries, want 2", len(entries))
	}
}

func TestLocalFileStore_CopyFile(t *testing.T) {
	store := NewLocalFileStore()
	dir := t.TempDir()

	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	store.WriteFile(src, []byte("copy me"), 0644)

	n, err := store.CopyFile(src, dst)
	if err != nil {
		t.Fatalf("CopyFile error: %v", err)
	}
	if n != 7 {
		t.Errorf("CopyFile copied %d bytes, want 7", n)
	}
}
