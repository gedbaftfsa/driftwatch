// Package cache provides a short-lived in-memory cache for drift summaries,
// reducing redundant provider calls during high-frequency polling.
package cache

import (
	"sync"
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Entry holds a cached value with an expiry timestamp.
type Entry struct {
	Summaries []drift.Summary
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the entry is past its TTL.
func (e Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Store is a thread-safe in-memory cache keyed by an arbitrary string.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New returns a Store with the given TTL applied to every cached entry.
func New(ttl time.Duration) *Store {
	return &Store{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Set stores summaries under key, overwriting any existing entry.
func (s *Store) Set(key string, summaries []drift.Summary) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.entries[key] = Entry{
		Summaries: summaries,
		CachedAt:  now,
		ExpiresAt: now.Add(s.ttl),
	}
}

// Get returns the cached summaries and true if the entry exists and has not
// expired. Otherwise it returns nil and false.
func (s *Store) Get(key string) ([]drift.Summary, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key]
	if !ok || e.IsExpired() {
		return nil, false
	}
	return e.Summaries, true
}

// Invalidate removes the entry for key, if present.
func (s *Store) Invalidate(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Flush removes all entries from the store.
func (s *Store) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[string]Entry)
}

// Len returns the number of entries currently held (including expired ones
// that have not yet been evicted).
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
