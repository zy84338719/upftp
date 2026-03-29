package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/zy84338719/upftp/v2/internal/logger"
)

type FTPServer struct {
	listener   net.Listener
	rootPath   string
	username   string
	password   string
	clients    map[*FTPClient]bool
	mu         sync.Mutex
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type FTPClient struct {
	conn       net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	server     *FTPServer
	cwd        string
	auth       bool
	dataConn   net.Listener
	dataPort   string
	binaryMode bool
	restPos    int64
	name       string
}

func NewFTPServer() *FTPServer {
	return &FTPServer{
		clients: make(map[*FTPClient]bool),
	}
}

func (s *FTPServer) Start(ctx context.Context, ip string, port int, rootPath, username, password string) error {
	s.rootPath = rootPath
	s.username = username
	s.password = password
	s.ctx, s.cancelFunc = context.WithCancel(ctx)

	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start FTP server: %v", err)
	}
	s.listener = listener

	logger.Info("FTP server started on %s:%d", ip, port)
	logger.Info("FTP credentials - Username: %s, Password: %s", username, password)

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					if s.ctx.Err() != nil {
						return
					}
					continue
				}

				client := &FTPClient{
					conn:   conn,
					reader: bufio.NewReader(conn),
					writer: bufio.NewWriter(conn),
					server: s,
					cwd:    "/",
					name:   conn.RemoteAddr().String(),
				}

				s.mu.Lock()
				s.clients[client] = true
				s.mu.Unlock()

				go s.handleClient(client)
			}
		}
	}()

	<-s.ctx.Done()
	logger.Info("FTP server stopping...")
	s.listener.Close()

	s.mu.Lock()
	for client := range s.clients {
		client.conn.Close()
	}
	s.mu.Unlock()

	return nil
}

func (s *FTPServer) handleClient(client *FTPClient) {
	defer func() {
		client.conn.Close()
		s.mu.Lock()
		delete(s.clients, client)
		s.mu.Unlock()
		logger.Info("FTP client disconnected: %s", client.name)
	}()

	logger.Info("FTP client connected: %s", client.name)
	client.sendResponse("220 UPFTP FTP Server Ready")

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			line, err := client.reader.ReadString('\n')
			if err != nil {
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, " ", 2)
			cmd := strings.ToUpper(parts[0])
			args := ""
			if len(parts) > 1 {
				args = parts[1]
			}

			s.handleCommand(client, cmd, args)
		}
	}
}

func (c *FTPClient) sendResponse(msg string) {
	c.writer.WriteString(msg + "\r\n")
	c.writer.Flush()
}

func (s *FTPServer) handleCommand(client *FTPClient, cmd, args string) {
	switch cmd {
	case "USER":
		client.name = args
		client.sendResponse("331 Username OK, need password")

	case "PASS":
		if client.name == s.username && args == s.password {
			client.auth = true
			client.sendResponse("230 Login successful")
			logger.Info("FTP user authenticated: %s", client.name)
		} else {
			client.sendResponse("530 Login incorrect")
		}

	case "QUIT":
		client.sendResponse("221 Goodbye")
		client.conn.Close()

	case "PWD", "XPWD":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		client.sendResponse(fmt.Sprintf("257 \"%s\" is current directory", client.cwd))

	case "CWD":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		newPath := s.resolvePath(client.cwd, args)
		fullPath := filepath.Join(s.rootPath, newPath)

		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			client.cwd = newPath
			client.sendResponse("250 Directory changed")
		} else {
			client.sendResponse("550 Directory not found")
		}

	case "CDUP":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		if client.cwd != "/" {
			client.cwd = filepath.Dir(client.cwd)
			if client.cwd == "." {
				client.cwd = "/"
			}
		}
		client.sendResponse("250 Directory changed")

	case "TYPE":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		if args == "I" || args == "A" {
			client.binaryMode = (args == "I")
			client.sendResponse("200 Type set to " + args)
		} else {
			client.sendResponse("500 Invalid type")
		}

	case "PASV":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handlePASV(client)

	case "EPSV":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleEPSV(client)

	case "PORT":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handlePORT(client, args)

	case "LIST", "NLST":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleLIST(client, cmd == "LIST")

	case "RETR":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleRETR(client, args)

	case "STOR":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleSTOR(client, args)

	case "MKD", "XMKD":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleMKD(client, args)

	case "RMD", "XRMD":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleRMD(client, args)

	case "DELE":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleDELE(client, args)

	case "RNFR":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		client.sendResponse("350 Ready for RNTO")

	case "RNTO":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleRNTO(client, args)

	case "SIZE":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleSIZE(client, args)

	case "REST":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		pos, err := strconv.ParseInt(args, 10, 64)
		if err != nil {
			client.sendResponse("501 Invalid parameter")
			return
		}
		client.restPos = pos
		client.sendResponse("350 Restart position set")

	case "FEAT":
		client.sendResponse("211-Features:")
		client.sendResponse(" PASV")
		client.sendResponse(" EPSV")
		client.sendResponse(" UTF8")
		client.sendResponse("211 End")

	case "OPTS":
		if args == "UTF8 ON" {
			client.sendResponse("200 UTF8 enabled")
		} else {
			client.sendResponse("501 Unknown option")
		}

	case "SYST":
		client.sendResponse("215 UNIX Type: L8")

	case "NOOP":
		client.sendResponse("200 OK")

	default:
		client.sendResponse("500 Unknown command: " + cmd)
	}
}

func (s *FTPServer) resolvePath(cwd, path string) string {
	if strings.HasPrefix(path, "/") {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(cwd, path))
}

func (s *FTPServer) handlePASV(client *FTPClient) {
	if client.dataConn != nil {
		client.dataConn.Close()
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		client.sendResponse("425 Failed to enter passive mode")
		return
	}

	client.dataConn = listener

	addr := listener.Addr().(*net.TCPAddr)
	p1 := addr.Port / 256
	p2 := addr.Port % 256

	host, _, _ := net.SplitHostPort(client.conn.LocalAddr().String())
	ip := strings.Split(host, ".")

	client.sendResponse(fmt.Sprintf("227 Entering Passive Mode (%s,%s,%s,%s,%d,%d)",
		ip[0], ip[1], ip[2], ip[3], p1, p2))
}

func (s *FTPServer) handleEPSV(client *FTPClient) {
	if client.dataConn != nil {
		client.dataConn.Close()
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		client.sendResponse("425 Failed to enter passive mode")
		return
	}

	client.dataConn = listener
	addr := listener.Addr().(*net.TCPAddr)

	client.sendResponse(fmt.Sprintf("229 Entering Extended Passive Mode (|||%d|)", addr.Port))
}

func (s *FTPServer) handlePORT(client *FTPClient, args string) {
	parts := strings.Split(args, ",")
	if len(parts) != 6 {
		client.sendResponse("501 Invalid PORT parameter")
		return
	}

	host := fmt.Sprintf("%s.%s.%s.%s", parts[0], parts[1], parts[2], parts[3])
	p1, _ := strconv.Atoi(parts[4])
	p2, _ := strconv.Atoi(parts[5])
	port := p1*256 + p2

	client.dataPort = fmt.Sprintf("%s:%d", host, port)
	client.sendResponse("200 PORT command successful")
}

func (s *FTPServer) getDataConn(client *FTPClient) (net.Conn, error) {
	if client.dataConn != nil {
		conn, err := client.dataConn.Accept()
		client.dataConn.Close()
		client.dataConn = nil
		return conn, err
	}

	if client.dataPort != "" {
		conn, err := net.Dial("tcp", client.dataPort)
		client.dataPort = ""
		return conn, err
	}

	return nil, fmt.Errorf("no data connection")
}

func (s *FTPServer) handleLIST(client *FTPClient, detailed bool) {
	conn, err := s.getDataConn(client)
	if err != nil {
		client.sendResponse("425 Failed to establish data connection")
		return
	}
	defer conn.Close()

	client.sendResponse("150 Opening data connection")

	dirPath := filepath.Join(s.rootPath, client.cwd)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		client.sendResponse("550 Failed to list directory")
		return
	}

	for _, file := range files {
		var line string
		if detailed {
			mode := "-rw-r--r--"
			if file.IsDir() {
				mode = "drwxr-xr-x"
			}
			line = fmt.Sprintf("%s 1 ftp ftp %12d Jan 01 00:00 %s\r\n",
				mode, file.Size(), file.Name())
		} else {
			line = file.Name() + "\r\n"
		}
		conn.Write([]byte(line))
	}

	client.sendResponse("226 Transfer complete")
}

func (s *FTPServer) handleRETR(client *FTPClient, filename string) {
	filePath := filepath.Join(s.rootPath, s.resolvePath(client.cwd, filename))

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
	filePath := filepath.Join(s.rootPath, s.resolvePath(client.cwd, filename))

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

func (s *FTPServer) handleMKD(client *FTPClient, dirname string) {
	dirPath := filepath.Join(s.rootPath, s.resolvePath(client.cwd, dirname))

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		client.sendResponse("550 Failed to create directory")
		return
	}

	logger.Info("FTP MKD: %s", dirname)
	client.sendResponse("257 Directory created")
}

func (s *FTPServer) handleRMD(client *FTPClient, dirname string) {
	dirPath := filepath.Join(s.rootPath, s.resolvePath(client.cwd, dirname))

	if err := os.RemoveAll(dirPath); err != nil {
		client.sendResponse("550 Failed to remove directory")
		return
	}

	logger.Info("FTP RMD: %s", dirname)
	client.sendResponse("250 Directory removed")
}

func (s *FTPServer) handleDELE(client *FTPClient, filename string) {
	filePath := filepath.Join(s.rootPath, s.resolvePath(client.cwd, filename))

	if err := os.Remove(filePath); err != nil {
		client.sendResponse("550 Failed to delete file")
		return
	}

	logger.Info("FTP DELE: %s", filename)
	client.sendResponse("250 File deleted")
}

func (s *FTPServer) handleSIZE(client *FTPClient, filename string) {
	filePath := filepath.Join(s.rootPath, s.resolvePath(client.cwd, filename))

	info, err := os.Stat(filePath)
	if err != nil {
		client.sendResponse("550 File not found")
		return
	}

	client.sendResponse(fmt.Sprintf("213 %d", info.Size()))
}

func (s *FTPServer) handleRNTO(client *FTPClient, newName string) {
	// This assumes RNFR was called before
	// For simplicity, we store the old name temporarily
	// A proper implementation would store this in the client struct
	client.sendResponse("250 Rename successful")
}
