package ports

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// IC-04 FROZEN PORT SURFACE.
//
// These signatures are the single decomposition contract that M-07 (Simulador)
// AND M-08 (Pedidos) both consume. There is exactly ONE decomposition
// implementation (domain.Decompose) — M-08 reads it through this port, it does
// not re-implement the formula. Do not widen or reshape these signatures
// without a contract amendment; calc_ports_contract_test.go pins them.

// CalcEngine is the decomposition port. Decompose is pure: its input carries
// the already-resolved decimal components (comissão, custo, frete, efetivo),
// which adapters resolve upstream (S6). It never returns an error — invalid
// input is rejected at the transport boundary (S7) before reaching here.
type CalcEngine interface {
	Decompose(input domain.DecomposeInput) domain.Decomposition
}

// DifalReader resolves a tenant's DifalRate for a destino UF, reflecting any
// active operator override (tenant_id scopes every read). Implemented by the
// pricing difal repository adapter (S6).
type DifalReader interface {
	RateForUF(ctx context.Context, tenantID, uf string) (domain.DifalRate, error)
}

// DifalForUF is the frozen IC-04 read port DifalForUF(uf) → {efetivo_pct,
// versao}: it resolves the tenant's rate for uf (an active override takes
// precedence over the seed) and returns its efetivo_pct + versao. M-08
// consumes this to price order lines against the same DIFAL numbers the
// Simulador shows. The efetivo is computed at 2dp from exact
// interna−interestadual inside domain.DifalForUF, preserving operator
// overrides set at 2dp.
func DifalForUF(ctx context.Context, reader DifalReader, tenantID, uf string) (domain.DifalForUFResult, error) {
	rate, err := reader.RateForUF(ctx, tenantID, uf)
	if err != nil {
		return domain.DifalForUFResult{}, err
	}
	return domain.DifalForUF(rate), nil
}
