package utils

import (
	"fmt"
	"net"
)

func GetMyIpAddrs() map[string]string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		return map[string]string{}
	}

	ipMap := make(map[string]string)
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP
			// 排除回环地址和无效地址
			if !ip.IsLoopback() && ip.IsGlobalUnicast() {
				// fmt.Printf("%-8s %s\n", iface.Name+":", ip)
				ipMap[iface.Name] = ip.String()
			}
		}
	}
	return ipMap
}
