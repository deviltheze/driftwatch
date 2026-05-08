package ratelimit_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/ratelimit"
)

func TestNew_DefaultValues(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{})
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	if got := l.Available(); got != 10 {
		t.Errorf("expected 10 default tokens, got %d", got)
	}
}

func TestNew_CustomValues(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{
		MaxTokens:      5,
		RefillInterval: 500 * time.Millisecond,
	})
	if got := l.Available(); got != 5 {
		t.Errorf("expected 5 tokens, got %d", got)
	}
}

func TestAllow_ConsumesToken(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{MaxTokens: 3, RefillInterval: time.Hour})
	if !l.Allow() {
		t.Fatal("expected first Allow() to return true")
	}
	if got := l.Available(); got != 2 {
		t.Errorf("expected 2 tokens after one Allow, got %d", got)
	}
}

func TestAllow_ExhaustsTokens(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{MaxTokens: 2, RefillInterval: time.Hour})
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Error("expected Allow() to return false when tokens exhausted")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{MaxTokens: 2, RefillInterval: 100 * time.Millisecond})
	l.Allow()
	l.Allow()

	// Simulate time passing by waiting for refill.
	time.Sleep(250 * time.Millisecond)

	if !l.Allow() {
		t.Error("expected Allow() to return true after refill interval")
	}
}

func TestAvailable_DoesNotConsumeToken(t *testing.T) {
	l := ratelimit.New(ratelimit.Config{MaxTokens: 4, RefillInterval: time.Hour})
	before := l.Available()
	_ = l.Available()
	after := l.Available()
	if before != after {
		t.Errorf("Available() should not consume tokens: before=%d after=%d", before, after)
	}
}
