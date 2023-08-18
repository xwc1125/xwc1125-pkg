// Package email
package email

import (
	"bytes"
	"html/template"
	"testing"
	"time"

	"github.com/chain5j/logger"
	"github.com/jordan-wright/email"
)

var Template = `
<div style="background-color:white;border-top:2px solid #12ADDB;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
    <h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;">
        <span style="color: #12ADDB;font-weight:bold;">
            {{.Title}}
        </span>
    </h2>
    <div style="padding:0 12px 0 12px; margin-top:18px;">
        {{if .Content}}
		<p>
            {{.Content}}
        </p>
		{{end}}
		{{if .QuoteContent}}
		<div style="background-color: #f5f5f5;padding: 10px 15px;margin:18px 0;word-wrap:break-word;">
            {{.QuoteContent}}
        </div>
		{{end}}
       
		{{if .Url}}
        <p>
            <a style="text-decoration:none; color:#12addb" href="{{.Url}}" target="_blank" rel="noopener">点击查看详情</a>
        </p>
		{{end}}
    </div>
</div>
`

// GetTemplateEmail 获取模版邮件
func GetTemplateEmail(title, content, quoteContent, url string) string {
	tpl, err := template.New("emailTemplate").Parse(Template)
	if err != nil {
		logger.Error("邮件模版", "err", err)
		return ""
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, map[string]interface{}{
		"Title":        title,
		"Content":      content,
		"QuoteContent": quoteContent,
		"Url":          url,
	})
	if err != nil {
		logger.Error("邮件发送", "err", err)
		return ""
	}
	html := b.String()
	return html
}

func TestEmail(t *testing.T) {
	config := Config{
		Addr:     "smtp.qq.com",
		Port:     "25",
		Username: "",
		Password: "",
	}
	newEmail, err := NewEmail(config)
	if err != nil {
		panic(err)
	}
	html := GetTemplateEmail(
		"邮件测试",
		"内容",
		"111",
		"https://www.baidu.com")
	err = newEmail.Send(&email.Email{
		From:    config.Username,
		To:      []string{"admin@qq.com"},
		Subject: "主题",
		// Text:        nil,
		HTML: []byte(html),
		// Sender:      "",
		// Headers:     nil,
		// Attachments: nil,
		// ReadReceipt: nil,
	}, 10*time.Second)
	if err != nil {
		panic(err)
	}
}
