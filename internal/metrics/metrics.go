// Package metrics provides lightweight in-memory counters for tracking
// drift check activity across the driftwatch daemon lifecycle.
package metrics

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Snapshot holds a point-in-time copy of all tracked counters.
type Snapshot struct {
	ChecksTotal   int64
	DriftDetected int64
	HostsChecked  int64
	LastCheckAt   time.Time
}

// Collector accumulates runtime metrics in a thread-safe manner.
type Collector struct {
	mu            sync.RWMutex
	checksTotal   int64
	driftDetected int64
	hostsChecked  int64
	lastCheckAt   time.Time
}

// NewCollector returns an initialised, zero-valued Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// RecordCheck increments the total check counter and records the timestamp.
func (c *Collector) RecordCheck(hostsInRun int64, drifted int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checksTotal++
	c.hostsChecked += hostsInRun
	c.driftDetected += drifted
	c.lastCheckAt = time.Now()
}

// Snapshot returns a consistent copy of the current counters.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ChecksTotal:   c.checksTotal,
		DriftDetected: c.driftDetected,
		HostsChecked:  c.hostsChecked,
		LastCheckAt:   c.lastCheckAt,
	}
}

// Write prints a human-readable summary of the snapshot to w.
func (c *Collector) Write(w io.Writer) {
	s := c.Snapshot()
	last := "never"
	if !s.LastCheckAt.IsZero() {
		last = s.LastCheckAt.Format(time.RFC3339)
	}
	fmt.Fprintf(w, "checks_total=%d hosts_checked=%d drift_detected=%d last_check=%s\n",
		s.ChecksTotal, s.HostsChecked, s.DriftDetected, last)
}
