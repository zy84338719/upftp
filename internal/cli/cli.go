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
	lang     string
}

var translations = map[string]map[string]string{
	"en": {
		"tagline":              "// ai-first file sharing server",
		"http_server":          "http_server",
		"ftp_server":           "ftp_server",
		"shared_path":          "shared_path",
		"files_found":          "files_found",
		"features":             "// features",
		"web_interface":        "web_interface",
		"mcp_support":          "mcp_support",
		"file_upload":          "file_upload",
		"qr_access":            "qr_access",
		"upload_enabled":       "upload: enabled",
		"auth_enabled":         "auth: enabled",
		"commands":             "// commands",
		"search_files":         "search_files",
		"list_all_files":       "list_all_files",
		"download_examples":    "download_examples",
		"refresh_list":         "refresh_list",
		"server_status":        "server_status",
		"configuration":        "configuration",
		"ftp_info":             "ftp_info",
		"about":                "about",
		"version_info":         "version_info",
		"quit_server":          "quit_server",
		"enter_command":        "enter command",
		"lang_hint":            "// lang:",
		"lang_switch":          "press [l] to switch language",
		"shutting_down":        "shutting down server...",
		"stdin_closed":         "stdin closed, running in headless mode.",
		"headless_note":        "HTTP/FTP services continue. Send SIGTERM to stop.",
		"invalid_option":       "invalid option, please try again.",
		"refreshing":           "refreshing file list...",
		"refresh_done":         "found %d files in directory.",
		"search_prompt":        "enter search term (or press Enter to show all): ",
		"search_results":       "SEARCH RESULTS",
		"no_match":             "no matching files found.",
		"found_count":          "found %d files.",
		"all_files":            "ALL FILES",
		"total_files":          "total: %d files",
		"and_more":             "... and %d more files",
		"showing_first":        "... and more files (showing first 20)",
		"no_files_example":     "no files available for download examples.",
		"dl_examples":          "DOWNLOAD EXAMPLES",
		"example_file":         "example_file",
		"browser":              "browser",
		"cli_tools":            "command line tools",
		"ftp_client":           "ftp client",
		"mcp_integration":      "mcp integration (claude desktop)",
		"version_info_box":     "VERSION INFORMATION",
		"version_label":        "version",
		"git_commit":           "git_commit",
		"build_date":           "build_date",
		"go_version":           "go_version",
		"platform":             "platform",
		"project_homepage":     "project_homepage",
		"server_status_box":    "SERVER STATUS",
		"server_config":        "server configuration",
		"http_port":            "http_port",
		"ftp_port":             "ftp_port",
		"root_directory":       "root_directory",
		"features_status":      "features status",
		"http_auth":            "http_auth",
		"username":             "username",
		"password":             "password",
		"file_upload_feat":     "file_upload",
		"ftp_server_feat":      "ftp_server",
		"https":                "https",
		"statistics":           "statistics",
		"files_available":      "files_available",
		"upload_max_size":      "upload_max_size",
		"about_box":            "ABOUT",
		"about_desc":           "is an ai-first lightweight file sharing server.",
		"key_features":         "key features",
		"feat_web_ui":          "modern responsive web ui with file preview",
		"feat_mcp":             "mcp (model context protocol) for ai integration",
		"feat_upload":          "file upload with drag & drop support",
		"feat_http_auth":       "http basic authentication",
		"feat_ftp":             "full ftp server implementation",
		"feat_qr":              "qr code for easy mobile access",
		"feat_tree":            "directory tree navigation",
		"feat_lang":            "multi-language support (en/zh)",
		"feat_yaml":            "yaml configuration file support",
		"feat_https":           "https support with custom certificates",
		"feat_cli":             "interactive cli configuration",
		"links":                "links",
		"license":              "license: mit",
		"build_info":           "build info",
		"config_menu":          "CONFIGURATION MENU",
		"ftp_credentials_menu": "ftp credentials",
		"http_auth_menu":       "http authentication",
		"ftp_server_menu":      "ftp server",
		"mcp_server_menu":      "mcp server",
		"server_ports_menu":    "server ports",
		"save_config_menu":     "save configuration",
		"view_config_menu":     "view current config",
		"back_main":            "back to main menu",
		"quit":                 "quit server",
		"current_config":       "CURRENT CONFIGURATION",
		"server_settings":      "server settings",
		"auto_select_ip":       "auto_select_ip",
		"mcp_server":           "mcp_server",
		"http_authentication":  "http_authentication",
		"upload_settings":      "upload settings",
		"enabled":              "enabled",
		"disabled":             "disabled",
		"logging":              "logging",
		"level":                "level",
		"format":               "format",
		"config_file":          "config_file",
		"note_save":            "use 'save configuration' to persist changes to file.",
		"ftp_connection_info":  "FTP CONNECTION INFO",
		"restart_required":     "note: changes require server restart to take effect",
	},
	"zh": {
		"tagline":              "// AI 时代的文件共享服务器",
		"http_server":          "HTTP 服务",
		"ftp_server":           "FTP 服务",
		"shared_path":          "共享目录",
		"files_found":          "文件数量",
		"features":             "// 功能特性",
		"web_interface":        "Web 界面",
		"mcp_support":          "MCP 支持",
		"file_upload":          "文件上传",
		"qr_access":            "二维码访问",
		"upload_enabled":       "上传: 已启用",
		"auth_enabled":         "认证: 已启用",
		"commands":             "// 命令菜单",
		"search_files":         "搜索文件",
		"list_all_files":       "列出文件",
		"download_examples":    "下载示例",
		"refresh_list":         "刷新列表",
		"server_status":        "服务器状态",
		"configuration":        "配置管理",
		"ftp_info":             "FTP 信息",
		"about":                "关于",
		"version_info":         "版本信息",
		"quit_server":          "退出服务器",
		"enter_command":        "请输入命令",
		"lang_hint":            "// 语言:",
		"lang_switch":          "按 [l] 切换语言",
		"shutting_down":        "正在关闭服务器...",
		"stdin_closed":         "标准输入已关闭，进入无头模式。",
		"headless_note":        "HTTP/FTP 服务继续运行。发送 SIGTERM 以停止。",
		"invalid_option":       "无效选项，请重试。",
		"refreshing":           "正在刷新文件列表...",
		"refresh_done":         "在目录中找到 %d 个文件。",
		"search_prompt":        "输入搜索关键词 (按 Enter 显示全部): ",
		"search_results":       "搜索结果",
		"no_match":             "未找到匹配的文件。",
		"found_count":          "找到 %d 个文件。",
		"all_files":            "所有文件",
		"total_files":          "共 %d 个文件",
		"and_more":             "... 还有 %d 个文件",
		"showing_first":        "... 更多文件 (显示前 20 个)",
		"no_files_example":     "没有可用于下载示例的文件。",
		"dl_examples":          "下载示例",
		"example_file":         "示例文件",
		"browser":              "浏览器",
		"cli_tools":            "命令行工具",
		"ftp_client":           "FTP 客户端",
		"mcp_integration":      "MCP 集成 (Claude Desktop)",
		"version_info_box":     "版本信息",
		"version_label":        "版本",
		"git_commit":           "Git 提交",
		"build_date":           "构建日期",
		"go_version":           "Go 版本",
		"platform":             "平台",
		"project_homepage":     "项目主页",
		"server_status_box":    "服务器状态",
		"server_config":        "服务器配置",
		"http_port":            "HTTP 端口",
		"ftp_port":             "FTP 端口",
		"root_directory":       "根目录",
		"features_status":      "功能状态",
		"http_auth":            "HTTP 认证",
		"username":             "用户名",
		"password":             "密码",
		"file_upload_feat":     "文件上传",
		"ftp_server_feat":      "FTP 服务",
		"https":                "HTTPS",
		"statistics":           "统计信息",
		"files_available":      "可用文件",
		"upload_max_size":      "上传大小限制",
		"about_box":            "关于",
		"about_desc":           "是一个 AI 时代的轻量级文件共享服务器。",
		"key_features":         "核心功能",
		"feat_web_ui":          "现代化响应式 Web 界面，支持文件预览",
		"feat_mcp":             "MCP (模型上下文协议) 用于 AI 集成",
		"feat_upload":          "文件上传，支持拖拽",
		"feat_http_auth":       "HTTP 基础认证",
		"feat_ftp":             "完整的 FTP 服务器实现",
		"feat_qr":              "二维码方便移动端访问",
		"feat_tree":            "目录树导航",
		"feat_lang":            "多语言支持 (EN/中文)",
		"feat_yaml":            "YAML 配置文件支持",
		"feat_https":           "HTTPS 自定义证书支持",
		"feat_cli":             "交互式 CLI 配置",
		"links":                "链接",
		"license":              "许可证: MIT",
		"build_info":           "构建信息",
		"config_menu":          "配置菜单",
		"ftp_credentials_menu": "FTP 凭据",
		"http_auth_menu":       "HTTP 认证",
		"ftp_server_menu":      "FTP 服务器",
		"mcp_server_menu":      "MCP 服务器",
		"server_ports_menu":    "服务器端口",
		"save_config_menu":     "保存配置",
		"view_config_menu":     "查看当前配置",
		"back_main":            "返回主菜单",
		"quit":                 "退出服务器",
		"current_config":       "当前配置",
		"server_settings":      "服务器设置",
		"auto_select_ip":       "自动选择 IP",
		"mcp_server":           "MCP 服务器",
		"http_authentication":  "HTTP 认证",
		"upload_settings":      "上传设置",
		"enabled":              "已启用",
		"disabled":             "已禁用",
		"logging":              "日志",
		"level":                "级别",
		"format":               "格式",
		"config_file":          "配置文件",
		"note_save":            "使用「保存配置」将更改持久化到文件。",
		"ftp_connection_info":  "FTP 连接信息",
		"restart_required":     "注意: 更改需要重启服务器才能生效",
	},
}

func t(lang, key string) string {
	if m, ok := translations[lang]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	if m, ok := translations["en"]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return key
}

func NewCLI() *CLI {
	lang := "en"
	home, _ := os.UserHomeDir()
	if home != "" {
		if data, err := ioutil.ReadFile(home + "/.upftp_lang"); err == nil {
			l := strings.TrimSpace(string(data))
			if l == "en" || l == "zh" {
				lang = l
			}
		}
	}
	if lang == "en" {
		if browserLang := os.Getenv("LANG"); strings.HasPrefix(browserLang, "zh") {
			lang = "zh"
		}
	}
	return &CLI{
		fileMap: make(map[string]string),
		lang:    lang,
	}
}

func (c *CLI) SetServerIP(ip string) {
	c.serverIP = ip
}

func (c *CLI) toggleLanguage() {
	if c.lang == "en" {
		c.lang = "zh"
	} else {
		c.lang = "en"
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		_ = ioutil.WriteFile(home+"/.upftp_lang", []byte(c.lang), 0644)
	}
}

func (c *CLI) ScanDirectory(pathname string) map[string]string {
	files := make(map[string]string)
	c.scanDir(pathname, files)
	return files
}

func (c *CLI) scanDir(pathname string, m map[string]string) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
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
					c.serverIP, fi.Name())
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
	lang := c.lang
	G := "\033[32m"
	Y := "\033[33m"
	Bd := "\033[1m"
	D := "\033[2m"
	W := "\033[37m"
	RS := "\033[0m"

	fmt.Printf("\n")
	fmt.Printf("  %s~%s %supftp%s %s%s%s  %sv%s%s\n",
		G+Bd, RS, Bd+W, RS, D, t(lang, "tagline"), RS,
		G, config.AppConfig.Version, RS)
	fmt.Printf("  %s%s: %s%s%s\n",
		D, t(lang, "build"), RS,
		config.AppConfig.LastCommit, D+" @ "+config.AppConfig.BuildDate+RS)
	fmt.Printf("\n")
	fmt.Printf("  %s────────────────────────────────────────────────────────────────%s\n", D, RS)
	fmt.Printf("  %s %-14s %shttp://%s%s\n", D, t(lang, "http_server")+":", G, c.serverIP+config.AppConfig.Port, RS)
	if config.AppConfig.EnableFTP {
		fmt.Printf("  %s %-14s %sftp://%s%s\n", D, t(lang, "ftp_server")+":", G, c.serverIP+config.AppConfig.FTPPort, RS)
	}
	fmt.Printf("  %s %-14s %s%s\n", D, t(lang, "shared_path")+":", RS, truncatePath(config.AppConfig.Root, 47))
	fmt.Printf("  %s %-14s %s%d%s\n", D, t(lang, "files_found")+":", W, len(c.fileMap), RS)
	fmt.Printf("  %s────────────────────────────────────────────────────────────────%s\n", D, RS)
	fmt.Printf("\n")
	fmt.Printf("  %s%s%s\n", D, t(lang, "features"), RS)
	fmt.Printf("    %s●%s %-14s %s●%s %-12s %s●%s %-12s %s●%s %s\n",
		G, RS, t(lang, "web_interface"),
		G, RS, t(lang, "mcp_support"),
		G, RS, t(lang, "file_upload"),
		G, RS, t(lang, "qr_access"))
	if config.AppConfig.Upload.Enabled {
		fmt.Printf("    %s●%s %-14s  %s●%s %s\n",
			Y, RS, t(lang, "upload_enabled"),
			Y, RS, t(lang, "auth_enabled"))
	}
	fmt.Printf("\n")
	fmt.Printf("  %s────────────────────────────────────────────────────────────────%s\n", D, RS)
	enLabel := fmt.Sprintf("%sEN%s", D, RS)
	zhLabel := fmt.Sprintf("%s中文%s", D, RS)
	if lang == "en" {
		enLabel = fmt.Sprintf("%s[%sEN%s]%s", G, Bd, G, RS)
	} else {
		zhLabel = fmt.Sprintf("%s[%s中文%s]%s", G, Bd, G, RS)
	}
	fmt.Printf("  %s %s %s  %s  %s%s\n",
		D+t(lang, "lang_hint")+RS, enLabel, zhLabel,
		D, t(lang, "lang_switch"), RS)
	fmt.Printf("  %s────────────────────────────────────────────────────────────────%s\n", D, RS)
	fmt.Printf("\n")
	c.printMenu()
}

func (c *CLI) printMenu() {
	lang := c.lang
	G := "\033[32m"
	R := "\033[31m"
	B := "\033[34m"
	Bd := "\033[1m"
	D := "\033[2m"
	RS := "\033[0m"

	fmt.Printf("  %s%s%s\n", D, t(lang, "commands"), RS)
	fmt.Printf("  ┌─────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │  %s[1]%s %-18s  %s[2]%s %-18s  %s[3]%s %-16s  │\n",
		G, RS, t(lang, "search_files"),
		G, RS, t(lang, "list_all_files"),
		G, RS, t(lang, "download_examples"))
	fmt.Printf("  │  %s[4]%s %-18s  %s[5]%s %-18s  %s[6]%s %-16s  │\n",
		G, RS, t(lang, "refresh_list"),
		G, RS, t(lang, "server_status"),
		G, RS, t(lang, "configuration"))
	if config.AppConfig.EnableFTP {
		fmt.Printf("  │  %s[7]%s %-18s  %s[8]%s %-18s                     │\n",
			G, RS, t(lang, "ftp_info"),
			G, RS, t(lang, "about"))
	} else {
		fmt.Printf("  │  %s[7]%s %-18s                                          │\n",
			G, RS, t(lang, "about"))
	}
	fmt.Printf("  │  %s[v]%s %-18s  %s[q]%s %-18s                     │\n",
		B, RS, t(lang, "version_info"),
		R, RS, t(lang, "quit_server"))
	fmt.Printf("  │  %s[l]%s %-18s                                        │\n",
		B, RS, t(lang, "lang_switch"))
	fmt.Printf("  └─────────────────────────────────────────────────────────────────┘\n")
	fmt.Printf("\n  %s>%s ", G+Bd, RS)
}

func (c *CLI) runLoop(s chan os.Signal) {
	for {
		if !c.handleInput(s) {
			fmt.Printf("\n  %s%s%s\n", "\033[2m", t(c.lang, "stdin_closed"), "\033[0m")
			fmt.Printf("  %s%s%s\n", "\033[2m", t(c.lang, "headless_note"), "\033[0m")
			return
		}
		fmt.Println()
	}
}

func (c *CLI) handleInput(s chan os.Signal) bool {
	var option string
	n, err := fmt.Scanln(&option)
	if err != nil || n == 0 {
		if err != nil {
			return false
		}
	}

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
		c.showServerStatus()
	case "6":
		c.showConfigMenu(s)
	case "7":
		if config.AppConfig.EnableFTP {
			c.showFTPInfo()
		} else {
			c.showAbout()
		}
	case "8":
		c.showAbout()
	case "v", "version":
		c.showVersion()
	case "l", "lang":
		c.toggleLanguage()
		c.printBanner()
	case "q", "quit", "exit":
		fmt.Printf("\n  %s%s%s\n", "\033[2m", t(c.lang, "shutting_down"), "\033[0m")
		s <- syscall.SIGQUIT
	default:
		fmt.Printf("  %s%s%s\n", "\033[31m", t(c.lang, "invalid_option"), "\033[0m")
	}
	return true
}

func (c *CLI) searchFiles() {
	lang := c.lang
	G := "\033[32m"
	D := "\033[2m"
	R := "\033[31m"
	RS := "\033[0m"

	fmt.Printf("  %s>%s %s", G, RS, t(lang, "search_prompt"))
	reader := bufio.NewReader(os.Stdin)
	term, _ := reader.ReadString('\n')
	term = strings.TrimSpace(term)

	found := false
	count := 0
	fmt.Printf("\n  %s%s%s\n", D, t(lang, "search_results"), RS)
	fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")

	for filename, url := range c.fileMap {
		if term == "" || strings.Contains(strings.ToLower(filename), strings.ToLower(term)) {
			found = true
			count++
			fmt.Printf("  │  %s●%s %-58s│\n", G, RS, truncateString(filename, 56))
			fmt.Printf("  │  %s%-58s│\n", D+truncateString(url, 58)+RS)
			fmt.Printf("  ├────────────────────────────────────────────────────────────────┤\n")
			if count >= 20 {
				fmt.Printf("  │  %s%-53s│\n", D+t(lang, "showing_first")+RS)
				break
			}
		}
	}

	if !found {
		fmt.Printf("  │  %s%-52s│\n", R+t(lang, "no_match")+RS)
	}
	fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
	fmt.Printf("\n  %s%s%s\n", D, fmt.Sprintf(t(lang, "found_count"), count), RS)
}

func (c *CLI) listAllFiles() {
	lang := c.lang
	G := "\033[32m"
	D := "\033[2m"
	RS := "\033[0m"

	count := 0
	fmt.Printf("\n  %s%s%s\n", D, t(lang, "all_files"), RS)
	fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")

	for filename, url := range c.fileMap {
		count++
		fmt.Printf("  │  %s●%s %-58s│\n", G, RS, truncateString(filename, 56))
		fmt.Printf("  │  %s%-58s│\n", D+truncateString(url, 58)+RS)
		fmt.Printf("  ├────────────────────────────────────────────────────────────────┤\n")
		if count >= 20 {
			fmt.Printf("  │  %s%-53s│\n", D+fmt.Sprintf(t(lang, "and_more"), len(c.fileMap)-20)+RS)
			break
		}
	}
	fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
	fmt.Printf("\n  %s%s%s\n", D, fmt.Sprintf(t(lang, "total_files"), len(c.fileMap)), RS)
}

func (c *CLI) showDownloadExamples() {
	lang := c.lang
	G := "\033[32m"
	D := "\033[2m"
	RS := "\033[0m"

	if len(c.fileMap) == 0 {
		fmt.Printf("  %s%s%s\n", "\033[31m", t(lang, "no_files_example"), RS)
		return
	}

	var exampleFile, exampleURL string
	for filename, url := range c.fileMap {
		exampleFile = filename
		exampleURL = url
		break
	}

	fmt.Printf("\n  %s%s%s\n", D, t(lang, "dl_examples"), RS)
	fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s %-14s %s%-40s│\n", D, t(lang, "example_file")+":", RS, truncateString(exampleFile, 40))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s\n", D, t(lang, "browser")+":"+RS)
	fmt.Printf("  │    %shttp://%s%s\n", G, c.serverIP+config.AppConfig.Port, RS)
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s\n", D, t(lang, "cli_tools")+":"+RS)
	fmt.Printf("  │    curl -O \"%s\n", truncateString(exampleURL, 52))
	fmt.Printf("  │    wget \"%s\n", truncateString(exampleURL, 54))
	fmt.Printf("  │                                                                │\n")
	if config.AppConfig.EnableFTP {
		fmt.Printf("  │  %s%s\n", D, t(lang, "ftp_client")+":"+RS)
		fmt.Printf("  │    ftp %s\n", c.serverIP)
		fmt.Printf("  │    Username: %s\n", config.AppConfig.Username)
		fmt.Printf("  │    Password: %s\n", config.AppConfig.Password)
		fmt.Printf("  │                                                                │\n")
	}
	fmt.Printf("  │  %s%s\n", D, t(lang, "mcp_integration")+":"+RS)
	fmt.Printf("  │    \"command\": \"upftp\",                                        │\n")
	fmt.Printf("  │    \"args\": [\"-enable-mcp\"]                                    │\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
}

func (c *CLI) showFTPInfo() {
	lang := c.lang
	D := "\033[2m"
	RS := "\033[0m"

	fmt.Printf("\n  %s%s%s\n", D, t(lang, "ftp_connection_info"), RS)
	fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s %-14s %s%-40s│\n", D, t(lang, "ftp_server")+":", RS, c.serverIP)
	fmt.Printf("  │  %s %-14s %-40s│\n", D, t(lang, "ftp_port")+":", config.AppConfig.FTPPort[1:])
	fmt.Printf("  │  %s %-14s %-40s│\n", D, t(lang, "username")+":", config.AppConfig.Username)
	fmt.Printf("  │  %s %-14s %-40s│\n", D, t(lang, "password")+":", config.AppConfig.Password)
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  $ ftp %s\n", c.serverIP)
	fmt.Printf("  │  Name: %s\n", config.AppConfig.Username)
	fmt.Printf("  │  Password: %s\n", config.AppConfig.Password)
	fmt.Printf("  │  ftp> ls                                                       │\n")
	fmt.Printf("  │  ftp> get filename                                             │\n")
	fmt.Printf("  │  ftp> put localfile                                            │\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
}

func (c *CLI) showVersion() {
	lang := c.lang
	G := "\033[32m"
	D := "\033[2m"
	RS := "\033[0m"

	fmt.Printf("\n  %s%s%s\n", D, t(lang, "version_info_box"), RS)
	fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s %-14s %-40s│\n", D, t(lang, "version_label")+":", config.AppConfig.Version)
	fmt.Printf("  │  %s %-14s %-40s│\n", D, t(lang, "git_commit")+":", config.AppConfig.LastCommit)
	fmt.Printf("  │  %s %-14s %-40s│\n", D, t(lang, "build_date")+":", config.AppConfig.BuildDate)
	fmt.Printf("  │  %s %-14s %-40s│\n", D, t(lang, "go_version")+":", config.AppConfig.GoVersion)
	fmt.Printf("  │  %s %-14s %-40s│\n", D, t(lang, "platform")+":", config.AppConfig.Platform)
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s %-14s %s%-40s│\n", D, t(lang, "project_homepage")+":", G, truncateString(config.AppConfig.ProjectURL, 40))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
}

func (c *CLI) showServerStatus() {
	lang := c.lang
	D := "\033[2m"
	RS := "\033[0m"

	fmt.Printf("\n  %s%s%s\n", D, t(lang, "server_status_box"), RS)
	fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "server_config")+":"+RS)
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "http_port")+":", config.AppConfig.Port[1:])
	if config.AppConfig.EnableFTP {
		fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "ftp_port")+":", config.AppConfig.FTPPort[1:])
	}
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "root_directory")+":", truncatePath(config.AppConfig.Root, 36))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "features_status")+":"+RS)
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "http_auth")+":", boolToStr(config.AppConfig.HTTPAuth.Enabled))
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "file_upload_feat")+":", boolToStr(config.AppConfig.Upload.Enabled))
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "ftp_server_feat")+":", boolToStr(config.AppConfig.EnableFTP))
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "https")+":", boolToStr(config.AppConfig.HTTPS.Enabled))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "statistics")+":"+RS)
	fmt.Printf("  │    %s %-14s %-36d│\n", D, t(lang, "files_available")+":", len(c.fileMap))
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "upload_max_size")+":", formatSize(config.AppConfig.Upload.MaxSize))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
}

func (c *CLI) showAbout() {
	lang := c.lang
	G := "\033[32m"
	D := "\033[2m"
	RS := "\033[0m"

	fmt.Printf("\n  %s %s v%s %s\n", D, t(lang, "about_box"), config.AppConfig.Version, RS)
	fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s %s %s\n", G, config.AppConfig.ProjectName, t(lang, "about_desc"))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "key_features")+":"+RS)
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_web_ui"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_mcp"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_upload"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_http_auth"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_ftp"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_qr"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_tree"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_lang"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_yaml"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_https"))
	fmt.Printf("  │    ● %-56s│\n", t(lang, "feat_cli"))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "links")+":"+RS)
	fmt.Printf("  │    Homepage:  %-47s│\n", truncateString(config.AppConfig.ProjectURL, 47))
	fmt.Printf("  │    GitHub:    https://github.com/zy84338719/upftp             │\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s\n", D, t(lang, "license")+RS)
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "build_info")+":"+RS)
	fmt.Printf("  │    %-12s %-47s│\n", t(lang, "version_label")+" "+":", config.AppConfig.Version)
	fmt.Printf("  │    %-12s %-47s│\n", t(lang, "git_commit")+" "+":", config.AppConfig.LastCommit)
	fmt.Printf("  │    %-12s %-47s│\n", t(lang, "build_date")+" "+":", config.AppConfig.BuildDate)
	fmt.Printf("  │    %-12s %-47s│\n", t(lang, "platform")+" "+":", config.AppConfig.Platform)
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
}

func (c *CLI) refreshFileList() {
	lang := c.lang
	D := "\033[2m"
	G := "\033[32m"
	RS := "\033[0m"

	fmt.Printf("  %s%s%s\n", D, t(lang, "refreshing"), RS)
	c.fileMap = c.ScanDirectory(config.AppConfig.Root)
	fmt.Printf("  %s%s%s\n", G, fmt.Sprintf(t(lang, "refresh_done"), len(c.fileMap)), RS)
}

func (c *CLI) showConfigMenu(s chan os.Signal) {
	lang := c.lang
	G := "\033[32m"
	Y := "\033[33m"
	R := "\033[31m"
	Bd := "\033[1m"
	D := "\033[2m"
	RS := "\033[0m"

	for {
		fmt.Printf("\n  %s%s%s\n", D, t(lang, "config_menu"), RS)
		fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")
		fmt.Printf("  │                                                                │\n")
		fmt.Printf("  │  %s[1]%s %-22s %-34s│\n", G, RS, t(lang, "ftp_credentials_menu"), "")
		fmt.Printf("  │  %s[2]%s %-22s %-34s│\n", G, RS, t(lang, "http_auth_menu"), "")
		fmt.Printf("  │  %s[3]%s %-22s %s%-30s│\n", G, RS, t(lang, "ftp_server_menu"), Y, boolToStr(config.AppConfig.EnableFTP)+RS)
		fmt.Printf("  │  %s[4]%s %-22s %s%-30s│\n", G, RS, t(lang, "mcp_server_menu"), Y, boolToStr(config.AppConfig.EnableMCP)+RS)
		fmt.Printf("  │  %s[5]%s %-22s %-34s│\n", G, RS, t(lang, "server_ports_menu"), "")
		fmt.Printf("  │  %s[6]%s %-22s %-34s│\n", G, RS, t(lang, "save_config_menu"), "")
		fmt.Printf("  │  %s[7]%s %-22s %-34s│\n", G, RS, t(lang, "view_config_menu"), "")
		fmt.Printf("  │  %s[b]%s %-22s                                        │\n", G, RS, t(lang, "back_main"))
		fmt.Printf("  │  %s[q]%s %-22s                                        │\n", R, RS, t(lang, "quit"))
		fmt.Printf("  │                                                                │\n")
		fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
		fmt.Printf("\n  %s>%s ", G+Bd, RS)

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
			fmt.Printf("\n  %s%s%s\n", D, t(lang, "shutting_down"), RS)
			s <- syscall.SIGQUIT
			return
		default:
			fmt.Printf("  %s%s%s\n", R, t(lang, "invalid_option"), RS)
		}
		fmt.Println()
	}
}

func (c *CLI) configureFTPCredentials() {
	lang := c.lang
	D := "\033[2m"
	G := "\033[32m"
	RS := "\033[0m"
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n  %s FTP Credentials %s\n", D, RS)
	fmt.Printf("  ┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "username")+":", config.AppConfig.Username)
	fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "password")+":", maskPassword(config.AppConfig.Password))
	fmt.Printf("  └──────────────────────────────────────────────────────────────┘\n")

	fmt.Printf("\n  %s>%s ", G, RS)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if username != "" {
		config.AppConfig.Username = username
		fmt.Printf("  %s username → %s%s\n", G, username, RS)
	}

	fmt.Printf("  %s>%s ", G, RS)
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	if password != "" {
		config.AppConfig.Password = password
		fmt.Printf("  %s password updated%s\n", G, RS)
	}

	fmt.Printf("\n  %s%s%s\n", D, t(lang, "note_save"), RS)
}

func (c *CLI) configureHTTPAuth() {
	lang := c.lang
	D := "\033[2m"
	G := "\033[32m"
	RS := "\033[0m"
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n  %s HTTP Authentication %s\n", D, RS)
	fmt.Printf("  ┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "http_auth")+":", boolToStr(config.AppConfig.HTTPAuth.Enabled))
	if config.AppConfig.HTTPAuth.Enabled {
		fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "username")+":", config.AppConfig.HTTPAuth.Username)
		fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "password")+":", maskPassword(config.AppConfig.HTTPAuth.Password))
	}
	fmt.Printf("  └──────────────────────────────────────────────────────────────┘\n")

	fmt.Printf("\n  %s>%s (y/n) ", G, RS)
	enableStr, _ := reader.ReadString('\n')
	enableStr = strings.ToLower(strings.TrimSpace(enableStr))

	if enableStr == "y" || enableStr == "yes" {
		config.AppConfig.HTTPAuth.Enabled = true
		fmt.Printf("  %s>%s ", G, RS)
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username != "" {
			config.AppConfig.HTTPAuth.Username = username
		}
		fmt.Printf("  %s>%s ", G, RS)
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		if password != "" {
			config.AppConfig.HTTPAuth.Password = password
		}
		fmt.Printf("  %s http auth enabled%s\n", G, RS)
	} else if enableStr == "n" || enableStr == "no" {
		config.AppConfig.HTTPAuth.Enabled = false
		fmt.Printf("  %s http auth disabled%s\n", G, RS)
	}
}

func (c *CLI) toggleFTPServer() {
	lang := c.lang
	D := "\033[2m"
	G := "\033[32m"
	RS := "\033[0m"
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n  %s FTP Server %s\n", D, RS)
	fmt.Printf("  ┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "ftp_server_feat")+":", boolToStr(config.AppConfig.EnableFTP))
	fmt.Printf("  │  %s %s%s\n", D, t(lang, "restart_required"), RS)
	fmt.Printf("  └──────────────────────────────────────────────────────────────┘\n")

	fmt.Printf("\n  %s>%s (y/n) ", G, RS)
	toggleStr, _ := reader.ReadString('\n')
	toggleStr = strings.ToLower(strings.TrimSpace(toggleStr))

	if toggleStr == "y" || toggleStr == "yes" {
		config.AppConfig.EnableFTP = !config.AppConfig.EnableFTP
		fmt.Printf("  %s ftp server → %s%s\n", G, boolToStr(config.AppConfig.EnableFTP), RS)
	}
}

func (c *CLI) toggleMCPServer() {
	lang := c.lang
	D := "\033[2m"
	G := "\033[32m"
	RS := "\033[0m"
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n  %s MCP Server %s\n", D, RS)
	fmt.Printf("  ┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "mcp_server")+":", boolToStr(config.AppConfig.EnableMCP))
	fmt.Printf("  │  %s %s%s\n", D, t(lang, "restart_required"), RS)
	fmt.Printf("  └──────────────────────────────────────────────────────────────┘\n")

	fmt.Printf("\n  %s>%s (y/n) ", G, RS)
	toggleStr, _ := reader.ReadString('\n')
	toggleStr = strings.ToLower(strings.TrimSpace(toggleStr))

	if toggleStr == "y" || toggleStr == "yes" {
		config.AppConfig.EnableMCP = !config.AppConfig.EnableMCP
		fmt.Printf("  %s mcp server → %s%s\n", G, boolToStr(config.AppConfig.EnableMCP), RS)
	}
}

func (c *CLI) configurePorts() {
	lang := c.lang
	D := "\033[2m"
	G := "\033[32m"
	RS := "\033[0m"
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n  %s Port Configuration %s\n", D, RS)
	fmt.Printf("  ┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "http_port")+":", config.AppConfig.Port[1:])
	fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "ftp_port")+":", config.AppConfig.FTPPort[1:])
	fmt.Printf("  │  %s %s%s\n", D, t(lang, "restart_required"), RS)
	fmt.Printf("  └──────────────────────────────────────────────────────────────┘\n")

	fmt.Printf("\n  %s>%s ", G, RS)
	httpPort, _ := reader.ReadString('\n')
	httpPort = strings.TrimSpace(httpPort)
	if httpPort != "" {
		config.AppConfig.Port = ":" + httpPort
		fmt.Printf("  %s http port → %s%s\n", G, httpPort, RS)
	}

	fmt.Printf("  %s>%s ", G, RS)
	ftpPort, _ := reader.ReadString('\n')
	ftpPort = strings.TrimSpace(ftpPort)
	if ftpPort != "" {
		config.AppConfig.FTPPort = ":" + ftpPort
		fmt.Printf("  %s ftp port → %s%s\n", G, ftpPort, RS)
	}
}

func (c *CLI) saveConfiguration() {
	lang := c.lang
	D := "\033[2m"
	G := "\033[32m"
	R := "\033[31m"
	RS := "\033[0m"
	reader := bufio.NewReader(os.Stdin)

	currentPath := config.GetConfigPath()
	if currentPath == "defaults" {
		currentPath = config.GetDefaultConfigPath()
	}

	fmt.Printf("\n  %s Save Configuration %s\n", D, RS)
	fmt.Printf("  ┌──────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │  %s %-14s %-41s│\n", D, t(lang, "config_file")+":", truncateString(currentPath, 41))
	fmt.Printf("  └──────────────────────────────────────────────────────────────┘\n")

	fmt.Printf("\n  %s>%s %s? (y/n) ", G, RS, currentPath)
	saveStr, _ := reader.ReadString('\n')
	saveStr = strings.ToLower(strings.TrimSpace(saveStr))

	if saveStr == "y" || saveStr == "yes" {
		if err := config.SaveConfig(); err != nil {
			fmt.Printf("  %s failed: %v%s\n", R, err, RS)
		} else {
			fmt.Printf("  %s saved → %s%s\n", G, currentPath, RS)
		}
	}
}

func (c *CLI) showCurrentConfig() {
	lang := c.lang
	D := "\033[2m"
	RS := "\033[0m"

	fmt.Printf("\n  %s%s%s\n", D, t(lang, "current_config"), RS)
	fmt.Printf("  ┌────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "server_settings")+":"+RS)
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "http_port")+":", config.AppConfig.Port[1:])
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "ftp_port")+":", config.AppConfig.FTPPort[1:])
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "root_directory")+":", truncatePath(config.AppConfig.Root, 36))
	fmt.Printf("  │    %s %-14s %-36s│\n", D, t(lang, "auto_select_ip")+":", boolToStr(config.AppConfig.AutoSelect))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "ftp_server")+":"+RS)
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "enabled")+" "+":", boolToStr(config.AppConfig.EnableFTP))
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "username")+" "+":", config.AppConfig.Username)
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "password")+" "+":", maskPassword(config.AppConfig.Password))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "mcp_server")+":"+RS)
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "enabled")+" "+":", boolToStr(config.AppConfig.EnableMCP))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "http_authentication")+":"+RS)
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "enabled")+" "+":", boolToStr(config.AppConfig.HTTPAuth.Enabled))
	if config.AppConfig.HTTPAuth.Enabled {
		fmt.Printf("  │    %-16s %-36s│\n", t(lang, "username")+" "+":", config.AppConfig.HTTPAuth.Username)
		fmt.Printf("  │    %-16s %-36s│\n", t(lang, "password")+" "+":", maskPassword(config.AppConfig.HTTPAuth.Password))
	}
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "upload_settings")+":"+RS)
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "enabled")+" "+":", boolToStr(config.AppConfig.Upload.Enabled))
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "upload_max_size")+" "+":", formatSize(config.AppConfig.Upload.MaxSize))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s%s%s\n", D, t(lang, "logging")+":"+RS)
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "level")+" "+":", config.AppConfig.Logging.Level)
	fmt.Printf("  │    %-16s %-36s│\n", t(lang, "format")+" "+":", config.AppConfig.Logging.Format)
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  │  %s %-14s %-36s│\n", D, t(lang, "config_file")+":", truncateString(config.GetConfigPath(), 36))
	fmt.Printf("  │                                                                │\n")
	fmt.Printf("  └────────────────────────────────────────────────────────────────┘\n")
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
		return "\033[32m● enabled\033[0m"
	}
	return "\033[31m● disabled\033[0m"
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

func maskPassword(password string) string {
	if len(password) == 0 {
		return ""
	}
	if len(password) <= 2 {
		return "****"
	}
	return password[:2] + strings.Repeat("*", len(password)-2)
}
