package handler

import (
	"embed"

	"github.com/cloudwego/hertz/pkg/route"
	"github.com/zy84338719/upftp/internal/middleware"
	"github.com/zy84338719/upftp/internal/service"
)

//go:embed templates
var templates embed.FS

var svc *service.Service

func SetService(s *service.Service) {
	svc = s
}

func GetService() *service.Service {
	return svc
}

func RegisterRoutes(r *route.Engine) {
	r.Static("/static", "./internal/handler/templates")

	r.GET("/login", HandleLoginPage)
	r.POST("/api/login", HandleLogin)
	r.POST("/api/logout", HandleLogout)
	r.GET("/api/settings", HandleGetSettings)
	r.POST("/api/settings/language", HandleSetLanguage)
	r.POST("/api/settings/http-auth", HandleSetHTTPAuth)
	r.POST("/api/settings/ftp", HandleSetFTP)

	auth := r.Group("/", middleware.AuthMiddleware(svc.Config(), svc.Sessions()))

	auth.GET("/", HandleIndexPage)
	auth.GET("/api/info", HandleServerInfo)
	auth.GET("/api/files", HandleFileListAPI)
	auth.GET("/api/tree", HandleDirectoryTree)
	auth.GET("/api/qrcode", HandleQRCode)
	auth.POST("/api/upload", handleUpload)
	auth.POST("/api/create-folder", handleCreateFolder)
	auth.POST("/api/delete", handleDelete)
	auth.POST("/api/rename", handleRename)
	auth.GET("/download/*path", handleDownload)
	auth.GET("/preview/*path", handlePreview)
	auth.GET("/files/*path", handleFiles)
}
