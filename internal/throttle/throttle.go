// Package throttle provides per-host SSH connection throttling to prevent
// overwhelming remote targets during concurrent drift checks.
package throttle

import (
	"fmt"
	"sync"
	"time"
)

// Throttle limits concurrent operations per host using a semaphore-based
// approach with an optional minimum delay between acquisitions.
type Throttle struct {
	mu       sync.Mutex
	sems     map[string]chan struct{}
	maxConns int
	delay    time.Duration
}

// New creates a Throttle that allows at most maxConns concurrent operations
// per host. delay is the minimum wait between successive Acquire calls for
// the same host (0 means no delay).
func New(maxConns int, delay time.Duration) *Throttle {
	if maxConns <= 0 {
		maxConns = 1
	}
	return &Throttle{
		sems:     make(map[string]chan struct{}),
		maxConns: maxConns,
		delay:    delay,
	}
}

// Acquire blocks until a slot is available for host, then reserves it.
// The caller must call Release with the same host when done.
func (t *Throttle) Acquire(host string) error {
	if host == "" {
		return fmt.Errorf("throttle: host must not be empty")
	}
	t.mu.Lock()
	ch, ok := t.sems[host]
	if !ok {
		ch = make(chan struct{}, t.maxConns)
		t.sems[host] = ch
	}
	t.mu.Unlock()

	ch <- struct{}{}
	if t.delay > 0 {
		time.Sleep(t.delay)
	}
	return nil
}

// Release frees a previously acquired slot for host.
func (t *Throttle) Release(host string) {
	t.mu.Lock()
	ch, ok := t.sems[host]
	t.mu.Unlock()
	if !ok {
		return
	}
	select {
	case <-ch:
	default:
	}
}

// Active returns the number of currently held slots for host.
func (t *Throttle) Active(host string) int {
	t.mu.Lock()
	ch, ok := t.sems[host]
	t.mu.Unlock()
	if !ok {
		return 0
	}
	return len(ch)
}
