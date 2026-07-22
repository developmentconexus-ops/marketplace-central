package erpproductlinks

import (
	"context"
	"fmt"

	erpapp "marketplace-central/apps/server_core/internal/modules/erp_import/application"
	productlinksapp "marketplace-central/apps/server_core/internal/modules/product_links/application"
	productlinksdomain "marketplace-central/apps/server_core/internal/modules/product_links/domain"
)

type candidateEngine interface {
	GenerateLinkCandidates(ctx context.Context, input productlinksapp.GenerateLinkCandidatesInput) (productlinksdomain.LinkCandidateGenerationResult, error)
}

type LinkCandidateGenerator struct {
	engine candidateEngine
}

func NewLinkCandidateGenerator(engine candidateEngine) *LinkCandidateGenerator {
	return &LinkCandidateGenerator{engine: engine}
}

func (g *LinkCandidateGenerator) GenerateLinkCandidates(ctx context.Context, installationID string) error {
	_, err := g.engine.GenerateLinkCandidates(ctx, productlinksapp.GenerateLinkCandidatesInput{
		InstallationID: installationID,
		Limit:          0,
	})
	if err != nil {
		return fmt.Errorf("generate link candidates for installation %q: %w", installationID, err)
	}
	return nil
}

var _ erpapp.LinkCandidateGenerator = (*LinkCandidateGenerator)(nil)
