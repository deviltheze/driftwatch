// Package backoff provides exponential backoff strategies for SSH reconnection
// and retry scenarios within driftwatch.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Policy defines the parameters for an exponential backoff strategy.
type Policy struct {
	// InitialInterval is the starting wait duration.
	InitialInterval time.Duration
	// MaxInterval caps the computed wait duration.
	MaxInterval time.Duration
	// Multiplier is applied to the interval after each attempt.
	Multiplier float64
	// Jitter adds randomness as a fraction of the computed interval (0–1).
	Jitter float64
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      1.5,
		Jitter:          0.2,
	}
}

// Next returns the wait duration for the given attempt number (0-indexed).
// The duration grows exponentially up to MaxInterval, with optional jitter.
func (p Policy) Next(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	base := float64(p.InitialInterval) * math.Pow(p.Multiplier, float64(attempt))
	if base > float64(p.MaxInterval) {
		base = float64(p.MaxInterval)
	}

	if p.Jitter > 0 {
		// Apply symmetric jitter: base ± (jitter * base)
		delta := p.Jitter * base
		base = base - delta + rand.Float64()*2*delta //nolint:gosec
	}

	if base < 0 {
		base = 0
	}

	return time.Duration(base)
}

// Steps returns a slice of wait durations for n attempts.
func (p Policy) Steps(n int) []time.Duration {
	durations := make([]time.Duration, n)
	for i := 0; i < n; i++ {
		durations[i] = p.Next(i)
	}
	return durations
}
