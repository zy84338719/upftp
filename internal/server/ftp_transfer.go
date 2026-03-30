package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zy84338719/upftp/internal/logger"
)

func ftpJoin(root, relativePath string) string {
	return filepath.Join(root, filepath.FromSlash(relativePath))
}

func statPath(p string) (os.FileInfo, error) {
	return os.Stat(p)
}

func (s *FTPServer) handleLIST(client *FTPClient, detailed bool) {
	conn, err := s.getDataConn(client)
	if err != nil {
		client.sendResponse("425 Failed to establish data connection")
		return
	}
	defer conn.Close()

	client.sendResponse("150 Opening data connection")

	dirPath := ftpJoin(s.rootPath, client.cwd)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		client.sendResponse("550 Failed to list directory")
		return
	}

	for _, file := range files {
		var line string
		if detailed {
			mode := "-rw-r--r--"
			nlink := 1
			if file.IsDir() {
				mode = "drwxr-xr-x"
				nlink = 2
			}
			line = fmt.Sprintf("%s %d ftp ftp %12d %s %s\r\n",
				mode, nlink, file.Size(), formatFTPDate(file.ModTime()), file.Name())
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

	dirPath := ftpJoin(s.rootPath, client.cwd)
	files, err := ioutil.ReadDir(dirPath)
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
		if file.IsDir() {
			writeEntry("dir", file.Name(), 0, file.ModTime())
		} else {
			writeEntry("file", file.Name(), file.Size(), file.ModTime())
		}
	}

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleRETR(client *FTPClient, filename string) {
	filePath := ftpJoin(s.rootPath, s.resolvePath(client.cwd, filename))

	file, err := os.Open(filePath)
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

	io.Copy(conn, file)
	logger.Info("FTP RETR: %s", filename)

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleSTOR(client *FTPClient, filename string) {
	filePath := ftpJoin(s.rootPath, s.resolvePath(client.cwd, filename))

	file, err := os.Create(filePath)
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

	io.Copy(file, conn)
	logger.Info("FTP STOR: %s", filename)

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleAPPE(client *FTPClient, filename string) {
	filePath := ftpJoin(s.rootPath, s.resolvePath(client.cwd, filename))

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	io.Copy(file, conn)
	logger.Info("FTP APPE: %s", filename)

	client.sendResponse("226 Transfer complete")
}

var _ = strings.ReplaceAll
