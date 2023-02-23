// Package tcputil
//
// @author: xwc1125
package tcputil

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/xwc1125/xwc1125-pkg/utils/iputil"
)

// GetTcpInfo 解析获取TCP内容
// @request: *http.Request
// returns:
// #1: types.TcpInfo
func GetTcpInfo(request *http.Request) *TcpInfo {
	tcpInfo := &TcpInfo{}
	tcpInfo.Ip = iputil.GetClientIP(request)
	detail := iputil.GetIpLocation(tcpInfo.Ip)
	if detail != nil {
		tcpInfo.IpAddress = detail.String()
	}
	scheme := "http://"
	if request.TLS != nil {
		scheme = "https://"
	}
	tcpInfo.Url = scheme + request.Host + request.RequestURI
	tcpInfo.Api = request.URL.Path
	userAgent := request.Header.Get("User-Agent")
	tcpInfo.Ua = userAgent
	tcpInfo.Browser = GetBrowser(request)
	tcpInfo.Os = GetClientOs(request)
	tcpInfo.IsMobile = IsMobile(request)

	lang := request.Header.Get("lang")
	if lang == "" {
		lang = request.Header.Get("language")
	}
	tcpInfo.Lang = lang
	return tcpInfo
}

// GetBrowser 获取浏览器标示
func GetBrowser(req *http.Request) string {
	browserInfo := "other"
	userAgent := req.Header.Get("User-Agent")
	if userAgent == "" || len(strings.TrimSpace(userAgent)) == 0 {
		return browserInfo
	}
	ua := strings.ToLower(userAgent)

	msieP := "msie ([\\d.]+)"
	firefoxP := "firefox\\/([\\d.]+)"
	ieheighP := "rv:([\\d.]+)"
	chromeP := "chrome\\/([\\d.]+)"
	operaP := "opr.([\\d.]+)"
	safariP := "version\\/([\\d.]+).*safari"

	// 匹配
	matched, _ := regexp.MatchString(msieP, ua)
	if matched {
		reg := regexp.MustCompile(msieP)
		// fmt.Println(reg.NumSubexp())// 匹配的个数
		// s := reg.FindAllString(ua, -1)// 获取所有匹配内容
		s := reg.FindString(ua)
		version := strings.Split(s, " ")
		// s = mat.group()
		// version = s.split(" ")[1]
		// browserInfo = "ie "+ version.substring(0, version.indexOf("."))
		browserInfo = "ie " + version[1]
		return browserInfo
	}
	matched, _ = regexp.MatchString(firefoxP, ua)
	if matched {
		reg := regexp.MustCompile(firefoxP)
		s := reg.FindString(ua)
		version := strings.Split(s, "/")
		browserInfo = "firefox " + version[1]
		return browserInfo
	}
	matched, _ = regexp.MatchString(ieheighP, ua)
	if matched {
		reg := regexp.MustCompile(ieheighP)
		s := reg.FindString(ua)
		version := strings.Split(s, "：")
		browserInfo = "ie " + version[1]
		return browserInfo
	}
	matched, _ = regexp.MatchString(operaP, ua)
	if matched {
		reg := regexp.MustCompile(operaP)
		s := reg.FindString(ua)
		version := strings.Split(s, "/")
		browserInfo = "opera " + version[1]
		return browserInfo
	}
	matched, _ = regexp.MatchString(chromeP, ua)
	if matched {
		reg := regexp.MustCompile(chromeP)
		s := reg.FindString(ua)
		version := strings.Split(s, "/")
		browserInfo = "chrome " + version[1]
		return browserInfo
	}
	matched, _ = regexp.MatchString(safariP, ua)
	if matched {
		reg := regexp.MustCompile(safariP)
		s := reg.FindString(ua)
		version := strings.Split(s, " ")
		vStr := version[0]
		browserInfo = "safari " + vStr[strings.Index(vStr, "/")+1:]
		return browserInfo
	}
	return browserInfo
}

var (
	mobileDeviceMap map[string]string
	pcDeviceMap     map[string]string
)

func getKeywordsMap() {
	if mobileDeviceMap == nil || len(mobileDeviceMap) == 0 {
		mobileDeviceMap = map[string]string{
			"Android":        "Android",
			"ANDROID":        "ANDROID",
			"IOS":            "IOS",
			"iPhone":         "iPhone",
			"iPod":           "iPod",
			"iPad":           "iPad",
			"Windows Phone":  "Windows Phone",
			"MQQBrowser":     "MQQBrowser",
			"UCWEB":          "UCWEB",
			"UCBrowser":      "UCWEB",
			"MicroMessenger": "WeiXin",
			"Opera":          "Opera",
		}
	}
	if pcDeviceMap == nil || len(pcDeviceMap) == 0 {
		pcDeviceMap = map[string]string{
			// PC端
			// Windows
			"(WinNT|Windows NT)":                      "WinNT",
			"(Windows NT 6\\.2)":                      "Win 8",
			"(Windows NT 6\\.1)":                      "Win 7",
			"(Windows NT 5\\.1|Windows XP)":           "WinXP",
			"(Windows NT 5\\.2)":                      "Win2003",
			"(Win2000|Windows 2000|Windows NT 5\\.0)": "Win2000",

			"(9x 4.90|Win9(5|8)|Windows 9(5|8)|95/NT|Win32|32bit)": "Win9x",

			// mac
			"(Mac|apple|MacOS8)": "MAC",
			"(68k|68000)":        "Mac68k",

			// Linux
			"Linux": "Linux",
		}
	}
}

// GetClientOs 获取OS
func GetClientOs(req *http.Request) string {
	userAgent := req.Header.Get("User-Agent")
	getKeywordsMap()

	cos := "UNKNOWN"
	if userAgent == "" {
		return cos
	}
	for k, v := range mobileDeviceMap {
		matched, _ := regexp.MatchString(".*"+k+".*", userAgent)
		if matched {
			return v
		}
	}
	for k, v := range pcDeviceMap {
		matched, _ := regexp.MatchString(".*"+k+".*", userAgent)
		if matched {
			return v
		}
	}
	return cos
}

var isMobileRegex = regexp.MustCompile(`(?i)(android|avantgo|blackberry|bolt|boost|cricket|docomo|fone|hiptop|mini|mobi|palm|phone|pie|tablet|up\.browser|up\.link|webos|wos)`)

func IsMobile(req *http.Request) bool {
	s := req.Header.Get("User-Agent")
	return isMobileRegex.MatchString(s)
}
