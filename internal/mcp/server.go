package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/zy84338719/upftp/internal/config"
	"github.com/zy84338719/upftp/internal/filehandler"
)

type MCPServer struct {
	server *server.MCPServer
	root   string
}

func NewMCPServer() *MCPServer {
	s := server.NewMCPServer(
		"upftp",
		config.AppConfig.Version,
		server.WithToolCapabilities(true),
	)

	mcpServer := &MCPServer{
		server: s,
		root:   config.AppConfig.Root,
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

func (s *MCPServer) Start(ctx context.Context) error {
	return server.ServeStdio(s.server)
}

func (s *MCPServer) GetServer() *server.MCPServer {
	return s.server
}
