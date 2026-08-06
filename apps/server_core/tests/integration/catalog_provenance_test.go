//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace-central/apps/server_core/internal/contexts/catalog"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestRehydrationKeepsEverySourceKeyAndRealEvidence proves both halves of the
// defect at once: a product reached by two source keys must come back with two,
// and its evidence must name the system that actually observed it.
func TestRehydrationKeepsEverySourceKeyAndRealEvidence(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_catalog_provenance")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	module := catalog.New(pool)
	ctx := context.Background()

	stamp := time.Now().UTC().UnixNano()
	externalKey := "prov-" + time.Now().UTC().Format("150405.000000000")

	created, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, externalKey, "Cafeteira Provenance", "789prov"+itoa(stamp), "sha256:prov"+itoa(stamp)))
	if err != nil {
		t.Fatalf("IngestProduct: %v", err)
	}

	// A second source key already resolves to this same canonical product —
	// the state a future dedup/merge site (Tarefa 8) will eventually create.
	// Nothing in today's public Ingest API can reach it on its own: BySourceKey
	// only ever finds a key that ALREADY maps to the product, so Apply's own
	// merge logic never gets a chance to run for a brand-new key. Wired
	// directly at the storage boundary so what's under test here is the READ
	// path (rehydration), not the (separate, not-yet-built) linking site.
	linkMirrorSourceKey(t, ctx, pool, tid, created.ProductID, externalKey+"-mirror")

	got, found, err := module.Reader().ByProductID(ctx, tid, created.ProductID)
	if err != nil {
		t.Fatalf("ByProductID: %v", err)
	}
	if !found {
		t.Fatalf("ByProductID did not find %s", created.ProductID)
	}

	if len(got.SourceKeys) != 2 {
		t.Errorf("want 2 source keys, got %d", len(got.SourceKeys))
	}
	if got.DescriptionEvidenceSystem != "sankhya" {
		t.Errorf("system = %q, want %q", got.DescriptionEvidenceSystem, "sankhya")
	}

	// The write side must have actually recorded what it observed: rehydration
	// can only stop fabricating "catalog" as the system once there is a real
	// row to read instead.
	var observed int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM catalog.source_observations WHERE tenant_id = $1 AND product_id = $2`,
		tid.String(), created.ProductID).Scan(&observed); err != nil {
		t.Fatalf("count source_observations: %v", err)
	}
	if observed == 0 {
		t.Errorf("catalog.source_observations has no row for %s; the write side never recorded what it observed", created.ProductID)
	}
}

// linkMirrorSourceKey attaches a second source address to productID directly at
// the storage boundary, scoped to the tenant the same way Repository.withTenant
// does it. It exists because the public Ingest API cannot produce this state on
// its own today (see the comment at the call site).
func linkMirrorSourceKey(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tid tenant.ID, productID, externalKey string) {
	t.Helper()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", tid.String()); err != nil {
		t.Fatalf("scope tenant: %v", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO catalog.source_product_keys
			(tenant_id, source_system, source_instance, object_kind, external_key, product_id)
		VALUES ($1,'sankhya','sankhya-it-02','product',$2,$3)`,
		tid.String(), externalKey, productID); err != nil {
		t.Fatalf("insert mirror source key: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit: %v", err)
	}
}
