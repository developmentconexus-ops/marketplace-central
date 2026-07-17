package connectors

import (
	"context"
	"encoding/json"
	"errors"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type StockWriter struct {
	capabilities                              *connectorsapp.MarketplaceCapabilityService
	tenantID, providerCode, providerAccountID string
}

func NewStockWriter(capabilities *connectorsapp.MarketplaceCapabilityService, tenantID, providerCode, providerAccountID string) *StockWriter {
	return &StockWriter{capabilities: capabilities, tenantID: tenantID, providerCode: providerCode, providerAccountID: providerAccountID}
}

func (w *StockWriter) Apply(ctx context.Context, item ports.WriteItem) (ports.WriteOutcome, error) {
	writer, err := w.capabilities.StockWriter(w.providerCode)
	if err != nil {
		return failedWrite(err), nil
	}
	listing, err := listingsdomain.ParseListingID(item.ListingID)
	if err != nil {
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderInvalidReference, "listing reference is invalid")), nil
	}
	var after struct {
		PublishQuantity *int   `json:"publish_quantity"`
		Reason          string `json:"reason"`
	}
	if err := json.Unmarshal(item.After, &after); err != nil || after.PublishQuantity == nil {
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderValidation, "stock intent is invalid")), nil
	}
	result, err := writer.UpdateAvailableQuantity(ctx, connectorsdomain.StockWriteRequest{TenantID: w.tenantID, InstallationID: item.InstallationID, ProviderCode: w.providerCode, ProviderAccountID: w.providerAccountID, ProviderItemID: listing.ProviderListingID, ProviderVariationID: listing.VariationID, IdempotencyKey: item.IdempotencyKey, RequestedQuantity: *after.PublishQuantity, Reason: after.Reason})
	if err != nil {
		return failedWrite(err), nil
	}
	switch result.Result {
	case connectorsdomain.StockWriteResultApplied:
		return ports.WriteOutcome{}, nil
	case connectorsdomain.StockWriteResultRejected:
		return rejectedWrite(result.Message, false), nil
	case connectorsdomain.StockWriteResultTransientFailure:
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderTransient, result.Message)), nil
	case connectorsdomain.StockWriteResultUnsupportedShape:
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderUnsupportedShape, result.Message)), nil
	default:
		return failedWrite(errors.New("unknown provider stock outcome")), nil
	}
}
