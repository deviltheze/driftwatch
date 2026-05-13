package rollup_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/rollup"
)

func makeReport(host string, drifted bool) *drift.Report {
	results := []drift.Result{}
	if drifted {
		results = append(results, drift.Result{
			File:     "/etc/hosts",
			Expected: "abc123",
			Actual:   "def456",
		})
	}
	return &drift.Report{Host: host, Results: results}
}

func TestAggregate_NoReports(t *testing.T) {
	s := rollup.Aggregate(nil)
	if s.TotalHosts != 0 {
		t.Errorf("expected 0 total hosts, got %d", s.TotalHosts)
	}
	if s.HasDrift() {
		t.Error("expected no drift for empty report set")
	}
}

func TestAggregate_AllClean(t *testing.T) {
	reports := []*drift.Report{
		makeReport("host-a", false),
		makeReport("host-b", false),
	}
	s := rollup.Aggregate(reports)
	if s.TotalHosts != 2 {
		t.Errorf("expected 2 total hosts, got %d", s.TotalHosts)
	}
	if s.HasDrift() {
		t.Error("expected no drift")
	}
	if len(s.CleanHosts) != 2 {
		t.Errorf("expected 2 clean hosts, got %d", len(s.CleanHosts))
	}
}

func TestAggregate_SomeDrifted(t *testing.T) {
	reports := []*drift.Report{
		makeReport("host-a", true),
		makeReport("host-b", false),
		makeReport("host-c", true),
	}
	s := rollup.Aggregate(reports)
	if !s.HasDrift() {
		t.Error("expected drift to be detected")
	}
	if len(s.DriftedHosts) != 2 {
		t.Errorf("expected 2 drifted hosts, got %d", len(s.DriftedHosts))
	}
	if len(s.CleanHosts) != 1 {
		t.Errorf("expected 1 clean host, got %d", len(s.CleanHosts))
	}
}

func TestDriftRatio(t *testing.T) {
	reports := []*drift.Report{
		makeReport("host-a", true),
		makeReport("host-b", false),
	}
	s := rollup.Aggregate(reports)
	if got := s.DriftRatio(); got != 0.5 {
		t.Errorf("expected drift ratio 0.5, got %f", got)
	}
}

func TestSummary_String_NoDrift(t *testing.T) {
	s := rollup.Aggregate([]*drift.Report{makeReport("host-a", false)})
	str := s.String()
	if str == "" {
		t.Error("expected non-empty summary string")
	}
}

func TestSummary_String_WithDrift(t *testing.T) {
	s := rollup.Aggregate([]*drift.Report{makeReport("host-a", true)})
	str := s.String()
	if str == "" {
		t.Error("expected non-empty summary string")
	}
}
