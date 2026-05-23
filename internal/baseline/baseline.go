package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Baseline represents a saved snapshot of declared service configurations
// used as a reference point for drift comparison over time.
type Baseline struct {
	CreatedAt time.Time          `json:"created_at"`
	Version   string             `json:"version"`
	Services  map[string]Service `json:"services"`
}

// Service holds the baseline state for a single declared service.
type Service struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Environment map[string]string `json:"environment"`
	Replicas    int               `json:"replicas"`
}

// Store manages persistence of baselines to disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating it if necessary.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("baseline: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a baseline to disk under the given label.
func (s *Store) Save(label string, b *Baseline) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	path := filepath.Join(s.dir, label+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// Load reads a baseline from disk by label.
func (s *Store) Load(label string) (*Baseline, error) {
	path := filepath.Join(s.dir, label+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: %q not found", label)
		}
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// List returns all available baseline labels in the store.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("baseline: list dir: %w", err)
	}
	var labels []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			labels = append(labels, e.Name()[:len(e.Name())-5])
		}
	}
	return labels, nil
}
