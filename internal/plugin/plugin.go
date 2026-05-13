// Package plugin provides a lightweight hook system for extending
// driftwatch with custom check logic at runtime.
package plugin

import (
	"errors"
	"fmt"
	"sync"
)

// CheckFn is a function that receives raw command output from a remote host
// and returns an annotation or an error.
type CheckFn func(host, output string) (string, error)

// Plugin represents a named extension point.
type Plugin struct {
	Name string
	Check CheckFn
}

// Registry holds registered plugins and provides thread-safe access.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
}

// NewRegistry returns an initialised, empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

// Register adds a plugin to the registry. Returns an error if the name is
// already taken or the CheckFn is nil.
func (r *Registry) Register(p Plugin) error {
	if p.Name == "" {
		return errors.New("plugin name must not be empty")
	}
	if p.Check == nil {
		return fmt.Errorf("plugin %q: CheckFn must not be nil", p.Name)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[p.Name]; exists {
		return fmt.Errorf("plugin %q already registered", p.Name)
	}
	r.plugins[p.Name] = p
	return nil
}

// Get returns the plugin with the given name and a boolean indicating
// whether it was found.
func (r *Registry) Get(name string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plugins[name]
	return p, ok
}

// Names returns a sorted-order snapshot of all registered plugin names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.plugins))
	for n := range r.plugins {
		names = append(names, n)
	}
	return names
}

// Len returns the number of registered plugins.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.plugins)
}
