package cache

import (
	"testing"
	"time"
)

func TestNew_NotNil(t *testing.T) {
	c := New(time.Minute)
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestSet_Get_Hit(t *testing.T) {
	c := New(time.Minute)
	c.Set("key", "value")
	val, ok := c.Get("key")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if val != "value" {
		t.Fatalf("expected 'value', got %v", val)
	}
}

func TestGet_Miss(t *testing.T) {
	c := New(time.Minute)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestSet_Expired(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Set("key", "value")
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestSetTTL_OverridesDefault(t *testing.T) {
	c := New(time.Hour)
	c.SetTTL("key", 42, 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected custom TTL expiry")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c := New(time.Minute)
	c.Set("key", "value")
	c.Delete("key")
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected deleted entry to be missing")
	}
}

func TestFlush_RemovesExpired(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Set("a", 1)
	c.Set("b", 2)
	time.Sleep(5 * time.Millisecond)
	c.SetTTL("c", 3, time.Hour) // not expired
	c.Flush()
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry after flush, got %d", c.Len())
	}
}

func TestLen_ReturnsTotal(t *testing.T) {
	c := New(time.Minute)
	c.Set("x", 1)
	c.Set("y", 2)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}

func TestZeroTTL_NeverExpires(t *testing.T) {
	c := New(0)
	c.Set("key", "forever")
	time.Sleep(5 * time.Millisecond)
	_, ok := c.Get("key")
	if !ok {
		t.Fatal("expected zero-TTL entry to never expire")
	}
}
