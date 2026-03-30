package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/zy84338719/upftp/internal/auth"
	"github.com/zy84338719/upftp/internal/biz"
	"github.com/zy84338719/upftp/internal/cli"
	"github.com/zy84338719/upftp/internal/conf"
	"github.com/zy84338719/upftp/internal/dal"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/network"
	"github.com/zy84338719/upftp/internal/protocol/ftp"
	"github.com/zy84338719/upftp/internal/protocol/http"
	"github.com/zy84338719/upftp/internal/protocol/mcp"
	"github.com/zy84338719/upftp/internal/service"
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

	store := dal.NewLocalFileStore()
	pathResolver := dal.NewPathResolver(conf.AppConfig.Root)
	fileSvc := biz.NewFileService(store, pathResolver, conf.AppConfig)
	sessions := auth.NewSessionManager()
	appSvc := service.New(fileSvc, conf.AppConfig, sessions)

	if conf.AppConfig.EnableMCP {
		mcpServer := mcp.NewMCPServer(appSvc)
		if err := mcpServer.Start(context.Background()); err != nil {
			logger.Fatal("MCP server error: %v", err)
		}
		return
	}

	selectedIP, err := network.GetInfo(
		conf.AppConfig.AutoSelect,
		conf.AppConfig.GetHTTPPort(),
		conf.AppConfig.GetFTPPort(),
	)
	if err != nil {
		logger.Fatal("Failed to get network information: %v", err)
	}

	appSvc.SetServerInfo(selectedIP, conf.AppConfig.GetHTTPPort(), conf.AppConfig.GetFTPPort(), conf.AppConfig.Root)

	logger.Info("Starting UPFTP v%s", conf.AppConfig.Version)
	logger.Info("Configuration loaded from: %s", conf.GetConfigPath())
	logger.Info("Shared directory: %s", conf.AppConfig.Root)
	logger.Info("HTTP server: http://%s:%d", selectedIP, conf.AppConfig.GetHTTPPort())
	if conf.AppConfig.EnableFTP {
		logger.Info("FTP server: ftp://%s:%d", selectedIP, conf.AppConfig.GetFTPPort())
	}
	if conf.AppConfig.HTTPAuth.Enabled {
		logger.Info("HTTP authentication: enabled (user: %s)", conf.AppConfig.HTTPAuth.Username)
	}
	if conf.AppConfig.Upload.Enabled {
		logger.Info("File upload: enabled (max size: %d MB)", conf.AppConfig.Upload.MaxSize/1024/1024)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	httpServer := http.NewHTTPServer(appSvc)
	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logger.Error("HTTP server error: %v", err)
		}
	}()

	if conf.AppConfig.EnableFTP {
		ftpServer := ftp.NewFTPServer(appSvc)
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

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
