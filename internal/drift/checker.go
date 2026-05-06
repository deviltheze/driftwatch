package drift

import (
	"crypto/sha256"
	"fmt"

	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/ssh"
)

// Result holds the drift check outcome for a single file on a single host.
type Result struct {
	Host     string
	FilePath string
	Expected string
	Actual   string
	Err      error
}

// Drifted returns true if the actual content differs from expected or an error occurred.
func (r Result) Drifted() bool {
	return r.Err != nil || r.Expected != r.Actual
}

// Report aggregates results from all hosts.
type Report struct {
	Results []Result
}

// HasDrift returns true if any result in the report has drifted.
func (r Report) HasDrift() bool {
	for _, res := range r.Results {
		if res.Drifted() {
			return true
		}
	}
	return false
}

// DriftedCount returns the number of drifted results.
func (r Report) DriftedCount() int {
	count := 0
	for _, res := range r.Results {
		if res.Drifted() {
			count++
		}
	}
	return count
}

// WriteFunc is a function used by Reporter to output a report.
type WriteFunc func(Report) error

// Checker performs file hash comparisons against expected baselines.
type Checker struct {
	baselines map[string]string
}

// NewChecker creates a Checker with the given baseline hashes keyed by file path.
func NewChecker(baselines map[string]string) *Checker {
	if baselines == nil {
		baselines = make(map[string]string)
	}
	return &Checker{baselines: baselines}
}

// HashContent returns a SHA-256 hex digest of the given content.
func HashContent(content string) string {
	sum := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", sum)
}

// CheckFiles runs remote cat commands for each watched file on the given host.
func (c *Checker) CheckFiles(host config.Host, runner *ssh.Runner) []Result {
	var results []Result
	for _, fp := range host.WatchFiles {
		cmdResult := runner.Run([]string{"cat", fp})
		actualHash := HashContent(cmdResult.Stdout)
		expected := c.baselines[fp]
		results = append(results, Result{
			Host:     host.Address,
			FilePath: fp,
			Expected: expected,
			Actual:   actualHash,
			Err:      cmdResult.Err,
		})
	}
	return results
}

// BuildReport wraps a slice of results into a Report.
func (c *Checker) BuildReport(results []Result) Report {
	return Report{Results: results}
}
