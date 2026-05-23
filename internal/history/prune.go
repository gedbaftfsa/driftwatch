package history

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// PruneOptions controls how old history entries are removed.
type PruneOptions struct {
	// MaxEntries is the maximum number of history files to retain.
	// If zero, no pruning by count is performed.
	MaxEntries int
}

// Prune removes the oldest history entries exceeding the configured limits.
// It returns the number of files deleted.
func (s *Store) Prune(opts PruneOptions) (int, error) {
	if opts.MaxEntries <= 0 {
		return 0, nil
	}

	matches, err := filepath.Glob(filepath.Join(s.dir, "*.json"))
	if err != nil {
		return 0, fmt.Errorf("history prune: glob: %w", err)
	}

	// Sort ascending so oldest entries come first.
	sort.Strings(matches)

	toDelete := len(matches) - opts.MaxEntries
	if toDelete <= 0 {
		return 0, nil
	}

	deleted := 0
	for _, f := range matches[:toDelete] {
		if err := os.Remove(f); err != nil {
			return deleted, fmt.Errorf("history prune: remove %s: %w", f, err)
		}
		deleted++
	}
	return deleted, nil
}
