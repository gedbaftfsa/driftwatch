package webhook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/webhook"
)

func noSleep(_ time.Duration) {}

func TestRetrySender_SucceedsFirstAttempt(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	inner := webhook.NewSender(ts.URL, "", time.Second)
	cfg := webhook.RetryConfig{MaxAttempts: 3, Delay: 0}
	rs := webhook.NewRetrySender(inner, cfg)
	rs.SetSleep(noSleep)

	err := rs.Send(context.Background(), []drift.Summary{{ServiceName: "svc"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetrySender_RetriesOnFailure(t *testing.T) {
	var calls int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	inner := webhook.NewSender(ts.URL, "", time.Second)
	cfg := webhook.RetryConfig{MaxAttempts: 3, Delay: 0}
	rs := webhook.NewRetrySender(inner, cfg)
	rs.SetSleep(noSleep)

	err := rs.Send(context.Background(), []drift.Summary{})
	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetrySender_ExhaustsAttempts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer ts.Close()

	inner := webhook.NewSender(ts.URL, "", time.Second)
	cfg := webhook.RetryConfig{MaxAttempts: 2, Delay: 0}
	rs := webhook.NewRetrySender(inner, cfg)
	rs.SetSleep(noSleep)

	err := rs.Send(context.Background(), []drift.Summary{})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
}

func TestRetrySender_ContextCancelled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	inner := webhook.NewSender(ts.URL, "", time.Second)
	cfg := webhook.RetryConfig{MaxAttempts: 3, Delay: 0}
	rs := webhook.NewRetrySender(inner, cfg)
	rs.SetSleep(noSleep)

	err := rs.Send(ctx, []drift.Summary{})
	if err == nil {
		t.Fatal("expected context cancellation error")
	}
}
