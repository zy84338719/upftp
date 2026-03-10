package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/zy84338719/upftp/internal/cli"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/mcp"
	"github.com/zy84338719/upftp/internal/network"
	"github.com/zy84338719/upftp/internal/server"
)

var (
	Version    = "undefined"
	LastCommit = "undefined"
)

func main() {
	config.Init(Version, LastCommit)
	logger.Init(config.AppConfig.Logging.Level, config.AppConfig.Logging.Format)

	if config.AppConfig.EnableMCP {
		mcpServer := mcp.NewMCPServer()
		if err := mcpServer.Start(context.Background()); err != nil {
			logger.Fatal("MCP server error: %v", err)
		}
		return
	}

	selectedIP, err := network.GetInfo(
		config.AppConfig.AutoSelect,
		config.AppConfig.GetHTTPPort(),
		config.AppConfig.GetFTPPort(),
	)
	if err != nil {
		logger.Fatal("Failed to get network information: %v", err)
	}

	logger.Info("UPFTP %s starting...", config.AppConfig.Version)
	logger.Info("Shared directory: %s", config.AppConfig.Root)

	ctx, cancelFunc := context.WithCancel(context.Background())

	httpServer := server.NewHTTPServer()
	go func() {
		if err := httpServer.Start(ctx, selectedIP,
			config.AppConfig.GetHTTPPort(),
			config.AppConfig.GetFTPPort(),
			config.AppConfig.Root); err != nil {
			logger.Error("HTTP server error: %v", err)
		}
	}()

	if config.AppConfig.EnableFTP {
		ftpServer := server.NewFTPServer()
		go func() {
			if err := ftpServer.Start(ctx, selectedIP,
				config.AppConfig.GetFTPPort(),
				config.AppConfig.Root,
				config.AppConfig.Username,
				config.AppConfig.Password); err != nil {
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
		}
	}
}
