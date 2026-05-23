package audit_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/audit"
	"github.com/yourusername/driftwatch/internal/drift"
)

func makeSummaries(count int, drifted bool) []drift.Summary {
	var s []drift.Summary
	for i := 0; i < count; i++ {
		s = append(s, drift.Summary{
			ServiceName: fmt.Sprintf("svc-%d", i),
			HasDrift:    drifted,
		})
	}
	return s
}

func TestNewLogger_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "audit_logs")
	_, err := audit.NewLogger(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestLog_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	logger, err := audit.NewLogger(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	summaries := makeSummaries(3, false)
	if err := logger.Log("run-001", summaries); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	expected := filepath.Join(dir, fmt.Sprintf("%s.jsonl", time.Now().UTC().Format("2006-01-02")))
	if _, err := os.Stat(expected); os.IsNotExist(err) {
		t.Fatalf("expected log file %q to exist", expected)
	}
}

func TestLog_ReadAll_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	logger, err := audit.NewLogger(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := logger.Log("run-001", makeSummaries(5, false)); err != nil {
		t.Fatalf("Log failed: %v", err)
	}
	if err := logger.Log("run-002", makeSummaries(3, true)); err != nil {
		t.Fatalf("Log failed: %v", err)
	}

	filename := fmt.Sprintf("%s.jsonl", time.Now().UTC().Format("2006-01-02"))
	entries, err := logger.ReadAll(filename)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].RunID != "run-001" {
		t.Errorf("expected RunID run-001, got %q", entries[0].RunID)
	}
	if entries[1].DriftedCount != 3 {
		t.Errorf("expected DriftedCount 3, got %d", entries[1].DriftedCount)
	}
	if entries[0].TotalChecked != 5 {
		t.Errorf("expected TotalChecked 5, got %d", entries[0].TotalChecked)
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	dir := t.TempDir()
	logger, err := audit.NewLogger(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := logger.ReadAll("nonexistent.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing file, got %v", entries)
	}
}
