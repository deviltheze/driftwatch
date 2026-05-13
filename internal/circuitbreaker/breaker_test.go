package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/circuitbreaker"
)

func TestNew_DefaultValues(t *testing.T) {
	b := circuitbreaker.New(0, 0)
	if b == nil {
		t.Fatal("expected non-nil breaker")
	}
	// Breaker should start closed for any host.
	if !b.Allow("host1") {
		t.Error("expected Allow to return true for fresh host")
	}
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	if !b.Allow("server-a") {
		t.Error("new host should be allowed")
	}
	if b.StateOf("server-a") != circuitbreaker.StateClosed {
		t.Errorf("expected closed, got %s", b.StateOf("server-a"))
	}
}

func TestRecordFailure_OpensCircuit(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	host := "server-b"

	b.RecordFailure(host)
	b.RecordFailure(host)
	if b.StateOf(host) != circuitbreaker.StateClosed {
		t.Error("should still be closed below threshold")
	}

	b.RecordFailure(host) // hits threshold
	if b.StateOf(host) != circuitbreaker.StateOpen {
		t.Errorf("expected open after threshold, got %s", b.StateOf(host))
	}
	if b.Allow(host) {
		t.Error("open circuit should block requests")
	}
}

func TestRecordSuccess_ClosesCircuit(t *testing.T) {
	b := circuitbreaker.New(2, time.Minute)
	host := "server-c"

	b.RecordFailure(host)
	b.RecordFailure(host)
	if b.StateOf(host) != circuitbreaker.StateOpen {
		t.Fatal("expected open")
	}

	b.RecordSuccess(host)
	if b.StateOf(host) != circuitbreaker.StateClosed {
		t.Errorf("expected closed after success, got %s", b.StateOf(host))
	}
	if !b.Allow(host) {
		t.Error("closed circuit should allow requests")
	}
}

func TestHalfOpen_AfterTimeout(t *testing.T) {
	b := circuitbreaker.New(1, 50*time.Millisecond)
	host := "server-d"

	b.RecordFailure(host)
	if b.StateOf(host) != circuitbreaker.StateOpen {
		t.Fatal("expected open")
	}

	time.Sleep(60 * time.Millisecond)

	if !b.Allow(host) {
		t.Error("expected Allow after timeout (half-open)")
	}
	if b.StateOf(host) != circuitbreaker.StateHalfOpen {
		t.Errorf("expected half-open, got %s", b.StateOf(host))
	}
}

func TestReset_ClearsState(t *testing.T) {
	b := circuitbreaker.New(1, time.Minute)
	host := "server-e"

	b.RecordFailure(host)
	if b.StateOf(host) != circuitbreaker.StateOpen {
		t.Fatal("expected open")
	}

	b.Reset(host)
	if b.StateOf(host) != circuitbreaker.StateClosed {
		t.Errorf("expected closed after reset, got %s", b.StateOf(host))
	}
}

func TestErrOpen_Message(t *testing.T) {
	err := circuitbreaker.ErrOpen{Host: "my-host"}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestState_String(t *testing.T) {
	cases := []struct {
		state circuitbreaker.State
		want  string
	}{
		{circuitbreaker.StateClosed, "closed"},
		{circuitbreaker.StateOpen, "open"},
		{circuitbreaker.StateHalfOpen, "half-open"},
	}
	for _, tc := range cases {
		if got := tc.state.String(); got != tc.want {
			t.Errorf("State(%d).String() = %q, want %q", tc.state, got, tc.want)
		}
	}
}
