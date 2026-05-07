package index

import (
	"embed"
	"io/fs"
	"path"
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
// 使用 path.Join 而非 filepath.Join，因为 embed.FS 始终使用正斜杠作为路径分隔符
func ReadTemplateFile(name string) ([]byte, error) {
	return embeddedFS.ReadFile(path.Join("templates", name))
}
