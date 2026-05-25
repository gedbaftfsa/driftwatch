package resolve_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/provider"
	"github.com/yourorg/driftwatch/internal/resolve"
)

func makeServer(t *testing.T, states []provider.ServiceState) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(states)
	}))
}

func TestResolve_JoinsMatchingServices(t *testing.T) {
	svr := makeServer(t, []provider.ServiceState{
		{Name: "api", Image: "api:v2", Replicas: 2},
	})
	defer svr.Close()

	rt := provider.NewRuntimeProvider(svr.URL)
	res := resolve.New(rt)

	declared := []config.Service{{Name: "api", Image: "api:v1", Replicas: 1}}
	states, err := res.Resolve(context.Background(), declared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected 1 state, got %d", len(states))
	}
	if states[0].Live == nil {
		t.Fatal("expected Live to be populated")
	}
	if states[0].Live.Image != "api:v2" {
		t.Errorf("expected live image api:v2, got %s", states[0].Live.Image)
	}
}

func TestResolve_MissingLiveService(t *testing.T) {
	svr := makeServer(t, []provider.ServiceState{})
	defer svr.Close()

	rt := provider.NewRuntimeProvider(svr.URL)
	res := resolve.New(rt)

	declared := []config.Service{{Name: "worker", Image: "worker:v1", Replicas: 1}}
	states, err := res.Resolve(context.Background(), declared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected 1 state, got %d", len(states))
	}
	if states[0].Live != nil {
		t.Error("expected Live to be nil for missing service")
	}
}

func TestResolve_ExtraLiveServicesIgnored(t *testing.T) {
	svr := makeServer(t, []provider.ServiceState{
		{Name: "api", Image: "api:v1", Replicas: 1},
		{Name: "ghost", Image: "ghost:v9", Replicas: 3},
	})
	defer svr.Close()

	rt := provider.NewRuntimeProvider(svr.URL)
	res := resolve.New(rt)

	declared := []config.Service{{Name: "api", Image: "api:v1", Replicas: 1}}
	states, err := res.Resolve(context.Background(), declared)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(states) != 1 {
		t.Errorf("expected 1 state (ghost ignored), got %d", len(states))
	}
}

func TestResolve_ServerError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	rt := provider.NewRuntimeProvider(svr.URL)
	res := resolve.New(rt)

	_, err := res.Resolve(context.Background(), []config.Service{{Name: "api"}})
	if err == nil {
		t.Fatal("expected error from server failure")
	}
}
