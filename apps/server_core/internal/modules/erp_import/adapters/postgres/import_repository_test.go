//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	erppostgres "marketplace-central/apps/server_core/internal/modules/erp_import/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func integrationRepo(t *testing.T) (*erppostgres.Repository, string) {
	t.Helper()
	testpostgres.SkipWithoutTarget(t)
	tenant := "erp-import-" + time.Now().UTC().Format("150405.000000000")
	pool, _ := testpostgres.OpenPool(t, tenant)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM erp_import_protocols WHERE tenant_id=$1`, tenant)
	})
	return erppostgres.NewRepository(pool), tenant
}

func TestPersistSnapshotAtomicDuplicateRejectedRollbackAndLock(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)
	column, offending := "custo", "0"
	now := time.Now().UTC().Truncate(time.Microsecond)
	s := domain.ImportSnapshot{ID: "11111111-1111-1111-1111-111111111111", Protocol: "#101-E", FileSHA256: "hash-a", Source: domain.SourceXLSX, ImportedAt: now, Status: domain.ImportStatusCompleted, AcceptedRows: []domain.NormalizedRow{{Codprod: "A", Descrprod: "Alpha", Custo: "10.50", StockPhysical: "4"}}, Issues: []domain.Issue{{Row: 2, Column: &column, Kind: domain.Rejection, Code: domain.CodeInvalidCusto, Detail: "positive required", OffendingValue: &offending}}}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, s); err != nil {
		t.Fatal(err)
	}
	r, err := repo.GetImport(ctx, tenant, s.ID)
	if err != nil || r.AcceptedCount != 1 || r.RejectedCount != 1 || len(r.Issues) != 1 || r.Issues[0].Column == nil || *r.Issues[0].Column != column || r.Issues[0].OffendingValue == nil || *r.Issues[0].OffendingValue != offending {
		t.Fatalf("round trip: %#v err=%v", r, err)
	}
	err = repo.PersistSnapshotAtomically(ctx, tenant, domain.ImportSnapshot{ID: "22222222-2222-2222-2222-222222222222", Protocol: "#102-E", FileSHA256: s.FileSHA256, ImportedAt: now.Add(time.Second), Status: domain.ImportStatusCompleted})
	var duplicate *ports.DuplicateFileError
	if !errors.As(err, &duplicate) || duplicate.ExistingID != s.ID || duplicate.ExistingProtocol != s.Protocol {
		t.Fatalf("duplicate=%#v err=%v", duplicate, err)
	}

	rejected := domain.ImportSnapshot{ID: "33333333-3333-3333-3333-333333333333", Protocol: "#103-E", FileSHA256: "hash-rejected", ImportedAt: now.Add(2 * time.Second), Status: domain.ImportStatusRejected, Issues: []domain.Issue{{Row: 2, Kind: domain.Rejection, Code: domain.CodeInvalidCusto}}}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, rejected); err != nil {
		t.Fatal(err)
	}
	pool, _ := testpostgres.OpenPool(t, tenant)
	var products int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM erp_import_products WHERE tenant_id=$1 AND protocol_id=$2`, tenant, rejected.ID).Scan(&products); err != nil || products != 0 {
		t.Fatalf("products=%d err=%v", products, err)
	}

	bad := domain.ImportSnapshot{ID: "44444444-4444-4444-4444-444444444444", Protocol: "#104-E", FileSHA256: "hash-bad", ImportedAt: now, Status: domain.ImportStatusCompleted, AcceptedRows: []domain.NormalizedRow{{Codprod: "B", Descrprod: "Bad", Custo: "0", StockPhysical: "1"}}}
	if err := repo.PersistSnapshotAtomically(ctx, tenant, bad); err == nil {
		t.Fatal("constraint violation succeeded")
	}
	var protocols int
	_ = pool.QueryRow(ctx, `SELECT count(*) FROM erp_import_protocols WHERE tenant_id=$1 AND id=$2`, tenant, bad.ID).Scan(&protocols)
	if protocols != 0 {
		t.Fatalf("partial protocol count=%d", protocols)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	var held bool
	// Must hold the SAME advisory-lock key production uses ('erp_import:'||tenant), else no contention
	// is created and ErrImportInProgress below would never fire (the assertion would pass vacuously).
	if err := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock(hashtextextended('erp_import:'||$1, 0)) WHERE $1=$1`, tenant).Scan(&held); err != nil || !held {
		t.Fatalf("hold lock=%v err=%v", held, err)
	}
	err = repo.PersistSnapshotAtomically(ctx, tenant, domain.ImportSnapshot{ID: "55555555-5555-5555-5555-555555555555", Protocol: "#105-E", FileSHA256: "hash-lock", ImportedAt: now, Status: domain.ImportStatusCompleted})
	if !errors.Is(err, ports.ErrImportInProgress) {
		t.Fatalf("lock err=%v", err)
	}
}
