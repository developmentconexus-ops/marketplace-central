package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// TariffDefaultsStore is the persistence port for a tenant/installation's
// pricing tariff defaults (CHIP-T1 Slice A): classic/premium commission
// percentages and the shipping-estimate policy. Every method is
// tenant-scoped; installationID "" is the single-installation sentinel.
type TariffDefaultsStore interface {
	// GetTariffDefaults returns the stored row, materializing it on first
	// read (DB column DEFAULTs fill comissao_classico_pct/comissao_premium_pct
	// /frete_policy) so callers never see a zero row.
	GetTariffDefaults(ctx context.Context, tenantID, installationID string) (domain.TariffDefaults, error)
	// UpsertTariffDefaults validates in.FretePolicy (domain.ErrInvalidFretePolicy
	// otherwise) and persists the full row, returning it as stored.
	UpsertTariffDefaults(ctx context.Context, tenantID, installationID string, in domain.TariffDefaults) (domain.TariffDefaults, error)
}
