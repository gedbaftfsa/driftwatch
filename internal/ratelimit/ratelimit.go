// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently drift checks are dispatched per service.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate of operations using a token-bucket strategy.
type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     time.Duration // minimum interval between allowed calls
	maxBurst int
}

type bucket struct {
	tokens    int
	lastRefil time.Time
}

// New creates a Limiter that allows at most maxBurst calls per rate interval per key.
func New(rate time.Duration, maxBurst int) *Limiter {
	if maxBurst < 1 {
		maxBurst = 1
	}
	return &Limiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		maxBurst: maxBurst,
	}
}

// Allow reports whether a call for the given key is permitted at this time.
// If permitted, one token is consumed.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.maxBurst, lastRefil: now}
		l.buckets[key] = b
	}

	// Refill tokens based on elapsed time.
	if l.rate > 0 {
		elapsed := now.Sub(b.lastRefil)
		gained := int(elapsed / l.rate)
		if gained > 0 {
			b.tokens += gained
			if b.tokens > l.maxBurst {
				b.tokens = l.maxBurst
			}
			b.lastRefil = b.lastRefil.Add(time.Duration(gained) * l.rate)
		}
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// Reset clears the bucket for the given key, allowing immediate access.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Remaining returns the number of tokens available for the given key.
func (l *Limiter) Remaining(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[key]
	if !ok {
		return l.maxBurst
	}
	return b.tokens
}
