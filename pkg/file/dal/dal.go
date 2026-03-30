package dal

import (
	"os"
	"path/filepath"
)

type FileStore interface {
	ReadFile(fullPath string) ([]byte, error)
	WriteFile(fullPath string, data []byte, perm os.FileMode) error
	DeleteFile(fullPath string) error
	RenameFile(oldPath string, newPath string) error
	StatFile(fullPath string) (os.FileInfo, error)
	OpenFile(fullPath string) (*os.File, error)
	CreateFile(fullPath string) (*os.File, error)
	ListDir(fullPath string) ([]os.DirEntry, error)
	CreateDir(fullPath string, perm os.FileMode) error
	RemoveDir(fullPath string) error
	WalkDir(root string, fn filepath.WalkFunc) error
	CopyFile(src string, dst string) (int64, error)
}

type PathResolver interface {
	IsPathSafe(relativePath string) bool
	SafeJoin(relativePath string) string
	ResolvePath(relativePath string) (string, error)
	Root() string
}
