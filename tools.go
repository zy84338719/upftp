package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
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
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(i, " ", ipnet.IP.String())
				ips = append(ips, ipnet.IP.String())
				i++
			}
		}
	}
	return ips
}

func GinServer(ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 你的自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s \n",
			param.ClientIP,
			param.TimeStamp.Format("3:04:05PM"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	router.StaticFS("/", gin.Dir(root, true))
	if err := router.Run(port); err != nil {
		panic(fmt.Errorf("Ftp Server start error = %s", err))
	}
}


func getAllFile(pathname string) map[string]string {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println(fmt.Errorf("file menu ready error，err=%s", err))
		return map[string]string{}
	}
	m := map[string]string{}
	for _, fi := range rd {
		if fi.IsDir() {
			//fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			for k, v := range getAllFile(path.Join(pathname, fi.Name())) {
				m[k] = v
			}
		} else {
			dir := strings.Replace(pathname, root, "", 1)
			if len(dir) > 0 {
				dir = path.Join(dir)
			}
			m[fi.Name()] = ip + port + dir + "/" + fi.Name()
			fmt.Println(ip + port + dir + "/" + fi.Name())
		}
	}
	return m
}

func scanCmd(ctx context.Context, files map[string]string, s chan os.Signal) {
	if len(files) == 0 {
		s <- syscall.SIGQUIT
		return
	}
	var data string
	for {
		fmt.Println("place enter file name, exit and q is kill me")
		_,_=fmt.Scanln(&data)
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
	s <- syscall.SIGQUIT
}