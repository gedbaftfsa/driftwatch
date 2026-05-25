// Package rollup aggregates drift summaries across multiple services
// into a single rolled-up report for high-level observability.
package rollup

import (
	"time"

	"github.com/example/driftwatch/internal/drift"
)

// Severity levels for rollup classification.
const (
	SeverityOK    = "ok"
	SeverityWarn  = "warn"
	SeverityError = "error"
)

// ServiceRollup holds the rolled-up state for a single service.
type ServiceRollup struct {
	ServiceName  string
	DriftCount   int
	TopSeverity  string
	DriftedFields []string
	At           time.Time
}

// Report is the aggregated rollup across all services.
type Report struct {
	GeneratedAt   time.Time
	TotalServices int
	DriftedCount  int
	TopSeverity   string
	Services      []ServiceRollup
}

// Build constructs a rollup Report from a slice of drift summaries.
func Build(summaries []drift.Summary) Report {
	now := time.Now().UTC()
	report := Report{
		GeneratedAt:   now,
		TotalServices: len(summaries),
	}

	overall := SeverityOK

	for _, s := range summaries {
		sr := ServiceRollup{
			ServiceName: s.ServiceName,
			At:          now,
		}

		if len(s.Diffs) == 0 {
			sr.TopSeverity = SeverityOK
		} else {
			report.DriftedCount++
			sr.DriftCount = len(s.Diffs)
			sr.TopSeverity = topSeverity(s.Diffs)
			sr.DriftedFields = fieldNames(s.Diffs)
			overall = higherSeverity(overall, sr.TopSeverity)
		}

		report.Services = append(report.Services, sr)
	}

	report.TopSeverity = overall
	return report
}

func topSeverity(diffs []drift.FieldDiff) string {
	top := SeverityWarn
	for _, d := range diffs {
		if d.Severity == SeverityError {
			return SeverityError
		}
	}
	return top
}

func fieldNames(diffs []drift.FieldDiff) []string {
	names := make([]string, 0, len(diffs))
	for _, d := range diffs {
		names = append(names, d.Field)
	}
	return names
}

func higherSeverity(a, b string) string {
	rank := map[string]int{SeverityOK: 0, SeverityWarn: 1, SeverityError: 2}
	if rank[b] > rank[a] {
		return b
	}
	return a
}
