package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zy84338719/upftp/pkg/file/model"
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
		dirIcon := "📁"
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

// ListFilesLegacy returns legacy file info for compatibility.
func (s *FileService) ListFilesLegacy(path string) ([]model.LegacyFileInfo, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	entries, err := s.store.ListDir(fullPath)
	if err != nil {
		return nil, err
	}

	files := make([]model.LegacyFileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		ft := model.GetFileType(entry.Name())
		dirIcon := "📁"
		icon := model.GetFileIcon(ft)
		if entry.IsDir() {
			icon = dirIcon
		}
		files = append(files, model.LegacyFileInfo{
			Name:       entry.Name(),
			Size:       formatSize(info.Size()),
			ModTime:    info.ModTime().Format("2006-01-02 15:04"),
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

// GetFile returns a file reader, size and mime type.
func (s *FileService) GetFile(path string) (io.ReadCloser, int64, string, error) {
	fullPath, err := s.path.ResolvePath(path)
	if err != nil {
		return nil, 0, "", err
	}
	info, err := s.store.StatFile(fullPath)
	if err != nil {
		return nil, 0, "", err
	}
	if info.IsDir() {
		return nil, 0, "", fmt.Errorf("path is a directory")
	}
	file, err := s.store.OpenFile(fullPath)
	if err != nil {
		return nil, 0, "", err
	}
	ft := model.GetFileType(path)
	return file, info.Size(), model.GetMimeType(ft), nil
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
		if entry.IsDir() {
			// 递归处理子目录
			childPath := filepath.Join(path, entry.Name())
			child, err := s.BuildTree(childPath, depth+1)
			if err == nil && child != nil {
				node.Children = append(node.Children, child)
			}
		} else {
			// 添加文件节点
			fileNode := &model.TreeNode{
				Name:  entry.Name(),
				Path:  filepath.Join(path, entry.Name()),
				IsDir: false,
			}
			node.Children = append(node.Children, fileNode)
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

// formatSize formats file size to human readable string.
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
