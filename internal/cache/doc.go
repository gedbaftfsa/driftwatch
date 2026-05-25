// Package cache provides a lightweight, thread-safe, TTL-based in-memory
// cache designed to reduce redundant calls to the runtime provider during
// high-frequency drift-detection polling.
//
// Typical usage:
//
//	store := cache.New(30 * time.Second)
//
//	// Attempt to serve from cache.
//	if summaries, ok := store.Get("default"); ok {
//		return summaries, nil
//	}
//
//	// Cache miss — fetch from provider and populate.
//	summaries, err := provider.FetchAll(ctx)
//	if err != nil {
//		return nil, err
//	}
//	store.Set("default", summaries)
//	return summaries, nil
//
// Background eviction of stale entries can be started with StartEvictLoop.
package cache
