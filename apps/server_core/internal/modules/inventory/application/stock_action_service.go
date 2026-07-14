package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/ports"
	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

type StockActionService struct {
	store  ports.StockActionStore
	writer ports.StockWriter
	invalidator internalreadports.CacheInvalidator
}

func NewStockActionService(store ports.StockActionStore, writer ports.StockWriter) StockActionService {
	return NewStockActionServiceWithInvalidator(store, writer, nil)
}

func NewStockActionServiceWithInvalidator(store ports.StockActionStore, writer ports.StockWriter, invalidator internalreadports.CacheInvalidator) StockActionService {
	return StockActionService{store: store, writer: writer, invalidator: invalidator}
}

type ApplyManualStockActionInput struct {
	ActionID          string
	Risk              domain.RiskInput
	ProviderRef       domain.ProviderStockRef
	RequestedQuantity int
	Reason            string
	Approval          domain.ManualApproval
	Now               time.Time
}

func (s StockActionService) ApplyManual(ctx context.Context, input ApplyManualStockActionInput) (domain.StockAction, error) {
	if strings.TrimSpace(input.ActionID) == "" {
		return domain.StockAction{}, errors.New("stock action id is required")
	}
	if s.store == nil {
		return domain.StockAction{}, errors.New("stock action store is required")
	}
	if existing, ok, err := s.store.FindByID(ctx, input.ActionID); err != nil {
		return domain.StockAction{}, err
	} else if ok {
		return existing, nil
	}

	row := domain.ClassifyStockRisk(input.Risk)
	action := newStockAction(input, row)

	if !input.Approval.Approved {
		action.State = domain.StockActionStateSkipped
		action.BlockingReason = domain.BlockingReason{Code: "manual_approval_required", Message: "manual approval is required before provider write"}
		action = withActionEvent(action, domain.StockActionStateSkipped, action.BlockingReason.Message, input.Now)
		return action, s.store.Save(ctx, action)
	}
	if strings.TrimSpace(input.Approval.Operator.ActorID) == "" || input.Approval.ApprovedAt.IsZero() {
		action.State = domain.StockActionStateSkipped
		action.BlockingReason = domain.BlockingReason{Code: "manual_approval_incomplete", Message: "manual approval requires operator identity and approval timestamp"}
		action = withActionEvent(action, domain.StockActionStateSkipped, action.BlockingReason.Message, input.Now)
		return action, s.store.Save(ctx, action)
	}

	if row.RecommendedQuantity == nil {
		action.State = domain.StockActionStateBlocked
		action.BlockingReason = row.BlockingReason
		action = withActionEvent(action, domain.StockActionStateBlocked, action.BlockingReason.Message, input.Now)
		return action, s.store.Save(ctx, action)
	}
	if row.State == domain.StockRiskHealthy {
		action.State = domain.StockActionStateSkipped
		action.BlockingReason = domain.BlockingReason{Code: "no_stock_change_required", Message: "provider quantity already matches policy recommendation"}
		action = withActionEvent(action, domain.StockActionStateSkipped, action.BlockingReason.Message, input.Now)
		return action, s.store.Save(ctx, action)
	}
	if input.RequestedQuantity != *row.RecommendedQuantity {
		action.State = domain.StockActionStateBlocked
		action.BlockingReason = domain.BlockingReason{Code: "requested_quantity_mismatch", Message: "requested quantity does not match policy recommendation"}
		action = withActionEvent(action, domain.StockActionStateBlocked, action.BlockingReason.Message, input.Now)
		return action, s.store.Save(ctx, action)
	}
	if s.writer == nil {
		action.State = domain.StockActionStateFailed
		action.FailureReason = domain.BlockingReason{Code: "stock_writer_unavailable", Message: "stock writer is not configured"}
		action = withActionEvent(action, domain.StockActionStateFailed, action.FailureReason.Message, input.Now)
		return action, s.store.Save(ctx, action)
	}

	action.State = domain.StockActionStateApproved
	action = withActionEvent(action, domain.StockActionStateApproved, "manual approval accepted", input.Now)
	if err := s.store.Save(ctx, action); err != nil {
		return domain.StockAction{}, err
	}

	result, err := s.writer.UpdateAvailableQuantity(ctx, domain.StockWriteRequest{
		ProviderRef:       input.ProviderRef,
		IdempotencyKey:    input.ActionID,
		RequestedQuantity: input.RequestedQuantity,
		Reason:            input.Reason,
	})
	action.ProviderResult = result
	if err != nil {
		action.State = domain.StockActionStateFailed
		action.FailureReason = domain.BlockingReason{Code: "provider_error", Message: err.Error()}
		action = withActionEvent(action, domain.StockActionStateFailed, action.FailureReason.Message, input.Now)
		return action, s.store.Save(ctx, action)
	}
	switch result.Status {
	case domain.StockWriteResultApplied:
		action.State = domain.StockActionStateApplied
		action = withActionEvent(action, domain.StockActionStateApplied, result.Message, input.Now)
	default:
		action.State = domain.StockActionStateFailed
		action.FailureReason = domain.BlockingReason{Code: "provider_rejected", Message: result.Message}
		action = withActionEvent(action, domain.StockActionStateFailed, result.Message, input.Now)
	}
	if err := s.store.Save(ctx, action); err != nil {
		return domain.StockAction{}, err
	}
	if action.State == domain.StockActionStateApplied && s.invalidator != nil {
		s.invalidator.InvalidateClass("inventory")
	}
	return action, nil
}

func newStockAction(input ApplyManualStockActionInput, row domain.StockRiskRow) domain.StockAction {
	return domain.StockAction{
		ActionID:            input.ActionID,
		State:               domain.StockActionStateProposed,
		Trigger:             domain.StockActionTriggerManual,
		ProviderRef:         input.ProviderRef,
		BeforeQuantity:      row.ProviderQuantity,
		RequestedQuantity:   input.RequestedQuantity,
		RecommendedQuantity: row.RecommendedQuantity,
		PolicyID:            row.PolicyID,
		InternalObservedAt:  row.InternalObservedAt,
		ProviderObservedAt:  row.ProviderObservedAt,
		Operator:            input.Approval.Operator,
		Reason:              input.Reason,
		IdempotencyKey:      input.ActionID,
		AuditEvents: []domain.StockActionAuditEvent{
			{State: domain.StockActionStateProposed, Reason: input.Reason, OccurredAt: input.Now},
		},
		CreatedAt: input.Now,
		UpdatedAt: input.Now,
	}
}

func withActionEvent(action domain.StockAction, state domain.StockActionState, reason string, occurredAt time.Time) domain.StockAction {
	action.UpdatedAt = occurredAt
	action.AuditEvents = append(action.AuditEvents, domain.StockActionAuditEvent{
		State:          state,
		Reason:         reason,
		ProviderStatus: action.ProviderResult.Status,
		OccurredAt:     occurredAt,
	})
	return action
}
