package scheduler_test

import (
	"context"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/scheduler"
)

func silentLogger() *log.Logger {
	return log.New(os.Discard, "", 0)
}

func TestScheduler_RunsJobImmediately(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	s := scheduler.New(10*time.Second, job, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = s.Run(ctx)

	if atomic.LoadInt32(&count) < 1 {
		t.Error("expected job to run at least once immediately")
	}
}

func TestScheduler_RunsJobOnInterval(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	s := scheduler.New(20*time.Millisecond, job, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
	defer cancel()

	_ = s.Run(ctx)

	got := atomic.LoadInt32(&count)
	if got < 2 {
		t.Errorf("expected at least 2 job runs, got %d", got)
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	job := func(ctx context.Context) error { return nil }

	s := scheduler.New(5*time.Second, job, silentLogger())
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- s.Run(ctx)
	}()

	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("scheduler did not stop after context cancel")
	}
}

func TestScheduler_JobErrorDoesNotStop(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return errors.New("drift check failed")
	}

	s := scheduler.New(20*time.Millisecond, job, silentLogger())
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()

	_ = s.Run(ctx)

	if atomic.LoadInt32(&count) < 2 {
		t.Error("expected scheduler to continue running despite job errors")
	}
}
