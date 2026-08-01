//go:build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	syncpostgres "marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// F-01/M-09: ReadAll must return the sentinel-installation "erp" row (that is
// where the products entity lives) alongside rows under other installation
// ids for the same tenant — no installation_id filter at all.
func TestHealthReaderReadsAcrossEveryInstallation(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "sync_health_reader_installs")
	tenant := "tenant-health-installs-" + newToken(t)
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant) })

	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, last_full_sync_at, consecutive_failures)
		VALUES ($1, 'erp', 'products', $2, 0)
	`, tenant, time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, last_full_sync_at, consecutive_failures)
		VALUES ($1, 'ml-installation-1', 'listings', $2, 0)
	`, tenant, time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatal(err)
	}

	reader := syncpostgres.NewHealthReader(pool)
	rows, err := reader.ReadAll(ctx, tenant)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2 (one per installation, no installation_id filter): %+v", len(rows), rows)
	}
	entities := map[string]bool{}
	for _, row := range rows {
		entities[row.Entity] = true
	}
	if !entities["products"] || !entities["listings"] {
		t.Fatalf("entities = %+v, want both products (erp sentinel) and listings", entities)
	}
}

// Ordering, cursor pass-through, and last_error mapping against real
// Postgres — the same scanning idiom sync_state_repo.go uses, widened to a
// tenant-scoped multi-row read.
func TestHealthReaderOrderingAndFieldMapping(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "sync_health_reader_fields")
	tenant := "tenant-health-fields-" + newToken(t)
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant) })

	failAt := time.Date(2026, 7, 30, 8, 0, 0, 0, time.UTC)
	errPayload := fmt.Sprintf(`{"message":"provider timeout","at":%q}`, failAt.Format(time.RFC3339))

	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor, last_error, consecutive_failures)
		VALUES ($1, 'erp', 'orders', '{"phase":"sweep"}'::jsonb, $2::jsonb, 4)
	`, tenant, errPayload); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, consecutive_failures)
		VALUES ($1, 'erp', 'listings', 0)
	`, tenant); err != nil {
		t.Fatal(err)
	}

	reader := syncpostgres.NewHealthReader(pool)
	rows, err := reader.ReadAll(ctx, tenant)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}
	// entity ascending: listings < orders
	if rows[0].Entity != "listings" || rows[1].Entity != "orders" {
		t.Fatalf("order = [%s, %s], want [listings, orders]", rows[0].Entity, rows[1].Entity)
	}

	listings, orders := rows[0], rows[1]
	if listings.Cursor != nil {
		t.Errorf("listings.Cursor = %q, want nil (NULL column)", listings.Cursor)
	}
	if listings.LastError != nil {
		t.Errorf("listings.LastError = %+v, want nil", listings.LastError)
	}

	var cursor struct {
		Phase string `json:"phase"`
	}
	if err := json.Unmarshal(orders.Cursor, &cursor); err != nil {
		t.Fatalf("orders cursor %q not valid JSON: %v", orders.Cursor, err)
	}
	if cursor.Phase != "sweep" {
		t.Errorf("orders cursor phase = %q, want sweep", cursor.Phase)
	}
	if orders.LastError == nil || orders.LastError.Message != "provider timeout" {
		t.Fatalf("orders.LastError = %+v, want message=provider timeout", orders.LastError)
	}
	if !orders.LastError.At.Equal(failAt) {
		t.Errorf("orders.LastError.At = %v, want %v", orders.LastError.At, failAt)
	}
	if orders.ConsecutiveFailures != 4 {
		t.Errorf("orders.ConsecutiveFailures = %d, want 4", orders.ConsecutiveFailures)
	}
}

// An honest empty tenant: zero rows, zero error — never fabricated.
func TestHealthReaderEmptyTenant(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "sync_health_reader_empty")
	tenant := "tenant-health-empty-" + newToken(t)

	reader := syncpostgres.NewHealthReader(pool)
	rows, err := reader.ReadAll(ctx, tenant)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if rows == nil {
		t.Fatal("rows is nil, want a non-nil empty slice")
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %d, want 0", len(rows))
	}
}
