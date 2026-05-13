package redact

import (
	"regexp"
	"strings"
	"testing"
)

func TestDefaultRules_NotNil(t *testing.T) {
	r := DefaultRules()
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestApply_Password(t *testing.T) {
	r := DefaultRules()
	input := "password=supersecret"
	out := r.Apply(input)
	if strings.Contains(out, "supersecret") {
		t.Errorf("sensitive value not redacted: %s", out)
	}
	if !strings.Contains(out, placeholder) {
		t.Errorf("expected placeholder in output, got: %s", out)
	}
}

func TestApply_Token(t *testing.T) {
	r := DefaultRules()
	out := r.Apply("token=abc123xyz")
	if strings.Contains(out, "abc123xyz") {
		t.Errorf("token not redacted: %s", out)
	}
}

func TestApply_Secret(t *testing.T) {
	r := DefaultRules()
	out := r.Apply("SECRET=myvalue")
	if strings.Contains(out, "myvalue") {
		t.Errorf("secret not redacted: %s", out)
	}
}

func TestApply_APIKey(t *testing.T) {
	r := DefaultRules()
	out := r.Apply("api_key=key-9999")
	if strings.Contains(out, "key-9999") {
		t.Errorf("api_key not redacted: %s", out)
	}
}

func TestApply_NoSensitiveData(t *testing.T) {
	r := DefaultRules()
	input := "hostname=web01.example.com"
	out := r.Apply(input)
	if out != input {
		t.Errorf("non-sensitive input was modified: %s", out)
	}
}

func TestLines_MultiLine(t *testing.T) {
	r := DefaultRules()
	input := "host=web01\npassword=secret123\nport=22"
	out := r.Lines(input)
	if strings.Contains(out, "secret123") {
		t.Errorf("sensitive value survived multi-line redaction: %s", out)
	}
	if !strings.Contains(out, "host=web01") {
		t.Errorf("non-sensitive line was altered: %s", out)
	}
}

func TestNew_CustomRule(t *testing.T) {
	rules := []Rule{
		{
			Label:   "custom",
			Pattern: regexp.MustCompile(`(?i)(private\s*=\s*)\S+`),
		},
	}
	r := New(rules)
	out := r.Apply("private=topsecret")
	if strings.Contains(out, "topsecret") {
		t.Errorf("custom rule did not redact value: %s", out)
	}
}

func TestApply_PreservesKeyName(t *testing.T) {
	r := DefaultRules()
	out := r.Apply("password=hunter2")
	if !strings.HasPrefix(out, "password=") {
		t.Errorf("key name not preserved after redaction: %s", out)
	}
}
