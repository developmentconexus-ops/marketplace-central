package registry

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/marketplaces/domain"
)

func init() { register(&MercadoLivrePlugin{}) }

type MercadoLivrePlugin struct{}

func (p *MercadoLivrePlugin) Code() string { return "mercado_livre" }

func (p *MercadoLivrePlugin) Definition() domain.MarketplaceDefinition {
	return domain.MarketplaceDefinition{
		MarketplaceCode: "mercado_livre",
		DisplayName:     "Mercado Livre",
		FeeSource:       "api_sync",
		AuthStrategy:    "oauth2",
		CapabilityProfile: domain.CapabilityProfile{
			ListingRead:   domain.CapabilitySupported,
			StockRead:     domain.CapabilitySupported,
			StockWrite:    domain.CapabilitySupported,
			OrderRead:     domain.CapabilitySupported,
			Messages:      domain.CapabilityDegraded,
			Questions:     domain.CapabilitySupported,
			FreightQuotes: domain.CapabilityDegraded,
			Webhooks:      domain.CapabilitySupported,
			Sandbox:       domain.CapabilityBlocked,
		},
		Metadata: domain.PluginMetadata{
			RolloutStage:  "v1",
			ExecutionMode: "live",
		},
		CredentialSchema: []domain.CredentialField{
			{Key: "client_id", Label: "Client ID", Secret: false},
			{Key: "client_secret", Label: "Client Secret", Secret: true},
			{Key: "redirect_uri", Label: "Redirect URI", Secret: false},
		},
		Active: true,
	}
}

// SeedFees is a no-op: ML fees are read per listing from the Fees API (sale_fee,
// per unit), never seeded as a flat rate.
func (p *MercadoLivrePlugin) SeedFees(_ context.Context, _ *pgxpool.Pool) error { return nil }

func (p *MercadoLivrePlugin) NewConnector(_ map[string]string) (MarketplaceConnector, error) {
	return nil, ErrNotImplemented
}
