package ratelimit

import (
	"sync"
	"time"
)

type RateLimitWindow struct {
	counter     int
	limit       int
	window      time.Duration
	windowStart time.Time
	mu          sync.Mutex
}

func NewRateLimitWindow(limit int, window time.Duration) *RateLimitWindow {
	return &RateLimitWindow{
		counter:     0,
		limit:       limit,
		window:      window,
		windowStart: time.Now(),
	}
}

func (r *RateLimitWindow) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if time.Since(r.windowStart) >= r.window {
		r.windowStart = time.Now()
		r.counter = 0
	}
	if r.limit <= r.counter {
		return false
	}
	r.counter++
	return true

}
