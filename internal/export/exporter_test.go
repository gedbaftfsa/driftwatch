package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
)

func makeSummaries() []drift.Summary {
	return []drift.Summary{
		{ServiceName: "api", Drifted: true, Severity: drift.SeverityError, Description: "image mismatch"},
		{ServiceName: "worker", Drifted: false, Severity: drift.SeverityNone, Description: ""},
	}
}

func TestNewExporter_InvalidFormat(t *testing.T) {
	_, err := NewExporter("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestNewExporter_ValidFormats(t *testing.T) {
	for _, f := range []Format{FormatCSV, FormatJSON} {
		_, err := NewExporter(f, &bytes.Buffer{})
		if err != nil {
			t.Fatalf("unexpected error for format %s: %v", f, err)
		}
	}
}

func TestExport_CSV_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := NewExporter(FormatCSV, &buf)
	if err := ex.Export(makeSummaries()); err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "timestamp,service,drifted,severity,description") {
		t.Errorf("CSV missing header, got:\n%s", out)
	}
}

func TestExport_CSV_ContainsRows(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := NewExporter(FormatCSV, &buf)
	_ = ex.Export(makeSummaries())
	out := buf.String()
	if !strings.Contains(out, "api") || !strings.Contains(out, "worker") {
		t.Errorf("CSV missing service rows, got:\n%s", out)
	}
	if !strings.Contains(out, "image mismatch") {
		t.Errorf("CSV missing description, got:\n%s", out)
	}
}

func TestExport_JSON_ValidStructure(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := NewExporter(FormatJSON, &buf)
	if err := ex.Export(makeSummaries()); err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	var records []Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
}

func TestExport_JSON_FieldValues(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := NewExporter(FormatJSON, &buf)
	_ = ex.Export(makeSummaries())
	var records []Record
	_ = json.Unmarshal(buf.Bytes(), &records)
	if records[0].Service != "api" || !records[0].Drifted {
		t.Errorf("unexpected first record: %+v", records[0])
	}
	if records[1].Service != "worker" || records[1].Drifted {
		t.Errorf("unexpected second record: %+v", records[1])
	}
}

func TestExport_EmptySummaries(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := NewExporter(FormatJSON, &buf)
	if err := ex.Export(nil); err != nil {
		t.Fatalf("Export with nil summaries failed: %v", err)
	}
	var records []Record
	_ = json.Unmarshal(buf.Bytes(), &records)
	if len(records) != 0 {
		t.Errorf("expected empty records, got %d", len(records))
	}
}
