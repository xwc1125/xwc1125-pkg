// Package response
//
// @author: xwc1125
package response

import (
	"bytes"
	"fmt"
)

// HTTP Methods copied from `net/http`.
const (
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodHead    = "HEAD"
	MethodPatch   = "PATCH"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

// HttpCode
const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1

	StatusOK                   = 200 // RFC 7231, 6.3.1
	StatusCreated              = 201 // RFC 7231, 6.3.2
	StatusAccepted             = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
	StatusNoContent            = 204 // RFC 7231, 6.3.5
	StatusResetContent         = 205 // RFC 7231, 6.3.6
	StatusPartialContent       = 206 // RFC 7233, 4.1
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices  = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently = 301 // RFC 7231, 6.4.2
	StatusFound            = 302 // RFC 7231, 6.4.3
	StatusSeeOther         = 303 // RFC 7231, 6.4.4
	StatusNotModified      = 304 // RFC 7232, 4.1
	StatusUseProxy         = 305 // RFC 7231, 6.4.5

	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3

	StatusBadRequest                   = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401 // RFC 7235, 3.1
	StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
	StatusForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound                     = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
	StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
	StatusConflict                     = 409 // RFC 7231, 6.5.8
	StatusGone                         = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
	StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
	StatusTeapot                       = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusTooEarly                     = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

type StatusErr struct {
	Status   int    `json:"status"`   // 状态值
	LangKey  string `json:"langKey"`  // 语言KeyCode
	Message  string `json:"msg"`      // 默认消息
	MoreInfo string `json:"moreInfo"` // 其他消息
}

func NewStatusErr(err error) *StatusErr {
	return &StatusErr{
		Status:   FailStatus.Status,
		LangKey:  "",
		Message:  err.Error(),
		MoreInfo: "",
	}
}
func (e *StatusErr) Code() int {
	return e.Status
}
func (e *StatusErr) Error() string {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%d", e.Status))
	if len(e.Message) > 0 {
		buff.WriteString(fmt.Sprintf(":%s", e.Message))
	} else if len(e.LangKey) > 0 {
		buff.WriteString(fmt.Sprintf(":%s", e.LangKey))
	}
	if len(e.MoreInfo) > 0 {
		buff.WriteString(fmt.Sprintf("->%s", e.MoreInfo))
	}
	return buff.String()
}
func (e *StatusErr) Msg() string {
	var buff bytes.Buffer

	if len(e.Message) > 0 {
		buff.WriteString(fmt.Sprintf("%s", e.Message))
	} else if len(e.LangKey) > 0 {
		buff.WriteString(fmt.Sprintf("%s", e.LangKey))
	}
	if len(e.MoreInfo) > 0 {
		buff.WriteString(fmt.Sprintf("->%s", e.MoreInfo))
	}
	return buff.String()
}
