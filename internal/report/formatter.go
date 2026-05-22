package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Format controls the output format of the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes drift detection results to an io.Writer.
type Formatter struct {
	format Format
	writer io.Writer
}

// NewFormatter creates a new Formatter with the given format and writer.
func NewFormatter(format Format, w io.Writer) *Formatter {
	return &Formatter{format: format, writer: w}
}

// Write outputs the drift summary using the configured format.
func (f *Formatter) Write(summary drift.Summary) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(summary)
	default:
		return f.writeText(summary)
	}
}

func (f *Formatter) writeText(s drift.Summary) error {
	fmt.Fprintf(f.writer, "Drift Report\n%s\n", strings.Repeat("=", 40))
	fmt.Fprintf(f.writer, "Services checked : %d\n", s.Total)
	fmt.Fprintf(f.writer, "Drifted          : %d\n", s.Drifted)
	fmt.Fprintf(f.writer, "Missing          : %d\n", s.Missing)

	if len(s.Details) == 0 {
		fmt.Fprintln(f.writer, "\nNo drift detected. ✓")
		return nil
	}

	fmt.Fprintln(f.writer, "\nDrift Details:")
	for _, d := range s.Details {
		fmt.Fprintf(f.writer, "  [%s] %s\n", d.Service, d.Field)
		fmt.Fprintf(f.writer, "    declared : %s\n", d.Declared)
		fmt.Fprintf(f.writer, "    live     : %s\n", d.Live)
	}
	return nil
}

func (f *Formatter) writeJSON(s drift.Summary) error {
	details := "[]"
	if len(s.Details) > 0 {
		var parts []string
		for _, d := range s.Details {
			parts = append(parts, fmt.Sprintf(
				`{"service":%q,"field":%q,"declared":%q,"live":%q}`,
				d.Service, d.Field, d.Declared, d.Live,
			))
		}
		details = "[" + strings.Join(parts, ",") + "]"
	}
	_, err := fmt.Fprintf(f.writer,
		`{"total":%d,"drifted":%d,"missing":%d,"details":%s}\n`,
		s.Total, s.Drifted, s.Missing, details,
	)
	return err
}
