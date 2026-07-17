package application

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

type retryFake struct {
	protocol ports.Protocol
	found    bool
	eligible bool
	clone    ports.Protocol
	cloneAt  time.Time
}

func (f *retryFake) GetProtocol(context.Context, string) (ports.Protocol, bool, error) {
	return f.protocol, f.found, nil
}
func (f *retryFake) CloneRetry(_ context.Context, _ string, at time.Time) (ports.Protocol, bool, error) {
	f.cloneAt = at
	return f.clone, f.eligible, nil
}

func TestRetryClonesEligibleTerminalProtocolWithoutMutatingOriginal(t *testing.T) {
	now := time.Date(2026, 7, 16, 15, 0, 0, 0, time.FixedZone("BRT", -3*60*60))
	original := ports.Protocol{ProtocolID: "MP-000007", State: domain.ProtocolStatePartiallyFailed, Intent: []byte(`{ "price": "10.00" }`), Selection: []byte(`{"mode":"filter"}`)}
	snapshotIntent, snapshotSelection := append([]byte(nil), original.Intent...), append([]byte(nil), original.Selection...)
	from := original.ProtocolID
	fake := &retryFake{protocol: original, found: true, eligible: true, clone: ports.Protocol{ProtocolID: "MP-000008", State: domain.ProtocolStateDraft, RetriedFrom: &from, Intent: append([]byte(nil), original.Intent...), Selection: append([]byte(nil), original.Selection...)}}
	got, err := NewRetryService(fake, func() time.Time { return now }).Retry(context.Background(), original.ProtocolID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ProtocolID == original.ProtocolID || got.RetriedFrom == nil || *got.RetriedFrom != original.ProtocolID || got.State != domain.ProtocolStateDraft {
		t.Fatalf("clone=%+v", got)
	}
	if !bytes.Equal(original.Intent, snapshotIntent) || !bytes.Equal(original.Selection, snapshotSelection) || fake.cloneAt != now.UTC() {
		t.Fatalf("original changed or clock not normalized: original=%+v at=%s", original, fake.cloneAt)
	}
}

func TestRetryRejectsNothingEligibleAndNonTerminalState(t *testing.T) {
	fake := &retryFake{found: true, eligible: false, protocol: ports.Protocol{State: domain.ProtocolStateFailedPreserved}}
	_, err := NewRetryService(fake, time.Now).Retry(context.Background(), "MP-1")
	var gate *GateError
	if !errors.As(err, &gate) || gate.Code != FailureCodeNothingToRetry {
		t.Fatalf("error=%v", err)
	}
	fake.protocol.State = domain.ProtocolStateApplying
	_, err = NewRetryService(fake, time.Now).Retry(context.Background(), "MP-1")
	if !errors.Is(err, domain.ErrInvalidStateTransition) {
		t.Fatalf("error=%v", err)
	}
}

func TestRetryUnknownProtocol(t *testing.T) {
	_, err := NewRetryService(&retryFake{}, time.Now).Retry(context.Background(), "MP-missing")
	if !errors.Is(err, ErrProtocolNotFound) {
		t.Fatalf("error=%v", err)
	}
}
