package logic

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"syscall"
)

// 导出全局变量
var (
	Root string
	IP   string
	Port string
)

func GetIps() []string {
	fmt.Println("Your networks list")
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ips := []string{}
	i := 0
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); !ok || ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
			continue
		} else {
			fmt.Println(i, " ", ipnet.IP.String())
			ips = append(ips, ipnet.IP.String())
			i++
		}
	}
	return ips
}

func GetAllFile(pathname string, m map[string]string) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println(fmt.Errorf("file menu ready error，err=%s", err))
		return
	}
	for _, fi := range rd {
		if fi.IsDir() {
			GetAllFile(path.Join(pathname, fi.Name()), m)
		} else {
			dir := strings.Replace(pathname, Root, "", 1)
			if len(dir) > 0 {
				dir = path.Join(dir)
				m[fi.Name()] = IP + path.Join(Port, "download", dir, fi.Name())
				fmt.Println(IP + path.Join(Port, "download", dir, fi.Name()))
			} else {
				m[fi.Name()] = IP + path.Join(Port, "download", fi.Name())
				fmt.Println(IP + path.Join(Port, "download", fi.Name()))
			}
		}
	}
}

func ScanCmd(ctx context.Context, files map[string]string, s chan os.Signal) {
	if len(files) == 0 {
		s <- syscall.SIGQUIT
		return
	}

	fmt.Printf("\n=== upftp Server ===\n")
	fmt.Printf("Web Interface: %s%s\n", IP, Port)
	fmt.Printf("Root Directory: %s\n\n", Root)

	for {
		fmt.Printf("\nOptions:\n")
		fmt.Printf("1. Search files\n")
		fmt.Printf("2. Show download commands\n")
		fmt.Printf("3. Exit\n")
		fmt.Printf("Choose an option (1-3): ")

		var option string
		_, _ = fmt.Scanln(&option)

		switch option {
		case "1":
			searchFiles(files)
		case "2":
			showDownloadCommands()
		case "3", "exit", "q":
			fmt.Println("\nShutting down server...")
			s <- syscall.SIGQUIT
			return
		default:
			fmt.Println("\nInvalid option, please try again")
		}
	}
}

func searchFiles(files map[string]string) {
	fmt.Printf("\nEnter search term (or press Enter to list all): ")
	var term string
	reader := bufio.NewReader(os.Stdin)
	term, _ = reader.ReadString('\n')
	term = strings.TrimSpace(term)

	found := false
	fmt.Println("\nMatching files:")
	fmt.Println("---------------")
	for k, v := range files {
		if term == "" || strings.Contains(strings.ToLower(k), strings.ToLower(term)) {
			found = true
			fmt.Printf("\nFile: %s\n", k)
			fmt.Printf("Download URL: %s\n", v)
			fmt.Printf("curl command: curl -O %s\n", v)
			fmt.Printf("wget command: wget %s\n", v)
			fmt.Println("---------------")
		}
	}

	if !found {
		fmt.Println("No matching files found")
	}
}

func showDownloadCommands() {
	fmt.Printf("\nDownload Commands Examples:\n")
	fmt.Printf("----------------------------\n")
	fmt.Printf("Using curl:\n")
	fmt.Printf("  Single file: curl -O <url>\n")
	fmt.Printf("  Save as different name: curl -o newname.txt <url>\n")
	fmt.Printf("  Download with progress: curl -# -O <url>\n\n")

	fmt.Printf("Using wget:\n")
	fmt.Printf("  Single file: wget <url>\n")
	fmt.Printf("  Save as different name: wget -O newname.txt <url>\n")
	fmt.Printf("  Download with progress: wget -q --show-progress <url>\n\n")

	fmt.Printf("Note:\n")
	fmt.Printf("- Replace <url> with the actual download URL\n")
	fmt.Printf("- You can also use web browser to download files: %s%s\n", IP, Port)
	fmt.Printf("----------------------------\n")
}

// 添加一个辅助函数来打印分隔线
func printDivider() {
	fmt.Println("\n" + strings.Repeat("-", 50) + "\n")
}

func getFileType(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return "image"
	case ".txt", ".md", ".json", ".yaml", ".yml", ".go", ".js", ".html", ".css":
		return "text"
	default:
		return "binary"
	}
}

func canPreviewFile(filename string) bool {
	fileType := getFileType(filename)
	return fileType == "image" || fileType == "text"
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func isPathSafe(filename string) bool {
	cleanPath := path.Clean(filename)
	return !strings.Contains(cleanPath, "..") && !strings.HasPrefix(cleanPath, "/")
}
