// Package response
package response

import (
	"bytes"
	"fmt"
)

type StatusErr struct {
	Status   int    `json:"status"`   // 状态值
	LangKey  string `json:"langKey"`  // 语言KeyCode
	Message  string `json:"msg"`      // 默认消息
	MoreInfo string `json:"moreInfo"` // 其他消息
}

func NewStatusErr(err error) *StatusErr {
	return &StatusErr{
		Status:   FailStatus.Status,
		LangKey:  "",
		Message:  err.Error(),
		MoreInfo: "",
	}
}
func (e *StatusErr) Code() int {
	return e.Status
}
func (e *StatusErr) Error() string {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%d", e.Status))
	if len(e.Message) > 0 {
		buff.WriteString(fmt.Sprintf(":%s", e.Message))
	} else if len(e.LangKey) > 0 {
		buff.WriteString(fmt.Sprintf(":%s", e.LangKey))
	}
	if len(e.MoreInfo) > 0 {
		buff.WriteString(fmt.Sprintf("->%s", e.MoreInfo))
	}
	return buff.String()
}

func (e *StatusErr) Msg() string {
	var buff bytes.Buffer

	if len(e.Message) > 0 {
		buff.WriteString(fmt.Sprintf("%s", e.Message))
	} else if len(e.LangKey) > 0 {
		buff.WriteString(fmt.Sprintf("%s", e.LangKey))
	}
	if len(e.MoreInfo) > 0 {
		buff.WriteString(fmt.Sprintf("->%s", e.MoreInfo))
	}
	return buff.String()
}
