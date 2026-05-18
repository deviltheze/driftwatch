package watchdog_test

import (
	"context"
	"errors"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/watchdog"
)

// silentLogger suppresses log output during tests.
func silentLogger() *log.Logger {
	return log.New(os.Discard, "", 0)
}

type healthyPinger struct{}

func (h *healthyPinger) Ping() error { return nil }

type failingPinger struct{}

func (f *failingPinger) Ping() error { return errors.New("component stalled") }

func TestNew_NotNil(t *testing.T) {
	wd := watchdog.New(0, nil)
	if wd == nil {
		t.Fatal("expected non-nil Watchdog")
	}
}

func TestNew_DefaultHeartbeat(t *testing.T) {
	// Passing zero should not panic and the watchdog should be usable.
	wd := watchdog.New(0, silentLogger())
	if wd == nil {
		t.Fatal("expected non-nil Watchdog")
	}
}

func TestNew_NilLogger_UsesDefault(t *testing.T) {
	// Should not panic when logger is nil.
	wd := watchdog.New(10*time.Millisecond, nil)
	if wd == nil {
		t.Fatal("expected non-nil Watchdog")
	}
}

func TestOnFailure_CalledForFailingPinger(t *testing.T) {
	wd := watchdog.New(20*time.Millisecond, silentLogger())
	wd.Register("bad", &failingPinger{})

	var called atomic.Int32
	wd.OnFailure(func(name string, err error) {
		if name == "bad" {
			called.Add(1)
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	wd.Start(ctx)

	if called.Load() == 0 {
		t.Error("expected OnFailure to be called at least once")
	}
}

func TestOnFailure_NotCalledForHealthyPinger(t *testing.T) {
	wd := watchdog.New(20*time.Millisecond, silentLogger())
	wd.Register("ok", &healthyPinger{})

	var called atomic.Int32
	wd.OnFailure(func(name string, err error) {
		called.Add(1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	wd.Start(ctx)

	if called.Load() != 0 {
		t.Errorf("expected OnFailure not to be called, got %d calls", called.Load())
	}
}

func TestStart_StopsOnContextCancel(t *testing.T) {
	wd := watchdog.New(50*time.Millisecond, silentLogger())
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		wd.Start(ctx)
		close(done)
	}()

	cancel()
	select {
	case <-done:
		// success
	case <-time.After(200 * time.Millisecond):
		t.Error("Start did not return after context cancellation")
	}
}
