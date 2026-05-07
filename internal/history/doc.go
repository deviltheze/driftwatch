// Package history provides a lightweight append-only log of drift-check
// outcomes. Each time the drift engine completes a scan, callers may
// record an Entry (host, timestamp, whether drift was detected, and any
// human-readable details). Entries are persisted as a JSON array on disk
// so that driftwatch can surface trends across restarts — for example,
// identifying hosts that drift repeatedly or alerting only after N
// consecutive failures.
//
// Typical usage:
//
//	rec, err := history.NewRecorder("/var/lib/driftwatch/history.json")
//	if err != nil { ... }
//
//	rec.Record(history.Entry{
//		Host:    "web-01",
//		Drifted: true,
//		Details: []string{"/etc/nginx/nginx.conf changed"},
//	})
package history
