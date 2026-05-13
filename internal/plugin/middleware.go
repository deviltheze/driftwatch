package plugin

import (
	"fmt"
	"log"
)

// RunAll executes every registered plugin against the given host and output.
// Results are collected into a map keyed by plugin name. Errors from individual
// plugins are logged but do not abort the remaining plugins.
func (r *Registry) RunAll(logger *log.Logger, host, output string) map[string]string {
	if logger == nil {
		logger = log.Default()
	}

	r.mu.RLock()
	plugins := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		plugins = append(plugins, p)
	}
	r.mu.RUnlock()

	results := make(map[string]string, len(plugins))
	for _, p := range plugins {
		annotation, err := p.Check(host, output)
		if err != nil {
			logger.Printf("plugin %q failed for host %q: %v", p.Name, host, err)
			results[p.Name] = fmt.Sprintf("error: %v", err)
			continue
		}
		results[p.Name] = annotation
	}
	return results
}
