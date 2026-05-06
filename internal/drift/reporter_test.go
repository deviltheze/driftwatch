package drift

import (
	"strings"
	"testing"
	"time"
)

func TestReport_HasDrift_True(t *testing.T) {
	report := Report{
		Results: []Result{
			{Path: "/etc/hosts", Drifted: false},
			{Path: "/etc/passwd", Drifted: true},
		},
	}
	if !report.HasDrift() {
		t.Error("expected HasDrift to return true")
	}
}

func TestReport_HasDrift_False(t *testing.T) {
	report := Report{
		Results: []Result{
			{Path: "/etc/hosts", Drifted: false},
		},
	}
	if report.HasDrift() {
		t.Error("expected HasDrift to return false")
	}
}

func TestReport_DriftedCount(t *testing.T) {
	report := Report{
		Results: []Result{
			{Path: "/etc/hosts", Drifted: false},
			{Path: "/etc/passwd", Drifted: true},
			{Path: "/etc/ssh/sshd_config", Drifted: true},
		},
	}
	if count := report.DriftedCount(); count != 2 {
		t.Errorf("expected DriftedCount 2, got %d", count)
	}
}

func TestReporter_Write_NoDrift(t *testing.T) {
	var buf strings.Builder
	r := NewReporter(&buf)

	report := Report{
		Host:      "web-01",
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Results:   []Result{{Path: "/etc/hosts", Drifted: false}},
	}

	if err := r.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "web-01") {
		t.Error("expected output to contain host name")
	}
	if !strings.Contains(out, "OK") {
		t.Error("expected output to contain OK status")
	}
	if strings.Contains(out, "DRIFT DETECTED") {
		t.Error("did not expect DRIFT DETECTED in output")
	}
}

func TestReporter_Write_WithDrift(t *testing.T) {
	var buf strings.Builder
	r := NewReporter(&buf)

	report := Report{
		Host:      "db-01",
		Timestamp: time.Now(),
		Results: []Result{
			{Path: "/etc/my.cnf", Drifted: true, ExpectedHash: "abc123", ActualHash: "def456"},
		},
	}

	if err := r.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DRIFT DETECTED") {
		t.Error("expected DRIFT DETECTED in output")
	}
	if !strings.Contains(out, "/etc/my.cnf") {
		t.Error("expected path in output")
	}
	if !strings.Contains(out, "abc123") {
		t.Error("expected expected hash in output")
	}
	if !strings.Contains(out, "def456") {
		t.Error("expected actual hash in output")
	}
}
