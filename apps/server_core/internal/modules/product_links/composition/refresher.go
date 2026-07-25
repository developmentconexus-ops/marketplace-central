package composition

import (
	"context"
	"errors"
	"fmt"

	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type installationLister interface {
	List(ctx context.Context) ([]integrationsdomain.Installation, error)
}

type candidateEngine interface {
	GenerateLinkCandidates(ctx context.Context, input productlinksapp.GenerateLinkCandidatesInput) (productlinksdomain.LinkCandidateGenerationResult, error)
}

// LinkCandidateRefresher regenerates link candidates for every installation of
// the tenant. It exists because product data changing is exactly when the
// answer to "which product is this anúncio?" can change: an ERP sync that
// brings in the codprod an anúncio names must produce the candidate for it
// without an operator remembering to press a button.
type LinkCandidateRefresher struct {
	installations installationLister
	engine        candidateEngine
}

func NewLinkCandidateRefresher(installations installationLister, engine candidateEngine) *LinkCandidateRefresher {
	return &LinkCandidateRefresher{installations: installations, engine: engine}
}

// RefreshLinkCandidates runs generation per installation. One installation
// failing does not skip the others — the failures are joined and returned so
// the caller can report them, having already tried everything it could.
func (r *LinkCandidateRefresher) RefreshLinkCandidates(ctx context.Context) error {
	if r == nil || r.installations == nil || r.engine == nil {
		return errors.New("PRODUCT_LINKS_REFRESHER_NOT_CONFIGURED")
	}
	installations, err := r.installations.List(ctx)
	if err != nil {
		return fmt.Errorf("list installations for link candidate refresh: %w", err)
	}
	var failures []error
	for _, installation := range installations {
		if _, err := r.engine.GenerateLinkCandidates(ctx, productlinksapp.GenerateLinkCandidatesInput{
			InstallationID: installation.InstallationID,
			Limit:          0,
		}); err != nil {
			failures = append(failures, fmt.Errorf("generate link candidates for installation %q: %w", installation.InstallationID, err))
		}
	}
	return errors.Join(failures...)
}
