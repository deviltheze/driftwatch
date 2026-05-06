package drift

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/yourusername/driftwatch/internal/ssh"
)

// FileState represents the state of a file on a remote host.
type FileState struct {
	Host     string
	Path     string
	Checksum string
	Missing  bool
}

// Result holds the drift detection result for a single host+file pair.
type Result struct {
	Host     string
	Path     string
	Drifted  bool
	Baseline string
	Actual   string
	Err      error
}

// Checker detects configuration drift by comparing remote file checksums
// against a known baseline.
type Checker struct {
	runner *ssh.Runner
}

// NewChecker creates a new Checker using the provided SSH runner.
func NewChecker(runner *ssh.Runner) *Checker {
	return &Checker{runner: runner}
}

// FetchChecksum retrieves the SHA-256 checksum of a file on the remote host.
func (c *Checker) FetchChecksum(path string) (string, error) {
	results, err := c.runner.Run([]string{fmt.Sprintf("sha256sum %s", path)})
	if err != nil {
		return "", fmt.Errorf("failed to run checksum command: %w", err)
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no results returned for path %s", path)
	}
	result := results[0]
	if result.ExitCode != 0 {
		return "", fmt.Errorf("checksum command failed: %s", result.Stderr)
	}
	parts := strings.Fields(result.Stdout)
	if len(parts) == 0 {
		return "", fmt.Errorf("unexpected sha256sum output: %q", result.Stdout)
	}
	return parts[0], nil
}

// CheckDrift compares the remote file checksum against the provided baseline.
func (c *Checker) CheckDrift(host, path, baseline string) Result {
	actual, err := c.FetchChecksum(path)
	if err != nil {
		return Result{Host: host, Path: path, Err: err}
	}
	return Result{
		Host:     host,
		Path:     path,
		Drifted:  actual != baseline,
		Baseline: baseline,
		Actual:   actual,
	}
}

// HashContent computes the SHA-256 hex digest of the given content string.
// Useful for generating baselines from known-good file contents.
func HashContent(content string) string {
	sum := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", sum)
}
