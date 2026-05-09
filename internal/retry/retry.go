// Package retry provides configurable retry logic for transient failures
// encountered during SSH connections and remote command execution.
package retry

import (
	"errors"
	"log"
	"math"
	"time"
)

// Policy defines the retry behaviour for a given operation.
type Policy struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration
	// MaxDelay caps the exponential back-off ceiling.
	MaxDelay time.Duration
	// Multiplier is the exponential growth factor (default 2.0).
	Multiplier float64
}

// DefaultPolicy returns a sensible retry policy for SSH operations.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts:  3,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Multiplier:   2.0,
	}
}

// Do executes fn according to p, retrying on non-nil errors.
// It returns the last error if all attempts are exhausted.
func Do(p Policy, logger *log.Logger, fn func() error) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	mul := p.Multiplier
	if mul <= 0 {
		mul = 2.0
	}

	var err error
	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		if attempt == p.MaxAttempts {
			break
		}
		delay := delay(p.InitialDelay, mul, attempt, p.MaxDelay)
		if logger != nil {
			logger.Printf("retry: attempt %d/%d failed (%v); retrying in %s",
				attempt, p.MaxAttempts, err, delay)
		}
		time.Sleep(delay)
	}
	return errors.New("retry: all attempts exhausted: " + err.Error())
}

// delay calculates the back-off duration for the given attempt number (1-based).
func delay(initial time.Duration, multiplier float64, attempt int, max time.Duration) time.Duration {
	d := float64(initial) * math.Pow(multiplier, float64(attempt-1))
	if d > float64(max) {
		d = float64(max)
	}
	return time.Duration(d)
}
