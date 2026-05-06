package drift

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Report holds the drift check results for a single host.
type Report struct {
	Host      string
	Timestamp time.Time
	Results   []Result
}

// HasDrift returns true if any result in the report indicates drift.
func (r *Report) HasDrift() bool {
	for _, res := range r.Results {
		if res.Drifted {
			return true
		}
	}
	return false
}

// DriftedCount returns the number of drifted results.
func (r *Report) DriftedCount() int {
	count := 0
	for _, res := range r.Results {
		if res.Drifted {
			count++
		}
	}
	return count
}

// Reporter writes drift reports to an output writer.
type Reporter struct {
	out io.Writer
}

// NewReporter creates a new Reporter that writes to the given writer.
func NewReporter(out io.Writer) *Reporter {
	return &Reporter{out: out}
}

// Write formats and writes a Report to the output writer.
func (r *Reporter) Write(report Report) error {
	var sb strings.Builder

	status := "OK"
	if report.HasDrift() {
		status = "DRIFT DETECTED"
	}

	sb.WriteString(fmt.Sprintf("[%s] Host: %s — %s\n",
		report.Timestamp.Format(time.RFC3339),
		report.Host,
		status,
	))

	for _, res := range report.Results {
		if res.Drifted {
			sb.WriteString(fmt.Sprintf("  DRIFT  %s\n    expected: %s\n    actual:   %s\n",
				res.Path, res.ExpectedHash, res.ActualHash))
		} else {
			sb.WriteString(fmt.Sprintf("  OK     %s\n", res.Path))
		}
	}

	_, err := fmt.Fprint(r.out, sb.String())
	return err
}
