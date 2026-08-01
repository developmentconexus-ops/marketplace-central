// Package ports declares the channel_fees core interfaces (IC-01). M-03..M-08
// consumers depend on these interfaces, never on the Postgres implementation
// in adapters/postgres directly (feature.md Expected Output).
package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/channelfees/domain"
)

// ChannelFeeWriter upserts a channel_fees row on its natural key (IC-01
// Persistence Expectations): a second UpsertFee for the same subject
// collapses to the same row with value/detail/coletado_em/source_time
// advanced — never a duplicate insert.
type ChannelFeeWriter interface {
	UpsertFee(ctx context.Context, fee domain.Fee) error
}

// ChannelFeeReader resolves the current commission and freight fee for one
// listing, ledger-only (IC-01 Enums And Statuses): commission ladder is
// layer 2 → layer 1 → typed-absent; freight ladder is layer 2 →
// typed-absent (freight has no layer 1). Layer 3 never participates in
// either ladder — this is a pricing-read seam, not the realized-fee record.
// The `config`/pricing_tariff_defaults fallback step belongs to the M-07
// pricing consumer's composition, never to this reader (auditoria P5 F-10).
type ChannelFeeReader interface {
	ResolveListingFees(ctx context.Context, tenantID, provider, installationID, providerListingID string) (domain.ResolvedListingFees, error)
}
