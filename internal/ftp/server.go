package ftp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

const (
	connTimeout    = 5 * time.Minute  // 单连接空闲超时
	readIdle       = 30 * time.Second // 读取下一条命令的等待
	dataTimeout    = 2 * time.Minute  // 数据连接读写超时
	maxConnections = 100              // 最大并发客户端
	bufferSize     = 32 * 1024
)

// Options 配置一个 FTPServer 实例。
type Options struct {
	Root      string // 共享根目录(绝对路径)
	Port      int    // 监听端口;0 表示由调用方提供 Listener
	Anonymous bool   // 允许匿名
	User      string // 用户名(User 非空时启用凭据认证)
	Pass      string // 密码
	ReadOnly  bool   // 只读,禁止写操作
	Logger    *slog.Logger
}

// FTPServer 是一个可独立启停的 FTP 服务器实例。
type FTPServer struct {
	opts   Options
	logger *slog.Logger

	listener net.Listener
	cancel   context.CancelFunc
	done     chan struct{}

	mu              sync.Mutex
	clients         map[*client]struct{}
	connectionCount int
}

// New 创建一个 FTP 服务器(未启动)。
func New(opts Options) *FTPServer {
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	return &FTPServer{
		opts:    opts,
		logger:  opts.Logger,
		clients: make(map[*client]struct{}),
		done:    make(chan struct{}),
	}
}

// Addr 返回监听地址(启动后可用)。
func (s *FTPServer) Addr() string {
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

// Start 阻塞运行直到 ctx 取消或监听失败。
// 若 opts.Port == 0 且未设置 listener,返回错误。
func (s *FTPServer) Start(ctx context.Context) error {
	if s.opts.Root == "" {
		return errors.New("ftp: empty root")
	}
	addr := fmt.Sprintf(":%d", s.opts.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("ftp listen :%d: %w", s.opts.Port, err)
	}
	s.listener = ln
	s.logger.Info("FTP server listening", "addr", ln.Addr().String(), "root", s.opts.Root,
		"anonymous", s.opts.Anonymous, "readonly", s.opts.ReadOnly)

	ctx, s.cancel = context.WithCancel(ctx)
	go s.reapIdle()

	// accept 循环
	go func() {
		defer close(s.done)
		for {
			conn, err := ln.Accept()
			if err != nil {
				// listener 关闭或取消 → 正常退出
				if ctx.Err() != nil {
					return
				}
				s.logger.Error("ftp accept error", "err", err)
				continue
			}
			s.accept(ctx, conn)
		}
	}()

	<-ctx.Done()
	s.logger.Info("FTP server stopping")
	_ = ln.Close()
	s.kickAll()
	<-s.doneCh()
	return nil
}

func (s *FTPServer) doneCh() chan struct{} { return s.done }

// accept 在连接数允许时接纳一个新客户端,否则直接拒绝。
func (s *FTPServer) accept(ctx context.Context, conn net.Conn) {
	s.mu.Lock()
	if s.connectionCount >= maxConnections {
		s.mu.Unlock()
		fmt.Fprintf(conn, "421 Too many connections\r\n")
		_ = conn.Close()
		return
	}
	s.connectionCount++
	c := &client{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
		server: s,
		cwd:    "/",
		name:   conn.RemoteAddr().String(),
	}
	s.clients[c] = struct{}{}
	s.mu.Unlock()

	go s.handle(ctx, c)
}

// handle 处理单个客户端连接的整个生命周期。
func (s *FTPServer) handle(ctx context.Context, c *client) {
	defer func() {
		_ = c.conn.Close()
		s.mu.Lock()
		// 修复点 #3:连接计数仅在此处减少一次,旧代码在 reaper 与 handle 各减一次会双重减。
		if _, ok := s.clients[c]; ok {
			delete(s.clients, c)
			s.connectionCount--
		}
		s.mu.Unlock()
		s.logger.Info("ftp client disconnected", "remote", c.name)
	}()

	s.logger.Info("ftp client connected", "remote", c.name)
	c.send("220 upftp ready")

	for {
		select {
		case <-ctx.Done():
			c.send("421 Service closing")
			return
		default:
		}

		_ = c.conn.SetReadDeadline(time.Now().Add(readIdle))
		line, err := c.reader.ReadString('\n')
		if err != nil {
			// 空闲超时或对端关闭,静默退出
			return
		}
		c.lastActivity = time.Now()

		line = trimCRLF(line)
		if line == "" {
			continue
		}
		cmd, args := splitCmd(line)
		s.dispatch(ctx, c, cmd, args)
		if c.quit {
			return
		}
	}
}

// kickAll 主动断开所有客户端(用于优雅关闭)。
func (s *FTPServer) kickAll() {
	s.mu.Lock()
	for c := range s.clients {
		_ = c.conn.Close()
	}
	s.mu.Unlock()
}

// reapIdle 周期清理超过 connTimeout 未活动的连接。
func (s *FTPServer) reapIdle() {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-t.C:
			s.mu.Lock()
			now := time.Now()
			for c := range s.clients {
				if now.Sub(c.lastActivity) > connTimeout {
					s.logger.Info("ftp client idle timeout", "remote", c.name)
					_ = c.conn.Close()
					// 计数由 handle 的 defer 统一减少,这里不重复操作
					delete(s.clients, c)
				}
			}
			s.mu.Unlock()
		}
	}
}

// --- 客户端结构 ---

type client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	server *FTPServer

	cwd          string
	auth         bool
	name         string
	binaryMode   bool
	restPos      int64
	lastActivity time.Time
	quit         bool

	// 数据连接协商状态
	dataListener net.Listener // PASV/EPSV 创建
	dataPort     string       // PORT/EPRT 指定的远端地址(host:port)
	rnfr         string       // RNFR 记录的源路径
}

func (c *client) send(msg string) {
	_, _ = c.writer.WriteString(msg + "\r\n")
	_ = c.writer.Flush()
}

func trimCRLF(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\r' || s[len(s)-1] == '\n') {
		s = s[:len(s)-1]
	}
	return s
}

// splitCmd 把 "RETR /foo.txt" 拆成 ("RETR", "/foo.txt")。
func splitCmd(line string) (cmd, args string) {
	for i := 0; i < len(line); i++ {
		if line[i] == ' ' {
			return line[:i], line[i+1:]
		}
	}
	return line, ""
}
