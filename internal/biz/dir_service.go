package biz

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zy84338719/upftp/internal/model"
)

const MaxTreeDepth = 10

// ListFiles lists files and directories at the given path.
func (s *FileService) ListFiles(path string) ([]model.FileInfo, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	entries, err := s.store.ListDir(fullPath)
	if err != nil {
		return nil, err
	}

	files := make([]model.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		ft := model.GetFileType(entry.Name())
		dirIcon := "\U0001F4C1"
		icon := model.GetFileIcon(ft)
		if entry.IsDir() {
			icon = dirIcon
		}
		files = append(files, model.FileInfo{
			Name:       entry.Name(),
			Size:       info.Size(),
			ModTime:    info.ModTime(),
			IsDir:      entry.IsDir(),
			Path:       filepath.Join(path, entry.Name()),
			FileType:   ft,
			CanPreview: !entry.IsDir() && model.CanPreviewFile(ft),
			MimeType:   model.GetMimeType(ft),
			Icon:       icon,
		})
	}
	return files, nil
}

// CreateFolder creates a directory at the given path.
func (s *FileService) CreateFolder(path string) error {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return err
	}
	return s.store.CreateDir(fullPath, 0755)
}

// BuildTree builds a directory tree structure with depth limiting.
func (s *FileService) BuildTree(path string, depth int) (*model.TreeNode, error) {
	if depth > MaxTreeDepth {
		return nil, nil
	}
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	entries, err := s.store.ListDir(fullPath)
	if err != nil {
		return nil, err
	}

	node := &model.TreeNode{
		Name:  filepath.Base(path),
		Path:  "/" + strings.TrimPrefix(path, "/"),
		IsDir: true,
	}
	if path == "" || path == "/" {
		node.Name = "/"
		node.Path = "/"
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		childPath := filepath.Join(path, entry.Name())
		child, err := s.BuildTree(childPath, depth+1)
		if err == nil && child != nil {
			node.Children = append(node.Children, child)
		}
	}
	return node, nil
}

// SearchFiles searches for files matching a pattern.
func (s *FileService) SearchFiles(pattern string, path string) ([]string, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	var results []string
	err = s.store.WalkDir(fullPath, func(walkPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		matched, matchErr := filepath.Match(pattern, info.Name())
		if matchErr != nil || !matched {
			return nil
		}
		relPath, _ := filepath.Rel(s.path.Root(), walkPath)
		results = append(results, "/"+relPath)
		return nil
	})
	return results, err
}
