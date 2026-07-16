package ports

import (
	"context"
	"encoding/json"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

type WriteItem struct {
	ProtocolID     string
	InstallationID string
	ProtocolType   domain.ProtocolType
	ItemID         string
	ListingID      string
	IdempotencyKey string
	Before         json.RawMessage
	After          json.RawMessage
}

type WriteOutcome struct {
	Failure *domain.Failure
}

type WriterPort interface {
	Apply(context.Context, WriteItem) (WriteOutcome, error)
}
