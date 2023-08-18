package slidingwindow

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/emirpasic/gods/queues/arrayqueue"
)

type SlidingWindow interface {
	Allow() bool
	SetLimit(rate int)
	SetWindow(window time.Duration)
}

type node struct {
	t time.Time
}

type option struct {
	clock clock.Clock
}

type Option func(o *option)

func WithClock(clock clock.Clock) Option {
	return func(o *option) { o.clock = clock }
}

func defaultOption() *option {
	return &option{}
}

type slidingWindow struct {
	mu     sync.RWMutex
	opts   *option
	window time.Duration
	limit  int
	q      *arrayqueue.Queue
}

func NewSlidingWindow(limit int, window time.Duration, opts ...Option) *slidingWindow {
	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}

	return &slidingWindow{
		opts:   opt,
		window: window,
		limit:  limit,
		q:      arrayqueue.New(),
	}
}

func (sl *slidingWindow) Allow() bool {
	sl.mu.RLock()
	limit := sl.limit
	window := sl.window
	sl.mu.RUnlock()

	now := sl.now()

	// not full, access allowed
	if sl.q.Size() < limit {
		sl.q.Enqueue(&node{t: now})
		return true
	}

	// take out the earliest one
	early, _ := sl.q.Peek()
	first := early.(*node)

	// the first request is still in the time window, access denied
	if now.Add(-window).Before(first.t) {
		return false
	}

	// pop the first request
	_, _ = sl.q.Dequeue()
	sl.q.Enqueue(&node{t: now})

	return true
}

func (sl *slidingWindow) SetLimit(limit int) {
	sl.mu.Lock()
	sl.limit = limit
	sl.mu.Unlock()
}

func (sl *slidingWindow) SetWindow(window time.Duration) {
	sl.mu.Lock()
	sl.window = window
	sl.mu.Unlock()
}

func (sl *slidingWindow) now() time.Time {
	if sl.opts.clock == nil {
		return time.Now()
	}
	return sl.opts.clock.Now()
}
