package handlers

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zy84338719/upftp/v2/internal/config"
	"github.com/zy84338719/upftp/v2/internal/filehandler"
)

// handleModernIndex 处理新的现代化页面
func handleModernIndex(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	if !strings.HasPrefix(urlPath, "/modern/") {
		urlPath = "/"
	} else {
		urlPath = strings.TrimPrefix(urlPath, "/modern")
	}

	if urlPath == "" {
		urlPath = "/"
	}

	fsPath := path.Join(config.AppConfig.Root, urlPath)
	fileInfo, err := os.Stat(fsPath)
	if err != nil {
		if urlPath != "/" {
			http.NotFound(w, r)
			return
		}
		handleModernFileList(w, r, urlPath)
		return
	}

	if !fileInfo.IsDir() {
		http.ServeFile(w, r, fsPath)
		return
	}

	handleModernFileList(w, r, urlPath)
}

// handleModernFileList 处理文件列表显示
func handleModernFileList(w http.ResponseWriter, r *http.Request, urlPath string) {
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
		filePath := path.Join(urlPath, file.Name())
		fileType := filehandler.GetFileType(file.Name())

		info := filehandler.FileInfo{
			Name:       file.Name(),
			Size:       filehandler.FormatFileSize(file.Size()),
			ModTime:    file.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:      file.IsDir(),
			CanPreview: !file.IsDir() && filehandler.CanPreviewFile(fileType),
			FileType:   fileType,
			Path:       filePath,
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
	}

	tmpl, _ := template.ParseFS(templates, "templates/modern-index.html")
	tmpl.Execute(w, data)
}
