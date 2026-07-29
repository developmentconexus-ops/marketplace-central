//go:build integration

package mirror_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/mirror"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func ptrF(v float64) *float64 { return &v }
func ptrS(v string) *string   { return &v }

func TestPgWriterPersistsSellableAssortmentFields(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mirror_assortment")
	tenant := "mirror-assortment-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "DELETE FROM products_mirror WHERE tenant_id = $1", tenant)
	})

	w := mirror.NewPgWriter(pool)
	if _, err := w.ApplySnapshot(ctx, tenant, []mirror.Row{
		{CodigoProduto: "A", Usoprod: ptrS(" r "), ADEcommerce: ptrS(" n ")},
		{CodigoProduto: "B"},
	}, nil); err != nil {
		t.Fatalf("ApplySnapshot: %v", err)
	}

	var usoprod, adEcommerce string
	if err := pool.QueryRow(ctx, `SELECT usoprod, ad_ecommerce FROM products_mirror WHERE tenant_id=$1 AND codigo_produto='A'`, tenant).Scan(&usoprod, &adEcommerce); err != nil {
		t.Fatalf("read A: %v", err)
	}
	if usoprod != "R" || adEcommerce != "N" {
		t.Errorf("A assortment = %q/%q, want R/N", usoprod, adEcommerce)
	}
	var nullCount int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM products_mirror WHERE tenant_id=$1 AND codigo_produto='B' AND usoprod IS NULL AND ad_ecommerce IS NULL`, tenant).Scan(&nullCount); err != nil {
		t.Fatalf("read B NULL predicates: %v", err)
	}
	if nullCount != 1 {
		t.Errorf("B SQL NULL match count = %d, want 1", nullCount)
	}
}

// TestApplySnapshotKeepAbsentRoundTrips proves the load-bearing merge semantics
// against a real Postgres (M04-C5): honest-NULL round-trip, keep-absent flag on the
// present→absent transition with last-known values preserved (never deleted),
// reappear clearing the flag, tenant isolation, and the empty-snapshot guard.
func TestApplySnapshotKeepAbsentRoundTrips(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mirror_keepabsent")
	token := time.Now().UTC().Format("150405.000000000")
	tenant := "mirror-a-" + token
	other := "mirror-b-" + token
	t.Cleanup(func() {
		for _, tbl := range []string{"products_mirror_stock_locations", "products_mirror"} {
			_, _ = pool.Exec(context.Background(), "DELETE FROM "+tbl+" WHERE tenant_id = ANY($1)", []string{tenant, other})
		}
	})

	w := mirror.NewPgWriter(pool)

	// A neighbouring tenant row must never be touched by tenant's sweep.
	if _, err := pool.Exec(ctx, `INSERT INTO products_mirror(tenant_id, source, codigo_produto, custo) VALUES($1,'xlsx','X',7)`, other); err != nil {
		t.Fatalf("seed other tenant: %v", err)
	}

	readRow := func(code string) (custo *float64, absent bool, staleSet bool) {
		var stale *time.Time
		err := pool.QueryRow(ctx, `SELECT custo::double precision, absent_in_last_snapshot, stale_since FROM products_mirror WHERE tenant_id=$1 AND codigo_produto=$2`, tenant, code).Scan(&custo, &absent, &stale)
		if err != nil {
			t.Fatalf("read %s: %v", code, err)
		}
		return custo, absent, stale != nil
	}
	countRows := func() int {
		var n int
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM products_mirror WHERE tenant_id=$1`, tenant).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}

	// Round 1: A (custo=10, stock) + B (custo NULL — honest-unknown).
	if _, err := w.ApplySnapshot(ctx, tenant,
		[]mirror.Row{
			{CodigoProduto: "A", Custo: ptrF(10), EstoqueTotal: ptrF(5), Descricao: ptrS("Prod A")},
			{CodigoProduto: "B"},
		},
		[]mirror.StockLocation{{CodigoProduto: "A", LocalCodigo: "1", Quantidade: ptrF(5)}},
	); err != nil {
		t.Fatalf("round1: %v", err)
	}
	if custo, absent, _ := readRow("A"); absent || custo == nil || *custo != 10 {
		t.Errorf("A round1: custo=%v absent=%v, want 10/false", custo, absent)
	}
	if custo, absent, _ := readRow("B"); absent || custo != nil {
		t.Errorf("B round1: custo=%v absent=%v, want NULL/false (honest-unknown)", custo, absent)
	}

	// Round 2: A only. B must be flagged absent, its last-known custo (NULL) kept, not deleted.
	if _, err := w.ApplySnapshot(ctx, tenant, []mirror.Row{{CodigoProduto: "A", Custo: ptrF(10)}}, nil); err != nil {
		t.Fatalf("round2: %v", err)
	}
	if countRows() != 2 {
		t.Errorf("round2 row count = %d, want 2 (keep-absent never deletes)", countRows())
	}
	if _, absent, _ := readRow("A"); absent {
		t.Errorf("A round2: absent=true, want false (still present)")
	}
	if _, absent, staleSet := readRow("B"); !absent || !staleSet {
		t.Errorf("B round2: absent=%v staleSet=%v, want true/true", absent, staleSet)
	}

	// Round 3: B reappears → flag clears, stale_since back to NULL.
	if _, err := w.ApplySnapshot(ctx, tenant, []mirror.Row{{CodigoProduto: "A"}, {CodigoProduto: "B", Custo: ptrF(4)}}, nil); err != nil {
		t.Fatalf("round3: %v", err)
	}
	if custo, absent, staleSet := readRow("B"); absent || staleSet || custo == nil || *custo != 4 {
		t.Errorf("B round3: custo=%v absent=%v staleSet=%v, want 4/false/false", custo, absent, staleSet)
	}

	// Empty snapshot is refused fail-closed and must NOT flip A/B absent.
	if _, err := w.ApplySnapshot(ctx, tenant, nil, nil); !errors.Is(err, mirror.ErrEmptySnapshot) {
		t.Errorf("empty snapshot err = %v, want ErrEmptySnapshot", err)
	}
	if _, absent, _ := readRow("A"); absent {
		t.Errorf("A after empty snapshot: absent=true, want false (guard must protect presence)")
	}

	// Tenant isolation: the neighbour's row is untouched.
	var otherAbsent bool
	if err := pool.QueryRow(ctx, `SELECT absent_in_last_snapshot FROM products_mirror WHERE tenant_id=$1 AND codigo_produto='X'`, other).Scan(&otherAbsent); err != nil {
		t.Fatalf("read other: %v", err)
	}
	if otherAbsent {
		t.Errorf("other tenant row flagged absent — sweep leaked across tenants")
	}
}
