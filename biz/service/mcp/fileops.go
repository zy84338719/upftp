package mcp

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func (s *MCPServer) registerFileOpsTools() {
	s.server.AddTool(mcp.NewTool("upload_file",
		mcp.WithDescription("Upload a file to the server with base64 encoded content."),
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
		mcp.WithDescription("Delete a file or directory. Directories are deleted recursively."),
		mcp.WithString("path",
			mcp.Description("Path to the file or directory to delete"),
			mcp.Required(),
		),
	), s.handleDeleteFile)

	s.server.AddTool(mcp.NewTool("rename_file",
		mcp.WithDescription("Rename a file or directory."),
		mcp.WithString("path",
			mcp.Description("Current path of the file or directory"),
			mcp.Required(),
		),
		mcp.WithString("new_name",
			mcp.Description("New name for the file or directory"),
			mcp.Required(),
		),
	), s.handleRenameFile)

	s.server.AddTool(mcp.NewTool("move_file",
		mcp.WithDescription("Move a file or directory to a new location."),
		mcp.WithString("source",
			mcp.Description("Source path of the file or directory"),
			mcp.Required(),
		),
		mcp.WithString("destination",
			mcp.Description("Destination path"),
			mcp.Required(),
		),
	), s.handleMoveFile)

	s.server.AddTool(mcp.NewTool("copy_file",
		mcp.WithDescription("Copy a file to a new location."),
		mcp.WithString("source",
			mcp.Description("Source file path"),
			mcp.Required(),
		),
		mcp.WithString("destination",
			mcp.Description("Destination file path"),
			mcp.Required(),
		),
	), s.handleCopyFile)

	s.server.AddTool(mcp.NewTool("create_directory",
		mcp.WithDescription("Create a new directory (and parent directories if needed)."),
		mcp.WithString("path",
			mcp.Description("Path of the directory to create"),
			mcp.Required(),
		),
	), s.handleCreateDirectory)
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

	content, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid base64 content: %v", err)), nil
	}

	if err := s.svc.WriteFileContent(destPath, content); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("File uploaded successfully: %s (%d bytes)", destPath, len(content))), nil
}

func (s *MCPServer) handleDeleteFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil || path == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	if err := s.svc.Delete(path); err != nil {
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

	if err := s.svc.Rename(oldPath, newName); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to rename: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Renamed: %s -> %s", oldPath, newName)), nil
}

func (s *MCPServer) handleMoveFile(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	source, err := request.RequireString("source")
	if err != nil || source == "" {
		return mcp.NewToolResultError("Source parameter is required"), nil
	}

	destination, err := request.RequireString("destination")
	if err != nil || destination == "" {
		return mcp.NewToolResultError("Destination parameter is required"), nil
	}

	if err := s.svc.Move(source, destination); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to move: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Moved: %s -> %s", source, destination)), nil
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

	if err := s.svc.Copy(source, destination); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to copy file: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Copied: %s -> %s", source, destination)), nil
}

func (s *MCPServer) handleCreateDirectory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := request.RequireString("path")
	if err != nil || path == "" {
		return mcp.NewToolResultError("Path parameter is required"), nil
	}

	if err := s.svc.CreateFolder(path); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create directory: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Directory created: %s", path)), nil
}
