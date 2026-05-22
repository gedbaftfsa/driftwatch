package alert

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/driftwatch/internal/drift"
)

// Level represents the severity threshold for alerts.
type Level string

const (
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Config holds alert notifier configuration.
type Config struct {
	Threshold Level
	Output    io.Writer
}

// Notifier evaluates drift summaries and emits alerts.
type Notifier struct {
	cfg Config
}

// NewNotifier creates a Notifier with the given config.
// If Output is nil, os.Stderr is used.
func NewNotifier(cfg Config) *Notifier {
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}
	if cfg.Threshold == "" {
		cfg.Threshold = LevelWarn
	}
	return &Notifier{cfg: cfg}
}

// Notify writes alert lines for each drifted service that meets the threshold.
// Returns true if any alerts were emitted.
func (n *Notifier) Notify(summaries []drift.Summary) bool {
	var emitted bool
	for _, s := range summaries {
		if !n.shouldAlert(s) {
			continue
		}
		level := n.levelFor(s)
		fmt.Fprintf(n.cfg.Output, "[%s] service %q has drift: %s\n",
			strings.ToUpper(string(level)), s.ServiceName, strings.Join(s.Diffs, "; "))
		emitted = true
	}
	return emitted
}

// shouldAlert returns true when the summary has diffs and meets the threshold.
func (n *Notifier) shouldAlert(s drift.Summary) bool {
	if len(s.Diffs) == 0 {
		return false
	}
	level := n.levelFor(s)
	if n.cfg.Threshold == LevelError && level != LevelError {
		return false
	}
	return true
}

// levelFor assigns a severity level based on the nature of diffs.
func (n *Notifier) levelFor(s drift.Summary) Level {
	for _, d := range s.Diffs {
		if strings.Contains(d, "missing") || strings.Contains(d, "image") {
			return LevelError
		}
	}
	return LevelWarn
}
