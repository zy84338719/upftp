package handler

import (
	"context"
	"net/url"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/skip2/go-qrcode"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/model"
	"github.com/zy84338719/upftp/internal/response"
)

func HandleServerInfo(ctx context.Context, c *app.RequestContext) {
	response.Success(c, svc.GetServerInfo())
}

func HandleFileListAPI(ctx context.Context, c *app.RequestContext) {
	reqPath := c.DefaultQuery("path", "/")
	if reqPath == "" {
		reqPath = "/"
	}

	legacy, err := svc.ListFilesLegacy(reqPath)
	if err != nil {
		response.InternalError(c, "Failed to list files: "+err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"files": legacy,
		"path":  reqPath,
	})
}

func HandleDirectoryTree(ctx context.Context, c *app.RequestContext) {
	rootPath := c.DefaultQuery("path", "/")
	if rootPath == "" {
		rootPath = "/"
	}

	tree, err := svc.BuildTree(rootPath, 0)
	if err != nil {
		response.InternalError(c, "Failed to build tree: "+err.Error())
		return
	}

	if tree == nil {
		tree = &model.TreeNode{Name: "/", Path: "/", IsDir: true}
	}

	response.Success(c, tree)
}

func HandleQRCode(ctx context.Context, c *app.RequestContext) {
	rawURL := c.DefaultQuery("url", "")
	if rawURL == "" {
		c.String(consts.StatusBadRequest, "missing url parameter")
		return
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		c.String(consts.StatusBadRequest, "invalid url")
		return
	}

	png, err := qrcode.Encode(rawURL, qrcode.Medium, 256)
	if err != nil {
		logger.Error("QR code generation error: %v", err)
		c.String(consts.StatusInternalServerError, "failed to generate QR code")
		return
	}

	c.Header("Content-Type", "image/png")
	c.SetStatusCode(consts.StatusOK)
	c.Response.BodyWriter().Write(png)
}
