package baseline

import "fmt"

// Manager wraps a Store and provides higher-level operations used by the
// drift engine: recording a new baseline and comparing against an existing one.
type Manager struct {
	store *Store
}

// NewManager creates a Manager backed by the given Store.
func NewManager(store *Store) *Manager {
	return &Manager{store: store}
}

// Record saves hash as the current baseline for host+path.
func (m *Manager) Record(host, path, hash string) error {
	if host == "" || path == "" {
		return fmt.Errorf("baseline: host and path must not be empty")
	}
	return m.store.Set(host, path, hash)
}

// Compare checks whether hash differs from the stored baseline.
// Returns (drifted, baselineExists, error).
func (m *Manager) Compare(host, path, hash string) (bool, bool, error) {
	if host == "" || path == "" {
		return false, false, fmt.Errorf("baseline: host and path must not be empty")
	}
	e, ok := m.store.Get(host, path)
	if !ok {
		return false, false, nil
	}
	return e.Hash != hash, true, nil
}

// Snapshot returns all stored entries as a flat slice.
func (m *Manager) Snapshot() []Entry {
	m.store.mu.RLock()
	defer m.store.mu.RUnlock()
	out := make([]Entry, 0, len(m.store.entries))
	for _, e := range m.store.entries {
		out = append(out, e)
	}
	return out
}
