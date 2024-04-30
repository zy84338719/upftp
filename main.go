package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall" //nolint:gci
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
	p := flag.String("p", "10000", "-p 8888 prot default 10000")
	dir := flag.String("d", "./", "-d /opt dir default ./")
	flag.Parse()
	port = ":" + *p
	root = *dir

	// ip addr
	ips := getIps()
	for {
		ipnum := 0
		fmt.Println("Enter select ip number")
		_, _ = fmt.Scanln(&ipnum)
		if ipnum >= 0 && ipnum < len(ips) {
			ip = "http://" + ips[ipnum]
			break
		}
		fmt.Println("Enter select ip number err")
	}
}

func main() {
	files := map[string]string{}
	getAllFile(root, files)
	ctx := context.Background()
	ctx, cancelFunc := context.WithCancel(ctx)
	go GinServer(ctx)
	sigChan := make(chan os.Signal, 1)
	go scanCmd(ctx, files, sigChan)
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
