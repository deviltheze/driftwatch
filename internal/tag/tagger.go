// Package tag provides host tagging and label resolution for driftwatch.
// Tags are used by the filter and plugin systems to target specific hosts.
package tag

import "strings"

// Tagger resolves and normalises tags for a given host.
type Tagger struct {
	// global holds tags that apply to every host.
	global []string
	// overrides maps a hostname to its explicit tag list.
	overrides map[string][]string
}

// New returns a Tagger with the supplied global tags.
// Per-host overrides can be added via AddOverride.
func New(globalTags []string) *Tagger {
	norm := make([]string, 0, len(globalTags))
	for _, t := range globalTags {
		if v := normalise(t); v != "" {
			norm = append(norm, v)
		}
	}
	return &Tagger{
		global:    norm,
		overrides: make(map[string][]string),
	}
}

// AddOverride registers a set of tags that are merged with the global tags
// for the named host. Duplicate tags are deduplicated.
func (t *Tagger) AddOverride(host string, tags []string) {
	norm := make([]string, 0, len(tags))
	for _, tag := range tags {
		if v := normalise(tag); v != "" {
			norm = append(norm, v)
		}
	}
	t.overrides[host] = norm
}

// Resolve returns the deduplicated, sorted tag set for host.
// Global tags are always included; per-host overrides are merged on top.
func (t *Tagger) Resolve(host string) []string {
	seen := make(map[string]struct{}, len(t.global))
	result := make([]string, 0, len(t.global))

	for _, g := range t.global {
		if _, ok := seen[g]; !ok {
			seen[g] = struct{}{}
			result = append(result, g)
		}
	}

	for _, o := range t.overrides[host] {
		if _, ok := seen[o]; !ok {
			seen[o] = struct{}{}
			result = append(result, o)
		}
	}

	return result
}

// HasTag reports whether host carries the given tag (case-insensitive).
func (t *Tagger) HasTag(host, tag string) bool {
	target := normalise(tag)
	for _, v := range t.Resolve(host) {
		if v == target {
			return true
		}
	}
	return false
}

// normalise lower-cases and trims whitespace from a tag value.
func normalise(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
