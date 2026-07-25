//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

const summaryTenant = "tenant-e10-summary"

// The confirmation queue is work waiting on a human. A candidate in it resolved
// an anchor, so its STATE is exact_sku/exact_ean — the states the pending count
// reads to find unresolved and conflicting work cannot see it. Counting only
// those would report an installation with a full confirmation queue as having
// nothing pending, which is the one thing an operator uses this number for.
func TestPendingLinksCountsTheConfirmationQueue(t *testing.T) {
	testpostgres.SkipWithoutTarget(t)
	pool, _ := testpostgres.OpenPool(t, summaryTenant)
	ctx := context.Background()
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM product_link_candidates WHERE tenant_id = $1`, summaryTenant)
		pool.Close()
	})

	repo := productlinkspostgres.NewLinkCandidateRepository(pool, summaryTenant)
	at := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	const installation = "inst-summary"

	candidate := func(id, itemID string, state productlinksdomain.LinkCandidateState, status productlinksdomain.LinkCandidateMatchStatus) productlinksdomain.LinkCandidate {
		return productlinksdomain.LinkCandidate{
			CandidateID:             id,
			InstallationID:          installation,
			ProviderCode:            "mercado_livre",
			ProviderItemID:          itemID,
			State:                   state,
			MatchInput:              productlinksdomain.LinkCandidateMatchInputEAN,
			MatchValue:              "7909251260214",
			MatchStatus:             status,
			SourceSnapshotFetchedAt: &at,
			CreatedAt:               at,
			UpdatedAt:               at,
		}
	}

	identities := []productlinksdomain.ListingIdentity{
		{InstallationID: installation, ProviderItemID: "MLB-CONFIRM"},
		{InstallationID: installation, ProviderItemID: "MLB-REVIEW"},
	}
	candidates := []productlinksdomain.LinkCandidate{
		candidate("cand-confirm", "MLB-CONFIRM", productlinksdomain.LinkCandidateStateExactEAN, productlinksdomain.LinkCandidateMatchStatusConfirm),
		candidate("cand-review", "MLB-REVIEW", productlinksdomain.LinkCandidateStateUnresolved, productlinksdomain.LinkCandidateMatchStatusReview),
	}
	if err := repo.ReplaceLinkCandidates(ctx, installation, identities, candidates); err != nil {
		t.Fatalf("ReplaceLinkCandidates() error = %v", err)
	}

	summary, err := productlinkspostgres.NewSummaryReader(pool, summaryTenant).GetLinkageSummary(ctx, installation)
	if err != nil {
		t.Fatalf("GetLinkageSummary() error = %v", err)
	}
	if summary.PendingLinks != 2 {
		t.Fatalf("pending_links = %d, want 2: the confirmation and the review are both waiting on the operator", summary.PendingLinks)
	}
}
