//go:build integration

package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestRepositoryCompletedPullIsAtomicAndTenantScoped(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings")
	token := time.Now().UTC().Format("150405.000000000")
	tenantA, tenantB, installation := "listing-a-"+token, "listing-b-"+token, "installation-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_sync_events WHERE tenant_id = ANY($1) AND installation_id = $2`, []string{tenantA, tenantB}, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id = ANY($1) AND installation_id = $2`, []string{tenantA, tenantB}, installation)
	})

	created := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	seedListing(t, pool, tenantA, installation, "same", "-", "old", "active", created, created)
	seedListing(t, pool, tenantA, installation, "omitted", "-", "old omitted", "active", created, created)
	seedListing(t, pool, tenantB, installation, "same", "-", "tenant b", "active", created, created)
	beforeB := readRawListing(t, pool, tenantB, installation, "same", "-")

	completed := created.Add(24 * time.Hour)
	fetched := completed.Add(-time.Minute)
	price := domain.PriceAmount("123.45")
	currency := domain.PriceCurrency("BRL")
	quantity := 7
	rows := []domain.Listing{
		listing(t, tenantA, installation, "same", "-", domain.ListingStatusActive, "updated", &price, &currency, &quantity, fetched, 80),
		listing(t, tenantA, installation, "same", "variation-2", domain.ListingStatusPaused, "variation", nil, nil, nil, fetched, 70),
	}
	repo := listingspostgres.NewRepository(pool, tenantA)
	if err := repo.ApplyCompletedPull(ctx, installation, rows, completed); err != nil {
		t.Fatalf("ApplyCompletedPull(): %v", err)
	}

	assertListing(t, pool, tenantA, installation, "same", "-", "active", "updated", created, completed, fetched, "123.45", "BRL", 7)
	assertListing(t, pool, tenantA, installation, "same", "variation-2", "paused", "variation", fetched, completed, fetched, "", "", -1)
	assertStatus(t, pool, tenantA, installation, "omitted", "-", "closed")
	afterB := readRawListing(t, pool, tenantB, installation, "same", "-")
	if !reflect.DeepEqual(beforeB, afterB) {
		t.Fatalf("tenant B changed: before=%#v after=%#v", beforeB, afterB)
	}
	assertEventKinds(t, pool, tenantA, installation, map[string]string{"same/-": "synced", "same/variation-2": "paused", "omitted/-": "closed"})

	beforeRollback := readRawListing(t, pool, tenantA, installation, "same", "-")
	eventsBefore := eventCount(t, pool, tenantA, installation)
	invalid := listing(t, tenantA, installation, "failure", "-", domain.ListingStatusActive, "invalid", nil, nil, nil, fetched.Add(time.Hour), 101)
	if err := repo.ApplyCompletedPull(ctx, installation, []domain.Listing{rows[0], invalid}, completed.Add(time.Hour)); err == nil {
		t.Fatal("ApplyCompletedPull() error = nil, want check failure")
	}
	if got := readRawListing(t, pool, tenantA, installation, "same", "-"); !reflect.DeepEqual(beforeRollback, got) {
		t.Fatalf("listing effects did not roll back: before=%#v after=%#v", beforeRollback, got)
	}
	assertStatus(t, pool, tenantA, installation, "omitted", "-", "closed")
	if got := eventCount(t, pool, tenantA, installation); got != eventsBefore {
		t.Fatalf("event count after rollback = %d, want %d", got, eventsBefore)
	}

	if err := repo.ApplyCompletedPull(ctx, installation, nil, completed.Add(2*time.Hour)); err != nil {
		t.Fatalf("empty ApplyCompletedPull(): %v", err)
	}
	assertAllClosed(t, pool, tenantA, installation)
	if got := readRawListing(t, pool, tenantB, installation, "same", "-"); !reflect.DeepEqual(beforeB, got) {
		t.Fatalf("empty pull changed tenant B: before=%#v after=%#v", beforeB, got)
	}
}

func seedListing(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation, title, status string, created, updated time.Time) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `INSERT INTO listings
		(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,status,sync_state,created_at,updated_at,fetched_at)
		VALUES ($1,$2,'mercado_livre',$3,$4,$5,$6,'synced',$7,$8,$8)`, tenant, installation, item, variation, title, status, created, updated)
	if err != nil {
		t.Fatalf("seed listing: %v", err)
	}
}

func listing(t *testing.T, tenant, installation, item, variation string, status domain.ListingStatus, title string, price *domain.PriceAmount, currency *domain.PriceCurrency, quantity *int, fetched time.Time, quality int) domain.Listing {
	t.Helper()
	row, err := domain.NewListing(domain.ListingInput{TenantID: tenant, InstallationID: installation, Provider: "mercado_livre", ProviderListingID: item, VariationID: variation, Title: title, Status: status, PriceAmount: price, PriceCurrency: currency, PublishedQuantity: quantity, SyncState: domain.ListingSyncStateSynced, QualityScore: &quality, FetchedAt: &fetched, CreatedAt: fetched, UpdatedAt: fetched})
	if err != nil {
		t.Fatal(err)
	}
	return row
}

type rawListing struct {
	Provider, Title, Status, Price, Currency string
	Quantity                                 *int
	Created, Updated, Fetched                time.Time
}

func readRawListing(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation string) rawListing {
	t.Helper()
	var got rawListing
	var price, currency *string
	err := pool.QueryRow(context.Background(), `SELECT provider,title,status,price_amount::text,price_currency,published_quantity,created_at,updated_at,fetched_at FROM listings WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id=$3 AND variation_id=$4`, tenant, installation, item, variation).Scan(&got.Provider, &got.Title, &got.Status, &price, &currency, &got.Quantity, &got.Created, &got.Updated, &got.Fetched)
	if err != nil {
		t.Fatalf("read listing %s/%s: %v", item, variation, err)
	}
	if price != nil {
		got.Price = *price
	}
	if currency != nil {
		got.Currency = *currency
	}
	return got
}
func assertListing(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation, status, title string, created, updated, fetched time.Time, price, currency string, quantity int) {
	t.Helper()
	got := readRawListing(t, pool, tenant, installation, item, variation)
	if got.Status != status || got.Title != title || !got.Created.Equal(created) || !got.Updated.Equal(updated) || !got.Fetched.Equal(fetched) || got.Price != price || got.Currency != currency {
		t.Fatalf("listing %s/%s = %#v", item, variation, got)
	}
	if quantity < 0 {
		if got.Quantity != nil {
			t.Fatalf("quantity = %v, want nil", *got.Quantity)
		}
	} else if got.Quantity == nil || *got.Quantity != quantity {
		t.Fatalf("quantity = %v, want %d", got.Quantity, quantity)
	}
}
func assertStatus(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation, want string) {
	t.Helper()
	var got string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM listings WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id=$3 AND variation_id=$4`, tenant, installation, item, variation).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("status %s/%s = %s, want %s", item, variation, got, want)
	}
}
func assertEventKinds(t *testing.T, pool *pgxpool.Pool, tenant, installation string, want map[string]string) {
	t.Helper()
	rows, err := pool.Query(context.Background(), `SELECT provider_listing_id,variation_id,kind FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var item, v, k string
		if err := rows.Scan(&item, &v, &k); err != nil {
			t.Fatal(err)
		}
		got[item+"/"+v] = k
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("event kinds = %#v, want %#v", got, want)
	}
}
func eventCount(t *testing.T, pool *pgxpool.Pool, tenant, installation string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
func assertAllClosed(t *testing.T, pool *pgxpool.Pool, tenant, installation string) {
	t.Helper()
	var open int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM listings WHERE tenant_id=$1 AND installation_id=$2 AND status <> 'closed'`, tenant, installation).Scan(&open); err != nil {
		t.Fatal(err)
	}
	if open != 0 {
		t.Fatalf("non-closed rows = %d", open)
	}
}
