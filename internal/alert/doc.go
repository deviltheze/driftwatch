// Package alert implements threshold-based alerting for driftwatch.
//
// It evaluates drift.Report values and produces Alert events when
// configuration drift is detected on remote hosts. Alerts carry a
// severity Level (WARN or ALERT) determined by comparing the number
// of drifted checks against a configurable threshold.
//
// Suppression prevents duplicate alerts from being emitted repeatedly
// within a configurable cooldown window, reducing noise during
// sustained drift conditions.
//
// Typical usage:
//
//	handler := alert.NewHandler(threshold, os.Stdout)
//	suppression := alert.NewSuppression(15 * time.Minute)
//
//	if a := handler.Evaluate(report); suppression.Allow(a) {
//	    handler.Emit(a)
//	}
package alert
