// Package base
//
// @author: xwc1125
package base

import (
	"net/http"

	"github.com/xwc1125/xwc1125-pkg/jwtauth"
	"github.com/xwc1125/xwc1125-pkg/protocol/contextx"
	"github.com/xwc1125/xwc1125-pkg/types"
	"github.com/xwc1125/xwc1125-pkg/types/response"
)

// BaseApi api的base
type BaseApi struct {
	Uid      int64
	UserName string
}

// IsJwtAuth 进行auth验证
func (a *BaseApi) IsJwtAuth(ctx contextx.Context, jwtConfig jwtauth.JWTConfig) bool {
	v, _ := ctx.Get(jwtConfig.JwtContextKey)
	var (
		userToken *jwtauth.UserToken
		empty     interface{}
	)
	if v == nil || v == empty {
		a.Error(ctx, http.StatusUnauthorized, "登录失败")
		return false
	} else {
		userToken = v.(*jwtauth.UserToken)
	}
	a.Uid = userToken.Uid
	a.UserName = userToken.Username
	return true
}
func (a *BaseApi) OK(ctx contextx.Context, msg string, data ...interface{}) {
	response.Ok(ctx, msg, data...)
}
func (a *BaseApi) FailTemp(ctx contextx.Context, status types.JsonError, data ...interface{}) {
	response.FailTemp(ctx, status, data...)
}
func (a *BaseApi) Fail(ctx contextx.Context, status int, msg string, data ...interface{}) {
	response.Fail(ctx, status, msg, data...)
}
func (a *BaseApi) Error(ctx contextx.Context, code int, msg string, data ...interface{}) {
	response.Error(ctx, code, msg, data...)
}

func (a *BaseApi) PageOK(ctx contextx.Context, msg string, list interface{}, total int64, pageIndex int64, pageSize int64, filter interface{}) {
	response.Ok(ctx, msg, response.PageInfo{
		Total:     total,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		List:      list,
		Filter:    filter,
	})
}
