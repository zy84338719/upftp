// Package config 定义 upftp 的配置模型与加载逻辑。
//
// 配置来源优先级:命令行 flag > 配置文件(yaml) > 默认值。
// 默认值面向"即开即用":匿名访问、自动选端口、当前目录。
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 是 upftp 的完整配置。
type Config struct {
	// Dir 是共享根目录,默认当前目录。
	Dir string `yaml:"dir"`
	// FTPPort 是 FTP 监听端口,0 表示自动分配。
	FTPPort int `yaml:"ftp_port"`
	// HTTPPort 是 HTTP/Web 监听端口,0 表示自动分配。
	HTTPPort int `yaml:"http_port"`
	// Host 对外展示的地址(IP),用于拼接下载链接/二维码。
	Host string `yaml:"host"`
	// User 用户名,为空表示匿名访问。
	User string `yaml:"user"`
	// Pass 密码,User 为空时忽略。
	Pass string `yaml:"pass"`
	// Anonymous 是否允许匿名,User 为空时自动为 true。
	Anonymous bool `yaml:"anonymous"`
	// ReadOnly 只读模式,禁止上传/删除/改名。
	ReadOnly bool `yaml:"read_only"`
	// TUI 是否启动终端交互界面,非 TTY 时自动关闭。
	TUI bool `yaml:"tui"`
	// LogLevel: debug / info / warn / error
	LogLevel string `yaml:"log_level"`
}

// Default 返回面向即开即用的默认配置。
func Default() *Config {
	dir, _ := os.Getwd()
	return &Config{
		Dir:       dir,
		FTPPort:   2121,
		HTTPPort:  8080,
		Host:      "",
		Anonymous: true,
		TUI:       true,
		LogLevel:  "info",
	}
}

// AuthEnabled 表示是否启用了凭据认证(非匿名)。
func (c *Config) AuthEnabled() bool {
	return c.User != ""
}

// LoadFromFile 从 yaml 文件读取并合并到当前配置(仅覆盖非零字段)。
func (c *Config) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// 用临时结构解析,再按"非零值才覆盖"合并,避免 false/0 把默认值清掉。
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	if v, ok := raw["dir"].(string); ok && v != "" {
		c.Dir = v
	}
	if v, ok := raw["ftp_port"].(int); ok && v != 0 {
		c.FTPPort = v
	}
	if v, ok := raw["http_port"].(int); ok && v != 0 {
		c.HTTPPort = v
	}
	if v, ok := raw["host"].(string); ok && v != "" {
		c.Host = v
	}
	if v, ok := raw["user"].(string); ok && v != "" {
		c.User = v
		c.Anonymous = false
	}
	if v, ok := raw["pass"].(string); ok {
		c.Pass = v
	}
	if v, ok := raw["anonymous"].(bool); ok {
		c.Anonymous = v
	}
	if v, ok := raw["read_only"].(bool); ok {
		c.ReadOnly = v
	}
	if v, ok := raw["tui"].(bool); ok {
		c.TUI = v
	}
	if v, ok := raw["log_level"].(string); ok && v != "" {
		c.LogLevel = v
	}
	return nil
}

// DefaultConfigPaths 返回按优先级查找配置文件的候选路径。
func DefaultConfigPaths() []string {
	home, _ := os.UserHomeDir()
	return []string{
		filepath.Join(".", "upftp.yaml"),
		filepath.Join(".", "upftp.yml"),
		filepath.Join(home, ".upftp", "config.yaml"),
		filepath.Join("/etc", "upftp", "config.yaml"),
	}
}

// ParsePort 容错解析端口字符串(支持 "8080" / ":8080")。
func ParsePort(s string) (int, error) {
	for len(s) > 0 && (s[0] == ':' || s[0] == ' ') {
		s = s[1:]
	}
	p, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q", s)
	}
	if p < 0 || p > 65535 {
		return 0, fmt.Errorf("port %d out of range", p)
	}
	return p, nil
}
