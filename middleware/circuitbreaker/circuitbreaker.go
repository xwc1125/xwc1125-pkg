package circuitbreaker

import (
	"errors"
)

// ErrNotAllowed error not allowed.
var ErrNotAllowed = errors.New("circuit breaker: not allowed for circuit open")

// Acceptable 自定义判定执行结果
type Acceptable func(err error) bool

type Breaker interface {
	// Name 熔断器名称
	Name() string

	// Allow 熔断方法，执行请求时必须手动上报执行结果
	Allow() error

	// Do 熔断方法，
	// req - 执行的函数
	// fallback - 支持自定义快速失败
	// acceptable - 支持自定义判定执行结果
	// 如果fallback和acceptable都为nil，自动上报执行结果
	Do(req func() error, fallback func(err error) error, acceptable Acceptable) error
}
