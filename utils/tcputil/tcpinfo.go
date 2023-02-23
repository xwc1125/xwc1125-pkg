// Package tcputil
//
// @author: xwc1125
package tcputil

import "github.com/chain5j/chain5j-pkg/codec/json"

type TcpInfo struct {
	Ip        string `json:"ip"`
	IpAddress string `json:"ipAddress"`
	Url       string `json:"url"`
	Api       string `json:"api"`
	Ua        string `json:"ua"`
	Browser   string `json:"browser"`
	Os        string `json:"os"`
	IsMobile  bool   `json:"isMobile"`
	Lang      string `json:"lang"`
}

func (t TcpInfo) String() string {
	bytes, _ := json.Marshal(t)
	return string(bytes)
}
