package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Format represents the output format for exported drift summaries.
type Format string

const (
	FormatCSV  Format = "csv"
	FormatJSON Format = "json"
)

// Record is a flat representation of a drift summary suitable for export.
type Record struct {
	Timestamp   string `json:"timestamp"`
	Service     string `json:"service"`
	Drifted     bool   `json:"drifted"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

// Exporter writes drift summaries to an io.Writer in the requested format.
type Exporter struct {
	format Format
	writer io.Writer
}

// NewExporter creates a new Exporter for the given format and writer.
func NewExporter(format Format, w io.Writer) (*Exporter, error) {
	if format != FormatCSV && format != FormatJSON {
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
	return &Exporter{format: format, writer: w}, nil
}

// Export converts drift summaries into records and writes them.
func (e *Exporter) Export(summaries []drift.Summary) error {
	records := toRecords(summaries)
	switch e.format {
	case FormatCSV:
		return writeCSV(e.writer, records)
	case FormatJSON:
		return writeJSON(e.writer, records)
	}
	return nil
}

func toRecords(summaries []drift.Summary) []Record {
	ts := time.Now().UTC().Format(time.RFC3339)
	records := make([]Record, 0, len(summaries))
	for _, s := range summaries {
		rec := Record{
			Timestamp:   ts,
			Service:     s.ServiceName,
			Drifted:     s.Drifted,
			Severity:    string(s.Severity),
			Description: s.Description,
		}
		records = append(records, rec)
	}
	return records
}

func writeCSV(w io.Writer, records []Record) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"timestamp", "service", "drifted", "severity", "description"}); err != nil {
		return fmt.Errorf("csv header: %w", err)
	}
	for _, r := range records {
		row := []string{r.Timestamp, r.Service, fmt.Sprintf("%t", r.Drifted), r.Severity, r.Description}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeJSON(w io.Writer, records []Record) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(records); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}
