package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/driftwatch/internal/alert"
	"github.com/driftwatch/internal/drift"
)

func makeReport(host string, results []drift.Result) drift.Report {
	return drift.Report{
		Host:    host,
		Results: results,
	}
}

func TestAlert_String_ContainsLevel(t *testing.T) {
	a := alert.Alert{
		Level:     alert.LevelAlert,
		Host:      "web-01",
		Message:   "2 check(s) drifted",
		Timestamp: time.Now().UTC(),
	}
	got := a.String()
	if !strings.Contains(got, "ALERT") {
		t.Errorf("expected ALERT in string, got: %s", got)
	}
	if !strings.Contains(got, "web-01") {
		t.Errorf("expected host in string, got: %s", got)
	}
}

func TestNewHandler_DefaultThreshold(t *testing.T) {
	h := alert.NewHandler(0, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestEvaluate_NoDrift_ReturnsNil(t *testing.T) {
	h := alert.NewHandler(1, nil)
	report := makeReport("host-1", []drift.Result{
		{Check: "uptime", Expected: "ok", Actual: "ok"},
	})
	got := h.Evaluate(report)
	if got != nil {
		t.Errorf("expected nil alert for no-drift report, got: %v", got)
	}
}

func TestEvaluate_WithDrift_BelowThreshold_ReturnsWarn(t *testing.T) {
	h := alert.NewHandler(3, nil)
	report := makeReport("host-2", []drift.Result{
		{Check: "sshd", Expected: "running", Actual: "stopped"},
	})
	got := h.Evaluate(report)
	if got == nil {
		t.Fatal("expected alert, got nil")
	}
	if got.Level != alert.LevelWarn {
		t.Errorf("expected WARN, got %s", got.Level)
	}
}

func TestEvaluate_WithDrift_AtThreshold_ReturnsAlert(t *testing.T) {
	h := alert.NewHandler(1, nil)
	report := makeReport("host-3", []drift.Result{
		{Check: "sshd", Expected: "running", Actual: "stopped"},
	})
	got := h.Evaluate(report)
	if got == nil {
		t.Fatal("expected alert, got nil")
	}
	if got.Level != alert.LevelAlert {
		t.Errorf("expected ALERT, got %s", got.Level)
	}
}

func TestEmit_WritesToOutput(t *testing.T) {
	var buf bytes.Buffer
	h := alert.NewHandler(1, &buf)
	a := &alert.Alert{
		Level:     alert.LevelAlert,
		Host:      "db-01",
		Message:   "1 check(s) drifted",
		Timestamp: time.Now().UTC(),
	}
	if err := h.Emit(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "db-01") {
		t.Errorf("expected host in output, got: %s", buf.String())
	}
}

func TestEmit_NilAlert_NoError(t *testing.T) {
	var buf bytes.Buffer
	h := alert.NewHandler(1, &buf)
	if err := h.Emit(nil); err != nil {
		t.Fatalf("expected no error for nil alert, got: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for nil alert")
	}
}
