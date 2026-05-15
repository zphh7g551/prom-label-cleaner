package pipeline

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func newSchedulerTestServer(t *testing.T, calls *atomic.Int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		fmt.Fprintln(w, `# HELP up service up\n# TYPE up gauge\nup 1`)
	}))
}

func TestSchedulerRunsImmediately(t *testing.T) {
	var calls atomic.Int32
	srv := newSchedulerTestServer(t, &calls)
	defer srv.Close()

	cfg := DefaultConfig()
	cfg.TargetURL = srv.URL
	cfg.DryRun = true

	runner, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	sched := NewScheduler(10*time.Second, runner)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	sched.Run(ctx)

	if calls.Load() < 1 {
		t.Error("expected at least one run immediately")
	}
}

func TestSchedulerRepeats(t *testing.T) {
	var calls atomic.Int32
	srv := newSchedulerTestServer(t, &calls)
	defer srv.Close()

	cfg := DefaultConfig()
	cfg.TargetURL = srv.URL
	cfg.DryRun = true

	runner, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	sched := NewScheduler(50*time.Millisecond, runner)
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()

	sched.Run(ctx)

	if got := calls.Load(); got < 2 {
		t.Errorf("expected at least 2 runs, got %d", got)
	}
}

func TestSchedulerStopsOnContextCancel(t *testing.T) {
	var calls atomic.Int32
	srv := newSchedulerTestServer(t, &calls)
	defer srv.Close()

	cfg := DefaultConfig()
	cfg.TargetURL = srv.URL
	cfg.DryRun = true

	runner, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	sched := NewScheduler(20*time.Millisecond, runner)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		sched.Run(ctx)
		close(done)
	}()

	time.Sleep(60 * time.Millisecond)
	cancel()

	select {
	case <-done:
	// ok
	case <-time.After(500 * time.Millisecond):
		t.Error("scheduler did not stop after context cancellation")
	}
}
