package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
)

func TestOrderRepositoryLowLevelUpsertDoesNotRegressNewerSnapshot(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := newOrderRepositoryForTest(t, ctx)
	defer cleanup()

	newer := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	older := newer.Add(-time.Hour)
	fresh := testOrder(repo, "atomic", "fresh", &newer, "fresh-item", "fresh-payment")
	stale := testOrder(repo, "atomic", "stale", &older, "stale-item", "stale-payment")

	tx, err := repo.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		t.Fatalf("begin fresh transaction: %v", err)
	}
	won, err := repo.upsertOrder(ctx, tx, fresh)
	if err != nil {
		t.Fatalf("upsert fresh snapshot: %v", err)
	}
	if !won {
		t.Fatal("fresh snapshot did not win the upsert")
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit fresh snapshot: %v", err)
	}

	tx, err = repo.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		t.Fatalf("begin stale transaction: %v", err)
	}
	won, err = repo.upsertOrder(ctx, tx, stale)
	if err != nil {
		t.Fatalf("upsert stale snapshot: %v", err)
	}
	if won {
		t.Fatal("older snapshot won the low-level upsert")
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit stale snapshot: %v", err)
	}

	stored := readOrderForTest(t, ctx, repo, fresh.InstallationID, fresh.ProviderOrderID)
	if stored.ProviderStatus != fresh.ProviderStatus {
		t.Fatalf("stored status = %q, want newer snapshot status %q", stored.ProviderStatus, fresh.ProviderStatus)
	}
	if stored.ProviderUpdatedAt == nil || !stored.ProviderUpdatedAt.Equal(newer) {
		t.Fatalf("stored provider updated at = %v, want %v", stored.ProviderUpdatedAt, newer)
	}
}

func TestOrderRepositoryUpsertOrdersAppliesFreshnessAndReplacesChildren(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := newOrderRepositoryForTest(t, ctx)
	defer cleanup()

	base := time.Date(2026, 7, 9, 10, 0, 0, 0, time.UTC)
	newer := base.Add(time.Hour)
	installationID := repo.tenantID + "-installation"
	knownOrderID := repo.tenantID + "-known"

	absent := testOrderForIdentity(installationID, knownOrderID, "absent", &base, "absent-item", "absent-payment")
	assertUpsertCounts(t, ctx, repo, absent, 1, 0)
	assertOrderRows(t, ctx, repo, absent, 1, 1, 1)

	newerSnapshot := testOrderForIdentity(installationID, knownOrderID, "newer", &newer, "newer-item", "newer-payment")
	assertUpsertCounts(t, ctx, repo, newerSnapshot, 1, 0)
	assertOrderRows(t, ctx, repo, newerSnapshot, 1, 1, 1)

	equalSnapshot := testOrderForIdentity(installationID, knownOrderID, "equal", &newer, "equal-item", "equal-payment")
	assertUpsertCounts(t, ctx, repo, equalSnapshot, 1, 0)
	assertOrderRows(t, ctx, repo, equalSnapshot, 1, 1, 1)

	olderSnapshot := testOrderForIdentity(installationID, knownOrderID, "older", &base, "older-item", "older-payment")
	assertUpsertCounts(t, ctx, repo, olderSnapshot, 0, 1)
	assertOrderRows(t, ctx, repo, equalSnapshot, 1, 1, 1)

	unknownIncoming := testOrderForIdentity(installationID, knownOrderID, "unknown", nil, "unknown-item", "unknown-payment")
	assertUpsertCounts(t, ctx, repo, unknownIncoming, 0, 1)
	assertOrderRows(t, ctx, repo, equalSnapshot, 1, 1, 1)

	unknownOrderID := repo.tenantID + "-unknown"
	unknownAbsent := testOrderForIdentity(installationID, unknownOrderID, "unknown-absent", nil, "unknown-absent-item", "unknown-absent-payment")
	assertUpsertCounts(t, ctx, repo, unknownAbsent, 1, 0)
	unknownReplacement := testOrderForIdentity(installationID, unknownOrderID, "unknown-replacement", nil, "unknown-replacement-item", "unknown-replacement-payment")
	assertUpsertCounts(t, ctx, repo, unknownReplacement, 1, 0)
	assertOrderRows(t, ctx, repo, unknownReplacement, 1, 1, 1)
}

func newOrderRepositoryForTest(t *testing.T, ctx context.Context) (*OrderRepository, func()) {
	t.Helper()
	pgURL := os.Getenv("MC_DATABASE_URL")
	if pgURL == "" {
		t.Skip("set MC_DATABASE_URL to run Postgres order repository tests")
	}

	pool, err := pgxpool.New(ctx, pgURL)
	if err != nil {
		t.Fatalf("open Postgres pool: %v", err)
	}
	tenantID := fmt.Sprintf("f01-order-repo-%d", time.Now().UnixNano())
	repo := NewOrderRepository(pool, tenantID)
	installationID := tenantID + "-installation"
	orderIDs := []string{tenantID + "-atomic", tenantID + "-known", tenantID + "-unknown"}
	return repo, func() {
		defer pool.Close()
		for _, orderID := range orderIDs {
			if _, err := pool.Exec(ctx, `DELETE FROM orders_marketplace_order_items WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3`, tenantID, installationID, orderID); err != nil {
				t.Errorf("cleanup order items %q: %v", orderID, err)
			}
			if _, err := pool.Exec(ctx, `DELETE FROM orders_marketplace_order_payments WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3`, tenantID, installationID, orderID); err != nil {
				t.Errorf("cleanup order payments %q: %v", orderID, err)
			}
			if _, err := pool.Exec(ctx, `DELETE FROM orders_marketplace_orders WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3`, tenantID, installationID, orderID); err != nil {
				t.Errorf("cleanup order %q: %v", orderID, err)
			}
		}
	}
}

func testOrder(repo *OrderRepository, suffix, status string, providerUpdatedAt *time.Time, itemID, paymentID string) ordersdomain.MarketplaceOrder {
	installationID := repo.tenantID + "-installation"
	orderID := repo.tenantID + "-" + suffix
	return testOrderForIdentity(installationID, orderID, status, providerUpdatedAt, itemID, paymentID)
}

func testOrderForIdentity(installationID, orderID, status string, providerUpdatedAt *time.Time, itemID, paymentID string) ordersdomain.MarketplaceOrder {
	now := time.Date(2026, 7, 9, 13, 0, 0, 0, time.UTC)
	return ordersdomain.MarketplaceOrder{
		InstallationID:    installationID,
		ProviderCode:      "mercado_livre",
		ProviderOrderID:   orderID,
		ProviderStatus:    status,
		ProviderUpdatedAt: providerUpdatedAt,
		FetchedAt:         now,
		RawProviderRef:    "/orders/" + orderID,
		CreatedAt:         now,
		UpdatedAt:         now,
		Items: []ordersdomain.MarketplaceOrderItem{{
			ProviderItemID: itemID,
			Quantity:       1,
			LinkQuality:    ordersdomain.LinkQualityMissing,
		}},
		Payments: []ordersdomain.MarketplaceOrderPayment{{
			ProviderPaymentID: paymentID,
			ProviderStatus:    "approved",
		}},
	}
}

func assertUpsertCounts(t *testing.T, ctx context.Context, repo *OrderRepository, order ordersdomain.MarketplaceOrder, wantImported, wantSkipped int) {
	t.Helper()
	imported, skipped, err := repo.UpsertOrders(ctx, []ordersdomain.MarketplaceOrder{order})
	if err != nil {
		t.Fatalf("upsert order %q: %v", order.ProviderOrderID, err)
	}
	if imported != wantImported || skipped != wantSkipped {
		t.Fatalf("upsert counts = (%d imported, %d skipped), want (%d imported, %d skipped)", imported, skipped, wantImported, wantSkipped)
	}
}

func assertOrderRows(t *testing.T, ctx context.Context, repo *OrderRepository, want ordersdomain.MarketplaceOrder, wantOrders, wantItems, wantPayments int) {
	t.Helper()
	stored := readOrderForTest(t, ctx, repo, want.InstallationID, want.ProviderOrderID)
	if stored.ProviderStatus != want.ProviderStatus {
		t.Fatalf("stored status = %q, want %q", stored.ProviderStatus, want.ProviderStatus)
	}
	if !sameTime(stored.ProviderUpdatedAt, want.ProviderUpdatedAt) {
		t.Fatalf("stored provider updated at = %v, want %v", stored.ProviderUpdatedAt, want.ProviderUpdatedAt)
	}
	if stored.RawProviderRef != want.RawProviderRef {
		t.Fatalf("stored raw provider reference = %q, want safe reference %q", stored.RawProviderRef, want.RawProviderRef)
	}
	if len(stored.Items) != len(want.Items) || stored.Items[0].ProviderItemID != want.Items[0].ProviderItemID {
		t.Fatalf("stored items = %#v, want %#v", stored.Items, want.Items)
	}
	if len(stored.Payments) != len(want.Payments) || stored.Payments[0].ProviderPaymentID != want.Payments[0].ProviderPaymentID {
		t.Fatalf("stored payments = %#v, want %#v", stored.Payments, want.Payments)
	}

	var orders, items, payments int
	err := repo.pool.QueryRow(ctx, `
		SELECT
			(SELECT count(*) FROM orders_marketplace_orders WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3),
			(SELECT count(*) FROM orders_marketplace_order_items WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3),
			(SELECT count(*) FROM orders_marketplace_order_payments WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3)
	`, repo.tenantID, want.InstallationID, want.ProviderOrderID).Scan(&orders, &items, &payments)
	if err != nil {
		t.Fatalf("count persisted rows: %v", err)
	}
	if orders != wantOrders || items != wantItems || payments != wantPayments {
		t.Fatalf("persisted row counts = (%d orders, %d items, %d payments), want (%d, %d, %d)", orders, items, payments, wantOrders, wantItems, wantPayments)
	}
}

func readOrderForTest(t *testing.T, ctx context.Context, repo *OrderRepository, installationID, orderID string) ordersdomain.MarketplaceOrder {
	t.Helper()
	orders, err := repo.ListOrders(ctx, installationID, 10)
	if err != nil {
		t.Fatalf("list orders: %v", err)
	}
	for _, order := range orders {
		if order.ProviderOrderID == orderID {
			return order
		}
	}
	t.Fatalf("order %q not found", orderID)
	return ordersdomain.MarketplaceOrder{}
}

func sameTime(left, right *time.Time) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Equal(*right)
}
