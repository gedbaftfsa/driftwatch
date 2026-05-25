package watch_test

import (
	"testing"
	"time"

	"github.com/driftwatch/internal/watch"
)

func TestDefaultConfig_Defaults(t *testing.T) {
	cfg := watch.DefaultConfig()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if len(cfg.ServiceNames) != 0 {
		t.Errorf("expected no service names, got %v", cfg.ServiceNames)
	}
}

func TestWithInterval_SetsInterval(t *testing.T) {
	cfg := watch.DefaultConfig(watch.WithInterval(5 * time.Minute))
	if cfg.Interval != 5*time.Minute {
		t.Errorf("expected 5m, got %v", cfg.Interval)
	}
}

func TestWithInterval_ZeroIgnored(t *testing.T) {
	cfg := watch.DefaultConfig(watch.WithInterval(0))
	if cfg.Interval != 30*time.Second {
		t.Errorf("zero interval should not override default, got %v", cfg.Interval)
	}
}

func TestWithServiceFilter_AddsNames(t *testing.T) {
	cfg := watch.DefaultConfig(watch.WithServiceFilter("api", "worker"))
	if len(cfg.ServiceNames) != 2 {
		t.Fatalf("expected 2 names, got %d", len(cfg.ServiceNames))
	}
	if cfg.ServiceNames[0] != "api" || cfg.ServiceNames[1] != "worker" {
		t.Errorf("unexpected names: %v", cfg.ServiceNames)
	}
}

func TestWithServiceFilter_MultipleCallsAppend(t *testing.T) {
	cfg := watch.DefaultConfig(
		watch.WithServiceFilter("api"),
		watch.WithServiceFilter("worker"),
	)
	if len(cfg.ServiceNames) != 2 {
		t.Errorf("expected 2 names, got %d", len(cfg.ServiceNames))
	}
}
