package ssh

import (
	"errors"
	"strings"
	"testing"
)

func TestCommandResult_String_Success(t *testing.T) {
	cr := CommandResult{
		Host:    "10.0.0.1",
		Command: "uname -r",
		Output:  "5.15.0-generic",
		Err:     nil,
	}

	got := cr.String()
	if !strings.Contains(got, "10.0.0.1") {
		t.Errorf("expected host in output, got: %s", got)
	}
	if !strings.Contains(got, "uname -r") {
		t.Errorf("expected command in output, got: %s", got)
	}
	if !strings.Contains(got, "5.15.0-generic") {
		t.Errorf("expected output in result, got: %s", got)
	}
}

func TestCommandResult_String_Error(t *testing.T) {
	cr := CommandResult{
		Host:    "10.0.0.2",
		Command: "cat /etc/secret",
		Output:  "",
		Err:     errors.New("permission denied"),
	}

	got := cr.String()
	if !strings.Contains(got, "ERROR") {
		t.Errorf("expected ERROR label in output, got: %s", got)
	}
	if !strings.Contains(got, "permission denied") {
		t.Errorf("expected error message in output, got: %s", got)
	}
}

func TestNewRunner(t *testing.T) {
	r := NewRunner(nil, "host-a")
	if r == nil {
		t.Fatal("expected non-nil Runner")
	}
	if r.host != "host-a" {
		t.Errorf("expected host 'host-a', got %q", r.host)
	}
}

func TestRunner_Run_EmptyCommands(t *testing.T) {
	r := NewRunner(nil, "host-b")
	results := r.Run([]string{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty command list, got %d", len(results))
	}
}

func TestRunner_Run_ResultCountMatchesCommands(t *testing.T) {
	// Verify that Run returns exactly one result per command, even when
	// execution fails due to a nil client.
	r := NewRunner(nil, "host-c")
	cmds := []string{"uname -r", "uptime", "hostname"}
	results := r.Run(cmds)
	if len(results) != len(cmds) {
		t.Errorf("expected %d results, got %d", len(cmds), len(results))
	}
	for _, res := range results {
		if res.Host != "host-c" {
			t.Errorf("expected host 'host-c', got %q", res.Host)
		}
	}
}
