package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/zy84338719/upftp/logic" // 替换为实际的导入路径
)

var (
	Version    = "undefined"
	LastCommit = "undefined"
)
var (
	root, ip, port string
)

func init() {
	fmt.Println("version:", Version)
	fmt.Println("last_commit:", LastCommit)
	p := flag.String("p", "10000", "-p 8888 port default 10000")
	dir := flag.String("d", "./", "-d /opt dir default ./")
	autoIP := flag.Bool("auto", false, "-auto to automatically select the first non-loopback IP")
	flag.Parse()
	port = ":" + *p
	root = *dir

	ips := logic.GetIps()
	if len(ips) == 0 {
		log.Fatal("No available network interfaces found")
	}

	if *autoIP {
		ip = "http://" + ips[0]
		return
	}

	fmt.Println("\nAvailable IP addresses:")
	for i, addr := range ips {
		fmt.Printf("[%d] %s\n", i, addr)
	}

	for {
		fmt.Print("\nSelect IP number (or press Enter for first IP): ")
		var input string
		_, _ = fmt.Scanln(&input)

		if input == "" {
			ip = "http://" + ips[0]
			break
		}

		if ipnum, err := strconv.Atoi(input); err == nil && ipnum >= 0 && ipnum < len(ips) {
			ip = "http://" + ips[ipnum]
			break
		}
		fmt.Println("Invalid selection, please try again")
	}

	logic.Root = root
	logic.IP = ip
	logic.Port = port
}

func main() {
	files := map[string]string{}
	logic.GetAllFile(logic.Root, files)
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)
	go logic.GinServer(ctx)
	sigChan := make(chan os.Signal, 1)
	go logic.ScanCmd(ctx, files, sigChan)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-sigChan
		log.Printf("get a signal %s\n", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancelFunc()
			log.Println("upftp server exit now...")
			return
		case syscall.SIGHUP:
		default:
		}
	}
}
