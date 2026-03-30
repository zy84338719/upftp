package nfs

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/zy84338719/upftp/pkg/logger"
)

type NFSServer struct {
	listener   net.Listener
	rootPath   string
	username   string
	password   string
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewNFSServer() *NFSServer {
	return &NFSServer{}
}

func (s *NFSServer) Start(ctx context.Context, ip string, port int, rootPath, username, password string) error {
	s.rootPath = rootPath
	s.username = username
	s.password = password
	s.ctx, s.cancelFunc = context.WithCancel(ctx)

	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start NFS server: %v", err)
	}
	s.listener = listener

	logger.Info("NFS server started on %s:%d", ip, port)
	logger.Info("NFS credentials - Username: %s", username)

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					if s.ctx.Err() != nil {
						return
					}
					logger.Error("Failed to accept NFS connection: %v", err)
					continue
				}

				go s.handleConnection(conn)
			}
		}
	}()

	<-s.ctx.Done()
	logger.Info("NFS server stopping...")
	s.listener.Close()

	return nil
}

func (s *NFSServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	logger.Info("NFS client connected: %s", conn.RemoteAddr().String())

	// NFS协议实现比较复杂，这里只是一个简单的框架
	// 实际生产环境中应该使用完整的NFS协议实现

	// 读取客户端请求
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		logger.Error("Failed to read NFS request: %v", err)
		return
	}

	// 简单的响应
	response := []byte("NFS server is running")
	_, err = conn.Write(response)
	if err != nil {
		logger.Error("Failed to write NFS response: %v", err)
		return
	}

	logger.Info("NFS client disconnected: %s", conn.RemoteAddr().String())
}

// 辅助函数：解析NFS路径
func (s *NFSServer) resolvePath(path string) (string, error) {
	// 确保路径在根目录内
	fullPath := filepath.Join(s.rootPath, strings.TrimPrefix(path, "/"))

	// 确保路径在根目录内
	absRoot, _ := filepath.Abs(s.rootPath)
	absPath, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absPath, absRoot) {
		return "", fmt.Errorf("access denied")
	}

	return fullPath, nil
}
