package handlers

import (
	"archive/zip"
	"context"
	"embed"
	"encoding/base64"
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

	"github.com/skip2/go-qrcode"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/filehandler"
	"github.com/zy84338719/upftp/internal/logger"
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
	// 公开路由（无需认证）
	mux.HandleFunc("/login", HandleLoginPage)
	mux.HandleFunc("/api/login", HandleLogin)
	mux.HandleFunc("/logout", HandleLogout)

	// 现代化页面路由（暂时公开，	mux.HandleFunc("/modern/", handleModernIndex)
	mux.HandleFunc("/modern/*", handleModernIndex)

	// 保护的 API 端点
	mux.HandleFunc("/api/info", withAuth(HandleServerInfo))
	mux.HandleFunc("/api/tree", withAuth(HandleDirectoryTree))
	mux.HandleFunc("/api/upload", withAuth(handleUpload))
	mux.HandleFunc("/api/qrcode", withAuth(handleQRCode))
	mux.HandleFunc("/api/create-folder", withAuth(handleCreateFolder))
	mux.HandleFunc("/api/delete", withAuth(handleDelete))
	mux.HandleFunc("/api/rename", withAuth(handleRename))

	// 保护的文件访问
	mux.HandleFunc("/files/", withAuth(handleFiles))

	// 保护的主要路由
	mux.HandleFunc("/", withAuth(handleIndex))
	mux.HandleFunc("/download/", withAuth(handleDownload))
	mux.HandleFunc("/preview/", withAuth(handlePreview))
}

func withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !config.AppConfig.HTTPAuth.Enabled {
			next(w, r)
			return
		}

		// Try session-based authentication first
		if cookie, err := r.Cookie("auth_token"); err == nil {
			if session, valid := sessionManager.ValidateSession(cookie.Value); valid {
				// Add username to context
				ctx := context.WithValue(r.Context(), "username", session.Username)
				next(w, r.WithContext(ctx))
				return
			}
		}

		// Fall back to Basic Auth for API clients
		if user, pass, ok := r.BasicAuth(); ok {
			if user == config.AppConfig.HTTPAuth.Username && pass == config.AppConfig.HTTPAuth.Password {
				next(w, r)
				return
			}
		}

		// Authentication failed
		// For API requests, return JSON error
		if r.Header.Get("Accept") == "application/json" || r.Header.Get("Content-Type") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized",
			})
		} else {
			// For web requests, redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		logger.Warn("Unauthorized access attempt from %s", r.RemoteAddr)
	}
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
		logger.Info("Downloaded: %s", filename)
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
	logger.Info("Downloaded directory as ZIP: %s", filename)
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

// handleFiles securely serves static file access with authentication
func handleFiles(w http.ResponseWriter, r *http.Request) {
	// Extract file path
	filename := strings.TrimPrefix(r.URL.Path, "/files/")

	// Security check: prevent path traversal attacks
	if !filehandler.IsPathSafe(filename) {
		http.Error(w, "Access denied", http.StatusForbidden)
		logger.Warn("Blocked unsafe file access: %s from %s", filename, r.RemoteAddr)
		return
	}

	// Build full file path
	filePath := path.Join(config.AppConfig.Root, filename)

	// Serve the file
	http.ServeFile(w, r, filePath)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if !config.AppConfig.Upload.Enabled {
		http.Error(w, `{"error": "Upload disabled"}`, http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, config.AppConfig.Upload.MaxSize)

	if err := r.ParseMultipartForm(config.AppConfig.Upload.MaxSize); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "File too large. Max: %d bytes"}`, config.AppConfig.Upload.MaxSize), http.StatusBadRequest)
		return
	}

	uploadPath := r.FormValue("path")
	if uploadPath == "" {
		uploadPath = "/"
	}

	if !filehandler.IsPathSafe(uploadPath) && uploadPath != "/" {
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}

	targetDir := path.Join(config.AppConfig.Root, uploadPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		http.Error(w, `{"error": "Cannot create directory"}`, http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, `{"error": "No files uploaded"}`, http.StatusBadRequest)
		return
	}

	var uploaded []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		filePath := path.Join(targetDir, fileHeader.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			file.Close()
			continue
		}

		io.Copy(dst, file)
		dst.Close()
		file.Close()
		uploaded = append(uploaded, fileHeader.Filename)
		logger.Info("Uploaded: %s", path.Join(uploadPath, fileHeader.Filename))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"uploaded": uploaded,
		"count":    len(uploaded),
	})
}

func handleQRCode(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("http://%s:%d", serverInfo.IP, serverInfo.HTTPPort)

	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

func handleCreateFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	if !filehandler.IsPathSafe(req.Path) && req.Path != "/" {
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}

	folderPath := path.Join(config.AppConfig.Root, req.Path, req.Name)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	logger.Info("Created folder: %s", path.Join(req.Path, req.Name))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	if !filehandler.IsPathSafe(req.Path) {
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}

	fullPath := path.Join(config.AppConfig.Root, req.Path)
	if err := os.RemoveAll(fullPath); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	logger.Info("Deleted: %s", req.Path)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleRename(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Path    string `json:"path"`
		NewName string `json:"newName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	if !filehandler.IsPathSafe(req.Path) {
		http.Error(w, `{"error": "Invalid path"}`, http.StatusBadRequest)
		return
	}

	oldPath := path.Join(config.AppConfig.Root, req.Path)
	newPath := path.Join(path.Dir(oldPath), req.NewName)

	if err := os.Rename(oldPath, newPath); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	logger.Info("Renamed: %s -> %s", req.Path, req.NewName)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func HandleServerInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qrCodeBase64 := ""
	if png, err := qrcode.Encode(fmt.Sprintf("http://%s:%d", serverInfo.IP, serverInfo.HTTPPort), qrcode.Medium, 256); err == nil {
		qrCodeBase64 = base64.StdEncoding.EncodeToString(png)
	}

	response := map[string]interface{}{
		"version":         config.AppConfig.Version,
		"lastCommit":      config.AppConfig.LastCommit,
		"httpPort":        serverInfo.HTTPPort,
		"ftpPort":         serverInfo.FTPPort,
		"ftpEnabled":      config.AppConfig.EnableFTP,
		"rootPath":        config.AppConfig.Root,
		"uploadEnabled":   config.AppConfig.Upload.Enabled,
		"httpAuthEnabled": config.AppConfig.HTTPAuth.Enabled,
		"qrCode":          qrCodeBase64,
	}

	json.NewEncoder(w).Encode(response)
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
