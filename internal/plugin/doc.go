// Package plugin implements a named-plugin registry for driftwatch.
//
// Plugins extend the drift-check pipeline by allowing callers to register
// custom CheckFn functions that annotate or transform raw SSH command output
// before it is compared against the stored baseline.
//
// Usage:
//
//	reg := plugin.NewRegistry()
//	err := reg.Register(plugin.Plugin{
//		Name:  "uptime-check",
//		Check: func(host, out string) (string, error) { ... },
//	})
package plugin
