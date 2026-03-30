package mcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/logger"
	"github.com/zy84338719/upftp/internal/server"
)

type MCPServer struct {
	server     *mcpserver.MCPServer
	root       string
	httpServer *server.HTTPServer
	ftpServer  *server.FTPServer
	httpCancel context.CancelFunc
	ftpCancel  context.CancelFunc
	httpPort   int
	ftpPort    int
	ip         string
	mu         sync.Mutex
	ctx        context.Context
	renamePath string
}

func NewMCPServer() *MCPServer {
	s := mcpserver.NewMCPServer(
		"upftp",
		config.AppConfig.Version,
		mcpserver.WithToolCapabilities(true),
	)

	ctx := context.Background()

	mcpServer := &MCPServer{
		server: s,
		root:   config.AppConfig.Root,
		ctx:    ctx,
	}

	mcpServer.registerTools()
	return mcpServer
}

func (s *MCPServer) registerTools() {
	s.registerFilesystemTools()
	s.registerFileOpsTools()
	s.registerServerTools()
}

func (s *MCPServer) Start(ctx context.Context) error {
	s.ctx = ctx
	logger.Info("MCP server starting...")
	return mcpserver.ServeStdio(s.server)
}

func (s *MCPServer) GetServer() *mcpserver.MCPServer {
	return s.server
}

func (s *MCPServer) getLANIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no LAN IP found")
}

func (s *MCPServer) isPathSafe(relativePath string) bool {
	cleanPath := filepath.Clean(relativePath)
	absPath := filepath.Join(s.root, cleanPath)

	if !strings.HasPrefix(absPath, s.root) {
		return false
	}

	return !strings.Contains(cleanPath, "..")
}

func (s *MCPServer) resolveIP() string {
	if s.ip != "" {
		return s.ip
	}
	ip, err := s.getLANIP()
	if err != nil {
		return "127.0.0.1"
	}
	return ip
}

func (s *MCPServer) resolveHTTPPort() int {
	if s.httpPort != 0 {
		return s.httpPort
	}
	return 10000
}

var _ = os.DevNull
