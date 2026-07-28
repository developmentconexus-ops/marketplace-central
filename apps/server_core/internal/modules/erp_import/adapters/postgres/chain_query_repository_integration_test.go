//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestGetImportChainCountsCurrentQueueAcrossInstallations(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("61000000-0000-0000-0000-000000000001")

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id=$1`, tenant)
	})

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES ($1, $2, 'chain-counts-hash', '#610-E', 'xlsx', now(), 'COMPLETED', 99, 0, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES
			($2, $1, '101', 'Product 101', 1, 0),
			($2, $1, '102', 'Product 102', 1, 0),
			($2, $1, '103', 'Product 103', 1, 0),
			($2, $1, '104', 'Product 104', 1, 0);
	`, importID, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES
			($1, 'installation-a', 'mercadolivre', 'item-101', 'resolved', 101),
			($1, 'installation-a', 'mercadolivre', 'item-102', 'resolved', 102),
			($1, 'installation-a', 'mercadolivre', 'item-103', 'pending', 103);
	`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
		VALUES
			($1, 'installation-a', 'market', '{"pending":["101","102","OUTSIDE"]}'::jsonb),
			($1, 'installation-b', 'market', '{"pending":["101","103"]}'::jsonb)
	`, tenant); err != nil {
		t.Fatal(err)
	}

	before := time.Now()
	got, err := repo.GetImportChain(ctx, tenant, importID)
	after := time.Now()
	if err != nil {
		t.Fatal(err)
	}
	if got.Protocol != "#610-E" || got.Importados != 4 || got.Vinculados != 2 || got.Enfileirados != 3 {
		t.Fatalf("chain=%#v", got)
	}
	if got.QueueReadAt.Before(before) || got.QueueReadAt.After(after) {
		t.Fatalf("queue_read_at=%s outside [%s,%s]", got.QueueReadAt, before, after)
	}
}

func TestGetImportChainMissingAndProtocolWithoutSyncState(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	importID := domain.ImportID("62000000-0000-0000-0000-000000000001")
	otherImportID := domain.ImportID("62000000-0000-0000-0000-000000000002")
	otherTenant := tenant + "-other"

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id=$1`, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_products WHERE tenant_id IN ($1, $2)`, tenant, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE tenant_id=$1`, otherTenant)
	})

	if _, err := repo.GetImportChain(ctx, tenant, "62000000-0000-0000-0000-000000000099"); !errors.Is(err, ports.ErrImportNotFound) {
		t.Fatalf("missing err=%v", err)
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_protocols
			(id, tenant_id, file_sha256, protocol, source, imported_at, status, accepted_count, rejected_count, warning_count)
		VALUES
			($1, $2, 'chain-empty-hash', '#620-E', 'xlsx', now(), 'COMPLETED', 1, 0, 0),
			($4, $3, 'chain-other-hash', '#620-E', 'xlsx', now(), 'COMPLETED', 1, 0, 0);
	`, importID, tenant, otherTenant, otherImportID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO erp_import_products (tenant_id, protocol_id, codprod, descrprod, custo, stock_physical)
		VALUES
			($2, $1, '201', 'Main product', 1, 0),
			($3, $4, '201', 'Other product', 1, 0);
	`, importID, tenant, otherTenant, otherImportID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links
			(tenant_id, installation_id, provider_code, provider_item_id, state, internal_product_id)
		VALUES ($1, 'installation-other', 'mercadolivre', 'item-other-201', 'resolved', 201);
	`, otherTenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO sync_state (tenant_id, installation_id, entity, cursor)
		VALUES ($1, 'installation-other', 'market', '{"pending":["201"]}'::jsonb)
	`, otherTenant); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetImportChain(ctx, tenant, importID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Protocol != "#620-E" || got.Importados != 1 || got.Vinculados != 0 || got.Enfileirados != 0 {
		t.Fatalf("chain=%#v", got)
	}
}
