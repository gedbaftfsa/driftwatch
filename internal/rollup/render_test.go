package rollup_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/rollup"
)

func buildReport() rollup.Report {
	return rollup.Build([]drift.Summary{
		{
			ServiceName: "api",
			Diffs: []drift.FieldDiff{
				{Field: "image", Declared: "v1", Live: "v2", Severity: "warn"},
			},
		},
		{ServiceName: "db", Diffs: nil},
	})
}

func TestRender_Text_ContainsServiceNames(t *testing.T) {
	var buf bytes.Buffer
	if err := rollup.Render(&buf, buildReport(), rollup.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "api") {
		t.Error("expected 'api' in text output")
	}
	if !strings.Contains(out, "db") {
		t.Error("expected 'db' in text output")
	}
}

func TestRender_Text_ShowsDriftedFields(t *testing.T) {
	var buf bytes.Buffer
	rollup.Render(&buf, buildReport(), rollup.FormatText)
	if !strings.Contains(buf.String(), "image") {
		t.Error("expected drifted field 'image' in text output")
	}
}

func TestRender_JSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := rollup.Render(&buf, buildReport(), rollup.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out rollup.Report
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.TotalServices != 2 {
		t.Errorf("expected 2 services in JSON, got %d", out.TotalServices)
	}
}

func TestRender_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	err := rollup.Render(&buf, buildReport(), rollup.Format("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestRender_Text_NoDrift_ShowsOK(t *testing.T) {
	report := rollup.Build([]drift.Summary{
		{ServiceName: "svc", Diffs: nil},
	})
	var buf bytes.Buffer
	rollup.Render(&buf, report, rollup.FormatText)
	if !strings.Contains(buf.String(), "OK") {
		t.Error("expected 'OK' for clean service")
	}
}
