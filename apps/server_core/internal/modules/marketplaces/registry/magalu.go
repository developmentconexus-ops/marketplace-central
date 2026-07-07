package registry

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/marketplaces/domain"
)

func init() { register(&MagaluPlugin{}) }

type MagaluPlugin struct{}

func (p *MagaluPlugin) Code() string { return "magalu" }

func (p *MagaluPlugin) Definition() domain.MarketplaceDefinition {
	return domain.MarketplaceDefinition{
		MarketplaceCode: "magalu",
		DisplayName:     "Magalu",
		FeeSource:       "seed",
		AuthStrategy:    "oauth2",
		CapabilityProfile: domain.CapabilityProfile{
			ListingRead:   domain.CapabilityBlocked,
			StockRead:     domain.CapabilityBlocked,
			StockWrite:    domain.CapabilityBlocked,
			OrderRead:     domain.CapabilityBlocked,
			Messages:      domain.CapabilityBlocked,
			Questions:     domain.CapabilityBlocked,
			FreightQuotes: domain.CapabilityUnsupported,
			Webhooks:      domain.CapabilityBlocked,
			Sandbox:       domain.CapabilityBlocked,
		},
		Metadata: domain.PluginMetadata{
			RolloutStage:  "v1",
			ExecutionMode: "live",
		},
		CredentialSchema: []domain.CredentialField{
			{Key: "client_id", Label: "Client ID", Secret: false},
			{Key: "client_secret", Label: "Client Secret", Secret: true},
		},
		Active: true,
	}
}

// SeedFees is a no-op — Magalu fees are seeded by connectors/adapters/magalu.FeeSyncer.
func (p *MagaluPlugin) SeedFees(_ context.Context, _ *pgxpool.Pool) error { return nil }

func (p *MagaluPlugin) NewConnector(_ map[string]string) (MarketplaceConnector, error) {
	return nil, ErrNotImplemented
}
