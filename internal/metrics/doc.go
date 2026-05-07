// Package metrics provides a lightweight, thread-safe in-process metrics
// collector for the driftwatch daemon.
//
// It intentionally avoids external dependencies; counters are accumulated in
// memory and can be written to any io.Writer for logging or inspection.
//
// Usage:
//
//	col := metrics.NewCollector()
//	col.RecordCheck(hostsInRun, driftedCount)
//	col.Write(os.Stdout)
package metrics
