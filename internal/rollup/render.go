package rollup

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format controls the output format for rendering a rollup report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Render writes the rollup Report to w in the requested format.
func Render(w io.Writer, r Report, format Format) error {
	switch format {
	case FormatJSON:
		return renderJSON(w, r)
	case FormatText:
		return renderText(w, r)
	default:
		return fmt.Errorf("rollup: unknown format %q", format)
	}
}

func renderJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func renderText(w io.Writer, r Report) error {
	fmt.Fprintf(w, "Rollup Report — %s\n", r.GeneratedAt.Format("2006-01-02T15:04:05Z"))
	fmt.Fprintf(w, "Services: %d total, %d drifted, severity: %s\n",
		r.TotalServices, r.DriftedCount, strings.ToUpper(r.TopSeverity))
	fmt.Fprintln(w, strings.Repeat("-", 48))

	for _, s := range r.Services {
		if s.DriftCount == 0 {
			fmt.Fprintf(w, "  %-24s OK\n", s.ServiceName)
		} else {
			fields := strings.Join(s.DriftedFields, ", ")
			fmt.Fprintf(w, "  %-24s %-6s %d field(s): %s\n",
				s.ServiceName, strings.ToUpper(s.TopSeverity), s.DriftCount, fields)
		}
	}
	return nil
}
