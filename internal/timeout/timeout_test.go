package timeout_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/timeout"
)

var silent = log.New(io.Discard, "", 0)

func TestDefaultPolicy_Values(t *testing.T) {
	p := timeout.DefaultPolicy()
	if p.Connect != 10*time.Second {
		t.Errorf("expected connect 10s, got %s", p.Connect)
	}
	if p.Command != 30*time.Second {
		t.Errorf("expected command 30s, got %s", p.Command)
	}
}

func TestNew_NotNil(t *testing.T) {
	e := timeout.New(timeout.Policy{Logger: silent})
	if e == nil {
		t.Fatal("expected non-nil Enforcer")
	}
}

func TestNew_ZeroDurations_UsesDefaults(t *testing.T) {
	// zero durations should fall back to defaults without panicking
	e := timeout.New(timeout.Policy{Logger: silent})
	ctx, cancel := e.ConnectContext(context.Background(), "host1")
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if time.Until(deadline) <= 0 {
		t.Error("deadline should be in the future")
	}
}

func TestConnectContext_HasDeadline(t *testing.T) {
	p := timeout.Policy{Connect: 5 * time.Second, Command: 10 * time.Second, Logger: silent}
	e := timeout.New(p)
	ctx, cancel := e.ConnectContext(context.Background(), "srv1")
	defer cancel()
	_, ok := ctx.Deadline()
	if !ok {
		t.Error("ConnectContext should set a deadline")
	}
}

func TestCommandContext_HasDeadline(t *testing.T) {
	p := timeout.Policy{Connect: 5 * time.Second, Command: 10 * time.Second, Logger: silent}
	e := timeout.New(p)
	ctx, cancel := e.CommandContext(context.Background(), "srv1")
	defer cancel()
	_, ok := ctx.Deadline()
	if !ok {
		t.Error("CommandContext should set a deadline")
	}
}

func TestWrap_SuccessBeforeDeadline(t *testing.T) {
	p := timeout.Policy{Connect: 5 * time.Second, Command: 2 * time.Second, Logger: silent}
	e := timeout.New(p)
	err := e.Wrap(context.Background(), "host1", func(_ context.Context) error {
		return nil
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestWrap_PropagatesFnError(t *testing.T) {
	p := timeout.Policy{Connect: 5 * time.Second, Command: 2 * time.Second, Logger: silent}
	e := timeout.New(p)
	sentinel := errors.New("ssh failure")
	err := e.Wrap(context.Background(), "host1", func(_ context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestWrap_ExceedsDeadline(t *testing.T) {
	p := timeout.Policy{Connect: 1 * time.Second, Command: 50 * time.Millisecond, Logger: silent}
	e := timeout.New(p)
	err := e.Wrap(context.Background(), "slow-host", func(_ context.Context) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}
