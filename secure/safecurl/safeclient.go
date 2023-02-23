package safecurl

import (
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

type Ipv4Net struct {
	Ipv4            [4]byte
	SubnetPrefixLen int
}

func ipv4Net(a, b, c, d byte, subnetPrefixLen int) net.IPNet {
	return net.IPNet{
		IP:   net.IPv4(a, b, c, d),
		Mask: net.CIDRMask(96+subnetPrefixLen, 128)}
}

var reservedIPv4Nets = []net.IPNet{
	ipv4Net(0, 0, 0, 0, 8),       // Current network
	ipv4Net(10, 0, 0, 0, 8),      // Private
	ipv4Net(100, 64, 0, 0, 10),   // RFC6598
	ipv4Net(127, 0, 0, 0, 8),     // Loopback
	ipv4Net(169, 254, 0, 0, 16),  // Link-local
	ipv4Net(172, 16, 0, 0, 12),   // Private
	ipv4Net(192, 0, 0, 0, 24),    // RFC6890
	ipv4Net(192, 0, 2, 0, 24),    // Test, doc, examples
	ipv4Net(192, 88, 99, 0, 24),  // IPv6 to IPv4 relay
	ipv4Net(192, 168, 0, 0, 16),  // Private
	ipv4Net(198, 18, 0, 0, 15),   // Benchmarking tests
	ipv4Net(198, 51, 100, 0, 24), // Test, doc, examples
	ipv4Net(203, 0, 113, 0, 24),  // Test, doc, examples
	ipv4Net(224, 0, 0, 0, 4),     // Multicast
	ipv4Net(240, 0, 0, 0, 4),     // Reserved (includes broadcast / 255.255.255.255)
}

var globalUnicastIPv6Net = net.IPNet{net.IP{0x20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, net.CIDRMask(3, 128)}

func isIPv6GlobalUnicast(address net.IP) bool {
	return globalUnicastIPv6Net.Contains(address)
}

func isIPv4Reserved(address net.IP, innerIpv4Nets ...Ipv4Net) bool {
	var reservedIPv4Nets1 = make([]net.IPNet, 0, len(innerIpv4Nets)+len(reservedIPv4Nets))
	copy(reservedIPv4Nets1, reservedIPv4Nets)
	for _, innerIpv4Net := range innerIpv4Nets {
		reservedIPv4Nets1 = append(reservedIPv4Nets1, ipv4Net(innerIpv4Net.Ipv4[0], innerIpv4Net.Ipv4[1], innerIpv4Net.Ipv4[2], innerIpv4Net.Ipv4[3], innerIpv4Net.SubnetPrefixLen))
	}
	for _, reservedNet := range reservedIPv4Nets1 {
		if reservedNet.Contains(address) {
			return true
		}
	}

	return false
}

func isPublicIPAddress(address net.IP, innerIpv4Nets ...Ipv4Net) bool {
	if address.To4() != nil {
		return !isIPv4Reserved(address, innerIpv4Nets...)
	} else {
		return isIPv6GlobalUnicast(address)
	}
}

func safeSocketControl(innerIpv4Nets ...Ipv4Net) func(network, address string, c syscall.RawConn) error {
	return func(network, address string, c syscall.RawConn) error {
		if !(network == "tcp4" || network == "tcp6") {
			return fmt.Errorf("%s is not a safe network type", network)
		}
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return fmt.Errorf("%s is not a valid host/port pair: %s", address, err)
		}

		ipaddress := net.ParseIP(host)
		if ipaddress == nil {
			return fmt.Errorf("%s is not a valid IP address", host)
		}
		if !isPublicIPAddress(ipaddress, innerIpv4Nets...) {
			return fmt.Errorf("%s is not a public IP address", ipaddress)
		}
		if !(port == "80" || port == "443") {
			return fmt.Errorf("%s is not a safe port number", port)
		}
		return nil
	}
}

// NewSafeClient 安全的HTTP请求客户端
func NewSafeClient(innerIpv4Nets ...Ipv4Net) *http.Client {
	safeDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
		Control:   safeSocketControl(innerIpv4Nets...),
	}
	safeTransport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           safeDialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	safeClient := &http.Client{
		Transport: safeTransport,
	}
	return safeClient
}

// SafeCurl 安全地发起HTTP GET请求
func SafeCurl(untrustedURL string, innerIpv4Nets ...Ipv4Net) *http.Response {
	safeClient := NewSafeClient(innerIpv4Nets...)
	resp, err := safeClient.Get(untrustedURL)
	if err != nil {
		return nil
	}
	return resp
}
