package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"syscall"
)

func getIps() []string {
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

func GinServer(ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	router := gin.New()
	router.LoadHTMLGlob("templates/*")             // 加载模板文件
	router.StaticFS("/files", gin.Dir(root, true)) // 静态文件服务

	router.GET("/", func(c *gin.Context) {
		files, _ := ioutil.ReadDir(root)
		fileList := []map[string]interface{}{}
		for _, file := range files {
			fileMap := map[string]interface{}{
				"Name":    file.Name(),
				"IsImage": strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".jpg"),
			}
			fileList = append(fileList, fileMap)
		}
		c.HTML(http.StatusOK, "index.html", fileList)
	})

	router.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.FileAttachment(path.Join(root, filename), filename)
	})

	router.GET("/preview/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.File(path.Join(root, filename))
	})

	if err := router.Run(port); err != nil {
		panic(fmt.Errorf("Server start error = %s", err))
	}
}

func getAllFile(pathname string, m map[string]string) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println(fmt.Errorf("file menu ready error，err=%s", err))
		return
	}
	for i, fi := range rd {
		if fi.IsDir() {
			getAllFile(path.Join(pathname, fi.Name()), m)
		} else {
			dir := strings.Replace(pathname, root, "", 1)
			if len(dir) > 0 {
				dir = path.Join(dir)
				m[fi.Name()] = ip + path.Join(port, "download", dir, fi.Name())
				fmt.Println(ip + path.Join(port, "download", dir, fi.Name()))
			} else {
				m[fi.Name()] = ip + path.Join(port, "download", fi.Name())
				fmt.Println(ip + path.Join(port, "download", fi.Name()))
			}
		}
		if i > 20 {
			fmt.Println("file list too long, only show 20")
			break
		}
	}
}

func scanCmd(ctx context.Context, files map[string]string, s chan os.Signal) {
	if len(files) == 0 {
		s <- syscall.SIGQUIT
		return
	}
	var data string
	for {
		fmt.Println("place index " + ip + port + " in browser to view files")
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
