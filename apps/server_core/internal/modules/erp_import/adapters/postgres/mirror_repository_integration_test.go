//go:build integration

package postgres_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestMirrorMergeKeepsAbsentIsTenantScopedAndResurrects(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	otherTenant := tenant + "-other"
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_links WHERE tenant_id IN ($1,$2)`, tenant, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror_stock_locations WHERE tenant_id IN ($1,$2)`, tenant, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id IN ($1,$2)`, tenant, otherTenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE tenant_id=$1`, otherTenant)
	})

	now := time.Now().UTC()
	x := mirrorSnapshot("61111111-1111-1111-1111-111111111111", "mirror-x-1", now, []domain.NormalizedRow{{Codprod: "X", Descrprod: "X one", StockPhysical: "8", StockReserved: stringPtr("3")}})
	if err := repo.PersistSnapshotAtomically(ctx, tenant, x); err != nil {
		t.Fatal(err)
	}
	other := mirrorSnapshot("62222222-2222-2222-2222-222222222222", "mirror-other", now, []domain.NormalizedRow{{Codprod: "X", Descrprod: "Other X", Custo: "9", StockPhysical: "2", StockReserved: stringPtr("1")}})
	if err := repo.PersistSnapshotAtomically(ctx, otherTenant, other); err != nil {
		t.Fatal(err)
	}
	_, err := pool.Exec(ctx, `INSERT INTO product_links (tenant_id,installation_id,provider_code,provider_item_id,state,internal_product_id) VALUES ($1,'installation','mercadolivre','item','linked',42)`, tenant)
	if err != nil {
		t.Fatal(err)
	}

	absent := mirrorSnapshot("63333333-3333-3333-3333-333333333333", "mirror-x-2", now.Add(time.Second), nil)
	if err := repo.PersistSnapshotAtomically(ctx, tenant, absent); err != nil {
		t.Fatal(err)
	}
	var flagged bool
	var stale *time.Time
	var cost, stock *string
	if err := pool.QueryRow(ctx, `SELECT absent_in_last_snapshot,stale_since,custo::text,estoque_total::text FROM products_mirror WHERE tenant_id=$1 AND codigo_produto='X'`, tenant).Scan(&flagged, &stale, &cost, &stock); err != nil {
		t.Fatal(err)
	}
	if !flagged || stale == nil {
		t.Fatalf("absent state flagged=%v stale=%v", flagged, stale)
	}
	if cost != nil || stock == nil || *stock != "5" {
		t.Fatalf("honest numerics cost=%v stock=%v", cost, stock)
	}
	var links int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM product_links WHERE tenant_id=$1 AND internal_product_id=42`, tenant).Scan(&links); err != nil || links != 1 {
		t.Fatalf("links=%d err=%v", links, err)
	}
	var otherFlagged bool
	if err := pool.QueryRow(ctx, `SELECT absent_in_last_snapshot FROM products_mirror WHERE tenant_id=$1 AND codigo_produto='X'`, otherTenant).Scan(&otherFlagged); err != nil || otherFlagged {
		t.Fatalf("other tenant flagged=%v err=%v", otherFlagged, err)
	}

	resurrected := mirrorSnapshot("64444444-4444-4444-4444-444444444444", "mirror-x-3", now.Add(2*time.Second), []domain.NormalizedRow{{Codprod: "X", Descrprod: "X again"}})
	if err := repo.PersistSnapshotAtomically(ctx, tenant, resurrected); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT absent_in_last_snapshot,stale_since,custo::text,estoque_total::text FROM products_mirror WHERE tenant_id=$1 AND codigo_produto='X'`, tenant).Scan(&flagged, &stale, &cost, &stock); err != nil {
		t.Fatal(err)
	}
	if flagged || stale != nil || cost != nil || stock != nil {
		t.Fatalf("resurrected flagged=%v stale=%v cost=%v stock=%v", flagged, stale, cost, stock)
	}
}

func TestMirrorMergeRollbackAndAcceptedCount(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror_stock_locations WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
	})

	bad := mirrorSnapshot("65555555-5555-5555-5555-555555555555", "mirror-bad", time.Now().UTC(), []domain.NormalizedRow{{Codprod: "BAD", Descrprod: "Bad", Custo: "1", PrecoVenda: "not-numeric", StockPhysical: "1"}})
	if err := repo.PersistSnapshotAtomically(ctx, tenant, bad); err == nil {
		t.Fatal("invalid mirror numeric succeeded")
	}
	var protocols, mirrors int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM erp_import_protocols WHERE tenant_id=$1 AND id=$2`, tenant, bad.ID).Scan(&protocols); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM products_mirror WHERE tenant_id=$1`, tenant).Scan(&mirrors); err != nil {
		t.Fatal(err)
	}
	if protocols != 0 || mirrors != 0 {
		t.Fatalf("rollback protocols=%d mirrors=%d", protocols, mirrors)
	}

	rejectedColumn := "codprod"
	accepted := mirrorSnapshot("66666666-6666-6666-6666-666666666666", "mirror-accepted", time.Now().UTC().Add(time.Second), []domain.NormalizedRow{{Codprod: "OK", Descrprod: "Accepted", StockPhysical: "1", StockReserved: stringPtr("0")}})
	accepted.Issues = []domain.Issue{{Row: 3, Column: &rejectedColumn, Kind: domain.Rejection, Code: domain.CodeEmptyCodprod}}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, accepted); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM products_mirror WHERE tenant_id=$1`, tenant).Scan(&mirrors); err != nil || mirrors != len(accepted.AcceptedRows) {
		t.Fatalf("mirror count=%d accepted=%d err=%v", mirrors, len(accepted.AcceptedRows), err)
	}
}

func TestMirrorMergeReplacesTouchedLocationsAndRetainsAbsentLocations(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror_stock_locations WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
	})
	now := time.Now().UTC()
	first := mirrorSnapshot("67777777-7777-7777-7777-777777777777", "mirror-loc-1", now, []domain.NormalizedRow{
		{Codprod: "A", Descrprod: "A", StockPhysical: "7", StockReserved: stringPtr("2"), Local: stringPtr("OLD")},
		{Codprod: "B", Descrprod: "B", StockPhysical: "4", StockReserved: stringPtr("1"), Local: stringPtr("KEEP")},
	})
	if err := repo.PersistSnapshotAtomically(ctx, tenant, first); err != nil {
		t.Fatal(err)
	}
	second := mirrorSnapshot("68888888-8888-8888-8888-888888888888", "mirror-loc-2", now.Add(time.Second), []domain.NormalizedRow{{Codprod: "A", Descrprod: "A", StockPhysical: "9", StockReserved: stringPtr("3"), Local: stringPtr("NEW")}})
	if err := repo.PersistSnapshotAtomically(ctx, tenant, second); err != nil {
		t.Fatal(err)
	}
	rows, err := pool.Query(ctx, `SELECT codigo_produto,local_codigo,local_descricao,quantidade::text FROM products_mirror_stock_locations WHERE tenant_id=$1 ORDER BY codigo_produto,local_codigo`, tenant)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	type location struct {
		product, code         string
		description, quantity *string
	}
	var got []location
	for rows.Next() {
		var v location
		if err := rows.Scan(&v.product, &v.code, &v.description, &v.quantity); err != nil {
			t.Fatal(err)
		}
		got = append(got, v)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].product != "A" || got[0].code != "NEW" || got[0].description != nil || got[0].quantity == nil || *got[0].quantity != "6" || got[1].product != "B" || got[1].code != "KEEP" {
		t.Fatalf("locations=%#v", got)
	}
}

func TestSyncLatestCompletedSnapshotUsesSharedMerge(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror_stock_locations WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
	})
	snapshot := mirrorSnapshot("69999999-9999-9999-9999-999999999999", "mirror-sync", time.Now().UTC(), []domain.NormalizedRow{{Codprod: "SYNC", Descrprod: "Sync", Custo: "3", StockPhysical: "5", StockReserved: stringPtr("2")}})
	if err := repo.PersistSnapshotAtomically(ctx, tenant, snapshot); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant); err != nil {
		t.Fatal(err)
	}
	processed, err := repo.SyncLatestCompletedSnapshot(ctx, tenant, domain.SourceXLSX)
	if err != nil || processed != 1 {
		t.Fatalf("processed=%d err=%v", processed, err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM products_mirror WHERE tenant_id=$1 AND codigo_produto='SYNC' AND protocol_id=$2`, tenant, snapshot.ID).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
}

func TestMirrorCurrentReadsAreTenantAndSourceScoped(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	otherTenant := tenant + "-other"
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id IN ($1,$2)`, tenant, otherTenant)
	})

	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "1", "xlsx one", stringPtr("shared"), false, nil, nil)
	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "2", "xlsx two", stringPtr("shared"), false, nil, nil)
	seedMirrorProduct(t, pool, tenant, domain.SourceCatalogoCliente, "3", "catalog three", stringPtr("shared"), false, nil, nil)
	seedMirrorProduct(t, pool, otherTenant, domain.SourceXLSX, "4", "other four", stringPtr("shared"), false, nil, nil)

	rows, err := repo.MirrorRows(ctx, tenant, domain.SourceXLSX)
	if err != nil || !reflect.DeepEqual(mirrorCodes(rows), []string{"1", "2"}) {
		t.Fatalf("rows codes=%v err=%v", mirrorCodes(rows), err)
	}
	product, found, err := repo.MirrorProductByCode(ctx, tenant, domain.SourceXLSX, "1")
	if err != nil || !found || product.CodigoProduto != "1" {
		t.Fatalf("product=%+v found=%v err=%v", product, found, err)
	}
	if _, found, err := repo.MirrorProductByCode(ctx, tenant, domain.SourceCatalogoCliente, "1"); err != nil || found {
		t.Fatalf("cross-source found=%v err=%v", found, err)
	}
	page, err := repo.MirrorCatalogPage(ctx, tenant, domain.SourceXLSX, "xlsx", 0, 10)
	if err != nil || !reflect.DeepEqual(mirrorCodes(page), []string{"1", "2"}) {
		t.Fatalf("page codes=%v err=%v", mirrorCodes(page), err)
	}
	counts, err := repo.MirrorEANCollisionCounts(ctx, tenant, domain.SourceXLSX)
	if err != nil || !reflect.DeepEqual(counts, map[string]int{"shared": 2}) {
		t.Fatalf("counts=%v err=%v", counts, err)
	}
}

func TestMirrorExcludesAbsentRows(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
	})

	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "1", "current", stringPtr("dup"), false, nil, nil)
	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "2", "absent", stringPtr("dup"), true, nil, nil)

	rows, err := repo.MirrorRows(ctx, tenant, domain.SourceXLSX)
	if err != nil || !reflect.DeepEqual(mirrorCodes(rows), []string{"1"}) {
		t.Fatalf("rows codes=%v err=%v", mirrorCodes(rows), err)
	}
	if _, found, err := repo.MirrorProductByCode(ctx, tenant, domain.SourceXLSX, "2"); err != nil || found {
		t.Fatalf("absent product found=%v err=%v", found, err)
	}
	page, err := repo.MirrorCatalogPage(ctx, tenant, domain.SourceXLSX, "", 0, 10)
	if err != nil || !reflect.DeepEqual(mirrorCodes(page), []string{"1"}) {
		t.Fatalf("page codes=%v err=%v", mirrorCodes(page), err)
	}
	counts, err := repo.MirrorEANCollisionCounts(ctx, tenant, domain.SourceXLSX)
	if err != nil || len(counts) != 0 {
		t.Fatalf("counts=%v err=%v", counts, err)
	}
}

func TestMirrorCatalogPagesDoNotDuplicateOrSkip(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
	})

	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "1", "one", nil, false, nil, nil)
	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "2", "two", nil, false, stringPtr("12.50"), stringPtr("4"))
	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "10", "ten", nil, false, nil, nil)
	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "SKU", "non-int", nil, false, nil, nil)
	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "3", "absent", nil, true, nil, nil)

	first, err := repo.MirrorCatalogPage(ctx, tenant, domain.SourceXLSX, "", 0, 2)
	if err != nil || !reflect.DeepEqual(mirrorCodes(first), []string{"1", "2", "10"}) {
		t.Fatalf("first page codes=%v err=%v", mirrorCodes(first), err)
	}
	second, err := repo.MirrorCatalogPage(ctx, tenant, domain.SourceXLSX, "", 2, 2)
	if err != nil || !reflect.DeepEqual(mirrorCodes(second), []string{"10"}) {
		t.Fatalf("second page codes=%v err=%v", mirrorCodes(second), err)
	}
	if first[0].Custo != nil || first[0].EstoqueTotal != nil {
		t.Fatalf("null numerics custo=%v estoque=%v", first[0].Custo, first[0].EstoqueTotal)
	}
	if first[1].Custo == nil || *first[1].Custo != "12.50" || first[1].EstoqueTotal == nil || *first[1].EstoqueTotal != "4" {
		t.Fatalf("numeric text custo=%v estoque=%v", first[1].Custo, first[1].EstoqueTotal)
	}
}

func TestMirrorEANCollisionCountsAreGlobalAndGReq2(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM products_mirror WHERE tenant_id=$1`, tenant)
	})

	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "1", "one", stringPtr("duplicate"), false, nil, nil)
	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "2", "two", stringPtr("duplicate"), false, nil, nil)
	seedMirrorProduct(t, pool, tenant, domain.SourceXLSX, "3", "three", stringPtr("unique"), false, nil, nil)

	counts, err := repo.MirrorEANCollisionCounts(ctx, tenant, domain.SourceXLSX)
	if err != nil || !reflect.DeepEqual(counts, map[string]int{"duplicate": 2}) {
		t.Fatalf("counts=%v err=%v", counts, err)
	}
}

func seedMirrorProduct(t *testing.T, pool *pgxpool.Pool, tenant string, source domain.ImportSource, code, description string, ean *string, absent bool, cost, stock *string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `INSERT INTO products_mirror
		(tenant_id,source,codigo_produto,descricao,ean,custo,estoque_total,absent_in_last_snapshot)
		VALUES ($1,$2,$3,$4,$5,$6::numeric,$7::numeric,$8)`, tenant, source, code, description, ean, cost, stock, absent)
	if err != nil {
		t.Fatal(err)
	}
}

func mirrorCodes(rows []domain.MirrorProduct) []string {
	codes := make([]string, len(rows))
	for i := range rows {
		codes[i] = rows[i].CodigoProduto
	}
	return codes
}

func mirrorSnapshot(id, hash string, importedAt time.Time, rows []domain.NormalizedRow) domain.ImportSnapshot {
	return domain.ImportSnapshot{ID: domain.ImportID(id), Protocol: domain.Protocol("#" + hash), FileSHA256: domain.FileSHA256(hash), Source: domain.SourceXLSX, ImportedAt: importedAt, Status: domain.ImportStatusCompleted, AcceptedRows: rows}
}

func stringPtr(value string) *string { return &value }
