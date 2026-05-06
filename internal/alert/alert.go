// Package alert provides threshold-based alerting for drift detection results.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single alert event.
type Alert struct {
	Level     Level
	Host      string
	Message   string
	Timestamp time.Time
}

// String returns a formatted string representation of the alert.
func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s %s: %s", a.Level, a.Timestamp.Format(time.RFC3339), a.Host, a.Message)
}

// Handler processes alerts generated from drift reports.
type Handler struct {
	threshold int
	out       io.Writer
}

// NewHandler creates a new alert Handler. threshold is the minimum number
// of drifted checks required to trigger an ALERT level (vs WARN).
func NewHandler(threshold int, out io.Writer) *Handler {
	if out == nil {
		out = os.Stdout
	}
	if threshold <= 0 {
		threshold = 1
	}
	return &Handler{threshold: threshold, out: out}
}

// Evaluate inspects a drift report and returns an Alert if drift is detected.
// Returns nil if no drift is present.
func (h *Handler) Evaluate(report drift.Report) *Alert {
	if !report.HasDrift() {
		return nil
	}

	level := LevelWarn
	if report.DriftedCount() >= h.threshold {
		level = LevelAlert
	}

	return &Alert{
		Level:     level,
		Host:      report.Host,
		Message:   fmt.Sprintf("%d check(s) drifted", report.DriftedCount()),
		Timestamp: time.Now().UTC(),
	}
}

// Emit writes the alert to the configured output writer.
func (h *Handler) Emit(a *Alert) error {
	if a == nil {
		return nil
	}
	_, err := fmt.Fprintln(h.out, a.String())
	return err
}
