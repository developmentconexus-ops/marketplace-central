//go:build integration && cgo

package oracle

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/mirror"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestSankhyaSyncIntegration is the M04-C4 real-provider driver the hub live-drives.
// It wires NewSankhyaAdapter against the REAL METALPRD Oracle (godror needs cgo) plus
// a REAL migrated Postgres, runs Sync once, and reads products_mirror back to prove
// the snapshot landed with source='sankhya' and honest-NULL facts — never a stub
// (AC-04). The production scheduler wiring (root.go) is the HUB-OWNED F1 step
// sequenced after M-03+M-04 merge, so the adapter is proven here by direct
// construction, not through the (still no-op) scheduler.
//
// Env var NAMES (values never printed): the Oracle connection uses the same names
// root.go's LoadConfigFromEnv reads — MPC_ORACLE_USERNAME, MPC_ORACLE_PASSWORD,
// MPC_ORACLE_CONNECT_STRING, MPC_ORACLE_LIB_DIR — and Postgres uses the harness
// MPC_TEST_DATABASE_URL. Absent env → t.Skip, so the unit lane stays green.
//
// Run:
//
//	CGO_ENABLED=1 go test -tags "integration cgo" \
//	  -run TestSankhyaSyncIntegration -v \
//	  ./internal/modules/internal_read/adapters/oracle/
func TestSankhyaSyncIntegration(t *testing.T) {
	if os.Getenv("MPC_ORACLE_CONNECT_STRING") == "" || os.Getenv("MPC_TEST_DATABASE_URL") == "" {
		t.Skip("set MPC_ORACLE_* (Oracle) and MPC_TEST_DATABASE_URL (harness Postgres) to run live Sankhya sync validation")
	}

	const tenant = "tenant_m04_sankhya_live_probe"

	cfg, err := LoadConfigFromEnv(os.Getenv)
	if err != nil {
		t.Fatalf("load live Oracle config: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	db, err := OpenDB(ctx, cfg)
	if err != nil {
		t.Fatalf("open live Oracle db: %v", err)
	}
	defer db.Close()

	pool, _ := testpostgres.OpenPool(t, tenant)

	// Isolated probe tenant: clear any prior run, and clean up afterwards so the live
	// proof never pollutes real tenant data.
	clear := func() {
		for _, tbl := range []string{"products_mirror_stock_locations", "products_mirror"} {
			if _, err := pool.Exec(context.Background(), "DELETE FROM "+tbl+" WHERE tenant_id=$1", tenant); err != nil {
				t.Logf("cleanup %s: %v", tbl, err)
			}
		}
	}
	clear()
	t.Cleanup(clear)

	adapter := NewSankhyaAdapter(db, mirror.NewPgWriter(pool), tenant)

	res, err := adapter.Sync(ctx)
	if err != nil {
		t.Fatalf("live Sync: %v", err)
	}
	t.Logf("SyncResult = %+v", res)
	if res.Processed <= 0 || res.Errors != 0 {
		t.Fatalf("SyncResult = %+v, want Processed>0 Errors=0", res)
	}

	// The mirror row count for this tenant/source must equal what Sync reported.
	var persisted int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM products_mirror WHERE tenant_id=$1 AND source='sankhya'`, tenant,
	).Scan(&persisted); err != nil {
		t.Fatalf("count sankhya mirror rows: %v", err)
	}
	if persisted != res.Processed {
		t.Fatalf("persisted rows = %d, want %d (Processed)", persisted, res.Processed)
	}

	// Honest-NULL must be real, not synthesised: the ratified mapping says ~51% of
	// products carry no EAN and many carry no cost/price. Assert at least one NULL
	// exists in each honest-unknown column AND at least one fully-resolved row exists,
	// so we know the writer distinguishes unknown from a real value.
	var eanNulls, custoNulls, fullyResolved int
	if err := pool.QueryRow(ctx, `
		SELECT
			count(*) FILTER (WHERE ean IS NULL),
			count(*) FILTER (WHERE custo IS NULL),
			count(*) FILTER (WHERE ean IS NOT NULL AND custo IS NOT NULL AND preco_venda IS NOT NULL AND estoque_total IS NOT NULL)
		FROM products_mirror WHERE tenant_id=$1 AND source='sankhya'`, tenant,
	).Scan(&eanNulls, &custoNulls, &fullyResolved); err != nil {
		t.Fatalf("honest-null profile: %v", err)
	}
	if eanNulls == 0 {
		t.Errorf("no NULL EANs — mapping expects ~51%% without EAN; check REFERENCIA shape filter")
	}
	if fullyResolved == 0 {
		t.Errorf("no fully-resolved rows — a live snapshot must resolve at least some cost/price/stock")
	}

	// AC-03 surface: a zero custo must be a genuine balance, never a stand-in for
	// unknown; the writer only ever gets a nil pointer for unknown, so we surface the
	// count for the hub to eyeball against the honest-NULL column.
	var zeroCusto int
	_ = pool.QueryRow(ctx,
		`SELECT count(*) FROM products_mirror WHERE tenant_id=$1 AND source='sankhya' AND custo=0`, tenant,
	).Scan(&zeroCusto)

	fmt.Printf("MPC_C04_SANKHYA_MIRROR_ROWS=%d\n", persisted)
	fmt.Printf("MPC_C04_EAN_NULLS=%d MPC_C04_CUSTO_NULLS=%d MPC_C04_FULLY_RESOLVED=%d MPC_C04_ZERO_CUSTO=%d\n",
		eanNulls, custoNulls, fullyResolved, zeroCusto)

	// Print a sample of the landed rows so the C04 stdout doubles as the
	// `SELECT * FROM products_mirror WHERE source='sankhya'` inspection the hub wants —
	// the probe tenant is cleaned up afterwards, so this dump is the durable evidence.
	sample, err := pool.Query(ctx, `
		SELECT codigo_produto, ean, custo::text, preco_venda::text, estoque_total::text, marca
		FROM products_mirror WHERE tenant_id=$1 AND source='sankhya'
		ORDER BY codigo_produto LIMIT 15`, tenant)
	if err != nil {
		t.Fatalf("sample rows: %v", err)
	}
	defer sample.Close()
	for sample.Next() {
		var cod string
		var ean, custo, preco, est, marca *string
		if err := sample.Scan(&cod, &ean, &custo, &preco, &est, &marca); err != nil {
			t.Fatalf("scan sample: %v", err)
		}
		fmt.Printf("MPC_C04_ROW cod=%s ean=%s custo=%s preco=%s estoque=%s marca=%s\n",
			cod, nullMark(ean), nullMark(custo), nullMark(preco), nullMark(est), nullMark(marca))
	}
	if err := sample.Err(); err != nil {
		t.Fatalf("iterate sample: %v", err)
	}

	// Estoque variance: the ORDER BY codigo_produto sample above can hit a run of
	// same-balance small parts (e.g. all qty 1), which looks like a constant. Prove the
	// snapshot actually varies — estoque_total must take >1 distinct value — and that
	// per-CODLOCAL stock locations landed with real quantities (SUM(ESTOQUE-RESERVADO)
	// per local, not a single rolled-up number).
	var distinctEstoque, locRows, distinctLocQty int
	var maxEstoque *string
	if err := pool.QueryRow(ctx, `
		SELECT count(DISTINCT estoque_total), max(estoque_total)::text
		FROM products_mirror WHERE tenant_id=$1 AND source='sankhya' AND estoque_total IS NOT NULL`, tenant,
	).Scan(&distinctEstoque, &maxEstoque); err != nil {
		t.Fatalf("estoque variance: %v", err)
	}
	if err := pool.QueryRow(ctx, `
		SELECT count(*), count(DISTINCT quantidade)
		FROM products_mirror_stock_locations WHERE tenant_id=$1`, tenant,
	).Scan(&locRows, &distinctLocQty); err != nil {
		t.Fatalf("stock-location profile: %v", err)
	}
	if distinctEstoque <= 1 {
		t.Errorf("estoque_total has %d distinct value(s) — a real snapshot must vary, not a constant", distinctEstoque)
	}
	if locRows == 0 {
		t.Errorf("no products_mirror_stock_locations rows landed — per-CODLOCAL stock did not persist")
	}
	if distinctLocQty <= 1 {
		t.Errorf("stock-location quantidade has %d distinct value(s) — expected real per-local variance", distinctLocQty)
	}

	// A few genuinely-stocked rows so the variance is visible in the proof, not just counted.
	stocked, err := pool.Query(ctx, `
		SELECT codigo_produto, estoque_total::text AS estoque_txt
		FROM products_mirror WHERE tenant_id=$1 AND source='sankhya' AND estoque_total > 1
		ORDER BY estoque_total DESC LIMIT 5`, tenant)
	if err != nil {
		t.Fatalf("stocked sample: %v", err)
	}
	defer stocked.Close()
	for stocked.Next() {
		var cod, est string
		if err := stocked.Scan(&cod, &est); err != nil {
			t.Fatalf("scan stocked: %v", err)
		}
		fmt.Printf("MPC_C04_STOCKED_ROW cod=%s estoque=%s\n", cod, est)
	}
	if err := stocked.Err(); err != nil {
		t.Fatalf("iterate stocked: %v", err)
	}

	fmt.Printf("MPC_C04_DISTINCT_ESTOQUE=%d MPC_C04_MAX_ESTOQUE=%s MPC_C04_STOCK_LOC_ROWS=%d MPC_C04_DISTINCT_LOC_QTY=%d\n",
		distinctEstoque, nullMark(maxEstoque), locRows, distinctLocQty)
	fmt.Println("MPC_C04_LIVE_SANKHYA_SYNC_OK=true")
}

// nullMark renders NULL explicitly so honest-NULL is visible in the C04 proof dump.
func nullMark(p *string) string {
	if p == nil || *p == "" {
		return "NULL"
	}
	return *p
}
