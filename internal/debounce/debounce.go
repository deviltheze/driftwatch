// Package debounce provides a mechanism to suppress repeated drift events
// for the same host within a configurable time window, reducing alert noise.
package debounce

import (
	"sync"
	"time"
)

// Event represents a debounce-tracked occurrence for a host.
type Event struct {
	Host      string
	LastSeen  time.Time
	Count     int
}

// Debouncer suppresses repeated events within a window.
type Debouncer struct {
	mu     sync.Mutex
	window time.Duration
	events map[string]*Event
	now    func() time.Time
}

// New creates a Debouncer with the given suppression window.
// A zero or negative window disables suppression (all events pass).
func New(window time.Duration) *Debouncer {
	return &Debouncer{
		window: window,
		events: make(map[string]*Event),
		now:    time.Now,
	}
}

// Allow returns true if the event for host should be forwarded.
// Repeated events within the window are suppressed; the first occurrence
// always passes. Thread-safe.
func (d *Debouncer) Allow(host string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()

	if ev, ok := d.events[host]; ok {
		if d.window > 0 && now.Sub(ev.LastSeen) < d.window {
			ev.Count++
			return false
		}
		ev.LastSeen = now
		ev.Count = 1
		return true
	}

	d.events[host] = &Event{
		Host:     host,
		LastSeen: now,
		Count:    1,
	}
	return true
}

// Reset clears the debounce state for a specific host.
func (d *Debouncer) Reset(host string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.events, host)
}

// Stats returns a copy of the current event record for a host, and whether
// it exists.
func (d *Debouncer) Stats(host string) (Event, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ev, ok := d.events[host]
	if !ok {
		return Event{}, false
	}
	return *ev, true
}
