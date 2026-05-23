package policy_test

import (
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/policy"
)

func makeSummaries() []drift.Summary {
	return []drift.Summary{
		{
			Service:  "api",
			Drifted:  true,
			Severity: "error",
			Diffs: []drift.Diff{
				{Field: "image", Declared: "nginx:1.24", Actual: "nginx:1.25"},
			},
		},
		{
			Service:  "worker",
			Drifted:  true,
			Severity: "warn",
			Diffs: []drift.Diff{
				{Field: "replicas", Declared: "3", Actual: "2"},
			},
		},
		{
			Service:  "cache",
			Drifted:  false,
			Severity: "",
		},
	}
}

func TestEvaluate_NoRules_NoViolations(t *testing.T) {
	p := &policy.Policy{}
	violations := p.Evaluate(makeSummaries())
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestEvaluate_MatchBySeverityError(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "no-image-drift", Severity: "error", Field: "image", Action: "deny"},
		},
	}
	violations := p.Evaluate(makeSummaries())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Service != "api" {
		t.Errorf("expected service 'api', got %q", violations[0].Service)
	}
	if violations[0].Rule != "no-image-drift" {
		t.Errorf("expected rule 'no-image-drift', got %q", violations[0].Rule)
	}
}

func TestEvaluate_MatchByField_Replicas(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "replica-check", Severity: "warn", Field: "replicas", Action: "warn"},
		},
	}
	violations := p.Evaluate(makeSummaries())
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Service != "worker" {
		t.Errorf("expected service 'worker', got %q", violations[0].Service)
	}
}

func TestEvaluate_NoMatchOnCleanService(t *testing.T) {
	p := &policy.Policy{
		Rules: []policy.Rule{
			{Name: "any-drift", Severity: "error", Field: "", Action: "deny"},
		},
	}
	violations := p.Evaluate(makeSummaries())
	// only 'api' is severity=error
	for _, v := range violations {
		if v.Service == "cache" {
			t.Errorf("clean service 'cache' should not produce a violation")
		}
	}
}

func TestHasViolations(t *testing.T) {
	if policy.HasViolations(nil) {
		t.Error("expected false for nil violations")
	}
	if !policy.HasViolations([]policy.Violation{{Service: "api", Rule: "r", Message: "m"}}) {
		t.Error("expected true for non-empty violations")
	}
}
