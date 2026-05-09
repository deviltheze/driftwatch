// Package retry implements exponential back-off retry logic used by the
// driftwatch SSH subsystem when transient connection or command errors occur.
//
// Usage:
//
//	p := retry.DefaultPolicy()
//	err := retry.Do(p, logger, func() error {
//		return sshClient.Connect(host)
//	})
//
// The policy controls the maximum number of attempts, the initial delay,
// the exponential multiplier, and the maximum delay ceiling.  A nil logger
// suppresses per-attempt log output.
package retry
