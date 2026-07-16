package stub

import (
	"context"
	"errors"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

func TestWriterProgramsResultsAndRecordsOrder(t *testing.T) {
	failure := &domain.Failure{Code: domain.FailureCodeProviderRateLimited, Retryable: true}
	w := NewWriter(map[string]Result{"b": {Failure: failure}, "c": {Err: errors.New("boom")}})
	for _, key := range []string{"a", "b", "c"} {
		out, err := w.Apply(context.Background(), ports.WriteItem{IdempotencyKey: key})
		if key == "a" && (err != nil || out.Failure != nil) {
			t.Fatalf("default=%+v err=%v", out, err)
		}
		if key == "b" && out.Failure != failure {
			t.Fatalf("failure=%+v", out.Failure)
		}
		if key == "c" && err == nil {
			t.Fatal("expected error")
		}
	}
	if got := w.Keys(); len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("keys=%v", got)
	}
}
