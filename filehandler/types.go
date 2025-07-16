package filehandler

import (
	"fmt"
	"path/filepath"
	"strings"
)

type FileType int

const (
	FileTypeUnknown FileType = iota
	FileTypeImage
	FileTypeText
	FileTypeVideo
	FileTypeAudio
	FileTypePDF
	FileTypeDocument
	FileTypeArchive
	FileTypeCode
)

type FileInfo struct {
	Name        string
	Size        string
	ModTime     string
	IsDir       bool
	CanPreview  bool
	FileType    FileType
	FileTypeStr string
	Path        string
	Icon        string
	MimeType    string
}

var fileTypeMap = map[string]FileType{
	// Images
	".jpg":  FileTypeImage,
	".jpeg": FileTypeImage,
	".png":  FileTypeImage,
	".gif":  FileTypeImage,
	".bmp":  FileTypeImage,
	".webp": FileTypeImage,
	".svg":  FileTypeImage,
	".ico":  FileTypeImage,

	// Videos
	".mp4":  FileTypeVideo,
	".avi":  FileTypeVideo,
	".mov":  FileTypeVideo,
	".wmv":  FileTypeVideo,
	".flv":  FileTypeVideo,
	".webm": FileTypeVideo,
	".mkv":  FileTypeVideo,
	".m4v":  FileTypeVideo,

	// Audio
	".mp3":  FileTypeAudio,
	".wav":  FileTypeAudio,
	".flac": FileTypeAudio,
	".aac":  FileTypeAudio,
	".ogg":  FileTypeAudio,
	".wma":  FileTypeAudio,
	".m4a":  FileTypeAudio,

	// Text files
	".txt":  FileTypeText,
	".md":   FileTypeText,
	".json": FileTypeText,
	".xml":  FileTypeText,
	".yaml": FileTypeText,
	".yml":  FileTypeText,
	".csv":  FileTypeText,
	".log":  FileTypeText,
	".ini":  FileTypeText,
	".conf": FileTypeText,
	".cfg":  FileTypeText,

	// Code files
	".go":   FileTypeCode,
	".js":   FileTypeCode,
	".ts":   FileTypeCode,
	".html": FileTypeCode,
	".css":  FileTypeCode,
	".py":   FileTypeCode,
	".java": FileTypeCode,
	".cpp":  FileTypeCode,
	".c":    FileTypeCode,
	".php":  FileTypeCode,
	".rb":   FileTypeCode,
	".rs":   FileTypeCode,
	".sh":   FileTypeCode,
	".sql":  FileTypeCode,

	// Documents
	".pdf":  FileTypePDF,
	".doc":  FileTypeDocument,
	".docx": FileTypeDocument,
	".xls":  FileTypeDocument,
	".xlsx": FileTypeDocument,
	".ppt":  FileTypeDocument,
	".pptx": FileTypeDocument,
	".rtf":  FileTypeDocument,
	".odt":  FileTypeDocument,
	".ods":  FileTypeDocument,
	".odp":  FileTypeDocument,

	// Archives
	".zip": FileTypeArchive,
	".rar": FileTypeArchive,
	".7z":  FileTypeArchive,
	".tar": FileTypeArchive,
	".gz":  FileTypeArchive,
	".bz2": FileTypeArchive,
	".xz":  FileTypeArchive,
}

var iconMap = map[FileType]string{
	FileTypeImage:    "🖼️",
	FileTypeVideo:    "🎥",
	FileTypeAudio:    "🎵",
	FileTypeText:     "📝",
	FileTypeCode:     "💻",
	FileTypePDF:      "📄",
	FileTypeDocument: "📊",
	FileTypeArchive:  "📦",
	FileTypeUnknown:  "📄",
}

var mimeTypeMap = map[FileType]string{
	FileTypeImage:    "image/*",
	FileTypeVideo:    "video/*",
	FileTypeAudio:    "audio/*",
	FileTypeText:     "text/plain",
	FileTypeCode:     "text/plain",
	FileTypePDF:      "application/pdf",
	FileTypeDocument: "application/octet-stream",
	FileTypeArchive:  "application/zip",
	FileTypeUnknown:  "application/octet-stream",
}

func GetFileType(filename string) FileType {
	ext := strings.ToLower(filepath.Ext(filename))
	if fileType, ok := fileTypeMap[ext]; ok {
		return fileType
	}
	return FileTypeUnknown
}

func GetFileTypeString(fileType FileType) string {
	switch fileType {
	case FileTypeImage:
		return "image"
	case FileTypeVideo:
		return "video"
	case FileTypeAudio:
		return "audio"
	case FileTypeText:
		return "text"
	case FileTypeCode:
		return "code"
	case FileTypePDF:
		return "pdf"
	case FileTypeDocument:
		return "document"
	case FileTypeArchive:
		return "archive"
	default:
		return "unknown"
	}
}

func GetFileIcon(fileType FileType) string {
	if icon, ok := iconMap[fileType]; ok {
		return icon
	}
	return iconMap[FileTypeUnknown]
}

func GetMimeType(fileType FileType) string {
	if mimeType, ok := mimeTypeMap[fileType]; ok {
		return mimeType
	}
	return mimeTypeMap[FileTypeUnknown]
}

func CanPreviewFile(fileType FileType) bool {
	switch fileType {
	case FileTypeImage, FileTypeText, FileTypeCode, FileTypeVideo, FileTypeAudio:
		return true
	default:
		return false
	}
}

func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func IsPathSafe(path string) bool {
	cleanPath := filepath.Clean(path)
	return !strings.Contains(cleanPath, "..") && !strings.HasPrefix(cleanPath, "/")
}
