// Package circuitbreaker provides a simple circuit breaker to prevent
// repeated SSH connection attempts to hosts that are consistently failing.
package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// State represents the current state of a circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing, requests blocked
	StateHalfOpen              // testing if host recovered
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker tracks failure counts per host and opens the circuit when the
// threshold is exceeded, preventing further attempts until the reset timeout.
type Breaker struct {
	mu        sync.Mutex
	threshold int
	timeout   time.Duration
	hosts     map[string]*hostState
}

type hostState struct {
	failures  int
	state     State
	openedAt  time.Time
}

// New creates a Breaker with the given failure threshold and reset timeout.
// Zero values fall back to sensible defaults.
func New(threshold int, timeout time.Duration) *Breaker {
	if threshold <= 0 {
		threshold = 3
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Breaker{
		threshold: threshold,
		timeout:   timeout,
		hosts:     make(map[string]*hostState),
	}
}

// Allow returns true if a request to the given host should proceed.
func (b *Breaker) Allow(host string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	hs := b.getOrCreate(host)
	switch hs.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(hs.openedAt) >= b.timeout {
			hs.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess resets the failure count and closes the circuit for the host.
func (b *Breaker) RecordSuccess(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	hs := b.getOrCreate(host)
	hs.failures = 0
	hs.state = StateClosed
}

// RecordFailure increments the failure count and may open the circuit.
func (b *Breaker) RecordFailure(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	hs := b.getOrCreate(host)
	hs.failures++
	if hs.failures >= b.threshold {
		hs.state = StateOpen
		hs.openedAt = time.Now()
	}
}

// StateOf returns the current state for the given host.
func (b *Breaker) StateOf(host string) State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.getOrCreate(host).state
}

// Reset clears all state for the given host.
func (b *Breaker) Reset(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.hosts, host)
}

func (b *Breaker) getOrCreate(host string) *hostState {
	if hs, ok := b.hosts[host]; ok {
		return hs
	}
	hs := &hostState{state: StateClosed}
	b.hosts[host] = hs
	return hs
}

// ErrOpen is returned when the circuit is open for a host.
type ErrOpen struct{ Host string }

func (e ErrOpen) Error() string {
	return fmt.Sprintf("circuit open for host %q: too many recent failures", e.Host)
}
