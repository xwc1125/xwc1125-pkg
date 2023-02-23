// Package jwtauth
//
// @author: xwc1125
package jwtauth

import (
	"errors"

	"github.com/xwc1125/xwc1125-pkg/types/response"
)

var (
	TokenExactFailed           string = "token不存在或header设置不正确"
	TokenExpire                string = "回话已过期"
	TokenCreateFailed          string = "生成token错误"
	TokenParseFailed           string = "token解析错误"
	TokenParseFailedAndEmpty   string = "解析错误,token为空"
	TokenParseFailedAndInvalid string = "解析错误,token无效"
	TokenUnsupportedOptions           = errors.New("unsupported options")

	ErrTokenIllegal = &response.StatusErr{480, "err.ErrTokenIllegal", "Illegal token", ""}
	ErrTokenReplace = &response.StatusErr{481, "err.TokenReplace", "Other clients logged in", ""}
	ErrTokenExpired = &response.StatusErr{482, "err.TokenExpired", "Token expired", ""}
)
