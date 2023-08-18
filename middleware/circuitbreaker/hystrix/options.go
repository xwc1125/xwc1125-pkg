// Package hystrix
package hystrix

type Option func(breaker *breaker)

// WithTimeout 单次请求超时时间,默认时间是1000毫秒,milliseconds
func WithTimeout(timeout int) Option {
	return func(b *breaker) {
		b.timeout = timeout
	}
}

// WithMaxConcurrentRequests 最大并发量,默认值是10
func WithMaxConcurrentRequests(maxConcurrentRequests int) Option {
	return func(b *breaker) {
		b.maxConcurrentRequests = maxConcurrentRequests
	}
}

// WithRequestVolumeThreshold 一个统计窗口10秒内请求数量。达到这个请求数量后才去判断是否要开启熔断,默认值是20(比如10秒内接到了11个请求只超过了1个)
func WithRequestVolumeThreshold(requestVolumeThreshold int) Option {
	return func(b *breaker) {
		b.requestVolumeThreshold = requestVolumeThreshold
	}
}

// WithSleepWindow 熔断后多久去尝试服务是否可用,默认值是5000毫秒(熔断器打开到半打开的时间),milliseconds
func WithSleepWindow(sleepWindow int) Option {
	return func(b *breaker) {
		b.sleepWindow = sleepWindow
	}
}

// WithErrorPercentThreshold 错误百分比,默认值是50(当错误百分比超过这个限制时则进行熔断)。请求数量大于等于RequestVolumeThreshold并且错误率到达这个百分比后就会启动熔断
func WithErrorPercentThreshold(errorPercentThreshold int) Option {
	return func(b *breaker) {
		b.errorPercentThreshold = errorPercentThreshold
	}
}
