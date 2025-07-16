package logic

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"

	"github.com/zy84338719/upftp/config"
)

var fileMap map[string]string
var serverIP string

func SetServerIP(ip string) {
	serverIP = ip
}

func ScanDirectory(pathname string) map[string]string {
	files := make(map[string]string)
	scanDir(pathname, files)
	return files
}

func scanDir(pathname string, m map[string]string) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", pathname, err)
		return
	}

	for _, fi := range rd {
		if fi.IsDir() {
			scanDir(path.Join(pathname, fi.Name()), m)
		} else {
			dir := strings.Replace(pathname, config.AppConfig.Root, "", 1)
			var downloadURL string
			if len(dir) > 0 {
				dir = path.Join(dir)
				downloadURL = fmt.Sprintf("http://%s%s/download/%s/%s", 
					serverIP, config.AppConfig.Port, dir, fi.Name())
			} else {
				downloadURL = fmt.Sprintf("http://%s%s/download/%s", 
					serverIP, config.AppConfig.Port, fi.Name())
			}
			m[fi.Name()] = downloadURL
		}
	}
}

func StartCommandInterface(ctx context.Context, s chan os.Signal) {
	fileMap = ScanDirectory(config.AppConfig.Root)
	
	if len(fileMap) == 0 {
		fmt.Println("No files found in the directory")
		s <- syscall.SIGQUIT
		return
	}

	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                          UPFTP SERVER                         ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Web Interface: http://%-39s ║\n", serverIP + config.AppConfig.Port)
	if config.AppConfig.EnableFTP {
		fmt.Printf("║ FTP Server:    ftp://%-40s ║\n", serverIP + config.AppConfig.FTPPort)
		fmt.Printf("║ FTP Login:     %s / %s%-26s ║\n", 
			config.AppConfig.Username, 
			config.AppConfig.Password,
			strings.Repeat(" ", 26-len(config.AppConfig.Username)-len(config.AppConfig.Password)))
	}
	fmt.Printf("║ Root Path:     %-47s ║\n", config.AppConfig.Root)
	fmt.Printf("║ Files Found:   %-47d ║\n", len(fileMap))
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n\n")

	for {
		fmt.Printf("Commands:\n")
		fmt.Printf("  [1] Search files\n")
		fmt.Printf("  [2] List all files\n")
		fmt.Printf("  [3] Show download examples\n")
		fmt.Printf("  [4] Refresh file list\n")
		if config.AppConfig.EnableFTP {
			fmt.Printf("  [5] FTP connection info\n")
		}
		fmt.Printf("  [q] Quit server\n")
		fmt.Printf("\nEnter command: ")

		var option string
		_, _ = fmt.Scanln(&option)

		switch strings.ToLower(option) {
		case "1":
			searchFiles()
		case "2":
			listAllFiles()
		case "3":
			showDownloadExamples()
		case "4":
			refreshFileList()
		case "5":
			if config.AppConfig.EnableFTP {
				showFTPInfo()
			} else {
				fmt.Println("Invalid option. FTP is not enabled.")
			}
		case "q", "quit", "exit":
			fmt.Println("\nShutting down server...")
			s <- syscall.SIGQUIT
			return
		default:
			fmt.Println("Invalid option, please try again.")
		}
		fmt.Println()
	}
}

func searchFiles() {
	fmt.Print("Enter search term (or press Enter to show all): ")
	reader := bufio.NewReader(os.Stdin)
	term, _ := reader.ReadString('\n')
	term = strings.TrimSpace(term)

	found := false
	fmt.Printf("\n%-50s %-60s\n", "FILE NAME", "DOWNLOAD URL")
	fmt.Println(strings.Repeat("=", 115))
	
	for filename, url := range fileMap {
		if term == "" || strings.Contains(strings.ToLower(filename), strings.ToLower(term)) {
			found = true
			fmt.Printf("%-50s %-60s\n", truncateString(filename, 48), url)
		}
	}

	if !found {
		fmt.Println("No matching files found.")
	}
}

func listAllFiles() {
	fmt.Printf("\n%-50s %-60s\n", "FILE NAME", "DOWNLOAD URL")
	fmt.Println(strings.Repeat("=", 115))
	
	for filename, url := range fileMap {
		fmt.Printf("%-50s %-60s\n", truncateString(filename, 48), url)
	}
}

func showDownloadExamples() {
	if len(fileMap) == 0 {
		fmt.Println("No files available for download examples.")
		return
	}

	// 获取第一个文件作为示例
	var exampleFile, exampleURL string
	for filename, url := range fileMap {
		exampleFile = filename
		exampleURL = url
		break
	}

	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                        DOWNLOAD EXAMPLES                      ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Example file: %-47s ║\n", truncateString(exampleFile, 45))
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Browser:                                                       ║\n")
	fmt.Printf("║   Open: http://%-46s ║\n", serverIP + config.AppConfig.Port)
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║ Command Line Tools:                                            ║\n")
	fmt.Printf("║   curl -O \"%-53s\" ║\n", truncateString(exampleURL, 51))
	fmt.Printf("║   wget \"%-55s\" ║\n", truncateString(exampleURL, 53))
	fmt.Printf("║                                                                ║\n")
	if config.AppConfig.EnableFTP {
		fmt.Printf("║ FTP Client:                                                    ║\n")
		fmt.Printf("║   ftp %-56s ║\n", serverIP)
		fmt.Printf("║   Username: %-50s ║\n", config.AppConfig.Username)
		fmt.Printf("║   Password: %-50s ║\n", config.AppConfig.Password)
		fmt.Printf("║                                                                ║\n")
	}
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

func showFTPInfo() {
	fmt.Printf("\n╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                        FTP CONNECTION INFO                    ║\n")
	fmt.Printf("╠════════════════════════════════════════════════════════════════╣\n")
	fmt.Printf("║ Server:    %-51s ║\n", serverIP)
	fmt.Printf("║ Port:      %-51s ║\n", config.AppConfig.FTPPort[1:])
	fmt.Printf("║ Username:  %-51s ║\n", config.AppConfig.Username)
	fmt.Printf("║ Password:  %-51s ║\n", config.AppConfig.Password)
	fmt.Printf("║                                                                ║\n")
	fmt.Printf("║ Example FTP commands:                                          ║\n")
	fmt.Printf("║   ftp %s                                              ║\n", serverIP)
	fmt.Printf("║   > Name: %s                                          ║\n", config.AppConfig.Username)
	fmt.Printf("║   > Password: %s                                      ║\n", config.AppConfig.Password)
	fmt.Printf("║   > ls                                                         ║\n")
	fmt.Printf("║   > get filename                                               ║\n")
	fmt.Printf("║   > put filename                                               ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n")
}

func refreshFileList() {
	fmt.Println("Refreshing file list...")
	fileMap = ScanDirectory(config.AppConfig.Root)
	fmt.Printf("Found %d files in directory.\n", len(fileMap))
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
