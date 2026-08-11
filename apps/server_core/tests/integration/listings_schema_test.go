//go:build integration

package integration

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/listings"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// withTenant runs fn inside a transaction with app.tenant_id pinned, which is
// what every FORCE ROW LEVEL SECURITY table in the listings schema requires of
// a reader or a writer alike. It returns the transaction's error verbatim: a
// test about a CHECK constraint has to see the constraint's own error, not a
// wrapped substitute.
func withTenant(t *testing.T, pool *pgxpool.Pool, tid tenant.ID, fn func(pgx.Tx) error) error {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tid.String()); err != nil {
		t.Fatalf("set_config: %v", err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	// A statement that succeeded has to survive the function, or a row seeded
	// for a later assertion is silently rolled back and the assertion that
	// depends on it fails against the wrong constraint.
	return tx.Commit(ctx)
}

// TestSourceObservationKeepsTheChannelsExactBytes is the falsifier for the
// column type of listings.source_observations.payload.
//
// The table exists to hold "the real bytes" so a human reconciling against
// Mercado Livre can hash what is stored and get the payload_hash stored beside
// it. jsonb cannot do that job: it is a parsed value, so it reorders object
// keys and normalises whitespace, and the round trip returns a document that
// hashes to something else. The payload below is deliberately written in
// non-canonical form — keys out of alphabetical order, irregular spacing — so
// that a jsonb column is observably wrong rather than accidentally right.
func TestSourceObservationKeepsTheChannelsExactBytes(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_listings_bytes")
	tid, err := tenant.Parse("tenant_listings_bytes")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	module := listings.New(pool)

	raw := []byte(`{"z_last":1,   "a_first":"x",  "nested":{"b":2,"a":1}}`)
	sum := sha256.Sum256(raw)
	hash := hex.EncodeToString(sum[:])

	// The helper records hash as the evidence's payload_hash; substituting the
	// payload it built for the bytes that actually hash to it is what makes
	// the round-trip assertion below meaningful.
	obs := listingObservation(t, tid, "MLB-BYTES-1", hash)
	obs.RawPayload = raw

	if _, err := module.IngestListing(context.Background(), obs); err != nil {
		t.Fatalf("ingest: %v", err)
	}

	var stored []byte
	if err := withTenant(t, pool, tid, func(tx pgx.Tx) error {
		return tx.QueryRow(context.Background(),
			`SELECT payload FROM listings.source_observations
			 WHERE tenant_id=$1 AND listing_id='MLB-BYTES-1' AND payload_hash=$2`,
			tid.String(), hash).Scan(&stored)
	}); err != nil {
		t.Fatalf("read back payload: %v", err)
	}

	if !bytes.Equal(stored, raw) {
		t.Fatalf("stored payload = %q, want the channel's own bytes %q", stored, raw)
	}
	// Stated as its own assertion because it is the property that matters: the
	// stored bytes must reproduce the hash stored next to them. bytes.Equal
	// above could be satisfied by a column that happened to preserve this one
	// payload; this one says why the column has to.
	roundTrip := sha256.Sum256(stored)
	if got := hex.EncodeToString(roundTrip[:]); got != hash {
		t.Fatalf("sha256(stored) = %s, want the recorded payload_hash %s", got, hash)
	}
}

// insertListingRow writes one listings.listings row with every state column
// forced to the caller's values, bypassing the domain. That is deliberate:
// fact.NewKnown cannot construct a known-with-no-value, so the only thing that
// can prove the database refuses one is SQL. These CHECKs are the backstop for
// a writer bug, and a backstop is only real if something has tried to get past
// it.
func insertListingRow(t *testing.T, pool *pgxpool.Pool, tid tenant.ID, listingID string, cols map[string]string) error {
	t.Helper()
	base := map[string]string{
		"tenant_id":            "'" + tid.String() + "'",
		"channel":              "'mercadolivre'",
		"account_external_id":  "'ml-account-01'",
		"listing_id":           "'" + listingID + "'",
		"version":              "1",
		"title_state":          "'known'",
		"title_value":          "'t'",
		"title_reason":         "NULL",
		"status_state":         "'known'",
		"status_value":         "'active'",
		"status_reason":        "NULL",
		"listing_type_state":   "'known'",
		"listing_type_value":   "'gold_special'",
		"listing_type_reason":  "NULL",
		"price_state":          "'known'",
		"price_amount":         "'1.00'",
		"price_currency":       "'BRL'",
		"price_reason":         "NULL",
		"available_qty_state":  "'known'",
		"available_qty_value":  "1",
		"available_qty_reason": "NULL",
		"seller_sku_state":     "'known'",
		"seller_sku_value":     "'SKU'",
		"seller_sku_reason":    "NULL",
		"gtin_state":           "'known'",
		"gtin_value":           "'789'",
		"gtin_reason":          "NULL",
		"last_payload_hash":    "'h'",
	}
	for k, v := range cols {
		if _, ok := base[k]; !ok {
			t.Fatalf("override names column %q, which listings.listings does not have", k)
		}
		base[k] = v
	}
	names := make([]string, 0, len(base))
	values := make([]string, 0, len(base))
	for k, v := range base {
		names = append(names, k)
		values = append(values, v)
	}
	return withTenant(t, pool, tid, func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO listings.listings (`+strings.Join(names, ",")+`) VALUES (`+strings.Join(values, ",")+`)`)
		return err
	})
}

// TestListingsRejectsKnownFactWithNoValue walks every fact triple on
// listings.listings. Constraining title and price alone would have let the
// schema say a malformed status is acceptable and a malformed title is not —
// a distinction nothing in the context draws.
func TestListingsRejectsKnownFactWithNoValue(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_listings_checks")
	tid, err := tenant.Parse("tenant_listings_checks")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}

	cases := []struct {
		field      string
		override   map[string]string
		constraint string
	}{
		{"title", map[string]string{"title_value": "NULL"}, "listings_title_consistent"},
		{"status", map[string]string{"status_value": "NULL"}, "listings_status_consistent"},
		{"listing_type", map[string]string{"listing_type_value": "NULL"}, "listings_listing_type_consistent"},
		{"price", map[string]string{"price_amount": "NULL"}, "listings_price_consistent"},
		{"available_qty", map[string]string{"available_qty_value": "NULL"}, "listings_available_qty_consistent"},
		{"seller_sku", map[string]string{"seller_sku_value": "NULL"}, "listings_seller_sku_consistent"},
		{"gtin", map[string]string{"gtin_value": "NULL"}, "listings_gtin_consistent"},
	}
	for _, tc := range cases {
		t.Run(tc.field, func(t *testing.T) {
			err := insertListingRow(t, pool, tid, "MLB-CHK-"+tc.field, tc.override)
			if err == nil {
				t.Fatalf("insert with %s_state='known' and no value succeeded; the schema accepts a fact that claims to be known and carries nothing", tc.field)
			}
			// Naming the constraint, not just "some error": a NOT NULL
			// violation or an RLS refusal would also be non-nil, and neither
			// would mean the triple is enforced.
			if !strings.Contains(err.Error(), tc.constraint) {
				t.Fatalf("error = %v, want it to name constraint %s", err, tc.constraint)
			}
		})
	}
}

// TestListingsRejectsUnknownFactWithNoReason is the other half of the triple.
// An unknown with no reason is the exact shape global constraint 9 forbids: it
// is indistinguishable from a column nobody got round to filling in.
func TestListingsRejectsUnknownFactWithNoReason(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_listings_checks")
	tid, err := tenant.Parse("tenant_listings_checks")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	err = insertListingRow(t, pool, tid, "MLB-CHK-unknown-noreason", map[string]string{
		"gtin_state":  "'unknown'",
		"gtin_value":  "NULL",
		"gtin_reason": "NULL",
	})
	if err == nil {
		t.Fatal("insert with gtin_state='unknown' and gtin_reason NULL succeeded; an unknown with no reason is a blank pretending to be a fact")
	}
	if !strings.Contains(err.Error(), "listings_gtin_consistent") {
		t.Fatalf("error = %v, want it to name constraint listings_gtin_consistent", err)
	}
}

// TestListingVariationsRejectsInconsistentFacts covers the child table, which
// carried no fact CHECK of any kind. A multi-variation seller's real prices
// live here, not in the parent row, so leaving it unconstrained meant the
// schema enforced the triple precisely where ML puts a placeholder and stopped
// enforcing it precisely where ML puts the numbers that matter.
func TestListingVariationsRejectsInconsistentFacts(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_listings_checks")
	tid, err := tenant.Parse("tenant_listings_checks")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	const parent = "MLB-CHK-VAR-PARENT"
	if err := insertListingRow(t, pool, tid, parent, nil); err != nil &&
		!strings.Contains(err.Error(), "listings_pkey") {
		t.Fatalf("seed parent listing: %v", err)
	}

	cases := []struct {
		name       string
		cols       string
		values     string
		constraint string
	}{
		{
			"price known with no amount",
			"price_state,price_amount,price_currency,price_reason",
			"'known',NULL,'BRL',NULL",
			"listing_variations_price_consistent",
		},
		{
			"available_qty unknown with no reason",
			"available_qty_state,available_qty_value,available_qty_reason",
			"'unknown',NULL,NULL",
			"listing_variations_available_qty_consistent",
		},
		{
			"seller_sku known with no value",
			"seller_sku_state,seller_sku_value,seller_sku_reason",
			"'known',NULL,NULL",
			"listing_variations_seller_sku_consistent",
		},
		{
			// A state outside the enum violates both gtin CHECKs at once —
			// it is in no branch of the consistency rule either — and
			// Postgres reports whichever it evaluates first. The expectation
			// is the shared prefix, which still separates a fact constraint
			// firing from an RLS refusal, a NOT NULL, or a foreign key.
			"gtin state is not a state at all",
			"gtin_state,gtin_value,gtin_reason",
			"'probably',NULL,NULL",
			"listing_variations_gtin_",
		},
	}
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Every column not under test is written in a valid shape, so the
			// only thing that can reject the row is the constraint named.
			defaults := map[string]string{
				"price_state":          "'known'",
				"price_amount":         "'1.00'",
				"price_currency":       "'BRL'",
				"price_reason":         "NULL",
				"available_qty_state":  "'known'",
				"available_qty_value":  "1",
				"available_qty_reason": "NULL",
				"seller_sku_state":     "'known'",
				"seller_sku_value":     "'SKU'",
				"seller_sku_reason":    "NULL",
				"gtin_state":           "'known'",
				"gtin_value":           "'789'",
				"gtin_reason":          "NULL",
			}
			overrideCols := strings.Split(tc.cols, ",")
			overrideVals := strings.Split(tc.values, ",")
			if len(overrideCols) != len(overrideVals) {
				t.Fatalf("case %d lists %d columns and %d values", i, len(overrideCols), len(overrideVals))
			}
			for j, c := range overrideCols {
				if _, ok := defaults[c]; !ok {
					t.Fatalf("case %d names column %q, which listing_variations does not have", i, c)
				}
				defaults[c] = overrideVals[j]
			}
			names := []string{"tenant_id", "channel", "account_external_id", "listing_id", "variation_id"}
			values := []string{"'" + tid.String() + "'", "'mercadolivre'", "'ml-account-01'", "'" + parent + "'", "'V" + tc.name + "'"}
			for k, v := range defaults {
				names = append(names, k)
				values = append(values, v)
			}
			err := withTenant(t, pool, tid, func(tx pgx.Tx) error {
				_, err := tx.Exec(context.Background(),
					`INSERT INTO listings.listing_variations (`+strings.Join(names, ",")+`) VALUES (`+strings.Join(values, ",")+`)`)
				return err
			})
			if err == nil {
				t.Fatalf("insert (%s) succeeded; listing_variations enforces nothing about its facts", tc.name)
			}
			if !strings.Contains(err.Error(), tc.constraint) {
				t.Fatalf("error = %v, want it to name constraint %s", err, tc.constraint)
			}
		})
	}
}
