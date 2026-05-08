package metrics

import (
	"time"
)

// CheckFunc is a function that performs a drift check and returns whether
// drift was detected and any error encountered.
type CheckFunc func(host string) (drifted bool, err error)

// Instrumented wraps a CheckFunc with metrics recording. It records the
// check duration, increments counters, and updates the last-check timestamp
// on the provided Collector.
func Instrumented(c *Collector, fn CheckFunc) CheckFunc {
	return func(host string) (bool, error) {
		start := time.Now()
		drifted, err := fn(host)
		duration := time.Since(start)
		c.RecordCheck(host, drifted, duration)
		return drifted, err
	}
}
