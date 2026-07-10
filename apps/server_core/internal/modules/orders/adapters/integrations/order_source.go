package integrations

import (
	"context"

	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
)

type ProviderOrderReader interface {
	ListOrders(ctx context.Context, installationID string, limit int) ([]connectorsdomain.OrderSnapshot, error)
}

type OrderSource struct {
	reader ProviderOrderReader
}

func NewOrderSource(reader ProviderOrderReader) OrderSource {
	return OrderSource{reader: reader}
}

func (s OrderSource) ListOrders(ctx context.Context, installationID string, limit int) ([]ordersdomain.OrderIngestionSnapshot, error) {
	snapshots, err := s.reader.ListOrders(ctx, installationID, limit)
	if err != nil {
		return nil, err
	}

	orders := make([]ordersdomain.OrderIngestionSnapshot, 0, len(snapshots))
	for _, snapshot := range snapshots {
		order := ordersdomain.OrderIngestionSnapshot{
			ProviderCode:         snapshot.ProviderCode,
			ProviderOrderID:      snapshot.ProviderOrderID,
			ProviderStatus:       snapshot.ProviderStatus,
			ProviderStatusDetail: snapshot.ProviderStatusDetail,
			ProviderCreatedAt:    snapshot.ProviderCreatedAt,
			ProviderClosedAt:     snapshot.ProviderClosedAt,
			ProviderUpdatedAt:    snapshot.ProviderUpdatedAt,
			FetchedAt:            snapshot.FetchedAt,
			SaleFeeAmount:        snapshot.SaleFeeAmount,
			ShippingID:           snapshot.ShippingID,
			CancellationDetail:   snapshot.CancellationDetail,
			Tags:                 snapshot.Tags,
		}
		for _, item := range snapshot.Items {
			order.Items = append(order.Items, ordersdomain.OrderIngestionItem{
				ProviderItemID:      item.ProviderItemID,
				ProviderVariationID: item.ProviderVariationID,
				SellerSKU:           item.SellerSKU,
				EAN:                 item.EAN,
				Title:               item.Title,
				Quantity:            item.Quantity,
				UnitPrice:           item.UnitPrice,
				SaleFeeAmount:       item.SaleFeeAmount,
			})
		}
		for _, payment := range snapshot.Payments {
			order.Payments = append(order.Payments, ordersdomain.OrderIngestionPayment{
				PaymentID:         payment.PaymentID,
				Status:            payment.Status,
				Amount:            payment.Amount,
				TransactionAmount: payment.TransactionAmount,
				TotalPaidAmount:   payment.TotalPaidAmount,
			})
		}
		orders = append(orders, order)
	}

	return orders, nil
}
