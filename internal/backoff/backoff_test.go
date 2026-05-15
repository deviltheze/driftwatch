package backoff_test

import (
	"testing"
	"time"

	"github.com/driftwatch/driftwatch/internal/backoff"
)

func TestDefaultPolicy_Values(t *testing.T) {
	p := backoff.DefaultPolicy()

	if p.InitialInterval != 500*time.Millisecond {
		t.Errorf("expected InitialInterval 500ms, got %v", p.InitialInterval)
	}
	if p.MaxInterval != 30*time.Second {
		t.Errorf("expected MaxInterval 30s, got %v", p.MaxInterval)
	}
	if p.Multiplier != 1.5 {
		t.Errorf("expected Multiplier 1.5, got %v", p.Multiplier)
	}
	if p.Jitter != 0.2 {
		t.Errorf("expected Jitter 0.2, got %v", p.Jitter)
	}
}

func TestNext_ZeroAttempt_ReturnsInitialRange(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     60 * time.Second,
		Multiplier:      2.0,
		Jitter:          0, // no jitter for deterministic test
	}

	d := p.Next(0)
	if d != 1*time.Second {
		t.Errorf("expected 1s for attempt 0, got %v", d)
	}
}

func TestNext_GrowsExponentially(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     60 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}

	d0 := p.Next(0)
	d1 := p.Next(1)
	d2 := p.Next(2)

	if d1 <= d0 {
		t.Errorf("expected attempt 1 (%v) > attempt 0 (%v)", d1, d0)
	}
	if d2 <= d1 {
		t.Errorf("expected attempt 2 (%v) > attempt 1 (%v)", d2, d1)
	}
}

func TestNext_CapsAtMaxInterval(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 1 * time.Second,
		MaxInterval:     4 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}

	// After enough attempts the value should not exceed MaxInterval.
	for attempt := 0; attempt < 20; attempt++ {
		d := p.Next(attempt)
		if d > p.MaxInterval {
			t.Errorf("attempt %d: duration %v exceeds MaxInterval %v", attempt, d, p.MaxInterval)
		}
	}
}

func TestNext_NegativeAttempt_TreatedAsZero(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}

	d := p.Next(-5)
	if d != 500*time.Millisecond {
		t.Errorf("expected 500ms for negative attempt, got %v", d)
	}
}

func TestSteps_ReturnsCorrectLength(t *testing.T) {
	p := backoff.DefaultPolicy()
	steps := p.Steps(5)

	if len(steps) != 5 {
		t.Errorf("expected 5 steps, got %d", len(steps))
	}
}

func TestSteps_ZeroSteps_ReturnsEmpty(t *testing.T) {
	p := backoff.DefaultPolicy()
	steps := p.Steps(0)

	if len(steps) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(steps))
	}
}
