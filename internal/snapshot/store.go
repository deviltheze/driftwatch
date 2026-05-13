package snapshot

import (
	"sync"
)

// Store holds the most recent snapshot entry per host+path key.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// NewStore initialises an empty in-memory snapshot store.
func NewStore() *Store {
	return &Store{entries: make(map[string]Entry)}
}

func key(host, path string) string {
	return host + ":" + path
}

// Set stores an entry, overwriting any previous value.
func (s *Store) Set(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[key(e.Host, e.Path)] = e
}

// Get retrieves an entry. The second return value is false if not found.
func (s *Store) Get(host, path string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key(host, path)]
	return e, ok
}

// All returns a copy of every stored entry.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

// Size returns the number of entries currently held.
func (s *Store) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
