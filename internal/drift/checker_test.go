package drift

import (
	"testing"
)

func TestHashContent_Deterministic(t *testing.T) {
	content := "[server]\nport=8080\n"
	first := HashContent(content)
	second := HashContent(content)
	if first != second {
		t.Errorf("expected deterministic hash, got %q and %q", first, second)
	}
}

func TestHashContent_DifferentInputs(t *testing.T) {
	a := HashContent("config_a=1")
	b := HashContent("config_b=2")
	if a == b {
		t.Error("expected different hashes for different inputs")
	}
}

func TestHashContent_KnownValue(t *testing.T) {
	// echo -n "hello" | sha256sum
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	got := HashContent("hello")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestResult_Drifted(t *testing.T) {
	r := Result{
		Host:     "web-01",
		Path:     "/etc/app/config.toml",
		Drifted:  true,
		Baseline: "abc123",
		Actual:   "def456",
	}
	if !r.Drifted {
		t.Error("expected result to be drifted")
	}
	if r.Baseline == r.Actual {
		t.Error("expected baseline and actual to differ")
	}
}

func TestResult_NoDrift(t *testing.T) {
	checksum := HashContent("stable config content")
	r := Result{
		Host:     "db-01",
		Path:     "/etc/mysql/my.cnf",
		Drifted:  false,
		Baseline: checksum,
		Actual:   checksum,
	}
	if r.Drifted {
		t.Error("expected result to not be drifted")
	}
	if r.Baseline != r.Actual {
		t.Errorf("expected baseline %q to equal actual %q", r.Baseline, r.Actual)
	}
}

func TestResult_WithError(t *testing.T) {
	r := Result{
		Host: "unreachable-host",
		Path: "/etc/config",
		Err:  fmt.Errorf("connection refused"),
	}
	if r.Err == nil {
		t.Error("expected error to be set")
	}
}
