package server

import (
	"context"
	"log"
)

type FTPServer struct{}

func NewFTPServer() *FTPServer {
	return &FTPServer{}
}

func (s *FTPServer) Start(ctx context.Context, ip string, port int, rootPath, username, password string) error {
	log.Printf("FTP server would start on %s:%d", ip, port)
	log.Printf("FTP credentials - Username: %s, Password: %s", username, password)
	log.Printf("FTP root path: %s", rootPath)

	go func() {
		<-ctx.Done()
		log.Println("FTP server context cancelled")
	}()

	log.Println("FTP server functionality is under development")
	return nil
}
