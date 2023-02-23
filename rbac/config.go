// Package rbac
//
// @author: xwc1125
package rbac

import (
	"fmt"
)

const (
	PrefixUserID  = "u"
	PrefixRoleID  = "r"
	PrefixGroupID = "g"
	PrefixPermID  = "p"
)

func GetUserKey(uid int64) string {
	return fmt.Sprintf("%s%d", PrefixUserID, uid)
}
func GetRoleKey(rid int64) string {
	return fmt.Sprintf("%s%d", PrefixRoleID, rid)
}
func GetGroupKey(gid int64) string {
	return fmt.Sprintf("%s%d", PrefixGroupID, gid)
}
func GetPermKey(pid int64) string {
	return fmt.Sprintf("%s%d", PrefixPermID, pid)
}

const (
	RType = "r"
	PType = "p"
	EType = "e"
	MType = "m"
)

const (
	TYPE_MODULE  = 1 // 模块
	TYPE_MENU    = 2 // 菜单
	TYPE_OPERATE = 3 // 操作
)
