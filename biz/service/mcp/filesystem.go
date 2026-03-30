package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/zy84338719/upftp/pkg/file/model"
)

func (s *MCPServer) registerFilesystemTools() {
	s.server.AddTool(mcp.NewTool("list_files",
		mcp.WithDescription("List files and directories in the specified path."),
		mcp.WithString("path", mcp.Description("Relative path from root directory (default: /)")),
	), s.handleListFiles)

	s.server.AddTool(mcp.NewTool("get_file_info",
		mcp.WithDescription("Get detailed information about a file or directory."),
		mcp.WithString("path", mcp.Description("Relative path"), mcp.Required()),
	), s.handleGetFileInfo)

	s.server.AddTool(mcp.NewTool("read_file",
		mcp.WithDescription("Read the content of a text file (max 10MB)."),
		mcp.WithString("path", mcp.Description("Relative path to the text file"), mcp.Required()),
	), s.handleReadFile)

	s.server.AddTool(mcp.NewTool("write_file",
		mcp.WithDescription("Write content to a text file."),
		mcp.WithString("path", mcp.Description("Destination path"), mcp.Required()),
		mcp.WithString("content", mcp.Description("Text content"), mcp.Required()),
	), s.handleWriteFile)

	s.server.AddTool(mcp.NewTool("download_file",
		mcp.WithDescription("Get download URL and base64 content for a file (max 10MB)."),
		mcp.WithString("path", mcp.Description("Relative path"), mcp.Required()),
	), s.handleDownloadFile)

	s.server.AddTool(mcp.NewTool("search_files",
		mcp.WithDescription("Search for files matching a pattern."),
		mcp.WithString("pattern", mcp.Description("Search pattern"), mcp.Required()),
		mcp.WithString("path", mcp.Description("Base path (default: /)")),
	), s.handleSearchFiles)

	s.server.AddTool(mcp.NewTool("get_directory_tree",
		mcp.WithDescription("Get directory tree structure."),
		mcp.WithString("path", mcp.Description("Root path (default: /)")),
	), s.handleGetDirectoryTree)
}

func (s *MCPServer) handleListFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relativePath := request.GetString("path", "/")

	files, err := s.svc.ListFiles(relativePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read directory: %v", err)), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Directory: %s\n\n", relativePath))
	result.WriteString(fmt.Sprintf("%-40s %-15s %-20s %-15s\n", "Name", "Size", "Modified", "Type"))
	result.WriteString(strings.Repeat("-", 95) + "\n")

	for _, f := range files {
		size := "-"
		fileTypeStr := "directory"
		if !f.IsDir {
			size = model.FormatFileSize(f.Size)
			fileTypeStr = model.GetFileTypeString(f.FileType)
		}

		name := f.Name
		if len(name) > 37 {
			name = name[:34] + "..."
		}

		result.WriteString(fmt.Sprintf("%-40s %-15s %-20s %-15s\n",
			name, size, f.ModTime.Format("2006-01-02 15:04"), fileTypeStr))
	}

	result.WriteString(fmt.Sprintf("\nTotal: %d items", len(files)))
	return mcp.NewToolResultText(result.String()), nil
}

func (s *MCPServer) handleGetFileInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relativePath, err := request.RequireString("path")
	if err != nil || relativePath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	info, err := s.svc.Stat(relativePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get file info: %v", err)), nil
	}

	fileType := model.GetFileType(relativePath)
	fileTypeStr := "directory"
	if !info.IsDir() {
		fileTypeStr = model.GetFileTypeString(fileType)
	}

	result := fmt.Sprintf(`Path: %s
Name: %s
Type: %s
Size: %s
Modified: %s
Is Directory: %t
Permissions: %s
`, relativePath, info.Name(), fileTypeStr, model.FormatFileSize(info.Size()),
		info.ModTime().Format("2006-01-02 15:04:05"), info.IsDir(), info.Mode().String())

	return mcp.NewToolResultText(result), nil
}

func (s *MCPServer) handleReadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relativePath, err := request.RequireString("path")
	if err != nil || relativePath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	content, err := s.svc.ReadFileContent(relativePath)
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

	if err := s.svc.WriteFileContent(destPath, []byte(content)); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("File written successfully: %s (%d bytes)", destPath, len(content))), nil
}

func (s *MCPServer) handleDownloadFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relativePath, err := request.RequireString("path")
	if err != nil || relativePath == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	content, err := s.svc.ReadFileContent(relativePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	info, err := s.svc.Stat(relativePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to stat file: %v", err)), nil
	}

	fileType := model.GetFileType(relativePath)
	mimeType := model.GetMimeType(fileType)

	result := fmt.Sprintf("Path: %s\nSize: %s\nMIME Type: %s\n\nBase64 Content:\n%s\n",
		relativePath, model.FormatFileSize(info.Size()), mimeType, base64.StdEncoding.EncodeToString(content))

	return mcp.NewToolResultText(result), nil
}

func (s *MCPServer) handleSearchFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pattern, err := request.RequireString("pattern")
	if err != nil || pattern == "" {
		return mcp.NewToolResultError("Pattern parameter is required"), nil
	}
	basePath := request.GetString("path", "/")

	results, err := s.svc.SearchFiles(pattern, basePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No files found matching pattern: %s", pattern)), nil
	}

	var output strings.Builder
	output.WriteString("=== Search Results ===\n")
	output.WriteString(fmt.Sprintf("Pattern: %s\n", pattern))
	output.WriteString(fmt.Sprintf("Found: %d items\n\n", len(results)))
	output.WriteString(strings.Join(results, "\n"))

	return mcp.NewToolResultText(output.String()), nil
}

func (s *MCPServer) handleGetDirectoryTree(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	basePath := request.GetString("path", "/")

	tree, err := s.svc.BuildTree(basePath, 0)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build tree: %v", err)), nil
	}

	if tree == nil {
		return mcp.NewToolResultText("Empty directory tree"), nil
	}

	var result strings.Builder
	result.WriteString("=== Directory Tree ===\n")
	result.WriteString(fmt.Sprintf("Root: %s\n\n", basePath))
	renderMcpTree(result, tree, 0)

	return mcp.NewToolResultText(result.String()), nil
}

func renderMcpTree(b strings.Builder, node *model.TreeNode, depth int) {
	prefix := ""
	for i := 0; i < depth; i++ {
		prefix += "|   "
	}
	connector := "|-- "
	icon := "file"
	if node.IsDir {
		icon = "dir"
	}
	b.WriteString(fmt.Sprintf("%s%s%s %s\n", prefix, connector, icon, node.Name))
	for _, child := range node.Children {
		renderMcpTree(b, child, depth+1)
	}
}
