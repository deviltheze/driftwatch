// Package watchdog provides a self-monitoring component for the driftwatch
// daemon. It periodically pings registered components via the Pinger interface
// and invokes a configurable failure callback when a component stops
// responding.
//
// Usage:
//
//	 wd := watchdog.New(30*time.Second, logger)
//	 wd.Register("engine", myEngine)
//	 wd.OnFailure(func(name string, err error) {
//	     // alert or restart logic
//	 })
//	 go wd.Start(ctx)
package watchdog
