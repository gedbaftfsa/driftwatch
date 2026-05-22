package drift

import (
	"strings"
	"testing"
)

func TestDetect_NoDrift(t *testing.T) {
	declared := []ServiceState{
		{Name: "api", Image: "api:v1", Replicas: 2, Env: map[string]string{"PORT": "8080"}},
	}
	deployed := map[string]DeployedState{
		"api": {Name: "api", Image: "api:v1", Replicas: 2, Env: map[string]string{"PORT": "8080"}},
	}

	results, err := Detect(declared, deployed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift, got diffs: %v", results[0].Diffs)
	}
}

func TestDetect_ImageDrift(t *testing.T) {
	declared := []ServiceState{
		{Name: "api", Image: "api:v2", Replicas: 1, Env: map[string]string{}},
	}
	deployed := map[string]DeployedState{
		"api": {Name: "api", Image: "api:v1", Replicas: 1, Env: map[string]string{}},
	}

	results, err := Detect(declared, deployed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Drifted {
		t.Error("expected drift due to image mismatch")
	}
	if len(results[0].Diffs) != 1 || !strings.Contains(results[0].Diffs[0], "image") {
		t.Errorf("expected image diff, got: %v", results[0].Diffs)
	}
}

func TestDetect_ServiceMissing(t *testing.T) {
	declared := []ServiceState{
		{Name: "worker", Image: "worker:latest", Replicas: 3, Env: map[string]string{}},
	}
	deployed := map[string]DeployedState{}

	results, err := Detect(declared, deployed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Drifted {
		t.Error("expected drift for missing service")
	}
	if !strings.Contains(results[0].Diffs[0], "not found") {
		t.Errorf("expected 'not found' message, got: %v", results[0].Diffs)
	}
}

func TestDetect_EnvDrift(t *testing.T) {
	declared := []ServiceState{
		{Name: "svc", Image: "svc:v1", Replicas: 1, Env: map[string]string{"LOG_LEVEL": "debug"}},
	}
	deployed := map[string]DeployedState{
		"svc": {Name: "svc", Image: "svc:v1", Replicas: 1, Env: map[string]string{"LOG_LEVEL": "info"}},
	}

	results, err := Detect(declared, deployed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Drifted {
		t.Error("expected drift due to env mismatch")
	}
}

func TestDetect_NilDeclared(t *testing.T) {
	_, err := Detect(nil, map[string]DeployedState{})
	if err == nil {
		t.Error("expected error for nil declared states")
	}
}

func TestSummary(t *testing.T) {
	results := []DriftResult{
		{ServiceName: "api", Drifted: false},
		{ServiceName: "worker", Drifted: true, Diffs: []string{"replicas: declared=3 actual=1"}},
	}
	out := Summary(results)
	if !strings.Contains(out, "[OK]    api") {
		t.Errorf("expected OK line for api, got: %s", out)
	}
	if !strings.Contains(out, "[DRIFT] worker") {
		t.Errorf("expected DRIFT line for worker, got: %s", out)
	}
	if !strings.Contains(out, "replicas") {
		t.Errorf("expected replicas diff in summary, got: %s", out)
	}
}
