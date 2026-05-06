package schedule_test

import (
	"context"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/schedule"
)

type mockRunner struct {
	calls atomic.Int32
	err   error
}

func (m *mockRunner) Run(_ context.Context) error {
	m.calls.Add(1)
	return m.err
}

func silentLogger() *log.Logger {
	return log.New(os.Discard, "", 0)
}

func TestNew_NotNil(t *testing.T) {
	s := schedule.New(time.Second, &mockRunner{}, silentLogger())
	if s == nil {
		t.Fatal("expected non-nil Scheduler")
	}
}

func TestNew_NilLogger_UsesDefault(t *testing.T) {
	// Should not panic when logger is nil.
	s := schedule.New(time.Second, &mockRunner{}, nil)
	if s == nil {
		t.Fatal("expected non-nil Scheduler")
	}
}

func TestStart_RunsImmediately(t *testing.T) {
	runner := &mockRunner{}
	s := schedule.New(10*time.Second, runner, silentLogger())

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- s.Start(ctx)
	}()

	// Give the scheduler time to fire the initial run.
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	if runner.calls.Load() < 1 {
		t.Errorf("expected at least 1 call, got %d", runner.calls.Load())
	}
}

func TestStart_TicksRepeatedly(t *testing.T) {
	runner := &mockRunner{}
	s := schedule.New(30*time.Millisecond, runner, silentLogger())

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- s.Start(ctx)
	}()

	time.Sleep(120 * time.Millisecond)
	cancel()
	<-done

	if runner.calls.Load() < 2 {
		t.Errorf("expected at least 2 calls, got %d", runner.calls.Load())
	}
}

func TestStart_ContinuesOnRunnerError(t *testing.T) {
	runner := &mockRunner{err: errors.New("boom")}
	s := schedule.New(30*time.Millisecond, runner, silentLogger())

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- s.Start(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	<-done

	if runner.calls.Load() < 1 {
		t.Errorf("expected calls despite error, got %d", runner.calls.Load())
	}
}
