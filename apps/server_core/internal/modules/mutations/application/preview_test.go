package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsadapter "marketplace-central/apps/server_core/internal/modules/mutations/adapters/listings"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
)

func TestPreviewReplacesSnapshotAfterListingChanges(t *testing.T) {
	ctx := context.Background()
	firstSource := time.Date(2026, 7, 16, 10, 0, 0, 0, time.UTC)
	secondSource := firstSource.Add(time.Minute)
	store := &previewStore{protocol: testProtocol(domain.ProtocolStateDraft)}
	resolver := &previewResolver{ids: []string{"inst-1~MLB1~-"}}
	reader := &previewListingReader{models: []listingsdomain.ListingReadModel{
		{ListingID: "inst-1~MLB1~-", Price: &listingsdomain.Money{Amount: "89.00", Currency: "BRL"}, FetchedAt: &firstSource},
		{ListingID: "inst-1~MLB1~-", Price: &listingsdomain.Money{Amount: "92.00", Currency: "BRL"}, FetchedAt: &secondSource},
	}}
	now := firstSource.Add(2 * time.Minute)
	service := NewService(store, resolver, reader, func() time.Time { return now })

	first, err := service.Preview(ctx, store.protocol.ProtocolID)
	if err != nil {
		t.Fatal(err)
	}
	store.protocol = first.Protocol
	second, err := service.Preview(ctx, store.protocol.ProtocolID)
	if err != nil {
		t.Fatal(err)
	}
	if store.replaceCalls != 2 || len(store.items) != 1 {
		t.Fatalf("replace calls=%d items=%d", store.replaceCalls, len(store.items))
	}
	assertJSON(t, store.items[0].Before, `{"price":{"amount":"92.00","currency":"BRL"}}`)
	if second.Protocol.SourceAsOf == nil || !second.Protocol.SourceAsOf.Equal(secondSource) {
		t.Fatalf("source_as_of=%v, want %v", second.Protocol.SourceAsOf, secondSource)
	}
}

func TestPreviewRejectsMissingSourceTimeWithoutReplacingItems(t *testing.T) {
	store := &previewStore{protocol: testProtocol(domain.ProtocolStateDraft)}
	service := NewService(store, &previewResolver{ids: []string{"inst-1~MLB1~-"}}, &previewListingReader{models: []listingsdomain.ListingReadModel{{ListingID: "inst-1~MLB1~-"}}}, time.Now)
	_, err := service.Preview(context.Background(), store.protocol.ProtocolID)
	assertPreviewGateCode(t, err, FailureCodeSourceTimeUnavailable)
	if store.replaceCalls != 0 {
		t.Fatalf("ReplacePreview calls=%d, want 0", store.replaceCalls)
	}
}

func TestPreviewRejectsEmptyAndOversizedSelections(t *testing.T) {
	for _, tc := range []struct {
		name  string
		count int
		code  domain.FailureCode
	}{{"empty", 0, FailureCodeEmptySelection}, {"2001 cap", 2001, FailureCodeSelectionTooLarge}} {
		t.Run(tc.name, func(t *testing.T) {
			ids := make([]string, tc.count)
			for i := range ids {
				ids[i] = fmt.Sprintf("inst-1~MLB%d~-", i)
			}
			store := &previewStore{protocol: testProtocol(domain.ProtocolStateDraft)}
			reader := &previewListingReader{}
			_, err := NewService(store, &previewResolver{ids: ids}, reader, time.Now).Preview(context.Background(), store.protocol.ProtocolID)
			assertPreviewGateCode(t, err, tc.code)
			if reader.calls != 0 || store.replaceCalls != 0 {
				t.Fatalf("listing reads=%d replacements=%d", reader.calls, store.replaceCalls)
			}
		})
	}
}

func TestCreateValidatesAndPersistsSelectionAsDraft(t *testing.T) {
	store := &previewStore{}
	now := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	input := CreateInput{InstallationID: "inst-1", Type: domain.ProtocolTypePriceUpdate, Actor: "operator_supplied_unverified", Intent: json.RawMessage(`{"new_price":{"amount":"49.90","currency":"BRL"}}`), Selection: json.RawMessage(`{"mode":"filter","filter":{"status":"active"}}`)}
	created, err := NewService(store, &previewResolver{}, &previewListingReader{}, func() time.Time { return now }).Create(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if created.State != domain.ProtocolStateDraft || created.InstallationID != "inst-1" || !created.CreatedAt.Equal(now) {
		t.Fatalf("created=%+v", created)
	}
}

func TestCreateRejectsBlankActorThroughGate(t *testing.T) {
	store := &previewStore{}
	input := CreateInput{InstallationID: "inst-1", Type: domain.ProtocolTypePriceUpdate, Actor: "   ", Intent: json.RawMessage(`{}`), Selection: json.RawMessage(`{"mode":"filter","filter":{"status":"active"}}`)}
	_, err := NewService(store, &previewResolver{}, &previewListingReader{}, time.Now).Create(context.Background(), input)
	assertPreviewGateCode(t, err, FailureCodeActorRequired)
	if store.protocol.ProtocolID != "" {
		t.Fatalf("protocol persisted despite blank actor: %+v", store.protocol)
	}
}

type previewStore struct {
	protocol     ports.Protocol
	items        []ports.ReplaceItemInput
	replaceCalls int
}

func (s *previewStore) CreateProtocol(_ context.Context, input ports.CreateProtocolInput) (ports.Protocol, error) {
	s.protocol = ports.Protocol{ProtocolID: "MP-000001", InstallationID: input.InstallationID, Type: input.Type, State: domain.ProtocolStateDraft, Actor: input.Actor, Intent: input.Intent, Selection: input.Selection, CreatedAt: input.CreatedAt}
	return s.protocol, nil
}
func (s *previewStore) GetProtocol(context.Context, string) (ports.Protocol, bool, error) {
	return s.protocol, s.protocol.ProtocolID != "", nil
}
func (s *previewStore) ReplacePreview(_ context.Context, _ string, inputs []ports.ReplaceItemInput, _ time.Time, _ time.Time) ([]ports.MutationItem, error) {
	s.replaceCalls++
	s.items = append([]ports.ReplaceItemInput(nil), inputs...)
	items := make([]ports.MutationItem, len(inputs))
	for i, input := range inputs {
		items[i] = ports.MutationItem{Seq: i + 1, ListingID: input.ListingID, Before: input.Before, After: input.After, State: domain.ItemStatePreviewed}
	}
	return items, nil
}

type previewResolver struct{ ids []string }

func (r *previewResolver) Resolve(context.Context, string, listingsadapter.Selection) ([]string, error) {
	return append([]string(nil), r.ids...), nil
}

type previewListingReader struct {
	models []listingsdomain.ListingReadModel
	calls  int
}

func (r *previewListingReader) Get(context.Context, listingsdomain.ListingID) (listingsdomain.ListingReadModel, []listingsdomain.TimelineEvent, error) {
	model := r.models[r.calls]
	r.calls++
	return model, nil, nil
}

func testProtocol(state domain.ProtocolState) ports.Protocol {
	return ports.Protocol{ProtocolID: "MP-000001", InstallationID: "inst-1", Type: domain.ProtocolTypePriceUpdate, State: state, Intent: json.RawMessage(`{"new_price":{"amount":"75.40","currency":"BRL"}}`), Selection: json.RawMessage(`{"mode":"explicit","listing_ids":["inst-1~MLB1~-"]}`)}
}

func assertPreviewGateCode(t *testing.T, err error, code domain.FailureCode) {
	t.Helper()
	var gate *GateError
	if !errors.As(err, &gate) || gate.Code != code {
		t.Fatalf("error=%v, want gate code %q", err, code)
	}
}

func assertJSON(t *testing.T, got json.RawMessage, want string) {
	t.Helper()
	var gotValue, wantValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatal(err)
	}
	gotJSON, _ := json.Marshal(gotValue)
	wantJSON, _ := json.Marshal(wantValue)
	if string(gotJSON) != string(wantJSON) {
		t.Fatalf("JSON=%s, want %s", gotJSON, wantJSON)
	}
}
