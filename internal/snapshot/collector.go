package snapshot

import (
	"fmt"
	"log"
	"time"

	"github.com/user/driftwatch/internal/ssh"
)

// Collector gathers file snapshots from a remote host via SSH.
type Collector struct {
	runner *ssh.Runner
	logger *log.Logger
}

// NewCollector creates a Collector using the provided SSH runner and logger.
func NewCollector(runner *ssh.Runner, logger *log.Logger) *Collector {
	if logger == nil {
		logger = log.Default()
	}
	return &Collector{runner: runner, logger: logger}
}

// Collect fetches the SHA-256 hash of a remote file at the given path.
func (c *Collector) Collect(host, path string) (Entry, error) {
	cmd := fmt.Sprintf("sha256sum %s | awk '{print $1}'", path)
	results, err := c.runner.Run([]string{cmd})
	if err != nil {
		return Entry{}, fmt.Errorf("snapshot collect %s@%s: %w", path, host, err)
	}
	if len(results) == 0 || results[0].Output == "" {
		return Entry{}, fmt.Errorf("snapshot collect %s@%s: empty output", path, host)
	}
	return Entry{
		Host:       host,
		Path:       path,
		Hash:       results[0].Output,
		CapturedAt: time.Now(),
	}, nil
}

// CollectMany collects snapshots for multiple paths on the same host.
func (c *Collector) CollectMany(host string, paths []string) ([]Entry, []error) {
	entries := make([]Entry, 0, len(paths))
	var errs []error
	for _, p := range paths {
		e, err := c.Collect(host, p)
		if err != nil {
			c.logger.Printf("snapshot: skipping %s: %v", p, err)
			errs = append(errs, err)
			continue
		}
		entries = append(entries, e)
	}
	return entries, errs
}
