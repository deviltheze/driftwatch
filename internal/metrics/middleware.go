package metrics

import (
	"log"
	"time"
)

// RunFunc is the signature of a single drift-check execution.
type RunFunc func() (hostsChecked int64, drifted int64, err error)

// Instrumented wraps fn so that every invocation is recorded in col.
// If fn returns an error it is logged via logger (or the default logger when
// nil) and the counters are still updated with whatever partial results
// fn returned before failing.
func Instrumented(col *Collector, logger *log.Logger, fn RunFunc) RunFunc {
	if col == nil {
		panic("metrics: Instrumented requires a non-nil Collector")
	}
	if logger == nil {
		logger = log.Default()
	}
	return func() (int64, int64, error) {
		start := time.Now()
		hosts, drifted, err := fn()
		col.RecordCheck(hosts, drifted)
		if err != nil {
			logger.Printf("metrics: check error after %s: %v", time.Since(start).Round(time.Millisecond), err)
		}
		return hosts, drifted, err
	}
}
