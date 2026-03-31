package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/zy84338719/upftp/biz/handler/static"
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

	// 注册测试路由
	h.GET("/test", func(ctx context.Context, c *app.RequestContext) {
		c.String(consts.StatusOK, "Test route works!")
	})

	// 注册静态文件路由
	static.RegisterStaticRoutes(h)

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
