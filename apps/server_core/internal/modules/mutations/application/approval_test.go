package application

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

func TestApproveRejectsExecuteAbsentFalseAndCoercionAttempts(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	for _, payload := range []string{`{}`, `{"execute":false}`, `{"execute":"true"}`, `{"execute":1}`, `{"execute":true,"extra":1}`, `{"execute":true} {}`} {
		t.Run(payload, func(t *testing.T) {
			store := approvalTestStore{protocol: previewedProtocol(now.Add(-time.Minute))}
			_, err := NewApprovalService(&store, func() time.Time { return now }).Approve(context.Background(), "MP-000001", json.RawMessage(payload))
			assertGateCode(t, err, FailureCodeExecuteRequired)
			if store.approveCalls != 0 {
				t.Fatalf("ApproveItems calls=%d", store.approveCalls)
			}
		})
	}
}

func TestApproveExactly15MinutesIsNotStale(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	store := approvalTestStore{protocol: previewedProtocol(now.Add(-15 * time.Minute))}
	got, err := NewApprovalService(&store, func() time.Time { return now }).Approve(context.Background(), "MP-000001", json.RawMessage(`{"execute":true}`))
	if err != nil || got.State != domain.ProtocolStateApproved || store.approveCalls != 1 {
		t.Fatalf("state=%q calls=%d err=%v", got.State, store.approveCalls, err)
	}
}

func TestApproveOver15MinutesReturnsPreviewStale(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	store := approvalTestStore{protocol: previewedProtocol(now.Add(-15*time.Minute - time.Nanosecond))}
	_, err := NewApprovalService(&store, func() time.Time { return now }).Approve(context.Background(), "MP-000001", json.RawMessage(`{"execute":true}`))
	assertGateCode(t, err, FailureCodePreviewStale)
	if store.approveCalls != 0 {
		t.Fatalf("ApproveItems calls=%d", store.approveCalls)
	}
}

func TestApprovePromotesExistingSnapshotWithoutListingReads(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	store := approvalTestStore{protocol: previewedProtocol(now.Add(-time.Minute))}
	_, err := NewApprovalService(&store, func() time.Time { return now }).Approve(context.Background(), "MP-000001", json.RawMessage(`{"execute":true}`))
	if err != nil || store.approveCalls != 1 || store.approvedID != "MP-000001" || !store.approvedAt.Equal(now) {
		t.Fatalf("id=%q at=%v calls=%d err=%v", store.approvedID, store.approvedAt, store.approveCalls, err)
	}
}

func TestCancelAfterApprovedReturnsTypedIllegalState(t *testing.T) {
	store := approvalTestStore{protocol: ports.Protocol{ProtocolID: "MP-000001", State: domain.ProtocolStateApproved}}
	_, err := NewApprovalService(&store, time.Now).Cancel(context.Background(), "MP-000001")
	var stateErr *domain.InvalidStateTransitionError
	if !errors.As(err, &stateErr) || stateErr.From != domain.ProtocolStateApproved || stateErr.To != domain.ProtocolStateCancelled {
		t.Fatalf("error=%T %v", err, err)
	}
	if store.cancelCalls != 0 {
		t.Fatalf("CancelProtocol calls=%d", store.cancelCalls)
	}
}

func TestCancelPersistsInjectedTimestamp(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.FixedZone("BRT", -3*60*60))
	store := approvalTestStore{protocol: ports.Protocol{ProtocolID: "MP-000001", State: domain.ProtocolStateDraft}}
	got, err := NewApprovalService(&store, func() time.Time { return now }).Cancel(context.Background(), "MP-000001")
	if err != nil || got.FinishedAt == nil || !got.FinishedAt.Equal(now) || !store.cancelledAt.Equal(now.UTC()) {
		t.Fatalf("protocol=%+v persisted=%v err=%v", got, store.cancelledAt, err)
	}
}

type approvalTestStore struct {
	protocol     ports.Protocol
	approveCalls int
	cancelCalls  int
	approvedID   string
	approvedAt   time.Time
	cancelledAt  time.Time
}

func (s *approvalTestStore) GetProtocol(context.Context, string) (ports.Protocol, bool, error) {
	return s.protocol, s.protocol.ProtocolID != "", nil
}
func (s *approvalTestStore) ApproveItems(_ context.Context, id string, at time.Time) error {
	s.approveCalls++
	s.approvedID, s.approvedAt = id, at
	return nil
}
func (s *approvalTestStore) CancelProtocol(_ context.Context, _ string, at time.Time) error {
	s.cancelCalls++
	s.cancelledAt = at
	return nil
}

func previewedProtocol(previewedAt time.Time) ports.Protocol {
	return ports.Protocol{ProtocolID: "MP-000001", State: domain.ProtocolStatePreviewed, PreviewedAt: &previewedAt}
}
