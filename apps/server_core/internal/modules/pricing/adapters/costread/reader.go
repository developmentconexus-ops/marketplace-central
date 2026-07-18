// Package costread adapts IC-02 internal_read's cost facts to the pricing
// engine's decimal cost port.
package costread

import (
	"context"
	"strconv"
	"time"

	internalreaddomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/pricing/domain"
)

// factReader is the narrow slice of internal_read this adapter needs (mirrors
// profitability's InternalFactReader). Keeping it local decouples pricing from
// the full internal_read port surface.
type factReader interface {
	GetCostAsOf(ctx context.Context, productID int, effectiveAt time.Time) (internalreaddomain.CostAsOf, error)
}

// Reader maps internal_read CostAsOf to the pricing CostReader port.
type Reader struct {
	facts factReader
}

func NewReader(facts factReader) Reader {
	return Reader{facts: facts}
}

// CostForProduct returns the product's custo_erp as decimal Money, or nil when
// the fact is absent. internal_read exposes the amount as *float64 (IC-02's
// contract) — this adapter is the single boundary where that float is converted
// to the pricing module's decimal-string money, at 2dp. A nil amount is
// propagated as nil Money (ADR-17: never 0).
func (r Reader) CostForProduct(ctx context.Context, productID int, asOf time.Time) (*domain.Money, error) {
	cost, err := r.facts.GetCostAsOf(ctx, productID, asOf)
	if err != nil {
		return nil, err
	}
	if cost.Amount == nil {
		return nil, nil
	}
	return &domain.Money{Amount: formatCostAmount(*cost.Amount), Currency: "BRL"}, nil
}

// formatCostAmount converts the upstream float64 cost to a fixed 2dp decimal
// string. This is the only float64→decimal crossing in the pricing path; it is
// confined to the IC-02 boundary and never re-enters as float.
func formatCostAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
