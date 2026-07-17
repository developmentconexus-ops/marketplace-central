package ports

import (
	"context"
	"encoding/json"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

type Protocol struct {
	ProtocolID     string
	InstallationID string
	Type           domain.ProtocolType
	State          domain.ProtocolState
	Actor          string
	Intent         json.RawMessage
	Selection      json.RawMessage
	Totals         json.RawMessage
	SourceAsOf     *time.Time
	RetriedFrom    *string
	CreatedAt      time.Time
	PreviewedAt    *time.Time
	ApprovedAt     *time.Time
	FinishedAt     *time.Time
}

type CreateProtocolInput struct {
	InstallationID string
	Type           domain.ProtocolType
	Actor          string
	Intent         json.RawMessage
	Selection      json.RawMessage
	CreatedAt      time.Time
}

type MutationItem struct {
	Seq            int
	ItemID         string
	ListingID      string
	IdempotencyKey string
	Before         json.RawMessage
	After          json.RawMessage
	State          domain.ItemState
	Failure        json.RawMessage
	AppliedAt      *time.Time
}

type ReplaceItemInput struct {
	ListingID string
	Before    json.RawMessage
	After     json.RawMessage
}

type ItemOutcome struct {
	State     domain.ItemState
	Failure   json.RawMessage
	AppliedAt time.Time
}

type ProtocolClaim interface {
	Protocol() Protocol
	FetchPendingItems(context.Context) ([]MutationItem, error)
	MarkItemApplying(context.Context, string) error
	WriteItemOutcome(context.Context, string, ItemOutcome) error
	AppliedIdempotencyKeys(context.Context, []string) ([]string, error)
	ItemStateCounts(context.Context) (map[domain.ItemState]int, error)
	Finish(context.Context, domain.ProtocolState, time.Time) error
	Commit(context.Context) error
	Rollback(context.Context) error
}

// ProtocolRepository serves mutation application commands and protocol polling reads.
type ProtocolRepository interface {
	CreateProtocol(context.Context, CreateProtocolInput) (Protocol, error)
	GetProtocol(context.Context, string) (Protocol, bool, error)
	ReplaceItems(context.Context, string, []ReplaceItemInput, *time.Time) ([]MutationItem, error)
	ApproveItems(context.Context, string, time.Time) error
	ClaimProtocol(context.Context, string) (ProtocolClaim, bool, error)
}
