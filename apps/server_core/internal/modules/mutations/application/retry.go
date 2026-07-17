package application

import (
	"context"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

const FailureCodeNothingToRetry domain.FailureCode = "nothing_to_retry"

type retryStore interface {
	GetProtocol(context.Context, string) (ports.Protocol, bool, error)
	CloneRetry(context.Context, string, time.Time) (ports.Protocol, bool, error)
}

type RetryService struct {
	protocols retryStore
	now       func() time.Time
}

func NewRetryService(protocols retryStore, now func() time.Time) RetryService {
	return RetryService{protocols: protocols, now: now}
}

func (s RetryService) Retry(ctx context.Context, protocolID string) (ports.Protocol, error) {
	protocol, found, err := s.protocols.GetProtocol(ctx, protocolID)
	if err != nil {
		return ports.Protocol{}, err
	}
	if !found {
		return ports.Protocol{}, ErrProtocolNotFound
	}
	if protocol.State != domain.ProtocolStatePartiallyFailed && protocol.State != domain.ProtocolStateFailedPreserved {
		return ports.Protocol{}, &domain.InvalidStateTransitionError{Code: domain.InvalidStateCode, From: protocol.State, To: domain.ProtocolStateDraft}
	}
	clone, eligible, err := s.protocols.CloneRetry(ctx, protocolID, s.now().UTC())
	if err != nil {
		return ports.Protocol{}, err
	}
	if !eligible {
		return ports.Protocol{}, gateError(FailureCodeNothingToRetry, "não há falhas elegíveis para nova tentativa")
	}
	return clone, nil
}
