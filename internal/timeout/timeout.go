// Package timeout provides per-host SSH operation timeout enforcement.
package timeout

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Policy holds timeout configuration for SSH operations.
type Policy struct {
	Connect time.Duration
	Command time.Duration
	Logger  *log.Logger
}

// DefaultPolicy returns a Policy with sensible defaults.
func DefaultPolicy() Policy {
	return Policy{
		Connect: 10 * time.Second,
		Command: 30 * time.Second,
		Logger:  log.Default(),
	}
}

// Enforcer wraps operations with deadline-bound contexts.
type Enforcer struct {
	policy Policy
}

// New creates an Enforcer using the given Policy.
func New(p Policy) *Enforcer {
	if p.Logger == nil {
		p.Logger = log.Default()
	}
	if p.Connect <= 0 {
		p.Connect = DefaultPolicy().Connect
	}
	if p.Command <= 0 {
		p.Command = DefaultPolicy().Command
	}
	return &Enforcer{policy: p}
}

// ConnectContext returns a context that expires after the connect timeout.
func (e *Enforcer) ConnectContext(parent context.Context, host string) (context.Context, context.CancelFunc) {
	e.policy.Logger.Printf("[timeout] connect deadline %s for host %s", e.policy.Connect, host)
	return context.WithTimeout(parent, e.policy.Connect)
}

// CommandContext returns a context that expires after the command timeout.
func (e *Enforcer) CommandContext(parent context.Context, host string) (context.Context, context.CancelFunc) {
	e.policy.Logger.Printf("[timeout] command deadline %s for host %s", e.policy.Command, host)
	return context.WithTimeout(parent, e.policy.Command)
}

// Wrap executes fn within a command-scoped context. Returns an error if the
// deadline is exceeded or fn itself returns an error.
func (e *Enforcer) Wrap(parent context.Context, host string, fn func(ctx context.Context) error) error {
	ctx, cancel := e.CommandContext(parent, host)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- fn(ctx) }()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("timeout: command exceeded %s on host %s", e.policy.Command, host)
	}
}
