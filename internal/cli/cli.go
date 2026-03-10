package cli

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/zy84338719/upftp/internal/config"
)

type CLI struct {
	serverIP string
	fileMap  map[string]string
}

func NewCLI() *CLI {
	return &CLI{
		fileMap: make(map[string]string),
	}
}

func (c *CLI) SetServerIP(ip string) {
	c.serverIP = ip
}

func (c *CLI) ScanDirectory(pathname string) map[string]string {
	files := make(map[string]string)
	c.scanDir(pathname, files)
	return files
}

func (c *CLI) scanDir(pathname string, m map[string]string) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", pathname, err)
		return
	}

	for _, fi := range rd {
		if fi.IsDir() {
			c.scanDir(path.Join(pathname, fi.Name()), m)
		} else {
			dir := strings.Replace(pathname, config.AppConfig.Root, "", 1)
			var downloadURL string
			if len(dir) > 0 {
				dir = path.Join(dir)
				downloadURL = fmt.Sprintf("http://%s%s/download/%s/%s",
					c.serverIP, config.AppConfig.Port, dir, fi.Name())
			} else {
				downloadURL = fmt.Sprintf("http://%s%s/download/%s",
					c.serverIP, config.AppConfig.Port, fi.Name())
			}
			m[fi.Name()] = downloadURL
		}
	}
}

func (c *CLI) Start(ctx context.Context, s chan os.Signal) {
	c.fileMap = c.ScanDirectory(config.AppConfig.Root)

	c.printBanner()
	c.runLoop(s)
}

func (c *CLI) printBanner() {
	fmt.Printf("\n")
	fmt.Printf("╔══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                                                                      ║\n")
	fmt.Printf("║    ██╗   ██╗██████╗ ███████╗██████╗ ███████╗██████╗                  ║\n")
	fmt.Printf("║    ██║   ██║██╔══██╗██╔════╝██╔══██╗██╔════╝██╔══██╗                 ║\n")
	fmt.Printf("║    ██║   ██║██████╔╝█████╗  ██████╔╝█████╗  ██║  ██║                 ║\n")
	fmt.Printf("║    ╚██╗ ██╔╝██╔══██╗██╔══╝  ██╔══██╗██╔══╝  ██║  ██║                 ║\n")
	fmt.Printf("║     ╚████╔╝ ██████╔╝███████╗██║  ██║███████╗██████╔╝                 ║\n")
	fmt.Printf("║      ╚═══╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝╚══════╝╚═════╝                  ║\n")
	fmt.Printf("║                                                                      ║\n")
	fmt.Printf("║              AI-First File Sharing Server v%-24s║\n", config.AppConfig.Version)
	fmt.Printf("╠══════════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                      ║\n")
	fmt.Printf("║  🌐 HTTP Server:    http://%-40s║\n", c.serverIP+config.AppConfig.Port)
	if config.AppConfig.EnableFTP {
		fmt.Printf("║  📁 FTP Server:     ftp://%-41s║\n", c.serverIP+config.AppConfig.FTPPort)
		fmt.Printf("║  🔑 FTP Credentials: %-47s║\n", config.AppConfig.Username+" / "+config.AppConfig.Password)
	}
	fmt.Printf("║  📂 Shared Path:    %-47s║\n", truncatePath(config.AppConfig.Root, 47))
	fmt.Printf("║  📄 Files Found:    %-47d║\n", len(c.fileMap))
	fmt.Printf("║                                                                      ║\n")
	fmt.Printf("╠══════════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║  ✨ Features:                                                        ║\n")
	fmt.Printf("║     • Modern Web Interface with Preview                              ║\n")
	if config.AppConfig.Upload.Enabled {
		fmt.Printf("║     • File Upload Enabled (Max: %-36s║\n", formatSize(config.AppConfig.Upload.MaxSize)+")")
	}
	if config.AppConfig.HTTPAuth.Enabled {
		fmt.Printf("║     • HTTP Basic Auth Enabled                                        ║\n")
	}
	fmt.Printf("║     • MCP Server for AI Integration                                  ║\n")
	fmt.Printf("║     • QR Code for Mobile Access                                      ║\n")
	fmt.Printf("║                                                                      ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════════════════════╝\n")
}

func (c *CLI) runLoop(s chan os.Signal) {
	for {
		c.printMenu()
		c.handleInput(s)
		fmt.Println()
	}
}

func (c *CLI) printMenu() {
	fmt.Printf("┌─────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│                        COMMAND MENU                             │\n")
	fmt.Printf("├─────────────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│  [1] 🔍 Search files          [2] 📋 List all files            │\n")
	fmt.Printf("│  [3] 📥 Download examples     [4] 🔄 Refresh file list         │\n")
	if config.AppConfig.EnableFTP {
		fmt.Printf("│  [5] 📁 FTP connection info   [6] ℹ️  Server status            │\n")
	} else {
		fmt.Printf("│  [5] ℹ️  Server status        [7] 📖 About & Features         │\n")
	}
	fmt.Printf("│  [v] 📌 Version info          [q] 🚪 Quit server               │\n")
	fmt.Printf("└─────────────────────────────────────────────────────────────────┘\n")
	fmt.Printf("\nEnter command: ")
}

func (c *CLI) handleInput(s chan os.Signal) {
	var option string
	_, _ = fmt.Scanln(&option)

	switch strings.ToLower(option) {
	case "1":
		c.searchFiles()
	case "2":
		c.listAllFiles()
	case "3":
		c.showDownloadExamples()
	case "4":
		c.refreshFileList()
	case "5":
		if config.AppConfig.EnableFTP {
			c.showFTPInfo()
		} else {
			c.showServerStatus()
		}
	case "6":
		if config.AppConfig.EnableFTP {
			c.showServerStatus()
		} else {
			c.showAbout()
		}
	case "7":
		c.showAbout()
	case "v", "version":
		c.showVersion()
	case "q", "quit", "exit":
		fmt.Println("\n👋 Shutting down server...")
		s <- syscall.SIGQUIT
	default:
		fmt.Println("❌ Invalid option, please try again.")
	}
}

func (c *CLI) showVersion() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                      VERSION INFORMATION                      ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  UPFTP Version:    %-43s║\n", config.AppConfig.Version)
	fmt.Printf("║  Git Commit:       %-43s║\n", config.AppConfig.LastCommit)
	fmt.Printf("║  Go Version:       %-43s║\n", "go1.23+")
	fmt.Printf("║  Build Date:       %-43s║\n", getBuildDate())
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  Repository:       https://github.com/zy84338719/upftp         ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

func (c *CLI) showServerStatus() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                        SERVER STATUS                           ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📊 Server Configuration:                                      ║\n")
	fmt.Printf("║     HTTP Port:      %-41s║\n", config.AppConfig.Port[1:])
	if config.AppConfig.EnableFTP {
		fmt.Printf("║     FTP Port:       %-41s║\n", config.AppConfig.FTPPort[1:])
	}
	fmt.Printf("║     Root Directory: %-41s║\n", truncatePath(config.AppConfig.Root, 41))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  ⚙️  Features Status:                                           ║\n")
	fmt.Printf("║     HTTP Auth:      %-41s║\n", boolToStr(config.AppConfig.HTTPAuth.Enabled))
	fmt.Printf("║     File Upload:    %-41s║\n", boolToStr(config.AppConfig.Upload.Enabled))
	fmt.Printf("║     FTP Server:     %-41s║\n", boolToStr(config.AppConfig.EnableFTP))
	fmt.Printf("║     HTTPS:          %-41s║\n", boolToStr(config.AppConfig.HTTPS.Enabled))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📁 Statistics:                                                ║\n")
	fmt.Printf("║     Files Available: %-39d║\n", len(c.fileMap))
	fmt.Printf("║     Upload Max Size: %-39s║\n", formatSize(config.AppConfig.Upload.MaxSize))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

func (c *CLI) showAbout() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                    ABOUT UPFTP v%-29s║\n", config.AppConfig.Version)
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  UPFTP is an AI-first lightweight file sharing server.        ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  ✨ Key Features:                                              ║\n")
	fmt.Printf("║     • Modern responsive Web UI with file preview              ║\n")
	fmt.Printf("║     • MCP (Model Context Protocol) for AI integration         ║\n")
	fmt.Printf("║     • File upload with drag & drop support                    ║\n")
	fmt.Printf("║     • HTTP Basic Authentication                               ║\n")
	fmt.Printf("║     • Full FTP server implementation                          ║\n")
	fmt.Printf("║     • QR code for easy mobile access                          ║\n")
	fmt.Printf("║     • Directory tree navigation                               ║\n")
	fmt.Printf("║     • Multi-language support (EN/ZH)                          ║\n")
	fmt.Printf("║     • YAML configuration file support                         ║\n")
	fmt.Printf("║     • HTTPS support with custom certificates                  ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  🔗 Links:                                                     ║\n")
	fmt.Printf("║     GitHub:    https://github.com/zy84338719/upftp            ║\n")
	fmt.Printf("║     Issues:    https://github.com/zy84338719/upftp/issues     ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📄 License: MIT                                               ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

func (c *CLI) searchFiles() {
	fmt.Print("🔍 Enter search term (or press Enter to show all): ")
	reader := bufio.NewReader(os.Stdin)
	term, _ := reader.ReadString('\n')
	term = strings.TrimSpace(term)

	found := false
	count := 0
	fmt.Printf("\n┌────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│                      SEARCH RESULTS                            │\n")
	fmt.Printf("├────────────────────────────────────────────────────────────────┤\n")

	for filename, url := range c.fileMap {
		if term == "" || strings.Contains(strings.ToLower(filename), strings.ToLower(term)) {
			found = true
			count++
			fmt.Printf("│ 📄 %-58s│\n", truncateString(filename, 56))
			fmt.Printf("│    %-58s│\n", truncateString(url, 56))
			fmt.Printf("├────────────────────────────────────────────────────────────────┤\n")
			if count >= 20 {
				fmt.Printf("│ ... and more files (showing first 20)                         │\n")
				break
			}
		}
	}

	if !found {
		fmt.Printf("│ ❌ No matching files found.                                    │\n")
	}
	fmt.Printf("└────────────────────────────────────────────────────────────────┘\n")
	fmt.Printf("\nFound %d files.\n", count)
}

func (c *CLI) listAllFiles() {
	count := 0
	fmt.Printf("\n┌────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│                        ALL FILES                                │\n")
	fmt.Printf("├────────────────────────────────────────────────────────────────┤\n")

	for filename, url := range c.fileMap {
		count++
		fmt.Printf("│ 📄 %-58s│\n", truncateString(filename, 56))
		fmt.Printf("│    %-58s│\n", truncateString(url, 56))
		fmt.Printf("├────────────────────────────────────────────────────────────────┤\n")
		if count >= 20 {
			fmt.Printf("│ ... and %d more files                                           │\n", len(c.fileMap)-20)
			break
		}
	}
	fmt.Printf("└────────────────────────────────────────────────────────────────┘\n")
	fmt.Printf("\nTotal: %d files\n", len(c.fileMap))
}

func (c *CLI) showDownloadExamples() {
	if len(c.fileMap) == 0 {
		fmt.Println("❌ No files available for download examples.")
		return
	}

	var exampleFile, exampleURL string
	for filename, url := range c.fileMap {
		exampleFile = filename
		exampleURL = url
		break
	}

	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                      DOWNLOAD EXAMPLES                         ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📄 Example file: %-44s║\n", truncateString(exampleFile, 44))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  🌐 Browser:                                                   ║\n")
	fmt.Printf("║     http://%s%s\n", c.serverIP, config.AppConfig.Port)
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  💻 Command Line Tools:                                        ║\n")
	fmt.Printf("║     curl -O \"%s\n", truncateString(exampleURL, 52))
	fmt.Printf("║     wget \"%s\n", truncateString(exampleURL, 54))
	fmt.Printf("║                                                                ║\n")
	if config.AppConfig.EnableFTP {
		fmt.Printf("║  📁 FTP Client:                                                ║\n")
		fmt.Printf("║     ftp %s\n", c.serverIP)
		fmt.Printf("║     Username: %s\n", config.AppConfig.Username)
		fmt.Printf("║     Password: %s\n", config.AppConfig.Password)
		fmt.Printf("║                                                                ║\n")
	}
	fmt.Printf("║  🤖 MCP Integration (Claude Desktop):                          ║\n")
	fmt.Printf("║     Add to Claude config:                                      ║\n")
	fmt.Printf("║     \"command\": \"upftp\",                                        ║\n")
	fmt.Printf("║     \"args\": [\"-enable-mcp\"]                                    ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

func (c *CLI) showFTPInfo() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                     FTP CONNECTION INFO                        ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📁 Server:      %-44s║\n", c.serverIP)
	fmt.Printf("║  🔌 Port:        %-44s║\n", config.AppConfig.FTPPort[1:])
	fmt.Printf("║  👤 Username:    %-44s║\n", config.AppConfig.Username)
	fmt.Printf("║  🔑 Password:    %-44s║\n", config.AppConfig.Password)
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  💻 Example FTP commands:                                      ║\n")
	fmt.Printf("║     $ ftp %s\n", c.serverIP)
	fmt.Printf("║     Name: %s\n", config.AppConfig.Username)
	fmt.Printf("║     Password: %s\n", config.AppConfig.Password)
	fmt.Printf("║     ftp> ls                                                     ║\n")
	fmt.Printf("║     ftp> get filename                                           ║\n")
	fmt.Printf("║     ftp> put localfile                                          ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

func (c *CLI) refreshFileList() {
	fmt.Println("🔄 Refreshing file list...")
	c.fileMap = c.ScanDirectory(config.AppConfig.Root)
	fmt.Printf("✅ Found %d files in directory.\n", len(c.fileMap))
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func truncatePath(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return "..." + s[len(s)-maxLen+3:]
}

func boolToStr(b bool) string {
	if b {
		return "✅ Enabled"
	}
	return "❌ Disabled"
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getBuildDate() string {
	return "2026-03-11"
}
