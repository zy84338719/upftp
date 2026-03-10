package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zy84338719/upftp/internal/cli"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/network"
	"github.com/zy84338719/upftp/internal/server"
)

var (
	Version    = "undefined"
	LastCommit = "undefined"
)

func main() {
	config.Init(Version, LastCommit)

	selectedIP, err := network.GetInfo(
		config.AppConfig.AutoSelect,
		config.AppConfig.GetHTTPPort(),
		config.AppConfig.GetFTPPort(),
	)
	if err != nil {
		log.Fatal("Failed to get network information:", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	httpServer := server.NewHTTPServer()
	go func() {
		if err := httpServer.Start(ctx, selectedIP,
			config.AppConfig.GetHTTPPort(),
			config.AppConfig.GetFTPPort(),
			config.AppConfig.Root); err != nil {
			log.Printf("HTTP server error: %v", err)
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
				log.Printf("FTP server error: %v", err)
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
		log.Printf("Received signal: %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancelFunc()
			log.Println("UPFTP server shutdown complete.")
			return
		case syscall.SIGHUP:
			log.Println("Received SIGHUP, reloading configuration...")
		}
	}
}
