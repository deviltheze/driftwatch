// Package baseline provides persistent storage for configuration file hashes
// used as reference points for drift detection.
//
// A Store holds host+path -> hash mappings on disk (JSON). On each drift check
// cycle the engine compares the freshly computed hash against the stored
// baseline; a mismatch indicates configuration drift.
//
// Typical usage:
//
//	store, err := baseline.NewStore("/var/lib/driftwatch/baseline.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Record initial baseline
//	_ = store.Set("web01", "/etc/nginx/nginx.conf", computedHash)
//
//	// Later, compare
//	entry, ok := store.Get("web01", "/etc/nginx/nginx.conf")
//	if ok && entry.Hash != newHash {
//		// drift detected
//	}
package baseline
