package registry

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/marketplaces/domain"
)

func init() { register(&LeroyPlugin{}) }

type LeroyPlugin struct{}

func (p *LeroyPlugin) Code() string { return "leroy_merlin" }

func (p *LeroyPlugin) Definition() domain.MarketplaceDefinition {
	return domain.MarketplaceDefinition{
		MarketplaceCode: "leroy_merlin",
		DisplayName:     "Leroy Merlin",
		FeeSource:       "seed",
		AuthStrategy:    "api_key",
		CapabilityProfile: domain.CapabilityProfile{
			ListingRead:   domain.CapabilityBlocked,
			StockRead:     domain.CapabilityBlocked,
			StockWrite:    domain.CapabilityBlocked,
			OrderRead:     domain.CapabilityBlocked,
			Messages:      domain.CapabilityBlocked,
			Questions:     domain.CapabilityBlocked,
			FreightQuotes: domain.CapabilityUnsupported,
			Webhooks:      domain.CapabilityBlocked,
			Sandbox:       domain.CapabilityUnsupported,
		},
		Metadata: domain.PluginMetadata{
			RolloutStage:  "wave_2",
			ExecutionMode: "live",
		},
		CredentialSchema: []domain.CredentialField{
			{Key: "api_key", Label: "API Key", Secret: true},
			{Key: "shop_id", Label: "Shop ID", Secret: false},
		},
		Active: true,
	}
}

// SeedFees inserts a stub default fee row for Leroy Merlin (Mirakl Seller API).
func (p *LeroyPlugin) SeedFees(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO marketplace_fee_schedules
			(marketplace_code, category_id, listing_type, commission_percent, fixed_fee_amount, notes, source, synced_at)
		VALUES ('leroy_merlin', 'default', NULL, 0.18, 0, 'stub — to be filled with official per-category rates', 'seed', NOW())
		ON CONFLICT DO NOTHING
	`)
	return err
}

func (p *LeroyPlugin) NewConnector(_ map[string]string) (MarketplaceConnector, error) {
	return nil, ErrNotImplemented
}
