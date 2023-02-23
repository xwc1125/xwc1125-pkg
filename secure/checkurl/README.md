## checkurl

本函数类用于校验URL的合法性，用于防止跳转漏洞、SSRF漏洞。

### 简介

| 函数                                             | 说明            |
|------------------------------------------------|---------------|
| secapi.CheckUrl(urlToCheck, ... newSecOptions) | 根据规则检查URL的合法性 |
| secapi.NewSchemeConfigure(schemeArr)           | 构造URL协议白名单    |
| secapi.NewHostConfigure(hostMap)               | 构造URL主机域名白名单  |

提供三种模式，用法如下：

| 模式      | 介绍      |
|---------|---------|
| equal   | 域名全等匹配  |
| subhost | 根域名后缀匹配 |
| regex   | 正则匹配    |

入参介绍如下：

| 参数            | 介绍                                 | 样例                                                                                                                                |
|---------------|------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------|
| urlToCheck    | 待校验的URL（需含协议+域名），如传入相对路径，会直接返回校验通过 | http://www.xwc1125.com                                                                                                            |
| newSecOptions | 校验规则                               | schemeArr := []string{"http","https"}<br/>hostMap := &hostOptions{ CheckMod: "subhost",   <br/>RuleArr:  []string{"xwc1125.com"}} |

### 示例

**secapi.CheckUrl**

```go
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
```
