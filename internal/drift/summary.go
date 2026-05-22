package drift

// Summary holds the result of comparing a single service's
// declared config against its live runtime state.
type Summary struct {
	ServiceName string
	Drifted     bool
	Diffs       []string
}

// Summarize converts a slice of DetectResult into a slice of Summary,
// keeping only entries that have at least one diff or are explicitly missing.
func Summarize(results []DetectResult) []Summary {
	var summaries []Summary
	for _, r := range results {
		s := Summary{
			ServiceName: r.ServiceName,
			Drifted:     len(r.Diffs) > 0,
			Diffs:       r.Diffs,
		}
		summaries = append(summaries, s)
	}
	return summaries
}
