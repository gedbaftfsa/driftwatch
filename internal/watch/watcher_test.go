package watch_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/provider"
	"github.com/driftwatch/internal/watch"
)

func silentLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func makeServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func TestWatcher_CallsHandlerOnDrift(t *testing.T) {
	body := `[{"name":"api","image":"nginx:1.25","replicas":2,"env":{}}]`
	srv := makeServer(t, body)
	defer srv.Close()

	p := provider.NewRuntimeProvider(srv.URL, http.DefaultClient)

	declared := []drift.ServiceSpec{
		{Name: "api", Image: "nginx:1.20", Replicas: 2, Env: map[string]string{}},
	}

	var called int32
	handler := func(results []drift.Result) {
		atomic.AddInt32(&called, 1)
	}

	cfg := watch.Config{Interval: 50 * time.Millisecond}
	w := watch.New(cfg, p, handler, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx, declared)

	if atomic.LoadInt32(&called) == 0 {
		t.Error("expected handler to be called on drift")
	}
}

func TestWatcher_NoHandlerOnNoDrift(t *testing.T) {
	body := `[{"name":"api","image":"nginx:1.20","replicas":2,"env":{}}]`
	srv := makeServer(t, body)
	defer srv.Close()

	p := provider.NewRuntimeProvider(srv.URL, http.DefaultClient)

	declared := []drift.ServiceSpec{
		{Name: "api", Image: "nginx:1.20", Replicas: 2, Env: map[string]string{}},
	}

	var called int32
	handler := func(results []drift.Result) {
		atomic.AddInt32(&called, 1)
	}

	cfg := watch.Config{Interval: 50 * time.Millisecond}
	w := watch.New(cfg, p, handler, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	_ = w.Run(ctx, declared)

	if atomic.LoadInt32(&called) != 0 {
		t.Errorf("expected no handler calls, got %d", called)
	}
}

func TestWatcher_StopsOnContextCancel(t *testing.T) {
	body := `[]`
	srv := makeServer(t, body)
	defer srv.Close()

	p := provider.NewRuntimeProvider(srv.URL, http.DefaultClient)
	cfg := watch.Config{Interval: 10 * time.Millisecond}
	w := watch.New(cfg, p, func(_ []drift.Result) {}, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := w.Run(ctx, nil)
	if err == nil {
		t.Error("expected context error")
	}
}
