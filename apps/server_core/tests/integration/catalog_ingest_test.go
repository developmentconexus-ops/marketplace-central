//go:build integration

package integration

import (
	"context"
	"strconv"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/catalog"
	"marketplace-central/apps/server_core/internal/contexts/catalog/contracts"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func catalogObservation(t *testing.T, tid tenant.ID, externalKey, description, ean, hash string) contracts.ProductObservation {
	t.Helper()
	e, err := provenance.NewEvidence("sankhya", "product", externalKey, time.Now().UTC(), hash)
	if err != nil {
		t.Fatalf("provenance.NewEvidence: %v", err)
	}
	desc, err := fact.NewKnown(description, e)
	if err != nil {
		t.Fatalf("fact.NewKnown: %v", err)
	}
	key, err := contracts.NewSourceProductKey(tid, "sankhya", "sankhya-it-01", "product", externalKey)
	if err != nil {
		t.Fatalf("contracts.NewSourceProductKey: %v", err)
	}
	id, err := contracts.NewIdentifier(contracts.IdentifierEAN, ean)
	if err != nil {
		t.Fatalf("contracts.NewIdentifier: %v", err)
	}
	return contracts.ProductObservation{Key: key, Description: desc,
		Identifiers: []contracts.Identifier{id}, Evidence: e}
}

func TestCatalogIngestPersistsAndIsIdempotent(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_catalog_ingest")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	module := catalog.New(pool)

	ctx := context.Background()
	stamp := time.Now().UTC().UnixNano()
	externalKey := "it-" + time.Now().UTC().Format("150405.000000000")

	first, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, externalKey, "Cafeteira Eletrica", "789ean"+itoa(stamp), "sha256:one"+itoa(stamp)))
	if err != nil {
		t.Fatalf("first IngestProduct: %v", err)
	}
	if first.Disposition != contracts.DispositionCreated || first.Version != 1 {
		t.Fatalf("first = %+v, want created/1", first)
	}

	// Same payload hash: nothing moves, and the id does not.
	second, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, externalKey, "Cafeteira Eletrica", "789ean"+itoa(stamp), "sha256:one"+itoa(stamp)))
	if err != nil {
		t.Fatalf("second IngestProduct: %v", err)
	}
	if second.Disposition != contracts.DispositionIdempotent {
		t.Fatalf("second Disposition = %q, want idempotent", second.Disposition)
	}
	if second.ProductID != first.ProductID {
		t.Fatalf("ProductID moved: %q then %q", first.ProductID, second.ProductID)
	}

	// Changed payload: version 2, same identity.
	third, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, externalKey, "Cafeteira Eletrica Inox", "789ean"+itoa(stamp), "sha256:two"+itoa(stamp)))
	if err != nil {
		t.Fatalf("third IngestProduct: %v", err)
	}
	if third.Version != 2 || third.ProductID != first.ProductID {
		t.Fatalf("third = %+v, want version 2 on %s", third, first.ProductID)
	}

	got, found, err := module.Reader().ByProductID(ctx, tid, first.ProductID)
	if err != nil {
		t.Fatalf("ByProductID: %v", err)
	}
	if !found {
		t.Fatalf("ByProductID did not find %s", first.ProductID)
	}
	if got.Description != "Cafeteira Eletrica Inox" || got.Version != 2 {
		t.Fatalf("summary = %+v, want the version 2 description", got)
	}
}

// The RLS proof. Reading with a different tenant scope must return nothing —
// not filtered by the application, filtered by the database.
func TestCatalogIsInvisibleToAnotherTenant(t *testing.T) {
	pool, cfg := testpostgres.OpenPool(t, "tenant_catalog_rls")
	tid, err := tenant.Parse(cfg.DefaultTenantID)
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	module := catalog.New(pool)

	ctx := context.Background()
	stamp := time.Now().UTC().UnixNano()
	created, err := module.IngestProduct(ctx,
		catalogObservation(t, tid, "rls-"+itoa(stamp), "Cafeteira", "789rls"+itoa(stamp), "sha256:rls"+itoa(stamp)))
	if err != nil {
		t.Fatalf("IngestProduct: %v", err)
	}

	intruder, err := tenant.Parse("tnt_intruder")
	if err != nil {
		t.Fatalf("tenant.Parse: %v", err)
	}
	_, found, err := module.Reader().ByProductID(ctx, intruder, created.ProductID)
	if err != nil {
		t.Fatalf("cross-tenant ByProductID: %v", err)
	}
	if found {
		t.Fatalf("tenant %s read product %s belonging to %s", intruder, created.ProductID, tid)
	}
}

// itoa keeps the fixture strings unique per run, so a re-run does not collide
// with the rows the previous run left behind.
func itoa(n int64) string { return strconv.FormatInt(n, 10) }
