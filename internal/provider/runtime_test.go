package provider_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/driftwatch/internal/provider"
)

func makeTestStates() []provider.ServiceState {
	return []provider.ServiceState{
		{
			Name:        "api",
			Image:       "myrepo/api:v1.2.3",
			Replicas:    3,
			Env:         map[string]string{"PORT": "8080", "LOG_LEVEL": "info"},
			LastUpdated: time.Now(),
		},
		{
			Name:        "worker",
			Image:       "myrepo/worker:v0.9.0",
			Replicas:    1,
			Env:         map[string]string{"QUEUE": "default"},
			LastUpdated: time.Now(),
		},
	}
}

func TestFetchAll_Success(t *testing.T) {
	states := makeTestStates()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(states)
	}))
	defer ts.Close()

	p := provider.NewRuntimeProvider(ts.URL)
	got, err := p.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 services, got %d", len(got))
	}
	if got[0].Name != "api" {
		t.Errorf("expected first service name 'api', got %q", got[0].Name)
	}
}

func TestFetchAll_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	p := provider.NewRuntimeProvider(ts.URL)
	_, err := p.FetchAll(context.Background())
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestFetchByName_Found(t *testing.T) {
	state := makeTestStates()[0]
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state)
	}))
	defer ts.Close()

	p := provider.NewRuntimeProvider(ts.URL)
	got, err := p.FetchByName(context.Background(), "api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected a service state, got nil")
	}
	if got.Image != "myrepo/api:v1.2.3" {
		t.Errorf("unexpected image: %q", got.Image)
	}
}

func TestFetchByName_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	p := provider.NewRuntimeProvider(ts.URL)
	got, err := p.FetchByName(context.Background(), "missing-service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for missing service, got %+v", got)
	}
}
