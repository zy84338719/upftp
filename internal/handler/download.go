package handler

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/model"
)

func handleDownload(ctx context.Context, c *app.RequestContext) {
	filename := string(c.Param("path"))
	if filename == "" {
		filename = string(c.Param("filename"))
	}

	info, err := svc.Stat(filename)
	if err != nil {
		c.String(consts.StatusNotFound, "File not found")
		return
	}

	fullPath, err := svc.SafePath(filename)
	if err != nil {
		c.String(consts.StatusNotFound, "File not found")
		return
	}

	if !info.IsDir() {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filename)))
		c.File(fullPath)
		logger.Info("Downloaded: %s", filename)
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", filepath.Base(filename)))

	zipWriter := zip.NewWriter(c.Response.BodyWriter())
	defer zipWriter.Close()

	err = filepath.Walk(fullPath, func(walkPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		relPath, relErr := filepath.Rel(fullPath, walkPath)
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
			_, copyErr := io.Copy(writer, file)
			file.Close()
			if copyErr != nil {
				return copyErr
			}
		}
		return nil
	})

	if err != nil {
		logger.Error("ZIP creation error: %v", err)
		c.String(consts.StatusInternalServerError, "Error creating zip file")
		return
	}
	logger.Info("Downloaded directory as ZIP: %s", filename)
}

func handlePreview(ctx context.Context, c *app.RequestContext) {
	filename := string(c.Param("path"))
	if filename == "" {
		filename = string(c.Param("filename"))
	}

	fullPath, err := svc.Download(filename)
	if err != nil {
		c.String(consts.StatusNotFound, "File not found")
		return
	}

	fileType := model.GetFileType(filename)
	mimeType := model.GetMimeType(fileType)
	c.Header("Content-Type", mimeType)
	c.File(fullPath)
}

func handleFiles(ctx context.Context, c *app.RequestContext) {
	filename := string(c.Param("path"))
	if filename == "" {
		filename = string(c.Param("filename"))
	}

	fullPath, err := svc.Download(filename)
	if err != nil {
		c.String(consts.StatusForbidden, "Access denied")
		logger.Warn("Blocked file access: %s from %s", filename, c.ClientIP())
		return
	}

	c.File(fullPath)
}
