//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestPersistSnapshotHandlesAFullCatalogWellInsideTheRouteDeadline pins the
// reason the ERP upload used to answer 504: a real client catalog is a couple of
// thousand rows, and one INSERT round trip per row did not fit in the 120s batch
// route deadline. The persist path now COPYs both the product rows and the issue
// rows. This is a budget guard plus a correctness-at-scale check: it asserts a
// full catalog lands intact (including honest-unknown NULLs) and well inside the
// deadline. It does not by itself measure the COPY-vs-per-row delta — read the
// logged duration for that.
func TestPersistSnapshotHandlesAFullCatalogWellInsideTheRouteDeadline(t *testing.T) {
	ctx := context.Background()
	repo, tenant := integrationRepo(t)

	const productCount = 3000
	const issueCount = 500
	rows := make([]domain.NormalizedRow, 0, productCount)
	marca, grupo, descrGrupo := "ACME", "10", "EPI"
	for i := 0; i < productCount; i++ {
		ean := fmt.Sprintf("789%010d", i)
		rows = append(rows, domain.NormalizedRow{
			Codprod:   fmt.Sprintf("SCALE-%05d", i),
			Descrprod: fmt.Sprintf("Produto de escala %05d", i),
			// Honest-unknown custo/estoque (the lenient client-catalog shape):
			// these must survive COPY as SQL NULL, never a fabricated zero.
			EAN:        &ean,
			Marca:      &marca,
			Grupo:      &grupo,
			DescrGrupo: &descrGrupo,
		})
	}
	column := "custo"
	offending := "abc"
	issues := make([]domain.Issue, 0, issueCount)
	for i := 0; i < issueCount; i++ {
		issues = append(issues, domain.Issue{
			Row:            i + 2,
			Column:         &column,
			Kind:           domain.Warning,
			Code:           domain.CodeInvalidCusto,
			Detail:         "unreadable cost cell",
			OffendingValue: &offending,
		})
	}

	snapshot := domain.ImportSnapshot{
		ID:           "55555555-5555-5555-5555-555555555555",
		Protocol:     "#900-E",
		FileSHA256:   "hash-scale",
		Source:       domain.SourceCatalogoCliente,
		ImportedAt:   time.Now().UTC().Truncate(time.Microsecond),
		Status:       domain.ImportStatusCompleted,
		AcceptedRows: rows,
		Issues:       issues,
	}

	started := time.Now()
	if err := repo.PersistSnapshotAtomically(ctx, tenant, snapshot); err != nil {
		t.Fatalf("persist %d-row snapshot: %v", productCount, err)
	}
	elapsed := time.Since(started)
	// A quarter of the route deadline, leaving room for parse, validation and
	// the market enqueue that share the same request.
	if elapsed > 30*time.Second {
		t.Fatalf("persisting %d rows took %s, want well under the 120s route deadline", productCount, elapsed)
	}
	t.Logf("persisted %d products + %d issues in %s", productCount, issueCount, elapsed)

	pool, _ := testpostgres.OpenPool(t, tenant)
	var products, storedIssues, nullCost int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM erp_import_products WHERE tenant_id=$1 AND protocol_id=$2`, tenant, snapshot.ID).Scan(&products); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM erp_import_issues WHERE tenant_id=$1 AND protocol_id=$2`, tenant, snapshot.ID).Scan(&storedIssues); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM erp_import_products WHERE tenant_id=$1 AND protocol_id=$2 AND custo IS NULL AND stock_physical IS NULL`, tenant, snapshot.ID).Scan(&nullCost); err != nil {
		t.Fatal(err)
	}
	if products != productCount || storedIssues != issueCount || nullCost != productCount {
		t.Fatalf("products=%d issues=%d honest-null=%d, want %d/%d/%d", products, storedIssues, nullCost, productCount, issueCount, productCount)
	}
}
