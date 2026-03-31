package index

import (
	"embed"
	"io/fs"
	"path/filepath"
)

//go:embed templates/*
var embeddedFS embed.FS

// GetTemplatesFS 返回嵌入的模板文件系统
func GetTemplatesFS() fs.FS {
	fs, err := fs.Sub(embeddedFS, "templates")
	if err != nil {
		panic(err)
	}
	return fs
}

// ReadTemplateFile 读取嵌入的模板文件
func ReadTemplateFile(path string) ([]byte, error) {
	return embeddedFS.ReadFile(filepath.Join("templates", path))
}
