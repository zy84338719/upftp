package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/filehandler"
	"github.com/zy84338719/upftp/internal/logger"
)

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
