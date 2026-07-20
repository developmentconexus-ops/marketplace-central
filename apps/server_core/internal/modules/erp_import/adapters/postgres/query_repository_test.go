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

func TestQueriesNewestCompletedAndTenantIsolation(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	other := tenant + "-other"
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE tenant_id=$1`, other)
	})
	now := time.Now().UTC().Truncate(time.Microsecond)
	grupo, descrGrupo := "07", "Ferramentas"
	completed := domain.ImportSnapshot{ID: "66666666-6666-6666-6666-666666666666", Protocol: "#201-E", FileSHA256: "same-hash", ImportedAt: now, Status: domain.ImportStatusCompleted, AcceptedRows: []domain.NormalizedRow{{Codprod: "A", Descrprod: "Alpha", Custo: "2.5", StockPhysical: "3", Grupo: &grupo, DescrGrupo: &descrGrupo}}}
	rejected := domain.ImportSnapshot{ID: "77777777-7777-7777-7777-777777777777", Protocol: "#202-E", FileSHA256: "later-rejected", ImportedAt: now.Add(time.Second), Status: domain.ImportStatusRejected, Issues: []domain.Issue{{Row: 2, Kind: domain.Rejection, Code: domain.CodeInvalidCusto}}}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, completed); err != nil {
		t.Fatal(err)
	}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, rejected); err != nil {
		t.Fatal(err)
	}
	if err := repo.PersistSnapshotAtomically(ctx, other, domain.ImportSnapshot{ID: "88888888-8888-8888-8888-888888888888", Protocol: "#201-E", FileSHA256: completed.FileSHA256, ImportedAt: now, Status: domain.ImportStatusCompleted}); err != nil {
		t.Fatal(err)
	}
	list, err := repo.ListImports(ctx, tenant)
	if err != nil || len(list) != 2 || list[0].ID != rejected.ID {
		t.Fatalf("list=%#v err=%v", list, err)
	}
	found, err := repo.FindByFileSHA256(ctx, tenant, completed.FileSHA256)
	if err != nil || found == nil || found.ID != completed.ID {
		t.Fatalf("found=%#v err=%v", found, err)
	}
	latest, err := repo.LatestCompletedSnapshot(ctx, tenant, domain.SourceXLSX)
	if err != nil || latest.ID != completed.ID || len(latest.AcceptedRows) != 1 {
		t.Fatalf("latest=%#v err=%v", latest, err)
	}
	if got := latest.AcceptedRows[0]; got.Grupo == nil || *got.Grupo != grupo || got.DescrGrupo == nil || *got.DescrGrupo != descrGrupo {
		t.Fatalf("grupo round trip lost: %#v", got)
	}
	if _, err := repo.GetImport(ctx, other, completed.ID); !errors.Is(err, ports.ErrImportNotFound) {
		t.Fatalf("cross tenant get err=%v", err)
	}
	missing, err := repo.FindByFileSHA256(ctx, other, rejected.FileSHA256)
	if err != nil || missing != nil {
		t.Fatalf("cross tenant find=%#v err=%v", missing, err)
	}
	if _, err := repo.LatestCompletedSnapshot(ctx, tenant+"-empty", domain.SourceXLSX); !errors.Is(err, ports.ErrImportNotFound) {
		t.Fatalf("empty latest err=%v", err)
	}
}

// Regression: a lenient client-catalog snapshot stores NULL custo/stock_physical
// (ADR-17 honest-unknown). LatestCompletedSnapshot must read those NULLs back as
// empty textual fields rather than failing the row scan (which previously surfaced
// to callers as a masked source_unavailable and broke the whole catalog page).
func TestLatestCompletedSnapshotReadsNullCostAndStock(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	now := time.Now().UTC().Truncate(time.Microsecond)
	ean := "7891234567895"
	snap := domain.ImportSnapshot{
		ID: "99999999-9999-9999-9999-999999999999", Protocol: "#209-E",
		FileSHA256: "null-cost-stock", ImportedAt: now, Source: domain.SourceCatalogoCliente,
		Status:       domain.ImportStatusCompleted,
		AcceptedRows: []domain.NormalizedRow{{Codprod: "500", Descrprod: "Faca sem custo", EAN: &ean}},
	}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, snap); err != nil {
		t.Fatal(err)
	}
	latest, err := repo.LatestCompletedSnapshot(ctx, tenant, domain.SourceCatalogoCliente)
	if err != nil {
		t.Fatalf("read back NULL cost/stock snapshot errored: %v", err)
	}
	if len(latest.AcceptedRows) != 1 {
		t.Fatalf("rows=%d", len(latest.AcceptedRows))
	}
	if got := latest.AcceptedRows[0]; got.Custo != "" || got.StockPhysical != "" {
		t.Fatalf("NULL custo/stock must read back empty (honest-unknown), got custo=%q stock=%q", got.Custo, got.StockPhysical)
	}
}
