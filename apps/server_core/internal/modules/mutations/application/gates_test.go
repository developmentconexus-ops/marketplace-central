package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
)

func TestActorGateRejectsBlankBeforeInsert(t *testing.T) {
	inserted := false
	err := InsertWithActorGate(context.Background(), " \t", func(context.Context) error { inserted = true; return nil })
	assertGateCode(t, err, FailureCodeActorRequired)
	if inserted {
		t.Fatal("insert called for blank actor")
	}
}

func TestIdempotencyGateSkipsAppliedKeyWithoutWriterCall(t *testing.T) {
	written := false
	skipped, err := WriteThroughGates(context.Background(), "MP-1:list-1",
		func(context.Context, string) (bool, error) { return true, nil },
		func(context.Context) error { t.Fatal("audit called for duplicate"); return nil },
		func(context.Context) error { written = true; return nil })
	if err != nil || !skipped || written {
		t.Fatalf("skipped=%v written=%v err=%v", skipped, written, err)
	}
}

func TestExecuteGateRejectsAnythingExceptLiteralTrue(t *testing.T) {
	for _, execute := range []*bool{nil, boolPtr(false)} {
		assertGateCode(t, RequireExecute(execute), FailureCodeExecuteRequired)
	}
	if err := RequireExecute(boolPtr(true)); err != nil {
		t.Fatalf("literal true rejected: %v", err)
	}
}

func TestResolvedLinkGateRejectsWriteIntentsWithoutLink(t *testing.T) {
	for _, typ := range []domain.ProtocolType{domain.ProtocolTypePriceUpdate, domain.ProtocolTypeStockCorrect, domain.ProtocolTypeListingEdit} {
		err := RequireResolvedLink(context.Background(), typ, "list-1", func(context.Context, string) (bool, error) { return false, nil })
		assertGateCode(t, err, domain.FailureCodeLinkUnresolved)
	}
}

func TestPolicyGateRejectsPriceAndStockWithoutPolicy(t *testing.T) {
	for _, typ := range []domain.ProtocolType{domain.ProtocolTypePriceUpdate, domain.ProtocolTypeStockCorrect} {
		called := 0
		err := RequirePolicy(context.Background(), typ, "list-1", func(context.Context, string) (bool, error) { called++; return false, nil })
		assertGateCode(t, err, domain.FailureCodePolicyMissing)
		if called != 1 {
			t.Fatalf("policy reader calls=%d", called)
		}
	}
}

func TestSourceTimestampGateRejectsAbsentAndStale(t *testing.T) {
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	assertGateCode(t, RequireFreshSource(nil, now, 15*time.Minute), FailureCodeSourceTimeUnavailable)
	stale := now.Add(-15*time.Minute - time.Nanosecond)
	assertGateCode(t, RequireFreshSource(&stale, now, 15*time.Minute), domain.FailureCodeStaleSource)
}

func TestBeforeAfterAuditGateBlocksProviderWhenPersistenceFails(t *testing.T) {
	auditErr := errors.New("audit unavailable")
	written := false
	skipped, err := WriteThroughGates(context.Background(), "MP-1:list-1",
		func(context.Context, string) (bool, error) { return false, nil },
		func(context.Context) error { return auditErr },
		func(context.Context) error { written = true; return nil })
	if !errors.Is(err, auditErr) || skipped || written {
		t.Fatalf("skipped=%v written=%v err=%v", skipped, written, err)
	}
}

func assertGateCode(t *testing.T, err error, want domain.FailureCode) {
	t.Helper()
	var gateErr *GateError
	if !errors.As(err, &gateErr) || gateErr.Code != want {
		t.Fatalf("error=%v, want GateError code %q", err, want)
	}
	failure := gateErr.Failure()
	if failure.Code != want || failure.Retryable != want.RetryableByDefault() {
		t.Fatalf("failure=%+v", failure)
	}
}

func boolPtr(value bool) *bool { return &value }
