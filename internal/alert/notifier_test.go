package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftwatch/internal/alert"
	"github.com/driftwatch/internal/drift"
)

func makeSummaries(names []string, diffs [][]string) []drift.Summary {
	var out []drift.Summary
	for i, name := range names {
		out = append(out, drift.Summary{ServiceName: name, Diffs: diffs[i]})
	}
	return out
}

func TestNotify_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(alert.Config{Output: &buf})
	summaries := makeSummaries([]string{"svc-a"}, [][]string{{}})
	emitted := n.Notify(summaries)
	if emitted {
		t.Error("expected no alerts for zero diffs")
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %s", buf.String())
	}
}

func TestNotify_WarnDrift(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(alert.Config{Threshold: alert.LevelWarn, Output: &buf})
	summaries := makeSummaries([]string{"svc-b"}, [][]string{{"env PORT: want 8080 got 9090"}})
	emitted := n.Notify(summaries)
	if !emitted {
		t.Error("expected alert to be emitted")
	}
	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected WARN level, got: %s", out)
	}
	if !strings.Contains(out, "svc-b") {
		t.Errorf("expected service name in output, got: %s", out)
	}
}

func TestNotify_ErrorDrift(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(alert.Config{Threshold: alert.LevelWarn, Output: &buf})
	summaries := makeSummaries([]string{"svc-c"}, [][]string{{"image mismatch: want nginx:1.25 got nginx:1.21"}})
	emitted := n.Notify(summaries)
	if !emitted {
		t.Error("expected alert to be emitted")
	}
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Errorf("expected ERROR level, got: %s", buf.String())
	}
}

func TestNotify_ErrorThreshold_SuppressesWarn(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewNotifier(alert.Config{Threshold: alert.LevelError, Output: &buf})
	summaries := makeSummaries([]string{"svc-d"}, [][]string{{"env PORT: want 8080 got 9090"}})
	emitted := n.Notify(summaries)
	if emitted {
		t.Error("expected warn-level drift to be suppressed at error threshold")
	}
}

func TestNotify_DefaultOutput(t *testing.T) {
	// Should not panic when Output is nil (falls back to os.Stderr)
	n := alert.NewNotifier(alert.Config{})
	summaries := makeSummaries([]string{"svc-e"}, [][]string{{}})
	n.Notify(summaries) // no assertion, just ensure no panic
}
