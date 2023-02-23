// Package pmodel
//
// @author: xwc1125
package pmodel

import "github.com/chain5j/chain5j-pkg/codec/json"

type ApkInfo struct {
	PkName       string `json:"pkName"`       // App的包名
	Label        string `json:"label"`        // App名称
	VersionName  string `json:"versionName"`  // App版本
	VersionCode  int    `json:"versionCode"`  // App版本号
	ApkFile      int    `json:"apkFile"`      // App的位置
	Installed    int    `json:"installed"`    // 0:未安装，1:已安装
	InstallState int    `json:"installState"` // -1:未知，0：失败，1：成功，2：升级失败，3：升级成功
}

func (a ApkInfo) ToString() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
