package application

import (
	"context"
	"errors"
	"testing"
	"time"

	internalreadports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
)

func TestApplyManualStockActionBlocksUnsafeRowsWithoutProviderWrite(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)

	cases := []struct {
		name     string
		risk     domain.RiskInput
		mutate   func(*domain.RiskInput)
		wantCode string
	}{
		{name: "unresolved", risk: actionRiskInput(now, domain.LinkStateUnresolved, intPtr(8), intPtr(9), &fresh, &fresh), wantCode: "unresolved_link"},
		{name: "rejected", risk: actionRiskInput(now, domain.LinkStateRejected, intPtr(8), intPtr(9), &fresh, &fresh), wantCode: "rejected_link"},
		{name: "conflict", risk: actionRiskInput(now, domain.LinkStateConflict, intPtr(8), intPtr(9), &fresh, &fresh), wantCode: "conflict_link"},
		{name: "stale internal", risk: actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), nil, &fresh), wantCode: "stale_internal_source"},
		{name: "stale provider", risk: actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, nil), wantCode: "stale_provider_source"},
		{name: "unsupported provider", risk: actionRiskInput(now, domain.LinkStateResolved, intPtr(8), nil, &fresh, &fresh), wantCode: "missing_provider_quantity"},
		{
			name: "ineligible",
			risk: actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh),
			mutate: func(input *domain.RiskInput) {
				input.Policy.Eligibility = domain.EligibilityPolicy{
					Rules: []domain.EligibilityRule{domain.GroupExclusion(77, "showroom-only group")},
				}
				input.Product.GroupID = intPtr(77)
			},
			wantCode: "ineligible_product",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			writer := &fakeStockWriter{}
			service := NewStockActionService(&memoryActionStore{}, writer)
			risk := tt.risk
			if tt.mutate != nil {
				tt.mutate(&risk)
			}

			action, err := service.ApplyManual(context.Background(), applyInput("act-"+tt.name, risk, 7, now))
			if err != nil {
				t.Fatalf("ApplyManual() error = %v", err)
			}

			if action.State != domain.StockActionStateBlocked {
				t.Fatalf("state=%s, want blocked", action.State)
			}
			if action.BlockingReason.Code != tt.wantCode {
				t.Fatalf("blocking code=%q, want %q", action.BlockingReason.Code, tt.wantCode)
			}
			if writer.calls != 0 {
				t.Fatalf("writer calls=%d, want 0", writer.calls)
			}
		})
	}
}

func TestApplyManualStockActionSkipsWithoutManualApproval(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	writer := &fakeStockWriter{}
	service := NewStockActionService(&memoryActionStore{}, writer)
	input := applyInput("act-unapproved", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh), 7, now)
	input.Approval.Approved = false

	action, err := service.ApplyManual(context.Background(), input)
	if err != nil {
		t.Fatalf("ApplyManual() error = %v", err)
	}

	if action.State != domain.StockActionStateSkipped {
		t.Fatalf("state=%s, want skipped", action.State)
	}
	if action.BlockingReason.Code != "manual_approval_required" {
		t.Fatalf("blocking code=%q, want manual_approval_required", action.BlockingReason.Code)
	}
	if writer.calls != 0 {
		t.Fatalf("writer calls=%d, want 0", writer.calls)
	}
}

func TestApplyManualStockActionSkipsWhenManualApprovalIsIncomplete(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	writer := &fakeStockWriter{}
	service := NewStockActionService(&memoryActionStore{}, writer)
	input := applyInput("act-incomplete-approval", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh), 7, now)
	input.Approval.Operator.ActorID = ""

	action, err := service.ApplyManual(context.Background(), input)
	if err != nil {
		t.Fatalf("ApplyManual() error = %v", err)
	}

	if action.State != domain.StockActionStateSkipped {
		t.Fatalf("state=%s, want skipped", action.State)
	}
	if action.BlockingReason.Code != "manual_approval_incomplete" {
		t.Fatalf("blocking code=%q, want manual_approval_incomplete", action.BlockingReason.Code)
	}
	if writer.calls != 0 {
		t.Fatalf("writer calls=%d, want 0", writer.calls)
	}
}

func TestApplyManualStockActionSkipsHealthyNoopRows(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	writer := &fakeStockWriter{}
	service := NewStockActionService(&memoryActionStore{}, writer)

	action, err := service.ApplyManual(context.Background(), applyInput("act-healthy", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(7), &fresh, &fresh), 7, now))
	if err != nil {
		t.Fatalf("ApplyManual() error = %v", err)
	}

	if action.State != domain.StockActionStateSkipped {
		t.Fatalf("state=%s, want skipped", action.State)
	}
	if action.BlockingReason.Code != "no_stock_change_required" {
		t.Fatalf("blocking code=%q, want no_stock_change_required", action.BlockingReason.Code)
	}
	if writer.calls != 0 {
		t.Fatalf("writer calls=%d, want 0", writer.calls)
	}
}

func TestApplyManualStockActionWritesWithIdempotencyAndAuditEvidence(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	writer := &fakeStockWriter{
		result: domain.StockWriteResult{
			Status:              domain.StockWriteResultApplied,
			IdempotencyKey:      "act-apply",
			ProviderCode:        "mercado_livre",
			ProviderItemID:      "MLB123",
			ProviderVariationID: "987",
			Message:             "updated",
			ResponseSummary:     "200 OK",
		},
	}
	store := &memoryActionStore{}
	service := NewStockActionService(store, writer)

	action, err := service.ApplyManual(context.Background(), applyInput("act-apply", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh), 7, now))
	if err != nil {
		t.Fatalf("ApplyManual() error = %v", err)
	}

	if action.State != domain.StockActionStateApplied {
		t.Fatalf("state=%s, want applied", action.State)
	}
	if writer.calls != 1 {
		t.Fatalf("writer calls=%d, want 1", writer.calls)
	}
	if writer.lastRequest.IdempotencyKey != "act-apply" {
		t.Fatalf("idempotency key=%q, want action id", writer.lastRequest.IdempotencyKey)
	}
	if action.BeforeQuantity == nil || *action.BeforeQuantity != 9 {
		t.Fatalf("before quantity=%v, want 9", action.BeforeQuantity)
	}
	if action.RequestedQuantity != 7 {
		t.Fatalf("requested quantity=%d, want 7", action.RequestedQuantity)
	}
	if action.PolicyID != "stock-seguro-default" {
		t.Fatalf("policy id=%q, want stock-seguro-default", action.PolicyID)
	}
	if action.InternalObservedAt == nil || action.ProviderObservedAt == nil {
		t.Fatalf("expected source timestamps in audit: %+v", action)
	}
	if action.Operator.ActorID != "leandro" || action.Trigger != domain.StockActionTriggerManual {
		t.Fatalf("operator/trigger not audited: %+v", action)
	}
	if action.ProviderResult.Status != domain.StockWriteResultApplied || action.ProviderResult.ResponseSummary != "200 OK" {
		t.Fatalf("provider result not audited: %+v", action.ProviderResult)
	}
	if !hasActionEvent(action, domain.StockActionStateApproved) || !hasActionEvent(action, domain.StockActionStateApplied) {
		t.Fatalf("expected approved and applied audit events, got %+v", action.AuditEvents)
	}
	for _, event := range action.AuditEvents {
		if event.OccurredAt.IsZero() {
			t.Fatalf("audit event has zero timestamp: %+v", event)
		}
	}
}

func TestApplyManualStockActionPersistsProviderFailureAudit(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	writer := &fakeStockWriter{
		result: domain.StockWriteResult{
			Status:          domain.StockWriteResultRejected,
			IdempotencyKey:  "act-fail",
			ProviderCode:    "mercado_livre",
			ProviderItemID:  "MLB123",
			Message:         "provider rejected quantity",
			ResponseSummary: "400 validation",
		},
	}
	service := NewStockActionService(&memoryActionStore{}, writer)

	action, err := service.ApplyManual(context.Background(), applyInput("act-fail", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh), 7, now))
	if err != nil {
		t.Fatalf("ApplyManual() error = %v", err)
	}

	if action.State != domain.StockActionStateFailed {
		t.Fatalf("state=%s, want failed", action.State)
	}
	if action.ProviderResult.Status != domain.StockWriteResultRejected {
		t.Fatalf("provider status=%s, want rejected", action.ProviderResult.Status)
	}
	if action.FailureReason.Code != "provider_rejected" {
		t.Fatalf("failure code=%q, want provider_rejected", action.FailureReason.Code)
	}
	if !hasActionEvent(action, domain.StockActionStateFailed) {
		t.Fatalf("expected failed audit event, got %+v", action.AuditEvents)
	}
}

func TestApplyManualStockActionPersistsWriterErrorAudit(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	writer := &fakeStockWriter{err: errors.New("network timeout")}
	service := NewStockActionService(&memoryActionStore{}, writer)

	action, err := service.ApplyManual(context.Background(), applyInput("act-error", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh), 7, now))
	if err != nil {
		t.Fatalf("ApplyManual() error = %v", err)
	}

	if action.State != domain.StockActionStateFailed {
		t.Fatalf("state=%s, want failed", action.State)
	}
	if action.FailureReason.Code != "provider_error" || action.FailureReason.Message != "network timeout" {
		t.Fatalf("failure reason=%+v, want provider_error network timeout", action.FailureReason)
	}
	if !hasActionEvent(action, domain.StockActionStateApproved) || !hasActionEvent(action, domain.StockActionStateFailed) {
		t.Fatalf("expected approved and failed audit events, got %+v", action.AuditEvents)
	}
}

func TestApplyManualStockActionDoesNotInvalidateFailedOrRejectedWrites(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	t.Run("persistence failure", func(t *testing.T) {
		invalidator := &recordingCacheInvalidator{}
		service := NewStockActionServiceWithInvalidator(&memoryActionStore{err: errors.New("rolled back")}, &fakeStockWriter{}, invalidator)
		if _, err := service.ApplyManual(context.Background(), applyInput("act-persist-fail", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh), 7, now)); err == nil {
			t.Fatal("ApplyManual() error = nil, want persistence failure")
		}
		if len(invalidator.classes) != 0 {
			t.Fatalf("invalidations=%v after persistence failure, want none", invalidator.classes)
		}
	})
	t.Run("provider rejection", func(t *testing.T) {
		invalidator := &recordingCacheInvalidator{}
		writer := &fakeStockWriter{result: domain.StockWriteResult{Status: domain.StockWriteResultRejected, Message: "rejected"}}
		service := NewStockActionServiceWithInvalidator(&memoryActionStore{}, writer, invalidator)
		action, err := service.ApplyManual(context.Background(), applyInput("act-provider-reject", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh), 7, now))
		if err != nil {
			t.Fatal(err)
		}
		if action.State == domain.StockActionStateApplied || len(invalidator.classes) != 0 {
			t.Fatalf("state=%s invalidations=%v, want non-applied and none", action.State, invalidator.classes)
		}
	})
}

func TestApplyManualStockActionDoesNotRepeatProviderIntentForSameActionID(t *testing.T) {
	now := time.Date(2026, 7, 9, 14, 0, 0, 0, time.UTC)
	fresh := now.Add(-time.Minute)
	writer := &fakeStockWriter{result: domain.StockWriteResult{Status: domain.StockWriteResultApplied, IdempotencyKey: "act-repeat"}}
	store := &memoryActionStore{}
	service := NewStockActionService(store, writer)
	input := applyInput("act-repeat", actionRiskInput(now, domain.LinkStateResolved, intPtr(8), intPtr(9), &fresh, &fresh), 7, now)

	first, err := service.ApplyManual(context.Background(), input)
	if err != nil {
		t.Fatalf("first ApplyManual() error = %v", err)
	}
	second, err := service.ApplyManual(context.Background(), input)
	if err != nil {
		t.Fatalf("second ApplyManual() error = %v", err)
	}

	if writer.calls != 1 {
		t.Fatalf("writer calls=%d, want 1", writer.calls)
	}
	if first.ActionID != second.ActionID || second.State != domain.StockActionStateApplied {
		t.Fatalf("duplicate did not return existing action: first=%+v second=%+v", first, second)
	}
}

func actionRiskInput(now time.Time, linkState domain.LinkState, internalQty *int, providerQty *int, internalObservedAt *time.Time, providerObservedAt *time.Time) domain.RiskInput {
	return domain.RiskInput{
		Policy: domain.DefaultStockPolicy(),
		Link: domain.LinkEvidence{
			State:             linkState,
			InternalProductID: intPtr(20318),
		},
		InternalStock: domain.InternalStockEvidence{
			Quantity: internalQty,
			Source:   domain.SourceEvidence{ObservedAt: internalObservedAt},
		},
		ProviderStock: domain.ProviderStockEvidence{
			AvailableQuantity: providerQty,
			Source:            domain.SourceEvidence{ObservedAt: providerObservedAt},
		},
		Product: domain.ProductEvidence{ProductID: 20318},
		Now:     now,
	}
}

func applyInput(actionID string, risk domain.RiskInput, requested int, now time.Time) ApplyManualStockActionInput {
	return ApplyManualStockActionInput{
		ActionID:          actionID,
		Risk:              risk,
		RequestedQuantity: requested,
		Reason:            "operator confirmed Stock Seguro recommendation",
		Approval: domain.ManualApproval{
			Approved: true,
			Operator: domain.StockActionOperator{
				ActorType: "operator",
				ActorID:   "leandro",
				ActorName: "Leandro",
			},
			ApprovedAt: now,
		},
		ProviderRef: domain.ProviderStockRef{
			TenantID:            "tenant-default",
			InstallationID:      "inst-ml",
			ProviderCode:        "mercado_livre",
			ProviderAccountID:   "seller-1",
			ProviderItemID:      "MLB123",
			ProviderVariationID: "987",
		},
		Now: now,
	}
}

func intPtr(value int) *int {
	return &value
}

func hasActionEvent(action domain.StockAction, state domain.StockActionState) bool {
	for _, event := range action.AuditEvents {
		if event.State == state {
			return true
		}
	}
	return false
}

type fakeStockWriter struct {
	calls       int
	lastRequest domain.StockWriteRequest
	result      domain.StockWriteResult
	err         error
}

func (w *fakeStockWriter) UpdateAvailableQuantity(_ context.Context, request domain.StockWriteRequest) (domain.StockWriteResult, error) {
	w.calls++
	w.lastRequest = request
	if w.err != nil {
		return domain.StockWriteResult{}, w.err
	}
	return w.result, nil
}

type memoryActionStore struct {
	items map[string]domain.StockAction
	err   error
}

type recordingCacheInvalidator struct{ classes []string }

func (r *recordingCacheInvalidator) InvalidateClass(class string) { r.classes = append(r.classes, class) }

var _ internalreadports.CacheInvalidator = (*recordingCacheInvalidator)(nil)

func (s *memoryActionStore) FindByID(_ context.Context, actionID string) (domain.StockAction, bool, error) {
	if s.err != nil {
		return domain.StockAction{}, false, s.err
	}
	if s.items == nil {
		return domain.StockAction{}, false, nil
	}
	item, ok := s.items[actionID]
	return item, ok, nil
}

func (s *memoryActionStore) Save(_ context.Context, action domain.StockAction) error {
	if s.err != nil {
		return s.err
	}
	if action.ActionID == "" {
		return errors.New("missing action id")
	}
	if s.items == nil {
		s.items = map[string]domain.StockAction{}
	}
	s.items[action.ActionID] = action
	return nil
}
