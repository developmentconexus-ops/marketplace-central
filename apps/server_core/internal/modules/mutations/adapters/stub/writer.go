package stub

import (
	"context"
	"sync"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type Result struct {
	Failure *domain.Failure
	Err     error
}

type Writer struct {
	mu      sync.Mutex
	results map[string]Result
	keys    []string
}

func NewWriter(results map[string]Result) *Writer { return &Writer{results: results} }

func (w *Writer) Apply(_ context.Context, item ports.WriteItem) (ports.WriteOutcome, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.keys = append(w.keys, item.IdempotencyKey)
	result := w.results[item.IdempotencyKey]
	return ports.WriteOutcome{Failure: result.Failure}, result.Err
}

func (w *Writer) Keys() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string(nil), w.keys...)
}
