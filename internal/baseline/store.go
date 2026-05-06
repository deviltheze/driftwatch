package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Entry represents a stored baseline hash for a host+path combination.
type Entry struct {
	Host string `json:"host"`
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// Store persists baseline hashes to disk for drift comparison.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Entry
}

// NewStore creates or loads a baseline store from the given file path.
func NewStore(path string) (*Store, error) {
	s := &Store{
		path:    path,
		entries: make(map[string]Entry),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("baseline: load store: %w", err)
	}
	return s, nil
}

func storeKey(host, path string) string {
	return host + "::" + path
}

// Get retrieves the stored hash for a host+path pair.
func (s *Store) Get(host, path string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[storeKey(host, path)]
	return e, ok
}

// Set stores or updates the hash for a host+path pair and persists to disk.
func (s *Store) Set(host, path, hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[storeKey(host, path)] = Entry{Host: host, Path: path, Hash: hash}
	return s.save()
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	var list []Entry
	if err := json.NewDecoder(f).Decode(&list); err != nil {
		return fmt.Errorf("baseline: decode: %w", err)
	}
	for _, e := range list {
		s.entries[storeKey(e.Host, e.Path)] = e
	}
	return nil
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("baseline: mkdir: %w", err)
	}
	tmp := s.path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("baseline: create tmp: %w", err)
	}
	list := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		list = append(list, e)
	}
	if err := json.NewEncoder(f).Encode(list); err != nil {
		f.Close()
		return fmt.Errorf("baseline: encode: %w", err)
	}
	f.Close()
	return os.Rename(tmp, s.path)
}
