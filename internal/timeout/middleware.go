package timeout

import (
	"context"
	"fmt"
)

// CheckFn is the signature for a single drift check operation.
type CheckFn func(ctx context.Context, host string) error

// Instrumented wraps a CheckFn so that every invocation is bounded by the
// Enforcer's command timeout. If the underlying check exceeds the deadline,
// a descriptive error is returned instead of blocking indefinitely.
//
// Usage:
//
//	guarded := timeout.Instrumented(enforcer, originalCheck)
//	err := guarded(ctx, "10.0.0.1")
func Instrumented(e *Enforcer, fn CheckFn) CheckFn {
	if e == nil {
		panic("timeout: Instrumented requires a non-nil Enforcer")
	}
	return func(ctx context.Context, host string) error {
		if host == "" {
			return fmt.Errorf("timeout: host must not be empty")
		}
		return e.Wrap(ctx, host, func(inner context.Context) error {
			return fn(inner, host)
		})
	}
}
