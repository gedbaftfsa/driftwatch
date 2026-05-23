package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/history"
)

func appendN(t *testing.T, store *history.Store, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		if err := store.Append(makeSummaries(false)); err != nil {
			t.Fatalf("Append %d: %v", i, err)
		}
		// Ensure unique filenames based on timestamp.
		time.Sleep(5 * time.Millisecond)
	}
}

func countFiles(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	return len(entries)
}

func TestPrune_RemovesOldestEntries(t *testing.T) {
	dir := t.TempDir()
	store, _ := history.NewStore(dir)
	appendN(t, store, 5)

	deleted, err := store.Prune(history.PruneOptions{MaxEntries: 3})
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if deleted != 2 {
		t.Errorf("expected 2 deleted, got %d", deleted)
	}
	if n := countFiles(t, dir); n != 3 {
		t.Errorf("expected 3 files remaining, got %d", n)
	}
}

func TestPrune_NoDeletionNeeded(t *testing.T) {
	dir := t.TempDir()
	store, _ := history.NewStore(dir)
	appendN(t, store, 2)

	deleted, err := store.Prune(history.PruneOptions{MaxEntries: 5})
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if deleted != 0 {
		t.Errorf("expected 0 deleted, got %d", deleted)
	}
	if n := countFiles(t, dir); n != 2 {
		t.Errorf("expected 2 files, got %d", n)
	}
}

func TestPrune_ZeroMaxEntries_NoOp(t *testing.T) {
	dir := t.TempDir()
	store, _ := history.NewStore(dir)
	appendN(t, store, 3)

	deleted, err := store.Prune(history.PruneOptions{MaxEntries: 0})
	if err != nil {
		t.Fatalf("Prune: %v", err)
	}
	if deleted != 0 {
		t.Errorf("expected 0 deleted, got %d", deleted)
	}
	if n := countFiles(t, dir); n != 3 {
		t.Errorf("expected 3 files, got %d", n)
	}
}
