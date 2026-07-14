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
	batch   internalreadports.BatchReader
}

func NewFactReader(service internalreadapp.Service, batch ...internalreadports.BatchReader) FactReader {
	var batchReader internalreadports.BatchReader
	if len(batch) > 0 {
		batchReader = batch[0]
	}
	return FactReader{service: service, batch: batchReader}
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

func (r FactReader) GetCostFactsByIDs(ctx context.Context, ids []int64) (map[int64]*internalreaddomain.CostAsOf, error) {
	if r.batch == nil {
		return nil, internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "oracle cost batch reader is unavailable", nil)
	}
	return r.batch.GetCostFactsByIDs(ctx, ids)
}

func (r FactReader) GetTaxFactsByIDs(ctx context.Context, ids []int64) (map[int64]*internalreaddomain.TaxInputs, error) {
	if r.batch == nil {
		return nil, internalreaddomain.NewReadError(internalreaddomain.ReadErrorSourceUnavailable, "oracle tax batch reader is unavailable", nil)
	}
	return r.batch.GetTaxFactsByIDs(ctx, ids)
}
