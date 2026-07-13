//go:build integration

package application

import (
	"context"
	"testing"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

type deterministicProductMatcher struct {
	product internalreaddomain.ProductCandidate
}

func (m deterministicProductMatcher) FindProductsForLinking(_ context.Context, input internalreadports.FindProductsInput) ([]internalreaddomain.ProductCandidate, error) {
	if input.SellerSKU != nil && *input.SellerSKU == "internal-4242" {
		return []internalreaddomain.ProductCandidate{m.product}, nil
	}
	return nil, nil
}

func TestGenerateLinkCandidatesPersistsDeterministicPostgresResult(t *testing.T) {
	ctx := context.Background()
	pool, cfg := testpostgres.OpenPool(t, "tenant_harness_product_links")
	testID := time.Now().UTC().Format("150405.000000000")
	installationID := "installation-" + testID
	providerItemID := "MLB-" + testID
	now := time.Date(2026, 7, 11, 12, 0, 0, 0, time.UTC)

	snapshotRepo := productlinkspostgres.NewListingSnapshotRepository(pool, cfg.DefaultTenantID)
	candidateRepo := productlinkspostgres.NewLinkCandidateRepository(pool, cfg.DefaultTenantID)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_link_candidates WHERE tenant_id=$1 AND installation_id=$2`, cfg.DefaultTenantID, installationID)
		_, _ = pool.Exec(context.Background(), `DELETE FROM product_link_listing_snapshots WHERE tenant_id=$1 AND installation_id=$2`, cfg.DefaultTenantID, installationID)
	})
	if err := snapshotRepo.UpsertListingSnapshots(ctx, []productlinksdomain.ListingSnapshot{{
		InstallationID: installationID, ProviderCode: "mercado_livre", ProviderItemID: providerItemID,
		SellerSKU: "internal-4242", Title: "Deterministic fixture", FetchedAt: now, CreatedAt: now, UpdatedAt: now,
	}}); err != nil {
		t.Fatalf("seed listing snapshot: %v", err)
	}

	reference := "FAB-2020.C.FLX"
	canonicalProductID := internalreaddomain.InternalProductID(4242)
	svc := NewGenerationService(GenerationServiceConfig{
		Snapshots: snapshotRepo,
		Matcher: deterministicProductMatcher{product: internalreaddomain.ProductCandidate{
			InternalProductID: &canonicalProductID, ProductID: 99, Name: "Deterministic product", ReferenceCode: &reference, IsActive: true,
		}},
		Store: candidateRepo,
		Now:   func() time.Time { return now },
	})
	result, err := svc.GenerateLinkCandidates(ctx, GenerateLinkCandidatesInput{InstallationID: installationID, Limit: 20})
	if err != nil {
		t.Fatalf("generate deterministic link candidates: %v", err)
	}
	if result.GeneratedCount != 1 || len(result.Items) != 1 || result.Items[0].State != productlinksdomain.LinkCandidateStateExactSKU {
		t.Fatalf("unexpected generation result: %#v", result)
	}
	persisted, err := svc.ListLinkCandidates(ctx, ListLinkCandidatesInput{InstallationID: installationID, Limit: 20})
	if err != nil {
		t.Fatalf("list persisted candidates: %v", err)
	}
	if len(persisted) != 1 || persisted[0].InternalProductID == nil || *persisted[0].InternalProductID != 4242 {
		t.Fatalf("unexpected persisted candidate: %#v", persisted)
	}
}
