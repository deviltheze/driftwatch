// Package ratelimit implements a token-bucket rate limiter used to control
// the frequency of SSH-based drift checks against remote hosts.
//
// Usage:
//
//	limiter := ratelimit.New(ratelimit.Config{
//		MaxTokens:      10,
//		RefillInterval: time.Second,
//	})
//
//	if limiter.Allow() {
//		// perform check
//	}
//
// The limiter is safe for concurrent use.
package ratelimit
