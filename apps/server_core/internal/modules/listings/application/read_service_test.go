package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

type fakeRows struct {
	pages []ports.ListingRowPage
	calls int
}

func (f *fakeRows) ListListingRows(_ context.Context, _ ports.ListingQuery) (ports.ListingRowPage, error) {
	p := f.pages[f.calls]
	f.calls++
	return p, nil
}
func (*fakeRows) ListListingGroupRows(context.Context, ports.ListingGroupQuery) (ports.ListingGroupRowPage, error) {
	return ports.ListingGroupRowPage{}, nil
}
func (*fakeRows) GetListingRow(context.Context, domain.ListingKey) (domain.ListingReadModel, bool, error) {
	return domain.ListingReadModel{}, false, nil
}
func (*fakeRows) GetListingsSummary(context.Context, ports.SummaryQuery) (ports.ListingSummaryRow, error) {
	return ports.ListingSummaryRow{}, nil
}
func (*fakeRows) ListListingTimeline(context.Context, domain.ListingKey, int) ([]domain.TimelineEvent, error) {
	return nil, nil
}

type fakeFacts struct {
	costs                   map[int64]*ports.CostFact
	ceilings                map[int64]*ports.ICMSCeiling
	err                     error
	costCalls, ceilingCalls int
}

func (f *fakeFacts) GetCostFactsByIDs(context.Context, []int64) (map[int64]*ports.CostFact, error) {
	f.costCalls++
	return f.costs, f.err
}
func (f *fakeFacts) GetICMSCeilingByOrigin(context.Context, int64) (map[int64]*ports.ICMSCeiling, error) {
	f.ceilingCalls++
	return f.ceilings, f.err
}

type fakePolicy struct{ found bool }
type fakeInstallation bool

func (f fakeInstallation) InstallationExists(context.Context, string) (bool, error) {
	return bool(f), nil
}

func (f fakePolicy) GetPricingPolicyForInstallation(context.Context, string) (ports.PricingPolicy, bool, error) {
	return ports.PricingPolicy{MinMarginPercent: 10}, f.found, nil
}
func money(v string) *domain.Money { return &domain.Money{Amount: v, Currency: "BRL"} }
func resolved(id string) domain.ListingReadModel {
	return domain.ListingReadModel{ListingID: id, InstallationID: "i", ProviderListingID: id, Title: id, Link: domain.ListingLink{State: domain.LinkStateResolved, ProductID: &id}, Price: money("90"), ICMSWorstCaseByUF: nil}
}

func TestReadServiceFastPathUsesOnePageAndEnrichesDisplay(t *testing.T) {
	r := &fakeRows{pages: []ports.ListingRowPage{{Items: []domain.ListingReadModel{resolved("42664")}}}}
	f := &fakeFacts{costs: map[int64]*ports.CostFact{42664: {Amount: ptr(65.6597)}}, ceilings: map[int64]*ports.ICMSCeiling{8: {Percent: ptr(22)}}}
	s := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), func() time.Time { return time.Date(2026, 7, 15, 12, 0, 0, 0, time.FixedZone("x", 3600)) })
	p, err := s.List(context.Background(), ports.ListingQuery{InstallationID: "i", Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	if r.calls != 1 || f.costCalls != 1 || f.ceilingCalls != 1 {
		t.Fatalf("calls repo=%d cost=%d ceiling=%d", r.calls, f.costCalls, f.ceilingCalls)
	}
	if p.Items[0].BelowMarginWorstCase == nil || !*p.Items[0].BelowMarginWorstCase {
		t.Fatalf("golden below margin=%v", p.Items[0].BelowMarginWorstCase)
	}
	if p.Items[0].Cost == nil || p.Items[0].Cost.Currency != "BRL" {
		t.Fatalf("paired cost=%+v", p.Items[0].Cost)
	}
	if p.AsOf.Location() != time.UTC {
		t.Fatalf("as_of=%v", p.AsOf)
	}
}

func TestReadServiceScanCapReturnsLastScannedCursor(t *testing.T) {
	pages := make([]ports.ListingRowPage, maxBelowMarginScanPages)
	for i := range pages {
		c := ports.ListingCursor{LastTitle: "z", ListingID: "i~z~-"}
		pages[i] = ports.ListingRowPage{Items: []domain.ListingReadModel{resolved("42664")}, NextCursor: &c}
	}
	r := &fakeRows{pages: pages}
	f := &fakeFacts{costs: map[int64]*ports.CostFact{42664: {Amount: ptr(1)}}, ceilings: map[int64]*ports.ICMSCeiling{8: {Percent: ptr(22)}}}
	s := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), time.Now)
	p, err := s.List(context.Background(), ports.ListingQuery{InstallationID: "i", Limit: 2, Filter: domain.ListingFilter{Exception: domain.ListingExceptionBelowMargin}})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Items) != 0 || p.NextCursor == nil || r.calls != maxBelowMarginScanPages {
		t.Fatalf("items=%d cursor=%v calls=%d", len(p.Items), p.NextCursor, r.calls)
	}
}

func TestReadServiceOracleFailureFailsRequest(t *testing.T) {
	r := &fakeRows{pages: []ports.ListingRowPage{{Items: []domain.ListingReadModel{resolved("42664")}}}}
	f := &fakeFacts{err: errors.New("oracle down")}
	_, err := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), time.Now).List(context.Background(), ports.ListingQuery{InstallationID: "i", Limit: 1})
	if err == nil {
		t.Fatal("expected source failure")
	}
}

func TestReadServiceRejectsUnknownInstallationBeforeReads(t *testing.T) {
	r := &fakeRows{}
	f := &fakeFacts{}
	_, err := NewReadService(r, f, fakePolicy{}, fakeInstallation(false), time.Now).List(context.Background(), ports.ListingQuery{InstallationID: "missing", Limit: 1})
	if !errors.Is(err, ErrInstallationNotFound) || r.calls != 0 || f.ceilingCalls != 0 {
		t.Fatalf("err=%v repo=%d ceiling=%d", err, r.calls, f.ceilingCalls)
	}
}

func TestBelowMarginTriState(t *testing.T) {
	if got := belowMargin(money("90"), ptr(65.6597), ptr(22), ptr(10)); got == nil || !*got {
		t.Fatalf("golden=%v", got)
	}
	if got := belowMargin(money("100"), ptr(10), ptr(22), ptr(10)); got == nil || *got {
		t.Fatalf("false=%v", got)
	}
	if got := belowMargin(money("100"), nil, ptr(22), ptr(10)); got != nil {
		t.Fatalf("unknown=%v", got)
	}
}
func ptr(v float64) *float64 { return &v }
