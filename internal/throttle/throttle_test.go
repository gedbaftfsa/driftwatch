package throttle_test

import (
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/throttle"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	th := throttle.New(5 * time.Second)
	if !th.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("svc-a")
	if th.Allow("svc-a") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_DifferentKeysIndependent(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("svc-a")
	if !th.Allow("svc-b") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestAllow_AfterCooldown_Permitted(t *testing.T) {
	th := throttle.New(10 * time.Millisecond)
	th.Allow("svc-a")
	time.Sleep(20 * time.Millisecond)
	if !th.Allow("svc-a") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("svc-a")
	th.Reset("svc-a")
	if !th.Allow("svc-a") {
		t.Fatal("expected reset key to be allowed immediately")
	}
}

func TestResetAll_ClearsAllKeys(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("svc-a")
	th.Allow("svc-b")
	th.ResetAll()
	if !th.Allow("svc-a") || !th.Allow("svc-b") {
		t.Fatal("expected all keys to be cleared after ResetAll")
	}
}

func TestActiveKeys_ReturnsTrackedKeys(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("svc-a")
	th.Allow("svc-b")
	keys := th.ActiveKeys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 active keys, got %d", len(keys))
	}
}

func TestActiveKeys_EmptyAfterResetAll(t *testing.T) {
	th := throttle.New(5 * time.Second)
	th.Allow("svc-a")
	th.ResetAll()
	if len(th.ActiveKeys()) != 0 {
		t.Fatal("expected no active keys after ResetAll")
	}
}
