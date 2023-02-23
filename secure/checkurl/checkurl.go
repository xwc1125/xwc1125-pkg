package secapi

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type HostOptions struct {
	// 支持equal - 代表域名全匹配，subhost - 域名后缀匹配，regex - 代表使用正则匹配；默认是equal模式
	CheckMod string
	// 支持传入多个完整域名及正则表达式
	RuleArr []string
}

type secOptions struct {
	SchemeArr []string
	HostMap   *HostOptions
}

// defaultSecOptions 默认URL检查配置
var defaultSecOptions = secOptions{
	SchemeArr: []string{"http", "https"},
	HostMap: &HostOptions{
		CheckMod: "equal",
		RuleArr:  []string{"www.xwc1125.com"},
	},
}

type newSecOptions func(*secOptions)

// NewSchemeConfigure 添加协议白名单
func NewSchemeConfigure(schemeArr []string) newSecOptions {
	return func(o *secOptions) {
		o.SchemeArr = schemeArr
	}
}

// NewHostConfigure 添加域名白名单
func NewHostConfigure(hostMap *HostOptions) newSecOptions {
	return func(o *secOptions) {
		o.HostMap = hostMap
	}
}

// CheckUrl 校验URL的合法性
func CheckUrl(urlToCheck string, opts ...newSecOptions) (bool, error) {
	secOptions := defaultSecOptions
	isSchemeValid := false
	isHostValid := false
	for _, o := range opts {
		o(&secOptions)
	}
	hostRes, err := url.Parse(urlToCheck)
	if err != nil {
		return false, errors.New("invalid url")
	}

	schemeToCheck := hostRes.Scheme
	hostToCheck := strings.Split(hostRes.Host, ":")[0]

	if schemeToCheck == "" && hostToCheck == "" {
		return true, nil
	}
	for _, schemeValue := range secOptions.SchemeArr {
		if schemeToCheck == schemeValue || schemeToCheck == "" {
			isSchemeValid = true
			break
		}
	}

	if secOptions.HostMap.CheckMod == "equal" {
		for _, hostValue := range secOptions.HostMap.RuleArr {
			if hostToCheck == hostValue || hostToCheck == "" {
				isHostValid = true
				break
			}
		}
	} else if secOptions.HostMap.CheckMod == "regex" {
		for _, regexValue := range secOptions.HostMap.RuleArr {
			reg := regexp.MustCompile(regexValue)
			if reg.MatchString(hostToCheck) {
				isHostValid = true
				break
			}
		}
	} else if secOptions.HostMap.CheckMod == "subhost" {
		for _, subHostValue := range secOptions.HostMap.RuleArr {
			domainWithDot := fmt.Sprintf(".%s", subHostValue)
			if strings.HasSuffix(hostToCheck, domainWithDot) || hostToCheck == subHostValue {
				isHostValid = true
				break
			}
		}
	}

	if !isSchemeValid {
		return false, errors.New("invalid scheme")
	}
	if !isHostValid {
		return false, errors.New("invalid host")
	}
	return true, nil
}
