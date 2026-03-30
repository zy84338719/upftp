package model

import "time"

type FileInfo struct {
	Name       string
	Size       int64
	ModTime    time.Time
	IsDir      bool
	Path       string
	FileType   FileType
	CanPreview bool
	MimeType   string
	Icon       string
}

type LegacyFileInfo struct {
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

func NewLegacyFileInfo(fi FileInfo) LegacyFileInfo {
	fileTypeStr := "directory"
	if !fi.IsDir {
		fileTypeStr = GetFileTypeString(fi.FileType)
	}
	return LegacyFileInfo{
		Name:        fi.Name,
		Size:        FormatFileSize(fi.Size),
		ModTime:     fi.ModTime.Format("2006-01-02 15:04:05"),
		IsDir:       fi.IsDir,
		CanPreview:  fi.CanPreview,
		FileType:    fi.FileType,
		FileTypeStr: fileTypeStr,
		Path:        fi.Path,
		Icon:        fi.Icon,
		MimeType:    fi.MimeType,
	}
}
