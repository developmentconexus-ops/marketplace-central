package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type Poller struct {
	repository ports.ProtocolRepository
	writer     ports.WriterPort
	now        func() time.Time
}

const protocolSourceMaxAge = 15 * time.Minute

func NewPoller(repository ports.ProtocolRepository, writer ports.WriterPort, now func() time.Time) *Poller {
	return &Poller{repository: repository, writer: writer, now: now}
}

func (p *Poller) Pass(ctx context.Context, installationID string) (worked bool, err error) {
	claim, found, err := p.repository.ClaimProtocol(ctx, installationID)
	if err != nil || !found {
		return false, err
	}
	defer func() { _ = claim.Rollback(ctx) }()

	protocol := claim.Protocol()
	execute := true
	if err := ValidateProtocol(ProtocolWriteValidation{Actor: protocol.Actor, Execute: &execute, SourceAsOf: protocol.SourceAsOf, Now: p.now(), MaxAge: protocolSourceMaxAge}); err != nil {
		failure, marshalErr := marshalFailure(*itemFailure(err))
		if marshalErr != nil {
			return true, marshalErr
		}
		for {
			items, fetchErr := claim.FetchPendingItems(ctx)
			if fetchErr != nil {
				return true, fmt.Errorf("fetch pending mutation items after protocol gate failure: %w", fetchErr)
			}
			if len(items) == 0 {
				break
			}
			for _, item := range items {
				if writeErr := claim.WriteItemOutcome(ctx, item.ItemID, ports.ItemOutcome{State: domain.ItemStateFailed, Failure: failure}); writeErr != nil {
					return true, fmt.Errorf("write mutation item gate failure outcome: %w", writeErr)
				}
			}
		}
		if finishErr := p.finish(ctx, claim); finishErr != nil {
			return true, finishErr
		}
		return true, nil
	}
	for {
		items, fetchErr := claim.FetchPendingItems(ctx)
		if fetchErr != nil {
			return true, fmt.Errorf("fetch pending mutation items: %w", fetchErr)
		}
		if len(items) == 0 {
			break
		}
		keys := make([]string, len(items))
		for i := range items {
			keys[i] = items[i].IdempotencyKey
		}
		applied, keysErr := claim.AppliedIdempotencyKeys(ctx, keys)
		if keysErr != nil {
			return true, fmt.Errorf("read applied idempotency keys: %w", keysErr)
		}
		for _, item := range items {
			outcome := ports.ItemOutcome{State: domain.ItemStateSkipped}
			write := ports.WriteOutcome{}
			skipped, gateErr := WriteThroughGates(ctx, item.IdempotencyKey,
				func(context.Context, string) (bool, error) {
					for _, key := range applied {
						if key == item.IdempotencyKey {
							return true, nil
						}
					}
					return false, nil
				},
				func(ctx context.Context) error { return claim.MarkItemApplying(ctx, item.ItemID) },
				func(ctx context.Context) error {
					var writeErr error
					write, writeErr = p.writer.Apply(ctx, ports.WriteItem{ProtocolID: protocol.ProtocolID, InstallationID: protocol.InstallationID, ProtocolType: protocol.Type, ItemID: item.ItemID, ListingID: item.ListingID, IdempotencyKey: item.IdempotencyKey, Before: item.Before, After: item.After})
					if writeErr != nil {
						write.Failure = &domain.Failure{Code: domain.FailureCodeInternal, MessagePT: "Falha interna ao aplicar alteração."}
					}
					return nil
				})
			if gateErr != nil {
				return true, gateErr
			}
			if !skipped {
				if write.Failure == nil {
					outcome.State, outcome.AppliedAt = domain.ItemStateApplied, p.now()
				} else {
					outcome.State = domain.ItemStateFailed
					outcome.Failure, err = marshalFailure(*write.Failure)
					if err != nil {
						return true, err
					}
				}
			}
			if err = claim.WriteItemOutcome(ctx, item.ItemID, outcome); err != nil {
				return true, fmt.Errorf("write mutation item outcome: %w", err)
			}
		}
	}
	if err = p.finish(ctx, claim); err != nil {
		return true, err
	}
	return true, nil
}

func (p *Poller) finish(ctx context.Context, claim ports.ProtocolClaim) error {
	counts, err := claim.ItemStateCounts(ctx)
	if err != nil {
		return fmt.Errorf("read mutation item state counts: %w", err)
	}
	states := make([]domain.ItemState, 0)
	for _, state := range []domain.ItemState{domain.ItemStateApplied, domain.ItemStateFailed, domain.ItemStateSkipped} {
		for range counts[state] {
			states = append(states, state)
		}
	}
	terminal, err := domain.TerminalProtocolState(states)
	if err != nil {
		return fmt.Errorf("compute mutation protocol terminal state: %w", err)
	}
	if err = claim.Finish(ctx, terminal, p.now()); err != nil {
		return fmt.Errorf("finish mutation protocol: %w", err)
	}
	if err = claim.Commit(ctx); err != nil {
		return fmt.Errorf("commit mutation protocol: %w", err)
	}
	return nil
}

func marshalFailure(f domain.Failure) (json.RawMessage, error) {
	f.Retryable = f.Code.RetryableByDefault()
	value := struct {
		Code            domain.FailureCode `json:"code"`
		MessagePT       string             `json:"message_pt"`
		MessageProvider string             `json:"message_provider"`
		Retryable       bool               `json:"retryable"`
	}{f.Code, f.MessagePT, f.MessageProvider, f.Retryable}
	b, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal mutation failure: %w", err)
	}
	return b, nil
}
