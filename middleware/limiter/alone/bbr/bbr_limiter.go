// Package bbr
package bbr

import (
	"context"
	"sync"

	"github.com/xwc1125/xwc1125-pkg/middleware/limiter"
)

var (
	_ limiter.Entry = new(bbrEntry)
)

type bbrEntry struct {
	bbr *BBR
}

func (b *bbrEntry) Allow() bool {
	return b.bbr.Allow()
}

func (b *bbrEntry) Finish() {

}

func (b *bbrEntry) Error() error {
	return nil
}

var (
	_ limiter.Limiter = new(bbrLimiter)
)

type bbrLimiter struct {
	limiters sync.Map // key resource name,value *rate.limiter
}

func NewLimiter() *bbrLimiter {
	return &bbrLimiter{
		limiters: sync.Map{},
	}
}

func (b *bbrLimiter) Check(ctx context.Context, r limiter.Resource) limiter.Entry {
	return &bbrEntry{
		bbr: b.getLimiter(r),
	}
}

func (b *bbrLimiter) SetLimit(ctx context.Context, r limiter.Resource) {

}

func (b *bbrLimiter) SetBurst(ctx context.Context, r limiter.Resource) {
	b.getLimiter(r).SetBurst(r.Burst)
}

func (b *bbrLimiter) SetWindow(ctx context.Context, r limiter.Resource) {
	b.getLimiter(r).SetWindow(r.Window)
}
func (b *bbrLimiter) getLimiter(r limiter.Resource) (lim *BBR) {
	val, ok := b.limiters.Load(r.Name)
	if !ok {
		lim = NewBBR(WithWindow(r.Window), WithBucket(r.Burst))
		b.limiters.Store(r.Name, lim)
		return
	}

	if lim, ok = val.(*BBR); !ok {
		lim = NewBBR(WithWindow(r.Window), WithBucket(r.Burst))
		b.limiters.Store(r.Name, lim)
		return
	}

	return
}
