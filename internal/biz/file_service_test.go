package biz

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zy84338719/upftp/internal/conf"
	"github.com/zy84338719/upftp/internal/dal"
)

func newTestFileService(t *testing.T) (*FileService, string) {
	t.Helper()
	root := t.TempDir()
	store := dal.NewLocalFileStore()
	resolver := dal.NewPathResolver(root)
	cfg := &conf.Config{Upload: struct {
		Enabled    bool   `yaml:"enabled"`
		MaxSize    int64  `yaml:"max_size"`
		AllowTypes string `yaml:"allow_types"`
	}{MaxSize: 1024 * 1024}}
	svc := NewFileService(store, resolver, cfg)
	return svc, root
}

func TestFileService_Upload_SizeLimit(t *testing.T) {
	svc, _ := newTestFileService(t)
	svc.cfg.Upload.MaxSize = 10

	err := svc.Upload("/test.txt", strings.NewReader("hello world this is more than 10 bytes"), 30)
	if err == nil {
		t.Error("expected size limit error, got nil")
	}
}

func TestFileService_Upload_Success(t *testing.T) {
	svc, root := newTestFileService(t)
	err := svc.Upload("/test.txt", strings.NewReader("hello"), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "test.txt"))
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("got %q, want %q", data, "hello")
	}
}

func TestFileService_Upload_NestedPath(t *testing.T) {
	svc, root := newTestFileService(t)
	err := svc.Upload("/sub/dir/test.txt", strings.NewReader("nested"), 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "sub", "dir", "test.txt"))
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(data) != "nested" {
		t.Errorf("got %q, want %q", data, "nested")
	}
}

func TestFileService_Delete_File(t *testing.T) {
	svc, root := newTestFileService(t)
	os.WriteFile(filepath.Join(root, "del.txt"), []byte("x"), 0644)

	err := svc.Delete("/del.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "del.txt")); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestFileService_Delete_Directory(t *testing.T) {
	svc, root := newTestFileService(t)
	os.MkdirAll(filepath.Join(root, "subdir"), 0755)
	os.WriteFile(filepath.Join(root, "subdir", "a.txt"), []byte("a"), 0644)

	err := svc.Delete("/subdir")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "subdir")); !os.IsNotExist(err) {
		t.Error("directory should be deleted")
	}
}

func TestFileService_Rename(t *testing.T) {
	svc, root := newTestFileService(t)
	os.WriteFile(filepath.Join(root, "old.txt"), []byte("data"), 0644)

	err := svc.Rename("/old.txt", "new.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "old.txt")); !os.IsNotExist(err) {
		t.Error("old file should be gone")
	}
	data, _ := os.ReadFile(filepath.Join(root, "new.txt"))
	if string(data) != "data" {
		t.Errorf("got %q, want %q", data, "data")
	}
}

func TestFileService_Rename_PathTraversal(t *testing.T) {
	svc, _ := newTestFileService(t)
	err := svc.Rename("/test.txt", "../../etc/passwd")
	if err == nil {
		t.Error("expected error for path traversal, got nil")
	}
}

func TestFileService_ListFiles(t *testing.T) {
	svc, root := newTestFileService(t)
	os.WriteFile(filepath.Join(root, "a.txt"), []byte("a"), 0644)
	os.MkdirAll(filepath.Join(root, "sub"), 0755)

	files, err := svc.ListFiles("/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	names := make(map[string]bool)
	for _, f := range files {
		names[f.Name] = true
		if f.Name == "sub" && !f.IsDir {
			t.Error("sub should be a directory")
		}
		if f.Name == "a.txt" && f.IsDir {
			t.Error("a.txt should not be a directory")
		}
	}
	if !names["a.txt"] || !names["sub"] {
		t.Error("missing expected files")
	}
}

func TestFileService_BuildTree(t *testing.T) {
	svc, root := newTestFileService(t)
	os.MkdirAll(filepath.Join(root, "sub1", "nested"), 0755)
	os.MkdirAll(filepath.Join(root, "sub2"), 0755)
	os.WriteFile(filepath.Join(root, "sub1", "file.txt"), []byte("x"), 0644)

	tree, err := svc.BuildTree("/", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tree.Name != "/" {
		t.Errorf("root name = %q, want %q", tree.Name, "/")
	}
	if len(tree.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(tree.Children))
	}
}

func TestFileService_ReadFileContent(t *testing.T) {
	svc, root := newTestFileService(t)
	os.WriteFile(filepath.Join(root, "read.txt"), []byte("content"), 0644)

	data, err := svc.ReadFileContent("/read.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "content" {
		t.Errorf("got %q, want %q", data, "content")
	}
}

func TestFileService_ReadFileContent_TooLarge(t *testing.T) {
	svc, root := newTestFileService(t)
	largeData := make([]byte, 11*1024*1024)
	os.WriteFile(filepath.Join(root, "large.bin"), largeData, 0644)

	_, err := svc.ReadFileContent("/large.bin")
	if err == nil {
		t.Error("expected size limit error, got nil")
	}
}

func TestFileService_CreateFolder(t *testing.T) {
	svc, root := newTestFileService(t)
	err := svc.CreateFolder("/new/folder")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "new", "folder")); err != nil {
		t.Fatalf("folder not created: %v", err)
	}
}

func TestFileService_Move(t *testing.T) {
	svc, root := newTestFileService(t)
	os.WriteFile(filepath.Join(root, "src.txt"), []byte("move"), 0644)

	err := svc.Move("/src.txt", "/dst.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(root, "dst.txt"))
	if string(data) != "move" {
		t.Errorf("got %q, want %q", data, "move")
	}
}

func TestFileService_Copy(t *testing.T) {
	svc, root := newTestFileService(t)
	os.WriteFile(filepath.Join(root, "src.txt"), []byte("copy"), 0644)

	err := svc.Copy("/src.txt", "/dst.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(filepath.Join(root, "dst.txt"))
	if string(data) != "copy" {
		t.Errorf("got %q, want %q", data, "copy")
	}
}

func TestFileService_PathTraversal(t *testing.T) {
	svc, _ := newTestFileService(t)

	paths := []string{
		"../../etc/passwd",
		"../secret",
		"sub/../../../etc",
	}

	for _, p := range paths {
		_, err := svc.Download(p)
		if err == nil {
			t.Errorf("expected error for path %q, got nil", p)
		}
	}
}
