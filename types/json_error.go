// Package types
//
// @author: xwc1125
package types

type JsonError interface {
	Code() int
	Error() string
	Msg() string
}
