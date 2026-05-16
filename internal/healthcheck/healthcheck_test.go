package healthcheck_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/healthcheck"
)

// stubProber returns a preset Result regardless of host.
type stubProber struct {
	result healthcheck.Result
}

func (s *stubProber) Probe(_ context.Context, host string) healthcheck.Result {
	r := s.result
	r.Host = host
	return r
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    healthcheck.Status
		want string
	}{
		{healthcheck.StatusHealthy, "healthy"},
		{healthcheck.StatusUnhealthy, "unhealthy"},
		{healthcheck.StatusUnknown, "unknown"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestNew_DefaultTimeout(t *testing.T) {
	r := healthcheck.New(&stubProber{}, 0, nil)
	if r == nil {
		t.Fatal("expected non-nil Runner")
	}
}

func TestCheck_AllHealthy(t *testing.T) {
	prober := &stubProber{result: healthcheck.Result{Status: healthcheck.StatusHealthy}}
	r := healthcheck.New(prober, 2*time.Second, nil)
	hosts := []string{"host1", "host2", "host3"}
	results := r.Check(context.Background(), hosts)
	if len(results) != len(hosts) {
		t.Fatalf("expected %d results, got %d", len(hosts), len(results))
	}
	for _, res := range results {
		if res.Status != healthcheck.StatusHealthy {
			t.Errorf("host %s: expected healthy, got %s", res.Host, res.Status)
		}
	}
}

func TestCheck_SomeUnhealthy(t *testing.T) {
	call := 0
	prober := &stubProber{result: healthcheck.Result{Status: healthcheck.StatusUnhealthy}}
	_ = call
	r := healthcheck.New(prober, 2*time.Second, nil)
	results := r.Check(context.Background(), []string{"bad"})
	if healthcheck.Healthy(results) {
		t.Error("expected Healthy to return false")
	}
}

func TestHealthy_EmptySlice(t *testing.T) {
	if !healthcheck.Healthy(nil) {
		t.Error("empty slice should be considered healthy")
	}
}

func TestSummary_Format(t *testing.T) {
	results := []healthcheck.Result{
		{Status: healthcheck.StatusHealthy},
		{Status: healthcheck.StatusHealthy},
		{Status: healthcheck.StatusUnhealthy},
	}
	s := healthcheck.Summary(results)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestTCPProber_UnreachableHost(t *testing.T) {
	// Use a port that is almost certainly not listening.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skip("cannot bind listener")
	}
	addr := ln.Addr().(*net.TCPAddr)
	_ = ln.Close() // close immediately so the port is gone

	prober := healthcheck.NewTCPProber(fmt.Sprint(addr.Port))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res := prober.Probe(ctx, "127.0.0.1")
	if res.Status != healthcheck.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}
