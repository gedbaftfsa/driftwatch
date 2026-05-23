package filter_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/filter"
)

func makeSummaries() []drift.Summary {
	return []drift.Summary{
		{ServiceName: "api", HasDrift: false, Severity: ""},
		{ServiceName: "worker", HasDrift: true, Severity: "warn"},
		{ServiceName: "gateway", HasDrift: true, Severity: "error"},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	summaries := makeSummaries()
	result := filter.Apply(summaries, filter.Options{})
	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}
}

func TestApply_OnlyDrifted(t *testing.T) {
	result := filter.Apply(makeSummaries(), filter.Options{OnlyDrifted: true})
	if len(result) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(result))
	}
	for _, s := range result {
		if !s.HasDrift {
			t.Errorf("expected only drifted summaries, got %+v", s)
		}
	}
}

func TestApply_MinSeverityWarn(t *testing.T) {
	result := filter.Apply(makeSummaries(), filter.Options{MinSeverity: "warn"})
	if len(result) != 2 {
		t.Fatalf("expected 2 results for min severity warn, got %d", len(result))
	}
}

func TestApply_MinSeverityError(t *testing.T) {
	result := filter.Apply(makeSummaries(), filter.Options{MinSeverity: "error"})
	if len(result) != 1 {
		t.Fatalf("expected 1 result for min severity error, got %d", len(result))
	}
	if result[0].ServiceName != "gateway" {
		t.Errorf("expected gateway, got %s", result[0].ServiceName)
	}
}

func TestApply_ServiceFilter(t *testing.T) {
	result := filter.Apply(makeSummaries(), filter.Options{Services: []string{"api", "worker"}})
	if len(result) != 2 {
		t.Fatalf("expected 2 results for service filter, got %d", len(result))
	}
}

func TestApply_ServiceFilter_CaseInsensitive(t *testing.T) {
	result := filter.Apply(makeSummaries(), filter.Options{Services: []string{"GATEWAY"}})
	if len(result) != 1 || result[0].ServiceName != "gateway" {
		t.Errorf("expected gateway (case-insensitive match), got %+v", result)
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	opts := filter.Options{
		Services:    []string{"worker", "gateway"},
		OnlyDrifted: true,
		MinSeverity: "error",
	}
	result := filter.Apply(makeSummaries(), opts)
	if len(result) != 1 || result[0].ServiceName != "gateway" {
		t.Errorf("expected only gateway with combined filters, got %+v", result)
	}
}
