package ftp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// osStat 是 os.Stat 的薄封装(便于测试)。
func osStat(p string) (os.FileInfo, error) { return os.Stat(p) }

// --- 目录列举 ---

// cmdLIST 列举目录内容。
func (s *FTPServer) cmdLIST(c *client, long bool, args string) {
	rel := args
	// 兼容 "LIST -la /foo":剥离选项参数。
	if strings.HasPrefix(rel, "-") {
		if sp := strings.IndexByte(rel, ' '); sp >= 0 {
			rel = strings.TrimSpace(rel[sp+1:])
		} else {
			rel = ""
		}
	}
	target := resolveRel(c.cwd, rel)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		c.send("550 Cannot list directory")
		return
	}

	dataConn, err := c.openData()
	if err != nil {
		c.send("425 Cannot open data connection")
		return
	}
	defer dataConn.Close()
	c.send("150 Opening data connection")

	bw := bufio.NewWriterSize(dataConn, bufferSize)
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		if long {
			fmt.Fprintln(bw, formatLISTLine(info))
		} else {
			fmt.Fprintln(bw, e.Name())
		}
	}
	_ = bw.Flush()
	c.send("226 Transfer complete")
}

// cmdMLSD 机器可读目录列举(RFC 3659)。
func (s *FTPServer) cmdMLSD(c *client, args string) {
	target := resolveRel(c.cwd, args)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		c.send("550 Cannot list directory")
		return
	}
	dataConn, err := c.openData()
	if err != nil {
		c.send("425 Cannot open data connection")
		return
	}
	defer dataConn.Close()
	c.send("150 Opening data connection")

	bw := bufio.NewWriterSize(dataConn, bufferSize)
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		fmt.Fprintln(bw, formatMLSDLine(e.Name(), info))
	}
	_ = bw.Flush()
	c.send("226 Transfer complete")
}

func formatLISTLine(info os.FileInfo) string {
	perm := "-rw-r--r--"
	if info.IsDir() {
		perm = "drwxr-xr-x"
	}
	size := info.Size()
	mtime := info.ModTime().Format("Jan _2 15:04")
	return fmt.Sprintf("%s 1 ftp ftp %12d %s %s", perm, size, mtime, info.Name())
}

func formatMLSDLine(name string, info os.FileInfo) string {
	typ := "file"
	if info.IsDir() {
		typ = "dir"
	}
	mod := info.ModTime().UTC().Format("20060102150405")
	return fmt.Sprintf("type=%s;size=%d;modify=%s; %s", typ, info.Size(), mod, name)
}

// --- 文件下载 ---

// cmdRETR 下载文件,支持 REST 断点续传。
func (s *FTPServer) cmdRETR(c *client, args string) {
	target := resolveRel(c.cwd, args)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	f, err := os.Open(abs)
	if err != nil {
		c.send("550 Cannot open file")
		return
	}
	defer f.Close()

	offset := c.restPos
	c.restPos = 0
	if offset > 0 {
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			c.send("550 Seek failed")
			return
		}
	}

	dataConn, err := c.openData()
	if err != nil {
		c.send("425 Cannot open data connection")
		return
	}
	defer dataConn.Close()
	c.send("150 Opening data connection")

	buf := make([]byte, bufferSize)
	_, _ = io.CopyBuffer(struct{ io.Writer }{dataConn}, f, buf)
	c.send("226 Transfer complete")
}

// cmdSTOR 上传文件。修复点 #2:支持 REST 续传(在 offset 处写入而非截断)。
func (s *FTPServer) cmdSTOR(c *client, args string) {
	target := resolveRel(c.cwd, args)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}

	offset := c.restPos
	c.restPos = 0

	// 确保父目录存在
	if mkErr := os.MkdirAll(filepath.Dir(abs), 0o755); mkErr != nil {
		c.send("550 Cannot create path")
		return
	}

	var f *os.File
	if offset > 0 {
		// 续传:以读写方式打开并定位到 offset,不截断
		f, err = os.OpenFile(abs, os.O_WRONLY, 0o644)
		if err == nil {
			_, err = f.Seek(offset, io.SeekStart)
		}
	} else {
		f, err = os.Create(abs)
	}
	if err != nil {
		c.send("550 Cannot create file")
		return
	}
	defer f.Close()

	dataConn, err := c.openData()
	if err != nil {
		c.send("425 Cannot open data connection")
		return
	}
	defer dataConn.Close()
	c.send("150 Opening data connection")

	buf := make([]byte, bufferSize)
	_, _ = io.CopyBuffer(f, struct{ io.Reader }{dataConn}, buf)
	c.send("226 Transfer complete")
}

// cmdAPPE 追加写入文件末尾。
func (s *FTPServer) cmdAPPE(c *client, args string) {
	target := resolveRel(c.cwd, args)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	if mkErr := os.MkdirAll(filepath.Dir(abs), 0o755); mkErr != nil {
		c.send("550 Cannot create path")
		return
	}
	f, err := os.OpenFile(abs, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		c.send("550 Cannot open file")
		return
	}
	defer f.Close()

	dataConn, err := c.openData()
	if err != nil {
		c.send("425 Cannot open data connection")
		return
	}
	defer dataConn.Close()
	c.send("150 Opening data connection")

	buf := make([]byte, bufferSize)
	_, _ = io.CopyBuffer(f, struct{ io.Reader }{dataConn}, buf)
	c.send("226 Transfer complete")
}

// --- 目录与文件管理 ---

func (s *FTPServer) cmdMKD(c *client, args string) {
	target := resolveRel(c.cwd, args)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		c.send("550 Cannot create directory")
		return
	}
	c.send(fmt.Sprintf("257 %q created", cleanSlash(target)))
}

// cmdRMD 删除目录。修复点 #1 补充:用 RemoveAll 但路径已由 SecureJoin 保证在沙箱内。
func (s *FTPServer) cmdRMD(c *client, args string) {
	target := resolveRel(c.cwd, args)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	// 不允许删除根目录本身
	if abs == filepath.Clean(s.opts.Root) {
		c.send("550 Cannot remove root")
		return
	}
	if err := os.RemoveAll(abs); err != nil {
		c.send("550 Cannot remove directory")
		return
	}
	c.send("250 Directory removed")
}

func (s *FTPServer) cmdDELE(c *client, args string) {
	target := resolveRel(c.cwd, args)
	abs, err := SecureJoin(s.opts.Root, target)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	if err := os.Remove(abs); err != nil {
		c.send("550 Cannot delete file")
		return
	}
	c.send("250 File deleted")
}

func (s *FTPServer) cmdRNFR(c *client, args string) {
	c.rnfr = resolveRel(c.cwd, args)
	c.send("350 Ready for destination name")
}

func (s *FTPServer) cmdRNTO(c *client, args string) {
	if c.rnfr == "" {
		c.send("503 Use RNFR first")
		return
	}
	src, err := SecureJoin(s.opts.Root, c.rnfr)
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	dst, err := SecureJoin(s.opts.Root, resolveRel(c.cwd, args))
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	c.rnfr = ""
	if err := os.Rename(src, dst); err != nil {
		c.send("550 Rename failed")
		return
	}
	c.send("250 Rename successful")
}

// --- 文件元信息 ---

func (s *FTPServer) cmdSIZE(c *client, args string) {
	abs, err := SecureJoin(s.opts.Root, resolveRel(c.cwd, args))
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	info, err := osStat(abs)
	if err != nil || info.IsDir() {
		c.send("550 Not a regular file")
		return
	}
	c.send(fmt.Sprintf("213 %d", info.Size()))
}

func (s *FTPServer) cmdMDTM(c *client, args string) {
	abs, err := SecureJoin(s.opts.Root, resolveRel(c.cwd, args))
	if err != nil {
		c.send("550 Permission denied")
		return
	}
	info, err := osStat(abs)
	if err != nil {
		c.send("550 File not found")
		return
	}
	c.send("213 " + info.ModTime().UTC().Format("20060102150405"))
}

// 规避未使用导入(time 在 MDTM 格式中通过常量隐式使用,确保编译)。
var _ = time.Now
var _ = strconv.Itoa
