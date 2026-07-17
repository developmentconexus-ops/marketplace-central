package application

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type ProtocolCursor struct {
	CreatedAt  time.Time
	ProtocolID string
}

type ItemCursor struct{ Seq int }

type ProtocolQuery struct {
	InstallationID string
	State          domain.ProtocolState
	Type           domain.ProtocolType
	Cursor         ProtocolCursor
	Limit          int
}

type ItemQuery struct {
	ProtocolID string
	State      domain.ItemState
	Cursor     ItemCursor
	Limit      int
}

type ProtocolPage struct {
	Items      []ports.Protocol
	NextCursor *ProtocolCursor
}

type ItemPage struct {
	Items      []ports.MutationItem
	NextCursor *ItemCursor
}

type readStore interface {
	ListProtocols(context.Context, ProtocolQuery) (ProtocolPage, error)
	GetProtocol(context.Context, string) (ports.Protocol, bool, error)
	ListItems(context.Context, ItemQuery) (ItemPage, error)
}

type ReadService struct{ protocols readStore }

func NewReadService(protocols readStore) ReadService { return ReadService{protocols: protocols} }

func (s ReadService) ListProtocols(ctx context.Context, query ProtocolQuery) (ProtocolPage, error) {
	return s.protocols.ListProtocols(ctx, query)
}

func (s ReadService) GetProtocol(ctx context.Context, protocolID string) (ports.Protocol, error) {
	protocol, found, err := s.protocols.GetProtocol(ctx, protocolID)
	if err != nil {
		return ports.Protocol{}, err
	}
	if !found {
		return ports.Protocol{}, ErrProtocolNotFound
	}
	return protocol, nil
}

func (s ReadService) ListItems(ctx context.Context, query ItemQuery) (ItemPage, error) {
	if _, err := s.GetProtocol(ctx, query.ProtocolID); err != nil {
		return ItemPage{}, err
	}
	return s.protocols.ListItems(ctx, query)
}
