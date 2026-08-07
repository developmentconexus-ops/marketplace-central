//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestCatalogRLSBlocksTheOtherTenant is the only proof that FORCE RLS on the
// catalog schema is not decorative. The application connects as the table
// owner, which is also a superuser with rolbypassrls: every policy evaluated
// against that role is skipped, not enforced. This test connects as the
// least-privilege role (mpc_app) via SET LOCAL ROLE inside a transaction on
// the existing owner connection, because mpc_app is NOLOGIN and cannot be
// dialed directly -- repointing the production DSN is out of scope for this
// task. PostgreSQL evaluates rolbypassrls and policy application against the
// CURRENT effective user, so SET LOCAL ROLE does subject the statement to RLS.
func TestCatalogRLSBlocksTheOtherTenant(t *testing.T) {
	pool, _ := testpostgres.OpenPool(t, "tenant_catalog_rls_role")
	ctx := context.Background()

	stamp := itoa(time.Now().UTC().UnixNano())
	tenantA := "tnt_rls_a_" + stamp
	tenantB := "tnt_rls_b_" + stamp
	productA := "prd_rlsa" + stamp
	productB := "prd_rlsb" + stamp

	// arrange: as the writer (the owning superuser bypasses RLS), insert one
	// product for tenant A and one for tenant B.
	insertRLSFixtureProduct(t, ctx, pool, tenantA, productA)
	insertRLSFixtureProduct(t, ctx, pool, tenantB, productB)
	t.Cleanup(func() {
		cleanupCtx := context.Background()
		if _, err := pool.Exec(cleanupCtx,
			`DELETE FROM catalog.products WHERE product_id = ANY($1)`,
			[]string{productA, productB}); err != nil {
			t.Errorf("cleanup rls fixture products: %v", err)
		}
	})

	bothIDs := []string{productA, productB}

	// Positive control: the SAME query, run as the owning superuser (no SET
	// ROLE, no app.tenant_id), must return BOTH rows. Without this control an
	// empty result from mpc_app below would be indistinguishable from an
	// empty table, a failed insert, or a typo'd tenant/product id -- and would
	// prove nothing about RLS at all.
	uncontrolled := countRLSFixtureProducts(t, ctx, pool, bothIDs)
	if uncontrolled != 2 {
		t.Fatalf("positive control (owner connection, no RLS in effect) count = %d, want 2; "+
			"fixture is broken so the mpc_app assertions below would prove nothing", uncontrolled)
	}

	// act: connect as mpc_app scoped to tenant A. Only tenant A's row must be
	// visible even though the query does not filter by tenant_id itself.
	gotA := countAsCatalogAppRole(t, ctx, pool, &tenantA, bothIDs)
	if gotA != 1 {
		t.Fatalf("mpc_app scoped to tenant A: count = %d, want 1", gotA)
	}

	// act: same role, tenant B's scope, in a NEW transaction -- SET LOCAL dies
	// with the transaction, so this is not a leftover from tenant A's scope.
	gotB, idB := countAndOneIDAsCatalogAppRole(t, ctx, pool, &tenantB, bothIDs)
	if gotB != 1 {
		t.Fatalf("mpc_app scoped to tenant B: count = %d, want 1", gotB)
	}
	if idB != productB {
		t.Fatalf("mpc_app scoped to tenant B: saw product_id %q, want %q (tenant A's row leaked)", idB, productB)
	}

	// negative control: mpc_app with NO app.tenant_id set must fail closed
	// (current_setting('app.tenant_id', true) is NULL, and tenant_id = NULL
	// is NULL, so no row qualifies) rather than fail open.
	gotNone := countAsCatalogAppRole(t, ctx, pool, nil, bothIDs)
	if gotNone != 0 {
		t.Fatalf("mpc_app with no tenant scope: count = %d, want 0 (RLS must fail closed, not open)", gotNone)
	}
}

// insertRLSFixtureProduct writes directly as the owning connection (which
// bypasses RLS), so the arrange step does not itself depend on the role or
// policy under test.
func insertRLSFixtureProduct(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenantID, productID string) {
	t.Helper()
	if _, err := pool.Exec(ctx, `
		INSERT INTO catalog.products
			(tenant_id, product_id, version, description_state, description_value,
			 description_reason, last_payload_hash)
		VALUES ($1, $2, 1, 'unknown', NULL, 'rls-role fixture', 'sha256:rlsrolefixture')`,
		tenantID, productID); err != nil {
		t.Fatalf("insert rls fixture product %s: %v", productID, err)
	}
}

// countRLSFixtureProducts counts on the pool's own (owner) connection, with
// no SET ROLE and no tenant scope -- the uncontrolled positive-control read.
func countRLSFixtureProducts(t *testing.T, ctx context.Context, pool *pgxpool.Pool, productIDs []string) int {
	t.Helper()
	var count int
	if err := pool.QueryRow(ctx,
		`SELECT count(*) FROM catalog.products WHERE product_id = ANY($1)`,
		productIDs).Scan(&count); err != nil {
		t.Fatalf("positive control count: %v", err)
	}
	return count
}

// countAsCatalogAppRole runs the count in a fresh transaction, switched into
// mpc_app via SET LOCAL ROLE. A nil tenantID leaves app.tenant_id unset,
// exercising the negative control.
func countAsCatalogAppRole(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenantID *string, productIDs []string) int {
	t.Helper()
	count, _ := countAndOneIDAsCatalogAppRole(t, ctx, pool, tenantID, productIDs)
	return count
}

// countAndOneIDAsCatalogAppRole returns both the row count and (if exactly
// one row is visible) its product_id, so the caller can prove which tenant's
// row RLS let through, not just how many.
func countAndOneIDAsCatalogAppRole(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenantID *string, productIDs []string) (int, string) {
	t.Helper()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, "SET LOCAL ROLE mpc_app"); err != nil {
		t.Fatalf("SET LOCAL ROLE mpc_app: %v", err)
	}
	if tenantID != nil {
		if _, err := tx.Exec(ctx, "SELECT set_config('app.tenant_id', $1, true)", *tenantID); err != nil {
			t.Fatalf("scope tenant: %v", err)
		}
	}

	rows, err := tx.Query(ctx,
		`SELECT product_id FROM catalog.products WHERE product_id = ANY($1) ORDER BY product_id`,
		productIDs)
	if err != nil {
		t.Fatalf("select as mpc_app: %v", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("scan as mpc_app: %v", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows as mpc_app: %v", err)
	}

	one := ""
	if len(ids) == 1 {
		one = ids[0]
	}
	return len(ids), one
}
