package drift

// Summary aggregates the results of a drift detection run.
type Summary struct {
	Total   int       `json:"total"`
	Drifted int       `json:"drifted"`
	Clean   int       `json:"clean"`
	Results []Result  `json:"results"`
}

// Result holds drift information for a single service.
type Result struct {
	Service string `json:"service"`
	Drifted bool   `json:"drifted"`
	Reasons []string `json:"reasons,omitempty"`
}

// Summarize builds a Summary from a slice of Results.
func Summarize(results []Result) Summary {
	s := Summary{
		Total:   len(results),
		Results: results,
	}
	for _, r := range results {
		if r.Drifted {
			s.Drifted++
		} else {
			s.Clean++
		}
	}
	return s
}

// HasDrift returns true if any service in the summary has drifted.
func (s Summary) HasDrift() bool {
	return s.Drifted > 0
}
