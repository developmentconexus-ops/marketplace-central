package application

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/orders/domain"
)

// EnrichedOrder wraps a canonical read model with derived, read-time-only
// enrichment facts. It is an internal application value — it is never
// marshaled directly onto the HTTP transport contract or the SDK.
type EnrichedOrder struct {
	Order         domain.OrderReadModel
	Buyer         domain.MaskedBuyer
	VinculoStatus domain.VinculoStatus
}

// EnrichService adds read-time enrichment (masked buyer identity, vinculo
// status) to canonical order read models. Every fact it derives already
// exists on the read model or its items; it makes no external ERP/ML/
// marketplace calls, and tenant scoping remains the repo layer's
// responsibility — installationID is accepted for future enrichment
// concerns but drives no query here.
type EnrichService struct{}

func NewEnrichService() EnrichService {
	return EnrichService{}
}

func (s EnrichService) Enrich(ctx context.Context, installationID string, orders []domain.OrderReadModel) []EnrichedOrder {
	enriched := make([]EnrichedOrder, 0, len(orders))
	for _, order := range orders {
		var nickname string
		if order.BuyerNickname != nil {
			nickname = *order.BuyerNickname
		}
		enriched = append(enriched, EnrichedOrder{
			Order:         order,
			Buyer:         domain.MaskBuyer(nickname),
			VinculoStatus: domain.DeriveVinculoStatus(order.Items),
		})
	}
	return enriched
}
