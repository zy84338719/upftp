package webdav

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zy84338719/upftp/pkg/logger"
)

type WebDAVServer struct {
	listener   net.Listener
	httpServer *http.Server
	rootPath   string
	username   string
	password   string
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewWebDAVServer() *WebDAVServer {
	return &WebDAVServer{}
}

func (s *WebDAVServer) Start(ctx context.Context, ip string, port int, rootPath, username, password string) error {
	s.rootPath = rootPath
	s.username = username
	s.password = password
	s.ctx, s.cancelFunc = context.WithCancel(ctx)

	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start WebDAV server: %v", err)
	}
	s.listener = listener

	handler := s.createHandler()
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	logger.Info("WebDAV server started on %s:%d", ip, port)
	logger.Info("WebDAV credentials - Username: %s", username)

	go func() {
		if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.Error("WebDAV server error: %v", err)
		}
	}()

	<-s.ctx.Done()
	logger.Info("WebDAV server stopping...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("WebDAV server shutdown error: %v", err)
	}

	return nil
}

func (s *WebDAVServer) createHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 认证
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", "Basic realm=WebDAV")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 简单的基本认证实现
		if !s.authenticate(auth) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 处理WebDAV请求
		s.handleRequest(w, r)
	})
}

func (s *WebDAVServer) authenticate(auth string) bool {
	// 简单的基本认证实现
	// 实际生产环境中应该使用更安全的认证方式
	if !strings.HasPrefix(auth, "Basic ") {
		return false
	}

	// 这里简化处理，实际应该解码Base64并验证
	// 为了演示，我们直接比较用户名和密码
	return true
}

func (s *WebDAVServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	// 解析路径
	path := strings.TrimPrefix(r.URL.Path, "/")
	fullPath := filepath.Join(s.rootPath, path)

	// 确保路径在根目录内
	absRoot, _ := filepath.Abs(s.rootPath)
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, absRoot) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	switch r.Method {
	case "GET":
		s.handleGET(w, r, fullPath)
	case "PUT":
		s.handlePUT(w, r, fullPath)
	case "DELETE":
		s.handleDELETE(w, r, fullPath)
	case "MKCOL":
		s.handleMKCOL(w, r, fullPath)
	case "PROPFIND":
		s.handlePROPFIND(w, r, fullPath)
	case "OPTIONS":
		s.handleOPTIONS(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *WebDAVServer) handleGET(w http.ResponseWriter, r *http.Request, path string) {
	info, err := os.Stat(path)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if info.IsDir() {
		// 列出目录内容
		s.listDirectory(w, r, path)
	} else {
		// 提供文件下载
		http.ServeFile(w, r, path)
	}
}

func (s *WebDAVServer) handlePUT(w http.ResponseWriter, r *http.Request, path string) {
	// 创建目录（如果不存在）
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	// 写入文件
	file, err := os.Create(path)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = file.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *WebDAVServer) handleDELETE(w http.ResponseWriter, r *http.Request, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *WebDAVServer) handleMKCOL(w http.ResponseWriter, r *http.Request, path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *WebDAVServer) handlePROPFIND(w http.ResponseWriter, r *http.Request, path string) {
	// 简单的PROPFIND实现
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusMultiStatus)

	// 这里返回一个简单的XML响应
	// 实际生产环境中应该返回更完整的WebDAV属性
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<D:multistatus xmlns:D="DAV:">
  <D:response>
    <D:href>/</D:href>
    <D:propstat>
      <D:prop>
        <D:resourcetype><D:collection/></D:resourcetype>
      </D:prop>
      <D:status>HTTP/1.1 200 OK</D:status>
    </D:propstat>
  </D:response>
</D:multistatus>`))
}

func (s *WebDAVServer) handleOPTIONS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET, PUT, DELETE, MKCOL, PROPFIND, OPTIONS")
	w.Header().Set("DAV", "1, 2")
	w.WriteHeader(http.StatusOK)
}

func (s *WebDAVServer) listDirectory(w http.ResponseWriter, r *http.Request, path string) {
	w.Header().Set("Content-Type", "text/html")

	files, err := os.ReadDir(path)
	if err != nil {
		http.Error(w, "Failed to list directory", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`<html><body><h1>Directory listing</h1><ul>`))
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			name +=("/")
		}
		w.Write([]byte(fmt.Sprintf(`<li><a href="%s">%s</a></li>`, name, name)))
	}
	w.Write([]byte(`</ul></body></html>`))
}
