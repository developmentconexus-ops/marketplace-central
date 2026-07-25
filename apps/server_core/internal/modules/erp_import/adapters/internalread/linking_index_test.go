package internalread

import (
	"context"
	"testing"
	"time"

	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

func indexedReader(ttl time.Duration) (*Reader, *fakeRepo, *time.Time) {
	repo := &fakeRepo{rows: mirrorRows()}
	clock := fetchedAt
	reader := NewReader(repo, "tenant-a",
		WithClock(func() time.Time { return clock }),
		WithLinkingIndexTTL(ttl),
	)
	return reader, repo, &clock
}

// A candidate-generation run asks the matcher about every listing in an account.
// Re-reading the whole mirror per question is what made a full run impossible,
// so one loaded snapshot must serve every anchor.
func TestLinkingLookupsShareOneMirrorRead(t *testing.T) {
	reader, repo, _ := indexedReader(time.Minute)
	ctx := WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente)

	byEAN, err := reader.FindProductsForLinking(ctx, readports.FindProductsInput{EAN: ptr("7894900011517")})
	if err != nil || len(byEAN) != 2 {
		t.Fatalf("EAN lookup = %+v err=%v, want both colliding CODPRODs", byEAN, err)
	}
	bySKU, err := reader.FindProductsForLinking(ctx, readports.FindProductsInput{SellerSKU: ptr("10")})
	if err != nil || len(bySKU) != 1 || bySKU[0].ProductID != 10 {
		t.Fatalf("SKU lookup = %+v err=%v", bySKU, err)
	}
	byTitle, err := reader.FindProductsForLinking(ctx, readports.FindProductsInput{Title: ptr("blue")})
	if err != nil || len(byTitle) != 2 {
		t.Fatalf("title lookup = %+v err=%v", byTitle, err)
	}

	if repo.mirrorRowsCalls != 1 {
		t.Fatalf("mirror reads = %d, want 1 shared snapshot", repo.mirrorRowsCalls)
	}
}

// The memo is a read-through cache, never a second source of truth: a sync that
// rewrites the mirror must become visible without a restart.
func TestLinkingIndexIsRebuiltAfterItsTTL(t *testing.T) {
	reader, repo, clock := indexedReader(time.Minute)
	ctx := WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente)

	if _, err := reader.FindProductsForLinking(ctx, readports.FindProductsInput{Title: ptr("blue")}); err != nil {
		t.Fatalf("first lookup error = %v", err)
	}
	*clock = clock.Add(time.Minute + time.Second)
	repo.rows = append(mirrorRows(), erpdomain.MirrorProduct{
		CodigoProduto: "11", Descricao: ptr("Blue Arrival"),
		ImportedAt: timePtr(importedAt), UpdatedAt: importedAt,
	})

	got, err := reader.FindProductsForLinking(ctx, readports.FindProductsInput{Title: ptr("blue")})
	if err != nil || len(got) != 3 {
		t.Fatalf("post-TTL lookup = %+v err=%v, want the newly synced row", got, err)
	}
	if repo.mirrorRowsCalls != 2 {
		t.Fatalf("mirror reads = %d, want a rebuild after the TTL", repo.mirrorRowsCalls)
	}
}

// Sources never fall back to one another; a cached snapshot of one source must
// never answer for another.
func TestLinkingIndexIsRebuiltWhenTheActiveSourceChanges(t *testing.T) {
	reader, repo, _ := indexedReader(time.Hour)

	if _, err := reader.FindProductsForLinking(WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente), readports.FindProductsInput{Title: ptr("blue")}); err != nil {
		t.Fatalf("catalogo_cliente lookup error = %v", err)
	}
	if _, err := reader.FindProductsForLinking(WithActiveSource(context.Background(), erpdomain.SourceXLSX), readports.FindProductsInput{Title: ptr("blue")}); err != nil {
		t.Fatalf("xlsx lookup error = %v", err)
	}
	if repo.mirrorRowsCalls != 2 || repo.source != erpdomain.SourceXLSX {
		t.Fatalf("reads=%d last source=%q, want a per-source snapshot", repo.mirrorRowsCalls, repo.source)
	}
}

// Narrowing is an optimization only: for every anchor, the rows handed to
// matches() must still contain every row a full scan would have matched.
func TestNarrowedCandidateRowsAgreeWithAFullScan(t *testing.T) {
	rows := append(mirrorRows(), erpdomain.MirrorProduct{
		CodigoProduto: "12", Descricao: ptr("UPPERCASE BLUE ITEM"), EAN: ptr("7896004000015"),
		ImportedAt: timePtr(importedAt), UpdatedAt: importedAt,
	})
	index := buildIndex(t, rows)

	inputs := []readports.FindProductsInput{
		{},
		{ProductID: ptr(2)},
		{SellerSKU: ptr("10")},
		{EAN: ptr("7894900011517")},
		{EAN: ptr("7896004000015")},
		{EAN: ptr("bad")},
		{Title: ptr("blue")},
		{Title: ptr("BLUE")},
		{Title: ptr("nothing matches this")},
	}
	for _, input := range inputs {
		want := map[string]bool{}
		for _, row := range rows {
			if matches(row, input) {
				want[row.CodigoProduto] = true
			}
		}
		got := map[string]bool{}
		for _, row := range index.candidateRows(input) {
			if matches(row, input) {
				got[row.CodigoProduto] = true
			}
		}
		if len(got) != len(want) {
			t.Fatalf("input %+v: narrowed matches = %v, full scan = %v", input, got, want)
		}
		for code := range want {
			if !got[code] {
				t.Fatalf("input %+v: narrowing dropped CODPROD %s", input, code)
			}
		}
	}
}

func buildIndex(t *testing.T, rows []erpdomain.MirrorProduct) *linkingIndex {
	t.Helper()
	reader := NewReader(&fakeRepo{rows: rows}, "tenant-a", WithClock(func() time.Time { return fetchedAt }))
	index, err := reader.linkingIndex(WithActiveSource(context.Background(), erpdomain.SourceCatalogoCliente), erpdomain.SourceCatalogoCliente)
	if err != nil {
		t.Fatalf("linkingIndex() error = %v", err)
	}
	return index
}
