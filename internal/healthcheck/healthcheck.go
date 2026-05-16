// Package healthcheck provides SSH-based liveness probing for remote hosts.
package healthcheck

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Status represents the health state of a remote host.
type Status int

const (
	StatusUnknown Status = iota
	StatusHealthy
	StatusUnhealthy
)

func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusUnhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

// Result holds the outcome of a single host probe.
type Result struct {
	Host      string
	Status    Status
	Latency   time.Duration
	Err       error
	CheckedAt time.Time
}

// Prober defines the interface for checking host reachability.
type Prober interface {
	Probe(ctx context.Context, host string) Result
}

// Runner executes health probes against a set of hosts.
type Runner struct {
	prober  Prober
	timeout time.Duration
	logger  *log.Logger
}

// New creates a Runner with the given Prober, timeout, and logger.
// If logger is nil a default is used.
func New(prober Prober, timeout time.Duration, logger *log.Logger) *Runner {
	if logger == nil {
		logger = log.Default()
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Runner{prober: prober, timeout: timeout, logger: logger}
}

// Check probes all hosts concurrently and returns a slice of Results.
func (r *Runner) Check(ctx context.Context, hosts []string) []Result {
	results := make([]Result, len(hosts))
	type indexed struct {
		i int
		res Result
	}
	ch := make(chan indexed, len(hosts))
	for i, h := range hosts {
		go func(idx int, host string) {
			tctx, cancel := context.WithTimeout(ctx, r.timeout)
			defer cancel()
			res := r.prober.Probe(tctx, host)
			if res.Err != nil {
				r.logger.Printf("healthcheck: host %s unhealthy: %v", host, res.Err)
			}
			ch <- indexed{idx, res}
		}(i, h)
	}
	for range hosts {
		v := <-ch
		results[v.i] = v.res
	}
	return results
}

// Healthy returns true when all results report StatusHealthy.
func Healthy(results []Result) bool {
	for _, r := range results {
		if r.Status != StatusHealthy {
			return false
		}
	}
	return true
}

// Summary returns a human-readable overview of the results.
func Summary(results []Result) string {
	healthy, unhealthy := 0, 0
	for _, r := range results {
		if r.Status == StatusHealthy {
			healthy++
		} else {
			unhealthy++
		}
	}
	return fmt.Sprintf("healthcheck: %d healthy, %d unhealthy", healthy, unhealthy)
}
