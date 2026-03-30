package ftp

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/zy84338719/upftp/internal/logger"
)

func (s *FTPServer) handleLIST(client *FTPClient, detailed bool) {
	conn, err := s.getDataConn(client)
	if err != nil {
		client.sendResponse("425 Failed to establish data connection")
		return
	}
	defer conn.Close()

	client.sendResponse("150 Opening data connection")

	files, err := s.svc.ListDir(client.cwd)
	if err != nil {
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
		conn.Write([]byte(line))
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

	client.sendResponse("150 Opening data connection")

	files, err := s.svc.ListDir(client.cwd)
	if err != nil {
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
		conn.Write([]byte(entry))
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
	file, err := s.svc.OpenFile(path)
	if err != nil {
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

	client.sendResponse("150 Opening data connection")

	if client.restPos > 0 {
		file.Seek(client.restPos, 0)
		client.restPos = 0
	}

	if _, copyErr := io.Copy(conn, file); copyErr != nil {
		logger.Error("FTP RETR error: %v", copyErr)
	}
	logger.Info("FTP RETR: %s", filename)

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleSTOR(client *FTPClient, filename string) {
	path := s.resolvePath(client.cwd, filename)
	file, err := s.svc.CreateFileForWrite(path)
	if err != nil {
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

	client.sendResponse("150 Opening data connection")

	if _, copyErr := io.Copy(file, conn); copyErr != nil {
		logger.Error("FTP STOR error: %v", copyErr)
	}
	logger.Info("FTP STOR: %s", filename)

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleAPPE(client *FTPClient, filename string) {
	path := s.resolvePath(client.cwd, filename)

	fullPath, err := s.svc.SafePath(path)
	if err != nil {
		client.sendResponse("550 Permission denied")
		return
	}

	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
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

	client.sendResponse("150 Opening data connection")

	if _, copyErr := io.Copy(file, conn); copyErr != nil {
		logger.Error("FTP APPE error: %v", copyErr)
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
