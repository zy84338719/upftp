package dal

import (
	"io"
	"os"
	"path/filepath"
)

type localFileStore struct{}

func NewLocalFileStore() FileStore {
	return &localFileStore{}
}

func (s *localFileStore) ReadFile(fullPath string) ([]byte, error) {
	return os.ReadFile(fullPath)
}

func (s *localFileStore) WriteFile(fullPath string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, data, perm)
}

func (s *localFileStore) DeleteFile(fullPath string) error {
	return os.Remove(fullPath)
}

func (s *localFileStore) RenameFile(oldPath string, newPath string) error {
	dir := filepath.Dir(newPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.Rename(oldPath, newPath)
}

func (s *localFileStore) StatFile(fullPath string) (os.FileInfo, error) {
	return os.Stat(fullPath)
}

func (s *localFileStore) OpenFile(fullPath string) (*os.File, error) {
	return os.Open(fullPath)
}

func (s *localFileStore) CreateFile(fullPath string) (*os.File, error) {
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return os.Create(fullPath)
}

func (s *localFileStore) ListDir(fullPath string) ([]os.DirEntry, error) {
	return os.ReadDir(fullPath)
}

func (s *localFileStore) CreateDir(fullPath string, perm os.FileMode) error {
	return os.MkdirAll(fullPath, perm)
}

func (s *localFileStore) RemoveDir(fullPath string) error {
	return os.RemoveAll(fullPath)
}

func (s *localFileStore) WalkDir(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}

func (s *localFileStore) CopyFile(src string, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}
