// Package audit provides structured, append-only audit logging for
// driftwatch events. Each event is written as a newline-delimited JSON
// record containing a UTC timestamp, event type, optional host, human-
// readable message, and an arbitrary metadata map.
//
// Typical usage:
//
//	logger := audit.New(file)
//	logger.Log(audit.EventDriftDetected, "web-01", "config drift found", map[string]string{
//		"file": "/etc/nginx/nginx.conf",
//	})
package audit
