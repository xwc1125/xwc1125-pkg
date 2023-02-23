// Package pmodel
//
// @author: xwc1125
package pmodel

import "github.com/chain5j/chain5j-pkg/codec/json"

type AppInfo struct {
	L        string `json:"l"`        // app名称
	C        int    `json:"c"`        // app版本号
	V        string `json:"v"`        // app版本名称
	Pk       string `json:"pk"`       // app包名
	Apk      string `json:"apk"`      // app密钥
	Cpk      string `json:"cpk"`      // company密钥;
	Platform string `json:"platform"` // 平台
	Sign     string `json:"sign"`     // 包签名内容
	R        int64  `json:"r"`        // 随机码;
}

func (a AppInfo) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
