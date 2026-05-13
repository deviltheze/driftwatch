package audit

import (
	"fmt"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// CheckFunc is the signature for a drift engine run.
type CheckFunc func() (*drift.Report, error)

// Instrumented wraps a CheckFunc, emitting audit events before and after
// execution and recording any detected drift per host.
func Instrumented(logger *Logger, host string, fn CheckFunc) CheckFunc {
	return func() (*drift.Report, error) {
		start := time.Now()
		_ = logger.Log(EventCheckStarted, host, "drift check started", nil)

		report, err := fn()
		if err != nil {
			_ = logger.Log(EventCheckFinished, host, fmt.Sprintf("check error: %v", err), map[string]string{
				"duration_ms": fmt.Sprintf("%d", time.Since(start).Milliseconds()),
			})
			return report, err
		}

		if report != nil && report.HasDrift() {
			_ = logger.Log(EventDriftDetected, host, "drift detected", map[string]string{
				"drifted_count": fmt.Sprintf("%d", report.DriftedCount()),
			})
		}

		_ = logger.Log(EventCheckFinished, host, "drift check finished", map[string]string{
			"duration_ms": fmt.Sprintf("%d", time.Since(start).Milliseconds()),
			"has_drift":   fmt.Sprintf("%v", report != nil && report.HasDrift()),
		})
		return report, nil
	}
}
