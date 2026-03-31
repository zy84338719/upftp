package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/biz/handler/index"
	"github.com/zy84338719/upftp/biz/router"
	"github.com/zy84338719/upftp/biz/service/ftp"
	"github.com/zy84338719/upftp/biz/service/nfs"
	"github.com/zy84338719/upftp/biz/service/webdav"
	"github.com/zy84338719/upftp/pkg/cli"
	"github.com/zy84338719/upftp/pkg/conf"
	"github.com/zy84338719/upftp/pkg/logger"
	"github.com/zy84338719/upftp/pkg/network"
)

var (
	Version     = "undefined"
	LastCommit  = "undefined"
	BuildDate   = "undefined"
	GoVersion   = runtime.Version()
	Platform    = runtime.GOOS + "/" + runtime.GOARCH
	ProjectURL  = "https://github.com/zy84338719/upftp"
	ProjectName = "UPFTP"
)

// serveEmbeddedStatic 从嵌入文件系统提供静态文件
func serveEmbeddedStatic(ctx context.Context, c *app.RequestContext) {
	path := c.Param("filepath")
	if path == "" {
		path = string(c.Request.URI().Path()[1:])
	}

	// 从嵌入文件系统读取
	templateFS := index.GetTemplatesFS()
	data, err := fs.ReadFile(templateFS, path)
	if err != nil {
		// 如果文件不存在，返回 404
		c.String(consts.StatusNotFound, "File not found")
		return
	}

	// 设置正确的 Content-Type
	var contentType string
	switch {
	case strings.HasSuffix(path, ".css"):
		contentType = "text/css"
	case strings.HasSuffix(path, ".js"):
		contentType = "application/javascript"
	case strings.HasSuffix(path, ".ico"):
		contentType = "image/x-icon"
	case strings.HasSuffix(path, ".png"):
		contentType = "image/png"
	case strings.HasSuffix(path, ".jpg"), strings.HasSuffix(path, ".jpeg"):
		contentType = "image/jpeg"
	default:
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.SetStatusCode(consts.StatusOK)
	c.Response.BodyWriter().Write(data)
}

func main() {
	conf.Init(Version, LastCommit, BuildDate, GoVersion, Platform, ProjectURL, ProjectName)
	logger.Init(conf.AppConfig.Logging.Level, conf.AppConfig.Logging.Format)

	selectedIP, err := network.GetInfo(
		conf.AppConfig.AutoSelect,
		conf.AppConfig.GetHTTPPort(),
		conf.AppConfig.GetFTPPort(),
	)
	if err != nil {
		logger.Fatal("Failed to get network information: %v", err)
	}

	logger.Info("Starting UPFTP v%s", conf.AppConfig.Version)
	logger.Info("Configuration loaded from: %s", conf.GetConfigPath())
	logger.Info("Shared directory: %s", conf.AppConfig.Root)
	logger.Info("HTTP server: http://%s:%d", selectedIP, conf.AppConfig.GetHTTPPort())
	if conf.AppConfig.EnableFTP {
		logger.Info("FTP server: ftp://%s:%d", selectedIP, conf.AppConfig.GetFTPPort())
	}
	if conf.AppConfig.EnableWebDAV {
		logger.Info("WebDAV server: http://%s:%d", selectedIP, conf.AppConfig.GetWebDAVPort())
	}
	if conf.AppConfig.EnableNFS {
		logger.Info("NFS server: nfs://%s:%d", selectedIP, conf.AppConfig.GetNFSPort())
	}
	if conf.AppConfig.HTTPAuth.Enabled {
		logger.Info("HTTP authentication: enabled (user: %s)", conf.AppConfig.HTTPAuth.Username)
	}
	if conf.AppConfig.Upload.Enabled {
		logger.Info("File upload: enabled (max size: %d MB)", conf.AppConfig.Upload.MaxSize/1024/1024)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	// 启动 HTTP 服务
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf(":%d", conf.AppConfig.GetHTTPPort())),
	)

	// 注册静态文件路由 - 处理 /assets/*
	h.GET("/assets/*filepath", serveEmbeddedStatic)
	h.GET("/favicon.ico", func(ctx context.Context, c *app.RequestContext) {
		templateFS := index.GetTemplatesFS()
		data, err := fs.ReadFile(templateFS, "favicon.ico")
		if err != nil {
			c.String(consts.StatusNotFound, "File not found")
			return
		}
		c.Header("Content-Type", "image/x-icon")
		c.SetStatusCode(consts.StatusOK)
		c.Response.BodyWriter().Write(data)
	})

	// 注册路由（包括根路径 / 的处理）
	router.GeneratedRegister(h)

	go func() {
		h.Spin()
	}()

	// 启动 FTP 服务
	if conf.AppConfig.EnableFTP {
		ftpServer := ftp.NewFTPServer()
		go func() {
			if err := ftpServer.Start(ctx, selectedIP,
				conf.AppConfig.GetFTPPort(),
				conf.AppConfig.Root,
				conf.AppConfig.Username,
				conf.AppConfig.Password); err != nil {
				logger.Error("FTP server error: %v", err)
			}
		}()
	}

	// 启动 WebDAV 服务
	if conf.AppConfig.EnableWebDAV {
		webdavServer := webdav.NewWebDAVServer()
		go func() {
			if err := webdavServer.Start(ctx, selectedIP,
				conf.AppConfig.GetWebDAVPort(),
				conf.AppConfig.Root,
				conf.AppConfig.Username,
				conf.AppConfig.Password); err != nil {
				logger.Error("WebDAV server error: %v", err)
			}
		}()
	}

	// 启动 NFS 服务
	if conf.AppConfig.EnableNFS {
		nfsServer := nfs.NewNFSServer()
		go func() {
			if err := nfsServer.Start(ctx, selectedIP,
				conf.AppConfig.GetNFSPort(),
				conf.AppConfig.Root,
				conf.AppConfig.Username,
				conf.AppConfig.Password); err != nil {
				logger.Error("NFS server error: %v", err)
			}
		}()
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	// 启动 TUI 界面
	cliApp := cli.NewCLI()
	cliApp.SetServerIP(selectedIP)
	go cliApp.Start(ctx, sigChan)

	for {
		s := <-sigChan
		logger.Info("Received signal: %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancelFunc()
			logger.Info("UPFTP server shutdown complete.")
			return
		case syscall.SIGHUP:
			logger.Info("Received SIGHUP, reloading configuration...")
			conf.ReloadConfig()
			logger.Info("Configuration reloaded from: %s", conf.GetConfigPath())
		}
	}
}
