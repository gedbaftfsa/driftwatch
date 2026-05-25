package ratelimit_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/ratelimit"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := ratelimit.New(time.Second, 1)
	if !l.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_ExceedsBurst_Blocked(t *testing.T) {
	l := ratelimit.New(time.Hour, 2)
	if !l.Allow("svc-a") {
		t.Fatal("expected first call allowed")
	}
	if !l.Allow("svc-a") {
		t.Fatal("expected second call allowed (burst=2)")
	}
	if l.Allow("svc-a") {
		t.Fatal("expected third call to be blocked")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.New(time.Hour, 1)
	l.Allow("svc-a") // consume token for svc-a
	if !l.Allow("svc-b") {
		t.Fatal("svc-b should have its own bucket")
	}
}

func TestAllow_TokensRefillOverTime(t *testing.T) {
	l := ratelimit.New(10*time.Millisecond, 1)
	l.Allow("svc-a") // consume
	if l.Allow("svc-a") {
		t.Fatal("should be blocked immediately after consuming")
	}
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("svc-a") {
		t.Fatal("expected token to refill after rate interval")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	l := ratelimit.New(time.Hour, 1)
	l.Allow("svc-a") // consume
	l.Reset("svc-a")
	if !l.Allow("svc-a") {
		t.Fatal("expected allow after reset")
	}
}

func TestRemaining_FullBucketForNewKey(t *testing.T) {
	l := ratelimit.New(time.Hour, 3)
	if got := l.Remaining("svc-new"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestRemaining_DecreasesAfterAllow(t *testing.T) {
	l := ratelimit.New(time.Hour, 3)
	l.Allow("svc-a")
	if got := l.Remaining("svc-a"); got != 2 {
		t.Fatalf("expected 2 remaining, got %d", got)
	}
}

func TestNew_MinBurstIsOne(t *testing.T) {
	l := ratelimit.New(time.Second, 0)
	if !l.Allow("svc-a") {
		t.Fatal("expected at least 1 burst even when 0 specified")
	}
}
