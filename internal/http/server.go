// Package httpserver 提供 REST API 与极简 Web 界面,控制 FTP 会话并直接操作文件。
package httpserver

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/zy84338719/upftp/internal/core"
)

//go:embed assets/index.html
var indexHTML []byte

// Server 是 HTTP 控制端。
type Server struct {
	mgr     *core.SessionManager
	files   *core.Files
	auth    *core.Auth
	logger  *slog.Logger
	host    string
	ftpPort int // 默认 FTP 端口(信息展示用)
}

// Options 构造 HTTP 服务。
type Options struct {
	Manager *core.SessionManager
	Files   *core.Files
	Auth    *core.Auth
	Logger  *slog.Logger
	Host    string
	FTPPort int
}

// New 创建 HTTP 服务(未启动)。
func New(o Options) *Server {
	if o.Logger == nil {
		o.Logger = slog.Default()
	}
	return &Server{
		mgr: o.Manager, files: o.Files, auth: o.Auth,
		logger: o.Logger, host: o.Host, ftpPort: o.FTPPort,
	}
}

// Handler 返回路由配置好的 http.Handler。
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	api := http.NewServeMux()
	api.HandleFunc("GET /api/health", s.health)
	api.HandleFunc("GET /api/sessions", s.listSessions)
	api.HandleFunc("POST /api/sessions", s.createSession)
	api.HandleFunc("DELETE /api/sessions/{id}", s.stopSession)
	api.HandleFunc("GET /api/files", s.listFiles)
	api.HandleFunc("GET /api/download", s.download)
	api.HandleFunc("POST /api/upload", s.upload)
	api.HandleFunc("POST /api/mkdir", s.mkdir)
	api.HandleFunc("DELETE /api/files", s.deleteFile)
	mux.Handle("/api/", s.auth.HTTPMiddleware(api))

	// Web 页面 + 其余路径回退到首页(单页应用)
	mux.Handle("/", s.auth.HTTPMiddleware(http.HandlerFunc(s.serveUI)))

	return s.logMiddleware(mux)
}

// Start 阻塞运行 HTTP 服务,直到 ctx 取消。
func (s *Server) Start(ctx context.Context, port int) error {
	srv := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: s.Handler()}
	go func() {
		<-ctx.Done()
		shut, cancel := context.WithCancel(context.Background())
		defer cancel()
		_ = srv.Shutdown(shut)
	}()
	s.logger.Info("HTTP server listening", "addr", srv.Addr, "auth", s.auth.Enabled())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) serveUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(indexHTML)
}

// logMiddleware 记录每个请求的方法、路径、状态码。
func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)
		s.logger.Info("http", "method", r.Method, "path", r.URL.Path, "status", rw.status)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
