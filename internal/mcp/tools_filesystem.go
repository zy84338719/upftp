package mcp

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/zy84338719/upftp/internal/filehandler"
)

func (s *MCPServer) registerFilesystemTools() {
	s.server.AddTool(mcp.NewTool("list_files",
		mcp.WithDescription("List files and directories in the specified path. Returns a formatted table with name, size, modified date, and type."),
		mcp.WithString("path",
			mcp.Description("Relative path from root directory (default: /)"),
		),
	), s.handleListFiles)

	s.server.AddTool(mcp.NewTool("get_file_info",
		mcp.WithDescription("Get detailed information about a file or directory including size, type, modification time, and permissions."),
		mcp.WithString("path",
			mcp.Description("Relative path to the file or directory"),
			mcp.Required(),
		),
	), s.handleGetFileInfo)

	s.server.AddTool(mcp.NewTool("read_file",
		mcp.WithDescription("Read the content of a text file. Only works for text and code files (max 10MB)."),
		mcp.WithString("path",
			mcp.Description("Relative path to the text file"),
			mcp.Required(),
		),
	), s.handleReadFile)

	s.server.AddTool(mcp.NewTool("write_file",
		mcp.WithDescription("Write content to a text file. Creates the file if it doesn't exist, overwrites if it does."),
		mcp.WithString("path",
			mcp.Description("Destination path including filename"),
			mcp.Required(),
		),
		mcp.WithString("content",
			mcp.Description("Text content to write to the file"),
			mcp.Required(),
		),
	), s.handleWriteFile)

	s.server.AddTool(mcp.NewTool("download_file",
		mcp.WithDescription("Get download URL and base64 encoded content for a file (max 10MB)."),
		mcp.WithString("path",
			mcp.Description("Relative path to the file"),
			mcp.Required(),
		),
	), s.handleDownloadFile)

	s.server.AddTool(mcp.NewTool("search_files",
		mcp.WithDescription("Search for files matching a pattern. Supports wildcards like *.txt or file*.pdf"),
		mcp.WithString("pattern",
			mcp.Description("Search pattern (supports wildcards: *, ?)"),
			mcp.Required(),
		),
		mcp.WithString("path",
			mcp.Description("Base path to search from (default: /)"),
		),
	), s.handleSearchFiles)

	s.server.AddTool(mcp.NewTool("get_directory_tree",
		mcp.WithDescription("Get the directory tree structure in a visual format."),
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
	result.WriteString(fmt.Sprintf("%-40s %-15s %-20s %-15s\n", "Name", "Size", "Modified", "Type"))
	result.WriteString(strings.Repeat("-", 95) + "\n")

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

		name := file.Name()
		if len(name) > 37 {
			name = name[:34] + "..."
		}

		result.WriteString(fmt.Sprintf("%-40s %-15s %-20s %-15s\n",
			name,
			size,
			file.ModTime().Format("2006-01-02 15:04"),
			fileTypeStr,
		))
	}

	result.WriteString(fmt.Sprintf("\nTotal: %d items", len(files)))

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

	result := fmt.Sprintf(`=== File Information ===
Path: %s
Name: %s
Type: %s
Size: %s
Modified: %s
Is Directory: %t
Permissions: %s
`,
		relativePath,
		fileInfo.Name(),
		fileTypeStr,
		filehandler.FormatFileSize(fileInfo.Size()),
		fileInfo.ModTime().Format("2006-01-02 15:04:05"),
		fileInfo.IsDir(),
		fileInfo.Mode().String(),
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
		return mcp.NewToolResultError("Can only read text files. Use download_file for binary files."), nil
	}

	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	return mcp.NewToolResultText(string(content)), nil
}

func (s *MCPServer) handleWriteFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	destPath, err := request.RequireString("path")
	if err != nil || destPath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	content, err := request.RequireString("content")
	if err != nil {
		return mcp.NewToolResultError("Content parameter is required"), nil
	}

	if !s.isPathSafe(destPath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	fullPath := filepath.Join(s.root, destPath)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create directory: %v", err)), nil
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("File written successfully: %s (%d bytes)", destPath, len(content))), nil
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

	result := fmt.Sprintf(`=== File Download ===
Path: %s
Size: %s
MIME Type: %s

Base64 Content:
%s
`,
		relativePath,
		filehandler.FormatFileSize(fileInfo.Size()),
		mimeType,
		encodeBase64(content),
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
		return mcp.NewToolResultText(fmt.Sprintf("No files found matching pattern: %s", pattern)), nil
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("=== Search Results ===\n"))
	output.WriteString(fmt.Sprintf("Pattern: %s\n", pattern))
	output.WriteString(fmt.Sprintf("Found: %d items\n\n", len(results)))
	output.WriteString(strings.Join(results, "\n"))

	return mcp.NewToolResultText(output.String()), nil
}

func (s *MCPServer) handleGetDirectoryTree(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	basePath := request.GetString("path", "/")

	if !s.isPathSafe(basePath) {
		return mcp.NewToolResultError("Access denied: invalid path"), nil
	}

	searchPath := filepath.Join(s.root, basePath)
	var result strings.Builder

	result.WriteString(fmt.Sprintf("=== Directory Tree ===\n"))
	result.WriteString(fmt.Sprintf("Root: %s\n\n", basePath))

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
		isLast := i == len(files)-1

		connector := "├── "
		if isLast {
			connector = "└── "
		}

		icon := "📄"
		if file.IsDir() {
			icon = "📁"
		}

		result.WriteString(prefix + connector + icon + " " + file.Name() + "\n")

		if file.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			s.buildTree(filepath.Join(currentPath, file.Name()), newPrefix, result)
		}
	}

	return nil
}
