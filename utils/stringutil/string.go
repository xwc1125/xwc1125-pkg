// Package stringutil
//
// @author: xwc1125
// @date: 2021/3/20
package stringutil

import "strings"

func IsEmpty(str string) bool {
	str = strings.TrimSpace(str)
	if str == "" {
		return true
	}
	str = strings.ToLower(str)
	if str == "null" {
		return true
	}
	return false
}
