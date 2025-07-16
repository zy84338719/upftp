package network

import (
	"fmt"
	"net"
	"strconv"
)

type NetworkInfo struct {
	IP   string
	Port int
}

func GetAvailableIPs() []string {
	fmt.Println("Available network interfaces:")
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return nil
	}

	var ips []string
	index := 0
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			fmt.Printf("[%d] %s\n", index, ipnet.IP.String())
			ips = append(ips, ipnet.IP.String())
			index++
		}
	}
	return ips
}

func SelectIP(ips []string, autoSelect bool) string {
	if len(ips) == 0 {
		return ""
	}

	if autoSelect {
		return ips[0]
	}

	for {
		fmt.Print("\nSelect IP number (or press Enter for first IP): ")
		var input string
		_, _ = fmt.Scanln(&input)

		if input == "" {
			return ips[0]
		}

		if ipnum, err := strconv.Atoi(input); err == nil && ipnum >= 0 && ipnum < len(ips) {
			return ips[ipnum]
		}
		fmt.Println("Invalid selection, please try again")
	}
}

func GetNetworkInfo(autoSelect bool, httpPort, ftpPort int) (string, error) {
	ips := GetAvailableIPs()
	if len(ips) == 0 {
		return "", fmt.Errorf("no available network interfaces found")
	}

	selectedIP := SelectIP(ips, autoSelect)
	return selectedIP, nil
}
