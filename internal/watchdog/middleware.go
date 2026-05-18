package watchdog

import (
	"fmt"
	"sync/atomic"
	"time"
)

// PingerFunc is a convenience adapter allowing plain functions to satisfy
// the Pinger interface.
type PingerFunc func() error

// Ping calls the underlying function.
func (f PingerFunc) Ping() error { return f() }

// LivenessProbe wraps an arbitrary function and tracks the last successful
// ping time. It implements Pinger and can be registered with a Watchdog.
type LivenessProbe struct {
	name    string
	probe   func() error
	lastOK  atomic.Int64 // unix nano
	stallAfter time.Duration
}

// NewLivenessProbe creates a LivenessProbe that considers itself stalled if
// the underlying probe has not succeeded within stallAfter.
func NewLivenessProbe(name string, probe func() error, stallAfter time.Duration) *LivenessProbe {
	return &LivenessProbe{
		name:       name,
		probe:      probe,
		stallAfter: stallAfter,
	}
}

// Ping executes the probe function. On success the last-OK timestamp is
// updated. If the probe itself succeeds but the stall window has elapsed
// since the previous success, an error is returned.
func (lp *LivenessProbe) Ping() error {
	if err := lp.probe(); err != nil {
		return fmt.Errorf("probe %q: %w", lp.name, err)
	}
	now := time.Now().UnixNano()
	prev := lp.lastOK.Swap(now)
	if prev != 0 && lp.stallAfter > 0 {
		elapsed := time.Duration(now - prev)
		if elapsed > lp.stallAfter {
			return fmt.Errorf("probe %q stalled: last success %s ago", lp.name, elapsed.Round(time.Second))
		}
	}
	return nil
}
