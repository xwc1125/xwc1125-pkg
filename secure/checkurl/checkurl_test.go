// Package secapi
//
// @author: xwc1125
package secapi

import (
	"fmt"
	"testing"
)

func TestCheck(t *testing.T) {
	urlToCheck := "http://xwc1125.com@abc.com:80"

	// schemeArr 允许的URL协议
	schemeArr := []string{"http", "https"}
	// hostMap 允许的URL域名
	hostMap := &HostOptions{
		CheckMod: "subhost",               // 根域名后缀匹配模式，支持equal、subhost和regex三种模式
		RuleArr:  []string{"xwc1125.com"}, // 允许xwc1125.com及其子域名
	}

	schemeWhiteList := NewSchemeConfigure(schemeArr)
	hostWhiteList := NewHostConfigure(hostMap)

	result, err := CheckUrl(urlToCheck, schemeWhiteList, hostWhiteList)
	if err != nil {
		fmt.Print(err)
		fmt.Print(result)
	}
}
