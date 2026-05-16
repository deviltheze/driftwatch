package healthcheck

import (
	"context"
	"net"
	"time"
)

// TCPProber checks host reachability by opening a TCP connection.
type TCPProber struct {
	Port string
}

// NewTCPProber returns a TCPProber that dials the given port (default "22").
func NewTCPProber(port string) *TCPProber {
	if port == "" {
		port = "22"
	}
	return &TCPProber{Port: port}
}

// Probe attempts a TCP dial to host:port and returns a Result.
func (p *TCPProber) Probe(ctx context.Context, host string) Result {
	addr := net.JoinHostPort(host, p.Port)
	start := time.Now()
	res := Result{
		Host:      host,
		CheckedAt: start,
	}

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	res.Latency = time.Since(start)

	if err != nil {
		res.Status = StatusUnhealthy
		res.Err = err
		return res
	}
	_ = conn.Close()
	res.Status = StatusHealthy
	return res
}
