package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileEntry 是文件列表中的一个条目,用于 HTTP API 响应与 Web 页面。
type FileEntry struct {
	Name    string `json:"name"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size"`
	ModTime string `json:"mod_time"`
}

// Files 封装对共享目录的文件操作。Root 在构造时固定,所有方法都做沙箱校验。
type Files struct {
	Root     string
	ReadOnly bool
}

// SecureJoin 把相对路径安全映射到 Root 之下,防止目录穿越。
func (f *Files) SecureJoin(rel string) (string, error) {
	root := filepath.Clean(f.Root)
	rel = strings.TrimPrefix(rel, "/")
	joined := filepath.Join(root, rel)
	cleaned := filepath.Clean(joined)
	if cleaned != root && !strings.HasPrefix(cleaned, root+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes root: %s", cleaned)
	}
	return cleaned, nil
}

// List 列出 rel 目录下的条目,按"目录在前、名字升序"排列。
func (f *Files) List(rel string) ([]FileEntry, error) {
	abs, err := f.SecureJoin(rel)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return nil, err
	}
	out := make([]FileEntry, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, FileEntry{
			Name:    e.Name(),
			IsDir:   e.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04"),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsDir != out[j].IsDir {
			return out[i].IsDir // 目录在前
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

// Open 打开 rel 路径的文件用于读取。
func (f *Files) Open(rel string) (*os.File, os.FileInfo, error) {
	abs, err := f.SecureJoin(rel)
	if err != nil {
		return nil, nil, err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, nil, err
	}
	if info.IsDir() {
		return nil, nil, fmt.Errorf("is a directory")
	}
	file, err := os.Open(abs)
	if err != nil {
		return nil, nil, err
	}
	return file, info, nil
}

// Save 从 r 读取内容写入 rel 路径。ReadOnly 时拒绝。
func (f *Files) Save(rel string, r io.Reader) error {
	if f.ReadOnly {
		return fmt.Errorf("read-only mode")
	}
	abs, err := f.SecureJoin(rel)
	if err != nil {
		return err
	}
	if mkErr := os.MkdirAll(filepath.Dir(abs), 0o755); mkErr != nil {
		return mkErr
	}
	tmp := abs + ".upftp-tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, r); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, abs)
}

// Remove 删除 rel(文件或目录)。ReadOnly 时拒绝。
func (f *Files) Remove(rel string) error {
	if f.ReadOnly {
		return fmt.Errorf("read-only mode")
	}
	abs, err := f.SecureJoin(rel)
	if err != nil {
		return err
	}
	if abs == filepath.Clean(f.Root) {
		return fmt.Errorf("cannot remove root")
	}
	return os.RemoveAll(abs)
}

// Mkdir 创建 rel 目录。
func (f *Files) Mkdir(rel string) error {
	if f.ReadOnly {
		return fmt.Errorf("read-only mode")
	}
	abs, err := f.SecureJoin(rel)
	if err != nil {
		return err
	}
	return os.MkdirAll(abs, 0o755)
}
