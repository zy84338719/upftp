package ftp

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/zy84338719/upftp/pkg/logger"
)

func (s *FTPServer) handleLIST(client *FTPClient, detailed bool) {
	conn, err := s.getDataConn(client)
	if err != nil {
		client.sendResponse("425 Failed to establish data connection")
		return
	}
	defer conn.Close()

	// 设置数据连接超时
	conn.SetDeadline(time.Now().Add(dataConnTimeout))

	client.sendResponse("150 Opening data connection")

	fullPath := s.getFullPath(client.cwd)
	files, err := os.ReadDir(fullPath)
	if err != nil {
		logger.Error("Failed to list directory: %v", err)
		client.sendResponse("550 Failed to list directory")
		return
	}

	for _, file := range files {
		var line string
		if detailed {
			info, infoErr := file.Info()
			if infoErr != nil {
				continue
			}
			mode := "-rw-r--r--"
			nlink := 1
			if file.IsDir() {
				mode = "drwxr-xr-x"
				nlink = 2
			}
			line = fmt.Sprintf("%s %d ftp ftp %12d %s %s\r\n",
				mode, nlink, info.Size(), formatFTPDate(info.ModTime()), file.Name())
		} else {
			line = file.Name() + "\r\n"
		}
		// 写入数据
		n, err := conn.Write([]byte(line))
		if err != nil {
			logger.Error("Failed to write LIST data: %v", err)
			client.sendResponse("426 Transfer aborted")
			return
		}
		if n != len(line) {
			logger.Error("Incomplete write for LIST data")
			client.sendResponse("426 Transfer aborted")
			return
		}
	}

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleMLSD(client *FTPClient) {
	conn, err := s.getDataConn(client)
	if err != nil {
		client.sendResponse("425 Failed to establish data connection")
		return
	}
	defer conn.Close()

	// 设置数据连接超时
	conn.SetDeadline(time.Now().Add(dataConnTimeout))

	client.sendResponse("150 Opening data connection")

	fullPath := s.getFullPath(client.cwd)
	files, err := os.ReadDir(fullPath)
	if err != nil {
		logger.Error("Failed to list directory: %v", err)
		client.sendResponse("550 Failed to list directory")
		return
	}

	writeEntry := func(ftype string, name string, size int64, modTime time.Time) {
		perm := "el"
		if ftype == "dir" {
			perm = "elcmp"
		} else {
			perm = "adfr"
		}
		entry := fmt.Sprintf("type=%s;size=%d;modify=%s;perm=%s; %s\r\n",
			ftype, size, modTime.Format("20060102150405"), perm, name)
		n, err := conn.Write([]byte(entry))
		if err != nil {
			logger.Error("Failed to write MLSD data: %v", err)
			client.sendResponse("426 Transfer aborted")
			return
		}
		if n != len(entry) {
			logger.Error("Incomplete write for MLSD data")
			client.sendResponse("426 Transfer aborted")
			return
		}
	}

	writeEntry("dir", ".", 0, time.Now())
	writeEntry("dir", "..", 0, time.Now())

	for _, file := range files {
		info, infoErr := file.Info()
		if infoErr != nil {
			continue
		}
		if file.IsDir() {
			writeEntry("dir", file.Name(), 0, info.ModTime())
		} else {
			writeEntry("file", file.Name(), info.Size(), info.ModTime())
		}
	}

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleRETR(client *FTPClient, filename string) {
	path := s.resolvePath(client.cwd, filename)
	fullPath := s.getFullPath(path)
	file, err := os.Open(fullPath)
	if err != nil {
		logger.Error("Failed to open file for RETR: %v", err)
		client.sendResponse("550 File not found")
		return
	}
	defer file.Close()

	conn, err := s.getDataConn(client)
	if err != nil {
		client.sendResponse("425 Failed to establish data connection")
		return
	}
	defer conn.Close()

	// 设置数据连接超时
	conn.SetDeadline(time.Now().Add(dataConnTimeout))

	client.sendResponse("150 Opening data connection")

	if client.restPos > 0 {
		_, err := file.Seek(client.restPos, 0)
		if err != nil {
			logger.Error("Failed to seek file: %v", err)
			client.sendResponse("550 Failed to seek file")
			return
		}
		client.restPos = 0
	}

	// 使用缓冲区提高传输性能
	buffer := make([]byte, bufferSize)
	_, copyErr := io.CopyBuffer(conn, file, buffer)
	if copyErr != nil {
		if netErr, ok := copyErr.(net.Error); ok && netErr.Timeout() {
			logger.Info("FTP RETR timeout: %s", filename)
		} else {
			logger.Error("FTP RETR error: %v", copyErr)
		}
		client.sendResponse("426 Transfer aborted")
		return
	}
	logger.Info("FTP RETR: %s", filename)

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleSTOR(client *FTPClient, filename string) {
	path := s.resolvePath(client.cwd, filename)
	fullPath := s.getFullPath(path)

	// 创建目录（如果不存在）
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Failed to create directory: %v", err)
		client.sendResponse("550 Failed to create directory")
		return
	}

	file, err := os.Create(fullPath)
	if err != nil {
		logger.Error("Failed to create file: %v", err)
		client.sendResponse("550 Failed to create file")
		return
	}
	defer file.Close()

	conn, err := s.getDataConn(client)
	if err != nil {
		client.sendResponse("425 Failed to establish data connection")
		return
	}
	defer conn.Close()

	// 设置数据连接超时
	conn.SetDeadline(time.Now().Add(dataConnTimeout))

	client.sendResponse("150 Opening data connection")

	// 使用缓冲区提高传输性能
	buffer := make([]byte, bufferSize)
	_, copyErr := io.CopyBuffer(file, conn, buffer)
	if copyErr != nil {
		if netErr, ok := copyErr.(net.Error); ok && netErr.Timeout() {
			logger.Info("FTP STOR timeout: %s", filename)
		} else {
			logger.Error("FTP STOR error: %v", copyErr)
		}
		client.sendResponse("426 Transfer aborted")
		return
	}
	logger.Info("FTP STOR: %s", filename)

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleAPPE(client *FTPClient, filename string) {
	path := s.resolvePath(client.cwd, filename)
	fullPath := s.getFullPath(path)

	// 创建目录（如果不存在）
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Failed to create directory: %v", err)
		client.sendResponse("550 Failed to create directory")
		return
	}

	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Error("Failed to open file: %v", err)
		client.sendResponse("550 Failed to open file")
		return
	}
	defer file.Close()

	conn, err := s.getDataConn(client)
	if err != nil {
		client.sendResponse("425 Failed to establish data connection")
		return
	}
	defer conn.Close()

	// 设置数据连接超时
	conn.SetDeadline(time.Now().Add(dataConnTimeout))

	client.sendResponse("150 Opening data connection")

	// 使用缓冲区提高传输性能
	buffer := make([]byte, bufferSize)
	_, copyErr := io.CopyBuffer(file, conn, buffer)
	if copyErr != nil {
		if netErr, ok := copyErr.(net.Error); ok && netErr.Timeout() {
			logger.Info("FTP APPE timeout: %s", filename)
		} else {
			logger.Error("FTP APPE error: %v", copyErr)
		}
		client.sendResponse("426 Transfer aborted")
		return
	}
	logger.Info("FTP APPE: %s", filename)

	client.sendResponse("226 Transfer complete")
}

func formatFTPDate(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() {
		return t.Format("Jan 02 15:04")
	}
	return t.Format("Jan 02  2006")
}
