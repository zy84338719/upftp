package handlers

import (
	"archive/zip"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/filehandler"
)

//go:embed templates/*
var templates embed.FS

type ServerInfo struct {
	IP       string
	HTTPPort int
	FTPPort  int
	Root     string
}

var serverInfo *ServerInfo

func SetServerInfo(ip string, httpPort, ftpPort int, root string) {
	serverInfo = &ServerInfo{
		IP:       ip,
		HTTPPort: httpPort,
		FTPPort:  ftpPort,
		Root:     root,
	}
}

func GetServerInfo() *ServerInfo {
	return serverInfo
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/info", HandleServerInfo)
	mux.HandleFunc("/api/tree", HandleDirectoryTree)
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(config.AppConfig.Root))))
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/download/", handleDownload)
	mux.HandleFunc("/preview/", handlePreview)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}

	fsPath := path.Join(config.AppConfig.Root, urlPath)
	fileInfo, err := os.Stat(fsPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if !fileInfo.IsDir() {
		http.ServeFile(w, r, fsPath)
		return
	}

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
			Name:        file.Name(),
			Size:        filehandler.FormatFileSize(file.Size()),
			ModTime:     file.ModTime().Format("2006-01-02 15:04:05"),
			IsDir:       file.IsDir(),
			CanPreview:  !file.IsDir() && filehandler.CanPreviewFile(fileType),
			FileType:    fileType,
			FileTypeStr: filehandler.GetFileTypeString(fileType),
			Path:        filePath,
			Icon:        getFileIcon(file.IsDir(), fileType),
			MimeType:    filehandler.GetMimeType(fileType),
		}
		fileList = append(fileList, info)
	}

	data := struct {
		Files       []filehandler.FileInfo
		ServerInfo  *ServerInfo
		CurrentPath string
	}{
		Files:       fileList,
		ServerInfo:  serverInfo,
		CurrentPath: urlPath,
	}

	tmpl, _ := template.ParseFS(templates, "templates/index.html")
	tmpl.Execute(w, data)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if !filehandler.IsPathSafe(filename) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	filePath := path.Join(config.AppConfig.Root, filename)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if !fileInfo.IsDir() {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filename)))
		http.ServeFile(w, r, filePath)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", filepath.Base(filename)))

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(filePath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	if err != nil {
		http.Error(w, "Error creating zip file", http.StatusInternalServerError)
		return
	}
}

func handlePreview(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/preview/")
	if !filehandler.IsPathSafe(filename) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	filePath := path.Join(config.AppConfig.Root, filename)

	fileType := filehandler.GetFileType(filename)
	mimeType := filehandler.GetMimeType(fileType)
	w.Header().Set("Content-Type", mimeType)

	http.ServeFile(w, r, filePath)
}

func HandleServerInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := fmt.Sprintf(`{
		"version": "%s",
		"lastCommit": "%s",
		"httpPort": %d,
		"ftpPort": %d,
		"ftpEnabled": %t,
		"rootPath": "%s"
	}`,
		config.AppConfig.Version,
		config.AppConfig.LastCommit,
		serverInfo.HTTPPort,
		serverInfo.FTPPort,
		config.AppConfig.EnableFTP,
		config.AppConfig.Root,
	)

	w.Write([]byte(response))
}

func getFileIcon(isDir bool, fileType filehandler.FileType) string {
	if isDir {
		return "📁"
	}
	return filehandler.GetFileIcon(fileType)
}

type TreeNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"isDir"`
	Children []*TreeNode `json:"children,omitempty"`
	Expanded bool        `json:"expanded,omitempty"`
}

func HandleDirectoryTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tree, err := buildDirectoryTree(config.AppConfig.Root, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytes)
}

func buildDirectoryTree(rootPath string, relativePath string) (*TreeNode, error) {
	fullPath := path.Join(rootPath, relativePath)

	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	node := &TreeNode{
		Name:  path.Base(relativePath),
		Path:  "/" + strings.TrimPrefix(relativePath, "/"),
		IsDir: true,
	}

	if relativePath == "" {
		node.Name = "/"
		node.Path = "/"
	}

	for _, file := range files {
		if file.IsDir() {
			childRelativePath := path.Join(relativePath, file.Name())
			child, err := buildDirectoryTree(rootPath, childRelativePath)
			if err == nil {
				node.Children = append(node.Children, child)
			}
		}
	}

	return node, nil
}
