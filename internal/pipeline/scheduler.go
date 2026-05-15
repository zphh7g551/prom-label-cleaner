package pipeline

import (
	"context"
	"log"
	"time"
)

// Scheduler periodically runs the pipeline at a configured interval.
type Scheduler struct {
	interval time.Duration
	runner   *Runner
}

// NewScheduler creates a Scheduler that will invoke the given Runner
// every interval duration.
func NewScheduler(interval time.Duration, runner *Runner) *Scheduler {
	return &Scheduler{
		interval: interval,
		runner:   runner,
	}
}

// Run starts the scheduling loop. It blocks until ctx is cancelled.
// The first run is executed immediately, then repeated every interval.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.tick(ctx)

	for {
		select {
		case <-ticker.C:
			s.tick(ctx)
		case <-ctx.Done():
			log.Println("scheduler: context cancelled, stopping")
			return
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	result, err := s.runner.Run(ctx)
	if err != nil {
		log.Printf("scheduler: run error: %v", err)
		return
	}
	log.Printf("scheduler: run complete — families=%d labelsRemoved=%d",
		result.FamilyCount, result.LabelsRemoved)
}
