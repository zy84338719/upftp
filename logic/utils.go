package logic

import (
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
	var data string
	for {
		fmt.Println("place index " + IP + Port + " in browser to view files")
		fmt.Println("place enter file name, exit and q is kill me")
		_, _ = fmt.Scanln(&data)
		if data == "exit" || data == "q" {
			break
		}
		for k, v := range files {
			if strings.Contains(k, data) {
				fmt.Println(v)
			}
		}
		fmt.Println()
		fmt.Println()
	}
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
