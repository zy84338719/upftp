package handlers

import (
	"archive/zip"
	"context"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
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

func RegisterRoutes(h *route.Engine) {
	h.GET("/login", HandleLoginPage)
	h.POST("/api/login", HandleLogin)
	h.Any("/logout", HandleLogout)

	h.GET("/static/css/styles.css", handleStaticCSS("templates/css/styles.css"))
	h.GET("/static/js/app.js", handleStaticJS("templates/js/app.js"))
	h.GET("/static/css/login.css", handleStaticCSS("templates/css/login.css"))
	h.GET("/static/js/login.js", handleStaticJS("templates/js/login.js"))

	auth := h.Group("/")
	auth.Use(withAuth())
	{
		auth.GET("/api/info", HandleServerInfo)
		auth.GET("/api/tree", HandleDirectoryTree)
		auth.POST("/api/upload", handleUpload)
		auth.GET("/api/qrcode", handleQRCode)
		auth.POST("/api/create-folder", handleCreateFolder)
		auth.POST("/api/delete", handleDelete)
		auth.POST("/api/rename", handleRename)
		auth.GET("/api/files", handleFileListAPI)

		auth.GET("/api/settings", HandleGetSettings)
		auth.POST("/api/settings/language", HandleSetLanguage)
		auth.POST("/api/settings/http-auth", HandleSetHTTPAuth)
		auth.POST("/api/settings/ftp", HandleSetFTP)

		auth.GET("/files/:path", handleFiles)
		auth.GET("/", handleModernIndex)
		auth.GET("/download/:path", handleDownload)
		auth.GET("/preview/:path", handlePreview)
	}
}

func withAuth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !config.AppConfig.HTTPAuth.Enabled {
			c.Next(ctx)
			return
		}

		if cookie := string(c.Cookie("auth_token")); cookie != "" {
			if session, valid := sessionManager.ValidateSession(cookie); valid {
				c.Set("username", session.Username)
				c.Next(ctx)
				return
			}
		}

		if user, pass, ok := c.Request.BasicAuth(); ok {
			if user == config.AppConfig.HTTPAuth.Username && pass == config.AppConfig.HTTPAuth.Password {
				c.Next(ctx)
				return
			}
		}

		accept := string(c.GetHeader("Accept"))
		contentType := string(c.GetHeader("Content-Type"))
		if strings.Contains(accept, "application/json") || strings.Contains(contentType, "application/json") {
			c.Header("Content-Type", "application/json")
			c.SetStatusCode(consts.StatusUnauthorized)
			_, _ = c.WriteString(`{"error":"Unauthorized"}`)
		} else {
			c.Redirect(consts.StatusSeeOther, []byte("/login"))
		}
		logger.Warn("Unauthorized access attempt from %s", c.ClientIP())
	}
}

func handleDownload(ctx context.Context, c *app.RequestContext) {
	filename := c.Param("path")
	if !filehandler.IsPathSafe(filename) {
		c.String(consts.StatusForbidden, "Access denied")
		return
	}

	filePath := path.Join(config.AppConfig.Root, filename)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		c.String(consts.StatusNotFound, "File not found")
		return
	}

	if !fileInfo.IsDir() {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filename)))
		c.File(filePath)
		logger.Info("Downloaded: %s", filename)
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", filepath.Base(filename)))

	zipWriter := zip.NewWriter(c.Response.BodyWriter())
	defer zipWriter.Close()

	err = filepath.Walk(filePath, func(walkPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relPath, relErr := filepath.Rel(filePath, walkPath)
		if relErr != nil {
			return relErr
		}

		if relPath == "." {
			return nil
		}

		header, hdrErr := zip.FileInfoHeader(info)
		if hdrErr != nil {
			return hdrErr
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, createErr := zipWriter.CreateHeader(header)
		if createErr != nil {
			return createErr
		}

		if !info.IsDir() {
			file, openErr := os.Open(walkPath)
			if openErr != nil {
				return openErr
			}
			defer file.Close()
			_, _ = io.Copy(writer, file)
		}
		return nil
	})

	if err != nil {
		c.String(consts.StatusInternalServerError, "Error creating zip file")
		return
	}
	logger.Info("Downloaded directory as ZIP: %s", filename)
}

func handlePreview(ctx context.Context, c *app.RequestContext) {
	filename := c.Param("path")
	if !filehandler.IsPathSafe(filename) {
		c.String(consts.StatusForbidden, "Access denied")
		return
	}

	filePath := path.Join(config.AppConfig.Root, filename)
	fileType := filehandler.GetFileType(filename)
	mimeType := filehandler.GetMimeType(fileType)
	c.Header("Content-Type", mimeType)
	c.File(filePath)
}

func handleFiles(ctx context.Context, c *app.RequestContext) {
	filename := c.Param("path")
	if !filehandler.IsPathSafe(filename) {
		c.String(consts.StatusForbidden, "Access denied")
		logger.Warn("Blocked unsafe file access: %s from %s", filename, c.ClientIP())
		return
	}

	filePath := path.Join(config.AppConfig.Root, filename)
	c.File(filePath)
}

func handleUpload(ctx context.Context, c *app.RequestContext) {
	if !config.AppConfig.Upload.Enabled {
		c.String(consts.StatusForbidden, `{"error": "Upload disabled"}`)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.String(consts.StatusBadRequest, fmt.Sprintf(`{"error": "Failed to parse form: %s"}`, err.Error()))
		return
	}

	uploadPath := form.Value["path"]
	var targetPath string
	if len(uploadPath) > 0 && uploadPath[0] != "" {
		targetPath = uploadPath[0]
	} else {
		targetPath = "/"
	}

	if !filehandler.IsPathSafe(targetPath) && targetPath != "/" {
		c.String(consts.StatusBadRequest, `{"error": "Invalid path"}`)
		return
	}

	targetDir := path.Join(config.AppConfig.Root, targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		c.String(consts.StatusInternalServerError, `{"error": "Cannot create directory"}`)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.String(consts.StatusBadRequest, `{"error": "No files uploaded"}`)
		return
	}

	var uploaded []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		dstPath := path.Join(targetDir, fileHeader.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			file.Close()
			continue
		}

		_, _ = io.Copy(dst, file)
		dst.Close()
		file.Close()
		uploaded = append(uploaded, fileHeader.Filename)
		logger.Info("Uploaded: %s", path.Join(targetPath, fileHeader.Filename))
	}

	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]interface{}{
		"success":  true,
		"uploaded": uploaded,
		"count":    len(uploaded),
	})
}

func handleQRCode(ctx context.Context, c *app.RequestContext) {
	target := string(c.Query("url"))
	if target == "" {
		target = fmt.Sprintf("http://%s:%d", serverInfo.IP, serverInfo.HTTPPort)
	}

	png, err := qrcode.Encode(target, qrcode.Medium, 256)
	if err != nil {
		c.String(consts.StatusInternalServerError, "Failed to generate QR code")
		return
	}

	c.Header("Content-Type", "image/png")
	c.Response.BodyWriter().Write(png)
}

func handleCreateFolder(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error": "Invalid request"}`)
		return
	}
	var req struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error": "Invalid request"}`)
		return
	}

	if !filehandler.IsPathSafe(req.Path) && req.Path != "/" {
		c.String(consts.StatusBadRequest, `{"error": "Invalid path"}`)
		return
	}

	folderPath := path.Join(config.AppConfig.Root, req.Path, req.Name)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	logger.Info("Created folder: %s", path.Join(req.Path, req.Name))
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}

func handleDelete(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error": "Invalid request"}`)
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error": "Invalid request"}`)
		return
	}

	if !filehandler.IsPathSafe(req.Path) {
		c.String(consts.StatusBadRequest, `{"error": "Invalid path"}`)
		return
	}

	fullPath := path.Join(config.AppConfig.Root, req.Path)
	if err := os.RemoveAll(fullPath); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	logger.Info("Deleted: %s", req.Path)
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}

func handleRename(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error": "Invalid request"}`)
		return
	}
	var req struct {
		Path    string `json:"path"`
		NewName string `json:"newName"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error": "Invalid request"}`)
		return
	}

	if !filehandler.IsPathSafe(req.Path) {
		c.String(consts.StatusBadRequest, `{"error": "Invalid path"}`)
		return
	}

	oldPath := path.Join(config.AppConfig.Root, req.Path)
	newPath := path.Join(path.Dir(oldPath), req.NewName)

	if err := os.Rename(oldPath, newPath); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
		return
	}

	logger.Info("Renamed: %s -> %s", req.Path, req.NewName)
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}

func HandleServerInfo(ctx context.Context, c *app.RequestContext) {
	c.Header("Content-Type", "application/json")

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
		"httpAuthUser":    config.AppConfig.HTTPAuth.Username,
		"httpAuthPass":    config.AppConfig.HTTPAuth.Password,
		"qrCode":          qrCodeBase64,
	}

	json.NewEncoder(c.Response.BodyWriter()).Encode(response)
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

func HandleDirectoryTree(ctx context.Context, c *app.RequestContext) {
	c.Header("Content-Type", "application/json")

	tree, err := buildDirectoryTree(config.AppConfig.Root, "")
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	jsonBytes, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	c.Response.BodyWriter().Write(jsonBytes)
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

func handleFileListAPI(ctx context.Context, c *app.RequestContext) {
	urlPath := string(c.Query("path"))
	if urlPath == "" {
		urlPath = "/"
	}

	fsPath := path.Join(config.AppConfig.Root, urlPath)
	fileInfo, err := os.Stat(fsPath)
	if err != nil || !fileInfo.IsDir() {
		c.Header("Content-Type", "application/json")
		json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]interface{}{"error": "not found"})
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

	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]interface{}{
		"files": fileList,
		"path":  urlPath,
	})
}

func HandleGetSettings(ctx context.Context, c *app.RequestContext) {
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]interface{}{
		"language":     config.AppConfig.GetLanguage(),
		"httpAuth":     config.AppConfig.HTTPAuth.Enabled,
		"httpAuthUser": config.AppConfig.HTTPAuth.Username,
		"httpAuthPass": config.AppConfig.HTTPAuth.Password,
		"ftpUser":      config.AppConfig.Username,
		"ftpPass":      config.AppConfig.Password,
		"ftpEnabled":   config.AppConfig.EnableFTP,
		"mcpEnabled":   config.AppConfig.EnableMCP,
		"httpPort":     config.AppConfig.GetHTTPPort(),
		"ftpPort":      config.AppConfig.GetFTPPort(),
	})
}

func HandleSetLanguage(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	var req struct {
		Language string `json:"language"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	if req.Language != "en" && req.Language != "zh" {
		c.String(consts.StatusBadRequest, `{"error":"Invalid language"}`)
		return
	}
	config.AppConfig.Language = req.Language
	if err := config.SaveConfig(); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}

func HandleSetHTTPAuth(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	var req struct {
		Enabled  bool   `json:"enabled"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	config.AppConfig.HTTPAuth.Enabled = req.Enabled
	if req.Username != "" {
		config.AppConfig.HTTPAuth.Username = req.Username
	}
	if req.Password != "" {
		config.AppConfig.HTTPAuth.Password = req.Password
	}
	if err := config.SaveConfig(); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}

func handleStaticCSS(filepath string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		data, err := templates.ReadFile(filepath)
		if err != nil {
			c.String(consts.StatusNotFound, "Not found")
			return
		}
		c.Header("Content-Type", "text/css; charset=utf-8")
		c.Response.BodyWriter().Write(data)
	}
}

func handleStaticJS(filepath string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		data, err := templates.ReadFile(filepath)
		if err != nil {
			c.String(consts.StatusNotFound, "Not found")
			return
		}
		c.Header("Content-Type", "application/javascript; charset=utf-8")
		c.Response.BodyWriter().Write(data)
	}
}

func HandleSetFTP(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		c.String(consts.StatusBadRequest, `{"error":"Invalid request"}`)
		return
	}
	if req.Username != "" {
		config.AppConfig.Username = req.Username
	}
	if req.Password != "" {
		config.AppConfig.Password = req.Password
	}
	if err := config.SaveConfig(); err != nil {
		c.String(consts.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}
	c.Header("Content-Type", "application/json")
	json.NewEncoder(c.Response.BodyWriter()).Encode(map[string]bool{"success": true})
}
