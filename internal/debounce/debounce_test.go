package debounce

import (
	"testing"
	"time"
)

func TestNew_NotNil(t *testing.T) {
	d := New(5 * time.Second)
	if d == nil {
		t.Fatal("expected non-nil Debouncer")
	}
}

func TestAllow_FirstEvent_Passes(t *testing.T) {
	d := New(10 * time.Second)
	if !d.Allow("host-a") {
		t.Error("first event should always be allowed")
	}
}

func TestAllow_DuplicateWithinWindow_Suppressed(t *testing.T) {
	d := New(10 * time.Second)
	d.Allow("host-a")
	if d.Allow("host-a") {
		t.Error("duplicate within window should be suppressed")
	}
}

func TestAllow_AfterWindow_Passes(t *testing.T) {
	now := time.Now()
	d := New(5 * time.Second)
	d.now = func() time.Time { return now }

	d.Allow("host-b")

	// advance clock beyond window
	d.now = func() time.Time { return now.Add(6 * time.Second) }

	if !d.Allow("host-b") {
		t.Error("event after window expiry should be allowed")
	}
}

func TestAllow_ZeroWindow_AllPass(t *testing.T) {
	d := New(0)
	d.Allow("host-c")
	if !d.Allow("host-c") {
		t.Error("zero window should allow all events")
	}
}

func TestAllow_DifferentHosts_Independent(t *testing.T) {
	d := New(10 * time.Second)
	d.Allow("host-x")
	if !d.Allow("host-y") {
		t.Error("different hosts should be tracked independently")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New(10 * time.Second)
	d.Allow("host-d")
	d.Reset("host-d")
	if !d.Allow("host-d") {
		t.Error("after reset, event should be allowed again")
	}
}

func TestStats_ReturnsCount(t *testing.T) {
	d := New(10 * time.Second)
	d.Allow("host-e")
	d.Allow("host-e") // suppressed
	d.Allow("host-e") // suppressed

	ev, ok := d.Stats("host-e")
	if !ok {
		t.Fatal("expected stats to exist for host-e")
	}
	if ev.Count != 3 {
		t.Errorf("expected count 3, got %d", ev.Count)
	}
}

func TestStats_MissingHost_ReturnsFalse(t *testing.T) {
	d := New(10 * time.Second)
	_, ok := d.Stats("nonexistent")
	if ok {
		t.Error("expected false for unknown host")
	}
}
