package connectors

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type PriceWriter struct {
	capabilities           *connectorsapp.MarketplaceCapabilityService
	tenantID, providerCode string
}

func NewPriceWriter(capabilities *connectorsapp.MarketplaceCapabilityService, tenantID, providerCode string) *PriceWriter {
	return &PriceWriter{capabilities: capabilities, tenantID: tenantID, providerCode: providerCode}
}

func (w *PriceWriter) Apply(ctx context.Context, item ports.WriteItem) (ports.WriteOutcome, error) {
	writer, err := w.capabilities.PriceWriter(w.providerCode)
	if err != nil {
		return failedWrite(err), nil
	}
	listing, err := listingsdomain.ParseListingID(item.ListingID)
	if err != nil {
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderInvalidReference, "listing reference is invalid")), nil
	}
	var after struct {
		NewPrice connectorsdomain.Price `json:"new_price"`
	}
	if err := json.Unmarshal(item.After, &after); err != nil {
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderValidation, "price intent is invalid")), nil
	}
	result, err := writer.UpdatePrice(ctx, connectorsdomain.PriceWriteRequest{TenantID: strings.TrimSpace(w.tenantID), InstallationID: item.InstallationID, ListingID: listing.ProviderListingID, IdempotencyKey: item.IdempotencyKey, Price: after.NewPrice})
	if err != nil {
		return failedWrite(err), nil
	}
	switch result.Result {
	case connectorsdomain.WriteResultApplied:
		return ports.WriteOutcome{}, nil
	case connectorsdomain.WriteResultRejected:
		return rejectedWrite(result.Message, false), nil
	case connectorsdomain.WriteResultTransientFailure:
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderTransient, result.Message)), nil
	case connectorsdomain.WriteResultUnsupportedShape:
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderUnsupportedShape, result.Message)), nil
	default:
		return failedWrite(errors.New("unknown provider price outcome")), nil
	}
}

func failedWrite(err error) ports.WriteOutcome {
	failure := MapFailure(err)
	return ports.WriteOutcome{Failure: &failure}
}
func rejectedWrite(message string, paused bool) ports.WriteOutcome {
	failure := MapRejected(message, paused)
	return ports.WriteOutcome{Failure: &failure}
}
