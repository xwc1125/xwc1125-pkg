// Package response
//
// @author: xwc1125
package response

import (
	"github.com/chain5j/chain5j-pkg/util/dateutil"
	contextx2 "github.com/xwc1125/xwc1125-pkg/protocol/contextx"
	"github.com/xwc1125/xwc1125-pkg/types"
)

var (
	Err404       = &StatusErr{404, "err.Err404", "No found", ""}
	FailStatus   = &StatusErr{2001, "fail", "Failed", ""}
	NoPermission = &StatusErr{4001, "err.NoPermission", "No Permission", ""}
)

type Map map[string]interface{}

type Response struct {
	Code int         `json:"code,omitempty"`
	Msg  interface{} `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
	T    interface{} `json:"t,omitempty"`
}

// JSON httpStatus=code
func JSON(ctx contextx2.Context, code int, data interface{}) {
	contextx2.JSON(ctx, code, data)
}

// Ok httpStatus=200,resp.Code=200
func Ok(ctx contextx2.Context, msg string, data ...interface{}) {
	result := &Response{
		Code: StatusOK,
		Msg:  msg,
		T:    dateutil.CurrentTimeSecond(),
	}
	if len(data) == 1 {
		result.Data = data[0]
	} else {
		result.Data = data
	}
	JSON(ctx, StatusOK, result)
}

// OkDefault httpStatus=200,resp.Code=200
func OkDefault(ctx contextx2.Context) {
	Ok(ctx, "Success")
}

// OkData httpStatus=200,resp.Code=200
func OkData(ctx contextx2.Context, data interface{}) {
	Ok(ctx, "Success", data)
}

// Fail httpStatus=200,resp.Code=status
func Fail(ctx contextx2.Context, status int, msg string, data ...interface{}) {
	ctx.Abort()
	result := &Response{
		Code: status,
		Msg:  msg,
		T:    dateutil.CurrentTimeSecond(),
	}

	if len(data) == 1 {
		result.Data = data[0]
	} else {
		result.Data = data
	}
	JSON(ctx, StatusOK, result)
}

// FailDefault httpStatus=200,resp.Code=2001
func FailDefault(ctx contextx2.Context) {
	Fail(ctx, FailStatus.Status, "Failed")
}

// FailTemp httpStatus=200,resp.Code=status
func FailTemp(ctx contextx2.Context, status types.JsonError, data ...interface{}) {
	Fail(ctx, status.Code(), status.Msg(), data...)
}

// Error httpStatus=code,resp.Code=code
func Error(ctx contextx2.Context, code int, msg string, data ...interface{}) {
	ctx.Abort()
	result := &Response{
		Code: code,
		Msg:  msg,
		T:    dateutil.CurrentTimeSecond(),
	}

	if len(data) == 1 {
		result.Data = data[0]
	} else {
		result.Data = data
	}
	JSON(ctx, code, result)
}

// ErrorDefault httpStatus=code,resp.Code=code
func ErrorDefault(ctx contextx2.Context, code int) {
	Error(ctx, code, "Error")
}

// ErrorTemp httpStatus=code,resp.Code=code
func ErrorTemp(ctx contextx2.Context, code types.JsonError, data ...interface{}) {
	Error(ctx, code.Code(), code.Msg(), data...)
}

// Unauthorized 401 error define
// httpStatus=401,resp.Code=401
func Unauthorized(ctx contextx2.Context, msg string, data ...interface{}) {
	Error(ctx, StatusUnauthorized, msg, data...)
}
