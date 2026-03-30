package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/filehandler"
	"github.com/zy84338719/upftp/internal/logger"
)

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
