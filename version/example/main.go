// Package main
//
// @author: xwc1125
package main

import (
	"github.com/xwc1125/xwc1125-pkg/version"
)

func main() {
	version.FilePath = "./logs"
	if version.Build("App: versionApp") {
		return
	}
}
