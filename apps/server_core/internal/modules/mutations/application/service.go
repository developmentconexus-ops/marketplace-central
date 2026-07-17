package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	listingsadapter "marketplace-central/apps/server_core/internal/modules/mutations/adapters/listings"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

const (
	FailureCodeEmptySelection    domain.FailureCode = "empty_selection"
	FailureCodeSelectionTooLarge domain.FailureCode = "selection_too_large"
)

var ErrProtocolNotFound = errors.New("mutation protocol not found")

type protocolStore interface {
	CreateProtocol(context.Context, ports.CreateProtocolInput) (ports.Protocol, error)
	GetProtocol(context.Context, string) (ports.Protocol, bool, error)
	ReplacePreview(context.Context, string, []ports.ReplaceItemInput, time.Time, time.Time) ([]ports.MutationItem, error)
}

type selectionResolver interface {
	Resolve(context.Context, string, listingsadapter.Selection) ([]string, error)
}

type Service struct {
	protocols  protocolStore
	selections selectionResolver
	listings   listingReader
	now        func() time.Time
}

func NewService(protocols protocolStore, selections selectionResolver, listings listingReader, now func() time.Time) Service {
	return Service{protocols: protocols, selections: selections, listings: listings, now: now}
}

type CreateInput struct {
	InstallationID string
	Type           domain.ProtocolType
	Actor          string
	Intent         json.RawMessage
	Selection      json.RawMessage
}

func (s Service) Create(ctx context.Context, input CreateInput) (ports.Protocol, error) {
	if _, err := listingsadapter.ParseSelectionJSON(input.Selection); err != nil {
		return ports.Protocol{}, err
	}
	var created ports.Protocol
	err := InsertWithActorGate(ctx, input.Actor, func(ctx context.Context) error {
		var insertErr error
		created, insertErr = s.protocols.CreateProtocol(ctx, ports.CreateProtocolInput{
			InstallationID: input.InstallationID,
			Type:           input.Type,
			Actor:          input.Actor,
			Intent:         append(json.RawMessage(nil), input.Intent...),
			Selection:      append(json.RawMessage(nil), input.Selection...),
			CreatedAt:      s.now().UTC(),
		})
		return insertErr
	})
	return created, err
}
