package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/divergences/domain"
)

// These are pure unit tests: the recorder is built over a nil pool, and
// every case here must return before ever touching it — proving the
// validation/skip short-circuits happen strictly before any DB call
// (round-trip/DB-touching behavior is proven separately by
// recorder_integration_test.go against a real hermetic Postgres).

func validInput() domain.EvaluateInput {
	expected := time.Date(2026, 7, 31, 9, 0, 12, 0, time.UTC)
	observed := time.Date(2026, 7, 31, 9, 3, 44, 0, time.UTC)
	return domain.EvaluateInput{
		TenantID:           "tenant1",
		Provider:           "mercado_livre",
		EntityType:         domain.EntityTypeListing,
		EntityID:           "MLB123",
		Kind:               domain.KindEstoque,
		Linked:             true,
		ExpectedValue:      "12",
		ObservedValue:      "9",
		ExpectedSource:     "mirror:sankhya:vendavel",
		ObservedSource:     "api_listing:available_quantity",
		ExpectedObservedAt: &expected,
		ObservedObservedAt: &observed,
	}
}

// IC-02 Error Cases + feature.md Negative Scenarios: Evaluate without
// expected_observed_at must be rejected by a named error, never default to
// now().
func TestEvaluateRejectsMissingExpectedObservedAt(t *testing.T) {
	rec := NewRecorder(nil)
	input := validInput()
	input.ExpectedObservedAt = nil

	_, err := rec.Evaluate(context.Background(), input)
	if !errors.Is(err, domain.ErrExpectedObservedAtRequired) {
		t.Fatalf("expected ErrExpectedObservedAtRequired, got %v", err)
	}
}

func TestEvaluateRejectsMissingObservedObservedAt(t *testing.T) {
	rec := NewRecorder(nil)
	input := validInput()
	input.ObservedObservedAt = nil

	_, err := rec.Evaluate(context.Background(), input)
	if !errors.Is(err, domain.ErrObservedObservedAtRequired) {
		t.Fatalf("expected ErrObservedObservedAtRequired, got %v", err)
	}
}

// IC-02 Persistence Expectations: no linkage → skip, not an error, not an
// evaluation against a fabricated "expected".
func TestEvaluateSkipsWithoutLinkage(t *testing.T) {
	rec := NewRecorder(nil)
	input := validInput()
	input.Linked = false

	outcome, err := rec.Evaluate(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error for a skip, got %v", err)
	}
	if !outcome.Skipped {
		t.Fatal("expected Skipped=true when Linked=false")
	}
	if outcome.Opened || outcome.Updated || outcome.Resolved {
		t.Fatalf("expected a skip to touch nothing, got %+v", outcome)
	}
}

// A skip must win over a missing-timestamp rejection: if the entity has no
// linkage at all, Evaluate must not even reach the timestamp validation —
// there is nothing to evaluate.
func TestEvaluateSkipTakesPriorityOverMissingTimestamps(t *testing.T) {
	rec := NewRecorder(nil)
	input := validInput()
	input.Linked = false
	input.ExpectedObservedAt = nil

	outcome, err := rec.Evaluate(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !outcome.Skipped {
		t.Fatal("expected Skipped=true")
	}
}
