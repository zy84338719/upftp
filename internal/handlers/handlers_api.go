package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/skip2/go-qrcode"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/filehandler"
)

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
