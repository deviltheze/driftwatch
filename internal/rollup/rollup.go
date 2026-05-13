// Package rollup aggregates drift results across multiple hosts into
// a summary suitable for reporting and alerting decisions.
package rollup

import (
	"fmt"
	"strings"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Summary holds aggregated drift information for a single check cycle.
type Summary struct {
	Timestamp   time.Time
	TotalHosts  int
	DriftedHosts []string
	CleanHosts   []string
	Reports     []*drift.Report
}

// HasDrift returns true if any host has drifted.
func (s *Summary) HasDrift() bool {
	return len(s.DriftedHosts) > 0
}

// DriftRatio returns the fraction of hosts that have drifted (0.0–1.0).
func (s *Summary) DriftRatio() float64 {
	if s.TotalHosts == 0 {
		return 0.0
	}
	return float64(len(s.DriftedHosts)) / float64(s.TotalHosts)
}

// String returns a human-readable one-line summary.
func (s *Summary) String() string {
	if !s.HasDrift() {
		return fmt.Sprintf("[%s] all %d host(s) clean",
			s.Timestamp.Format(time.RFC3339), s.TotalHosts)
	}
	return fmt.Sprintf("[%s] drift detected on %d/%d host(s): %s",
		s.Timestamp.Format(time.RFC3339),
		len(s.DriftedHosts),
		s.TotalHosts,
		strings.Join(s.DriftedHosts, ", "))
}

// Aggregate combines a slice of drift reports into a Summary.
func Aggregate(reports []*drift.Report) *Summary {
	s := &Summary{
		Timestamp: time.Now().UTC(),
		TotalHosts: len(reports),
		Reports:    reports,
	}
	for _, r := range reports {
		if r.HasDrift() {
			s.DriftedHosts = append(s.DriftedHosts, r.Host)
		} else {
			s.CleanHosts = append(s.CleanHosts, r.Host)
		}
	}
	return s
}
