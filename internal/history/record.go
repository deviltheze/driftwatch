// Package history tracks drift check results over time,
// allowing driftwatch to surface trends and repeated violations.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single recorded drift-check outcome for one host.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Drifted   bool      `json:"drifted"`
	Details   []string  `json:"details,omitempty"`
}

// Recorder persists drift history to a JSON file.
type Recorder struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// NewRecorder loads (or creates) a history file at the given path.
func NewRecorder(path string) (*Recorder, error) {
	r := &Recorder{path: path}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return r, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &r.entries); err != nil {
		return nil, err
	}
	return r, nil
}

// Record appends a new entry and flushes to disk.
func (r *Recorder) Record(e Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	r.entries = append(r.entries, e)
	return r.flush()
}

// Entries returns a copy of all recorded entries.
func (r *Recorder) Entries() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Entry, len(r.entries))
	copy(out, r.entries)
	return out
}

// DriftedHosts returns the distinct set of hosts that have drifted at least once.
func (r *Recorder) DriftedHosts() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	seen := make(map[string]struct{})
	var hosts []string
	for _, e := range r.entries {
		if e.Drifted {
			if _, ok := seen[e.Host]; !ok {
				seen[e.Host] = struct{}{}
				hosts = append(hosts, e.Host)
			}
		}
	}
	return hosts
}

func (r *Recorder) flush() error {
	data, err := json.MarshalIndent(r.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0o644)
}
