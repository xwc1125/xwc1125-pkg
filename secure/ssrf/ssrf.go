// Package ssrf
//
// @author: xwc1125
// @date: 2021/4/28
package ssrf

import (
	"net"
	"net/http"
	"net/url"

	"github.com/chain5j/logger"
	"github.com/xwc1125/xwc1125-pkg/secure/safecurl"
)

// IsSSRF 判断是否为ssrf
func IsSSRF(urlPath string, headers map[string]string, innerHosts ...string) bool {
	uri, err := url.Parse(urlPath)
	if err != nil {
		logger.Error("url.Parse err", "url", urlPath, "err", err)
		return true
	}
	ips, err := net.LookupIP(uri.Hostname())
	if err != nil {
		logger.Error("LookupIP err", "uri.Hostname()", uri.Hostname(), "err", err)
		return true
	}
	if ips == nil || len(ips) == 0 {
		logger.Error("ips empty")
		return true
	}
	for _, ip := range ips {
		ipStr := ip.String()
		if safecurl.IsInnerIp(ipStr) {
			logger.Error("[ssrf]url is innerIp", "url", urlPath, "ip", ipStr)
			return true
		}

	}
	safeClient := safecurl.NewSafeClient()
	httpReqInValid, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		logger.Error("new request err", "url", urlPath, "err", err)
		return true
	}
	for key, value := range headers {
		httpReqInValid.Header.Add(key, value)
	}
	respInValid, err := safeClient.Do(httpReqInValid)
	if respInValid != nil {
		logger.Debug("[ssrf] safeClient respond", "resp", respInValid)
		if respInValid.StatusCode != 200 && respInValid.StatusCode != 201 {
			logger.Error("[ssrf]respInValid", "url", urlPath)
			return true
		}
	}
	return false
}
