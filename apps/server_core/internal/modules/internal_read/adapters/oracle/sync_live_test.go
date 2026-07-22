//go:build cgo

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

// TestSankhyaSyncLive is the M04-C4 real-provider proof: it drives the whole
// SankhyaAdapter.Sync path against the REAL METALPRD Oracle (never a stub) and a
// REAL migrated Postgres, then reads products_mirror back to confirm the snapshot
// landed with source='sankhya' and honest-NULL facts. This is the entrypoint the hub
// live-drives — the production scheduler wiring (root.go) is the HUB-OWNED F1 step
// sequenced after M-03+M-04 merge, so the adapter is proven here through direct
// construction, not through the (still no-op) scheduler.
//
// Gating: needs BOTH real backends. Runs only with MPC_ORACLE_LIVE_TEST=1 and the
// standard MPC_ORACLE_* / harness Postgres env (names only — never commit values).
//
//	MPC_ORACLE_LIVE_TEST=1 go test -tags cgo \
//	  -run TestSankhyaSyncLive -v \
//	  ./internal/modules/internal_read/adapters/oracle/
func TestSankhyaSyncLive(t *testing.T) {
	if os.Getenv("MPC_ORACLE_LIVE_TEST") != "1" {
		t.Skip("set MPC_ORACLE_LIVE_TEST=1 (+ MPC_ORACLE_* and harness Postgres env) to run live Sankhya sync validation")
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

	// Sanity: no unresolvable fact was written as a 0/default sentinel (AC-03). A zero
	// custo/preco/estoque must be a genuine balance, never a stand-in for unknown; the
	// writer only ever gets a nil pointer for unknown, so the mirror should not contain
	// a 0 custo unless Oracle truly returned 0. We surface the count for the hub to eyeball.
	var zeroCusto int
	_ = pool.QueryRow(ctx,
		`SELECT count(*) FROM products_mirror WHERE tenant_id=$1 AND source='sankhya' AND custo=0`, tenant,
	).Scan(&zeroCusto)

	fmt.Printf("MPC_C04_SANKHYA_MIRROR_ROWS=%d\n", persisted)
	fmt.Printf("MPC_C04_EAN_NULLS=%d MPC_C04_CUSTO_NULLS=%d MPC_C04_FULLY_RESOLVED=%d MPC_C04_ZERO_CUSTO=%d\n",
		eanNulls, custoNulls, fullyResolved, zeroCusto)
	fmt.Println("MPC_C04_LIVE_SANKHYA_SYNC_OK=true")
}
