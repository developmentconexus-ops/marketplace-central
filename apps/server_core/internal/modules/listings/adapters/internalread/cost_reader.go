package internalread

import (
	"context"

	internaldomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	internalports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	listingports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type CostReader struct{ reader internalports.BatchReader }

func NewCostReader(reader internalports.BatchReader) CostReader { return CostReader{reader: reader} }
func (a CostReader) GetCostFactsByIDs(ctx context.Context, ids []int64) (map[int64]*listingports.CostFact, error) {
	facts, err := a.reader.GetCostFactsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make(map[int64]*listingports.CostFact, len(facts))
	for id, f := range facts {
		if f != nil {
			out[id] = &listingports.CostFact{Amount: f.Amount}
		}
	}
	return out, nil
}
func (a CostReader) GetICMSCeilingByOrigin(ctx context.Context, origin int64) (map[int64]*listingports.ICMSCeiling, error) {
	facts, err := a.reader.GetICMSCeilingByOrigin(ctx, internaldomain.UF(origin))
	if err != nil {
		return nil, err
	}
	out := make(map[int64]*listingports.ICMSCeiling, len(facts))
	for uf, f := range facts {
		if f != nil {
			out[int64(uf)] = &listingports.ICMSCeiling{Percent: f.CeilingPercent}
		}
	}
	return out, nil
}

var _ listingports.CostReader = CostReader{}
