package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port       string `yaml:"port"`
	FTPPort    string `yaml:"ftp_port"`
	Root       string `yaml:"root"`
	AutoSelect bool   `yaml:"auto_select"`
	EnableFTP  bool   `yaml:"enable_ftp"`
	EnableMCP  bool   `yaml:"enable_mcp"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`

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
		Root:        "./",
		Username:    "admin",
		Password:    "admin",
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
		fmt.Fprintf(os.Stderr, "  -d <dir>        Share directory (default: current directory)\n")
		fmt.Fprintf(os.Stderr, "  -auto           Automatically select first available network interface\n")
		fmt.Fprintf(os.Stderr, "  -enable-ftp     Enable FTP server (default: false)\n")
		fmt.Fprintf(os.Stderr, "  -enable-mcp     Enable MCP server for AI integration (default: false)\n")
		fmt.Fprintf(os.Stderr, "  -config <file>  Configuration file path (default: ./upftp.yaml)\n")
		fmt.Fprintf(os.Stderr, "  -user <name>    FTP username (default: admin)\n")
		fmt.Fprintf(os.Stderr, "  -pass <pass>    FTP password (default: admin)\n")
		fmt.Fprintf(os.Stderr, "  -h              Show this help message\n")
	}

	p := flag.String("p", "", "HTTP server port")
	ftpPort := flag.String("ftp", "", "FTP server port")
	dir := flag.String("d", "", "Share directory")
	autoIP := flag.Bool("auto", false, "Automatically select first available network interface")
	enableFTP := flag.Bool("enable-ftp", false, "Enable FTP server")
	enableMCP := flag.Bool("enable-mcp", false, "Enable MCP server for AI integration")
	configFile := flag.String("config", "", "Configuration file path")
	user := flag.String("user", "", "FTP username")
	pass := flag.String("pass", "", "FTP password")

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
	if *user != "" {
		AppConfig.Username = *user
	}
	if *pass != "" {
		AppConfig.Password = *pass
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
	AppConfig.Upload.Enabled = true
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

func (c *Config) HTTPAddr() string {
	return c.Port
}

func (c *Config) FTPAddr() string {
	return c.FTPPort
}

func GenerateSampleConfig() string {
	return `# UPFTP Configuration File
# https://github.com/zy84338719/upftp

# Server settings
port: "10000"
ftp_port: "2121"
root: "./"
auto_select: false
enable_ftp: false
enable_mcp: false

# FTP credentials
username: "admin"
password: "admin"

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
