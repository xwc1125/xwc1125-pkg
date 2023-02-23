// Package iputil
//
// @author: xwc1125
package iputil

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/xwc1125/xwc1125-pkg/utils/iputil/ip2region"
)

// GetIpLocation 获取IP归属地
func GetIpLocation(ip string) *ip2region.IpInfo {
	ipInfo, err := ip2region.GetIpInfo(ip)
	if err != nil {
		return nil
	}
	return &ipInfo
}

// GetIntranetIp 获取内网IP
func GetIntranetIp(prefix string) []string {
	ips := make([]string, 0)
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ips
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipAddr := ipNet.IP.String()
				if prefix == "" {
					ips = append(ips, ipAddr)
				} else {
					if strings.HasPrefix(ipAddr, prefix) {
						ips = append(ips, ipAddr)
					}
				}
			}
		}
	}
	return ips
}

// GetLocalhost 获取局域网ip地址
func GetLocalhost() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						return ipNet.IP.String()
					}
				}
			}
		}

	}
	return ""
}

// GetClientIP 获取客户端的IP
func GetClientIP(request *http.Request) string {
	clientIP := request.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = request.Header.Get("X-real-ip")
	}
	ip := request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = clientIP
	}

	remoteIP := GetRemoteIP(request)
	if ip == "127.0.0.1" || ip == "" {
		ip = remoteIP
	}
	return ip
}

func GetRemoteIP(request *http.Request) string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(request.RemoteAddr))
	if err != nil {
		return ""
	}
	remoteIP := net.ParseIP(ip)
	if remoteIP == nil {
		return ""
	}

	return remoteIP.String()
}
