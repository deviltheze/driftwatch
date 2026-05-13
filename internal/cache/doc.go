// Package cache provides a lightweight in-memory TTL cache used by
// driftwatch to avoid repeated SSH calls for unchanged remote files
// within a configurable window.
//
// Usage:
//
//	c := cache.New(30 * time.Second)
//	c.Set("host:path", checksum)
//	if val, ok := c.Get("host:path"); ok {
//	    // use cached checksum
//	}
//
// Entries are lazily evicted on Get; call Flush to proactively remove
// expired items.
package cache
