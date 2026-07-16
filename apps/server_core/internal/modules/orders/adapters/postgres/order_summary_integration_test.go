//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	orderspostgres "marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres"
	ordersapplication "marketplace-central/apps/server_core/internal/modules/orders/application"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestOrderSummaryCountsTodayAndSevenDaysTenantScoped(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "orders_summary_scope")
	token := fmt.Sprint(time.Now().UnixNano())
	tenant, otherTenant := "summary-tenant-a-"+token, "summary-tenant-b-"+token
	installation, otherInstallation := "summary-installation-a-"+token, "summary-installation-b-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM orders_marketplace_orders WHERE tenant_id = ANY($1)`, []string{tenant, otherTenant})
	})

	referenceTime := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	seedOrderSummaryRow(t, pool, tenant, installation, "today", timestampPtr(referenceTime.Add(-time.Hour)))
	seedOrderSummaryRow(t, pool, tenant, installation, "seven-days", timestampPtr(referenceTime.AddDate(0, 0, -6)))
	seedOrderSummaryRow(t, pool, tenant, installation, "older", timestampPtr(referenceTime.AddDate(0, 0, -7)))
	seedOrderSummaryRow(t, pool, tenant, installation, "unknown", nil)
	seedOrderSummaryRow(t, pool, otherTenant, installation, "other-tenant", timestampPtr(referenceTime.Add(-time.Hour)))
	seedOrderSummaryRow(t, pool, tenant, otherInstallation, "other-installation", timestampPtr(referenceTime.Add(-time.Hour)))

	service := ordersapplication.NewSummaryService(orderspostgres.NewOrderRepository(pool, tenant))
	got, err := service.Summary(ctx, installation, referenceTime)
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.Today != 1 || got.SevenDays != 2 {
		t.Fatalf("summary = %+v, want today=1 seven_days=2", got)
	}
}

func TestOrderSummaryDoesNotTreatUnknownTimestampAsToday(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "orders_summary_unknown_timestamp")
	token := fmt.Sprint(time.Now().UnixNano())
	tenant := "summary-unknown-" + token
	installation := "summary-installation-" + token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM orders_marketplace_orders WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})

	seedOrderSummaryRow(t, pool, tenant, installation, "unknown", nil)
	referenceTime := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	service := ordersapplication.NewSummaryService(orderspostgres.NewOrderRepository(pool, tenant))

	got, err := service.Summary(ctx, installation, referenceTime)
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.Today != 0 || got.SevenDays != 0 {
		t.Fatalf("summary = %+v, want unknown timestamp excluded from both windows", got)
	}
}

func seedOrderSummaryRow(t *testing.T, pool *pgxpool.Pool, tenant, installation, orderID string, providerCreatedAt *time.Time) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `
		INSERT INTO orders_marketplace_orders(
			tenant_id, installation_id, provider_code, provider_order_id,
			provider_status, provider_created_at, fetched_at, created_at, updated_at
		) VALUES($1, $2, 'mercado_livre', $3, 'paid', $4, now(), now(), now())
	`, tenant, installation, orderID, providerCreatedAt)
	if err != nil {
		t.Fatalf("seed order %q: %v", orderID, err)
	}
}

func timestampPtr(value time.Time) *time.Time {
	return &value
}
