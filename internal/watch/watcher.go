package watch

import (
	"context"
	"log"
	"time"

	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/provider"
)

// Config holds watcher configuration.
type Config struct {
	Interval    time.Duration
	ServiceNames []string
}

// DriftHandler is called whenever drift is detected.
type DriftHandler func(results []drift.Result)

// Watcher continuously polls the runtime provider and reports drift.
type Watcher struct {
	cfg      Config
	provider *provider.RuntimeProvider
	handler  DriftHandler
	logger   *log.Logger
}

// New creates a new Watcher.
func New(cfg Config, p *provider.RuntimeProvider, handler DriftHandler, logger *log.Logger) *Watcher {
	if logger == nil {
		logger = log.Default()
	}
	return &Watcher{
		cfg:      cfg,
		provider: p,
		handler:  handler,
		logger:   logger,
	}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context, declared []drift.ServiceSpec) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	w.poll(ctx, declared)

	for {
		select {
		case <-ticker.C:
			w.poll(ctx, declared)
		case <-ctx.Done():
			w.logger.Println("watcher: context cancelled, stopping")
			return ctx.Err()
		}
	}
}

func (w *Watcher) poll(ctx context.Context, declared []drift.ServiceSpec) {
	states, err := w.provider.FetchAll(ctx)
	if err != nil {
		w.logger.Printf("watcher: fetch error: %v", err)
		return
	}

	results := drift.Detect(declared, states)
	if len(results) > 0 {
		w.handler(results)
	}
}
