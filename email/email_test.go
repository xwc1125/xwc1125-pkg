// Package email
//
// @author: xwc1125
package email

import "testing"

func TestEmail(t *testing.T) {
	SendTemplateEmail(EmailConfig{
		Addr:     "smtp.qq.com",
		Port:     "25",
		Username: "",
		Password: "",
	}, "xwc1125@qq.com",
		"主题",
		"邮件测试",
		"内容",
		"111",
		"https://www.baidu.com")
}
