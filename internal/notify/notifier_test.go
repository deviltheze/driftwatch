package notify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/notify"
)

func makeReport(results []drift.Result) drift.Report {
	return drift.Report{Results: results}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := notify.New(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNotify_NoDrift_InfoLevel(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	report := makeReport([]drift.Result{
		{File: "/etc/hosts", Drifted: false},
	})

	if err := n.Notify("server1", report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, string(notify.LevelInfo)) {
		t.Errorf("expected INFO level in output, got: %s", out)
	}
	if !strings.Contains(out, "host=server1") {
		t.Errorf("expected host in output, got: %s", out)
	}
}

func TestNotify_WithDrift_AlertLevel(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	report := makeReport([]drift.Result{
		{File: "/etc/ssh/sshd_config", Drifted: true},
		{File: "/etc/hosts", Drifted: false},
	})

	if err := n.Notify("server2", report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, string(notify.LevelAlert)) {
		t.Errorf("expected ALERT level in output, got: %s", out)
	}
	if !strings.Contains(out, "drifted=1") {
		t.Errorf("expected drifted=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "total=2") {
		t.Errorf("expected total=2 in output, got: %s", out)
	}
}

func TestNotify_EmptyReport(t *testing.T) {
	var buf bytes.Buffer
	n := notify.New(&buf)

	report := makeReport([]drift.Result{})
	if err := n.Notify("server3", report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "drifted=0") {
		t.Errorf("expected drifted=0 in output, got: %s", out)
	}
}
