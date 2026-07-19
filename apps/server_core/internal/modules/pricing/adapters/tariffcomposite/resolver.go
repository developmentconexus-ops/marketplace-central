// Package tariffcomposite chains the tariff ladder behind the single
// ports.TariffResolver port: it resolves the degrau-4 config defaults, then
// upgrades the commission with a degrau-3 live COTAÇÃO when one is available.
// Frete stays degrau 4 in this scope (degrau 3 quotes commission only).
package tariffcomposite

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// commissionResolver is the degrau-3 assembler seam: it resolves a live
// commission component, or reports a miss (found=false). Satisfied by
// tarifflive.Resolver.
type commissionResolver interface {
	ResolveCommission(ctx context.Context, req ports.TariffRequest) (domain.ComponentResolution, bool, error)
}

// Resolver composes degrau 3 over a degrau-4 base, both behind
// ports.TariffResolver.
type Resolver struct {
	base ports.TariffResolver
	live commissionResolver
}

var _ ports.TariffResolver = (*Resolver)(nil)

// NewResolver builds a composite that upgrades base's commission with live's
// degrau-3 quote when present.
func NewResolver(base ports.TariffResolver, live commissionResolver) *Resolver {
	return &Resolver{base: base, live: live}
}

// Resolve resolves the degrau-4 base then, if degrau 3 quotes a commission,
// swaps it in. A degrau-4 error is a real failure and propagates. A degrau-3
// miss OR error falls silently to the degrau-4 commission: a live-probe
// failure must never break a decompose and never fabricates a fact (ADR-17) —
// it degrades to the honest, lower-provenance config default.
func (r *Resolver) Resolve(ctx context.Context, req ports.TariffRequest) (domain.TariffResolution, error) {
	base, err := r.base.Resolve(ctx, req)
	if err != nil {
		return domain.TariffResolution{}, err
	}

	comissao, ok, liveErr := r.live.ResolveCommission(ctx, req)
	if liveErr != nil || !ok {
		return base, nil
	}

	base.Comissao = comissao
	return base, nil
}
