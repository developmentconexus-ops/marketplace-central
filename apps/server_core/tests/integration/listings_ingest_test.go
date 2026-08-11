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
		Key: key,
		State: contracts.ListingState{
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

// TestListingsCurrentRoundTripsState proves Current() reconstructs the exact
// facts SaveVersion wrote -- a mix of Known, Unknown and NotApplicable states
// plus a variation -- rather than just the version and a hash. A column read
// in the wrong order, or a state string that does not map back to the right
// constructor, shows up here as Equal returning false, not as a panic or a
// wrong type: this is the test that catches it.
func TestListingsCurrentRoundTripsState(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_listings_current")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	module := listings.New(pool)
	ctx := context.Background()

	e, err := provenance.NewEvidence("mercadolivre", "listing", "MLB-CUR-1", time.Now().UTC(), "hash-current")
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
	key, err := contracts.NewSourceListingKey(tid, account, "MLB-CUR-1")
	if err != nil {
		t.Fatalf("contracts.NewSourceListingKey: %v", err)
	}
	title, err := fact.NewKnown("Titulo current", e)
	if err != nil {
		t.Fatalf("fact.NewKnown title: %v", err)
	}
	status, err := fact.NewUnknown[string]("ml omitted status", e)
	if err != nil {
		t.Fatalf("fact.NewUnknown status: %v", err)
	}
	listingType, err := fact.NewNotApplicable[string]("classified listing has no type", e)
	if err != nil {
		t.Fatalf("fact.NewNotApplicable listing_type: %v", err)
	}
	currency, err := exact.ParseCurrency("BRL")
	if err != nil {
		t.Fatalf("exact.ParseCurrency: %v", err)
	}
	amount, err := exact.ParseMoney("149.50", currency)
	if err != nil {
		t.Fatalf("exact.ParseMoney: %v", err)
	}
	price, err := fact.NewKnown(amount, e)
	if err != nil {
		t.Fatalf("fact.NewKnown price: %v", err)
	}
	qty, err := fact.NewUnknown[int]("ml omitted available_quantity", e)
	if err != nil {
		t.Fatalf("fact.NewUnknown qty: %v", err)
	}
	sku, err := fact.NewKnown("SKU-CUR-1", e)
	if err != nil {
		t.Fatalf("fact.NewKnown sku: %v", err)
	}
	gtin, err := fact.NewUnknown[string]("ml omitted gtin attribute", e)
	if err != nil {
		t.Fatalf("fact.NewUnknown gtin: %v", err)
	}
	vPrice, err := fact.NewKnown(amount, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation price: %v", err)
	}
	vQty, err := fact.NewKnown(3, e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation qty: %v", err)
	}
	vSku, err := fact.NewKnown("SKU-CUR-1-V1", e)
	if err != nil {
		t.Fatalf("fact.NewKnown variation sku: %v", err)
	}
	vGtin, err := fact.NewNotApplicable[string]("variation has no gtin of its own", e)
	if err != nil {
		t.Fatalf("fact.NewNotApplicable variation gtin: %v", err)
	}

	state := contracts.ListingState{
		Title: title, Status: status, ListingType: listingType,
		Price: price, AvailableQuantity: qty, SellerSKU: sku, GTIN: gtin,
		Variations: []contracts.VariationObservation{
			{VariationID: "V1", Price: vPrice, AvailableQuantity: vQty, SellerSKU: vSku, GTIN: vGtin},
		},
	}
	observation := contracts.ListingObservation{
		Key: key, State: state,
		RawPayload: []byte(`{"listing_id":"MLB-CUR-1"}`), Evidence: e,
	}

	if _, err := module.IngestListing(ctx, observation); err != nil {
		t.Fatalf("ingest: %v", err)
	}

	// The round-trip is proven twice, and neither proof needs a public method
	// that exists only for this test to call. A façade method with no
	// production caller is the orphan surface this slice's own plan forbids,
	// and it buys nothing here: both assertions below are strictly stronger.
	//
	// First through the door production uses. Folding the identical
	// observation again must report idempotent, and that verdict IS the
	// round-trip: Ingest reaches it by reconstructing the stored row through
	// Current and comparing it to this observation with ListingState.Equal. A
	// reader that dropped a fact, mangled a state, or lost the variation would
	// report changed here and bump the version.
	again, err := module.IngestListing(ctx, observation)
	if err != nil {
		t.Fatalf("second ingest: %v", err)
	}
	if again.Disposition != contracts.DispositionIdempotent {
		t.Fatalf("second ingest of an identical observation = %q at version %d, want %q at version 1 -- the stored row did not read back equal to what was written",
			again.Disposition, again.Version, contracts.DispositionIdempotent)
	}
	if again.Version != 1 {
		t.Fatalf("version = %d after an idempotent fold, want 1", again.Version)
	}

	// Second on the columns themselves, in SQL none of this context's Go code
	// touches. The fold above compares a reconstruction against the
	// observation, so a defect that writes a fact to the wrong column and
	// reads it back from that same wrong column is invisible to it. This is
	// not: it names the column it expected.
	var (
		skuState, skuValue     string
		gtinState, gtinReason  string
		priceAmount, priceCurr string
		variationSKU           string
	)
	if err := queryTenantRow(pool, tid,
		`SELECT seller_sku_state, coalesce(seller_sku_value,''), gtin_state, coalesce(gtin_reason,''),
		        coalesce(price_amount::text,''), coalesce(price_currency,'')
		 FROM listings.listings WHERE tenant_id=$1 AND listing_id=$2`,
		[]any{tid.String(), "MLB-CUR-1"},
		&skuState, &skuValue, &gtinState, &gtinReason, &priceAmount, &priceCurr); err != nil {
		t.Fatalf("read listing columns: %v", err)
	}
	if skuState != "known" || skuValue != "SKU-CUR-1" {
		t.Fatalf("seller_sku columns = state=%q value=%q, want known/SKU-CUR-1", skuState, skuValue)
	}
	// An unknown fact keeps its reason and stores no value: the knowledge
	// state is data, not an absence to be read as zero.
	if gtinState != "unknown" || gtinReason != "ml omitted gtin attribute" {
		t.Fatalf("gtin columns = state=%q reason=%q, want unknown with the reason preserved", gtinState, gtinReason)
	}
	if priceAmount != "149.50" || priceCurr != "BRL" {
		t.Fatalf("price columns = amount=%q currency=%q, want 149.50/BRL", priceAmount, priceCurr)
	}
	if err := queryTenantRow(pool, tid,
		`SELECT coalesce(seller_sku_value,'') FROM listings.listing_variations
		 WHERE tenant_id=$1 AND listing_id=$2 AND variation_id=$3`,
		[]any{tid.String(), "MLB-CUR-1", "V1"}, &variationSKU); err != nil {
		t.Fatalf("read variation columns: %v", err)
	}
	if variationSKU != "SKU-CUR-1-V1" {
		t.Fatalf("variation seller_sku_value = %q, want SKU-CUR-1-V1", variationSKU)
	}
}
