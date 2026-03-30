package biz

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zy84338719/upftp/internal/conf"
	"github.com/zy84338719/upftp/internal/dal"
)

type FileService struct {
	store dal.FileStore
	path  dal.PathResolver
	cfg   *conf.Config
}

func NewFileService(store dal.FileStore, path dal.PathResolver, cfg *conf.Config) *FileService {
	return &FileService{store: store, path: path, cfg: cfg}
}

func (s *FileService) PathResolver() dal.PathResolver {
	return s.path
}

func (s *FileService) Upload(path string, reader io.Reader, size int64) error {
	if s.cfg.Upload.MaxSize > 0 && size > s.cfg.Upload.MaxSize {
		return fmt.Errorf("file size %d exceeds limit %d", size, s.cfg.Upload.MaxSize)
	}
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return err
	}
	f, err := s.store.CreateFile(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	written, err := io.Copy(f, reader)
	if err != nil {
		return fmt.Errorf("write file: %w (written %d bytes)", err, written)
	}
	return nil
}

func (s *FileService) Download(path string) (string, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return "", err
	}
	_, err = s.store.StatFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("file not found: %w", err)
	}
	return fullPath, nil
}

func (s *FileService) Delete(path string) error {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return err
	}
	info, err := s.store.StatFile(fullPath)
	if err != nil {
		return fmt.Errorf("not found: %w", err)
	}
	if info.IsDir() {
		return s.store.RemoveDir(fullPath)
	}
	return s.store.DeleteFile(fullPath)
}

func (s *FileService) Rename(path string, newName string) error {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return err
	}
	newPath := filepath.Join(filepath.Dir(fullPath), newName)
	if !strings.HasPrefix(newPath, s.path.Root()) {
		return fmt.Errorf("new name results in path outside root")
	}
	return s.store.RenameFile(fullPath, newPath)
}

func (s *FileService) Move(src string, dst string) error {
	fullSrc, err := s.path.ResolvePath(src)
	if err != nil {
		return err
	}
	fullDst, err := s.path.ResolvePath(dst)
	if err != nil {
		return err
	}
	return s.store.RenameFile(fullSrc, fullDst)
}

func (s *FileService) Copy(src string, dst string) error {
	fullSrc, err := s.path.ResolvePath(src)
	if err != nil {
		return err
	}
	fullDst, err := s.path.ResolvePath(dst)
	if err != nil {
		return err
	}
	_, err = s.store.CopyFile(fullSrc, fullDst)
	return err
}

func (s *FileService) ReadFileContent(path string) ([]byte, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	info, err := s.store.StatFile(fullPath)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("cannot read directory")
	}
	if info.Size() > 10*1024*1024 {
		return nil, fmt.Errorf("file too large (max 10MB)")
	}
	return s.store.ReadFile(fullPath)
}

func (s *FileService) WriteFileContent(path string, content []byte) error {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return err
	}
	return s.store.WriteFile(fullPath, content, 0644)
}

func (s *FileService) Stat(path string) (os.FileInfo, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	return s.store.StatFile(fullPath)
}

// SafePath resolves and validates a path, returning the absolute path.
// Used by protocol layers that need the resolved path for streaming.
func (s *FileService) SafePath(relativePath string) (string, error) {
	return s.path.ResolvePath(relativePath)
}

// OpenFile opens a file for reading after path validation.
func (s *FileService) OpenFile(path string) (*os.File, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	return s.store.OpenFile(fullPath)
}

// CreateFileForWrite creates a file for writing after path validation.
func (s *FileService) CreateFileForWrite(path string) (*os.File, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	return s.store.CreateFile(fullPath)
}

// ListDir returns directory entries after path validation.
func (s *FileService) ListDir(path string) ([]os.DirEntry, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	return s.store.ListDir(fullPath)
}
