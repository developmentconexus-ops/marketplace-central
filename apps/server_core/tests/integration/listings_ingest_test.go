//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/listings"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// listingObservation builds a complete observation: every fact known, exactly
// as the plan requires, because 0099's CHECK constraints on title/price reject
// a zero-value Fact (state=unknown, reason="" -> title_reason NULL). The raw
// payload varies with hash so two distinct ingests leave two distinct
// source_observations rows, not one deduplicated by an accidental collision.
func listingObservation(t *testing.T, tid tenant.ID, listingID, hash string) contracts.ListingObservation {
	t.Helper()
	e, err := provenance.NewEvidence("mercadolivre", "listing", listingID, time.Now().UTC(), hash)
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	code, err := channel.ParseCode("mercadolivre")
	if err != nil {
		t.Fatalf("channel.ParseCode: %v", err)
	}
	account, err := channel.NewAccountRef(code, "ml-account-01")
	if err != nil {
		t.Fatalf("channel.NewAccountRef: %v", err)
	}
	key, err := contracts.NewSourceListingKey(tid, account, listingID)
	if err != nil {
		t.Fatalf("contracts.NewSourceListingKey: %v", err)
	}
	title, err := fact.NewKnown("Titulo "+hash, e)
	if err != nil {
		t.Fatalf("fact.NewKnown title: %v", err)
	}
	status, err := fact.NewKnown("active", e)
	if err != nil {
		t.Fatalf("fact.NewKnown status: %v", err)
	}
	listingType, err := fact.NewKnown("gold_special", e)
	if err != nil {
		t.Fatalf("fact.NewKnown listing type: %v", err)
	}
	currency, err := exact.ParseCurrency("BRL")
	if err != nil {
		t.Fatalf("exact.ParseCurrency: %v", err)
	}
	amount, err := exact.ParseMoney("199.90", currency)
	if err != nil {
		t.Fatalf("exact.ParseMoney: %v", err)
	}
	price, err := fact.NewKnown(amount, e)
	if err != nil {
		t.Fatalf("fact.NewKnown price: %v", err)
	}
	qty, err := fact.NewKnown(10, e)
	if err != nil {
		t.Fatalf("fact.NewKnown qty: %v", err)
	}
	sku, err := fact.NewKnown("SKU-"+listingID, e)
	if err != nil {
		t.Fatalf("fact.NewKnown sku: %v", err)
	}
	gtin, err := fact.NewKnown("7890000000"+hash, e)
	if err != nil {
		t.Fatalf("fact.NewKnown gtin: %v", err)
	}

	vPrice, err := fact.NewKnown(amount, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation price: %v", err)
	}
	vQty, err := fact.NewKnown(5, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation qty: %v", err)
	}
	vSku, err := fact.NewKnown("SKU-"+listingID+"-V1", e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation sku: %v", err)
	}
	vGtin, err := fact.NewKnown("7890000001"+hash, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation gtin: %v", err)
	}

	return contracts.ListingObservation{
		Key:               key,
		Title:             title,
		Status:            status,
		ListingType:       listingType,
		Price:             price,
		AvailableQuantity: qty,
		SellerSKU:         sku,
		GTIN:              gtin,
		Variations: []contracts.VariationObservation{
			{
				VariationID:       "V1",
				Price:             vPrice,
				AvailableQuantity: vQty,
				SellerSKU:         vSku,
				GTIN:              vGtin,
			},
		},
		RawPayload: []byte(`{"listing_id":"` + listingID + `","hash":"` + hash + `"}`),
		Evidence:   e,
	}
}

// scanOne runs query inside a transaction that first pins app.tenant_id, so a
// count against a FORCE ROW LEVEL SECURITY table sees the tenant's rows
// instead of the empty world an unscoped session would see (repository.go:49
// of catalog is the ratified precedent this mirrors).
func scanOne(pool *pgxpool.Pool, tid tenant.ID, dest *int, query string) error {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tid.String()); err != nil {
		return err
	}
	if err := tx.QueryRow(ctx, query, tid.String()).Scan(dest); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func TestListingsIngestPersistsVersionAndVariations(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_listings_ingest")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	module := listings.New(pool)
	ctx := context.Background()

	first := listingObservation(t, tid, "MLB-IT-1", "hash-a")
	res, err := module.IngestListing(ctx, first)
	if err != nil {
		t.Fatalf("ingest v1: %v", err)
	}
	if res.Disposition != contracts.DispositionCreated || res.Version != 1 {
		t.Fatalf("got %+v, want created v1", res)
	}

	// A aceitação é o BANCO, não o retorno (aceite = observável).
	var rows int
	if err := scanOne(pool, tid, &rows,
		`SELECT count(*) FROM listings.listings WHERE tenant_id=$1 AND listing_id='MLB-IT-1'`); err != nil {
		t.Fatalf("count listings: %v", err)
	}
	if rows != 1 {
		t.Fatalf("listings rows = %d, want 1", rows)
	}
	if err := scanOne(pool, tid, &rows,
		`SELECT count(*) FROM listings.listing_variations WHERE tenant_id=$1 AND listing_id='MLB-IT-1'`); err != nil {
		t.Fatalf("count variations: %v", err)
	}
	if rows != 1 {
		t.Fatalf("variation rows = %d, want 1", rows)
	}

	second := listingObservation(t, tid, "MLB-IT-1", "hash-b")
	res, err = module.IngestListing(ctx, second)
	if err != nil {
		t.Fatalf("ingest v2: %v", err)
	}
	if res.Disposition != contracts.DispositionChanged || res.Version != 2 {
		t.Fatalf("got %+v, want changed v2", res)
	}
	if err := scanOne(pool, tid, &rows,
		`SELECT count(*) FROM listings.source_observations WHERE tenant_id=$1 AND listing_id='MLB-IT-1'`); err != nil {
		t.Fatalf("count observations: %v", err)
	}
	if rows != 2 {
		t.Fatalf("observation rows = %d, want 2 (one per distinct payload)", rows)
	}
}
