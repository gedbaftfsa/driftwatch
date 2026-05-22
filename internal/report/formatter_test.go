package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/report"
)

func makeSummary(drifted bool) drift.Summary {
	s := drift.Summary{Total: 2, Drifted: 0, Missing: 0}
	if drifted {
		s.Drifted = 1
		s.Details = []drift.Detail{
			{Service: "api", Field: "image", Declared: "app:1.0", Live: "app:1.1"},
		}
	}
	return s
}

func TestFormatter_Text_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatText, &buf)
	if err := f.Write(makeSummary(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got:\n%s", out)
	}
}

func TestFormatter_Text_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatText, &buf)
	if err := f.Write(makeSummary(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"api", "image", "app:1.0", "app:1.1", "Drifted", "1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestFormatter_JSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatJSON, &buf)
	if err := f.Write(makeSummary(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"drifted":0`) {
		t.Errorf("expected drifted:0 in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"details":[]`) {
		t.Errorf("expected empty details array in JSON, got: %s", out)
	}
}

func TestFormatter_JSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := report.NewFormatter(report.FormatJSON, &buf)
	if err := f.Write(makeSummary(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{`"service":"api"`, `"field":"image"`, `"declared":"app:1.0"`, `"live":"app:1.1"`} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output, got: %s", want, out)
		}
	}
}

func TestFormatter_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	// empty Format string should fall back to text
	f := report.NewFormatter("", &buf)
	if err := f.Write(makeSummary(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Drift Report") {
		t.Error("expected text format fallback")
	}
}
