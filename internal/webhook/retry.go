package webhook

import (
	"context"
	"fmt"
	"time"

	"github.com/driftwatch/internal/drift"
)

// RetryConfig controls retry behaviour for webhook delivery.
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
}

// DefaultRetryConfig returns sensible retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       2 * time.Second,
	}
}

// RetrySender wraps a Sender and retries on failure.
type RetrySender struct {
	inner  *Sender
	cfg    RetryConfig
	sleep  func(time.Duration) // injectable for tests
}

// NewRetrySender creates a RetrySender with the given config.
func NewRetrySender(inner *Sender, cfg RetryConfig) *RetrySender {
	return &RetrySender{
		inner: inner,
		cfg:   cfg,
		sleep: time.Sleep,
	}
}

// Send attempts delivery up to MaxAttempts times, sleeping Delay between tries.
func (r *RetrySender) Send(ctx context.Context, summaries []drift.Summary) error {
	var lastErr error
	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("webhook retry: context cancelled: %w", err)
		}
		lastErr = r.inner.Send(ctx, summaries)
		if lastErr == nil {
			return nil
		}
		if attempt < r.cfg.MaxAttempts {
			r.sleep(r.cfg.Delay)
		}
	}
	return fmt.Errorf("webhook retry: all %d attempts failed: %w", r.cfg.MaxAttempts, lastErr)
}
