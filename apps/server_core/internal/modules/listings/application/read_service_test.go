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
	pages        []ports.ListingRowPage
	groupPages   []ports.ListingGroupRowPage
	calls        int
	groupQueries []ports.ListingGroupQuery
	getModel     domain.ListingReadModel
	getFound     bool
	getErr       error
	timeline     []domain.TimelineEvent
	summary      ports.ListingSummaryRow
	summaryErr   error
	summaryCalls int
}

func (f *fakeRows) ListListingRows(_ context.Context, _ ports.ListingQuery) (ports.ListingRowPage, error) {
	p := f.pages[f.calls]
	f.calls++
	return p, nil
}
func (f *fakeRows) ListListingGroupRows(_ context.Context, q ports.ListingGroupQuery) (ports.ListingGroupRowPage, error) {
	f.groupQueries = append(f.groupQueries, q)
	p := f.groupPages[len(f.groupQueries)-1]
	return p, nil
}
func (f *fakeRows) GetListingRow(context.Context, domain.ListingKey) (domain.ListingReadModel, bool, error) {
	return f.getModel, f.getFound, f.getErr
}
func (f *fakeRows) GetListingsSummary(context.Context, ports.SummaryQuery) (ports.ListingSummaryRow, error) {
	f.summaryCalls++
	return f.summary, f.summaryErr
}
func (f *fakeRows) ListListingTimeline(context.Context, domain.ListingKey, int) ([]domain.TimelineEvent, error) {
	return f.timeline, nil
}

type fakeFacts struct {
	costs                   map[int64]*ports.CostFact
	ceilings                map[int64]*ports.ICMSCeiling
	err                     error
	costErr                 error
	ceilingErr              error
	costCalls, ceilingCalls int
}

func (f *fakeFacts) GetCostFactsByIDs(context.Context, []int64) (map[int64]*ports.CostFact, error) {
	f.costCalls++
	if f.costErr != nil {
		return nil, f.costErr
	}
	return f.costs, f.err
}
func (f *fakeFacts) GetICMSCeilingByOrigin(context.Context, int64) (map[int64]*ports.ICMSCeiling, error) {
	f.ceilingCalls++
	if f.ceilingErr != nil {
		return nil, f.ceilingErr
	}
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

func TestReadServiceSummaryCountsMarginsAndBatchesCosts(t *testing.T) {
	r := &fakeRows{summary: ports.ListingSummaryRow{Total: 4, Active: 1, Paused: 1, SyncError: 1, Stale: 1, Unlinked: 1, Linked: []ports.SummaryLinkedRow{
		{CostID: 1, Price: money("90")}, {CostID: 2, Price: money("100")}, {CostID: 3, Price: nil}, {CostID: 1, Price: money("90")},
	}}}
	f := &fakeFacts{costs: map[int64]*ports.CostFact{1: {Amount: ptr(65.6597)}, 2: {Amount: ptr(10)}, 3: nil}, ceilings: map[int64]*ports.ICMSCeiling{8: {Percent: ptr(22)}}}
	now := time.Date(2026, 7, 15, 13, 0, 0, 0, time.FixedZone("x", 3600))
	got, err := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), func() time.Time { return now }).Summary(context.Background(), ports.SummaryQuery{InstallationID: "i"})
	if err != nil {
		t.Fatal(err)
	}
	if got.BelowMarginWorstCase == nil || *got.BelowMarginWorstCase != 2 || got.MarginUnknown == nil || *got.MarginUnknown != 1 {
		t.Fatalf("summary=%+v", got)
	}
	if got.Total != 4 || got.Active != 1 || got.Paused != 1 || got.SyncError != 1 || got.Stale != 1 || got.Unlinked != 1 {
		t.Fatalf("counters=%+v", got)
	}
	if f.costCalls != 1 || f.ceilingCalls != 1 || r.summaryCalls != 1 || !got.AsOf.Equal(now.UTC()) {
		t.Fatalf("calls cost=%d ceiling=%d repo=%d asof=%s", f.costCalls, f.ceilingCalls, r.summaryCalls, got.AsOf)
	}
}

func TestReadServiceSummaryUnknownInputsAndSourceOutage(t *testing.T) {
	base := ports.ListingSummaryRow{Total: 2, Unlinked: 1, Linked: []ports.SummaryLinkedRow{{CostID: 1, Price: money("90")}}}
	for _, tc := range []struct {
		name    string
		policy  bool
		facts   *fakeFacts
		wantNil bool
	}{
		{"missing policy", false, &fakeFacts{costs: map[int64]*ports.CostFact{1: {Amount: ptr(1)}}, ceilings: map[int64]*ports.ICMSCeiling{8: {Percent: ptr(22)}}}, false},
		{"missing cost", true, &fakeFacts{costs: map[int64]*ports.CostFact{}, ceilings: map[int64]*ports.ICMSCeiling{8: {Percent: ptr(22)}}}, false},
		{"cost outage", true, &fakeFacts{costErr: errors.New("down"), ceilings: map[int64]*ports.ICMSCeiling{8: {Percent: ptr(22)}}}, true},
		{"ceiling outage", true, &fakeFacts{ceilingErr: errors.New("down")}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewReadService(&fakeRows{summary: base}, tc.facts, fakePolicy{found: tc.policy}, fakeInstallation(true), time.Now).Summary(context.Background(), ports.SummaryQuery{InstallationID: "i"})
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantNil {
				if got.BelowMarginWorstCase != nil || got.MarginUnknown != nil {
					t.Fatalf("summary=%+v", got)
				}
			} else if got.MarginUnknown == nil || *got.MarginUnknown != 1 || got.BelowMarginWorstCase == nil || *got.BelowMarginWorstCase != 0 {
				t.Fatalf("summary=%+v", got)
			}
			if got.Unlinked != 1 {
				t.Fatalf("unlinked=%d", got.Unlinked)
			}
		})
	}
}

func TestReadServiceSummaryValidatesInstallationFirst(t *testing.T) {
	for _, tc := range []struct {
		id           string
		installation fakeInstallation
		want         error
	}{{"", true, ErrInstallationIDRequired}, {"i", false, ErrInstallationNotFound}} {
		r, f := &fakeRows{}, &fakeFacts{}
		_, err := NewReadService(r, f, fakePolicy{found: true}, tc.installation, time.Now).Summary(context.Background(), ports.SummaryQuery{InstallationID: tc.id})
		if !errors.Is(err, tc.want) || r.summaryCalls != 0 || f.costCalls != 0 || f.ceilingCalls != 0 {
			t.Fatalf("err=%v calls repo=%d cost=%d ceiling=%d", err, r.summaryCalls, f.costCalls, f.ceilingCalls)
		}
	}
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

func TestReadServiceGetNotFoundReturnsListingNotFoundWithoutEnrichment(t *testing.T) {
	r := &fakeRows{}
	f := &fakeFacts{}
	_, _, err := NewReadService(r, f, fakePolicy{}, fakeInstallation(true), time.Now).Get(context.Background(), domain.ListingID{InstallationID: "i", ProviderListingID: "missing", VariationID: "-"})
	var notFound *domain.ListingNotFoundError
	if !errors.As(err, &notFound) || r.getErr != nil || f.costCalls != 0 || f.ceilingCalls != 0 {
		t.Fatalf("err=%v repo_err=%v cost=%d ceiling=%d", err, r.getErr, f.costCalls, f.ceilingCalls)
	}
}

func TestReadServiceGetBuildsDeterministicICMSMatrixAndTimeline(t *testing.T) {
	model := resolved("i~42664~-")
	model.InstallationID = "i"
	model.ProviderListingID = "42664"
	model.Link.ProductID = ptrString("42664")
	r := &fakeRows{
		getModel: model,
		getFound: true,
		timeline: []domain.TimelineEvent{{Kind: "synced", MessagePT: "Sincronizado"}},
	}
	rate22, rateNil, rate10 := 22.0, (*float64)(nil), 10.0
	f := &fakeFacts{
		costs:    map[int64]*ports.CostFact{42664: {Amount: ptr(65.6597)}},
		ceilings: map[int64]*ports.ICMSCeiling{1: {Percent: rateNil}, 5: {Percent: &rate10}, 8: {Percent: &rate22}},
	}
	got, timeline, err := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), time.Now).Get(context.Background(), domain.ListingID{InstallationID: "i", ProviderListingID: "42664", VariationID: "-"})
	if err != nil {
		t.Fatal(err)
	}
	if len(timeline) != 1 || timeline[0].Kind != "synced" {
		t.Fatalf("timeline=%+v", timeline)
	}
	if got.ICMSWorstCaseByUF == nil || len(*got.ICMSWorstCaseByUF) != 3 {
		t.Fatalf("matrix=%+v", got.ICMSWorstCaseByUF)
	}
	rows := *got.ICMSWorstCaseByUF
	if rows[0].DestinationUF != "1" || rows[1].DestinationUF != "5" || rows[2].DestinationUF != "8" {
		t.Fatalf("matrix order=%+v", rows)
	}
	if rows[2].WorstCaseICMSPct == nil || *rows[2].WorstCaseICMSPct != "22" || rows[2].PriceNetBasis == nil || *rows[2].PriceNetBasis != "70.2" || rows[2].BelowMarginAtUF == nil || !*rows[2].BelowMarginAtUF {
		t.Fatalf("golden row=%+v", rows[2])
	}
	if rows[0].WorstCaseICMSPct != nil || rows[0].PriceNetBasis != nil || rows[0].BelowMarginAtUF != nil {
		t.Fatalf("nil ceiling row=%+v", rows[0])
	}
}

func TestReadServiceGetMissingCostLeavesMatrixMarginUnknown(t *testing.T) {
	model := resolved("i~42664~-")
	model.InstallationID = "i"
	model.ProviderListingID = "42664"
	model.Link.ProductID = ptrString("42664")
	rate := 22.0
	r := &fakeRows{getModel: model, getFound: true}
	f := &fakeFacts{ceilings: map[int64]*ports.ICMSCeiling{8: {Percent: &rate}}}
	got, _, err := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), time.Now).Get(context.Background(), domain.ListingID{InstallationID: "i", ProviderListingID: "42664", VariationID: "-"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ICMSWorstCaseByUF == nil || len(*got.ICMSWorstCaseByUF) != 1 || (*got.ICMSWorstCaseByUF)[0].BelowMarginAtUF != nil {
		t.Fatalf("matrix=%+v", got.ICMSWorstCaseByUF)
	}
}

func TestReadServiceGetOracleFailureFailsRequest(t *testing.T) {
	model := resolved("i~42664~-")
	r := &fakeRows{getModel: model, getFound: true}
	f := &fakeFacts{err: errors.New("oracle down")}
	_, _, err := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), time.Now).Get(context.Background(), domain.ListingID{InstallationID: "i", ProviderListingID: "42664", VariationID: "-"})
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

func TestByProductFastPathEnrichesAndComputesWorstChild(t *testing.T) {
	id := "42664"
	r := &fakeRows{groupPages: []ports.ListingGroupRowPage{{Groups: []domain.ListingGroup{{ProductID: &id, ProductTitle: ptrString("Mesmo"), Listings: []domain.ListingReadModel{resolved(id), {Link: domain.ListingLink{State: domain.LinkStateResolved}, SyncState: domain.ListingSyncStateError}}}}}}}
	f := &fakeFacts{costs: map[int64]*ports.CostFact{42664: {Amount: ptr(1)}}, ceilings: map[int64]*ports.ICMSCeiling{}}
	p, err := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), time.Now).ByProduct(context.Background(), ports.ListingGroupQuery{InstallationID: "i", Limit: 50})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.groupQueries) != 1 || f.costCalls != 1 || len(p.Groups) != 1 || p.Groups[0].GroupState != domain.GroupStateError || p.Groups[0].ListingCount != len(p.Groups[0].Listings) {
		t.Fatalf("page=%+v queries=%d costs=%d", p, len(r.groupQueries), f.costCalls)
	}
}

func TestByProductScanDropsGroupsAndHandlesSemProdutoDefinitionally(t *testing.T) {
	falseValue := false
	nullTitle := "sem produto"
	c1 := ports.GroupCursor{ProductTitle: "A", ProductID: "1"}
	r := &fakeRows{groupPages: []ports.ListingGroupRowPage{
		{Groups: []domain.ListingGroup{{ProductID: ptrString("1"), ProductTitle: ptrString("A"), Listings: []domain.ListingReadModel{resolved("1")}}}, NextCursor: &c1},
		{Groups: []domain.ListingGroup{{ProductTitle: &nullTitle, Listings: []domain.ListingReadModel{{Link: domain.ListingLink{State: domain.LinkStateUnresolved}}}}}},
	}}
	f := &fakeFacts{costs: map[int64]*ports.CostFact{1: {Amount: ptr(1)}}, ceilings: map[int64]*ports.ICMSCeiling{8: {Percent: ptr(22)}}}
	p, err := NewReadService(r, f, fakePolicy{found: true}, fakeInstallation(true), time.Now).ByProduct(context.Background(), ports.ListingGroupQuery{InstallationID: "i", Limit: 2, Filter: domain.ListingFilter{HasException: &falseValue}})
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Groups) != 1 || p.Groups[0].ProductID == nil || p.Groups[0].ListingCount != 1 || len(r.groupQueries) != 2 || !r.groupQueries[1].Cursor.IsFirstPage() && r.groupQueries[1].Cursor != c1 {
		t.Fatalf("page=%+v queries=%+v", p, r.groupQueries)
	}
}

func TestGroupStateWorstChildIsOrderIndependent(t *testing.T) {
	below := true
	cases := []struct {
		name  string
		items []domain.ListingReadModel
		want  domain.GroupState
	}{
		{"ok", []domain.ListingReadModel{{Link: domain.ListingLink{State: domain.LinkStateResolved}}}, domain.GroupStateOK},
		{"stale", []domain.ListingReadModel{{Link: domain.ListingLink{State: domain.LinkStateResolved}, SyncState: domain.ListingSyncStateStale}}, domain.GroupStateAttention},
		{"below", []domain.ListingReadModel{{Link: domain.ListingLink{State: domain.LinkStateResolved}, BelowMarginWorstCase: &below}}, domain.GroupStateAttention},
		{"unlinked", []domain.ListingReadModel{{Link: domain.ListingLink{State: domain.LinkStateConflict}}}, domain.GroupStateAttention},
		{"error-last", []domain.ListingReadModel{{Link: domain.ListingLink{State: domain.LinkStateConflict}}, {Link: domain.ListingLink{State: domain.LinkStateResolved}, SyncState: domain.ListingSyncStateError}}, domain.GroupStateError},
		{"error-first", []domain.ListingReadModel{{Link: domain.ListingLink{State: domain.LinkStateResolved}, SyncError: &domain.ReadSyncError{}}, {Link: domain.ListingLink{State: domain.LinkStateConflict}}}, domain.GroupStateError},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := groupState(tc.items); got != tc.want {
				t.Fatalf("got=%s want=%s", got, tc.want)
			}
		})
	}
}

func TestSemProdutoDependentFilterDefinitions(t *testing.T) {
	item := domain.ListingReadModel{Link: domain.ListingLink{State: domain.LinkStateUnresolved}}
	yes, no := true, false
	if matchesDependentFilter(item, domain.ListingFilter{Exception: domain.ListingExceptionBelowMargin}) || !matchesDependentFilter(item, domain.ListingFilter{HasException: &yes}) || matchesDependentFilter(item, domain.ListingFilter{HasException: &no}) {
		t.Fatal("sem produto definition mismatch")
	}
}
func ptr(v float64) *float64     { return &v }
func ptrString(v string) *string { return &v }
