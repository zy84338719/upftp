package handlers

import (
	"context"
	"embed"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zy84338719/upftp/internal/config"
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
