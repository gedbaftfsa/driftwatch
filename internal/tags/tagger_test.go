package tags

import (
	"testing"
)

func TestNewTagSet_Empty(t *testing.T) {
	ts := NewTagSet("svc-a")
	if ts.Service != "svc-a" {
		t.Fatalf("expected service svc-a, got %s", ts.Service)
	}
	if len(ts.Tags) != 0 {
		t.Fatalf("expected empty tags, got %d", len(ts.Tags))
	}
}

func TestSet_And_Get(t *testing.T) {
	ts := NewTagSet("svc-b")
	if err := ts.Set("env", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := ts.Get("env")
	if !ok || v != "production" {
		t.Fatalf("expected production, got %q (found=%v)", v, ok)
	}
}

func TestSet_EmptyKey_Error(t *testing.T) {
	ts := NewTagSet("svc-c")
	if err := ts.Set("  ", "val"); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestDelete_RemovesTag(t *testing.T) {
	ts := NewTagSet("svc-d")
	_ = ts.Set("team", "platform")
	ts.Delete("team")
	_, ok := ts.Get("team")
	if ok {
		t.Fatal("expected tag to be deleted")
	}
}

func TestKeys_Sorted(t *testing.T) {
	ts := NewTagSet("svc-e")
	_ = ts.Set("zone", "us-east")
	_ = ts.Set("env", "staging")
	_ = ts.Set("team", "infra")
	keys := ts.Keys()
	expected := []string{"env", "team", "zone"}
	for i, k := range keys {
		if k != expected[i] {
			t.Fatalf("keys mismatch at %d: want %s got %s", i, expected[i], k)
		}
	}
}

func TestMatches_AllPresent(t *testing.T) {
	ts := NewTagSet("svc-f")
	_ = ts.Set("env", "prod")
	_ = ts.Set("region", "eu-west")
	if !ts.Matches(map[string]string{"env": "prod", "region": "eu-west"}) {
		t.Fatal("expected match")
	}
}

func TestMatches_PartialFail(t *testing.T) {
	ts := NewTagSet("svc-g")
	_ = ts.Set("env", "prod")
	if ts.Matches(map[string]string{"env": "prod", "region": "eu-west"}) {
		t.Fatal("expected no match when filter has extra key")
	}
}

func TestString_Format(t *testing.T) {
	ts := NewTagSet("svc-h")
	_ = ts.Set("env", "dev")
	_ = ts.Set("team", "core")
	s := ts.String()
	expected := "svc-h[env=dev,team=core]"
	if s != expected {
		t.Fatalf("expected %q, got %q", expected, s)
	}
}
