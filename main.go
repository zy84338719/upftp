package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zy84338719/upftp/config"
	"github.com/zy84338719/upftp/ftp"
	"github.com/zy84338719/upftp/logic"
	"github.com/zy84338719/upftp/network"
)

var (
	Version    = "undefined"
	LastCommit = "undefined"
)

func main() {
	// 初始化配置
	config.Init(Version, LastCommit)

	// 获取网络信息
	selectedIP, err := network.GetNetworkInfo(
		config.AppConfig.AutoSelect,
		config.AppConfig.GetHTTPPort(),
		config.AppConfig.GetFTPPort(),
	)
	if err != nil {
		log.Fatal("Failed to get network information:", err)
	}

	// 设置服务器信息
	logic.SetServerInfo(
		selectedIP,
		config.AppConfig.GetHTTPPort(),
		config.AppConfig.GetFTPPort(),
		config.AppConfig.Root,
	)
	logic.SetServerIP(selectedIP)

	// 创建上下文
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)

	// 启动HTTP服务器
	go func() {
		if err := logic.StartHTTPServer(ctx); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// 启动FTP服务器（如果启用）
	if config.AppConfig.EnableFTP {
		go func() {
			if err := ftp.StartFTPServer(
				ctx,
				selectedIP,
				config.AppConfig.GetFTPPort(),
				config.AppConfig.Root,
				config.AppConfig.Username,
				config.AppConfig.Password,
			); err != nil {
				log.Printf("FTP server error: %v", err)
			}
		}()
	}

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	// 启动命令行界面
	go logic.StartCommandInterface(ctx, sigChan)

	// 等待退出信号
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
			// 这里可以添加重新加载配置的逻辑
		}
	}
}
