// Package jwtauth
//
// @author: xwc1125
package jwtauth

import (
	"github.com/dgrijalva/jwt-go"
)

// UserToken 用户信息
type UserToken struct {
	Uid       int64     `form:"uid"`      // UID
	Username  string    `form:"username"` // 用户名
	ExtraData MapClaims // 扩展内容
}

// JwtClaims jwt claims信息
type JwtClaims struct {
	UserToken
	jwt.StandardClaims
}

// JwtRefreshClaims 刷新claims
type JwtRefreshClaims struct {
	Token string `json:"token"`
	jwt.StandardClaims
}
