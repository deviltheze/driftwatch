// Package snapshot captures and compares remote file states.
package snapshot

import (
	"fmt"
	"time"
)

// Entry represents a single captured file state from a remote host.
type Entry struct {
	Host      string
	Path      string
	Hash      string
	CapturedAt time.Time
}

// Diff represents a detected change between two snapshots.
type Diff struct {
	Host     string
	Path     string
	OldHash  string
	NewHash  string
	DetectedAt time.Time
}

// Changed returns true if the hashes differ.
func (d Diff) Changed() bool {
	return d.OldHash != d.NewHash
}

// String returns a human-readable summary of the diff.
func (d Diff) String() string {
	return fmt.Sprintf("[%s] %s: %s -> %s", d.Host, d.Path, d.OldHash[:8], d.NewHash[:8])
}

// Compare takes a previous and current entry and returns a Diff if they differ.
func Compare(prev, curr Entry) (Diff, bool) {
	if prev.Hash == curr.Hash {
		return Diff{}, false
	}
	return Diff{
		Host:       curr.Host,
		Path:       curr.Path,
		OldHash:    prev.Hash,
		NewHash:    curr.Hash,
		DetectedAt: time.Now(),
	}, true
}
