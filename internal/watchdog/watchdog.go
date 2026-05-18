// Package watchdog provides a self-monitoring component that tracks the
// health and liveness of the driftwatch daemon itself, restarting stalled
// components and emitting heartbeat signals at a configurable interval.
package watchdog

import (
	"context"
	"log"
	"sync"
	"time"
)

// DefaultHeartbeat is the interval between liveness pings.
const DefaultHeartbeat = 30 * time.Second

// Pinger is implemented by any component that can report its own liveness.
type Pinger interface {
	Ping() error
}

// Watchdog monitors a set of named Pingers and calls the OnFailure hook
// when a component fails to respond within the heartbeat window.
type Watchdog struct {
	heartbeat  time.Duration
	pingers    map[string]Pinger
	mu         sync.RWMutex
	onFailure  func(name string, err error)
	logger     *log.Logger
}

// New creates a Watchdog with the given heartbeat interval.
// A nil logger falls back to the default logger.
func New(heartbeat time.Duration, logger *log.Logger) *Watchdog {
	if heartbeat <= 0 {
		heartbeat = DefaultHeartbeat
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Watchdog{
		heartbeat: heartbeat,
		pingers:   make(map[string]Pinger),
		logger:    logger,
		onFailure: func(name string, err error) {},
	}
}

// Register adds a named Pinger to the watch list.
func (w *Watchdog) Register(name string, p Pinger) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.pingers[name] = p
}

// OnFailure sets the callback invoked when a Pinger returns an error.
func (w *Watchdog) OnFailure(fn func(name string, err error)) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.onFailure = fn
}

// Start begins the heartbeat loop, blocking until ctx is cancelled.
func (w *Watchdog) Start(ctx context.Context) {
	ticker := time.NewTicker(w.heartbeat)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.check()
		}
	}
}

func (w *Watchdog) check() {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for name, p := range w.pingers {
		if err := p.Ping(); err != nil {
			w.logger.Printf("watchdog: component %q failed ping: %v", name, err)
			w.onFailure(name, err)
		} else {
			w.logger.Printf("watchdog: component %q alive", name)
		}
	}
}
