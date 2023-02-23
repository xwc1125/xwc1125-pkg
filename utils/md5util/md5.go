// Package md5util
//
// @author: xwc1125
package md5util

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5 returns the MD5 checksum string of the data.
func Md5(b []byte) string {
	checksum := md5.Sum(b)
	return hex.EncodeToString(checksum[:])
}
