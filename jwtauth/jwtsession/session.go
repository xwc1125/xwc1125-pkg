// Package jwtsession
//
// @author: xwc1125
package jwtsession

import (
	"github.com/chain5j/chain5j-pkg/util/convutil"
	"github.com/chain5j/logger"
	"github.com/xwc1125/xwc1125-pkg/jwtauth"
	"github.com/xwc1125/xwc1125-pkg/protocol/contextx"
)

var (
	JwtPayloadKey = "JWT_PAYLOAD"
	UserIDKey     = "identity"
	UserNameKey   = "user_name"
	RoleIDKey     = "role_id"
	RoleKey       = "role_key"
	RoleNameKey   = "role_name"
	DeptIDKey     = "dept_id"
	DeptNameKey   = "dept_name"
	DataScopeKey  = "data_scope_key"
)

// ExtractClaims 从request中获取jwtauth.MapClaims
func ExtractClaims(c contextx.Context) jwtauth.MapClaims {
	claims, exists := c.Get(JwtPayloadKey)
	if !exists {
		return make(jwtauth.MapClaims)
	}

	return claims.(jwtauth.MapClaims)
}

// Get 从jwt的claims中获取指定key的值
func Get(c contextx.Context, key string) interface{} {
	data := ExtractClaims(c)
	if data[key] != nil {
		return data[key]
	}
	logger.Warn("jwt claims miss key", "key", key, "url", c.Request().Method+":"+c.Request().URL.Path)
	return nil
}

// GetUserId 获取userId
func GetUserId(c contextx.Context) int {
	userId := Get(c, UserIDKey)
	if userId == nil {
		return 0
	}
	return convutil.ToInt(userId)
}

// GetUserIdStr 获取userId字符串
func GetUserIdStr(c contextx.Context) string {
	userId := Get(c, UserIDKey)
	if userId == nil {
		return ""
	}
	return convutil.ToString(userId)
}

// GetUserName 获取userName
func GetUserName(c contextx.Context) string {
	userName := Get(c, UserNameKey)
	if userName == nil {
		return ""
	}
	return convutil.ToString(userName)
}

// GetRoleId 获取roleId
func GetRoleId(c contextx.Context) int {
	data := Get(c, RoleIDKey)
	if data == nil {
		return 0
	}
	return convutil.ToInt(data)
}

// GetRoleName 获取角色名称
func GetRoleName(c contextx.Context) string {
	data := Get(c, RoleNameKey)
	if data == nil {
		return ""
	}
	return convutil.ToString(data)
}

// GetRoleKey 获取角色
func GetRoleKey(c contextx.Context) string {
	data := Get(c, RoleKey)
	if data == nil {
		return ""
	}
	return convutil.ToString(data)
}

// GetDeptId 获取部门ID
func GetDeptId(c contextx.Context) int {
	data := Get(c, DeptIDKey)
	if data == nil {
		return 0
	}
	return convutil.ToInt(data)
}

// GetDeptName 获取部门名称
func GetDeptName(c contextx.Context) string {
	data := Get(c, DeptNameKey)
	if data == nil {
		return ""
	}
	return convutil.ToString(data)
}
