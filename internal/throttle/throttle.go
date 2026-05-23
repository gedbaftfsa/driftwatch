package throttle

import (
	"sync"
	"time"
)

// Throttle limits how frequently an action can be triggered per key.
type Throttle struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
	cooldown time.Duration
}

// New creates a Throttle with the given cooldown duration.
func New(cooldown time.Duration) *Throttle {
	return &Throttle{
		lastSeen: make(map[string]time.Time),
		cooldown: cooldown,
	}
}

// Allow returns true if the key has not been seen within the cooldown window.
// If allowed, it records the current time for the key.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if last, ok := t.lastSeen[key]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.lastSeen[key] = now
	return true
}

// Reset clears the recorded time for a specific key, allowing it immediately.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSeen, key)
}

// ResetAll clears all recorded keys.
func (t *Throttle) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSeen = make(map[string]time.Time)
}

// ActiveKeys returns the list of keys currently tracked.
func (t *Throttle) ActiveKeys() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	keys := make([]string, 0, len(t.lastSeen))
	for k := range t.lastSeen {
		keys = append(keys, k)
	}
	return keys
}
