// Package filter provides utilities for selecting and excluding hosts
// from drift checks based on label selectors and tag patterns.
package filter

import "strings"

// HostFilter selects hosts based on include/exclude tag rules.
type HostFilter struct {
	include []string
	exclude []string
}

// New returns a HostFilter with the given include and exclude tag lists.
// An empty include list means all hosts are included by default.
func New(include, exclude []string) *HostFilter {
	return &HostFilter{
		include: include,
		exclude: exclude,
	}
}

// Match reports whether a host with the given tags passes the filter.
// A host passes when it matches at least one include tag (or include is empty)
// and does not match any exclude tag.
func (f *HostFilter) Match(tags []string) bool {
	if matchesAny(tags, f.exclude) {
		return false
	}
	if len(f.include) == 0 {
		return true
	}
	return matchesAny(tags, f.include)
}

// Filter returns only the host names whose tags satisfy the filter.
// hosts maps host name → its tags.
func (f *HostFilter) Filter(hosts map[string][]string) []string {
	var result []string
	for name, tags := range hosts {
		if f.Match(tags) {
			result = append(result, name)
		}
	}
	return result
}

// matchesAny returns true if any tag in haystack equals any pattern (case-insensitive).
func matchesAny(tags, patterns []string) bool {
	for _, t := range tags {
		for _, p := range patterns {
			if strings.EqualFold(t, p) {
				return true
			}
		}
	}
	return false
}
