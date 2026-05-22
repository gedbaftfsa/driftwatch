package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/snapshot"
)

func makeSummary(drifted int) drift.Summary {
	return drift.Summary{
		Total:   drifted + 2,
		Drifted: drifted,
		Clean:   2,
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	snap := snapshot.Snapshot{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Summary:   makeSummary(3),
	}

	path, err := store.Save(snap)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Summary.Drifted != snap.Summary.Drifted {
		t.Errorf("Drifted: got %d, want %d", loaded.Summary.Drifted, snap.Summary.Drifted)
	}
	if !loaded.Timestamp.Equal(snap.Timestamp) {
		t.Errorf("Timestamp: got %v, want %v", loaded.Timestamp, snap.Timestamp)
	}
}

func TestStore_SaveCreatesFile(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.NewStore(dir)

	snap := snapshot.Snapshot{
		Timestamp: time.Now().UTC(),
		Summary:   makeSummary(0),
	}
	path, err := store.Save(snap)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file %q to exist: %v", path, err)
	}
}

func TestStore_LoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.NewStore(dir)

	_, err := store.Load(filepath.Join(dir, "nonexistent.json"))
	if err == nil {
		t.Error("expected error loading missing file, got nil")
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "snapshots")
	_, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir %q to be created: %v", dir, err)
	}
}
