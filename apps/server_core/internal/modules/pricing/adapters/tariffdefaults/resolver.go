// Package tariffdefaults implements degrau 4 of the tariff resolver: it
// resolves commission + shipping tariff components purely from a
// tenant/installation's TariffDefaults config (Slice A), with no live or
// per-listing data. It implements ports.TariffResolver so a future composite
// resolver can chain degraus 1->2->3->4 behind the same port.
package tariffdefaults

import (
	"context"
	"fmt"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// degrau4 is this adapter's fixed position in the resolver chain.
const degrau4 = 4

// Resolver is the degrau-4 (config-defaults) implementation of
// ports.TariffResolver.
type Resolver struct {
	store          ports.TariffDefaultsStore
	tenantID       string
	installationID string
}

var _ ports.TariffResolver = (*Resolver)(nil)

// NewResolver builds a degrau-4 resolver bound to one tenant/installation.
func NewResolver(store ports.TariffDefaultsStore, tenantID, installationID string) *Resolver {
	return &Resolver{store: store, tenantID: tenantID, installationID: installationID}
}

// Resolve loads the tenant's TariffDefaults and maps them onto the requested
// modalidade's commission rate plus the configured frete estimate (or
// NO-DATA per ADR-17 — never a fabricated 0).
func (r *Resolver) Resolve(ctx context.Context, req ports.TariffRequest) (domain.TariffResolution, error) {
	defaults, err := r.store.GetTariffDefaults(ctx, r.tenantID, r.installationID)
	if err != nil {
		return domain.TariffResolution{}, fmt.Errorf("tariffdefaults: get tariff defaults: %w", err)
	}

	comissaoPct := defaults.ComissaoClassicoPct
	if req.Modalidade == domain.ModalidadePremium || req.Modalidade == domain.ModalidadeFull {
		comissaoPct = defaults.ComissaoPremiumPct
	}

	comissao := domain.ComponentResolution{
		Valor:      &comissaoPct,
		Fonte:      domain.FontePadrao,
		Degrau:     degrau4,
		Estimativa: true,
	}

	frete := domain.ComponentResolution{
		Valor:      nil,
		Fonte:      domain.FontePadrao,
		Degrau:     degrau4,
		Estimativa: true,
	}
	if defaults.FretePolicy == domain.FretePolicyEstimativa && defaults.FreteEstimativaAmount != nil {
		frete.Valor = defaults.FreteEstimativaAmount
	}

	return domain.TariffResolution{Comissao: comissao, Frete: frete}, nil
}
