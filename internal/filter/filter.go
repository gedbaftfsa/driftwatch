package filter

import (
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Options holds filtering criteria for drift summaries.
type Options struct {
	// Services restricts results to the named services (empty = all).
	Services []string
	// OnlyDrifted, when true, excludes summaries with no drift detected.
	OnlyDrifted bool
	// MinSeverity filters out summaries whose highest severity is below this
	// level. Accepted values: "", "warn", "error".
	MinSeverity string
}

// Apply returns the subset of summaries that match all criteria in opts.
func Apply(summaries []drift.Summary, opts Options) []drift.Summary {
	var out []drift.Summary
	for _, s := range summaries {
		if !matchesServiceFilter(s, opts.Services) {
			continue
		}
		if opts.OnlyDrifted && !s.HasDrift {
			continue
		}
		if !matchesSeverity(s, opts.MinSeverity) {
			continue
		}
		out = append(out, s)
	}
	return out
}

func matchesServiceFilter(s drift.Summary, services []string) bool {
	if len(services) == 0 {
		return true
	}
	for _, name := range services {
		if strings.EqualFold(s.ServiceName, name) {
			return true
		}
	}
	return false
}

func matchesSeverity(s drift.Summary, min string) bool {
	switch strings.ToLower(min) {
	case "error":
		return s.Severity == "error"
	case "warn":
		return s.Severity == "warn" || s.Severity == "error"
	default:
		return true
	}
}
