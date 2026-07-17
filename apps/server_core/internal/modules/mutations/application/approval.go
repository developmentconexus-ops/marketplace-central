package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

const FailureCodePreviewStale domain.FailureCode = "preview_stale"

type approvalStore interface {
	GetProtocol(context.Context, string) (ports.Protocol, bool, error)
	ApproveItems(context.Context, string, time.Time) error
	CancelProtocol(context.Context, string, time.Time) error
}

type ApprovalService struct {
	protocols approvalStore
	now       func() time.Time
}

func NewApprovalService(protocols approvalStore, now func() time.Time) ApprovalService {
	return ApprovalService{protocols: protocols, now: now}
}

func (s ApprovalService) Approve(ctx context.Context, protocolID string, payload json.RawMessage) (ports.Protocol, error) {
	var input struct {
		Execute *bool `json:"execute"`
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		return ports.Protocol{}, gateError(FailureCodeExecuteRequired, "aprovação exige execute=true explícito")
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return ports.Protocol{}, gateError(FailureCodeExecuteRequired, "aprovação exige execute=true explícito")
	}
	if err := RequireExecute(input.Execute); err != nil {
		return ports.Protocol{}, err
	}
	protocol, err := s.load(ctx, protocolID)
	if err != nil {
		return ports.Protocol{}, err
	}
	if err := domain.TransitionProtocolState(protocol.State, domain.ProtocolStateApproved); err != nil {
		return ports.Protocol{}, err
	}
	now := s.now().UTC()
	if protocol.PreviewedAt == nil || protocol.PreviewedAt.Before(now.Add(-15*time.Minute)) {
		return ports.Protocol{}, gateError(FailureCodePreviewStale, "prévia expirada; gere novamente")
	}
	if err := s.protocols.ApproveItems(ctx, protocolID, now); err != nil {
		return ports.Protocol{}, err
	}
	protocol.State = domain.ProtocolStateApproved
	protocol.ApprovedAt = &now
	return protocol, nil
}

func (s ApprovalService) Cancel(ctx context.Context, protocolID string) (ports.Protocol, error) {
	protocol, err := s.load(ctx, protocolID)
	if err != nil {
		return ports.Protocol{}, err
	}
	if err := domain.TransitionProtocolState(protocol.State, domain.ProtocolStateCancelled); err != nil {
		return ports.Protocol{}, err
	}
	now := s.now().UTC()
	if err := s.protocols.CancelProtocol(ctx, protocolID, now); err != nil {
		return ports.Protocol{}, err
	}
	protocol.State = domain.ProtocolStateCancelled
	protocol.FinishedAt = &now
	return protocol, nil
}

func (s ApprovalService) load(ctx context.Context, protocolID string) (ports.Protocol, error) {
	protocol, found, err := s.protocols.GetProtocol(ctx, protocolID)
	if err != nil {
		return ports.Protocol{}, err
	}
	if !found {
		return ports.Protocol{}, ErrProtocolNotFound
	}
	return protocol, nil
}
