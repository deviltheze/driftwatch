package cache

import (
	"fmt"
	"time"
)

// CheckFunc is the signature of a remote check operation.
type CheckFunc func(host, path string) (string, error)

// Cached wraps a CheckFunc with a Cache, returning a cached checksum
// when available and refreshing it on miss or expiry.
//
// The cache key is "host:path". This reduces SSH calls for files that
// have not changed between scheduler ticks.
func Cached(c *Cache, next CheckFunc) CheckFunc {
	return func(host, path string) (string, error) {
		key := fmt.Sprintf("%s:%s", host, path)
		if val, ok := c.Get(key); ok {
			if s, ok := val.(string); ok {
				return s, nil
			}
		}
		result, err := next(host, path)
		if err != nil {
			return "", err
		}
		c.Set(key, result)
		return result, nil
	}
}

// InvalidateHost removes all cache entries whose key starts with the
// given host prefix, forcing a fresh check on the next call.
func InvalidateHost(c *Cache, host string) {
	prefix := host + ":"
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.items {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(c.items, k)
		}
	}
}

// TTLFromDuration returns a time.Duration suitable for cache TTL based
// on a check interval, using half the interval as a conservative default.
func TTLFromDuration(interval time.Duration) time.Duration {
	if interval <= 0 {
		return 30 * time.Second
	}
	return interval / 2
}
