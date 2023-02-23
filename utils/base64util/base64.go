// Package base64util
//
// @author: xwc1125
package base64util

import "encoding/base64"

// Decode 解密
func Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Encode 加密
func Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}
