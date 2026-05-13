// Package cache provides a lightweight in-memory TTL cache for storing
// transient check results to reduce redundant SSH round-trips.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached value and its expiry time.
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Expired reports whether the entry has passed its TTL.
func (e Entry) Expired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache is a thread-safe in-memory store with per-entry TTL.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]Entry
	default TTL time.Duration
}

// New returns a Cache with the given default TTL.
// If ttl is zero, entries never expire.
func New(ttl time.Duration) *Cache {
	return &Cache{
		items:      make(map[string]Entry),
		defaultTTL: ttl,
	}
}

// Set stores value under key, expiring after the cache's default TTL.
func (c *Cache) Set(key string, value interface{}) {
	c.SetTTL(key, value, c.defaultTTL)
}

// SetTTL stores value under key with an explicit TTL.
func (c *Cache) SetTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.items[key] = Entry{Value: value, ExpiresAt: exp}
}

// Get retrieves a value by key. Returns (value, true) on a valid hit.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if !e.ExpiresAt.IsZero() && e.Expired() {
		return nil, false
	}
	return e.Value, true
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all expired entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, e := range c.items {
		if !e.ExpiresAt.IsZero() && e.Expired() {
			delete(c.items, k)
		}
	}
}

// Len returns the total number of entries, including expired ones.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
