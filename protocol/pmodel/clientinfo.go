// Package pmodel
//
// @author: xwc1125
package pmodel

import "github.com/chain5j/chain5j-pkg/codec/json"

type ClientInfo struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	AesKey       string `json:"aesKey"`
	R            int64  `json:"r"`
}

func (a ClientInfo) ToString() string {
	bytes, _ := json.Marshal(a)
	return string(bytes)
}
