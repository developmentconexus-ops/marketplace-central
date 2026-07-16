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
			if tt.name == "mixed" {
				var f domain.Failure
				if err := json.Unmarshal(r.claim.outcomes["item-2"].Failure, &f); err != nil || f.Code != domain.FailureCodeProviderValidation {
					t.Fatalf("failure=%+v err=%v", f, err)
				}
			}
		})
	}
}

func TestPollerPassSkipsAppliedKeyAndSanitizesUnknownError(t *testing.T) {
	r := newFakeRepo("p:a", "p:b")
	r.claim.applied = []string{"p:a"}
	w := stub.NewWriter(map[string]stub.Result{"p:b": {Err: errors.New("token=secret upstream dump")}})
	if _, err := NewPoller(r, w, time.Now).Pass(context.Background(), "inst"); err != nil {
		t.Fatal(err)
	}
	if got := w.Keys(); len(got) != 1 || got[0] != "p:b" {
		t.Fatalf("writer keys=%v", got)
	}
	if r.claim.outcomes["item-1"].State != domain.ItemStateSkipped {
		t.Fatalf("duplicate outcome=%+v", r.claim.outcomes["item-1"])
	}
	var f struct {
		Code            domain.FailureCode `json:"code"`
		MessagePT       string             `json:"message_pt"`
		MessageProvider string             `json:"message_provider"`
		Retryable       bool               `json:"retryable"`
	}
	if err := json.Unmarshal(r.claim.outcomes["item-2"].Failure, &f); err != nil {
		t.Fatal(err)
	}
	if f.Code != domain.FailureCodeInternal || f.Retryable || f.MessagePT != "Falha interna ao aplicar alteração." || f.MessageProvider != "" {
		t.Fatalf("failure=%+v", f)
	}
}

type fakeRepo struct{ claim *fakeClaim }

func newFakeRepo(keys ...string) *fakeRepo {
	items := make([]ports.MutationItem, len(keys))
	for i, key := range keys {
		items[i] = ports.MutationItem{Seq: i + 1, ItemID: "item-" + string(rune('1'+i)), ListingID: key[2:], IdempotencyKey: key, After: json.RawMessage(`{"value":1}`)}
	}
	return &fakeRepo{&fakeClaim{protocol: ports.Protocol{ProtocolID: "p", InstallationID: "inst", Type: domain.ProtocolTypePriceUpdate, State: domain.ProtocolStateApplying}, items: items, outcomes: map[string]ports.ItemOutcome{}}}
}
func (r *fakeRepo) CreateProtocol(context.Context, ports.CreateProtocolInput) (ports.Protocol, error) {
	panic("unused")
}
func (r *fakeRepo) GetProtocol(context.Context, string) (ports.Protocol, bool, error) {
	panic("unused")
}
func (r *fakeRepo) ReplaceItems(context.Context, string, []ports.ReplaceItemInput) ([]ports.MutationItem, error) {
	panic("unused")
}
func (r *fakeRepo) ApproveItems(context.Context, string, time.Time) error { panic("unused") }
func (r *fakeRepo) ClaimProtocol(context.Context, string) (ports.ProtocolClaim, bool, error) {
	return r.claim, true, nil
}

type fakeClaim struct {
	protocol              ports.Protocol
	items                 []ports.MutationItem
	applied               []string
	outcomes              map[string]ports.ItemOutcome
	finished              domain.ProtocolState
	committed, rolledBack bool
}

func (c *fakeClaim) Protocol() ports.Protocol { return c.protocol }
func (c *fakeClaim) FetchPendingItems(context.Context) ([]ports.MutationItem, error) {
	items := c.items
	c.items = nil
	return items, nil
}
func (c *fakeClaim) WriteItemOutcome(_ context.Context, id string, o ports.ItemOutcome) error {
	c.outcomes[id] = o
	return nil
}
func (c *fakeClaim) AppliedIdempotencyKeys(context.Context, []string) ([]string, error) {
	return c.applied, nil
}
func (c *fakeClaim) Finish(_ context.Context, s domain.ProtocolState, _ time.Time) error {
	c.finished = s
	return nil
}
func (c *fakeClaim) Commit(context.Context) error   { c.committed = true; return nil }
func (c *fakeClaim) Rollback(context.Context) error { c.rolledBack = true; return nil }
