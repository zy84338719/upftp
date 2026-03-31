package model

import (
	"testing"
)

func TestFileType(t *testing.T) {
	// Test GetFileType function
	testCases := []struct {
		name     string
		filename string
		expected FileType
	}{
		{"text file", "test.txt", FileTypeText},
		{"image file", "test.jpg", FileTypeImage},
		{"video file", "test.mp4", FileTypeVideo},
		{"audio file", "test.mp3", FileTypeAudio},
		{"archive file", "test.zip", FileTypeArchive},
		{"other file", "test.xyz", FileTypeOther},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetFileType(tc.filename)
			if result != tc.expected {
				t.Errorf("GetFileType(%s) = %v, want %v", tc.filename, result, tc.expected)
			}
		})
	}
}
