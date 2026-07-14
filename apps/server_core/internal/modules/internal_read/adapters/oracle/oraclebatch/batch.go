package oraclebatch

import (
	"context"
	"sync"
)

// Chunks splits ids into ordered chunks no larger than max.
func Chunks(ids []int64, max int) [][]int64 {
	if len(ids) == 0 || max <= 0 {
		return nil
	}
	chunks := make([][]int64, 0, (len(ids)+max-1)/max)
	for start := 0; start < len(ids); start += max {
		end := start + max
		if end > len(ids) {
			end = len(ids)
		}
		chunk := make([]int64, end-start)
		copy(chunk, ids[start:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}

// Semaphore bounds batch-class Oracle concurrency.
type Semaphore struct {
	permits chan struct{}
}

func NewSemaphore(permits int) *Semaphore {
	if permits <= 0 {
		permits = 1
	}
	return &Semaphore{permits: make(chan struct{}, permits)}
}

// Acquire waits for a batch permit or context cancellation.
func (s *Semaphore) Acquire(ctx context.Context) (release func(), err error) {
	if s == nil || s.permits == nil {
		return nil, context.Canceled
	}
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case s.permits <- struct{}{}:
		var once sync.Once
		return func() {
			once.Do(func() { <-s.permits })
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
