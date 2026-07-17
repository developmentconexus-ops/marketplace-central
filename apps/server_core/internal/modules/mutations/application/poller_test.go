package application

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/mutations/adapters/stub"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

func TestPollerPassTerminalStates(t *testing.T) {
	tests := []struct {
		name    string
		program map[string]stub.Result
		want    domain.ProtocolState
	}{
		{"all success", nil, domain.ProtocolStateApplied},
		{"mixed", map[string]stub.Result{"p:b": {Failure: &domain.Failure{Code: domain.FailureCodeProviderValidation, MessagePT: "inválido"}}}, domain.ProtocolStatePartiallyFailed},
		{"all fail", map[string]stub.Result{"p:a": {Failure: &domain.Failure{Code: domain.FailureCodeProviderAuth}}, "p:b": {Failure: &domain.Failure{Code: domain.FailureCodeProviderValidation}}}, domain.ProtocolStateFailedPreserved},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeRepo("p:a", "p:b")
			w := stub.NewWriter(tt.program)
			if worked, err := NewPoller(r, w, time.Now).Pass(context.Background(), "inst"); err != nil || !worked {
				t.Fatalf("Pass() worked=%v err=%v", worked, err)
			}
			if r.claim.finished != tt.want || !r.claim.committed || r.claim.rolledBack {
				t.Fatalf("finish=%q committed=%v rolledBack=%v", r.claim.finished, r.claim.committed, r.claim.rolledBack)
			}
			if !equalStrings(r.applying, []string{"item-1", "item-2"}) {
				t.Fatalf("applying marks=%v, want every sent item marked applying before send", r.applying)
			}
			if tt.name == "mixed" {
				var f domain.Failure
				if err := json.Unmarshal(r.outcomes["item-2"].Failure, &f); err != nil || f.Code != domain.FailureCodeProviderValidation {
					t.Fatalf("failure=%+v err=%v", f, err)
				}
			}
		})
	}
}

func TestPollerOwnsGateSequenceExactlyOncePerItem(t *testing.T) {
	r := newFakeRepo("p:a", "p:b")
	r.applied = []string{"p:a"}
	w := stub.NewWriter(map[string]stub.Result{"p:b": {Err: errors.New("token=secret upstream dump")}})
	if _, err := NewPoller(r, w, time.Now).Pass(context.Background(), "inst"); err != nil {
		t.Fatal(err)
	}
	if got := w.Keys(); len(got) != 1 || got[0] != "p:b" {
		t.Fatalf("writer keys=%v", got)
	}
	if r.outcomes["item-1"].State != domain.ItemStateSkipped {
		t.Fatalf("duplicate outcome=%+v", r.outcomes["item-1"])
	}
	if !equalStrings(r.applying, []string{"item-2"}) {
		t.Fatalf("applying marks=%v, skipped item must never be marked applying", r.applying)
	}
	if r.appliedReads != 1 {
		t.Fatalf("applied idempotency reads=%d, want one claim-backed read for the chunk", r.appliedReads)
	}
	var f struct {
		Code            domain.FailureCode `json:"code"`
		MessagePT       string             `json:"message_pt"`
		MessageProvider string             `json:"message_provider"`
		Retryable       bool               `json:"retryable"`
	}
	if err := json.Unmarshal(r.outcomes["item-2"].Failure, &f); err != nil {
		t.Fatal(err)
	}
	if f.Code != domain.FailureCodeInternal || f.Retryable || f.MessagePT != "Falha interna ao aplicar alteração." || f.MessageProvider != "" {
		t.Fatalf("failure=%+v", f)
	}
}

func TestPollerRejectsProtocolWithoutFreshSource(t *testing.T) {
	now := time.Date(2026, 7, 17, 12, 0, 0, 0, time.UTC)
	for _, tt := range []struct {
		name   string
		source *time.Time
		code   domain.FailureCode
	}{
		{name: "missing", code: FailureCodeSourceTimeUnavailable},
		{name: "stale", source: timePtr(now.Add(-15*time.Minute - time.Nanosecond)), code: domain.FailureCodeStaleSource},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := newFakeRepo("p:a")
			r.protocol.SourceAsOf = tt.source
			worked, err := NewPoller(r, stub.NewWriter(nil), func() time.Time { return now }).Pass(context.Background(), "inst")
			if !worked {
				t.Fatal("Pass() worked=false, want claimed protocol")
			}
			assertGateCode(t, err, tt.code)
			if r.claim.committed || len(r.outcomes) != 0 {
				t.Fatalf("committed=%v outcomes=%v", r.claim.committed, r.outcomes)
			}
		})
	}
}

func timePtr(value time.Time) *time.Time { return &value }

func TestPollerPassCrashResumeDoesNotResendDurableAppliedItems(t *testing.T) {
	r := newFakeRepo("p:a", "p:b", "p:c")
	r.failOutcomeAfter = 1
	w := stub.NewWriter(nil)
	p := NewPoller(r, w, time.Now)

	if worked, err := p.Pass(context.Background(), "inst"); err == nil || !worked {
		t.Fatalf("first Pass() worked=%v err=%v, want simulated crash error", worked, err)
	}
	if r.protocolState != domain.ProtocolStateApproved || r.outcomes["item-1"].State != domain.ItemStateApplied {
		t.Fatalf("after crash protocol=%q outcomes=%v", r.protocolState, r.outcomes)
	}
	r.failOutcomeAfter = 0
	if worked, err := p.Pass(context.Background(), "inst"); err != nil || !worked {
		t.Fatalf("resumed Pass() worked=%v err=%v", worked, err)
	}
	if got := w.Keys(); !equalStrings(got, []string{"p:a", "p:b", "p:c"}) {
		t.Fatalf("writer keys=%v, durable applied key was resent", got)
	}
	if r.protocolState != domain.ProtocolStateApplied || !r.claim.committed {
		t.Fatalf("terminal=%q committed=%v", r.protocolState, r.claim.committed)
	}
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

type fakeRepo struct {
	claim            *fakeClaim
	protocol         ports.Protocol
	protocolState    domain.ProtocolState
	items            []ports.MutationItem
	outcomes         map[string]ports.ItemOutcome
	applied          []string
	appliedReads     int
	applying         []string
	failOutcomeAfter int
}

func newFakeRepo(keys ...string) *fakeRepo {
	items := make([]ports.MutationItem, len(keys))
	for i, key := range keys {
		items[i] = ports.MutationItem{Seq: i + 1, ItemID: "item-" + string(rune('1'+i)), ListingID: key[2:], IdempotencyKey: key, After: json.RawMessage(`{"value":1}`)}
	}
	source := time.Now()
	protocol := ports.Protocol{ProtocolID: "p", InstallationID: "inst", Type: domain.ProtocolTypePriceUpdate, State: domain.ProtocolStateApproved, Actor: "operator", SourceAsOf: &source}
	return &fakeRepo{protocol: protocol, protocolState: domain.ProtocolStateApproved, items: items, outcomes: map[string]ports.ItemOutcome{}}
}
func (r *fakeRepo) CreateProtocol(context.Context, ports.CreateProtocolInput) (ports.Protocol, error) {
	panic("unused")
}
func (r *fakeRepo) GetProtocol(context.Context, string) (ports.Protocol, bool, error) {
	panic("unused")
}
func (r *fakeRepo) ReplaceItems(context.Context, string, []ports.ReplaceItemInput, *time.Time) ([]ports.MutationItem, error) {
	panic("unused")
}
func (r *fakeRepo) ApproveItems(context.Context, string, time.Time) error { panic("unused") }
func (r *fakeRepo) ClaimProtocol(context.Context, string) (ports.ProtocolClaim, bool, error) {
	r.protocolState = domain.ProtocolStateApplying
	r.protocol.State = domain.ProtocolStateApplying
	r.claim = &fakeClaim{repo: r, protocol: r.protocol}
	return r.claim, true, nil
}

type fakeClaim struct {
	repo                  *fakeRepo
	protocol              ports.Protocol
	finished              domain.ProtocolState
	writes                int
	committed, rolledBack bool
}

func (c *fakeClaim) Protocol() ports.Protocol { return c.protocol }
func (c *fakeClaim) FetchPendingItems(context.Context) ([]ports.MutationItem, error) {
	var items []ports.MutationItem
	for _, item := range c.repo.items {
		if _, terminal := c.repo.outcomes[item.ItemID]; !terminal {
			items = append(items, item)
		}
	}
	return items, nil
}
func (c *fakeClaim) MarkItemApplying(_ context.Context, id string) error {
	if _, terminal := c.repo.outcomes[id]; terminal {
		return errors.New("mutation item outcome is immutable")
	}
	c.repo.applying = append(c.repo.applying, id)
	return nil
}
func (c *fakeClaim) WriteItemOutcome(_ context.Context, id string, o ports.ItemOutcome) error {
	c.repo.outcomes[id] = o
	c.writes++
	if c.repo.failOutcomeAfter > 0 && c.writes == c.repo.failOutcomeAfter {
		return errors.New("simulated crash after durable outcome")
	}
	return nil
}
func (c *fakeClaim) AppliedIdempotencyKeys(context.Context, []string) ([]string, error) {
	c.repo.appliedReads++
	return c.repo.applied, nil
}
func (c *fakeClaim) ItemStateCounts(context.Context) (map[domain.ItemState]int, error) {
	counts := make(map[domain.ItemState]int)
	for _, outcome := range c.repo.outcomes {
		counts[outcome.State]++
	}
	return counts, nil
}
func (c *fakeClaim) Finish(_ context.Context, s domain.ProtocolState, _ time.Time) error {
	c.finished = s
	c.repo.protocolState = s
	return nil
}
func (c *fakeClaim) Commit(context.Context) error { c.committed = true; return nil }
func (c *fakeClaim) Rollback(context.Context) error {
	if c.committed {
		return nil
	}
	c.rolledBack = true
	c.repo.protocolState = domain.ProtocolStateApproved
	return nil
}
