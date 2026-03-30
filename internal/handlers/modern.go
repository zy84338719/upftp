package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"path"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/filehandler"
)

func handleModernIndex(ctx context.Context, c *app.RequestContext) {
	urlPath := string(c.Request.URI().Path())
	if urlPath == "" {
		urlPath = "/"
	}

	fsPath := path.Join(config.AppConfig.Root, urlPath)
	fileInfo, err := os.Stat(fsPath)
	if err != nil {
		if urlPath != "/" {
			c.String(404, "Not Found")
			return
		}
		handleModernFileList(ctx, c, urlPath)
		return
	}

	if !fileInfo.IsDir() {
		c.File(fsPath)
		return
	}

	handleModernFileList(ctx, c, urlPath)
}

func handleModernFileList(ctx context.Context, c *app.RequestContext, urlPath string) {
	fsPath := path.Join(config.AppConfig.Root, urlPath)
	files, _ := ioutil.ReadDir(fsPath)
	fileList := []filehandler.FileInfo{}

	if urlPath != "/" {
		fileList = append(fileList, filehandler.FileInfo{
			Name:  "..",
			IsDir: true,
			Path:  path.Dir(urlPath),
			Icon:  "📁",
		})
	}

	for _, file := range files {
		fp := path.Join(urlPath, file.Name())
		fileType := filehandler.GetFileType(file.Name())

		info := filehandler.FileInfo{
			Name:       file.Name(),
			Size:       filehandler.FormatFileSize(file.Size()),
			ModTime:    file.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:      file.IsDir(),
			CanPreview: !file.IsDir() && filehandler.CanPreviewFile(fileType),
			FileType:   fileType,
			Path:       fp,
			Icon:       getFileIcon(file.IsDir(), fileType),
		}
		fileList = append(fileList, info)
	}

	data := struct {
		Files         []filehandler.FileInfo
		ServerInfo    *ServerInfo
		CurrentPath   string
		UploadEnabled bool
		Version       string
		BuildDate     string
		GoVersion     string
		Platform      string
		ProjectURL    string
		ProjectName   string
		LastCommit    string
		Language      string
		HTTPAuthOn    bool
		HTTPAuthUser  string
		HTTPAuthPass  string
	}{
		Files:         fileList,
		ServerInfo:    serverInfo,
		CurrentPath:   urlPath,
		UploadEnabled: config.AppConfig.Upload.Enabled,
		Version:       config.AppConfig.Version,
		BuildDate:     config.AppConfig.BuildDate,
		GoVersion:     config.AppConfig.GoVersion,
		Platform:      config.AppConfig.Platform,
		ProjectURL:    config.AppConfig.ProjectURL,
		ProjectName:   config.AppConfig.ProjectName,
		LastCommit:    config.AppConfig.LastCommit,
		Language:      config.AppConfig.GetLanguage(),
		HTTPAuthOn:    config.AppConfig.HTTPAuth.Enabled,
		HTTPAuthUser:  config.AppConfig.HTTPAuth.Username,
		HTTPAuthPass:  config.AppConfig.HTTPAuth.Password,
	}

	funcMap := template.FuncMap{
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}

	content, err := templates.ReadFile("templates/index.html")
	if err != nil {
		c.String(500, "Template read error: "+err.Error())
		return
	}
	tmpl, err := template.New("index").Funcs(funcMap).Parse(string(content))
	if err != nil {
		c.String(500, "Template parse error: "+err.Error())
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(c.Response.BodyWriter(), data); err != nil {
		c.String(500, "Template execute error: "+err.Error())
	}
}
