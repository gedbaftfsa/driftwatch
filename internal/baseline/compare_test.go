package baseline_test

import (
	"testing"

	"github.com/example/driftwatch/internal/baseline"
)

func TestCompare_NoDrift(t *testing.T) {
	b := makeBaseline()
	current := map[string]baseline.Service{
		"api": {Name: "api", Image: "api:1.0", Replicas: 2, Environment: map[string]string{"LOG_LEVEL": "info"}},
	}
	deltas := baseline.Compare(b, current)
	if len(deltas) != 0 {
		t.Errorf("expected no deltas, got %d: %v", len(deltas), deltas)
	}
}

func TestCompare_ImageChanged(t *testing.T) {
	b := makeBaseline()
	current := map[string]baseline.Service{
		"api": {Name: "api", Image: "api:2.0", Replicas: 2, Environment: map[string]string{"LOG_LEVEL": "info"}},
	}
	deltas := baseline.Compare(b, current)
	if len(deltas) != 1 || deltas[0].Field != "image" {
		t.Errorf("expected image delta, got %v", deltas)
	}
}

func TestCompare_ReplicasDrift(t *testing.T) {
	b := makeBaseline()
	current := map[string]baseline.Service{
		"api": {Name: "api", Image: "api:1.0", Replicas: 5, Environment: map[string]string{"LOG_LEVEL": "info"}},
	}
	deltas := baseline.Compare(b, current)
	if len(deltas) != 1 || deltas[0].Field != "replicas" {
		t.Errorf("expected replicas delta, got %v", deltas)
	}
}

func TestCompare_ServiceMissing(t *testing.T) {
	b := makeBaseline()
	current := map[string]baseline.Service{}
	deltas := baseline.Compare(b, current)
	if len(deltas) != 1 || deltas[0].Field != "existence" {
		t.Errorf("expected existence delta for missing service, got %v", deltas)
	}
}

func TestCompare_NewServiceAppeared(t *testing.T) {
	b := makeBaseline()
	current := map[string]baseline.Service{
		"api":    {Name: "api", Image: "api:1.0", Replicas: 2, Environment: map[string]string{"LOG_LEVEL": "info"}},
		"worker": {Name: "worker", Image: "worker:1.0", Replicas: 1},
	}
	deltas := baseline.Compare(b, current)
	if len(deltas) != 1 || deltas[0].Service != "worker" {
		t.Errorf("expected delta for new service 'worker', got %v", deltas)
	}
}

func TestCompare_EnvDrift(t *testing.T) {
	b := makeBaseline()
	current := map[string]baseline.Service{
		"api": {Name: "api", Image: "api:1.0", Replicas: 2, Environment: map[string]string{"LOG_LEVEL": "debug"}},
	}
	deltas := baseline.Compare(b, current)
	if len(deltas) != 1 || deltas[0].Field != "env:LOG_LEVEL" {
		t.Errorf("expected env delta, got %v", deltas)
	}
}
