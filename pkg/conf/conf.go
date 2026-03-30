package conf

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         string `yaml:"port"`
	FTPPort      string `yaml:"ftp_port"`
	WebDAVPort   string `yaml:"webdav_port"`
	NFSPort      string `yaml:"nfs_port"`
	Root         string `yaml:"root"`
	AutoSelect   bool   `yaml:"auto_select"`
	EnableFTP    bool   `yaml:"enable_ftp"`
	EnableMCP    bool   `yaml:"enable_mcp"`
	EnableWebDAV bool   `yaml:"enable_webdav"`
	EnableNFS    bool   `yaml:"enable_nfs"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Language     string `yaml:"language"`

	Version     string `yaml:"-"`
	LastCommit  string `yaml:"-"`
	BuildDate   string `yaml:"-"`
	GoVersion   string `yaml:"-"`
	Platform    string `yaml:"-"`
	ProjectURL  string `yaml:"-"`
	ProjectName string `yaml:"-"`

	HTTPAuth struct {
		Enabled  bool   `yaml:"enabled"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"http_auth"`

	HTTPS struct {
		Enabled  bool   `yaml:"enabled"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"https"`

	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`

	Upload struct {
		Enabled    bool   `yaml:"enabled"`
		MaxSize    int64  `yaml:"max_size"`
		AllowTypes string `yaml:"allow_types"`
	} `yaml:"upload"`
}

var AppConfig *Config
var currentConfigPath string

var configPaths = []string{
	"./upftp.yaml",
	"./upftp.yml",
	"~/.upftp/config.yaml",
	"/etc/upftp/config.yaml",
}

func GetConfigPath() string {
	if currentConfigPath != "" {
		return currentConfigPath
	}
	return "defaults"
}

func Init(version, lastCommit, buildDate, goVersion, platform, projectURL, projectName string) {
	AppConfig = &Config{
		Version:     version,
		LastCommit:  lastCommit,
		BuildDate:   buildDate,
		GoVersion:   goVersion,
		Platform:    platform,
		ProjectURL:  projectURL,
		ProjectName: projectName,
		Port:        ":10000",
		FTPPort:     ":2121",
		WebDAVPort:  ":8080",
		NFSPort:     ":2049",
		Root:        "./",
		Username:    "admin",
		Password:    "admin",
		Language:    "",
	}

	loadConfigFile()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "upftp - AI-first lightweight file sharing server\n\n")
		fmt.Fprintf(os.Stderr, "Project: https://github.com/zy84338719/upftp\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  upftp [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -p <port>       HTTP server port (default: 10000)\n")
		fmt.Fprintf(os.Stderr, "  -ftp <port>     FTP server port (default: 2121)\n")
		fmt.Fprintf(os.Stderr, "  -webdav <port>  WebDAV server port (default: 8080)\n")
		fmt.Fprintf(os.Stderr, "  -nfs <port>     NFS server port (default: 2049)\n")
		fmt.Fprintf(os.Stderr, "  -d <dir>        Share directory (default: current directory)\n")
		fmt.Fprintf(os.Stderr, "  -auto           Automatically select first available network interface\n")
		fmt.Fprintf(os.Stderr, "  -enable-ftp     Enable FTP server (default: false)\n")
		fmt.Fprintf(os.Stderr, "  -enable-mcp     Enable MCP server for AI integration (default: false)\n")
		fmt.Fprintf(os.Stderr, "  -enable-webdav  Enable WebDAV server (default: false)\n")
		fmt.Fprintf(os.Stderr, "  -enable-nfs     Enable NFS server (default: false)\n")
		fmt.Fprintf(os.Stderr, "  -config <file>  Configuration file path (default: ./upftp.yaml)\n")
		fmt.Fprintf(os.Stderr, "  -user <name>    FTP username (default: admin)\n")
		fmt.Fprintf(os.Stderr, "  -pass <pass>    FTP password (default: admin)\n")
		fmt.Fprintf(os.Stderr, "  -http-auth      Enable HTTP authentication (default: false)\n")
		fmt.Fprintf(os.Stderr, "  -http-user <name>  HTTP auth username (default: admin)\n")
		fmt.Fprintf(os.Stderr, "  -http-pass <pass>  HTTP auth password (default: admin)\n")
		fmt.Fprintf(os.Stderr, "  -h              Show this help message\n")
	}

	p := flag.String("p", "", "HTTP server port")
	ftpPort := flag.String("ftp", "", "FTP server port")
	webdavPort := flag.String("webdav", "", "WebDAV server port")
	nfsPort := flag.String("nfs", "", "NFS server port")
	dir := flag.String("d", "", "Share directory")
	autoIP := flag.Bool("auto", false, "Automatically select first available network interface")
	enableFTP := flag.Bool("enable-ftp", false, "Enable FTP server")
	enableMCP := flag.Bool("enable-mcp", false, "Enable MCP server for AI integration")
	enableWebDAV := flag.Bool("enable-webdav", false, "Enable WebDAV server")
	enableNFS := flag.Bool("enable-nfs", false, "Enable NFS server")
	configFile := flag.String("config", "", "Configuration file path")
	user := flag.String("user", "", "FTP username")
	pass := flag.String("pass", "", "FTP password")
	httpAuthEnabled := flag.Bool("http-auth", false, "Enable HTTP authentication")
	httpAuthUser := flag.String("http-user", "", "HTTP auth username")
	httpAuthPass := flag.String("http-pass", "", "HTTP auth password")

	flag.Parse()

	if *configFile != "" {
		currentConfigPath = *configFile
		loadConfigFromFile(*configFile)
	}

	if *p != "" {
		AppConfig.Port = ":" + *p
	}
	if *ftpPort != "" {
		AppConfig.FTPPort = ":" + *ftpPort
	}
	if *webdavPort != "" {
		AppConfig.WebDAVPort = ":" + *webdavPort
	}
	if *nfsPort != "" {
		AppConfig.NFSPort = ":" + *nfsPort
	}
	if *dir != "" {
		AppConfig.Root = *dir
	}
	if *autoIP {
		AppConfig.AutoSelect = true
	}
	if *enableFTP {
		AppConfig.EnableFTP = true
	}
	if *enableMCP {
		AppConfig.EnableMCP = true
	}
	if *enableWebDAV {
		AppConfig.EnableWebDAV = true
	}
	if *enableNFS {
		AppConfig.EnableNFS = true
	}
	if *user != "" {
		AppConfig.Username = *user
	}
	if *pass != "" {
		AppConfig.Password = *pass
	}

	if *httpAuthEnabled {
		AppConfig.HTTPAuth.Enabled = true
	}
	if *httpAuthUser != "" {
		AppConfig.HTTPAuth.Username = *httpAuthUser
	}
	if *httpAuthPass != "" {
		AppConfig.HTTPAuth.Password = *httpAuthPass
	}

	if AppConfig.Language == "" {
		AppConfig.Language = "en"
	}
	if AppConfig.Logging.Level == "" {
		AppConfig.Logging.Level = "info"
	}
	if AppConfig.Logging.Format == "" {
		AppConfig.Logging.Format = "text"
	}
	if AppConfig.Upload.MaxSize == 0 {
		AppConfig.Upload.MaxSize = 100 * 1024 * 1024
	}
}

func loadConfigFile() {
	for _, path := range configPaths {
		expandedPath := expandPath(path)
		if _, err := os.Stat(expandedPath); err == nil {
			currentConfigPath = expandedPath
			loadConfigFromFile(expandedPath)
			return
		}
	}
}

func loadConfigFromFile(path string) {
	expandedPath := expandPath(path)
	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return
	}

	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse config file: %v\n", err)
	}
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

func (c *Config) GetHTTPPort() int {
	port := c.Port
	if strings.HasPrefix(port, ":") {
		port = port[1:]
	}
	var p int
	fmt.Sscanf(port, "%d", &p)
	if p == 0 {
		return 10000
	}
	return p
}

func (c *Config) GetFTPPort() int {
	port := c.FTPPort
	if strings.HasPrefix(port, ":") {
		port = port[1:]
	}
	var p int
	fmt.Sscanf(port, "%d", &p)
	if p == 0 {
		return 2121
	}
	return p
}

func (c *Config) GetWebDAVPort() int {
	port := c.WebDAVPort
	if strings.HasPrefix(port, ":") {
		port = port[1:]
	}
	var p int
	fmt.Sscanf(port, "%d", &p)
	if p == 0 {
		return 8080
	}
	return p
}

func (c *Config) GetNFSPort() int {
	port := c.NFSPort
	if strings.HasPrefix(port, ":") {
		port = port[1:]
	}
	var p int
	fmt.Sscanf(port, "%d", &p)
	if p == 0 {
		return 2049
	}
	return p
}

func (c *Config) HTTPAddr() string {
	addr := c.Port
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	return addr
}

func (c *Config) FTPAddr() string {
	addr := c.FTPPort
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	return addr
}

func (c *Config) GetLanguage() string {
	if c.Language == "" {
		return "en"
	}
	return c.Language
}

func (c *Config) SetLanguage(lang string) error {
	// 支持的语言列表
	supportedLangs := []string{"en", "zh", "zh-CN", "zh-TW"}
	for _, l := range supportedLangs {
		if l == lang {
			c.Language = lang
			return nil
		}
	}
	return fmt.Errorf("unsupported language: %s", lang)
}

func (c *Config) SetHTTPPort(port string) {
	c.Port = ":" + port
}

func (c *Config) SetFTPPort(port string) {
	c.FTPPort = ":" + port
}

func GenerateSampleConfig() string {
	return `# UPFTP Configuration File
# https://github.com/zy84338719/upftp

# Server settings
port: "10000"
ftp_port: "2121"
webdav_port: "8080"
nfs_port: "2049"
root: "./"
auto_select: false
enable_ftp: false
enable_mcp: false
enable_webdav: false
enable_nfs: false

# FTP credentials
username: "admin"
password: "admin"

# UI language (shared by TUI and Web): "en" or "zh"
language: "en"

# HTTP Basic Authentication
http_auth:
  enabled: false
  username: "admin"
  password: "admin123"

# Logging settings
logging:
  level: "info"      # debug, info, warn, error
  format: "text"     # text, json

# Upload settings
upload:
  enabled: true
  max_size: 104857600  # 100MB in bytes
  allow_types: ""      # empty = all types, or ".jpg,.png,.pdf"
`
}

func SaveConfig() error {
	return SaveConfigToPath(currentConfigPath)
}

func SaveConfigToPath(path string) error {
	if path == "" {
		path = "./upftp.yaml"
	}

	data, err := yaml.Marshal(AppConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	header := `# UPFTP Configuration File
# https://github.com/zy84338719/upftp

`

	content := header + string(data)

	expandedPath := expandPath(path)
	if err := os.WriteFile(expandedPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	if currentConfigPath == "" {
		currentConfigPath = expandedPath
	}

	return nil
}

func GetDefaultConfigPath() string {
	return "./upftp.yaml"
}

// ReloadConfig reloads configuration from the current config file
func ReloadConfig() {
	if currentConfigPath != "" {
		loadConfigFromFile(currentConfigPath)
	}
	if AppConfig.Language == "" {
		AppConfig.Language = "en"
	}
}
