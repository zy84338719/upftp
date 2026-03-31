package index

import (
	"html/template"
	"os"
	"path/filepath"
)

func getTemplatesFS() string {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return "./biz/handler/index/templates"
	}
	return filepath.Join(wd, "biz/handler/index/templates")
}

func parseTemplate(name string) (*template.Template, error) {
	templateDir := getTemplatesFS()
	templatePath := filepath.Join(templateDir, name)
	return template.ParseFiles(templatePath)
}
