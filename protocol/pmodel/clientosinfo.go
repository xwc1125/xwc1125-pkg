// Package pmodel
//
// @author: xwc1125
package pmodel

import "strings"

type ClientOsType int

const (
	UNKNOWN ClientOsType = -1
	ANDROID ClientOsType = 0
	IOS     ClientOsType = 1
	WEB     ClientOsType = 2
	WINDOWS ClientOsType = 3
	MACBOOK ClientOsType = 4
	LINUX   ClientOsType = 5
)

func (p ClientOsType) String() string {
	switch p {
	case UNKNOWN:
		return "UNKNOWN"
	case ANDROID:
		return "ANDROID"
	case IOS:
		return "IOS"
	case WEB:
		return "WEB"
	case WINDOWS:
		return "WINDOWS"
	case MACBOOK:
		return "MACBOOK"
	case LINUX:
		return "LINUX"
	default:
		return "UNKNOWN"
	}
}

func ParseOsTypeByType(t int) ClientOsType {
	switch t {
	case -1:
		return UNKNOWN
	case 0:
		return ANDROID
	case 1:
		return IOS
	case 2:
		return WEB
	case 3:
		return WINDOWS
	case 4:
		return MACBOOK
	case 5:
		return LINUX
	default:
		return UNKNOWN
	}
}

func ParseOsTypeByName(n string) ClientOsType {
	switch strings.ToUpper(n) {
	case "UNKNOWN":
		return UNKNOWN
	case "ANDROID":
		return ANDROID
	case "IOS":
		return IOS
	case "WEB":
		return WEB
	case "WINDOWS":
		return WINDOWS
	case "MACBOOK":
		return MACBOOK
	case "LINUX":
		return LINUX
	default:
		return UNKNOWN
	}
}
