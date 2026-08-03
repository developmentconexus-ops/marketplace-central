//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	orderspostgres "marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestReingestSameOrderKeepsOneRow ingere o mesmo provider_order_id duas vezes
// contra Postgres real (hermetic lane, sem fake/mocks) e conta as linhas. É o
// pré-requisito da janela com sobreposição de 5 minutos do scheduler de
// pedidos (F-00): sem upsert real, essa sobreposição duplicaria o mesmo
// pedido a cada reconsulta. order_repo.go declara ON CONFLICT (tenant_id,
// installation_id, provider_order_id) DO UPDATE — este teste PROVA o
// comportamento, não apenas lê a SQL.
func TestReingestSameOrderKeepsOneRow(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "orders_reingest")

	tenantID := fmt.Sprintf("f00-reingest-%d", time.Now().UnixNano())
	installationID := tenantID + "-installation"
	providerOrderID := tenantID + "-order-1"

	repo := orderspostgres.NewOrderRepository(pool, tenantID)

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		for _, stmt := range []string{
			`DELETE FROM orders_marketplace_order_items WHERE tenant_id = $1`,
			`DELETE FROM orders_marketplace_order_payments WHERE tenant_id = $1`,
			`DELETE FROM orders_marketplace_orders WHERE tenant_id = $1`,
		} {
			if _, err := pool.Exec(cleanupCtx, stmt, tenantID); err != nil {
				t.Errorf("cleanup %q: %v", stmt, err)
			}
		}
	})

	baseTime := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	firstUpdatedAt := baseTime
	first := ordersdomain.MarketplaceOrder{
		InstallationID:    installationID,
		ProviderCode:      "mercado_livre",
		ProviderOrderID:   providerOrderID,
		ProviderStatus:    "paid",
		ProviderUpdatedAt: &firstUpdatedAt,
		FetchedAt:         baseTime,
		CreatedAt:         baseTime,
		UpdatedAt:         baseTime,
	}

	// Step 1: ingest via the real production write path (OrderRepository.UpsertOrders,
	// the same method the orders sync job calls), exactly like the first pass of a
	// scheduler cursor window would.
	imported, skipped, err := repo.UpsertOrders(ctx, []ordersdomain.MarketplaceOrder{first})
	if err != nil {
		t.Fatalf("first ingest: %v", err)
	}
	if imported != 1 || skipped != 0 {
		t.Fatalf("first ingest imported=%d skipped=%d, want imported=1 skipped=0", imported, skipped)
	}

	// Step 2: reingest the SAME provider_order_id, simulating the scheduler's 5-minute
	// overlap window re-fetching an order it already saw. A real re-fetch carries a
	// changed provider_status and an advanced provider_updated_at (the freshness gate
	// in upsertOrder's WHERE clause requires provider_updated_at to not regress).
	secondUpdatedAt := baseTime.Add(2 * time.Minute)
	second := first
	second.ProviderStatus = "cancelled"
	second.ProviderUpdatedAt = &secondUpdatedAt
	second.UpdatedAt = secondUpdatedAt

	imported, skipped, err = repo.UpsertOrders(ctx, []ordersdomain.MarketplaceOrder{second})
	if err != nil {
		t.Fatalf("second ingest (reingest): %v", err)
	}
	if imported != 1 || skipped != 0 {
		t.Fatalf("second ingest imported=%d skipped=%d, want imported=1 skipped=0 (fresher snapshot must win the upsert)", imported, skipped)
	}

	// Assertion 1 (presence): reingesting the same order must not duplicate the row.
	var rowCount int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM orders_marketplace_orders
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
	`, tenantID, installationID, providerOrderID).Scan(&rowCount); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	if rowCount != 1 {
		t.Fatalf("row count = %d, want 1 (reingest must not duplicate the order row)", rowCount)
	}

	// Assertion 2 (value, not just presence): the stored provider_status must be the
	// SECOND write's value. A DO NOTHING conflict clause would also leave rowCount == 1
	// above, but would silently keep the FIRST write's "paid" forever — this is the
	// assertion that actually distinguishes DO UPDATE from DO NOTHING.
	var storedStatus string
	if err := pool.QueryRow(ctx, `
		SELECT provider_status FROM orders_marketplace_orders
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
	`, tenantID, installationID, providerOrderID).Scan(&storedStatus); err != nil {
		t.Fatalf("read provider_status: %v", err)
	}
	if storedStatus != second.ProviderStatus {
		t.Fatalf("provider_status = %q, want %q (the second write's value) — a DO NOTHING conflict clause would pass the row-count check above while leaving the FIRST write's value in place", storedStatus, second.ProviderStatus)
	}
}
