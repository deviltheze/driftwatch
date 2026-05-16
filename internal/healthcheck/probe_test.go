package healthcheck_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/healthcheck"
)

func TestNewTCPProber_DefaultPort(t *testing.T) {
	p := healthcheck.NewTCPProber("")
	if p.Port != "22" {
		t.Errorf("expected default port 22, got %s", p.Port)
	}
}

func TestNewTCPProber_CustomPort(t *testing.T) {
	p := healthcheck.NewTCPProber("2222")
	if p.Port != "2222" {
		t.Errorf("expected port 2222, got %s", p.Port)
	}
}

func TestTCPProber_HealthyHost(t *testing.T) {
	// Start a real TCP listener so the probe succeeds.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	port := fmt.Sprint(ln.Addr().(*net.TCPAddr).Port)
	prober := healthcheck.NewTCPProber(port)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res := prober.Probe(ctx, "127.0.0.1")
	if res.Status != healthcheck.StatusHealthy {
		t.Errorf("expected healthy, got %s (err: %v)", res.Status, res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
	if res.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", res.Host)
	}
}

func TestTCPProber_SetsCheckedAt(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	port := fmt.Sprint(ln.Addr().(*net.TCPAddr).Port)
	prober := healthcheck.NewTCPProber(port)

	before := time.Now()
	res := prober.Probe(context.Background(), "127.0.0.1")
	after := time.Now()

	if res.CheckedAt.Before(before) || res.CheckedAt.After(after) {
		t.Errorf("CheckedAt %v not in expected range [%v, %v]", res.CheckedAt, before, after)
	}
}
