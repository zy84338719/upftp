package ftp

import (
	"bufio"
	"bytes"
	"context"
	"net"
	"path/filepath"
	"testing"
	"time"
)

func TestCleanSlash(t *testing.T) {
	cases := map[string]string{
		"/":          "/",
		"":           "/",
		"//":         "/",
		"/foo/":      "/foo",
		"/foo//bar/": "/foo/bar",
	}
	for in, want := range cases {
		if got := cleanSlash(in); got != want {
			t.Errorf("cleanSlash(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestResolveRel(t *testing.T) {
	cases := []struct{ cwd, arg, want string }{
		{"/", "foo", "/foo"},
		{"/a", "b", "/a/b"},
		{"/a/b", "..", "/a"},
		{"/a/b", "/x/y", "/x/y"},
		{"/a", ".", "/a"},
		{"/a/b", "../../c", "/c"},
	}
	for _, c := range cases {
		if got := resolveRel(c.cwd, c.arg); got != c.want {
			t.Errorf("resolveRel(%q,%q) = %q, want %q", c.cwd, c.arg, got, c.want)
		}
	}
}

func TestSecureJoin_Normal(t *testing.T) {
	root := t.TempDir()
	got, err := SecureJoin(root, "/foo/bar.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(root, "foo", "bar.txt")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestSecureJoin_Escape 验证目录穿越被拦截(修复点 #1 的核心)。
func TestSecureJoin_Escape(t *testing.T) {
	root := t.TempDir()
	_, err := SecureJoin(root, "/../../../etc/passwd")
	if err == nil {
		t.Error("expected escape error for traversal path, got nil")
	}
}

func TestSecureJoin_RootItself(t *testing.T) {
	root := t.TempDir()
	got, err := SecureJoin(root, "/")
	if err != nil {
		t.Fatalf("root access failed: %v", err)
	}
	if got != root {
		t.Errorf("got %q, want root %q", got, root)
	}
}

func TestSecureJoin_DeepNested(t *testing.T) {
	root := t.TempDir()
	got, err := SecureJoin(root, "/a/b/../../c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(root, "c")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSplitCmd(t *testing.T) {
	cases := []struct {
		in        string
		cmd, args string
	}{
		{"RETR /foo.txt", "RETR", "/foo.txt"},
		{"NOOP", "NOOP", ""},
		{"PASS secret", "PASS", "secret"},
	}
	for _, c := range cases {
		cmd, args := splitCmd(c.in)
		if cmd != c.cmd || args != c.args {
			t.Errorf("splitCmd(%q) = (%q,%q), want (%q,%q)", c.in, cmd, args, c.cmd, c.args)
		}
	}
}

// TestFTPServer_AnonymousAuth 验证匿名登录流程(修复点 #4)。
func TestFTPServer_AnonymousAuth(t *testing.T) {
	root := t.TempDir()
	s := New(Options{Root: root, Anonymous: true})

	c := &client{name: "x", cwd: "/", server: s, writer: bufio.NewWriter(&bytes.Buffer{})}

	// 匿名模式:USER 直接 230
	s.cmdUser(c, "anonymous")
	if !c.auth {
		t.Error("anonymous USER should authenticate immediately")
	}
}

// TestFTPServer_PasswordAuth 验证密码登录:错密码不通过、对密码通过。
func TestFTPServer_PasswordAuth(t *testing.T) {
	root := t.TempDir()
	s := New(Options{Root: root, Anonymous: false, User: "admin", Pass: "secret"})

	c := &client{name: "admin", cwd: "/", server: s, writer: bufio.NewWriter(&bytes.Buffer{})}

	s.cmdPass(c, "wrong")
	if c.auth {
		t.Error("wrong password should not authenticate")
	}
	s.cmdPass(c, "secret")
	if !c.auth {
		t.Error("correct password should authenticate")
	}
}

// TestFTPServer_StartStop 烟雾测试:真实监听一个端口,启动后能干净关闭。
func TestFTPServer_StartStop(t *testing.T) {
	root := t.TempDir()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close() // Start 内部会重新 Listen 同端口,偶发冲突可接受

	s := New(Options{Root: root, Port: port, Anonymous: true})
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- s.Start(ctx) }()

	// 等服务器起来
	time.Sleep(100 * time.Millisecond)

	cancel()
	select {
	case <-s.done:
	case <-time.After(time.Second):
		t.Error("server did not stop in time")
	}
}
