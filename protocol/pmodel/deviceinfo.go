// Package pmodel
//
// @author: xwc1125
package pmodel

import "github.com/chain5j/chain5j-pkg/codec/json"

type DeviceInfo struct {
	UID string `json:"uid"` // 唯一码
	Dn  string `json:"dn"`  // 设备名称deviceName
	Dsn string `json:"dsn"` // 设备系统名称deviceSystemName
	Dsv string `json:"dsv"` // 设备系统版本名称DeviceSystemVersion
	Dm  string `json:"dm"`  // 设备model
	Dlm string `json:"dlm"` // 设备localizeModel
	R   int64  `json:"r"`   // 随机码
}

func (a DeviceInfo) String() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
