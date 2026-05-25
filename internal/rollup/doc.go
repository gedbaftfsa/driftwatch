// Package rollup provides aggregation of per-service drift summaries into
// a single high-level report suitable for dashboards and alerting pipelines.
//
// Usage:
//
//	// Collect drift summaries from the detector
//	summaries := detector.Detect(declared, live)
//
//	// Build a rolled-up report
//	report := rollup.Build(summaries)
//
//	// Render to stdout as plain text
//	rollup.Render(os.Stdout, report, rollup.FormatText)
//
//	// Or emit as JSON for downstream consumers
//	rollup.Render(w, report, rollup.FormatJSON)
//
// The Report captures overall severity (ok / warn / error), total service
// count, and per-service breakdowns including drifted field names.
package rollup
