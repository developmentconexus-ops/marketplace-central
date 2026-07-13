package orders

import (
	"context"

	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

type OrderReader struct {
	lister interface {
		List(ctx context.Context, input ordersapp.ListOrdersInput) ([]ordersdomain.MarketplaceOrder, error)
	}
}

func NewOrderReader(lister interface {
	List(ctx context.Context, input ordersapp.ListOrdersInput) ([]ordersdomain.MarketplaceOrder, error)
}) OrderReader {
	return OrderReader{lister: lister}
}

func (r OrderReader) ListOrders(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.OrderFact, error) {
	orders, err := r.lister.List(ctx, ordersapp.ListOrdersInput{InstallationID: installationID, Limit: limit})
	if err != nil {
		return nil, err
	}

	facts := make([]profitabilitydomain.OrderFact, len(orders))
	for orderIndex, order := range orders {
		items := make([]profitabilitydomain.OrderItemFact, len(order.Items))
		for itemIndex, item := range order.Items {
			items[itemIndex] = profitabilitydomain.OrderItemFact{
				MPCLineID:           string(item.MPCLineID),
				ReconciliationState: mapReconciliationState(item.ReconciliationState),
				ProviderItemID:      item.ProviderItemID,
				ProviderVariationID: item.ProviderVariationID,
				Quantity:            item.Quantity,
				UnitPrice:           item.UnitPrice,
				SaleFeeAmount:       item.SaleFeeAmount,
				LinkQuality:         mapLinkQuality(item.LinkQuality),
				InternalProductID:   item.InternalProductID,
			}
		}
		facts[orderIndex] = profitabilitydomain.OrderFact{
			InstallationID:    order.InstallationID,
			ProviderOrderID:   order.ProviderOrderID,
			RealizationState:  mapRealizationState(order.ProviderCode, order.ProviderStatus),
			ProviderCreatedAt: order.ProviderCreatedAt,
			ProviderClosedAt:  order.ProviderClosedAt,
			ProviderUpdatedAt: order.ProviderUpdatedAt,
			FetchedAt:         order.FetchedAt,
			Items:             items,
		}
	}
	return facts, nil
}

func mapReconciliationState(state ordersdomain.LineReconciliationState) profitabilitydomain.OrderLineReconciliationState {
	switch state {
	case ordersdomain.LineReconciliationStable:
		return profitabilitydomain.OrderLineStable
	case ordersdomain.LineReconciliationLegacyUnresolved:
		return profitabilitydomain.OrderLineLegacyUnresolved
	case ordersdomain.LineReconciliationAmbiguous:
		return profitabilitydomain.OrderLineAmbiguous
	default:
		return ""
	}
}

func mapRealizationState(providerCode, providerStatus string) profitabilitydomain.OrderRealizationState {
	if providerCode != "mercado_livre" {
		return profitabilitydomain.OrderRealizationUnknown
	}
	if providerStatus == "paid" {
		return profitabilitydomain.OrderRealizationRealized
	}
	if providerStatus == "cancelled" || providerStatus == "canceled" {
		return profitabilitydomain.OrderRealizationNotRealized
	}
	return profitabilitydomain.OrderRealizationUnknown
}

func mapLinkQuality(quality ordersdomain.LinkQuality) profitabilitydomain.OrderLinkQuality {
	switch quality {
	case ordersdomain.LinkQualityResolved:
		return profitabilitydomain.OrderLinkResolved
	case ordersdomain.LinkQualityRejected:
		return profitabilitydomain.OrderLinkRejected
	case ordersdomain.LinkQualityConflict:
		return profitabilitydomain.OrderLinkConflict
	case ordersdomain.LinkQualityUnresolved:
		return profitabilitydomain.OrderLinkUnresolved
	default:
		return profitabilitydomain.OrderLinkMissing
	}
}
