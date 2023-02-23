package safecurl

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Err 安全API校验错误值
var (
	errSsrfAttack    = fmt.Errorf("ssrf attack attempt")
	errInvalidDoamin = fmt.Errorf("invalid domain")
	errInvalidURL    = fmt.Errorf("invalid url")
)

// CheckIp 检查IP地址是否合法
func CheckIp(ipStr string) bool {
	address := net.ParseIP(ipStr)
	if address == nil {
		return false
	} else {
		return true
	}
}

// InetAton IP转int64
func InetAton(ipStr string) int64 {
	bits := strings.Split(ipStr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

// InetNtoa int64 to IP
func InetNtoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// IsInnerIp 检查是否内网IP地址
func IsInnerIp(ipStr string, innerHosts ...string) bool {
	if !CheckIp(ipStr) {
		return false
	}
	inputIpNum := InetAton(ipStr)
	innerIpA := InetAton("10.255.255.255")
	innerIpB := InetAton("172.16.255.255")
	innerIpC := InetAton("192.168.255.255")
	innerIpD := InetAton("100.64.255.255")
	innerIpE := InetAton("127.255.255.255")
	if inputIpNum>>24 == innerIpA>>24 || inputIpNum>>20 == innerIpB>>20 ||
		inputIpNum>>16 == innerIpC>>16 || inputIpNum>>22 == innerIpD>>22 ||
		inputIpNum>>24 == innerIpE>>24 {
		return true
	}
	for _, addr := range innerHosts {
		innerIpF := InetAton(addr)
		if inputIpNum>>24 == innerIpF>>24 {
			return true
		}
	}
	return false
}

// GetIpByDomain 获取域名指向IP
func GetIpByDomain(hostname string) (string, error) {
	addr, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	} else {
		return addr[0].String(), nil
	}
}

// GetSafeURL 校验URL的安全性
func GetSafeURL(inputUrl string, innerHosts ...string) (r string, errorReturn error) {
	var digIp, hostname string
	s := inputUrl

	u, err := url.Parse(s)
	if err != nil {
		errorReturn = errInvalidURL
		return
	}

	h := strings.Split(u.Host, ":")
	hostname = h[0]

	if CheckIp(hostname) {
		digIp = hostname
		if IsInnerIp(digIp, innerHosts...) {
			errorReturn = errSsrfAttack
			return
		}
	} else {
		digIp, err = GetIpByDomain(hostname)
		if err != nil {
			errorReturn = errInvalidDoamin
			return
		}
		if IsInnerIp(digIp, innerHosts...) {
			errorReturn = errSsrfAttack
			return
		}
	}

	return digIp, errorReturn
}
