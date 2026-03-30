package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/server"
)

func (s *MCPServer) registerServerTools() {
	s.server.AddTool(mcp.NewTool("start_server",
		mcp.WithDescription("Start HTTP/FTP server for LAN file sharing. Returns access URLs that can be shared with others."),
		mcp.WithNumber("http_port",
			mcp.Description("HTTP server port (default: 10000)"),
		),
		mcp.WithNumber("ftp_port",
			mcp.Description("FTP server port (default: 2121)"),
		),
		mcp.WithBoolean("enable_ftp",
			mcp.Description("Enable FTP server (default: false)"),
		),
		mcp.WithString("directory",
			mcp.Description("Directory to share (default: current shared directory)"),
		),
	), s.handleStartServer)

	s.server.AddTool(mcp.NewTool("stop_server",
		mcp.WithDescription("Stop HTTP/FTP servers."),
		mcp.WithBoolean("stop_http",
			mcp.Description("Stop HTTP server (default: true)"),
		),
		mcp.WithBoolean("stop_ftp",
			mcp.Description("Stop FTP server (default: true)"),
		),
	), s.handleStopServer)

	s.server.AddTool(mcp.NewTool("get_server_status",
		mcp.WithDescription("Get current server status including running state, ports, and access URLs."),
	), s.handleGetServerStatus)

	s.server.AddTool(mcp.NewTool("set_share_directory",
		mcp.WithDescription("Change the shared directory for the MCP server."),
		mcp.WithString("path",
			mcp.Description("Absolute or relative path to the directory to share"),
			mcp.Required(),
		),
	), s.handleSetShareDirectory)

	s.server.AddTool(mcp.NewTool("get_download_url",
		mcp.WithDescription("Get a direct download URL for a file that can be shared with others."),
		mcp.WithString("path",
			mcp.Description("Relative path to the file"),
			mcp.Required(),
		),
	), s.handleGetDownloadURL)
}

func (s *MCPServer) handleStartServer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	httpPort := request.GetInt("http_port", 10000)
	ftpPort := request.GetInt("ftp_port", 2121)
	enableFTP := request.GetBool("enable_ftp", false)
	newDir := request.GetString("directory", "")

	if newDir != "" {
		absPath, err := filepath.Abs(newDir)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid directory path: %v", err)), nil
		}
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return mcp.NewToolResultError(fmt.Sprintf("Directory does not exist: %s", absPath)), nil
		}
		s.root = absPath
		config.AppConfig.Root = absPath
	}

	ip := s.resolveIP()
	s.ip = ip

	var results []string
	results = append(results, fmt.Sprintf("Shared directory: %s", s.root))
	results = append(results, "")

	if s.httpServer == nil {
		s.httpServer = server.NewHTTPServer()
		s.httpPort = httpPort
		config.AppConfig.Port = ":" + strconv.Itoa(httpPort)

		ctx, cancel := context.WithCancel(s.ctx)
		s.httpCancel = cancel

		go func() {
			s.httpServer.Start(ctx, ip, httpPort, ftpPort, s.root)
		}()

		time.Sleep(100 * time.Millisecond)
		results = append(results, fmt.Sprintf("✅ HTTP Server started: http://%s:%d", ip, httpPort))
	} else {
		results = append(results, fmt.Sprintf("ℹ️  HTTP Server already running on port %d", s.httpPort))
	}

	if enableFTP {
		if s.ftpServer == nil {
			s.ftpServer = server.NewFTPServer()
			s.ftpPort = ftpPort

			ctx, cancel := context.WithCancel(s.ctx)
			s.ftpCancel = cancel

			go func() {
				s.ftpServer.Start(ctx, ip, ftpPort, s.root, config.AppConfig.Username, config.AppConfig.Password)
			}()

			time.Sleep(100 * time.Millisecond)
			results = append(results, fmt.Sprintf("✅ FTP Server started: ftp://%s:%d", ip, ftpPort))
			results = append(results, fmt.Sprintf("   Credentials: %s / %s", config.AppConfig.Username, config.AppConfig.Password))
		} else {
			results = append(results, fmt.Sprintf("ℹ️  FTP Server already running on port %d", s.ftpPort))
		}
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func (s *MCPServer) handleStopServer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stopHTTP := request.GetBool("stop_http", true)
	stopFTP := request.GetBool("stop_ftp", true)

	var results []string

	if stopHTTP && s.httpCancel != nil {
		s.httpCancel()
		s.httpServer = nil
		s.httpCancel = nil
		results = append(results, "✅ HTTP Server stopped")
	}

	if stopFTP && s.ftpCancel != nil {
		s.ftpCancel()
		s.ftpServer = nil
		s.ftpCancel = nil
		results = append(results, "✅ FTP Server stopped")
	}

	if len(results) == 0 {
		results = append(results, "ℹ️  No servers to stop")
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func (s *MCPServer) handleGetServerStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ip := s.resolveIP()

	var results []string
	results = append(results, "=== UPFTP Server Status ===")
	results = append(results, fmt.Sprintf("Shared Directory: %s", s.root))
	results = append(results, fmt.Sprintf("LAN IP: %s", ip))
	results = append(results, "")

	if s.httpServer != nil {
		results = append(results, fmt.Sprintf("HTTP Server: ✅ RUNNING on port %d", s.httpPort))
		results = append(results, fmt.Sprintf("  Access URL: http://%s:%d", ip, s.httpPort))
	} else {
		results = append(results, "HTTP Server: ⏹️  STOPPED")
	}

	if s.ftpServer != nil {
		results = append(results, fmt.Sprintf("FTP Server: ✅ RUNNING on port %d", s.ftpPort))
		results = append(results, fmt.Sprintf("  Access URL: ftp://%s:%d", ip, s.ftpPort))
		results = append(results, fmt.Sprintf("  Credentials: %s / %s", config.AppConfig.Username, config.AppConfig.Password))
	} else {
		results = append(results, "FTP Server: ⏹️  STOPPED")
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func (s *MCPServer) handleGetDownloadURL(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil || path == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	if !s.isPathSafe(path) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	s.mu.Lock()
	ip := s.resolveIP()
	port := s.resolveHTTPPort()
	s.mu.Unlock()

	url := fmt.Sprintf("http://%s:%d/download/%s", ip, port, path)

	var results []string
	results = append(results, "=== Download URL ===")
	results = append(results, url)
	results = append(results, "")
	results = append(results, "Share this URL with others to download the file.")
	results = append(results, "Note: The HTTP server must be running for this URL to work.")

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func (s *MCPServer) handleSetShareDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	newPath, err := request.RequireString("path")
	if err != nil || newPath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	absPath, err := filepath.Abs(newPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid path: %v", err)), nil
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return mcp.NewToolResultError(fmt.Sprintf("Directory does not exist: %s", absPath)), nil
	}

	s.mu.Lock()
	s.root = absPath
	config.AppConfig.Root = absPath
	s.mu.Unlock()

	return mcp.NewToolResultText(fmt.Sprintf("Shared directory changed to: %s\n\nNote: Restart servers for the change to take effect on running servers.", absPath)), nil
}
