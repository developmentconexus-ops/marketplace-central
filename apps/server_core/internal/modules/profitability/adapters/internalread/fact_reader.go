package internalread

import (
	"context"
	"time"

	internalreadapp "marketplace-central/apps/server_core/internal/modules/internal_read/application"
	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type FactReader struct {
	service internalreadapp.Service
}

func NewFactReader(service internalreadapp.Service) FactReader {
	return FactReader{service: service}
}

func (r FactReader) GetCostAsOf(ctx context.Context, productID int, effectiveAt time.Time) (internalreaddomain.CostAsOf, error) {
	return r.service.GetCostAsOf(ctx, internalreadports.CostAsOfInput{
		ProductID: productID,
		Policy: internalreaddomain.CostAsOfPolicy{
			CompanyID:   1,
			EffectiveAt: effectiveAt,
			Basis:       internalreaddomain.CostBasisCUSSEMICM,
		},
		Freshness: internalreaddomain.FreshnessPolicy{},
	})
}

func (r FactReader) GetTaxInputs(ctx context.Context, productID int, effectiveAt time.Time, source internalreaddomain.TaxSourceIdentity) (internalreaddomain.TaxInputs, error) {
	policy := internalreaddomain.DefaultTaxPolicy(effectiveAt)
	policy.Source = source
	return r.service.GetTaxInputs(ctx, internalreadports.TaxInput{
		ProductID: productID,
		Policy:    policy,
		Freshness: internalreaddomain.FreshnessPolicy{},
	})
}
