// Package hystrix
package hystrix

import (
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/xwc1125/xwc1125-pkg/middleware/circuitbreaker"
)

var (
	_ circuitbreaker.Breaker = new(breaker)
)

type breaker struct {
	name string

	timeout                int
	maxConcurrentRequests  int
	requestVolumeThreshold int
	sleepWindow            int
	errorPercentThreshold  int
}

func NewBreaker(name string, opts ...Option) circuitbreaker.Breaker {
	// 3.设置熔断器
	// 第一个参数:当前创建的熔断器名称
	// 第二个参数: hystrix.CommandConfig配置的限流规则
	b := &breaker{
		name: name,
	}
	for _, opt := range opts {
		opt(b)
	}

	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:                b.timeout,                // 单次请求超时时间,默认时间是1000毫秒。
		MaxConcurrentRequests:  b.maxConcurrentRequests,  // 最大并发量,默认值是10
		SleepWindow:            b.sleepWindow,            // 熔断后多久去尝试服务是否可用,默认值是5000毫秒(熔断器打开到半打开的时间)
		RequestVolumeThreshold: b.requestVolumeThreshold, // 一个统计窗口10秒内请求数量。达到这个请求数量后才去判断是否要开启熔断,默认值是20(比如10秒内接到了11个请求只超过了1个,当前为1没有超过熔断限制,则不熔断)
		ErrorPercentThreshold:  b.errorPercentThreshold,  // 错误百分比,默认值是50(当错误百分比超过这个限制时则进行熔断)。请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
	})
	return b
}

func getValueByDefault(value int, defaultVal int) int {
	if value == 0 {
		return defaultVal
	}
	return value
}

func (b *breaker) Name() string {
	return b.name
}

func (b *breaker) Allow() error {
	return fmt.Errorf("hystrix does not support allow")
}

func (b *breaker) Do(req func() error, fallback func(err error) error, acceptable circuitbreaker.Acceptable) error {
	err := hystrix.Do(b.name, req, fallback)
	if acceptable != nil {
		b := acceptable(err)
		if b {
			return nil
		}
		return err
	}
	return err
}
