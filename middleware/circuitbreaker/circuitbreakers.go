// Package circuitbreaker
package circuitbreaker

import (
	"fmt"
	"sync"
)

type Breakers struct {
	lock     sync.RWMutex
	breakers map[string]Breaker
}

func NewBreakers() *Breakers {
	return &Breakers{
		lock:     sync.RWMutex{},
		breakers: make(map[string]Breaker),
	}
}

func (b *Breakers) Add(breaker Breaker) {
	b.lock.Lock()
	b.breakers[breaker.Name()] = breaker
	b.lock.Unlock()
}

func (b *Breakers) Get(breakerName string) (Breaker, error) {
	b.lock.Lock()
	breaker, ok := b.breakers[breakerName]
	b.lock.Unlock()
	if ok {
		return breaker, nil
	}
	return nil, fmt.Errorf("breaker:{%s} does not exist", breakerName)
}

func (b *Breakers) Allow(breakerName string) error {
	breaker, err := b.Get(breakerName)
	if err != nil {
		return err
	}
	return breaker.Allow()
}

func (b *Breakers) Do(breakerName string, req func() error, fallback func(err error) error, acceptable Acceptable) error {
	breaker, err := b.Get(breakerName)
	if err != nil {
		return err
	}
	return breaker.Do(req, fallback, acceptable)
}
