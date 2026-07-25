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
		_, _ = pool.Exec(ctx, `DELETE FROM product_links WHERE tenant_id = $1`, summaryTenant)
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
		{InstallationID: installation, ProviderItemID: "MLB-ANSWERED"},
		{InstallationID: installation, ProviderItemID: "MLB-KILLED"},
	}
	candidates := []productlinksdomain.LinkCandidate{
		candidate("cand-confirm", "MLB-CONFIRM", productlinksdomain.LinkCandidateStateExactEAN, productlinksdomain.LinkCandidateMatchStatusConfirm),
		candidate("cand-review", "MLB-REVIEW", productlinksdomain.LinkCandidateStateUnresolved, productlinksdomain.LinkCandidateMatchStatusReview),
		// Approving or rejecting settles product_links and leaves the candidate
		// row exactly where it was; regeneration then rebuilds it from the same
		// anchors. A count that reads candidates alone never goes back down —
		// the operator answers the queue and the number does not move.
		candidate("cand-answered", "MLB-ANSWERED", productlinksdomain.LinkCandidateStateExactEAN, productlinksdomain.LinkCandidateMatchStatusConfirm),
		candidate("cand-killed", "MLB-KILLED", productlinksdomain.LinkCandidateStateUnresolved, productlinksdomain.LinkCandidateMatchStatusReview),
	}
	if err := repo.ReplaceLinkCandidates(ctx, installation, identities, candidates); err != nil {
		t.Fatalf("ReplaceLinkCandidates() error = %v", err)
	}

	for _, settled := range []struct {
		itemID string
		state  productlinksdomain.ProductLinkState
	}{
		{"MLB-ANSWERED", productlinksdomain.ProductLinkStateResolved},
		{"MLB-KILLED", productlinksdomain.ProductLinkStateRejected},
	} {
		if _, err := pool.Exec(ctx, `
			INSERT INTO product_links (
				tenant_id, installation_id, provider_code, provider_item_id, provider_variation_id,
				state, created_at, updated_at
			) VALUES ($1, $2, 'mercado_livre', $3, '', $4, $5, $5)
		`, summaryTenant, installation, settled.itemID, settled.state, at); err != nil {
			t.Fatalf("seed settled link %s: %v", settled.itemID, err)
		}
	}

	summary, err := productlinkspostgres.NewSummaryReader(pool, summaryTenant).GetLinkageSummary(ctx, installation)
	if err != nil {
		t.Fatalf("GetLinkageSummary() error = %v", err)
	}
	if summary.PendingLinks != 2 {
		t.Fatalf("pending_links = %d, want 2: the unanswered confirmation and review count, the answered and the rejected do not", summary.PendingLinks)
	}
}
