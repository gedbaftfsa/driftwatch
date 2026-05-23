package metrics_test

import (
	"strings"
	"testing"

	"github.com/driftwatch/internal/metrics"
)

// makeRegistry returns a fresh Registry for use in tests.
func makeRegistry() *metrics.Registry {
	return metrics.NewRegistry()
}

func TestNewRegistry_IsNotNil(t *testing.T) {
	r := makeRegistry()
	if r == nil {
		t.Fatal("expected non-nil Registry")
	}
}

func TestRegistry_RecordRun_IncrementsTotal(t *testing.T) {
	r := makeRegistry()
	r.RecordRun()
	r.RecordRun()

	snapshot := r.Snapshot()
	if snapshot.TotalRuns != 2 {
		t.Errorf("expected TotalRuns=2, got %d", snapshot.TotalRuns)
	}
}

func TestRegistry_RecordDrift_IncrementsCounter(t *testing.T) {
	r := makeRegistry()
	r.RecordDrift("svc-a", "image")
	r.RecordDrift("svc-a", "replicas")
	r.RecordDrift("svc-b", "env")

	snapshot := r.Snapshot()
	if snapshot.TotalDriftEvents != 3 {
		t.Errorf("expected TotalDriftEvents=3, got %d", snapshot.TotalDriftEvents)
	}
}

func TestRegistry_RecordError_IncrementsErrorCount(t *testing.T) {
	r := makeRegistry()
	r.RecordError()
	r.RecordError()
	r.RecordError()

	snapshot := r.Snapshot()
	if snapshot.TotalErrors != 3 {
		t.Errorf("expected TotalErrors=3, got %d", snapshot.TotalErrors)
	}
}

func TestRegistry_Snapshot_ZeroValues(t *testing.T) {
	r := makeRegistry()
	snapshot := r.Snapshot()

	if snapshot.TotalRuns != 0 {
		t.Errorf("expected TotalRuns=0, got %d", snapshot.TotalRuns)
	}
	if snapshot.TotalDriftEvents != 0 {
		t.Errorf("expected TotalDriftEvents=0, got %d", snapshot.TotalDriftEvents)
	}
	if snapshot.TotalErrors != 0 {
		t.Errorf("expected TotalErrors=0, got %d", snapshot.TotalErrors)
	}
}

func TestRegistry_DriftByService_Tracked(t *testing.T) {
	r := makeRegistry()
	r.RecordDrift("alpha", "image")
	r.RecordDrift("alpha", "image")
	r.RecordDrift("beta", "env")

	snapshot := r.Snapshot()

	if snapshot.DriftByService["alpha"] != 2 {
		t.Errorf("expected alpha=2, got %d", snapshot.DriftByService["alpha"])
	}
	if snapshot.DriftByService["beta"] != 1 {
		t.Errorf("expected beta=1, got %d", snapshot.DriftByService["beta"])
	}
}

func TestRegistry_DriftByField_Tracked(t *testing.T) {
	r := makeRegistry()
	r.RecordDrift("svc-a", "image")
	r.RecordDrift("svc-b", "image")
	r.RecordDrift("svc-c", "replicas")

	snapshot := r.Snapshot()

	if snapshot.DriftByField["image"] != 2 {
		t.Errorf("expected image=2, got %d", snapshot.DriftByField["image"])
	}
	if snapshot.DriftByField["replicas"] != 1 {
		t.Errorf("expected replicas=1, got %d", snapshot.DriftByField["replicas"])
	}
}

func TestRegistry_Format_ContainsKeyMetrics(t *testing.T) {
	r := makeRegistry()
	r.RecordRun()
	r.RecordDrift("svc-x", "image")
	r.RecordError()

	output := r.Format()

	for _, want := range []string{"total_runs", "total_drift_events", "total_errors"} {
		if !strings.Contains(output, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, output)
		}
	}
}

func TestRegistry_Reset_ClearsAllCounters(t *testing.T) {
	r := makeRegistry()
	r.RecordRun()
	r.RecordDrift("svc", "image")
	r.RecordError()
	r.Reset()

	snapshot := r.Snapshot()
	if snapshot.TotalRuns != 0 || snapshot.TotalDriftEvents != 0 || snapshot.TotalErrors != 0 {
		t.Errorf("expected all counters reset to 0 after Reset(), got %+v", snapshot)
	}
	if len(snapshot.DriftByService) != 0 {
		t.Errorf("expected DriftByService to be empty after Reset()")
	}
	if len(snapshot.DriftByField) != 0 {
		t.Errorf("expected DriftByField to be empty after Reset()")
	}
}
