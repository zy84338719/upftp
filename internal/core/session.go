package core

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zy84338719/upftp/internal/ftp"
)

// SessionSpec 描述创建一个 FTP 会话所需的参数。
type SessionSpec struct {
	Dir       string `json:"dir"`
	Port      int    `json:"port"` // 0 = 自动选可用端口
	Anonymous bool   `json:"anonymous"`
	User      string `json:"user"`
	Pass      string `json:"pass"`
	ReadOnly  bool   `json:"read_only"`
}

// Session 是一个运行中的 FTP 会话。
type Session struct {
	ID        string    `json:"id"`
	Dir       string    `json:"dir"`
	Port      int       `json:"port"`
	Anonymous bool      `json:"anonymous"`
	User      string    `json:"user,omitempty"`
	Pass      string    `json:"pass,omitempty"`
	ReadOnly  bool      `json:"read_only"`
	StartedAt time.Time `json:"started_at"`
	Host      string    `json:"host"` // 对外展示 IP

	FTPServer *ftp.FTPServer
	cancel    context.CancelFunc
}

// FTPURL 拼接对外可用的 ftp:// 访问地址。
func (s *Session) FTPURL() string {
	if s.Host == "" {
		return fmt.Sprintf("ftp://0.0.0.0:%d", s.Port)
	}
	if s.Anonymous {
		return fmt.Sprintf("ftp://%s:%d", s.Host, s.Port)
	}
	return fmt.Sprintf("ftp://%s:%s@%s:%d", s.User, s.Pass, s.Host, s.Port)
}

// SessionManager 管理一组 FTP 会话的生命周期,线程安全。
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	counter  atomic.Int64
	logger   *slog.Logger
	host     string // 对外展示 IP
}

// NewSessionManager 创建会话管理器。host 用于拼接展示地址。
func NewSessionManager(host string, logger *slog.Logger) *SessionManager {
	if logger == nil {
		logger = slog.Default()
	}
	return &SessionManager{
		sessions: make(map[string]*Session),
		host:     host,
		logger:   logger,
	}
}

// Create 启动一个新的 FTP 会话并返回它。Port 为 0 时自动分配。
func (m *SessionManager) Create(ctx context.Context, spec SessionSpec) (*Session, error) {
	port := spec.Port
	if port == 0 {
		var err error
		port, err = freePort()
		if err != nil {
			return nil, fmt.Errorf("no free port: %w", err)
		}
	}

	srv := ftp.New(ftp.Options{
		Root:      spec.Dir,
		Port:      port,
		Anonymous: spec.Anonymous,
		User:      spec.User,
		Pass:      spec.Pass,
		ReadOnly:  spec.ReadOnly,
		Logger:    m.logger,
	})

	subCtx, cancel := context.WithCancel(ctx)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start(subCtx)
	}()

	// 等待一小段,确认监听未立即失败(端口占用等)。
	select {
	case err := <-errCh:
		cancel()
		return nil, fmt.Errorf("ftp start failed: %w", err)
	case <-time.After(150 * time.Millisecond):
		// 正常:仍在运行
	case <-subCtx.Done():
		cancel()
		return nil, fmt.Errorf("cancelled before start")
	}

	id := fmt.Sprintf("s%d", m.counter.Add(1))
	sess := &Session{
		ID:        id,
		Dir:       spec.Dir,
		Port:      port,
		Anonymous: spec.Anonymous,
		User:      spec.User,
		Pass:      spec.Pass,
		ReadOnly:  spec.ReadOnly,
		StartedAt: time.Now(),
		Host:      m.host,
		FTPServer: srv,
		cancel:    cancel,
	}

	m.mu.Lock()
	m.sessions[id] = sess
	m.mu.Unlock()

	m.logger.Info("session created", "id", id, "dir", spec.Dir, "port", port, "url", sess.FTPURL())
	return sess, nil
}

// Get 按 ID 取会话。
func (m *SessionManager) Get(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	return s, ok
}

// List 返回所有会话。
func (m *SessionManager) List() []*Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		out = append(out, s)
	}
	return out
}

// Stop 停止并移除一个会话。
func (m *SessionManager) Stop(id string) bool {
	m.mu.Lock()
	sess, ok := m.sessions[id]
	if !ok {
		m.mu.Unlock()
		return false
	}
	delete(m.sessions, id)
	m.mu.Unlock()

	sess.cancel()
	m.logger.Info("session stopped", "id", id)
	return true
}

// StopAll 停止全部会话(用于进程退出)。
func (m *SessionManager) StopAll() {
	m.mu.Lock()
	ids := make([]string, 0, len(m.sessions))
	for id, s := range m.sessions {
		ids = append(ids, id)
		s.cancel()
	}
	m.sessions = make(map[string]*Session)
	m.mu.Unlock()
	m.logger.Info("stopped all sessions", "count", len(ids))
}

// freePort 让操作系统分配一个空闲 TCP 端口(通过 listen-then-close 取端口,再交由 FTP 重用)。
// 注:存在极小的端口复用竞争窗口,对内网分享场景可接受。
func freePort() (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	return port, nil
}
