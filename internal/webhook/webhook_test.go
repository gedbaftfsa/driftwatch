package webhook_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/webhook"
)

func makeSummaries(hasDrift bool) []drift.Summary {
	return []drift.Summary{
		{ServiceName: "api", HasDrift: hasDrift},
		{ServiceName: "worker", HasDrift: false},
	}
}

func TestSend_Success(t *testing.T) {
	var received webhook.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := webhook.NewSender(ts.URL, "", time.Second)
	err := s.Send(context.Background(), makeSummaries(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.TotalDrift != 1 {
		t.Errorf("expected TotalDrift=1, got %d", received.TotalDrift)
	}
	if len(received.Summaries) != 2 {
		t.Errorf("expected 2 summaries, got %d", len(received.Summaries))
	}
}

func TestSend_SecretHeader(t *testing.T) {
	var gotSecret string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSecret = r.Header.Get("X-Driftwatch-Secret")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	s := webhook.NewSender(ts.URL, "mysecret", time.Second)
	_ = s.Send(context.Background(), makeSummaries(false))
	if gotSecret != "mysecret" {
		t.Errorf("expected secret header 'mysecret', got %q", gotSecret)
	}
}

func TestSend_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	s := webhook.NewSender(ts.URL, "", time.Second)
	err := s.Send(context.Background(), makeSummaries(false))
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestSend_InvalidEndpoint(t *testing.T) {
	s := webhook.NewSender("http://127.0.0.1:0", "", 200*time.Millisecond)
	err := s.Send(context.Background(), makeSummaries(false))
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}

func TestSend_NoDriftCountsZero(t *testing.T) {
	var received webhook.Payload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	s := webhook.NewSender(ts.URL, "", time.Second)
	_ = s.Send(context.Background(), makeSummaries(false))
	if received.TotalDrift != 0 {
		t.Errorf("expected TotalDrift=0, got %d", received.TotalDrift)
	}
}
