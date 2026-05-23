package scheduler

import (
	"context"
	"log"
	"time"
)

// Job is a function that runs on a schedule.
type Job func(ctx context.Context) error

// Scheduler runs a job at a fixed interval.
type Scheduler struct {
	interval time.Duration
	job      Job
	logger   *log.Logger
}

// New creates a new Scheduler with the given interval and job.
func New(interval time.Duration, job Job, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		interval: interval,
		job:      job,
		logger:   logger,
	}
}

// Run starts the scheduler loop. It runs the job immediately, then repeats
// at the configured interval until the context is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	s.logger.Printf("scheduler: starting with interval %s", s.interval)

	if err := s.runJob(ctx); err != nil {
		s.logger.Printf("scheduler: job error: %v", err)
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("scheduler: context cancelled, stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := s.runJob(ctx); err != nil {
				s.logger.Printf("scheduler: job error: %v", err)
			}
		}
	}
}

func (s *Scheduler) runJob(ctx context.Context) error {
	s.logger.Println("scheduler: running drift check job")
	start := time.Now()
	err := s.job(ctx)
	s.logger.Printf("scheduler: job completed in %s", time.Since(start))
	return err
}
