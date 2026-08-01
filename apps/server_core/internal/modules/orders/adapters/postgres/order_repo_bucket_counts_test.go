//go:build integration

package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// newBucketCountsRepositoryForTest is a standalone fixture (own tenant, own
// cleanup) rather than reusing newOrderRepositoryForTest's fixed order-id set
// (order_repo_test.go), because this test needs its own set of provider_order_ids
// with a mix of populated/absent `bucket` columns.
func newBucketCountsRepositoryForTest(t *testing.T, ctx context.Context) (*OrderRepository, string, []string, func()) {
	t.Helper()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_orders")
	tenantID := fmt.Sprintf("f03-bucket-counts-%d", time.Now().UnixNano())
	repo := NewOrderRepository(pool, tenantID)
	installationID := tenantID + "-installation"
	orderIDs := []string{
		tenantID + "-persisted-enviado",
		tenantID + "-legacy-faturar",
		tenantID + "-persisted-novo",
		tenantID + "-invalid-bucket-faturar",
	}
	return repo, installationID, orderIDs, func() {
		for _, orderID := range orderIDs {
			if _, err := pool.Exec(ctx, `DELETE FROM orders_marketplace_order_items WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3`, tenantID, installationID, orderID); err != nil {
				t.Errorf("cleanup order items %q: %v", orderID, err)
			}
			if _, err := pool.Exec(ctx, `DELETE FROM orders_marketplace_orders WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3`, tenantID, installationID, orderID); err != nil {
				t.Errorf("cleanup order %q: %v", orderID, err)
			}
		}
	}
}

// bucketCountsFixtureOrder builds a minimal MarketplaceOrder with
// ProviderCreatedAt set (GetOrderBucketCounts filters on it) and an explicit
// Bucket. Passing bucket="" (ordersdomain.OrderBucket zero value) simulates a
// pre-F-02 legacy row: UpsertOrders (unlike IngestOrder) never computes a
// bucket itself, so the column round-trips exactly whatever the caller sets —
// an empty string here, same as a row this test's own fallback branch must
// treat identically to SQL NULL.
func bucketCountsFixtureOrder(installationID, orderID, status string, bucket ordersdomain.OrderBucket, tags []string, faturadoAt *time.Time) ordersdomain.MarketplaceOrder {
	created := time.Date(2026, 7, 20, 10, 0, 0, 0, time.UTC)
	return ordersdomain.MarketplaceOrder{
		InstallationID:    installationID,
		ProviderCode:      "mercado_livre",
		ProviderOrderID:   orderID,
		ProviderStatus:    status,
		ProviderCreatedAt: &created,
		ProviderUpdatedAt: &created,
		FetchedAt:         created,
		RawProviderRef:    "/orders/" + orderID,
		CreatedAt:         created,
		UpdatedAt:         created,
		Tags:              tags,
		Bucket:            bucket,
		Items: []ordersdomain.MarketplaceOrderItem{{
			ProviderItemID: orderID + "-item",
			Quantity:       1,
			LinkQuality:    ordersdomain.LinkQualityMissing,
		}},
	}
}

// TestGetOrderBucketCountsUsesPersistedBucketWithLegacyFallback is the F-03
// before/after evidence the brief requires: a row whose PERSISTED bucket
// reflects a real shipped status (order_repo.go:378's old hardcoded
// shipmentStatus="" could never produce "enviado" for a merely-"paid",
// not-yet-faturado order — re-deriving with an empty shipment status would
// have landed it in "faturar" instead) must count under its PERSISTED bucket,
// while a legacy row with no persisted bucket (NULL/empty, pre-F-02) must
// still fall back to the byte-identical old derivation.
//
// Before/after (documented per M03-U2's per-order attribution mandate):
//   - persisted-enviado: provider_status=paid, no delivered tag, bucket
//     column="enviado" (F-02 ingest saw a real shipment status "shipped").
//     OLD behavior (always re-derive with shipmentStatus=""): would have
//     counted as FATURAR (paid, not delivered-tagged, faturado_at NULL).
//     NEW behavior: counts as ENVIADO — this delta is a CORRECTION, not a
//     regression, because the persisted bucket saw the real "shipped" status
//     order_repo.go:378 (pre-F-03) could never see. Attribution: this order's
//     bucket moved faturar -> enviado, caused by shipment status "shipped"
//     recorded at F-02 ingest time.
//   - legacy-faturar: bucket column left "" (simulating a pre-F-02 row that
//     was never re-ingested). Falls back to DeriveOrderBucket(providerStatus,
//     "", tags, faturado) exactly as before F-03 -- counts FATURAR in both
//     old and new behavior. No delta.
//   - persisted-novo: provider_status=created, bucket column="novo" (agrees
//     with what re-derivation would also produce). No delta either way --
//     included to prove the persisted-bucket path is not special-cased to
//     only matter when it disagrees with re-derivation.
//   - invalid-bucket-faturar: bucket column set to "garbage_value" (not one of
//     the five ordersdomain.OrderBucket constants -- e.g. stale/corrupt data,
//     no CHECK constraint on migration 0089's bucket column). isValidOrderBucket
//     must reject it and fall back to DeriveOrderBucket, same provider_status/
//     tags/faturado_at inputs as legacy-faturar (paid, no tags, not faturado)
//     so it deterministically re-derives to FATURAR. Without the guard, this
//     row's raw "garbage_value" would match no case in the counting switch and
//     be silently dropped from every bucket -- the total across counts would
//     be 3 instead of 4.
func TestGetOrderBucketCountsUsesPersistedBucketWithLegacyFallback(t *testing.T) {
	ctx := context.Background()
	repo, installationID, orderIDs, cleanup := newBucketCountsRepositoryForTest(t, ctx)
	defer cleanup()

	persistedEnviado := bucketCountsFixtureOrder(installationID, orderIDs[0], "paid", ordersdomain.BucketEnviado, nil, nil)
	legacyFaturar := bucketCountsFixtureOrder(installationID, orderIDs[1], "paid", "", nil, nil)
	persistedNovo := bucketCountsFixtureOrder(installationID, orderIDs[2], "created", ordersdomain.BucketNovo, nil, nil)
	invalidBucketFaturar := bucketCountsFixtureOrder(installationID, orderIDs[3], "paid", ordersdomain.OrderBucket("garbage_value"), nil, nil)

	for _, order := range []ordersdomain.MarketplaceOrder{persistedEnviado, legacyFaturar, persistedNovo, invalidBucketFaturar} {
		if _, _, err := repo.UpsertOrders(ctx, []ordersdomain.MarketplaceOrder{order}); err != nil {
			t.Fatalf("seed order %q: %v", order.ProviderOrderID, err)
		}
	}

	// Sanity: the legacy row's bucket column really is empty/NULL (not
	// accidentally "faturar" already), the persisted rows really carry
	// their bucket verbatim, and the invalid-bucket row's column really
	// round-tripped the unrecognized garbage string verbatim (not silently
	// rejected/normalized at write time) -- otherwise this test would not be
	// exercising the two-tier branch it claims to.
	assertPersistedBucketColumn(t, ctx, repo, installationID, orderIDs[0], "enviado")
	assertPersistedBucketColumn(t, ctx, repo, installationID, orderIDs[1], "")
	assertPersistedBucketColumn(t, ctx, repo, installationID, orderIDs[2], "novo")
	assertPersistedBucketColumn(t, ctx, repo, installationID, orderIDs[3], "garbage_value")

	var store ports.OrderBucketStore = repo
	counts, err := store.GetOrderBucketCounts(ctx, installationID)
	if err != nil {
		t.Fatalf("GetOrderBucketCounts() error = %v", err)
	}

	want := ports.OrderBucketCounts{Novo: 1, Faturar: 2, Enviar: 0, Enviado: 1}
	if counts != want {
		t.Fatalf("GetOrderBucketCounts() = %+v, want %+v (persisted-enviado -> Enviado via persisted column, legacy-faturar -> Faturar via fallback, persisted-novo -> Novo, invalid-bucket-faturar -> Faturar via isValidOrderBucket rejecting garbage and falling back)", counts, want)
	}

	// Total across all four counted buckets must equal the number of seeded
	// orders (4). If isValidOrderBucket were removed and the raw
	// "garbage_value" string were switched on directly in GetOrderBucketCounts,
	// it would match no case and invalidBucketFaturar would be silently
	// dropped from every count -- this sum would be 3, not 4.
	total := counts.Novo + counts.Faturar + counts.Enviar + counts.Enviado
	if total != 4 {
		t.Fatalf("GetOrderBucketCounts() total across counted buckets = %d, want 4 (invalid-bucket-faturar must not be silently dropped)", total)
	}

	// Cross-check: re-deriving with the OLD hardcoded shipmentStatus="" call
	// for the persisted-enviado row proves the delta is real and attributable
	// -- it would have landed in faturar, not enviado, without the persisted
	// column this fix now reads.
	oldDerivation := ordersdomain.DeriveOrderBucket(persistedEnviado.ProviderStatus, "", persistedEnviado.Tags, false)
	if oldDerivation != ordersdomain.BucketFaturar {
		t.Fatalf("test premise broken: old shipmentStatus=\"\" derivation for persisted-enviado = %q, want %q (the fixture must actually disagree with re-derivation to prove the fix)", oldDerivation, ordersdomain.BucketFaturar)
	}
}

func assertPersistedBucketColumn(t *testing.T, ctx context.Context, repo *OrderRepository, installationID, orderID, want string) {
	t.Helper()
	var got *string
	err := repo.pool.QueryRow(ctx, `
		SELECT bucket FROM orders_marketplace_orders
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
	`, repo.tenantID, installationID, orderID).Scan(&got)
	if err != nil {
		t.Fatalf("read persisted bucket for %q: %v", orderID, err)
	}
	gotValue := ""
	if got != nil {
		gotValue = *got
	}
	if gotValue != want {
		t.Fatalf("persisted bucket column for %q = %q, want %q", orderID, gotValue, want)
	}
}
