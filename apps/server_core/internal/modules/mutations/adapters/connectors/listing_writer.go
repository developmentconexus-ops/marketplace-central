package connectors

import (
	"context"
	"encoding/json"
	"errors"

	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type ListingWriter struct {
	capabilities           *connectorsapp.MarketplaceCapabilityService
	tenantID, providerCode string
}

func NewListingWriter(capabilities *connectorsapp.MarketplaceCapabilityService, tenantID, providerCode string) *ListingWriter {
	return &ListingWriter{capabilities: capabilities, tenantID: tenantID, providerCode: providerCode}
}

func (w *ListingWriter) Apply(ctx context.Context, item ports.WriteItem) (ports.WriteOutcome, error) {
	writer, err := w.capabilities.ListingWriter(w.providerCode)
	if err != nil {
		return failedWrite(err), nil
	}
	listing, err := listingsdomain.ParseListingID(item.ListingID)
	if err != nil {
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderInvalidReference, "listing reference is invalid")), nil
	}
	request := connectorsdomain.ListingWriteRequest{TenantID: w.tenantID, InstallationID: item.InstallationID, ListingID: listing.ProviderListingID, IdempotencyKey: item.IdempotencyKey}
	switch item.ProtocolType {
	case mutationsdomain.ProtocolTypeListingPause:
		request.Action = connectorsdomain.ListingWritePause
	case mutationsdomain.ProtocolTypeListingEdit:
		request.Action = connectorsdomain.ListingWriteEdit
		var after struct {
			Attributes []connectorsdomain.ListingAttribute `json:"attributes"`
		}
		if err := json.Unmarshal(item.After, &after); err != nil {
			return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderValidation, "listing intent is invalid")), nil
		}
		request.Attributes = after.Attributes
	default:
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderValidation, "listing intent type is invalid")), nil
	}
	result, err := writer.UpdateListing(ctx, request)
	if err != nil {
		return failedWrite(err), nil
	}
	switch result.Result {
	case connectorsdomain.WriteResultApplied:
		return ports.WriteOutcome{}, nil
	case connectorsdomain.WriteResultRejected:
		return rejectedWrite(result.Message, item.ProtocolType == mutationsdomain.ProtocolTypeListingPause), nil
	case connectorsdomain.WriteResultTransientFailure:
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderTransient, result.Message)), nil
	case connectorsdomain.WriteResultUnsupportedShape:
		return failedWrite(connectorsdomain.NewCapabilityError(connectorsdomain.ErrCodeProviderUnsupportedShape, result.Message)), nil
	default:
		return failedWrite(errors.New("unknown provider listing outcome")), nil
	}
}
