// Package redact provides utilities for scrubbing sensitive values
// from configuration snapshots and drift reports before they are
// written to logs, audit trails, or alert outputs.
package redact

import (
	"regexp"
	"strings"
)

const placeholder = "[REDACTED]"

// Rule describes a single redaction pattern.
type Rule struct {
	Label   string
	Pattern *regexp.Regexp
}

// Redactor holds a set of rules and applies them to text.
type Redactor struct {
	rules []Rule
}

// DefaultRules returns a Redactor pre-loaded with common sensitive patterns
// such as passwords, tokens, and private keys.
func DefaultRules() *Redactor {
	return &Redactor{
		rules: []Rule{
			{
				Label:   "password",
				Pattern: regexp.MustCompile(`(?i)(password\s*=\s*)\S+`),
			},
			{
				Label:   "token",
				Pattern: regexp.MustCompile(`(?i)(token\s*=\s*)\S+`),
			},
			{
				Label:   "secret",
				Pattern: regexp.MustCompile(`(?i)(secret\s*=\s*)\S+`),
			},
			{
				Label:   "api_key",
				Pattern: regexp.MustCompile(`(?i)(api[_-]?key\s*=\s*)\S+`),
			},
		},
	}
}

// New returns a Redactor with the supplied rules.
func New(rules []Rule) *Redactor {
	return &Redactor{rules: rules}
}

// Apply replaces all sensitive matches in s with the redaction placeholder.
func (r *Redactor) Apply(s string) string {
	for _, rule := range r.rules {
		s = rule.Pattern.ReplaceAllStringFunc(s, func(match string) string {
			idx := strings.Index(strings.ToLower(match), "=")
			if idx == -1 {
				return placeholder
			}
			return match[:idx+1] + placeholder
		})
	}
	return s
}

// Lines applies redaction to each line of a multi-line string independently.
func (r *Redactor) Lines(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = r.Apply(l)
	}
	return strings.Join(lines, "\n")
}
