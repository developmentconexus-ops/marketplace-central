package composition_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	productlinkscomposition "marketplace-central/apps/server_core/internal/modules/product_links/composition"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type stubInstallations struct {
	installations []integrationsdomain.Installation
	err           error
}

func (s stubInstallations) List(ctx context.Context) ([]integrationsdomain.Installation, error) {
	return s.installations, s.err
}

type stubEngine struct {
	seen   []string
	failOn map[string]error
}

func (e *stubEngine) GenerateLinkCandidates(ctx context.Context, input productlinksapp.GenerateLinkCandidatesInput) (productlinksdomain.LinkCandidateGenerationResult, error) {
	e.seen = append(e.seen, input.InstallationID)
	if err, ok := e.failOn[input.InstallationID]; ok {
		return productlinksdomain.LinkCandidateGenerationResult{}, err
	}
	return productlinksdomain.LinkCandidateGenerationResult{}, nil
}

func installations(ids ...string) []integrationsdomain.Installation {
	out := make([]integrationsdomain.Installation, 0, len(ids))
	for _, id := range ids {
		out = append(out, integrationsdomain.Installation{InstallationID: id})
	}
	return out
}

// M05-C11: generation runs for every installation of the tenant, unbounded
// (Limit 0) — a partial regeneration would leave anúncios silently unmatched.
func TestRefreshGeneratesForEveryInstallation(t *testing.T) {
	engine := &stubEngine{}
	refresher := productlinkscomposition.NewLinkCandidateRefresher(
		stubInstallations{installations: installations("inst-a", "inst-b")}, engine)

	if err := refresher.RefreshLinkCandidates(context.Background()); err != nil {
		t.Fatal(err)
	}
	if strings.Join(engine.seen, ",") != "inst-a,inst-b" {
		t.Fatalf("generated for %v want both installations", engine.seen)
	}
}

// One installation failing must not skip the others: the failures are joined
// and reported after everything that could run, did.
func TestRefreshKeepsGoingAfterOneInstallationFails(t *testing.T) {
	engine := &stubEngine{failOn: map[string]error{"inst-a": errors.New("matcher exploded")}}
	refresher := productlinkscomposition.NewLinkCandidateRefresher(
		stubInstallations{installations: installations("inst-a", "inst-b")}, engine)

	err := refresher.RefreshLinkCandidates(context.Background())
	if err == nil {
		t.Fatal("want the installation failure reported")
	}
	if !strings.Contains(err.Error(), "inst-a") {
		t.Fatalf("error %q does not name the failing installation", err)
	}
	if len(engine.seen) != 2 {
		t.Fatalf("generated for %v want the second installation attempted anyway", engine.seen)
	}
}

// An unreadable installation list is an error, never an empty successful
// refresh: "nothing to match" and "we could not tell" are different facts.
func TestRefreshFailsWhenInstallationsCannotBeListed(t *testing.T) {
	engine := &stubEngine{}
	refresher := productlinkscomposition.NewLinkCandidateRefresher(
		stubInstallations{err: errors.New("db down")}, engine)

	if err := refresher.RefreshLinkCandidates(context.Background()); err == nil {
		t.Fatal("want the listing failure surfaced")
	}
	if len(engine.seen) != 0 {
		t.Fatalf("generated for %v want nothing", engine.seen)
	}
}
