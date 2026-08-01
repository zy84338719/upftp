// Package netinfo 检测本机可用的局域网 IP,用于拼接对外展示的访问地址。
package netinfo

import (
	"fmt"
	"net"
)

// AvailableIPs 返回所有非 loopback 的 IPv4 地址。
func AvailableIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips, nil
}

// BestIP 返回最可能可用的 IP:优先取第一个局域网地址,无则回退到 127.0.0.1。
func BestIP() string {
	ips, err := AvailableIPs()
	if err != nil || len(ips) == 0 {
		return "127.0.0.1"
	}
	return ips[0]
}

// AllHosts 返回用于展示的地址清单(供 TUI/API 打印多网卡情况)。
func AllHosts() []string {
	ips, err := AvailableIPs()
	if err != nil || len(ips) == 0 {
		return []string{"127.0.0.1"}
	}
	out := make([]string, len(ips))
	for i, ip := range ips {
		out[i] = fmt.Sprintf("%s", ip)
	}
	return out
}
