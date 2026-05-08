// Package ratelimit provides a token-bucket rate limiter for controlling
// how frequently SSH checks are performed against remote hosts.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate of operations using a simple token-bucket approach.
type Limiter struct {
	mu       sync.Mutex
	tokens   int
	max      int
	refillAt time.Duration
	last     time.Time
	clock    func() time.Time
}

// Config holds configuration for the Limiter.
type Config struct {
	// MaxTokens is the maximum number of tokens (burst capacity).
	MaxTokens int
	// RefillInterval is how often a token is added back.
	RefillInterval time.Duration
}

// New creates a Limiter with the given config. Defaults are applied for
// zero-value fields.
func New(cfg Config) *Limiter {
	if cfg.MaxTokens <= 0 {
		cfg.MaxTokens = 10
	}
	if cfg.RefillInterval <= 0 {
		cfg.RefillInterval = time.Second
	}
	return &Limiter{
		tokens:   cfg.MaxTokens,
		max:      cfg.MaxTokens,
		refillAt: cfg.RefillInterval,
		last:     time.Now(),
		clock:    time.Now,
	}
}

// Allow reports whether an operation is permitted. It consumes one token
// and refills based on elapsed time since the last call.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.last)
	refill := int(elapsed / l.refillAt)
	if refill > 0 {
		l.tokens += refill
		if l.tokens > l.max {
			l.tokens = l.max
		}
		l.last = now
	}

	if l.tokens <= 0 {
		return false
	}
	l.tokens--
	return true
}

// Available returns the current number of available tokens.
func (l *Limiter) Available() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.tokens
}
