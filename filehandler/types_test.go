package filehandler

import (
	"testing"
)

func TestGetFileType(t *testing.T) {
	tests := []struct {
		filename string
		expected FileType
	}{
		{"image.jpg", FileTypeImage},
		{"video.mp4", FileTypeVideo},
		{"audio.mp3", FileTypeAudio},
		{"text.txt", FileTypeText},
		{"code.go", FileTypeCode},
		{"document.pdf", FileTypePDF},
		{"unknown.xyz", FileTypeUnknown},
	}

	for _, test := range tests {
		result := GetFileType(test.filename)
		if result != test.expected {
			t.Errorf("GetFileType(%s) = %v, expected %v", test.filename, result, test.expected)
		}
	}
}

func TestCanPreviewFile(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected bool
	}{
		{FileTypeImage, true},
		{FileTypeVideo, true},
		{FileTypeAudio, true},
		{FileTypeText, true},
		{FileTypeCode, true},
		{FileTypePDF, false},
		{FileTypeDocument, false},
		{FileTypeUnknown, false},
	}

	for _, test := range tests {
		result := CanPreviewFile(test.fileType)
		if result != test.expected {
			t.Errorf("CanPreviewFile(%v) = %v, expected %v", test.fileType, result, test.expected)
		}
	}
}

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, test := range tests {
		result := FormatFileSize(test.size)
		if result != test.expected {
			t.Errorf("FormatFileSize(%d) = %s, expected %s", test.size, result, test.expected)
		}
	}
}

func TestIsPathSafe(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"normal/path", true},
		{"../dangerous", false},
		{"/absolute/path", false},
		{"safe.txt", true},
		{"folder/../file", false},
	}

	for _, test := range tests {
		result := IsPathSafe(test.path)
		if result != test.expected {
			t.Errorf("IsPathSafe(%s) = %v, expected %v", test.path, result, test.expected)
		}
	}
}
