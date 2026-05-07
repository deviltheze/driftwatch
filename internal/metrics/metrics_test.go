package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/metrics"
)

func TestNewCollector_ZeroValues(t *testing.T) {
	c := metrics.NewCollector()
	s := c.Snapshot()
	if s.ChecksTotal != 0 || s.DriftDetected != 0 || s.HostsChecked != 0 {
		t.Fatalf("expected zero snapshot, got %+v", s)
	}
	if !s.LastCheckAt.IsZero() {
		t.Fatalf("expected zero time, got %v", s.LastCheckAt)
	}
}

func TestRecordCheck_IncrementsCounters(t *testing.T) {
	c := metrics.NewCollector()
	c.RecordCheck(3, 1)
	s := c.Snapshot()
	if s.ChecksTotal != 1 {
		t.Errorf("ChecksTotal: want 1, got %d", s.ChecksTotal)
	}
	if s.HostsChecked != 3 {
		t.Errorf("HostsChecked: want 3, got %d", s.HostsChecked)
	}
	if s.DriftDetected != 1 {
		t.Errorf("DriftDetected: want 1, got %d", s.DriftDetected)
	}
}

func TestRecordCheck_Accumulates(t *testing.T) {
	c := metrics.NewCollector()
	c.RecordCheck(2, 0)
	c.RecordCheck(4, 2)
	s := c.Snapshot()
	if s.ChecksTotal != 2 {
		t.Errorf("ChecksTotal: want 2, got %d", s.ChecksTotal)
	}
	if s.HostsChecked != 6 {
		t.Errorf("HostsChecked: want 6, got %d", s.HostsChecked)
	}
	if s.DriftDetected != 2 {
		t.Errorf("DriftDetected: want 2, got %d", s.DriftDetected)
	}
}

func TestRecordCheck_UpdatesTimestamp(t *testing.T) {
	c := metrics.NewCollector()
	before := time.Now()
	c.RecordCheck(1, 0)
	after := time.Now()
	s := c.Snapshot()
	if s.LastCheckAt.Before(before) || s.LastCheckAt.After(after) {
		t.Errorf("LastCheckAt %v not in expected range [%v, %v]", s.LastCheckAt, before, after)
	}
}

func TestWrite_ContainsExpectedFields(t *testing.T) {
	c := metrics.NewCollector()
	c.RecordCheck(5, 2)
	var buf bytes.Buffer
	c.Write(&buf)
	out := buf.String()
	for _, want := range []string{"checks_total=1", "hosts_checked=5", "drift_detected=2", "last_check="} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got: %s", want, out)
		}
	}
}

func TestWrite_NeverChecked_ShowsNever(t *testing.T) {
	c := metrics.NewCollector()
	var buf bytes.Buffer
	c.Write(&buf)
	if !strings.Contains(buf.String(), "last_check=never") {
		t.Errorf("expected 'last_check=never', got: %s", buf.String())
	}
}
