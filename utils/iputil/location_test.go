// Package iputil
//
// @author: xwc1125
package iputil

import (
	"fmt"
	"testing"

	"github.com/xwc1125/xwc1125-pkg/utils/iputil/ip2region"
)

func TestLocation(t *testing.T) {
	ip, err := ip2region.GetIpInfo("127.0.0.1")
	fmt.Println(ip, err)
	ip, err = ip2region.GetIpInfo("127.0.0.1")
	fmt.Println(ip, err)
	ip, err = ip2region.GetIpInfo("61.148.16.170")
	fmt.Println(ip, err)
}

func TestIP(t *testing.T) {
	localhost := GetLocalhost()
	fmt.Println(localhost)
	ips := GetIntranetIp("192.")
	for i, ip := range ips {
		fmt.Println(i, ip)
	}
}
