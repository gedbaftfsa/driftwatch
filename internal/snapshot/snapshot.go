package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Snapshot holds a point-in-time record of drift detection results.
type Snapshot struct {
	Timestamp time.Time      `json:"timestamp"`
	Summary   drift.Summary  `json:"summary"`
}

// Store manages reading and writing snapshots to disk.
type Store struct {
	Dir string
}

// NewStore creates a Store that persists snapshots under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir %q: %w", dir, err)
	}
	return &Store{Dir: dir}, nil
}

// Save writes a snapshot to a timestamped JSON file.
func (s *Store) Save(snap Snapshot) (string, error) {
	filename := fmt.Sprintf("snapshot-%s.json", snap.Timestamp.UTC().Format("20060102T150405Z"))
	path := filepath.Join(s.Dir, filename)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("snapshot: write %q: %w", path, err)
	}
	return path, nil
}

// Load reads a snapshot from the given file path.
func (s *Store) Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read %q: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal %q: %w", path, err)
	}
	return snap, nil
}
