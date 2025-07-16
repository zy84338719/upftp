package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        string
	FTPPort     string
	Root        string
	AutoSelect  bool
	EnableFTP   bool
	Username    string
	Password    string
	Version     string
	LastCommit  string
}

var AppConfig *Config

func Init(version, lastCommit string) {
	AppConfig = &Config{
		Version:    version,
		LastCommit: lastCommit,
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "upftp - A lightweight file sharing server with FTP support\n\n")
		fmt.Fprintf(os.Stderr, "Project: https://github.com/zy84338719/upftp\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  upftp [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -p <port>       HTTP server port (default: 10000)\n")
		fmt.Fprintf(os.Stderr, "  -ftp <port>     FTP server port (default: 2121)\n")
		fmt.Fprintf(os.Stderr, "  -d <dir>        Share directory (default: current directory)\n")
		fmt.Fprintf(os.Stderr, "  -auto           Automatically select first available network interface\n")
		fmt.Fprintf(os.Stderr, "  -enable-ftp     Enable FTP server (default: false)\n")
		fmt.Fprintf(os.Stderr, "  -user <name>    FTP username (default: admin)\n")
		fmt.Fprintf(os.Stderr, "  -pass <pass>    FTP password (default: admin)\n")
		fmt.Fprintf(os.Stderr, "  -h              Show this help message\n")
	}

	p := flag.String("p", "10000", "HTTP server port")
	ftpPort := flag.String("ftp", "2121", "FTP server port")
	dir := flag.String("d", "./", "Share directory")
	autoIP := flag.Bool("auto", false, "Automatically select first available network interface")
	enableFTP := flag.Bool("enable-ftp", false, "Enable FTP server")
	user := flag.String("user", "admin", "FTP username")
	pass := flag.String("pass", "admin", "FTP password")

	flag.Parse()

	AppConfig.Port = ":" + *p
	AppConfig.FTPPort = ":" + *ftpPort
	AppConfig.Root = *dir
	AppConfig.AutoSelect = *autoIP
	AppConfig.EnableFTP = *enableFTP
	AppConfig.Username = *user
	AppConfig.Password = *pass
}

func (c *Config) GetHTTPPort() int {
	port := c.Port[1:] // 去掉冒号
	if p, err := strconv.Atoi(port); err == nil {
		return p
	}
	return 10000
}

func (c *Config) GetFTPPort() int {
	port := c.FTPPort[1:] // 去掉冒号
	if p, err := strconv.Atoi(port); err == nil {
		return p
	}
	return 2121
}
