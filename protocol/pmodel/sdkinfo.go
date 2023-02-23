// Package pmodel
//
// @author: xwc1125
package pmodel

import "github.com/chain5j/chain5j-pkg/codec/json"

type SdkInfo struct {
	C  int    `json:"c"`  // sdk版本号
	N  string `json:"n"`  // sdk名称
	V  string `json:"v"`  // sdk版本名称
	Cm string `json:"cm"` // sdk 定制方
	R  int64  `json:"r"`  // 随机码;
}

func (a SdkInfo) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
