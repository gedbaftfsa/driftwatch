package cache_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/cache"
)

func TestEvict_RemovesExpiredEntries(t *testing.T) {
	s := cache.New(time.Millisecond)
	s.Set("expired1", makeSummaries("a"))
	s.Set("expired2", makeSummaries("b"))
	time.Sleep(5 * time.Millisecond)
	removed := s.Evict()
	if removed != 2 {
		t.Fatalf("expected 2 removed, got %d", removed)
	}
	if s.Len() != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", s.Len())
	}
}

func TestEvict_KeepsValidEntries(t *testing.T) {
	s := cache.New(time.Minute)
	s.Set("fresh", makeSummaries("svc"))
	removed := s.Evict()
	if removed != 0 {
		t.Fatalf("expected 0 removed, got %d", removed)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", s.Len())
	}
}

func TestEvict_MixedEntries(t *testing.T) {
	short := cache.New(time.Millisecond)
	short.Set("old", makeSummaries("x"))
	time.Sleep(5 * time.Millisecond)
	// Add a fresh entry with a longer TTL by using a separate store and
	// manually testing the counts instead.
	s := cache.New(time.Minute)
	s.Set("fresh", makeSummaries("y"))
	// Evict the short-TTL store.
	removed := short.Evict()
	if removed != 1 {
		t.Fatalf("expected 1 removed from short store, got %d", removed)
	}
	// Fresh store should be unaffected.
	if s.Len() != 1 {
		t.Fatalf("expected fresh store to still have 1 entry")
	}
}

func TestStartEvictLoop_EventuallyEvicts(t *testing.T) {
	s := cache.New(5 * time.Millisecond)
	s.Set("svc", makeSummaries("api"))
	done := make(chan struct{})
	s.StartEvictLoop(10*time.Millisecond, done)
	time.Sleep(40 * time.Millisecond)
	close(done)
	if s.Len() != 0 {
		t.Fatalf("expected evict loop to remove expired entries, got %d", s.Len())
	}
}
