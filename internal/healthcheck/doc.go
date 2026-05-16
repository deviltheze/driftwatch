// Package healthcheck implements lightweight liveness probing for remote
// hosts managed by driftwatch.
//
// A Prober implementation (e.g. TCPProber) is injected into a Runner which
// fans out concurrent probes with per-host deadlines derived from a shared
// timeout.  Results carry latency, status and any error so that upstream
// components (alert, metrics) can act accordingly.
//
// Typical usage:
//
//	prober := healthcheck.NewTCPProber("22")
//	runner := healthcheck.New(prober, 5*time.Second, nil)
//	results := runner.Check(ctx, hosts)
//	if !healthcheck.Healthy(results) {
//		log.Println(healthcheck.Summary(results))
//	}
package healthcheck
