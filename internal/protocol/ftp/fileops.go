package ftp

import (
	"fmt"

	"github.com/zy84338719/upftp/internal/logger"
)

func (s *FTPServer) handleMKD(client *FTPClient, dirname string) {
	path := s.resolvePath(client.cwd, dirname)
	if err := s.svc.CreateFolder(path); err != nil {
		client.sendResponse("550 Failed to create directory")
		return
	}

	logger.Info("FTP MKD: %s", dirname)
	client.sendResponse(fmt.Sprintf("257 \"%s\" directory created", dirname))
}

func (s *FTPServer) handleRMD(client *FTPClient, dirname string) {
	path := s.resolvePath(client.cwd, dirname)
	if err := s.svc.Delete(path); err != nil {
		client.sendResponse("550 Failed to remove directory")
		return
	}

	logger.Info("FTP RMD: %s", dirname)
	client.sendResponse("250 Directory removed")
}

func (s *FTPServer) handleDELE(client *FTPClient, filename string) {
	path := s.resolvePath(client.cwd, filename)
	if err := s.svc.Delete(path); err != nil {
		client.sendResponse("550 Failed to delete file")
		return
	}

	logger.Info("FTP DELE: %s", filename)
	client.sendResponse("250 File deleted")
}

func (s *FTPServer) handleSIZE(client *FTPClient, filename string) {
	path := s.resolvePath(client.cwd, filename)
	info, err := s.svc.Stat(path)
	if err != nil {
		client.sendResponse("550 File not found")
		return
	}

	client.sendResponse(fmt.Sprintf("213 %d", info.Size()))
}

func (s *FTPServer) handleMDTM(client *FTPClient, filename string) {
	path := s.resolvePath(client.cwd, filename)
	info, err := s.svc.Stat(path)
	if err != nil {
		client.sendResponse("550 File not found")
		return
	}

	client.sendResponse(fmt.Sprintf("213 %s", info.ModTime().Format("20060102150405")))
}

func (s *FTPServer) handleRNTO(client *FTPClient, newName string) {
	if client.rnfrName == "" {
		client.sendResponse("503 Need RNFR before RNTO")
		return
	}

	if err := s.svc.Rename(client.rnfrName, newName); err != nil {
		client.sendResponse("550 Failed to rename")
		client.rnfrName = ""
		return
	}

	logger.Info("FTP RNTO: %s -> %s", client.rnfrName, newName)
	client.rnfrName = ""
	client.sendResponse("250 Rename successful")
}
