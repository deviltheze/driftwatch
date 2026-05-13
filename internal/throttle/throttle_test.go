package throttle_test

import (
	"sync"
	"testing"
	"time"

	"github.com/yourusername/driftwatch/internal/throttle"
)

func TestNew_DefaultsMaxConns(t *testing.T) {
	th := throttle.New(0, 0)
	if th == nil {
		t.Fatal("expected non-nil Throttle")
	}
}

func TestAcquire_EmptyHost_ReturnsError(t *testing.T) {
	th := throttle.New(2, 0)
	if err := th.Acquire(""); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestAcquire_Release_SingleSlot(t *testing.T) {
	th := throttle.New(1, 0)
	host := "server-01"

	if err := th.Acquire(host); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := th.Active(host); got != 1 {
		t.Fatalf("expected 1 active, got %d", got)
	}
	th.Release(host)
	if got := th.Active(host); got != 0 {
		t.Fatalf("expected 0 active after release, got %d", got)
	}
}

func TestActive_UnknownHost_ReturnsZero(t *testing.T) {
	th := throttle.New(3, 0)
	if got := th.Active("ghost"); got != 0 {
		t.Fatalf("expected 0 for unknown host, got %d", got)
	}
}

func TestAcquire_ConcurrentRespectsCap(t *testing.T) {
	const cap = 3
	th := throttle.New(cap, 0)
	host := "server-02"

	var wg sync.WaitGroup
	peak := make(chan int, 20)

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = th.Acquire(host)
			peak <- th.Active(host)
			time.Sleep(10 * time.Millisecond)
			th.Release(host)
		}()
	}
	wg.Wait()
	close(peak)

	for v := range peak {
		if v > cap {
			t.Fatalf("active count %d exceeded cap %d", v, cap)
		}
	}
}

func TestRelease_NoAcquire_NoOp(t *testing.T) {
	th := throttle.New(2, 0)
	// should not panic
	th.Release("phantom")
}
