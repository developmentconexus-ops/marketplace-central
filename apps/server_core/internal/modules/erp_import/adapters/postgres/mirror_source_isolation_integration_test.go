//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// dropProtocols removes protocols by ID. erp_import_protocols is keyed on id
// alone, so a fixture UUID is global across tenants and packages: without this
// a re-run (or a concurrent package) collides on the primary key even though
// the tenant differs.
func dropProtocols(t *testing.T, pool *pgxpool.Pool, ids ...string) {
	t.Helper()
	cleanup := func() {
		for _, id := range ids {
			_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE id=$1`, id)
		}
	}
	cleanup()
	t.Cleanup(cleanup)
}

// Two sources carrying the SAME product code must own SEPARATE mirror rows.
// Before 0078 the mirror was keyed on (tenant, codigo_produto) alone, so the
// second import hijacked the first source's row and flipped its source — the
// live D-121 failure where switching the active source showed the wrong data.
func TestMirrorKeepsOneRowPerSourceForTheSameProduct(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	dropProtocols(t, pool, "a1111111-1111-4111-8111-111111111101", "a1111111-1111-4111-8111-111111111102")
	now := time.Now().UTC().Truncate(time.Microsecond)

	erp := domain.ImportSnapshot{
		ID: "a1111111-1111-4111-8111-111111111101", Protocol: "#201-E", FileSHA256: "iso-erp",
		Source: domain.SourceXLSX, ImportedAt: now, Status: domain.ImportStatusCompleted,
		AcceptedRows: []domain.NormalizedRow{
			{Codprod: "SHARED", Descrprod: "ERP DESCRIPTION", Custo: "10.50", PrecoVenda: "19.90", StockPhysical: "7"},
		},
	}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, erp); err != nil {
		t.Fatal(err)
	}

	catalog := domain.ImportSnapshot{
		ID: "a1111111-1111-4111-8111-111111111102", Protocol: "#202-E", FileSHA256: "iso-catalog",
		Source: domain.SourceCatalogoCliente, ImportedAt: now.Add(time.Second), Status: domain.ImportStatusCompleted,
		AcceptedRows: []domain.NormalizedRow{
			{Codprod: "SHARED", Descrprod: "CATALOG DESCRIPTION"},
		},
	}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, catalog); err != nil {
		t.Fatal(err)
	}

	var rows int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM products_mirror WHERE tenant_id=$1 AND codigo_produto='SHARED'`, tenant).Scan(&rows); err != nil {
		t.Fatal(err)
	}
	if rows != 2 {
		t.Fatalf("mirror rows for SHARED=%d want 2 (one per source)", rows)
	}

	// The ERP row keeps its own values: the catalog import must not overwrite
	// cost/price/stock with the honest-unknown NULLs of the other source.
	var description string
	var custo, preco, estoque *string
	if err := pool.QueryRow(ctx, `SELECT descricao,custo::text,preco_venda::text,estoque_total::text
		FROM products_mirror WHERE tenant_id=$1 AND source=$2 AND codigo_produto='SHARED'`,
		tenant, domain.SourceXLSX).Scan(&description, &custo, &preco, &estoque); err != nil {
		t.Fatal(err)
	}
	if description != "ERP DESCRIPTION" {
		t.Fatalf("erp description=%q want %q", description, "ERP DESCRIPTION")
	}
	if custo == nil || *custo != "10.50" {
		t.Fatalf("erp custo=%v want 10.50", custo)
	}
	if preco == nil || *preco != "19.90" {
		t.Fatalf("erp preco_venda=%v want 19.90 (sale price must survive the other source's import)", preco)
	}
	// Physical 7 with no reported reservation is 7 sellable, not unknown.
	if estoque == nil || *estoque != "7" {
		t.Fatalf("erp estoque_total=%v want 7 (physical with no reservation reported)", estoque)
	}

	// And the ERP row must NOT have been flagged absent by the other source's
	// snapshot: keep-absent is scoped per source.
	var absent bool
	if err := pool.QueryRow(ctx, `SELECT absent_in_last_snapshot FROM products_mirror
		WHERE tenant_id=$1 AND source=$2 AND codigo_produto='SHARED'`, tenant, domain.SourceXLSX).Scan(&absent); err != nil {
		t.Fatal(err)
	}
	if absent {
		t.Fatal("erp row flagged absent by an unrelated source's import")
	}
}

// A mirror rebuild from the stored protocol history must restore the sale price.
// Before 0078 erp_import_products had no preco_venda column, so a rebuild could
// only ever restore NULL — the price silently disappeared on every resync.
func TestRebuildFromProtocolRestoresSalePrice(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	dropProtocols(t, pool, "a1111111-1111-4111-8111-111111111103")
	now := time.Now().UTC().Truncate(time.Microsecond)

	snapshot := domain.ImportSnapshot{
		ID: "a1111111-1111-4111-8111-111111111103", Protocol: "#203-E", FileSHA256: "iso-rebuild",
		Source: domain.SourceXLSX, ImportedAt: now, Status: domain.ImportStatusCompleted,
		AcceptedRows: []domain.NormalizedRow{
			{Codprod: "REBUILD", Descrprod: "PRICED PRODUCT", Custo: "42.10", PrecoVenda: "99.90", StockPhysical: "3"},
		},
	}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, snapshot); err != nil {
		t.Fatal(err)
	}

	// Wipe the mirror row so only the protocol history can restore it.
	if _, err := pool.Exec(ctx, `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.SyncLatestCompletedSnapshot(ctx, tenant, domain.SourceXLSX); err != nil {
		t.Fatal(err)
	}

	var preco *string
	if err := pool.QueryRow(ctx, `SELECT preco_venda::text FROM products_mirror
		WHERE tenant_id=$1 AND source=$2 AND codigo_produto='REBUILD'`, tenant, domain.SourceXLSX).Scan(&preco); err != nil {
		t.Fatal(err)
	}
	if preco == nil || *preco != "99.90" {
		t.Fatalf("rebuilt preco_venda=%v want 99.90", preco)
	}
}
