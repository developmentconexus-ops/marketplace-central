package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// TariffRequest carries the inputs needed to resolve one tariff. Degrau 4
// only needs the modalidade; later degraus (live ML API, per-listing
// overrides) will add fields here as they land — do not add them now.
type TariffRequest struct {
	Modalidade domain.Modalidade
}

// TariffResolver resolves the commission + shipping tariff components for a
// TariffRequest. Degrau 4 (adapters/tariffdefaults) implements this purely
// from tenant config defaults; a future composite resolver will chain
// degraus 1->2->3->4 behind the same port with no interface change.
type TariffResolver interface {
	Resolve(ctx context.Context, req TariffRequest) (domain.TariffResolution, error)
}
