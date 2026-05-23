package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/history"
)

func makeSummaries(drifted bool) []drift.Summary {
	return []drift.Summary{
		{
			ServiceName: "api",
			HasDrift:    drifted,
			Drifts:      nil,
		},
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/historytest"
	_, err := history.NewStore(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestStore_AppendAndList(t *testing.T) {
	dir := t.TempDir()
	store, err := history.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	if err := store.Append(makeSummaries(false)); err != nil {
		t.Fatalf("Append: %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := store.Append(makeSummaries(true)); err != nil {
		t.Fatalf("Append: %v", err)
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].DriftCount != 0 {
		t.Errorf("first entry drift count: want 0, got %d", entries[0].DriftCount)
	}
	if entries[1].DriftCount != 1 {
		t.Errorf("second entry drift count: want 1, got %d", entries[1].DriftCount)
	}
}

func TestStore_ListEmpty(t *testing.T) {
	dir := t.TempDir()
	store, err := history.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty list, got %d entries", len(entries))
	}
}

func TestStore_AppendCreatesFile(t *testing.T) {
	dir := t.TempDir()
	store, _ := history.NewStore(dir)
	_ = store.Append(makeSummaries(false))

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}
