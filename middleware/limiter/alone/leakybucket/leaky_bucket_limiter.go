package leakybucket

import (
	"context"
	"sync"

	"github.com/xwc1125/xwc1125-pkg/middleware/limiter"
)

type leakyBucketEntry struct {
	allow bool
}

var _ limiter.Entry = (*leakyBucketEntry)(nil)

func (e *leakyBucketEntry) Allow() bool {
	return e.allow
}

func (e *leakyBucketEntry) Finish() {}

func (e *leakyBucketEntry) Error() error { return nil }

type leakyBucketLimiter struct {
	limiters sync.Map
}

var _ limiter.Limiter = (*leakyBucketLimiter)(nil)

func NewLimiter() *leakyBucketLimiter {
	return &leakyBucketLimiter{
		limiters: sync.Map{},
	}
}

func (l *leakyBucketLimiter) Check(ctx context.Context, r limiter.Resource) limiter.Entry {
	return &leakyBucketEntry{allow: l.getLimiter(r).Allow()}
}

func (l *leakyBucketLimiter) SetLimit(ctx context.Context, r limiter.Resource) {
	l.getLimiter(r).SetRate(r.Limit)
}

func (l *leakyBucketLimiter) SetBurst(ctx context.Context, r limiter.Resource) {
	l.getLimiter(r).SetVolume(r.Burst)
}

func (l *leakyBucketLimiter) SetWindow(ctx context.Context, r limiter.Resource) {}

func (l *leakyBucketLimiter) getLimiter(r limiter.Resource) (lim LeakyBucket) {
	val, ok := l.limiters.Load(r.Name)
	if !ok {
		lim = NewLeakyBucket(r.Limit, r.Burst)
		l.limiters.Store(r.Name, lim)
		return
	}

	if lim, ok = val.(LeakyBucket); !ok {
		lim = NewLeakyBucket(r.Limit, r.Burst)
		l.limiters.Store(r.Name, lim)
		return
	}

	return
}
