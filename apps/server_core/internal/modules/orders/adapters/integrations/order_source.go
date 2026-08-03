package integrations

import (
	"context"
	"strconv"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
)

// ProviderOrderReader is the boundary this module owns onto the connectors capability. It takes
// connectorsdomain.ListOrdersInput (Cursor/Limit/UpdatedAfter — the window, per F-00 Task 1) plus
// installationID as a SEPARATE parameter: AccountRef is resolved from installationID by
// ProviderOperationService (F-00 Task 2), never built here, so this module never learns the
// installation->account mapping. ListOrderRefs is the id-only enumeration (F-00 Task 3): the old
// ListOrders reply hydrated every hit with a per-order provider call that orders.ImportService
// then discarded except for the id, so this module now asks for exactly what it consumes.
type ProviderOrderReader interface {
	ListOrderRefs(ctx context.Context, input connectorsdomain.ListOrdersInput, installationID string) ([]connectorsdomain.OrderSearchHit, error)
}

type OrderSource struct {
	reader ProviderOrderReader
}

func NewOrderSource(reader ProviderOrderReader) OrderSource {
	return OrderSource{reader: reader}
}

func (s OrderSource) ListOrders(ctx context.Context, input ordersports.ListOrdersInput) ([]ordersports.OrderRef, error) {
	hits, err := s.reader.ListOrderRefs(ctx, connectorsdomain.ListOrdersInput{
		Cursor:       strconv.Itoa(input.Offset),
		Limit:        input.Limit,
		UpdatedAfter: input.UpdatedAfter,
	}, input.InstallationID)
	if err != nil {
		return nil, err
	}

	refs := make([]ordersports.OrderRef, 0, len(hits))
	for _, hit := range hits {
		refs = append(refs, ordersports.OrderRef{
			ProviderOrderID:   hit.ProviderOrderID,
			ProviderUpdatedAt: hit.ProviderUpdatedAt,
		})
	}

	return refs, nil
}
