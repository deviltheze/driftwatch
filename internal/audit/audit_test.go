package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/audit"
)

func TestNew_DefaultsToStdout(t *testing.T) {
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLog_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	err := l.Log(audit.EventCheckStarted, "host-01", "check started", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if e.Type != audit.EventCheckStarted {
		t.Errorf("expected type %q, got %q", audit.EventCheckStarted, e.Type)
	}
	if e.Host != "host-01" {
		t.Errorf("expected host %q, got %q", "host-01", e.Host)
	}
	if e.Message != "check started" {
		t.Errorf("expected message %q, got %q", "check started", e.Message)
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLog_WithMeta(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	meta := map[string]string{"file": "/etc/hosts", "hash": "abc123"}
	if err := l.Log(audit.EventDriftDetected, "host-02", "drift found", meta); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if e.Meta["file"] != "/etc/hosts" {
		t.Errorf("expected meta file /etc/hosts, got %q", e.Meta["file"])
	}
}

func TestLog_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	events := []audit.EventType{
		audit.EventCheckStarted,
		audit.EventDriftDetected,
		audit.EventAlertSent,
		audit.EventCheckFinished,
	}
	for _, et := range events {
		if err := l.Log(et, "h", "msg", nil); err != nil {
			t.Fatalf("log error: %v", err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != len(events) {
		t.Errorf("expected %d lines, got %d", len(events), len(lines))
	}
}

func TestLog_NoHost(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)

	if err := l.Log(audit.EventCheckStarted, "", "global check", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if e.Host != "" {
		t.Errorf("expected empty host, got %q", e.Host)
	}
}
