package webhook

import "time"

// SetSleep allows tests to inject a no-op sleep function into RetrySender
// without exposing the field in the public API.
func (r *RetrySender) SetSleep(fn func(time.Duration)) {
	r.sleep = fn
}
