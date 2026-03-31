package nfs

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

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

				logger.Info("NFS client connected: %s", conn.RemoteAddr().String())

				// 处理NFS连接
				go s.handleConnection(conn)
			}
		}
	}()

	<-s.ctx.Done()
	logger.Info("NFS server stopping...")
	s.listener.Close()

	return nil
}

// 处理NFS连接
func (s *NFSServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	// 读取客户端请求
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		logger.Error("Failed to read NFS request: %v", err)
		return
	}

	// 简单的NFS协议实现
	// 模拟NFS挂载响应
	response := []byte("NFS3 OK\x00")
	_, err = conn.Write(response)
	if err != nil {
		logger.Error("Failed to write NFS response: %v", err)
		return
	}

	// 继续处理其他请求
	for {
		// 读取请求
		n, err := conn.Read(buffer)
		if err != nil {
			logger.Error("Failed to read NFS request: %v", err)
			break
		}

		// 解析请求
		request := string(buffer[:n])
		logger.Debug("Received NFS request: %s", request)

		// 处理不同类型的请求
		response := s.handleNFSRequest(request)

		// 发送响应
		_, err = conn.Write([]byte(response + "\x00"))
		if err != nil {
			logger.Error("Failed to write NFS response: %v", err)
			break
		}

		// 模拟延迟
		time.Sleep(100 * time.Millisecond)
	}

	logger.Info("NFS client disconnected: %s", conn.RemoteAddr().String())
}

// 处理NFS请求
func (s *NFSServer) handleNFSRequest(request string) string {
	// 简单的请求处理
	parts := strings.Fields(request)
	if len(parts) == 0 {
		return "NFS3 ERR"
	}

	command := parts[0]
	switch command {
	case "STAT":
		if len(parts) < 2 {
			return "NFS3 ERR"
		}
		return s.handleStatRequest(parts[1])
	case "READDIR":
		if len(parts) < 2 {
			return "NFS3 ERR"
		}
		return s.handleReaddirRequest(parts[1])
	case "READ":
		if len(parts) < 3 {
			return "NFS3 ERR"
		}
		return s.handleReadRequest(parts[1], parts[2])
	case "WRITE":
		if len(parts) < 3 {
			return "NFS3 ERR"
		}
		return s.handleWriteRequest(parts[1], parts[2])
	case "MKDIR":
		if len(parts) < 2 {
			return "NFS3 ERR"
		}
		return s.handleMkdirRequest(parts[1])
	case "RMDIR":
		if len(parts) < 2 {
			return "NFS3 ERR"
		}
		return s.handleRmdirRequest(parts[1])
	case "REMOVE":
		if len(parts) < 2 {
			return "NFS3 ERR"
		}
		return s.handleRemoveRequest(parts[1])
	case "RENAME":
		if len(parts) < 3 {
			return "NFS3 ERR"
		}
		return s.handleRenameRequest(parts[1], parts[2])
	default:
		return "NFS3 OK"
	}
}

// 处理STAT请求
func (s *NFSServer) handleStatRequest(path string) string {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return "NFS3 ERR"
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return "NFS3 ERR"
	}

	return fmt.Sprintf("NFS3 OK %s %d %s", info.Name(), info.Size(), info.Mode())
}

// 处理READDIR请求
func (s *NFSServer) handleReaddirRequest(path string) string {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return "NFS3 ERR"
	}

	dirents, err := os.ReadDir(fullPath)
	if err != nil {
		return "NFS3 ERR"
	}

	response := "NFS3 OK"
	for _, de := range dirents {
		info, err := de.Info()
		if err != nil {
			continue
		}
		response += fmt.Sprintf(" %s:%d:%s", de.Name(), info.Size(), info.Mode())
	}

	return response
}

// 处理READ请求
func (s *NFSServer) handleReadRequest(path string, offset string) string {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return "NFS3 ERR"
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "NFS3 ERR"
	}

	return fmt.Sprintf("NFS3 OK %s", string(data))
}

// 处理WRITE请求
func (s *NFSServer) handleWriteRequest(path string, data string) string {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return "NFS3 ERR"
	}

	err = os.WriteFile(fullPath, []byte(data), 0644)
	if err != nil {
		return "NFS3 ERR"
	}

	return "NFS3 OK"
}

// 处理MKDIR请求
func (s *NFSServer) handleMkdirRequest(path string) string {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return "NFS3 ERR"
	}

	err = os.Mkdir(fullPath, 0755)
	if err != nil {
		return "NFS3 ERR"
	}

	return "NFS3 OK"
}

// 处理RMDIR请求
func (s *NFSServer) handleRmdirRequest(path string) string {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return "NFS3 ERR"
	}

	err = os.RemoveAll(fullPath)
	if err != nil {
		return "NFS3 ERR"
	}

	return "NFS3 OK"
}

// 处理REMOVE请求
func (s *NFSServer) handleRemoveRequest(path string) string {
	fullPath, err := s.resolvePath(path)
	if err != nil {
		return "NFS3 ERR"
	}

	err = os.Remove(fullPath)
	if err != nil {
		return "NFS3 ERR"
	}

	return "NFS3 OK"
}

// 处理RENAME请求
func (s *NFSServer) handleRenameRequest(oldPath string, newPath string) string {
	oldFullPath, err := s.resolvePath(oldPath)
	if err != nil {
		return "NFS3 ERR"
	}

	newFullPath, err := s.resolvePath(newPath)
	if err != nil {
		return "NFS3 ERR"
	}

	err = os.Rename(oldFullPath, newFullPath)
	if err != nil {
		return "NFS3 ERR"
	}

	return "NFS3 OK"
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

// 辅助函数：获取文件信息
func (s *NFSServer) getFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// 辅助函数：读取目录
func (s *NFSServer) readDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

// 辅助函数：读取文件
func (s *NFSServer) readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// 辅助函数：写入文件
func (s *NFSServer) writeFile(path string, data []byte) (int, error) {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

// 辅助函数：创建文件
func (s *NFSServer) createFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// 辅助函数：删除文件
func (s *NFSServer) removeFile(path string) error {
	return os.Remove(path)
}

// 辅助函数：重命名文件
func (s *NFSServer) renameFile(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// 辅助函数：创建目录
func (s *NFSServer) mkdir(path string) error {
	return os.Mkdir(path, 0755)
}

// 辅助函数：删除目录
func (s *NFSServer) rmdir(path string) error {
	return os.RemoveAll(path)
}

// 辅助函数：创建符号链接
func (s *NFSServer) symlink(target, linkPath string) error {
	return os.Symlink(target, linkPath)
}

// 辅助函数：读取符号链接
func (s *NFSServer) readlink(path string) (string, error) {
	return os.Readlink(path)
}

// 辅助函数：创建硬链接
func (s *NFSServer) link(oldPath, newPath string) error {
	return os.Link(oldPath, newPath)
}
