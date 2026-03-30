package ftp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

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

	if strings.Contains(host, ":") {
		client.dataConn.Close()
		client.dataConn = nil
		client.sendResponse("425 IPv6 not supported by PASV, use EPSV")
		return
	}

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
	if client.dataConn != nil {
		client.dataConn.Close()
		client.dataConn = nil
	}
	client.sendResponse("200 PORT command successful")
}

func (s *FTPServer) handleEPRT(client *FTPClient, args string) {
	if len(args) < 3 || args[0] != '|' {
		client.sendResponse("501 Invalid EPRT parameter")
		return
	}

	parts := strings.Split(args[1:], "|")
	if len(parts) < 3 {
		client.sendResponse("501 Invalid EPRT parameter")
		return
	}

	protocol := parts[0]
	host := parts[1]
	portStr := parts[2]

	if protocol != "1" && protocol != "2" {
		client.sendResponse("522 Protocol not supported, use (1,2)")
		return
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		client.sendResponse("501 Invalid port")
		return
	}

	client.dataPort = fmt.Sprintf("%s:%d", host, port)
	if client.dataConn != nil {
		client.dataConn.Close()
		client.dataConn = nil
	}
	client.sendResponse("200 EPRT command successful")
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
