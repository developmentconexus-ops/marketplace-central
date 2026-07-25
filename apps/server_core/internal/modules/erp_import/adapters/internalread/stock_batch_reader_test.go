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

// The pricing screen works on the products that are linked to listings, and
// those ids are scattered across the whole catalog. Asking for them by id must
// return exactly the ones the active source carries, in id order, with the
// snapshot's own observation time.
func TestCatalogFactsByIDsReturnsOnlyTheProductsTheActiveSourceCarries(t *testing.T) {
	repo := &fakeRepo{rows: mirrorRows()}
	reader := NewReader(repo, "tenant-a", WithClock(func() time.Time { return fetchedAt }))
	ctx := WithActiveSource(context.Background(), erpdomain.SourceXLSX)

	page, err := reader.CatalogProductFactsByIDs(ctx, []int64{10, 1, 999, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(page.Items) != 2 {
		t.Fatalf("expected the 2 known products, got %+v", page.Items)
	}
	if page.Items[0].InternalProductID != 1 || page.Items[1].InternalProductID != 10 {
		t.Fatalf("expected id order, got %d then %d", page.Items[0].InternalProductID, page.Items[1].InternalProductID)
	}
	// A set is not a page of a longer sequence.
	if page.NextCursor != nil {
		t.Fatalf("next cursor = %+v, want none", page.NextCursor)
	}
	// as_of is the snapshot's data time, never the moment of the read.
	if !page.AsOf.Equal(importedAt.Add(2 * time.Minute)) {
		t.Fatalf("as_of = %s, want the newest observed row", page.AsOf)
	}
}

func TestCatalogFactsByIDsFailsClosedWithoutAnActiveSource(t *testing.T) {
	reader := NewReader(&fakeRepo{rows: mirrorRows()}, "tenant-a")
	if _, err := reader.CatalogProductFactsByIDs(context.Background(), []int64{1}); err == nil {
		t.Fatal("expected an honest failure without an active source")
	}
}
