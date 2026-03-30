package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/zy84338719/upftp/internal/logger"
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
	rnfrName   string
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
		}

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

func (c *FTPClient) sendResponse(msg string) {
	c.writer.WriteString(msg + "\r\n")
	c.writer.Flush()
}

func (s *FTPServer) resolvePath(cwd, pathArg string) string {
	if pathArg == "" {
		return cwd
	}
	if strings.HasPrefix(pathArg, "/") {
		return strings.ReplaceAll(filepathClean(pathArg), `\`, "/")
	}
	return strings.ReplaceAll(filepathClean(strings.Join([]string{cwd, pathArg}, "/")), `\`, "/")
}

func filepathClean(p string) string {
	cleaned := p
	for strings.Contains(cleaned, "//") {
		cleaned = strings.ReplaceAll(cleaned, "//", "/")
	}
	if cleaned != "/" {
		cleaned = strings.TrimRight(cleaned, "/")
	}
	return cleaned
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
		fullPath := ftpJoin(s.rootPath, newPath)

		if info, err := statPath(fullPath); err == nil && info.IsDir() {
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
			client.cwd = s.resolvePath(client.cwd, "..")
		}
		client.sendResponse("250 Directory changed")

	case "TYPE":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		if args == "I" || args == "A" || args == "A N" {
			t := args
			if strings.HasPrefix(args, "A") {
				t = "A"
				client.binaryMode = false
			} else {
				client.binaryMode = true
			}
			client.sendResponse("200 Type set to " + t)
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

	case "EPRT":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleEPRT(client, args)

	case "LIST", "NLST":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleLIST(client, cmd == "LIST")

	case "MLSD":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleMLSD(client)

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

	case "APPE":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleAPPE(client, args)

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
		client.rnfrName = s.resolvePath(client.cwd, args)
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

	case "MDTM":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		s.handleMDTM(client, args)

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
		client.sendResponse(" EPRT")
		client.sendResponse(" UTF8")
		client.sendResponse(" MLST type*;size*;modify*;perm*;")
		client.sendResponse(" MLSD")
		client.sendResponse(" REST STREAM")
		client.sendResponse(" SIZE")
		client.sendResponse(" MDTM")
		client.sendResponse("211 End")

	case "OPTS":
		if args == "UTF8 ON" {
			client.sendResponse("200 UTF8 enabled")
		} else {
			client.sendResponse("501 Unknown option")
		}

	case "SYST":
		client.sendResponse("215 UNIX Type: L8")

	case "STAT":
		if !client.auth {
			client.sendResponse("530 Not logged in")
			return
		}
		client.sendResponse("213-FTP Server Status")
		client.sendResponse(fmt.Sprintf(" Connected from %s", client.conn.RemoteAddr().String()))
		client.sendResponse(fmt.Sprintf(" Current directory: %s", client.cwd))
		client.sendResponse("213 End")

	case "NOOP":
		client.sendResponse("200 OK")

	case "ABOR":
		if client.dataConn != nil {
			client.dataConn.Close()
			client.dataConn = nil
		}
		client.sendResponse("226 Abort successful")

	default:
		client.sendResponse("502 Command not implemented: " + cmd)
	}
}
