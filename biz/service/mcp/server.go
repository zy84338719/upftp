package mcp

import (
	"context"
	"fmt"
	"net"
	"sync"

	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/zy84338719/upftp/biz/service/file"
	"github.com/zy84338719/upftp/pkg/conf"
	"github.com/zy84338719/upftp/pkg/logger"
)

type MCPServer struct {
	server     *mcpserver.MCPServer
	svc        *file.Service
	httpServer *httpAdapter
	ftpServer  *ftpAdapter
	httpCancel context.CancelFunc
	ftpCancel  context.CancelFunc
	httpPort   int
	ftpPort    int
	ip         string
	mu         sync.Mutex
	ctx        context.Context
	renamePath string
}

func NewMCPServer(svc *file.Service) *MCPServer {
	s := mcpserver.NewMCPServer(
		"upftp",
		conf.AppConfig.Version,
		mcpserver.WithToolCapabilities(true),
	)

	mcpServer := &MCPServer{
		server: s,
		svc:    svc,
		ctx:    context.Background(),
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
