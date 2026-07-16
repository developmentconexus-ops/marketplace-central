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

// ProtocolRepository serves mutation application commands and protocol polling reads.
type ProtocolRepository interface {
	CreateProtocol(context.Context, CreateProtocolInput) (Protocol, error)
	GetProtocol(context.Context, string) (Protocol, bool, error)
}
