package dal

import (
	"fmt"
	"path/filepath"
	"strings"
)

type localPathResolver struct {
	rootDir string
}

func NewPathResolver(rootDir string) PathResolver {
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		panic("invalid root dir: " + err.Error())
	}
	if !strings.HasSuffix(absRoot, string(filepath.Separator)) {
		absRoot += string(filepath.Separator)
	}
	return &localPathResolver{rootDir: absRoot}
}

func (r *localPathResolver) Root() string {
	return r.rootDir
}

func (r *localPathResolver) IsPathSafe(relativePath string) bool {
	if relativePath == "/" || relativePath == "" {
		return true
	}
	cleaned := filepath.Clean(filepath.FromSlash(relativePath))
	if cleaned == "." {
		return true
	}
	if cleaned == ".." || cleaned == string(filepath.Separator)+".." {
		return false
	}
	if strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return false
	}
	absPath := filepath.Join(r.rootDir, cleaned)
	return strings.HasPrefix(absPath, r.rootDir)
}

func (r *localPathResolver) SafeJoin(relativePath string) string {
	cleaned := filepath.Clean(filepath.FromSlash(relativePath))
	cleaned = strings.TrimPrefix(cleaned, string(filepath.Separator))
	return filepath.Join(r.rootDir, cleaned)
}

func (r *localPathResolver) ResolvePath(relativePath string) (string, error) {
	if !r.IsPathSafe(relativePath) {
		return "", fmt.Errorf("path traversal denied: %s", relativePath)
	}
	return r.SafeJoin(relativePath), nil
}
