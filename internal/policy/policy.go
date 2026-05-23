package policy

import (
	"fmt"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Rule defines a named policy rule applied to drift summaries.
type Rule struct {
	Name     string
	Severity string // "warn" or "error"
	Field    string // "image", "replicas", "env"
	Action   string // "deny" or "warn"
}

// Policy holds a set of rules to evaluate against drift summaries.
type Policy struct {
	Rules []Rule
}

// Violation represents a policy rule that was triggered.
type Violation struct {
	Service string
	Rule    string
	Message string
}

// Evaluate applies the policy rules to a slice of drift summaries and
// returns any violations found.
func (p *Policy) Evaluate(summaries []drift.Summary) []Violation {
	var violations []Violation
	for _, s := range summaries {
		for _, r := range p.Rules {
			if !matchesSeverity(s, r.Severity) {
				continue
			}
			if !matchesField(s, r.Field) {
				continue
			}
			violations = append(violations, Violation{
				Service: s.Service,
				Rule:    r.Name,
				Message: fmt.Sprintf("service %q violates rule %q: %s drift detected (severity=%s)", s.Service, r.Name, r.Field, r.Severity),
			})
		}
	}
	return violations
}

// HasViolations returns true if any violations exist.
func HasViolations(violations []Violation) bool {
	return len(violations) > 0
}

func matchesSeverity(s drift.Summary, severity string) bool {
	if severity == "" {
		return true
	}
	return strings.EqualFold(s.Severity, severity)
}

func matchesField(s drift.Summary, field string) bool {
	if field == "" {
		return true
	}
	for _, d := range s.Diffs {
		if strings.EqualFold(d.Field, field) {
			return true
		}
	}
	return false
}
