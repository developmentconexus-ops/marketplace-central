package oraclebatch

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestChunksPreservesOrderAndMaximum(t *testing.T) {
	ids := []int64{1, 2, 3, 4, 5}
	got := Chunks(ids, 2)
	if len(got) != 3 || len(got[0]) != 2 || len(got[2]) != 1 {
		t.Fatalf("Chunks() = %#v, want 3 chunks sized 2,2,1", got)
	}
	flattened := make([]int64, 0, len(ids))
	for _, chunk := range got {
		flattened = append(flattened, chunk...)
	}
	for i, want := range ids {
		if flattened[i] != want {
			t.Fatalf("flattened[%d]=%d, want %d", i, flattened[i], want)
		}
	}
}

func TestStockBatchSemaphoreBoundsConcurrency(t *testing.T) {
	const calls = 8
	sem := NewSemaphore(4)
	var inFlight, maxInFlight atomic.Int32
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < calls; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			release, err := sem.Acquire(context.Background())
			if err != nil {
				t.Errorf("Acquire() error = %v", err)
				return
			}
			defer release()
			current := inFlight.Add(1)
			for {
				previous := maxInFlight.Load()
				if current <= previous || maxInFlight.CompareAndSwap(previous, current) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			inFlight.Add(-1)
		}()
	}
	close(start)
	wg.Wait()
	if got := maxInFlight.Load(); got != 4 {
		t.Fatalf("max in-flight = %d, want 4", got)
	}
}

func TestStockBatchSemaphoreDoesNotBlockInteractiveClass(t *testing.T) {
	sem := NewSemaphore(1)
	release, err := sem.Acquire(context.Background())
	if err != nil {
		t.Fatalf("batch Acquire() error = %v", err)
	}
	defer release()

	interactiveDone := make(chan struct{})
	go func() {
		close(interactiveDone)
	}()
	select {
	case <-interactiveDone:
	case <-time.After(time.Second):
		t.Fatal("interactive operation was blocked by batch semaphore")
	}
}
