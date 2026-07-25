//go:build integration

package integration_test

import (
	"context"
	"strings"
	"testing"
	"time"

	productlinkspostgres "marketplace-central/apps/server_core/internal/modules/product_links/adapters/postgres"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

const decisionsTenant = "tenant-e10"

func decisionsRepo(t *testing.T) (*productlinkspostgres.LinkCandidateRepository, *pgxpool.Pool) {
	t.Helper()
	testpostgres.SkipWithoutTarget(t)
	pool, _ := testpostgres.OpenPool(t, decisionsTenant)
	t.Cleanup(func() {
		ctx := context.Background()
		// Decisions first: a superseded row points at its successor.
		_, _ = pool.Exec(ctx, `UPDATE product_link_decisions SET superseded_by = NULL WHERE tenant_id = $1`, decisionsTenant)
		_, _ = pool.Exec(ctx, `DELETE FROM product_link_decisions WHERE tenant_id = $1`, decisionsTenant)
		_, _ = pool.Exec(ctx, `DELETE FROM product_link_audit_entries WHERE tenant_id = $1`, decisionsTenant)
		_, _ = pool.Exec(ctx, `DELETE FROM product_links WHERE tenant_id = $1`, decisionsTenant)
		pool.Close()
	})
	return productlinkspostgres.NewLinkCandidateRepository(pool, decisionsTenant), pool
}

func approvingTransition(installationID, itemID, variationID, decisionID string, rule productlinksdomain.ProductLinkDecisionRule, actor string, collisions *int, at time.Time) productlinksdomain.ProductLinkTransition {
	productID := 100002
	return productlinksdomain.ProductLinkTransition{
		Link: productlinksdomain.ProductLink{
			InstallationID:      installationID,
			ProviderCode:        "mercado_livre",
			ProviderItemID:      itemID,
			ProviderVariationID: variationID,
			State:               productlinksdomain.ProductLinkStateResolved,
			InternalProductID:   &productID,
			InternalProductName: "PUXADOR DHARMA 192MM PINT TIT",
			CreatedAt:           at,
			UpdatedAt:           at,
		},
		Audit: productlinksdomain.ProductLinkAuditEntry{
			AuditID:               "audit-" + decisionID,
			InstallationID:        installationID,
			ProviderCode:          "mercado_livre",
			ProviderItemID:        itemID,
			ProviderVariationID:   variationID,
			Action:                productlinksdomain.ProductLinkActionApproveCandidate,
			Actor:                 productlinksdomain.ActorMetadata{ActorType: actor, ActorID: actor},
			PreviousState:         productlinksdomain.ProductLinkStateNone,
			NextState:             productlinksdomain.ProductLinkStateResolved,
			NextInternalProductID: &productID,
			CreatedAt:             at,
		},
		Decision: &productlinksdomain.ProductLinkDecision{
			DecisionID:           decisionID,
			InstallationID:       installationID,
			ProviderItemID:       itemID,
			ProviderVariationID:  variationID,
			LinkID:               productlinksdomain.LinkID(installationID, itemID, variationID),
			RuleMatched:          rule,
			Actor:                actor,
			CollisionsAtDecision: collisions,
			CreatedAt:            at,
		},
	}
}

// M05-C3/M05-C12: the corroborated path records what it decided on — the rule,
// the actor, and how many products the anchors matched at that moment.
func TestCorroboratedApprovalWritesTheE10Row(t *testing.T) {
	repo, pool := decisionsRepo(t)
	ctx := context.Background()
	at := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	collisions := 1

	if err := repo.ApplyProductLinkTransition(ctx, approvingTransition("inst-c3", "MLB-C3", "", "dec-c3", productlinksdomain.DecisionRuleConcordantCodprodEAN, productlinksdomain.DecisionActorSystem, &collisions, at)); err != nil {
		t.Fatalf("ApplyProductLinkTransition() error = %v", err)
	}

	var rule, actor, linkID string
	var got int
	var superseded *string
	if err := pool.QueryRow(ctx, `
		SELECT rule_matched, actor, collisions_at_decision, link_id, superseded_by
		FROM product_link_decisions WHERE tenant_id = $1 AND decision_id = $2
	`, decisionsTenant, "dec-c3").Scan(&rule, &actor, &got, &linkID, &superseded); err != nil {
		t.Fatalf("read decision: %v", err)
	}
	if rule != "concordant_codprod_ean" || actor != "system" || got != 1 {
		t.Errorf("decision = (%s, %s, %d), want (concordant_codprod_ean, system, 1)", rule, actor, got)
	}
	if linkID != "inst-c3|MLB-C3|" {
		t.Errorf("link_id = %q — generated from the identity, so it cannot drift from the link", linkID)
	}
	if superseded != nil {
		t.Errorf("a first decision is in force, not superseded: %q", *superseded)
	}
}

// ADR-17: a decision taken without reading the collision count stores NULL. A 0
// there would read as "the anchor matched no product", a fact nobody established.
func TestAnUnreadCollisionCountStaysNull(t *testing.T) {
	repo, pool := decisionsRepo(t)
	ctx := context.Background()
	at := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)

	if err := repo.ApplyProductLinkTransition(ctx, approvingTransition("inst-null", "MLB-NULL", "", "dec-null", productlinksdomain.DecisionRuleManual, productlinksdomain.DecisionActorOperator, nil, at)); err != nil {
		t.Fatalf("ApplyProductLinkTransition() error = %v", err)
	}
	var collisions *int
	if err := pool.QueryRow(ctx, `SELECT collisions_at_decision FROM product_link_decisions WHERE tenant_id = $1 AND decision_id = $2`, decisionsTenant, "dec-null").Scan(&collisions); err != nil {
		t.Fatalf("read decision: %v", err)
	}
	if collisions != nil {
		t.Fatalf("collisions_at_decision = %d, want NULL", *collisions)
	}
}

// D-121-2 at the schema level: corroboration is the only rule the system may
// apply alone, so the database refuses an actor='system' row on any other.
func TestTheDatabaseRefusesASystemDecisionOnASingleAnchor(t *testing.T) {
	repo, _ := decisionsRepo(t)
	ctx := context.Background()
	at := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	collisions := 1

	err := repo.ApplyProductLinkTransition(ctx, approvingTransition("inst-guard", "MLB-GUARD", "", "dec-guard", productlinksdomain.DecisionRuleExactEANUnique, productlinksdomain.DecisionActorSystem, &collisions, at))
	if err == nil {
		t.Fatalf("a single-anchor decision by the system was accepted")
	}
	if !strings.Contains(err.Error(), "product_link_decisions_system_is_corroborated_check") {
		t.Fatalf("rejected for the wrong reason: %v", err)
	}
}

// M05-C10: an operator override appends and stamps. The system decision stays
// readable and is never written back over the operator's.
func TestAnOperatorOverrideSupersedesTheAutomaticDecision(t *testing.T) {
	repo, pool := decisionsRepo(t)
	ctx := context.Background()
	at := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	collisions := 1

	if err := repo.ApplyProductLinkTransition(ctx, approvingTransition("inst-c10", "MLB-C10", "", "dec-auto", productlinksdomain.DecisionRuleConcordantCodprodEAN, productlinksdomain.DecisionActorSystem, &collisions, at)); err != nil {
		t.Fatalf("automatic approval: %v", err)
	}
	if err := repo.ApplyProductLinkTransition(ctx, approvingTransition("inst-c10", "MLB-C10", "", "dec-manual", productlinksdomain.DecisionRuleManual, productlinksdomain.DecisionActorOperator, nil, at.Add(time.Minute))); err != nil {
		t.Fatalf("operator override: %v", err)
	}

	decisions, err := repo.ListDecisionsForLink(ctx, productlinksdomain.ListingIdentity{InstallationID: "inst-c10", ProviderItemID: "MLB-C10"})
	if err != nil {
		t.Fatalf("ListDecisionsForLink() error = %v", err)
	}
	if len(decisions) != 2 {
		t.Fatalf("decisions = %d, want 2 (the trail is append-only)", len(decisions))
	}
	if decisions[0].Actor != "system" || decisions[0].SupersededBy != "dec-manual" {
		t.Errorf("first decision = %+v, want the system row stamped by the override", decisions[0])
	}
	if decisions[1].Actor != "operator" || decisions[1].SupersededBy != "" {
		t.Errorf("second decision = %+v, want the operator row in force", decisions[1])
	}

	var links int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM product_links WHERE tenant_id = $1 AND installation_id = $2`, decisionsTenant, "inst-c10").Scan(&links); err != nil {
		t.Fatalf("count links: %v", err)
	}
	if links != 1 {
		t.Fatalf("product_links rows = %d, want 1 — two decisions on one link, not two links", links)
	}
}

// M05-C8/M05-C9: idempotency comes from the key product_links has carried since
// 0025. Re-running the same identity updates one row; a different variation of
// the SAME listing is a different link, which the key the plan originally asked
// for would have collapsed.
func TestListingIdentityKeepsApprovalIdempotentPerVariation(t *testing.T) {
	repo, pool := decisionsRepo(t)
	ctx := context.Background()
	at := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	collisions := 1

	for i, decisionID := range []string{"dec-run1", "dec-run2"} {
		if err := repo.ApplyProductLinkTransition(ctx, approvingTransition("inst-c8", "MLB-C8", "", decisionID, productlinksdomain.DecisionRuleConcordantCodprodEAN, productlinksdomain.DecisionActorSystem, &collisions, at.Add(time.Duration(i)*time.Minute))); err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
	}
	if err := repo.ApplyProductLinkTransition(ctx, approvingTransition("inst-c8", "MLB-C8", "VAR-2", "dec-var2", productlinksdomain.DecisionRuleConcordantCodprodEAN, productlinksdomain.DecisionActorSystem, &collisions, at)); err != nil {
		t.Fatalf("second variation: %v", err)
	}

	var sameIdentity, wholeListing int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM product_links
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_item_id = $3 AND provider_variation_id = ''
	`, decisionsTenant, "inst-c8", "MLB-C8").Scan(&sameIdentity); err != nil {
		t.Fatalf("count identity: %v", err)
	}
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM product_links
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_item_id = $3
	`, decisionsTenant, "inst-c8", "MLB-C8").Scan(&wholeListing); err != nil {
		t.Fatalf("count listing: %v", err)
	}
	if sameIdentity != 1 {
		t.Errorf("re-running the same identity produced %d links, want 1", sameIdentity)
	}
	if wholeListing != 2 {
		t.Errorf("links for the listing = %d, want 2 — each variation is its own link", wholeListing)
	}
}

// AC-09: a product may be linked to as many distinct listings as the operator
// listed it under. Nothing caps or flags that.
func TestOneProductMayLinkToManyListings(t *testing.T) {
	repo, pool := decisionsRepo(t)
	ctx := context.Background()
	at := time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC)
	collisions := 1

	for i, itemID := range []string{"MLB-N1", "MLB-N2", "MLB-N3"} {
		if err := repo.ApplyProductLinkTransition(ctx, approvingTransition("inst-c20", itemID, "", "dec-"+itemID, productlinksdomain.DecisionRuleConcordantCodprodEAN, productlinksdomain.DecisionActorSystem, &collisions, at.Add(time.Duration(i)*time.Minute))); err != nil {
			t.Fatalf("listing %s: %v", itemID, err)
		}
	}
	var links int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM product_links
		WHERE tenant_id = $1 AND installation_id = $2 AND internal_product_id = 100002 AND state = 'resolved'
	`, decisionsTenant, "inst-c20").Scan(&links); err != nil {
		t.Fatalf("count links: %v", err)
	}
	if links != 3 {
		t.Fatalf("active links for one product = %d, want 3", links)
	}
}
