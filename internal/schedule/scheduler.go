package schedule

import (
	"context"
	"log"
	"time"
)

// Runner is anything that can execute a drift check cycle.
type Runner interface {
	Run(ctx context.Context) error
}

// Scheduler periodically triggers a Runner at a fixed interval.
type Scheduler struct {
	interval time.Duration
	runner   Runner
	logger   *log.Logger
}

// New creates a Scheduler with the given interval and runner.
func New(interval time.Duration, runner Runner, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		interval: interval,
		runner:   runner,
		logger:   logger,
	}
}

// Start runs the scheduler loop, blocking until ctx is cancelled.
// It executes the runner immediately, then on every tick.
func (s *Scheduler) Start(ctx context.Context) error {
	s.logger.Printf("scheduler: starting with interval %s", s.interval)

	if err := s.tick(ctx); err != nil {
		s.logger.Printf("scheduler: initial run error: %v", err)
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("scheduler: shutting down")
			return ctx.Err()
		case <-ticker.C:
			if err := s.tick(ctx); err != nil {
				s.logger.Printf("scheduler: run error: %v", err)
			}
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) error {
	s.logger.Println("scheduler: running drift check")
	return s.runner.Run(ctx)
}
