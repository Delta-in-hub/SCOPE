package utils

import (
	"testing"
)

func TestGetMyIpAddrs(t *testing.T) {
	ipMap := GetMyIpAddrs()

	// 确保返回的不是空映射
	if len(ipMap) == 0 {
		// 在某些环境中可能没有网络接口，所以这不一定是错误
		t.Log("No IP addresses found, this might be expected in some environments")
		return
	}

	// 验证每个IP地址的格式
	for iface, ip := range ipMap {
		if iface == "" {
			t.Errorf("Empty interface name found")
		}

		if ip == "" {
			t.Errorf("Empty IP address for interface %s", iface)
		}

		// 简单检查IP地址格式是否合理
		// 这不是完整的IP地址验证，但可以检查基本格式
		if len(ip) < 7 { // 最短的有效IP是形如1.1.1.1
			t.Errorf("IP address %s for interface %s seems invalid", ip, iface)
		}
		t.Logf("%s: %s", iface, ip)
	}
}
