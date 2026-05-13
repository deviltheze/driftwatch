// Package filter implements host-level tag filtering for driftwatch.
//
// Use filter.New to build a HostFilter from include/exclude tag lists,
// then call Match or Filter to determine which configured hosts should
// participate in a drift-check run.
//
// Example:
//
//	f := filter.New(
//		[]string{"production"},  // only production-tagged hosts
//		[]string{"maintenance"}, // skip hosts under maintenance
//	)
//
//	hosts := map[string][]string{
//		"web-01": {"production", "web"},
//		"db-01":  {"production", "maintenance"},
//		"dev-01": {"development"},
//	}
//
//	selected := f.Filter(hosts) // ["web-01"]
package filter
