package cache_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/cache"
	"github.com/example/driftwatch/internal/drift"
)

func makeSummaries(names ...string) []drift.Summary {
	out := make([]drift.Summary, 0, len(names))
	for _, n := range names {
		out = append(out, drift.Summary{ServiceName: n})
	}
	return out
}

func TestGet_MissOnEmpty(t *testing.T) {
	s := cache.New(time.Minute)
	_, ok := s.Get("svc")
	if ok {
		t.Fatal("expected cache miss on empty store")
	}
}

func TestSet_And_Get(t *testing.T) {
	s := cache.New(time.Minute)
	summaries := makeSummaries("api", "worker")
	s.Set("run1", summaries)
	got, ok := s.Get("run1")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(got))
	}
}

func TestGet_ExpiredEntry_ReturnsMiss(t *testing.T) {
	s := cache.New(time.Millisecond)
	s.Set("svc", makeSummaries("api"))
	time.Sleep(5 * time.Millisecond)
	_, ok := s.Get("svc")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	s := cache.New(time.Minute)
	s.Set("svc", makeSummaries("api"))
	s.Invalidate("svc")
	_, ok := s.Get("svc")
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	s := cache.New(time.Minute)
	s.Set("a", makeSummaries("x"))
	s.Set("b", makeSummaries("y"))
	s.Flush()
	if s.Len() != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", s.Len())
	}
}

func TestLen_ReflectsEntryCount(t *testing.T) {
	s := cache.New(time.Minute)
	if s.Len() != 0 {
		t.Fatal("expected 0 initially")
	}
	s.Set("a", makeSummaries("x"))
	s.Set("b", makeSummaries("y"))
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}
