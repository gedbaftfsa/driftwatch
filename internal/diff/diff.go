package diff

import (
	"fmt"
	"strings"
)

// FieldDiff represents a single field-level difference between declared and live state.
type FieldDiff struct {
	Field    string
	Declared string
	Live     string
}

// String returns a human-readable representation of the field diff.
func (f FieldDiff) String() string {
	return fmt.Sprintf("%s: declared=%q live=%q", f.Field, f.Declared, f.Live)
}

// ServiceDiff holds all field-level diffs for a single service.
type ServiceDiff struct {
	ServiceName string
	Fields      []FieldDiff
}

// HasDiff returns true if there are any field differences.
func (s ServiceDiff) HasDiff() bool {
	return len(s.Fields) > 0
}

// Summary returns a compact multi-line string of all field diffs.
func (s ServiceDiff) Summary() string {
	if !s.HasDiff() {
		return fmt.Sprintf("%s: no diff", s.ServiceName)
	}
	lines := make([]string, 0, len(s.Fields)+1)
	lines = append(lines, fmt.Sprintf("%s:", s.ServiceName))
	for _, f := range s.Fields {
		lines = append(lines, "  "+f.String())
	}
	return strings.Join(lines, "\n")
}

// CompareFields compares two maps of string fields and returns a list of FieldDiff.
// Keys present in declared but missing from live are reported with live value "<missing>".
// Keys present in live but missing from declared are reported with declared value "<untracked>".
func CompareFields(declared, live map[string]string) []FieldDiff {
	var diffs []FieldDiff

	for k, dv := range declared {
		lv, ok := live[k]
		if !ok {
			diffs = append(diffs, FieldDiff{Field: k, Declared: dv, Live: "<missing>"})
		} else if dv != lv {
			diffs = append(diffs, FieldDiff{Field: k, Declared: dv, Live: lv})
		}
	}

	for k, lv := range live {
		if _, ok := declared[k]; !ok {
			diffs = append(diffs, FieldDiff{Field: k, Declared: "<untracked>", Live: lv})
		}
	}

	return diffs
}
