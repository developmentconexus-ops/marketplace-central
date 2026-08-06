//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/composition"
	"marketplace-central/apps/server_core/internal/contexts/catalog"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/contexts/catalog/port"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// twoPageStubFeed is a legal stand-in for Oracle HERE because what is under
// test is the composition path and the writer, not the source adapter. It
// hands back the same three observations on every walk so a second run can
// prove idempotence.
type twoPageStubFeed struct {
	pages []port.Page
	calls int
}

func (f *twoPageStubFeed) NextPage(ctx context.Context, t tenant.ID, after port.Cursor, limit int) (port.Page, error) {
	if f.calls >= len(f.pages) {
		return port.Page{}, fmt.Errorf("stub feed: called past its last page (call %d)", f.calls+1)
	}
	page := f.pages[f.calls]
	f.calls++
	return page, nil
}

func (f *twoPageStubFeed) reset() { f.calls = 0 }

var _ port.ProductFeed = (*twoPageStubFeed)(nil)

// TestCatalogIngestWritesRows is the acceptance this whole foundation lacked:
// the slice had green lanes and four empty tables. A stub feed walked through
// the real composition path (RunCatalogIngest) against the real Postgres lane
// must land real rows, and re-polling the same feed must be free.
func TestCatalogIngestWritesRows(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_catalog_composition")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	module := catalog.New(pool)
	ctx := context.Background()

	stamp := itoa(time.Now().UTC().UnixNano())
	obs1 := catalogObservation(t, tid, "comp-1-"+stamp, "Produto Um", "eanc1"+stamp, "sha256:c1"+stamp)
	obs2 := catalogObservation(t, tid, "comp-2-"+stamp, "Produto Dois", "eanc2"+stamp, "sha256:c2"+stamp)
	obs3 := catalogObservation(t, tid, "comp-3-"+stamp, "Produto Tres", "eanc3"+stamp, "sha256:c3"+stamp)

	feed := &twoPageStubFeed{pages: []port.Page{
		{Observations: []contracts.ProductObservation{obs1, obs2}, Next: port.NewCursor("page-2"), Done: false},
		{Observations: []contracts.ProductObservation{obs3}, Done: true},
	}}

	report, err := composition.RunCatalogIngest(ctx, module, feed, tid, 2)
	if err != nil {
		t.Fatalf("RunCatalogIngest (first walk): %v", err)
	}
	if report.Created != 3 || report.Pages != 2 {
		t.Fatalf("first walk report = %+v, want Created=3 Pages=2", report)
	}

	var productCount int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM catalog.products WHERE tenant_id = $1`, tid.String()).Scan(&productCount); err != nil {
		t.Fatalf("count catalog.products: %v", err)
	}
	if productCount != 3 {
		t.Fatalf("catalog.products count = %d, want 3", productCount)
	}

	var observationCount int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM catalog.source_observations WHERE tenant_id = $1`, tid.String()).Scan(&observationCount); err != nil {
		t.Fatalf("count catalog.source_observations: %v", err)
	}
	if observationCount != 3 {
		t.Fatalf("catalog.source_observations count = %d, want 3", observationCount)
	}

	// Re-poll the same feed: re-polling must be free, not a second write.
	feed.reset()
	second, err := composition.RunCatalogIngest(ctx, module, feed, tid, 2)
	if err != nil {
		t.Fatalf("RunCatalogIngest (second walk): %v", err)
	}
	if second.Idempotent != 3 || second.Created != 0 {
		t.Fatalf("second walk report = %+v, want Idempotent=3 Created=0", second)
	}

	var productCountAfter int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM catalog.products WHERE tenant_id = $1`, tid.String()).Scan(&productCountAfter); err != nil {
		t.Fatalf("count catalog.products after re-poll: %v", err)
	}
	if productCountAfter != productCount {
		t.Fatalf("catalog.products count changed on re-poll: %d then %d", productCount, productCountAfter)
	}

	var observationCountAfter int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM catalog.source_observations WHERE tenant_id = $1`, tid.String()).Scan(&observationCountAfter); err != nil {
		t.Fatalf("count catalog.source_observations after re-poll: %v", err)
	}
	if observationCountAfter != observationCount {
		t.Fatalf("catalog.source_observations count changed on re-poll: %d then %d", observationCount, observationCountAfter)
	}
}
