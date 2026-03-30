package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/response"
)

type renameRequest struct {
	Path    string `json:"path"`
	NewName string `json:"newName"`
}

type createFolderRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func handleUpload(ctx context.Context, c *app.RequestContext) {
	if !svc.UploadEnabled() {
		c.Header("Content-Type", "application/json")
		c.String(consts.StatusForbidden, `{"error":"Upload disabled"}`)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.String(consts.StatusBadRequest, fmt.Sprintf(`{"error":"Failed to parse form: %s"}`, err.Error()))
		return
	}

	uploadPath := form.Value["path"]
	targetPath := "/"
	if len(uploadPath) > 0 && uploadPath[0] != "" {
		targetPath = uploadPath[0]
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.Header("Content-Type", "application/json")
		c.String(consts.StatusBadRequest, `{"error":"No files uploaded"}`)
		return
	}

	var uploaded []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			logger.Error("Failed to open uploaded file %s: %v", fileHeader.Filename, err)
			continue
		}

		fullPath := filepath.Join(targetPath, fileHeader.Filename)
		if err := svc.Upload(fullPath, file, fileHeader.Size); err != nil {
			logger.Error("Upload error for %s: %v", fileHeader.Filename, err)
			file.Close()
			continue
		}
		file.Close()
		uploaded = append(uploaded, fileHeader.Filename)
		logger.Info("Uploaded: %s", filepath.Join(targetPath, fileHeader.Filename))
	}

	response.Success(c, map[string]interface{}{
		"success":  true,
		"uploaded": uploaded,
		"count":    len(uploaded),
	})
}

func handleCreateFolder(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	var req createFolderRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	folderPath := filepath.Join(req.Path, req.Name)
	if err := svc.CreateFolder(folderPath); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	logger.Info("Created folder: %s", folderPath)
	response.Success(c, map[string]bool{"success": true})
}

func handleDelete(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}
	var req struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	if err := svc.Delete(req.Path); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	logger.Info("Deleted: %s", req.Path)
	response.Success(c, map[string]bool{"success": true})
}

func handleRename(ctx context.Context, c *app.RequestContext) {
	body, err := c.Body()
	if err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	var req renameRequest
	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(c, "Invalid request")
		return
	}

	if err := svc.Rename(req.Path, req.NewName); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	logger.Info("Renamed: %s -> %s", req.Path, req.NewName)
	response.Success(c, map[string]bool{"success": true})
}
