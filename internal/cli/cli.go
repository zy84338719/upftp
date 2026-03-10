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
	fmt.Printf("║              Build: %-48s║\n", config.AppConfig.LastCommit+" @ "+config.AppConfig.BuildDate)
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
	fmt.Printf("║     • Interactive CLI Configuration                                  ║\n")
	fmt.Printf("║                                                                      ║\n")
	fmt.Printf("╠══════════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║  🔗 Homepage: %-54s║\n", truncateString(config.AppConfig.ProjectURL, 54))
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
		fmt.Printf("│  [7] ⚙️  Configuration        [8] 📖 About & Features         │\n")
	} else {
		fmt.Printf("│  [5] ℹ️  Server status        [6] ⚙️  Configuration            │\n")
		fmt.Printf("│  [7] 📖 About & Features                                     │\n")
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
			c.showConfigMenu(s)
		}
	case "7":
		if config.AppConfig.EnableFTP {
			c.showConfigMenu(s)
		} else {
			c.showAbout()
		}
	case "8":
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
	fmt.Printf("║  🚀 %s Version:    %-42s║\n", config.AppConfig.ProjectName, config.AppConfig.Version)
	fmt.Printf("║  🔧 Git Commit:       %-41s║\n", config.AppConfig.LastCommit)
	fmt.Printf("║  📅 Build Date:       %-41s║\n", config.AppConfig.BuildDate)
	fmt.Printf("║  🐹 Go Version:       %-41s║\n", config.AppConfig.GoVersion)
	fmt.Printf("║  💻 Platform:         %-41s║\n", config.AppConfig.Platform)
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  🏠 Project Homepage: %-40s║\n", truncateString(config.AppConfig.ProjectURL, 40))
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
	fmt.Printf("║                  ABOUT %s v%-32s║\n", config.AppConfig.ProjectName, config.AppConfig.Version)
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  %s is an AI-first lightweight file sharing server.       ║\n", config.AppConfig.ProjectName)
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
	fmt.Printf("║     • Interactive CLI configuration                           ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  🔗 Links:                                                     ║\n")
	fmt.Printf("║     Homepage:  %-47s║\n", truncateString(config.AppConfig.ProjectURL, 47))
	fmt.Printf("║     GitHub:    https://github.com/zy84338719/upftp            ║\n")
	fmt.Printf("║     Issues:    https://github.com/zy84338719/upftp/issues     ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📄 License: MIT                                               ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  💻 Build Info:                                                ║\n")
	fmt.Printf("║     Version:   %-47s║\n", config.AppConfig.Version)
	fmt.Printf("║     Commit:    %-47s║\n", config.AppConfig.LastCommit)
	fmt.Printf("║     Built:     %-47s║\n", config.AppConfig.BuildDate)
	fmt.Printf("║     Platform:  %-47s║\n", config.AppConfig.Platform)
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

func (c *CLI) showConfigMenu(s chan os.Signal) {
	for {
		c.printConfigMenu()
		var option string
		_, _ = fmt.Scanln(&option)

		switch strings.ToLower(option) {
		case "1":
			c.configureFTPCredentials()
		case "2":
			c.configureHTTPAuth()
		case "3":
			c.toggleFTPServer()
		case "4":
			c.toggleMCPServer()
		case "5":
			c.configurePorts()
		case "6":
			c.saveConfiguration()
		case "7":
			c.showCurrentConfig()
		case "b", "back":
			return
		case "q", "quit", "exit":
			fmt.Println("\n👋 Shutting down server...")
			s <- syscall.SIGQUIT
			return
		default:
			fmt.Println("❌ Invalid option, please try again.")
		}
		fmt.Println()
	}
}

func (c *CLI) printConfigMenu() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                     CONFIGURATION MENU                         ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  [1] 👤 FTP Credentials        %30s║\n", "Set FTP username/password")
	fmt.Printf("║  [2] 🔐 HTTP Authentication    %30s║\n", "Set HTTP basic auth")
	fmt.Printf("║  [3] 📁 FTP Server             %30s║\n", boolToStr(config.AppConfig.EnableFTP))
	fmt.Printf("║  [4] 🤖 MCP Server             %30s║\n", boolToStr(config.AppConfig.EnableMCP))
	fmt.Printf("║  [5] 🔌 Server Ports           %30s║\n", "Configure ports")
	fmt.Printf("║  [6] 💾 Save Configuration     %30s║\n", "Save to YAML file")
	fmt.Printf("║  [7] 👁️  View Current Config   %30s║\n", "Show all settings")
	fmt.Printf("║  [b] ↩️  Back to Main Menu                                       ║\n")
	fmt.Printf("║  [q] 🚪 Quit Server                                            ║\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
	fmt.Printf("\nEnter command: ")
}

func (c *CLI) configureFTPCredentials() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│              FTP Credentials Configuration                   │\n")
	fmt.Printf("├──────────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│  Current Username: %-41s║\n", config.AppConfig.Username)
	fmt.Printf("│  Current Password: %-41s║\n", maskPassword(config.AppConfig.Password))
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")

	fmt.Print("\nEnter new FTP username (press Enter to keep current): ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username != "" {
		config.AppConfig.Username = username
		fmt.Printf("✅ FTP username updated to: %s\n", username)
	}

	fmt.Print("Enter new FTP password (press Enter to keep current): ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	if password != "" {
		config.AppConfig.Password = password
		fmt.Println("✅ FTP password updated.")
	}

	fmt.Println("\n💡 Note: FTP credentials are now in memory.")
	fmt.Println("   Use 'Save Configuration' to persist changes to file.")
}

func (c *CLI) configureHTTPAuth() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│            HTTP Authentication Configuration                 │\n")
	fmt.Printf("├──────────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│  Status: %-51s║\n", boolToStr(config.AppConfig.HTTPAuth.Enabled))
	if config.AppConfig.HTTPAuth.Enabled {
		fmt.Printf("│  Username: %-49s║\n", config.AppConfig.HTTPAuth.Username)
		fmt.Printf("│  Password: %-49s║\n", maskPassword(config.AppConfig.HTTPAuth.Password))
	}
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")

	fmt.Print("\nEnable HTTP authentication? (y/n): ")
	enableStr, _ := reader.ReadString('\n')
	enableStr = strings.ToLower(strings.TrimSpace(enableStr))

	if enableStr == "y" || enableStr == "yes" {
		config.AppConfig.HTTPAuth.Enabled = true

		fmt.Print("Enter HTTP auth username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username != "" {
			config.AppConfig.HTTPAuth.Username = username
		}

		fmt.Print("Enter HTTP auth password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		if password != "" {
			config.AppConfig.HTTPAuth.Password = password
		}

		fmt.Println("✅ HTTP authentication enabled and configured.")
	} else if enableStr == "n" || enableStr == "no" {
		config.AppConfig.HTTPAuth.Enabled = false
		fmt.Println("✅ HTTP authentication disabled.")
	}
}

func (c *CLI) toggleFTPServer() {
	fmt.Printf("\n┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│                    FTP Server Toggle                          │\n")
	fmt.Printf("├──────────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│  Current Status: %-43s║\n", boolToStr(config.AppConfig.EnableFTP))
	fmt.Printf("│                                                              ║\n")
	fmt.Printf("│  ⚠️  Note: Changes require server restart to take effect     ║\n")
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nToggle FTP server? (y/n): ")
	toggleStr, _ := reader.ReadString('\n')
	toggleStr = strings.ToLower(strings.TrimSpace(toggleStr))

	if toggleStr == "y" || toggleStr == "yes" {
		config.AppConfig.EnableFTP = !config.AppConfig.EnableFTP
		fmt.Printf("✅ FTP server is now: %s\n", boolToStr(config.AppConfig.EnableFTP))
		if config.AppConfig.EnableFTP {
			fmt.Println("   Restart server to start FTP service.")
		} else {
			fmt.Println("   FTP service will stop after restart.")
		}
	}
}

func (c *CLI) toggleMCPServer() {
	fmt.Printf("\n┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│                    MCP Server Toggle                          │\n")
	fmt.Printf("├──────────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│  Current Status: %-43s║\n", boolToStr(config.AppConfig.EnableMCP))
	fmt.Printf("│                                                              ║\n")
	fmt.Printf("│  ℹ️  MCP enables AI assistants (like Claude) to interact     ║\n")
	fmt.Printf("│     with files through the Model Context Protocol            ║\n")
	fmt.Printf("│  ⚠️  Note: Changes require server restart to take effect     ║\n")
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nToggle MCP server? (y/n): ")
	toggleStr, _ := reader.ReadString('\n')
	toggleStr = strings.ToLower(strings.TrimSpace(toggleStr))

	if toggleStr == "y" || toggleStr == "yes" {
		config.AppConfig.EnableMCP = !config.AppConfig.EnableMCP
		fmt.Printf("✅ MCP server is now: %s\n", boolToStr(config.AppConfig.EnableMCP))
		if config.AppConfig.EnableMCP {
			fmt.Println("   Restart with -enable-mcp flag to start MCP service.")
		}
	}
}

func (c *CLI) configurePorts() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│                   Port Configuration                          │\n")
	fmt.Printf("├──────────────────────────────────────────────────────────────┤\n")
	fmt.Printf("│  Current HTTP Port: %-40s║\n", config.AppConfig.Port[1:])
	fmt.Printf("│  Current FTP Port:  %-40s║\n", config.AppConfig.FTPPort[1:])
	fmt.Printf("│                                                              ║\n")
	fmt.Printf("│  ⚠️  Note: Changes require server restart to take effect     ║\n")
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")

	fmt.Print("\nEnter new HTTP port (press Enter to keep current): ")
	httpPort, _ := reader.ReadString('\n')
	httpPort = strings.TrimSpace(httpPort)
	if httpPort != "" {
		config.AppConfig.Port = ":" + httpPort
		fmt.Printf("✅ HTTP port updated to: %s\n", httpPort)
	}

	fmt.Print("Enter new FTP port (press Enter to keep current): ")
	ftpPort, _ := reader.ReadString('\n')
	ftpPort = strings.TrimSpace(ftpPort)
	if ftpPort != "" {
		config.AppConfig.FTPPort = ":" + ftpPort
		fmt.Printf("✅ FTP port updated to: %s\n", ftpPort)
	}
}

func (c *CLI) saveConfiguration() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│                  Save Configuration                           │\n")
	fmt.Printf("├──────────────────────────────────────────────────────────────┤\n")

	currentPath := config.GetConfigPath()
	if currentPath == "defaults" {
		currentPath = config.GetDefaultConfigPath()
	}
	fmt.Printf("│  Config Path: %-47s║\n", truncateString(currentPath, 47))
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")

	fmt.Printf("\nSave to %s? (y/n): ", currentPath)
	saveStr, _ := reader.ReadString('\n')
	saveStr = strings.ToLower(strings.TrimSpace(saveStr))

	if saveStr == "y" || saveStr == "yes" {
		if err := config.SaveConfig(); err != nil {
			fmt.Printf("❌ Failed to save configuration: %v\n", err)
		} else {
			fmt.Printf("✅ Configuration saved to: %s\n", currentPath)
		}
	}
}

func (c *CLI) showCurrentConfig() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                   CURRENT CONFIGURATION                        ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  🌐 Server Settings:                                           ║\n")
	fmt.Printf("║     HTTP Port:      %-41s║\n", config.AppConfig.Port[1:])
	fmt.Printf("║     FTP Port:       %-41s║\n", config.AppConfig.FTPPort[1:])
	fmt.Printf("║     Root Directory: %-41s║\n", truncatePath(config.AppConfig.Root, 41))
	fmt.Printf("║     Auto Select IP: %-41s║\n", boolToStr(config.AppConfig.AutoSelect))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📁 FTP Server:                                                ║\n")
	fmt.Printf("║     Enabled:        %-41s║\n", boolToStr(config.AppConfig.EnableFTP))
	fmt.Printf("║     Username:       %-41s║\n", config.AppConfig.Username)
	fmt.Printf("║     Password:       %-41s║\n", maskPassword(config.AppConfig.Password))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  🤖 MCP Server:                                                ║\n")
	fmt.Printf("║     Enabled:        %-41s║\n", boolToStr(config.AppConfig.EnableMCP))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  🔐 HTTP Authentication:                                       ║\n")
	fmt.Printf("║     Enabled:        %-41s║\n", boolToStr(config.AppConfig.HTTPAuth.Enabled))
	if config.AppConfig.HTTPAuth.Enabled {
		fmt.Printf("║     Username:       %-41s║\n", config.AppConfig.HTTPAuth.Username)
		fmt.Printf("║     Password:       %-41s║\n", maskPassword(config.AppConfig.HTTPAuth.Password))
	}
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📤 Upload Settings:                                           ║\n")
	fmt.Printf("║     Enabled:        %-41s║\n", boolToStr(config.AppConfig.Upload.Enabled))
	fmt.Printf("║     Max Size:       %-41s║\n", formatSize(config.AppConfig.Upload.MaxSize))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  📝 Logging:                                                   ║\n")
	fmt.Printf("║     Level:          %-41s║\n", config.AppConfig.Logging.Level)
	fmt.Printf("║     Format:         %-41s║\n", config.AppConfig.Logging.Format)
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║  💾 Config File: %-45s║\n", truncateString(config.GetConfigPath(), 45))
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

func maskPassword(password string) string {
	if len(password) == 0 {
		return ""
	}
	if len(password) <= 2 {
		return "****"
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}
