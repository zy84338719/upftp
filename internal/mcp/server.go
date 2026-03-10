package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/filehandler"
	"github.com/zy84338719/upftp/internal/server"
)

type ServerStatus struct {
	HTTPRunning bool
	HTTPPort    int
	FTPRunning  bool
	FTPPort     int
	Root        string
	IP          string
}

type MCPServer struct {
	server     *mcpserver.MCPServer
	root       string
	httpServer *server.HTTPServer
	ftpServer  *server.FTPServer
	httpCancel context.CancelFunc
	ftpCancel  context.CancelFunc
	httpPort   int
	ftpPort    int
	mu         sync.Mutex
	ctx        context.Context
}

func NewMCPServer() *MCPServer {
	s := mcpserver.NewMCPServer(
		"upftp",
		config.AppConfig.Version,
		mcpserver.WithToolCapabilities(true),
	)

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	mcpServer := &MCPServer{
		server: s,
		root:   config.AppConfig.Root,
		ctx:    ctx,
	}

	mcpServer.registerTools()
	return mcpServer
}

func (s *MCPServer) registerTools() {
	s.server.AddTool(mcp.NewTool("list_files",
		mcp.WithDescription("List files and directories in the specified path"),
		mcp.WithString("path",
			mcp.Description("Relative path from root directory (default: /)"),
		),
	), s.handleListFiles)

	s.server.AddTool(mcp.NewTool("get_file_info",
		mcp.WithDescription("Get detailed information about a file or directory"),
		mcp.WithString("path",
			mcp.Description("Relative path to the file or directory"),
			mcp.Required(),
		),
	), s.handleGetFileInfo)

	s.server.AddTool(mcp.NewTool("read_file",
		mcp.WithDescription("Read the content of a text file"),
		mcp.WithString("path",
			mcp.Description("Relative path to the text file"),
			mcp.Required(),
		),
	), s.handleReadFile)

	s.server.AddTool(mcp.NewTool("download_file",
		mcp.WithDescription("Get download URL and base64 encoded content for a file"),
		mcp.WithString("path",
			mcp.Description("Relative path to the file"),
			mcp.Required(),
		),
	), s.handleDownloadFile)

	s.server.AddTool(mcp.NewTool("search_files",
		mcp.WithDescription("Search for files matching a pattern"),
		mcp.WithString("pattern",
			mcp.Description("Search pattern (supports wildcards)"),
			mcp.Required(),
		),
		mcp.WithString("path",
			mcp.Description("Base path to search from (default: /)"),
		),
	), s.handleSearchFiles)

	s.server.AddTool(mcp.NewTool("get_directory_tree",
		mcp.WithDescription("Get the directory tree structure"),
		mcp.WithString("path",
			mcp.Description("Root path for the tree (default: /)"),
		),
	), s.handleGetDirectoryTree)

	s.server.AddTool(mcp.NewTool("start_server",
		mcp.WithDescription("Start HTTP/FTP server for LAN file sharing. Returns access URL."),
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
		mcp.WithDescription("Stop HTTP/FTP servers"),
		mcp.WithBoolean("stop_http",
			mcp.Description("Stop HTTP server (default: true)"),
		),
		mcp.WithBoolean("stop_ftp",
			mcp.Description("Stop FTP server (default: true)"),
		),
	), s.handleStopServer)

	s.server.AddTool(mcp.NewTool("get_server_status",
		mcp.WithDescription("Get current server status including running state, ports, and access URLs"),
	), s.handleGetServerStatus)

	s.server.AddTool(mcp.NewTool("set_share_directory",
		mcp.WithDescription("Change the shared directory"),
		mcp.WithString("path",
			mcp.Description("Absolute or relative path to the directory to share"),
			mcp.Required(),
		),
	), s.handleSetShareDirectory)

	s.server.AddTool(mcp.NewTool("upload_file",
		mcp.WithDescription("Upload a file to the server. Content should be base64 encoded."),
		mcp.WithString("path",
			mcp.Description("Destination path including filename (e.g., /folder/file.txt)"),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("Base64 encoded file content"),
			mcp.Required(),
		),
	), s.handleUploadFile)

	s.server.AddTool(mcp.NewTool("delete_file",
		mcp.WithDescription("Delete a file or directory"),
		mcp.WithString("path",
			mcp.Description("Path to the file or directory to delete"),
			mcp.Required(),
		),
	), s.handleDeleteFile)

	s.server.AddTool(mcp.NewTool("rename_file",
		mcp.WithDescription("Rename a file or directory"),
		mcp.WithString("path",
			mcp.Description("Current path of the file or directory"),
			mcp.Required(),
		),
		mcp.WithString("new_name",
			mcp.Description("New name for the file or directory"),
			mcp.Required(),
		),
	), s.handleRenameFile)

	s.server.AddTool(mcp.NewTool("create_directory",
		mcp.WithDescription("Create a new directory"),
		mcp.WithString("path",
			mcp.Description("Path of the directory to create"),
			mcp.Required(),
		),
	), s.handleCreateDirectory)

	s.server.AddTool(mcp.NewTool("copy_file",
		mcp.WithDescription("Copy a file to a new location"),
		mcp.WithString("source",
			mcp.Description("Source file path"),
			mcp.Required(),
		),
		mcp.WithString("destination",
			mcp.Description("Destination file path"),
			mcp.Required(),
		),
	), s.handleCopyFile)
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

	ip, err := s.getLANIP()
	if err != nil {
		ip = "127.0.0.1"
	}

	var results []string

	if s.httpServer == nil {
		s.httpServer = server.NewHTTPServer()
		s.httpPort = httpPort
		config.AppConfig.Port = ":" + strconv.Itoa(httpPort)

		ctx, cancel := context.WithCancel(s.ctx)
		s.httpCancel = cancel

		go func() {
			s.httpServer.Start(ctx, ip, httpPort, ftpPort, s.root)
		}()

		results = append(results, fmt.Sprintf("HTTP Server started: http://%s:%d", ip, httpPort))
	} else {
		results = append(results, fmt.Sprintf("HTTP Server already running on port %d", s.httpPort))
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

			results = append(results, fmt.Sprintf("FTP Server started: ftp://%s:%d (user: %s, pass: %s)",
				ip, ftpPort, config.AppConfig.Username, config.AppConfig.Password))
		} else {
			results = append(results, fmt.Sprintf("FTP Server already running on port %d", s.ftpPort))
		}
	}

	results = append(results, fmt.Sprintf("Shared directory: %s", s.root))

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
		results = append(results, "HTTP Server stopped")
	}

	if stopFTP && s.ftpCancel != nil {
		s.ftpCancel()
		s.ftpServer = nil
		s.ftpCancel = nil
		results = append(results, "FTP Server stopped")
	}

	if len(results) == 0 {
		results = append(results, "No servers to stop")
	}

	return mcp.NewToolResultText(strings.Join(results, "\n")), nil
}

func (s *MCPServer) handleGetServerStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ip, err := s.getLANIP()
	if err != nil {
		ip = "127.0.0.1"
	}

	var results []string
	results = append(results, "=== UPFTP Server Status ===")
	results = append(results, fmt.Sprintf("Shared Directory: %s", s.root))
	results = append(results, fmt.Sprintf("LAN IP: %s", ip))
	results = append(results, "")

	if s.httpServer != nil {
		results = append(results, fmt.Sprintf("HTTP Server: RUNNING on port %d", s.httpPort))
		results = append(results, fmt.Sprintf("  Access URL: http://%s:%d", ip, s.httpPort))
	} else {
		results = append(results, "HTTP Server: STOPPED")
	}

	if s.ftpServer != nil {
		results = append(results, fmt.Sprintf("FTP Server: RUNNING on port %d", s.ftpPort))
		results = append(results, fmt.Sprintf("  Access URL: ftp://%s:%d", ip, s.ftpPort))
		results = append(results, fmt.Sprintf("  Credentials: %s / %s", config.AppConfig.Username, config.AppConfig.Password))
	} else {
		results = append(results, "FTP Server: STOPPED")
	}

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

	return mcp.NewToolResultText(fmt.Sprintf("Shared directory changed to: %s\nNote: Restart servers for the change to take effect on running servers.", absPath)), nil
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

func (s *MCPServer) handleListFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relativePath := request.GetString("path", "/")

	if !s.isPathSafe(relativePath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	fullPath := filepath.Join(s.root, relativePath)

	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read directory: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Directory: %s\n\n", relativePath))
	result.WriteString("Name\t\tSize\t\tModified\t\tType\n")
	result.WriteString(strings.Repeat("-", 80) + "\n")

	for _, file := range files {
		fileType := filehandler.GetFileType(file.Name())
		fileTypeStr := "directory"
		if !file.IsDir() {
			fileTypeStr = filehandler.GetFileTypeString(fileType)
		}

		size := "-"
		if !file.IsDir() {
			size = filehandler.FormatFileSize(file.Size())
		}

		result.WriteString(fmt.Sprintf("%s\t\t%s\t\t%s\t\t%s\n",
			file.Name(),
			size,
			file.ModTime().Format("2006-01-02 15:04"),
			fileTypeStr,
		))
	}

	return mcp.NewToolResultText(result.String()), nil
}

func (s *MCPServer) handleGetFileInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relativePath, err := request.RequireString("path")
	if err != nil || relativePath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	if !s.isPathSafe(relativePath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	fullPath := filepath.Join(s.root, relativePath)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get file info: %v", err)), nil
	}

	fileType := filehandler.GetFileType(fileInfo.Name())
	fileTypeStr := "directory"
	if !fileInfo.IsDir() {
		fileTypeStr = filehandler.GetFileTypeString(fileType)
	}

	result := fmt.Sprintf(`File Information:
Path: %s
Name: %s
Type: %s
Size: %s
Modified: %s
Is Directory: %t
`,
		relativePath,
		fileInfo.Name(),
		fileTypeStr,
		filehandler.FormatFileSize(fileInfo.Size()),
		fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		fileInfo.IsDir(),
	)

	return mcp.NewToolResultText(result), nil
}

func (s *MCPServer) handleReadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relativePath, err := request.RequireString("path")
	if err != nil || relativePath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	if !s.isPathSafe(relativePath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	fullPath := filepath.Join(s.root, relativePath)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to access file: %v", err)), nil
	}

	if fileInfo.IsDir() {
		return mcp.NewToolResultError("Cannot read a directory"), nil
	}

	if fileInfo.Size() > 10*1024*1024 {
		return mcp.NewToolResultError("File too large (max 10MB)"), nil
	}

	fileType := filehandler.GetFileType(fileInfo.Name())
	if fileType != filehandler.FileTypeText && fileType != filehandler.FileTypeCode {
		return mcp.NewToolResultError("Can only read text files"), nil
	}

	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	return mcp.NewToolResultText(string(content)), nil
}

func (s *MCPServer) handleDownloadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relativePath, err := request.RequireString("path")
	if err != nil || relativePath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	if !s.isPathSafe(relativePath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	fullPath := filepath.Join(s.root, relativePath)

	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to access file: %v", err)), nil
	}

	if fileInfo.IsDir() {
		return mcp.NewToolResultError("Cannot download a directory directly. Use list_files instead."), nil
	}

	if fileInfo.Size() > 10*1024*1024 {
		return mcp.NewToolResultError("File too large for direct download (max 10MB)"), nil
	}

	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	fileType := filehandler.GetFileType(fileInfo.Name())
	mimeType := filehandler.GetMimeType(fileType)

	result := fmt.Sprintf(`File: %s
Size: %s
MIME Type: %s

Base64 Content:
%s
`,
		relativePath,
		filehandler.FormatFileSize(fileInfo.Size()),
		mimeType,
		base64.StdEncoding.EncodeToString(content),
	)

	return mcp.NewToolResultText(result), nil
}

func (s *MCPServer) handleSearchFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pattern, err := request.RequireString("pattern")
	if err != nil || pattern == "" {
		return mcp.NewToolResultError("Pattern parameter is required"), nil
	}

	basePath := request.GetString("path", "/")

	if !s.isPathSafe(basePath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	searchPath := filepath.Join(s.root, basePath)
	var results []string

	err = filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return nil
		}

		if matched {
			relPath, _ := filepath.Rel(s.root, path)
			fileType := "file"
			if info.IsDir() {
				fileType = "directory"
			}
			results = append(results, fmt.Sprintf("/%s (%s, %s)",
				relPath,
				fileType,
				filehandler.FormatFileSize(info.Size()),
			))
		}

		return nil
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No files found matching pattern: " + pattern), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Found %d files:\n%s",
		len(results),
		strings.Join(results, "\n"),
	)), nil
}

func (s *MCPServer) handleGetDirectoryTree(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	basePath := request.GetString("path", "/")

	if !s.isPathSafe(basePath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	searchPath := filepath.Join(s.root, basePath)
	var result strings.Builder

	result.WriteString(fmt.Sprintf("Directory Tree: %s\n", basePath))

	err := s.buildTree(searchPath, "", &result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build tree: %v", err)), nil
	}

	return mcp.NewToolResultText(result.String()), nil
}

func (s *MCPServer) buildTree(currentPath string, prefix string, result *strings.Builder) error {
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		return err
	}

	for i, file := range files {
		if i == len(files)-1 {
			result.WriteString(prefix + "└── " + file.Name() + "\n")
			if file.IsDir() {
				s.buildTree(filepath.Join(currentPath, file.Name()), prefix+"    ", result)
			}
		} else {
			result.WriteString(prefix + "├── " + file.Name() + "\n")
			if file.IsDir() {
				s.buildTree(filepath.Join(currentPath, file.Name()), prefix+"│   ", result)
			}
		}
	}

	return nil
}

func (s *MCPServer) isPathSafe(relativePath string) bool {
	cleanPath := filepath.Clean(relativePath)
	absPath := filepath.Join(s.root, cleanPath)

	if !strings.HasPrefix(absPath, s.root) {
		return false
	}

	return !strings.Contains(cleanPath, "..")
}

func (s *MCPServer) handleUploadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	destPath, err := request.RequireString("path")
	if err != nil || destPath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	contentBase64, err := request.RequireString("content")
	if err != nil || contentBase64 == "" {
		return mcp.NewToolResultError("Content parameter is required"), nil
	}

	if !s.isPathSafe(destPath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	content, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid base64 content: %v", err)), nil
	}

	fullPath := filepath.Join(s.root, destPath)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create directory: %v", err)), nil
	}

	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("File uploaded successfully: %s (%d bytes)", destPath, len(content))), nil
}

func (s *MCPServer) handleDeleteFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil || path == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	if !s.isPathSafe(path) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	fullPath := filepath.Join(s.root, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return mcp.NewToolResultError("File or directory does not exist"), nil
	}

	if err := os.RemoveAll(fullPath); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to delete: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted: %s", path)), nil
}

func (s *MCPServer) handleRenameFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	oldPath, err := request.RequireString("path")
	if err != nil || oldPath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	newName, err := request.RequireString("new_name")
	if err != nil || newName == "" {
		return mcp.NewToolResultError("New name parameter is required"), nil
	}

	if !s.isPathSafe(oldPath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	fullOldPath := filepath.Join(s.root, oldPath)
	fullNewPath := filepath.Join(filepath.Dir(fullOldPath), newName)

	if err := os.Rename(fullOldPath, fullNewPath); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to rename: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Renamed: %s -> %s", oldPath, newName)), nil
}

func (s *MCPServer) handleCreateDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil || path == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	if !s.isPathSafe(path) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	fullPath := filepath.Join(s.root, path)

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create directory: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Directory created: %s", path)), nil
}

func (s *MCPServer) handleCopyFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	source, err := request.RequireString("source")
	if err != nil || source == "" {
		return mcp.NewToolResultError("Source parameter is required"), nil
	}

	destination, err := request.RequireString("destination")
	if err != nil || destination == "" {
		return mcp.NewToolResultError("Destination parameter is required"), nil
	}

	if !s.isPathSafe(source) || !s.isPathSafe(destination) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	srcPath := filepath.Join(s.root, source)
	dstPath := filepath.Join(s.root, destination)

	content, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read source file: %v", err)), nil
	}

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create destination directory: %v", err)), nil
	}

	if err := ioutil.WriteFile(dstPath, content, 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write destination file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Copied: %s -> %s", source, destination)), nil
}

func (s *MCPServer) Start(ctx context.Context) error {
	s.ctx = ctx
	return mcpserver.ServeStdio(s.server)
}

func (s *MCPServer) GetServer() *mcpserver.MCPServer {
	return s.server
}
