package cache

import (
	"errors"
	"testing"
	"time"
)

func TestCached_HitSkipsNext(t *testing.T) {
	c := New(time.Minute)
	calls := 0
	next := func(host, path string) (string, error) {
		calls++
		return "abc123", nil
	}
	wrapped := Cached(c, next)

	// first call populates cache
	if _, err := wrapped("host1", "/etc/hosts"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// second call should hit cache
	val, err := wrapped("host1", "/etc/hosts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "abc123" {
		t.Fatalf("expected abc123, got %s", val)
	}
	if calls != 1 {
		t.Fatalf("expected 1 upstream call, got %d", calls)
	}
}

func TestCached_MissCallsNext(t *testing.T) {
	c := New(time.Minute)
	next := func(host, path string) (string, error) {
		return "deadbeef", nil
	}
	wrapped := Cached(c, next)
	val, err := wrapped("host2", "/etc/passwd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "deadbeef" {
		t.Fatalf("expected deadbeef, got %s", val)
	}
}

func TestCached_PropagatesError(t *testing.T) {
	c := New(time.Minute)
	next := func(host, path string) (string, error) {
		return "", errors.New("ssh timeout")
	}
	wrapped := Cached(c, next)
	_, err := wrapped("host3", "/etc/os-release")
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestInvalidateHost_RemovesMatchingKeys(t *testing.T) {
	c := New(time.Minute)
	c.Set("host1:/etc/hosts", "v1")
	c.Set("host1:/etc/passwd", "v2")
	c.Set("host2:/etc/hosts", "v3")
	InvalidateHost(c, "host1")
	if _, ok := c.Get("host1:/etc/hosts"); ok {
		t.Error("expected host1:/etc/hosts to be evicted")
	}
	if _, ok := c.Get("host1:/etc/passwd"); ok {
		t.Error("expected host1:/etc/passwd to be evicted")
	}
	if _, ok := c.Get("host2:/etc/hosts"); !ok {
		t.Error("expected host2:/etc/hosts to remain")
	}
}

func TestTTLFromDuration_HalfInterval(t *testing.T) {
	ttl := TTLFromDuration(60 * time.Second)
	if ttl != 30*time.Second {
		t.Fatalf("expected 30s, got %v", ttl)
	}
}

func TestTTLFromDuration_ZeroUsesDefault(t *testing.T) {
	ttl := TTLFromDuration(0)
	if ttl != 30*time.Second {
		t.Fatalf("expected default 30s, got %v", ttl)
	}
}
