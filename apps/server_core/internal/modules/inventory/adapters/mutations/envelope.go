package mutations

import (
	"context"
	"encoding/json"
	"strings"

	inventorydomain "marketplace-central/apps/server_core/internal/modules/inventory/domain"
	inventoryports "marketplace-central/apps/server_core/internal/modules/inventory/ports"
	mutationsdomain "marketplace-central/apps/server_core/internal/modules/mutations/domain"
	mutationsports "marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

var _ inventoryports.StockMutationEnvelope = (*Envelope)(nil)

type Envelope struct {
	repository mutationsports.ProtocolRepository
}

func NewEnvelope(repository mutationsports.ProtocolRepository) *Envelope {
	return &Envelope{repository: repository}
}

func (e *Envelope) CreateStockCorrection(ctx context.Context, action inventorydomain.StockAction) (string, error) {
	intent, err := json.Marshal(struct {
		PublishQuantity int `json:"publish_quantity"`
	}{PublishQuantity: action.RequestedQuantity})
	if err != nil {
		return "", err
	}
	listingID := stockListingID(action.ProviderRef)
	selection, err := json.Marshal(struct {
		Mode       string   `json:"mode"`
		ListingIDs []string `json:"listing_ids"`
	}{Mode: "explicit", ListingIDs: []string{listingID}})
	if err != nil {
		return "", err
	}
	protocol, err := e.repository.CreateProtocol(ctx, mutationsports.CreateProtocolInput{
		InstallationID: action.ProviderRef.InstallationID,
		Type:           mutationsdomain.ProtocolTypeStockCorrect,
		Actor:          action.Operator.ActorID,
		Intent:         intent,
		Selection:      selection,
		CreatedAt:      action.CreatedAt,
	})
	if err != nil {
		return "", err
	}
	before, err := json.Marshal(struct {
		AvailableQuantity *int `json:"available_quantity"`
	}{AvailableQuantity: action.BeforeQuantity})
	if err != nil {
		return "", err
	}
	if _, err := e.repository.ReplaceItems(ctx, protocol.ProtocolID, []mutationsports.ReplaceItemInput{{
		ListingID: listingID,
		Before:    before,
		After:     intent,
	}}); err != nil {
		return "", err
	}
	if err := e.repository.ApproveItems(ctx, protocol.ProtocolID, action.UpdatedAt); err != nil {
		return "", err
	}
	return protocol.ProtocolID, nil
}

func stockListingID(ref inventorydomain.ProviderStockRef) string {
	return strings.Join([]string{ref.InstallationID, ref.ProviderItemID, ref.ProviderVariationID}, "~")
}
