package ftp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// cmdPASV 被动模式(IPv4):服务端开一个随机端口监听,返回 227。
func (s *FTPServer) cmdPASV(c *client) {
	c.closeData()
	ip, port, err := s.openPassiveListener(c)
	if err != nil {
		c.send("425 Cannot open data connection")
		return
	}
	c.send(fmt.Sprintf("227 Entering Passive Mode (%s)", formatPASVAddr(ip, port)))
}

// cmdEPSV 扩展被动模式(支持 IPv6):返回 229。
func (s *FTPServer) cmdEPSV(c *client) {
	c.closeData()
	_, port, err := s.openPassiveListener(c)
	if err != nil {
		c.send("425 Cannot open data connection")
		return
	}
	c.send(fmt.Sprintf("229 Entering Extended Passive Mode (|||%d|)", port))
}

// openPassiveListener 在随机端口监听,记录到客户端,返回用于对外通告的 ip/port。
func (s *FTPServer) openPassiveListener(c *client) (string, int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", 0, err
	}
	c.dataListener = ln
	tcp := ln.Addr().(*net.TCPAddr)
	// 对外通告的 IP 取自控制连接的本地地址。
	advertiseIP := "127.0.0.1"
	if local, ok := c.conn.LocalAddr().(*net.TCPAddr); ok && len(local.IP) > 0 {
		advertiseIP = local.IP.String()
	}
	return advertiseIP, tcp.Port, nil
}

// cmdPORT 主动模式(IPv4):解析 h1,h2,h3,h4,p1,p2,记录客户端指定的数据端点。
func (s *FTPServer) cmdPORT(c *client, args string) {
	parts := strings.Split(args, ",")
	if len(parts) != 6 {
		c.send("501 Invalid PORT")
		return
	}
	var oct [4]int
	for i := 0; i < 4; i++ {
		n, err := strconv.Atoi(parts[i])
		if err != nil || n < 0 || n > 255 {
			c.send("501 Invalid PORT")
			return
		}
		oct[i] = n
	}
	p1, e1 := strconv.Atoi(parts[4])
	p2, e2 := strconv.Atoi(parts[5])
	if e1 != nil || e2 != nil {
		c.send("501 Invalid PORT")
		return
	}
	c.closeData()
	c.dataPort = fmt.Sprintf("%d.%d.%d.%d:%d", oct[0], oct[1], oct[2], oct[3], p1*256+p2)
	c.send("200 PORT command successful")
}

// cmdEPRT 扩展主动模式:解析 |proto|host|port|。
func (s *FTPServer) cmdEPRT(c *client, args string) {
	// 格式 |2|::1|12345|
	sep := string(args[0])
	fields := strings.Split(args[1:], sep)
	if len(fields) < 3 {
		c.send("501 Invalid EPRT")
		return
	}
	host := fields[1]
	port := fields[2]
	if _, err := strconv.Atoi(port); err != nil {
		c.send("501 Invalid EPRT port")
		return
	}
	c.closeData()
	c.dataPort = "[" + host + "]:" + port
	c.send("200 EPRT command successful")
}

// openData 打开数据连接:被动模式 Accept,主动模式 Dial。
func (c *client) openData() (net.Conn, error) {
	if c.dataListener != nil {
		_ = c.dataListener.(*net.TCPListener).SetDeadline(time.Now().Add(dataTimeout))
		return c.dataListener.Accept()
	}
	if c.dataPort != "" {
		return net.DialTimeout("tcp", c.dataPort, dataTimeout)
	}
	return nil, fmt.Errorf("no data connection negotiated")
}

// closeData 释放协商的数据连接资源。
func (c *client) closeData() {
	if c.dataListener != nil {
		_ = c.dataListener.Close()
		c.dataListener = nil
	}
	c.dataPort = ""
	c.restPos = 0
}

// formatPASVAddr 把 ip:port 编码为 PASV 响应所需的 h1,h2,h3,h4,p1,p2。
func formatPASVAddr(ip string, port int) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		parts = []string{"127", "0", "0", "1"}
	}
	return fmt.Sprintf("%s,%s,%s,%s,%d,%d", parts[0], parts[1], parts[2], parts[3], port/256, port%256)
}
