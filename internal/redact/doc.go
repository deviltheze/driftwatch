// Package redact provides pattern-based scrubbing of sensitive configuration
// values before they appear in logs, audit records, drift reports, or alert
// payloads.
//
// Usage:
//
//	r := redact.DefaultRules()
//	safeOutput := r.Lines(rawConfigContent)
//
// DefaultRules covers common patterns such as passwords, tokens, secrets, and
// API keys.  Custom rules can be supplied via New when the defaults are
// insufficient for a particular environment.
//
// Redaction is intentionally conservative: only the value portion of a
// key=value pair is replaced with the [REDACTED] placeholder so that key
// names remain visible for diagnostic purposes.
package redact
