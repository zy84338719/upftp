package static

import (
	"context"
	"io/fs"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/biz/handler/index"
)

// ServeAssets 处理 /assets/*filepath 路由，从嵌入文件系统提供静态文件
func ServeAssets(ctx context.Context, c *app.RequestContext) {
	path := c.Param("filepath")

	templateFS := index.GetTemplatesFS()
	data, err := fs.ReadFile(templateFS, "assets/"+path)
	if err != nil {
		c.String(consts.StatusNotFound, "File not found")
		return
	}

	contentType := getContentType(path)
	c.Header("Content-Type", contentType)
	c.SetStatusCode(consts.StatusOK)
	c.Response.BodyWriter().Write(data)
}

// ServeFavicon 处理 /favicon.ico 路由
func ServeFavicon(ctx context.Context, c *app.RequestContext) {
	templateFS := index.GetTemplatesFS()
	data, err := fs.ReadFile(templateFS, "favicon.ico")
	if err != nil {
		c.String(consts.StatusNotFound, "File not found")
		return
	}
	c.Header("Content-Type", "image/x-icon")
	c.SetStatusCode(consts.StatusOK)
	c.Response.BodyWriter().Write(data)
}

// getContentType 根据文件扩展名返回对应的 Content-Type
func getContentType(path string) string {
	switch {
	case strings.HasSuffix(path, ".css"):
		return "text/css"
	case strings.HasSuffix(path, ".js"):
		return "application/javascript"
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	case strings.HasSuffix(path, ".png"):
		return "image/png"
	case strings.HasSuffix(path, ".jpg"), strings.HasSuffix(path, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".woff"):
		return "font/woff"
	case strings.HasSuffix(path, ".woff2"):
		return "font/woff2"
	case strings.HasSuffix(path, ".ttf"):
		return "font/ttf"
	case strings.HasSuffix(path, ".eot"):
		return "application/vnd.ms-fontobject"
	default:
		return "application/octet-stream"
	}
}

// RegisterStaticRoutes 注册所有静态文件路由
func RegisterStaticRoutes(h *server.Hertz) {
	h.GET("/assets/*filepath", ServeAssets)
	h.GET("/favicon.ico", ServeFavicon)
}
