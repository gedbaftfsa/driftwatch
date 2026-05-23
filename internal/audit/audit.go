package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/driftwatch/internal/drift"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp  time.Time        `json:"timestamp"`
	RunID      string           `json:"run_id"`
	TotalChecked int            `json:"total_checked"`
	DriftedCount int            `json:"drifted_count"`
	Summaries  []drift.Summary  `json:"summaries"`
}

// Logger writes audit entries to a JSONL file.
type Logger struct {
	dir string
}

// NewLogger creates a new Logger, ensuring the directory exists.
func NewLogger(dir string) (*Logger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("audit: create dir %q: %w", dir, err)
	}
	return &Logger{dir: dir}, nil
}

// Log appends an audit entry to the current day's log file.
func (l *Logger) Log(runID string, summaries []drift.Summary) error {
	drifted := 0
	for _, s := range summaries {
		if s.HasDrift {
			drifted++
		}
	}

	entry := Entry{
		Timestamp:    time.Now().UTC(),
		RunID:        runID,
		TotalChecked: len(summaries),
		DriftedCount: drifted,
		Summaries:    summaries,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}

	filename := filepath.Join(l.dir, fmt.Sprintf("%s.jsonl", time.Now().UTC().Format("2006-01-02")))
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("audit: open file %q: %w", filename, err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", data)
	if err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}

// ReadAll reads all audit entries from a given log file.
func (l *Logger) ReadAll(filename string) ([]Entry, error) {
	path := filepath.Join(l.dir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: read file %q: %w", path, err)
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: unmarshal line: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
