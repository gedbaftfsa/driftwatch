package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Entry represents a single drift check result stored in history.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Summaries []drift.Summary  `json:"summaries"`
	DriftCount int             `json:"drift_count"`
}

// Store persists drift history entries to disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir %s: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Append writes a new history entry with the current timestamp.
func (s *Store) Append(summaries []drift.Summary) error {
	entry := Entry{
		Timestamp:  time.Now().UTC(),
		Summaries:  summaries,
		DriftCount: countDrifted(summaries),
	}
	filename := entry.Timestamp.Format("20060102T150405Z") + ".json"
	path := filepath.Join(s.dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("history: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(entry)
}

// List returns all history entries sorted by timestamp ascending.
func (s *Store) List() ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(s.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	var entries []Entry
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, fmt.Errorf("history: read %s: %w", m, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: parse %s: %w", m, err)
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
	return entries, nil
}

func countDrifted(summaries []drift.Summary) int {
	n := 0
	for _, s := range summaries {
		if s.HasDrift {
			n++
		}
	}
	return n
}
