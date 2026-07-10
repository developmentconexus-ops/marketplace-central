package connectors

import (
	"context"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

type StockWriter struct {
	capabilities *connectorsapp.MarketplaceCapabilityService
}

func NewStockWriter(capabilities *connectorsapp.MarketplaceCapabilityService) StockWriter {
	return StockWriter{capabilities: capabilities}
}

func (w StockWriter) UpdateAvailableQuantity(ctx context.Context, request inventorydomain.StockWriteRequest) (inventorydomain.StockWriteResult, error) {
	writer, err := w.capabilities.StockWriter(request.ProviderRef.ProviderCode)
	if err != nil {
		return inventorydomain.StockWriteResult{}, err
	}
	result, err := writer.UpdateAvailableQuantity(ctx, connectorsdomain.StockWriteRequest{
		TenantID:            request.ProviderRef.TenantID,
		InstallationID:      request.ProviderRef.InstallationID,
		ProviderCode:        request.ProviderRef.ProviderCode,
		ProviderAccountID:   request.ProviderRef.ProviderAccountID,
		ProviderItemID:      request.ProviderRef.ProviderItemID,
		ProviderVariationID: request.ProviderRef.ProviderVariationID,
		IdempotencyKey:      request.IdempotencyKey,
		RequestedQuantity:   request.RequestedQuantity,
		Reason:              request.Reason,
	})
	if err != nil {
		return inventorydomain.StockWriteResult{}, err
	}
	return inventorydomain.StockWriteResult{
		Status:              mapStockWriteResult(result.Result),
		IdempotencyKey:      result.IDempotencyKey,
		ProviderCode:        result.ProviderCode,
		ProviderItemID:      result.ProviderItemID,
		ProviderVariationID: result.ProviderVariationID,
		Message:             result.Message,
		ResponseSummary:     result.Message,
	}, nil
}

func mapStockWriteResult(status connectorsdomain.StockWriteResultStatus) inventorydomain.StockWriteResultStatus {
	switch status {
	case connectorsdomain.StockWriteResultApplied:
		return inventorydomain.StockWriteResultApplied
	case connectorsdomain.StockWriteResultRejected:
		return inventorydomain.StockWriteResultRejected
	case connectorsdomain.StockWriteResultUnsupportedShape:
		return inventorydomain.StockWriteResultUnsupportedShape
	default:
		return inventorydomain.StockWriteResultTransientFailure
	}
}
