package internalread

import (
	"context"
	"testing"
	"time"

	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

func batchReader(t *testing.T, repo *fakeRepo) *BatchStockReader {
	t.Helper()
	return NewBatchStockReader(NewReader(repo, "tenant-a", WithClock(func() time.Time { return fetchedAt })))
}

func TestBatchStockAnswersEveryKnownProductFromTheMirror(t *testing.T) {
	repo := &fakeRepo{rows: mirrorRows()}
	ctx := WithActiveSource(context.Background(), erpdomain.SourceXLSX)

	facts, err := batchReader(t, repo).GetStockFactsByIDs(ctx, []int64{1, 2, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(facts) != 2 {
		t.Fatalf("expected 2 facts, got %d", len(facts))
	}
	first := facts[1]
	if first == nil || first.Quantity == nil || *first.Quantity != 7 {
		t.Fatalf("product 1 stock: %+v", first)
	}
	if first.Source.System != string(erpdomain.SourceXLSX) {
		t.Fatalf("expected the active source as system, got %q", first.Source.System)
	}
	// The upload source's data time is the protocol imported_at, never the sync stamp.
	if first.Source.ObservedAt == nil || !first.Source.ObservedAt.Equal(importedAt) {
		t.Fatalf("expected observed_at %s, got %+v", importedAt, first.Source.ObservedAt)
	}
	if len(first.QualityFlags) != 0 {
		t.Fatalf("a known quantity carries no quality flag, got %v", first.QualityFlags)
	}
}

func TestBatchStockFlagsAProductWhoseStockTheSourceDidNotReport(t *testing.T) {
	repo := &fakeRepo{rows: mirrorRows()}
	ctx := WithActiveSource(context.Background(), erpdomain.SourceXLSX)

	facts, err := batchReader(t, repo).GetStockFactsByIDs(ctx, []int64{10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fact := facts[10]
	if fact == nil {
		t.Fatal("expected a fact for the known product 10")
	}
	// Honest unknown: never a zero, which the risk engine would read as sold out.
	if fact.Quantity != nil {
		t.Fatalf("expected an unknown quantity, got %v", *fact.Quantity)
	}
	if !readdomain.HasQualityFlag(fact.QualityFlags, readdomain.QualityMissingStock) {
		t.Fatalf("expected a missing-stock flag, got %v", fact.QualityFlags)
	}
}

func TestBatchStockOmitsProductsThatAreNotInTheActiveSource(t *testing.T) {
	repo := &fakeRepo{rows: mirrorRows()}
	ctx := WithActiveSource(context.Background(), erpdomain.SourceXLSX)

	facts, err := batchReader(t, repo).GetStockFactsByIDs(ctx, []int64{999})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(facts) != 0 {
		t.Fatalf("an unknown product has no fact, got %+v", facts)
	}
}

func TestBatchStockFailsClosedWithoutAnActiveSource(t *testing.T) {
	repo := &fakeRepo{rows: mirrorRows()}

	if _, err := batchReader(t, repo).GetStockFactsByIDs(context.Background(), []int64{1}); err == nil {
		t.Fatal("expected an unavailable-source error without an active source pin")
	}
}
