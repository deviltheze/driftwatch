// Package snapshot provides primitives for capturing, storing, and comparing
// remote file state hashes collected over SSH.
//
// Typical usage:
//
//	collector := snapshot.NewCollector(runner, logger)
//	store     := snapshot.NewStore()
//
//	curr, err := collector.Collect(host, "/etc/hosts")
//	if prev, ok := store.Get(host, "/etc/hosts"); ok {
//		if diff, changed := snapshot.Compare(prev, curr); changed {
//			log.Println(diff)
//		}
//	}
//	store.Set(curr)
package snapshot
