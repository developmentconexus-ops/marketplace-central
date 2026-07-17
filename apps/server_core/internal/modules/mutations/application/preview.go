package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsadapter "marketplace-central/apps/server_core/internal/modules/mutations/adapters/listings"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

const previewItemLimit = 2000

type listingReader interface {
	Get(context.Context, listingsdomain.ListingID) (listingsdomain.ListingReadModel, []listingsdomain.TimelineEvent, error)
}

type PreviewResult struct {
	Protocol ports.Protocol
	Items    []ports.MutationItem
}

func (s Service) Preview(ctx context.Context, protocolID string) (PreviewResult, error) {
	protocol, found, err := s.protocols.GetProtocol(ctx, protocolID)
	if err != nil {
		return PreviewResult{}, err
	}
	if !found {
		return PreviewResult{}, ErrProtocolNotFound
	}
	if err := domain.TransitionProtocolState(protocol.State, domain.ProtocolStatePreviewed); err != nil {
		return PreviewResult{}, err
	}
	selection, err := listingsadapter.ParseSelectionJSON(protocol.Selection)
	if err != nil {
		return PreviewResult{}, err
	}
	ids, err := s.selections.Resolve(ctx, protocol.InstallationID, selection)
	if err != nil {
		return PreviewResult{}, err
	}
	if len(ids) == 0 {
		return PreviewResult{}, gateError(FailureCodeEmptySelection, "a seleção não contém anúncios")
	}
	if len(ids) > previewItemLimit {
		return PreviewResult{}, gateError(FailureCodeSelectionTooLarge, "a seleção excede 2000 anúncios")
	}

	inputs := make([]ports.ReplaceItemInput, 0, len(ids))
	var sourceAsOf time.Time
	for _, rawID := range ids {
		id, err := listingsdomain.ParseListingID(rawID)
		if err != nil {
			return PreviewResult{}, err
		}
		if id.InstallationID != protocol.InstallationID {
			return PreviewResult{}, fmt.Errorf("listing %q is outside protocol installation", rawID)
		}
		listing, _, err := s.listings.Get(ctx, id)
		if err != nil {
			return PreviewResult{}, err
		}
		if listing.FetchedAt == nil || listing.FetchedAt.IsZero() {
			return PreviewResult{}, gateError(FailureCodeSourceTimeUnavailable, "horário da fonte indisponível para o anúncio")
		}
		if listing.FetchedAt.After(sourceAsOf) {
			sourceAsOf = listing.FetchedAt.UTC()
		}
		before, err := previewBefore(protocol.Type, listing)
		if err != nil {
			return PreviewResult{}, err
		}
		inputs = append(inputs, ports.ReplaceItemInput{ListingID: rawID, Before: before, After: append(json.RawMessage(nil), protocol.Intent...)})
	}
	previewedAt := s.now().UTC()
	items, err := s.protocols.ReplacePreview(ctx, protocolID, inputs, sourceAsOf, previewedAt)
	if err != nil {
		return PreviewResult{}, err
	}
	protocol.State = domain.ProtocolStatePreviewed
	protocol.SourceAsOf = &sourceAsOf
	protocol.PreviewedAt = &previewedAt
	protocol.Totals = totalsJSON(len(items))
	return PreviewResult{Protocol: protocol, Items: items}, nil
}

func previewBefore(protocolType domain.ProtocolType, listing listingsdomain.ListingReadModel) (json.RawMessage, error) {
	var value any
	switch protocolType {
	case domain.ProtocolTypePriceUpdate:
		value = struct {
			Price *listingsdomain.Money `json:"price"`
		}{listing.Price}
	case domain.ProtocolTypeStockCorrect:
		value = struct {
			PublishQuantity *int `json:"publish_quantity"`
		}{listing.PublishedQuantity}
	case domain.ProtocolTypeLinkApply:
		value = struct {
			Link listingsdomain.ListingLink `json:"link"`
		}{listing.Link}
	case domain.ProtocolTypeListingPause:
		value = struct {
			Status listingsdomain.ListingStatus `json:"status"`
		}{listing.Status}
	case domain.ProtocolTypeListingResync, domain.ProtocolTypeListingEdit:
		value = listing
	default:
		return nil, fmt.Errorf("preview unsupported for mutation type %q", protocolType)
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal preview before snapshot: %w", err)
	}
	return payload, nil
}

func totalsJSON(items int) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{"items":%d,"previewed":%d,"applied":0,"failed":0,"skipped":0}`, items, items))
}
