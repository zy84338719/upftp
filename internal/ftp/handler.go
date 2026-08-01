package ftp

import (
	"context"
	"fmt"
	"strconv"
)

// dispatch 把一条 FTP 命令路由到对应处理函数。
func (s *FTPServer) dispatch(ctx context.Context, c *client, cmd, args string) {
	switch cmd {
	// --- 认证 ---
	case "USER":
		s.cmdUser(c, args)
	case "PASS":
		s.cmdPass(c, args)

	// --- 无需登录 ---
	case "QUIT":
		c.send("221 Goodbye")
		c.quit = true
	case "NOOP":
		c.send("200 OK")
	case "FEAT":
		s.cmdFeat(c)
	case "OPTS":
		s.cmdOpts(c, args)
	case "SYST":
		c.send("215 UNIX Type: L8")
	case "ABOR":
		c.closeData()
		c.send("226 Abort successful")

	// --- 需要登录 ---
	case "PWD", "XPWD":
		s.requireAuth(c, func() { c.send(fmt.Sprintf("257 %q is current directory", c.cwd)) })
	case "CWD":
		s.requireAuth(c, func() { s.cmdCWD(c, args) })
	case "CDUP":
		s.requireAuth(c, s.cmdCDUP(c))
	case "TYPE":
		s.requireAuth(c, func() { s.cmdTYPE(c, args) })
	case "PASV":
		s.requireAuth(c, func() { s.cmdPASV(c) })
	case "EPSV":
		s.requireAuth(c, func() { s.cmdEPSV(c) })
	case "PORT":
		s.requireAuth(c, func() { s.cmdPORT(c, args) })
	case "EPRT":
		s.requireAuth(c, func() { s.cmdEPRT(c, args) })
	case "LIST", "NLST":
		s.requireAuth(c, func() { s.cmdLIST(c, cmd == "LIST", args) })
	case "MLSD":
		s.requireAuth(c, func() { s.cmdMLSD(c, args) })
	case "RETR":
		s.requireAuth(c, func() { s.cmdRETR(c, args) })
	case "STOR":
		s.requireAuth(c, func() { s.cmdSTOR(c, args) })
	case "APPE":
		s.requireAuth(c, func() { s.cmdAPPE(c, args) })
	case "MKD", "XMKD":
		s.requireAuthWritable(c, func() { s.cmdMKD(c, args) })
	case "RMD", "XRMD":
		s.requireAuthWritable(c, func() { s.cmdRMD(c, args) })
	case "DELE":
		s.requireAuthWritable(c, func() { s.cmdDELE(c, args) })
	case "RNFR":
		s.requireAuthWritable(c, func() { s.cmdRNFR(c, args) })
	case "RNTO":
		s.requireAuthWritable(c, func() { s.cmdRNTO(c, args) })
	case "SIZE":
		s.requireAuth(c, func() { s.cmdSIZE(c, args) })
	case "MDTM":
		s.requireAuth(c, func() { s.cmdMDTM(c, args) })
	case "REST":
		s.requireAuth(c, func() { s.cmdREST(c, args) })
	case "STAT":
		s.requireAuth(c, func() { s.cmdSTAT(c) })

	default:
		c.send("502 Command not implemented: " + cmd)
	}
}

// requireAuth 在已认证时执行 fn,否则拒绝。
func (s *FTPServer) requireAuth(c *client, fn func()) {
	if !c.auth {
		c.send("530 Not logged in")
		return
	}
	fn()
}

// requireAuthWritable 在已认证且非只读时执行 fn。
func (s *FTPServer) requireAuthWritable(c *client, fn func()) {
	if !c.auth {
		c.send("530 Not logged in")
		return
	}
	if s.opts.ReadOnly {
		c.send("550 Read-only filesystem")
		return
	}
	fn()
}

// --- 认证实现 ---

// cmdUser 处理 USER。修复点 #4:匿名模式下直接放行,不强制特定用户名。
func (s *FTPServer) cmdUser(c *client, name string) {
	c.name = name
	if s.opts.Anonymous {
		// 匿名模式:USER 即视为认证通过,PASS 阶段也直接 230。
		c.auth = true
		c.send("230 Logged in (anonymous)")
		return
	}
	c.send("331 Please send password")
}

// cmdPass 处理 PASS。
func (s *FTPServer) cmdPass(c *client, pass string) {
	if s.opts.Anonymous {
		c.auth = true
		c.send("230 Logged in (anonymous)")
		return
	}
	if c.name == s.opts.User && pass == s.opts.Pass {
		c.auth = true
		c.send("230 Login successful")
		s.logger.Info("ftp user authenticated", "user", c.name)
		return
	}
	c.send("530 Login incorrect")
}

// --- 杂项命令 ---

func (s *FTPServer) cmdFeat(c *client) {
	lines := []string{
		"211-Features:",
		" PASV",
		" EPSV",
		" EPRT",
		" UTF8",
		" MLST type*;size*;modify*;",
		" MLSD",
		" REST STREAM",
		" SIZE",
		" MDTM",
		"211 End",
	}
	for _, l := range lines {
		c.send(l)
	}
}

func (s *FTPServer) cmdOpts(c *client, args string) {
	if args == "UTF8 ON" {
		c.send("200 UTF8 enabled")
		return
	}
	c.send("501 Unknown option")
}

func (s *FTPServer) cmdTYPE(c *client, args string) {
	switch {
	case args == "I":
		c.binaryMode = true
		c.send("200 Type set to I")
	case len(args) >= 1 && args[0] == 'A':
		c.binaryMode = false
		c.send("200 Type set to A")
	default:
		c.send("500 Invalid type")
	}
}

func (s *FTPServer) cmdREST(c *client, args string) {
	pos, err := strconv.ParseInt(args, 10, 64)
	if err != nil || pos < 0 {
		c.send("501 Invalid restart position")
		return
	}
	c.restPos = pos
	c.send(fmt.Sprintf("350 Restart position %d", pos))
}

func (s *FTPServer) cmdSTAT(c *client) {
	c.send("213-FTP Server Status")
	c.send(fmt.Sprintf(" Connected from %s", c.conn.RemoteAddr()))
	c.send(fmt.Sprintf(" Current directory: %s", c.cwd))
	c.send("213 End")
}

// cmdCWD 切换目录。
func (s *FTPServer) cmdCWD(c *client, args string) {
	target := resolveRel(c.cwd, args)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	if info, e := osStat(abs); e == nil && info.IsDir() {
		c.cwd = cleanSlash(target)
		c.send(fmt.Sprintf("250 Directory changed to %s", c.cwd))
		return
	}
	c.send("550 Directory not found")
}

// cmdCDUP 返回上级目录。
func (s *FTPServer) cmdCDUP(c *client) func() {
	return func() {
		if c.cwd != "/" {
			c.cwd = resolveRel(c.cwd, "..")
		}
		c.send("250 Directory changed")
	}
}
