// Package metrics provides lightweight in-process counters and gauges
// for tracking driftwatch operational statistics such as scan counts,
// drift detections, and alert notifications.
package metrics

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing integer counter.
type Counter struct {
	value uint64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() {
	atomic.AddUint64(&c.value, 1)
}

// Add increments the counter by n.
func (c *Counter) Add(n uint64) {
	atomic.AddUint64(&c.value, n)
}

// Value returns the current counter value.
func (c *Counter) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

// Registry holds a named set of counters and gauges.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		counters: make(map[string]*Counter),
	}
}

// Counter returns the named counter, creating it if it does not exist.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()

	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Snapshot returns a point-in-time copy of all counter values keyed by name.
func (r *Registry) Snapshot() map[string]uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snap := make(map[string]uint64, len(r.counters))
	for name, c := range r.counters {
		snap[name] = c.Value()
	}
	return snap
}

// WriteTo writes a human-readable summary of all counters to w.
// Counters are emitted in alphabetical order.
func (r *Registry) WriteTo(w io.Writer) (int64, error) {
	snap := r.Snapshot()

	names := make([]string, 0, len(snap))
	for name := range snap {
		names = append(names, name)
	}
	sort.Strings(names)

	var total int64
	for _, name := range names {
		n, err := fmt.Fprintf(w, "%-40s %d\n", name, snap[name])
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// Default is the package-level registry used by the driftwatch application.
var Default = NewRegistry()

// Well-known counter names used across driftwatch.
const (
	CounterScansTotal     = "driftwatch_scans_total"
	CounterDriftDetected  = "driftwatch_drift_detected_total"
	CounterAlertsWarn     = "driftwatch_alerts_warn_total"
	CounterAlertsError    = "driftwatch_alerts_error_total"
	CounterSnapshotsSaved = "driftwatch_snapshots_saved_total"
	CounterExportsTotal   = "driftwatch_exports_total"
)
