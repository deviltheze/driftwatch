package retry

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"
)

func silentLogger() *log.Logger {
	return log.New(os.Stderr, "[test] ", 0)
}

func TestDefaultPolicy_Values(t *testing.T) {
	p := DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.InitialDelay != 500*time.Millisecond {
		t.Errorf("unexpected InitialDelay: %s", p.InitialDelay)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", p.Multiplier)
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := Do(DefaultPolicy(), nil, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	p := Policy{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Multiplier: 2.0}
	calls := 0
	err := Do(p, silentLogger(), func() error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success on third attempt, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	p := Policy{MaxAttempts: 2, InitialDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond, Multiplier: 2.0}
	calls := 0
	err := Do(p, nil, func() error {
		calls++
		return errors.New("permanent")
	})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestDo_ZeroMaxAttempts_RunsOnce(t *testing.T) {
	p := Policy{MaxAttempts: 0, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 2.0}
	calls := 0
	_ = Do(p, nil, func() error {
		calls++
		return errors.New("fail")
	})
	if calls != 1 {
		t.Errorf("expected 1 call for zero MaxAttempts, got %d", calls)
	}
}

func TestDelay_RespectsMaxDelay(t *testing.T) {
	max := 5 * time.Millisecond
	d := delay(time.Millisecond, 2.0, 10, max)
	if d > max {
		t.Errorf("delay %s exceeds max %s", d, max)
	}
}
