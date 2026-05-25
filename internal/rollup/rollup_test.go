package rollup_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/rollup"
)

func makeSummaries() []drift.Summary {
	return []drift.Summary{
		{
			ServiceName: "api",
			Diffs: []drift.FieldDiff{
				{Field: "image", Declared: "v1", Live: "v2", Severity: "warn"},
			},
		},
		{
			ServiceName: "worker",
			Diffs: []drift.FieldDiff{
				{Field: "replicas", Declared: "3", Live: "1", Severity: "error"},
			},
		},
		{
			ServiceName: "cache",
			Diffs:       nil,
		},
	}
}

func TestBuild_TotalServices(t *testing.T) {
	r := rollup.Build(makeSummaries())
	if r.TotalServices != 3 {
		t.Errorf("expected 3 total services, got %d", r.TotalServices)
	}
}

func TestBuild_DriftedCount(t *testing.T) {
	r := rollup.Build(makeSummaries())
	if r.DriftedCount != 2 {
		t.Errorf("expected 2 drifted services, got %d", r.DriftedCount)
	}
}

func TestBuild_TopSeverityError(t *testing.T) {
	r := rollup.Build(makeSummaries())
	if r.TopSeverity != rollup.SeverityError {
		t.Errorf("expected top severity 'error', got %q", r.TopSeverity)
	}
}

func TestBuild_NoDrift_TopSeverityOK(t *testing.T) {
	summaries := []drift.Summary{
		{ServiceName: "svc", Diffs: nil},
	}
	r := rollup.Build(summaries)
	if r.TopSeverity != rollup.SeverityOK {
		t.Errorf("expected 'ok', got %q", r.TopSeverity)
	}
	if r.DriftedCount != 0 {
		t.Errorf("expected 0 drifted, got %d", r.DriftedCount)
	}
}

func TestBuild_ServiceRollup_FieldNames(t *testing.T) {
	r := rollup.Build(makeSummaries())
	var api rollup.ServiceRollup
	for _, s := range r.Services {
		if s.ServiceName == "api" {
			api = s
		}
	}
	if len(api.DriftedFields) != 1 || api.DriftedFields[0] != "image" {
		t.Errorf("unexpected drifted fields: %v", api.DriftedFields)
	}
}

func TestBuild_Empty(t *testing.T) {
	r := rollup.Build(nil)
	if r.TotalServices != 0 || r.DriftedCount != 0 {
		t.Errorf("expected empty report, got %+v", r)
	}
}
