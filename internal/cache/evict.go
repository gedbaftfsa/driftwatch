package cache

import "time"

// Evict removes all expired entries from the store.
// It is safe to call concurrently and is intended to be invoked
// periodically (e.g. from a background goroutine) to reclaim memory.
func (s *Store) Evict() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	removed := 0
	now := time.Now()
	for k, e := range s.entries {
		if now.After(e.ExpiresAt) {
			delete(s.entries, k)
			removed++
		}
	}
	return removed
}

// StartEvictLoop launches a background goroutine that calls Evict every
// interval. The goroutine exits when done is closed.
func (s *Store) StartEvictLoop(interval time.Duration, done <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.Evict()
			case <-done:
				return
			}
		}
	}()
}
