//go:build integration

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/composition"
	"marketplace-central/apps/server_core/internal/contexts/listings"
	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/exact"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// queryTenantRow runs query inside a transaction that first pins
// app.tenant_id, mirroring scanOne (listings_ingest_test.go) but generalized
// to an arbitrary destination list -- this test needs to read a version, a
// state/value pair, a row count and a timestamp back from the database, not
// just a single int.
func queryTenantRow(pool *pgxpool.Pool, tid tenant.ID, query string, args []any, dest ...any) error {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tid.String()); err != nil {
		return err
	}
	if err := tx.QueryRow(ctx, query, args...).Scan(dest...); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// TestListingsReprocessHealsStaleFactsWithoutCallingTheChannel reproduces a
// real incident hermetically: 34 live listings were stored with seller_sku
// Unknown/"ml omitted seller_sku" because the old mapper never read the
// SELLER_SKU attribute, even though the payload sitting right next to each
// row carried the SKU. The mapper was corrected (api.ItemSellerSKU); the
// payload bytes never changed. A reprocess must heal those rows by
// re-deriving facts from the bytes already on disk -- no network call, no
// credential, nothing reachable through composition.StoredListings that can
// reach Mercado Livre.
func TestListingsReprocessHealsStaleFactsWithoutCallingTheChannel(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_listings_reprocess")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	module := listings.New(pool)
	ctx := context.Background()

	const listingID = "MLB-REPRO-1"

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

	// A realistic Mercado Livre item body: the SKU lives in the SELLER_SKU
	// attribute, exactly where the corrected mapper (api.ItemSellerSKU) reads
	// it and the old one did not.
	rawPayload := []byte(`{"id":"MLB-REPRO-1","title":"Caneca Termica Inox 300ml","status":"active","listing_type_id":"gold_special","price":199.90,"currency_id":"BRL","available_quantity":10,"attributes":[{"id":"SELLER_SKU","value_name":"SKU-MLB1"}]}`)
	sum := sha256.Sum256(rawPayload)
	payloadHash := hex.EncodeToString(sum[:])

	observedAt := time.Now().UTC()
	e, err := provenance.NewEvidence("mercadolivre", "item", listingID, observedAt, payloadHash)
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}

	title, err := fact.NewKnown("Caneca Termica Inox 300ml", e)
	if err != nil {
		t.Fatalf("fact.NewKnown title: %v", err)
	}
	status, err := fact.NewKnown("active", e)
	if err != nil {
		t.Fatalf("fact.NewKnown status: %v", err)
	}
	listingType, err := fact.NewKnown("gold_special", e)
	if err != nil {
		t.Fatalf("fact.NewKnown listing_type: %v", err)
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
	// The STALE fact: the old mapper never looked at the SELLER_SKU
	// attribute, so every listing it ingested carries this exact reason even
	// though the payload beside it names the SKU.
	staleSKU, err := fact.NewUnknown[string]("ml omitted seller_sku", e)
	if err != nil {
		t.Fatalf("fact.NewUnknown seller_sku: %v", err)
	}
	gtin, err := fact.NewUnknown[string]("ml omitted gtin attribute", e)
	if err != nil {
		t.Fatalf("fact.NewUnknown gtin: %v", err)
	}

	seeded := contracts.ListingObservation{
		Key: key,
		State: contracts.ListingState{
			Title: title, Status: status, ListingType: listingType,
			Price: price, AvailableQuantity: qty, SellerSKU: staleSKU, GTIN: gtin,
		},
		RawPayload: rawPayload,
		Evidence:   e,
	}

	res, err := module.IngestListing(ctx, seeded)
	if err != nil {
		t.Fatalf("seed ingest: %v", err)
	}
	if res.Disposition != contracts.DispositionCreated || res.Version != 1 {
		t.Fatalf("seed ingest = %+v, want created v1", res)
	}

	// Capture the observation's observed_at as the database recorded it, not
	// the Go value we constructed -- a reprocess that fabricated a new
	// observation would still leave THIS Go value untouched, so re-reading
	// from the database is what actually catches it.
	var seededObservedAt time.Time
	if err := queryTenantRow(pool, tid,
		`SELECT observed_at FROM listings.source_observations WHERE tenant_id=$1 AND listing_id=$2`,
		[]any{tid.String(), listingID}, &seededObservedAt); err != nil {
		t.Fatalf("read seeded observed_at: %v", err)
	}

	// Sanity check on the fixture itself: the row we are about to heal must
	// actually be broken, or a "healed" assertion below would prove nothing.
	var preSKUState, preSKUValue string
	if err := queryTenantRow(pool, tid,
		`SELECT seller_sku_state, coalesce(seller_sku_value, '') FROM listings.listings WHERE tenant_id=$1 AND listing_id=$2`,
		[]any{tid.String(), listingID}, &preSKUState, &preSKUValue); err != nil {
		t.Fatalf("read pre-reprocess seller_sku: %v", err)
	}
	if preSKUState != "unknown" || preSKUValue != "" {
		t.Fatalf("fixture is not stale: seller_sku_state=%q seller_sku_value=%q, want unknown/empty", preSKUState, preSKUValue)
	}

	wiring := composition.WireListingsReprocess(pool)
	report, err := composition.RunListingsReprocess(ctx, wiring.Module, wiring.Mapper, tid, 50)
	if err != nil {
		t.Fatalf("first reprocess: %v", err)
	}
	if report.Pages < 1 || report.Read != 1 || report.Changed != 1 || report.Idempotent != 0 || report.Created != 0 {
		t.Fatalf("first reprocess report = %+v, want Pages>=1 Read=1 Changed=1 Idempotent=0 Created=0", report)
	}

	var (
		skuState, skuValue string
		version            int
	)
	if err := queryTenantRow(pool, tid,
		`SELECT seller_sku_state, seller_sku_value, version FROM listings.listings WHERE tenant_id=$1 AND listing_id=$2`,
		[]any{tid.String(), listingID}, &skuState, &skuValue, &version); err != nil {
		t.Fatalf("read healed listing: %v", err)
	}
	if skuState != "known" || skuValue != "SKU-MLB1" {
		t.Fatalf("seller_sku = state=%q value=%q, want known/SKU-MLB1 (reprocess did not heal the row)", skuState, skuValue)
	}
	if version != 2 {
		t.Fatalf("version = %d, want 2", version)
	}

	var observationRows int
	if err := queryTenantRow(pool, tid,
		`SELECT count(*) FROM listings.source_observations WHERE tenant_id=$1 AND listing_id=$2`,
		[]any{tid.String(), listingID}, &observationRows); err != nil {
		t.Fatalf("count observations: %v", err)
	}
	if observationRows != 1 {
		t.Fatalf("source_observations rows = %d, want 1 (a reprocess observes nothing and must not fabricate a new observation)", observationRows)
	}

	var healedObservedAt time.Time
	if err := queryTenantRow(pool, tid,
		`SELECT observed_at FROM listings.source_observations WHERE tenant_id=$1 AND listing_id=$2`,
		[]any{tid.String(), listingID}, &healedObservedAt); err != nil {
		t.Fatalf("read post-reprocess observed_at: %v", err)
	}
	if !healedObservedAt.Equal(seededObservedAt) {
		t.Fatalf("observed_at changed from %v to %v: a reprocess must not mint a new observation moment", seededObservedAt, healedObservedAt)
	}

	// Run it again: the mapping is deterministic over unchanged bytes, so a
	// second reprocess must converge instead of minting a version every time
	// it runs.
	second, err := composition.RunListingsReprocess(ctx, wiring.Module, wiring.Mapper, tid, 50)
	if err != nil {
		t.Fatalf("second reprocess: %v", err)
	}
	if second.Read != 1 || second.Idempotent != 1 || second.Changed != 0 || second.Created != 0 {
		t.Fatalf("second reprocess report = %+v, want Read=1 Idempotent=1 Changed=0 Created=0", second)
	}

	var versionAfterSecond int
	if err := queryTenantRow(pool, tid,
		`SELECT version FROM listings.listings WHERE tenant_id=$1 AND listing_id=$2`,
		[]any{tid.String(), listingID}, &versionAfterSecond); err != nil {
		t.Fatalf("read version after second reprocess: %v", err)
	}
	if versionAfterSecond != 2 {
		t.Fatalf("version after second reprocess = %d, want 2 (idempotent must not bump the version)", versionAfterSecond)
	}
}
