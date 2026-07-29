//go:build integration

package postgres_test

import (
	"context"
	"sort"
	"testing"
	"time"

	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func candidateFixture(installation, itemID, candidateID string, now time.Time) productlinksdomain.LinkCandidate {
	return productlinksdomain.LinkCandidate{
		CandidateID:    candidateID,
		InstallationID: installation,
		ProviderCode:   "mercado_livre",
		ProviderItemID: itemID,
		State:          productlinksdomain.LinkCandidateStateExactSKU,
		MatchInput:     productlinksdomain.LinkCandidateMatchInputSellerSKU,
		MatchValue:     "SKU-" + itemID,
		Confidence:     70,
		ConfidenceBand: productlinksdomain.LinkCandidateConfidenceBandMedia,
		MatchStatus:    productlinksdomain.LinkCandidateMatchStatusConfirm,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// A capped generation run only clears the identities it processed, so an anúncio
// beyond the cap kept the previous run's candidate and the fila rendered it as
// the current answer. The run said nothing about that listing; presenting a
// stale row as fresh is the ADR-17 violation. A listing the operator already
// decided is NOT swept — it is no longer pending.
func TestReplaceLinkCandidatesClearsPendingCandidatesLeftOutOfTheRun(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_product_links_stale")
	token := time.Now().UTC().Format("150405.000000000")
	tenant := "stale-" + token
	installation := "stale-inst-" + token
	now := time.Date(2026, 7, 28, 12, 0, 0, 0, time.UTC)

	t.Cleanup(func() {
		for _, table := range []string{"product_link_candidates", "product_links"} {
			_, _ = pool.Exec(context.Background(), "DELETE FROM "+table+" WHERE tenant_id=$1 AND installation_id=$2", tenant, installation)
		}
	})

	repo := productlinkspostgres.NewLinkCandidateRepository(pool, tenant)

	identity := func(itemID string) productlinksdomain.ListingIdentity {
		return productlinksdomain.ListingIdentity{InstallationID: installation, ProviderItemID: itemID}
	}

	firstRun := []productlinksdomain.ListingIdentity{identity("MLB-in-cap"), identity("MLB-out-pending"), identity("MLB-out-decided")}
	firstCandidates := []productlinksdomain.LinkCandidate{
		candidateFixture(installation, "MLB-in-cap", "cand-in-cap-v1", now),
		candidateFixture(installation, "MLB-out-pending", "cand-out-pending-v1", now),
		candidateFixture(installation, "MLB-out-decided", "cand-out-decided-v1", now),
	}
	if err := repo.ReplaceLinkCandidates(ctx, installation, firstRun, firstCandidates); err != nil {
		t.Fatalf("seed first generation run: %v", err)
	}

	// The operator resolved one of the anúncios the second run will not revisit.
	if _, err := pool.Exec(ctx, `
		INSERT INTO product_links(tenant_id, installation_id, provider_code, provider_item_id, provider_variation_id, state)
		VALUES($1, $2, 'mercado_livre', 'MLB-out-decided', '', 'resolved')
	`, tenant, installation); err != nil {
		t.Fatalf("seed decided link: %v", err)
	}

	secondRun := []productlinksdomain.ListingIdentity{identity("MLB-in-cap")}
	secondCandidates := []productlinksdomain.LinkCandidate{
		candidateFixture(installation, "MLB-in-cap", "cand-in-cap-v2", now.Add(time.Hour)),
	}
	if err := repo.ReplaceLinkCandidates(ctx, installation, secondRun, secondCandidates); err != nil {
		t.Fatalf("capped generation run: %v", err)
	}

	rows, err := pool.Query(ctx, `
		SELECT candidate_id FROM product_link_candidates WHERE tenant_id=$1 AND installation_id=$2
	`, tenant, installation)
	if err != nil {
		t.Fatalf("read candidates: %v", err)
	}
	defer rows.Close()
	var stored []string
	for rows.Next() {
		var candidateID string
		if err := rows.Scan(&candidateID); err != nil {
			t.Fatalf("scan candidate: %v", err)
		}
		stored = append(stored, candidateID)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate candidates: %v", err)
	}
	sort.Strings(stored)

	want := []string{"cand-in-cap-v2", "cand-out-decided-v1"}
	if len(stored) != len(want) {
		t.Fatalf("stored candidates = %v, want %v", stored, want)
	}
	for index, candidateID := range want {
		if stored[index] != candidateID {
			t.Fatalf("stored candidates = %v, want %v", stored, want)
		}
	}
}
