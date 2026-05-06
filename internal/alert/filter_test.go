package alert_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/alert"
)

func makeAlert(host string) *alert.Alert {
	return &alert.Alert{
		Level:     alert.LevelWarn,
		Host:      host,
		Message:   "test",
		Timestamp: time.Now().UTC(),
	}
}

func TestNewSuppression_NotNil(t *testing.T) {
	s := alert.NewSuppression(time.Minute)
	if s == nil {
		t.Fatal("expected non-nil suppression")
	}
}

func TestAllow_FirstAlert_Allowed(t *testing.T) {
	s := alert.NewSuppression(time.Minute)
	a := makeAlert("web-01")
	if !s.Allow(a) {
		t.Error("expected first alert to be allowed")
	}
}

func TestAllow_DuplicateWithinCooldown_Suppressed(t *testing.T) {
	s := alert.NewSuppression(time.Hour)
	a := makeAlert("web-01")
	s.Allow(a)
	if s.Allow(a) {
		t.Error("expected duplicate alert within cooldown to be suppressed")
	}
}

func TestAllow_AfterCooldown_Allowed(t *testing.T) {
	s := alert.NewSuppression(time.Millisecond)
	a := makeAlert("web-02")
	s.Allow(a)
	time.Sleep(5 * time.Millisecond)
	if !s.Allow(a) {
		t.Error("expected alert after cooldown to be allowed")
	}
}

func TestAllow_NilAlert_NotAllowed(t *testing.T) {
	s := alert.NewSuppression(time.Minute)
	if s.Allow(nil) {
		t.Error("expected nil alert to not be allowed")
	}
}

func TestAllow_DifferentHosts_BothAllowed(t *testing.T) {
	s := alert.NewSuppression(time.Hour)
	a1 := makeAlert("host-a")
	a2 := makeAlert("host-b")
	if !s.Allow(a1) {
		t.Error("expected host-a to be allowed")
	}
	if !s.Allow(a2) {
		t.Error("expected host-b to be allowed")
	}
}

func TestReset_ClearsSuppression(t *testing.T) {
	s := alert.NewSuppression(time.Hour)
	a := makeAlert("web-03")
	s.Allow(a)
	s.Reset()
	if !s.Allow(a) {
		t.Error("expected alert to be allowed after reset")
	}
}
