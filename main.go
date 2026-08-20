// Command upftp 是一个快速、即开即用的 FTP 文件分享工具。
//
// 单二进制同时提供:FTP 服务器、REST API、极简 Web 界面。
// 默认匿名访问、自动选端口、共享当前目录。
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/zy84338719/upftp/internal/cli"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/core"
	httpserver "github.com/zy84338719/upftp/internal/http"
	"github.com/zy84338719/upftp/internal/netinfo"
)

var (
	version = "dev"
)

func main() {
	cfg := config.Default()

	// 先加载配置文件(若有),使其成为 flag 的默认值,这样命令行自然优先于文件。
	for _, p := range config.DefaultConfigPaths() {
		if _, err := os.Stat(p); err == nil {
			_ = cfg.LoadFromFile(p)
			break
		}
	}

	// 命令行 flag(以 cfg 当前值作默认 → 命令行 > 配置文件 > 内置默认)
	flag.StringVar(&cfg.Dir, "d", cfg.Dir, "共享目录")
	flag.IntVar(&cfg.FTPPort, "p", cfg.FTPPort, "FTP 端口")
	flag.IntVar(&cfg.HTTPPort, "http", cfg.HTTPPort, "HTTP/Web 端口")
	flag.StringVar(&cfg.User, "user", cfg.User, "用户名(留空=匿名)")
	flag.StringVar(&cfg.Pass, "pass", cfg.Pass, "密码")
	flag.BoolVar(&cfg.ReadOnly, "read-only", cfg.ReadOnly, "只读模式")
	flag.BoolVar(&cfg.TUI, "tui", cfg.TUI, "启用终端交互界面")
	flag.StringVar(&cfg.LogLevel, "log", cfg.LogLevel, "日志级别(debug/info/warn/error)")
	showVersion := flag.Bool("version", false, "打印版本")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "upftp %s — 快速 FTP 文件分享\n\n", version)
		fmt.Fprintf(os.Stderr, "用法: upftp [选项]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  upftp                         # 当前目录,匿名,默认端口\n")
		fmt.Fprintf(os.Stderr, "  upftp -d ~/share -p 2121      # 指定目录与端口\n")
		fmt.Fprintf(os.Stderr, "  upftp -user me -pass secret   # 启用密码\n")
	}
	flag.Parse()

	if *showVersion {
		fmt.Println("upftp", version)
		return
	}

	// 日志
	level := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	host := cfg.Host
	if host == "" {
		host = netinfo.BestIP()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	auth := &core.Auth{Anonymous: !cfg.AuthEnabled(), User: cfg.User, Pass: cfg.Pass}
	mgr := core.NewSessionManager(host, logger)
	files := &core.Files{Root: cfg.Dir, ReadOnly: cfg.ReadOnly}

	// 启动默认 FTP 会话(共享指定目录)
	defaultSess, err := mgr.Create(ctx, core.SessionSpec{
		Dir:       cfg.Dir,
		Port:      cfg.FTPPort,
		Anonymous: !cfg.AuthEnabled(),
		User:      cfg.User,
		Pass:      cfg.Pass,
		ReadOnly:  cfg.ReadOnly,
	})
	if err != nil {
		logger.Error("启动 FTP 失败", "err", err)
		os.Exit(1)
	}

	// 启动 HTTP
	httpSrv := httpserver.New(httpserver.Options{
		Manager: mgr, Files: files, Auth: auth, Logger: logger,
		Host: host, FTPPort: defaultSess.Port,
	})
	go func() {
		if err := httpSrv.Start(ctx, cfg.HTTPPort); err != nil {
			logger.Error("HTTP 服务退出", "err", err)
		}
	}()

	// 打印访问信息
	fmt.Println()
	fmt.Printf("  📁 upftp %s 已启动\n", version)
	fmt.Printf("  📂 共享目录: %s\n", cfg.Dir)
	fmt.Printf("  📡 FTP:   %s\n", defaultSess.FTPURL())
	fmt.Printf("  🌐 Web:   http://%s:%d\n", host, cfg.HTTPPort)
	if !auth.Enabled() {
		fmt.Printf("  🔓 匿名访问(无需密码)\n")
	} else {
		fmt.Printf("  🔒 凭据: %s / %s\n", cfg.User, cfg.Pass)
	}
	fmt.Println()

	// 信号处理:TUI 在前台运行会阻塞读键盘;这里并行监听信号,
	// 任一方结束(用户退出 TUI 或收到 Ctrl-C)都进入关闭流程。
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	if cfg.TUI {
		// TUI 阻塞运行,直到用户按 q 退出;之后转入纯后台模式。
		cli.Run(ctx, files, host, cfg.HTTPPort, defaultSess.FTPURL())
		fmt.Println("TUI 已退出,服务继续在后台运行。再次按 Ctrl-C 退出。")
	}

	<-sigCh
	fmt.Println("\n正在关闭...")
	mgr.StopAll()
	cancel()
}
