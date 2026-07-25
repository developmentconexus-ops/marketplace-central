package application

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	dashboarddomain "marketplace-central/apps/server_core/internal/modules/dashboard/domain"
	erpimportdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	integrationports "marketplace-central/apps/server_core/internal/modules/integrations/ports"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
	linkageports "marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

var summaryReferenceTime = time.Date(2026, 7, 16, 15, 4, 5, 0, time.UTC)

type installationSourceStub struct {
	found bool
	err   error
}

func (s installationSourceStub) Get(context.Context, string) (domain.Installation, bool, error) {
	return domain.Installation{}, s.found, s.err
}

type listingsSourceStub struct {
	row listingsports.ListingSummaryRow
	err error
}

func (s listingsSourceStub) Summary(context.Context, listingsports.SummaryQuery) (listingsports.ListingSummaryRow, error) {
	return s.row, s.err
}

type linkageSourceStub struct {
	summary linkageports.LinkageSummary
	err     error
}

func (s linkageSourceStub) Summary(context.Context, string) (linkageports.LinkageSummary, error) {
	return s.summary, s.err
}

type ordersSourceStub struct {
	summary ordersports.OrderSummary
	err     error
}

func (s ordersSourceStub) Summary(context.Context, string, time.Time) (ordersports.OrderSummary, error) {
	return s.summary, s.err
}

type syncSourceStub struct {
	runs []integrationports.LatestRunByModule
	err  error
}

func (s syncSourceStub) LatestRunsByModule(context.Context, string) ([]integrationports.LatestRunByModule, error) {
	return s.runs, s.err
}

type erpSourceStub struct {
	reports []erpimportdomain.ImportReport
	err     error
}

func (s erpSourceStub) ListImports(context.Context) ([]erpimportdomain.ImportReport, error) {
	return s.reports, s.err
}

func healthySources() (installationSourceStub, listingsSourceStub, linkageSourceStub, ordersSourceStub, syncSourceStub, erpSourceStub) {
	margin := 7
	lastSync := summaryReferenceTime.Add(-time.Hour)
	return installationSourceStub{found: true},
		listingsSourceStub{row: listingsports.ListingSummaryRow{SyncError: 2, Active: 42, Unlinked: 9, BelowMarginWorstCase: &margin}},
		linkageSourceStub{summary: linkageports.LinkageSummary{PendingLinks: 3, MissingGTIN: 4}},
		ordersSourceStub{summary: ordersports.OrderSummary{Today: 5, SevenDays: 11}},
		syncSourceStub{runs: []integrationports.LatestRunByModule{{OperationType: "listings_refresh", LatestAttemptedAt: &lastSync}}},
		erpSourceStub{reports: []erpimportdomain.ImportReport{
			{Protocol: "p-1", ImportedAt: summaryReferenceTime.Add(-2 * time.Hour), Status: erpimportdomain.ImportStatusCompleted},
		}}
}

func newTestService(
	installation installationSourceStub,
	listings listingsSourceStub,
	linkage linkageSourceStub,
	orders ordersSourceStub,
	sync syncSourceStub,
	erp erpSourceStub,
) Service {
	return NewService(installation, listings, linkage, orders, sync, erp, func() time.Time { return summaryReferenceTime })
}

func TestSummaryOneFailedLinkageSourceReturnsMissingGTINNull(t *testing.T) {
	installation, listings, _, orders, sync, erp := healthySources()
	service := newTestService(installation, listings, linkageSourceStub{err: errors.New("linkage unavailable")}, orders, sync, erp)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	// PendingLinks comes from the LISTINGS source (the /anuncios "sem vínculo"
	// predicate), so a down linkage source must not null it.
	if got.PendingLinks == nil || *got.PendingLinks != 9 {
		t.Fatalf("PendingLinks = %v, want 9", got.PendingLinks)
	}
	if got.MissingGTIN != nil {
		t.Fatalf("MissingGTIN = %v, want nil", *got.MissingGTIN)
	}
	if !reflect.DeepEqual(got.Degraded, []string{"linkage"}) {
		t.Fatalf("Degraded = %v, want [linkage]", got.Degraded)
	}
	if got.SyncErrors == nil || *got.SyncErrors != 2 {
		t.Fatalf("SyncErrors = %v, want 2", got.SyncErrors)
	}
	if got.BelowMargin == nil || *got.BelowMargin != 7 {
		t.Fatalf("BelowMargin = %v, want 7", got.BelowMargin)
	}
	if got.OrdersToday == nil || *got.OrdersToday != 5 {
		t.Fatalf("OrdersToday = %v, want 5", got.OrdersToday)
	}
	if got.Orders7d == nil || *got.Orders7d != 11 {
		t.Fatalf("Orders7d = %v, want 11", got.Orders7d)
	}
	lastSync, ok := got.LastSyncAt["listings_refresh"]
	if !ok || lastSync == nil {
		t.Fatalf("LastSyncAt[listings_refresh] = %v, want non-nil", lastSync)
	}
}

func TestSummaryAllSourcesFailReturnsAllNullAndFullDegradedList(t *testing.T) {
	service := newTestService(
		installationSourceStub{found: true},
		listingsSourceStub{err: errors.New("listings unavailable")},
		linkageSourceStub{err: errors.New("linkage unavailable")},
		ordersSourceStub{err: errors.New("orders unavailable")},
		syncSourceStub{err: errors.New("sync unavailable")},
		erpSourceStub{err: errors.New("erp unavailable")},
	)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.SyncErrors != nil || got.PendingLinks != nil || got.BelowMargin != nil || got.MissingGTIN != nil || got.OrdersToday != nil || got.Orders7d != nil || got.AnunciosAtivos != nil {
		t.Fatalf("summary counters = %+v, want all nil", got)
	}
	if got.LastSyncAt != nil {
		t.Fatalf("LastSyncAt = %v, want nil", got.LastSyncAt)
	}
	if got.LastImport != nil {
		t.Fatalf("LastImport = %v, want nil", got.LastImport)
	}
	if !reflect.DeepEqual(got.Degraded, []string{"listings", "linkage", "orders", "sync", "erp_import"}) {
		t.Fatalf("Degraded = %v, want full fixed order", got.Degraded)
	}
}

func TestSummaryUnknownInstallationReturnsInstallationNotFound(t *testing.T) {
	service := newTestService(installationSourceStub{}, listingsSourceStub{}, linkageSourceStub{}, ordersSourceStub{}, syncSourceStub{}, erpSourceStub{})

	_, err := service.Summary(context.Background(), "missing")
	if !errors.Is(err, dashboarddomain.ErrInstallationNotFound) {
		t.Fatalf("Summary() error = %v, want installation_not_found", err)
	}
	var typed *dashboarddomain.InstallationNotFoundError
	if !errors.As(err, &typed) {
		t.Fatalf("Summary() error type = %T, want InstallationNotFoundError", err)
	}
}

func TestSummaryKeepsNilBelowMarginFromHealthySourceNil(t *testing.T) {
	installation, listings, linkage, orders, sync, erp := healthySources()
	listings.row.BelowMarginWorstCase = nil
	service := newTestService(installation, listings, linkage, orders, sync, erp)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.BelowMargin != nil {
		t.Fatalf("BelowMargin = %v, want nil", *got.BelowMargin)
	}
	if got.SyncErrors == nil || *got.SyncErrors != 2 {
		t.Fatalf("SyncErrors = %v, want 2", got.SyncErrors)
	}
}

func TestSummaryDegradedDeterministicAndDeduplicated(t *testing.T) {
	service := newTestService(
		installationSourceStub{found: true},
		listingsSourceStub{err: errors.New("listings unavailable")},
		linkageSourceStub{err: errors.New("linkage unavailable")},
		ordersSourceStub{err: errors.New("orders unavailable")},
		syncSourceStub{err: errors.New("sync unavailable")},
		erpSourceStub{err: errors.New("erp unavailable")},
	)

	for i := 0; i < 3; i++ {
		got, err := service.Summary(context.Background(), "installation-1")
		if err != nil {
			t.Fatalf("Summary() error = %v", err)
		}
		if !reflect.DeepEqual(got.Degraded, []string{"listings", "linkage", "orders", "sync", "erp_import"}) {
			t.Fatalf("run %d Degraded = %v, want fixed deduplicated order", i, got.Degraded)
		}
	}
}

func TestSummaryIncludesLastImportFromNewestCompletedImport(t *testing.T) {
	installation, listings, linkage, orders, sync, _ := healthySources()
	older := summaryReferenceTime.Add(-72 * time.Hour)
	newestCompleted := summaryReferenceTime.Add(-30 * time.Minute)
	newestRejected := summaryReferenceTime.Add(-5 * time.Minute)
	erp := erpSourceStub{reports: []erpimportdomain.ImportReport{
		{Protocol: "p-older", ImportedAt: older, Status: erpimportdomain.ImportStatusCompleted},
		{Protocol: "p-newest-completed", ImportedAt: newestCompleted, Status: erpimportdomain.ImportStatusCompleted},
		{Protocol: "p-newest-rejected", ImportedAt: newestRejected, Status: erpimportdomain.ImportStatusRejected},
	}}
	service := newTestService(installation, listings, linkage, orders, sync, erp)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.LastImport == nil {
		t.Fatalf("LastImport = nil, want non-nil")
	}
	if !got.LastImport.At.Equal(newestCompleted.UTC()) {
		t.Fatalf("LastImport.At = %v, want %v", got.LastImport.At, newestCompleted.UTC())
	}
	wantAge := int64(summaryReferenceTime.Sub(newestCompleted).Seconds())
	if got.LastImport.AgeSeconds != wantAge {
		t.Fatalf("LastImport.AgeSeconds = %v, want %v", got.LastImport.AgeSeconds, wantAge)
	}
	if len(got.Degraded) != 0 {
		t.Fatalf("Degraded = %v, want empty", got.Degraded)
	}
}

func TestSummaryNoImportsReturnsNilLastImportNotDegraded(t *testing.T) {
	installation, listings, linkage, orders, sync, _ := healthySources()
	erp := erpSourceStub{reports: []erpimportdomain.ImportReport{}}
	service := newTestService(installation, listings, linkage, orders, sync, erp)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.LastImport != nil {
		t.Fatalf("LastImport = %v, want nil", got.LastImport)
	}
	for _, source := range got.Degraded {
		if source == "erp_import" {
			t.Fatalf("Degraded = %v, want erp_import absent (no-import-yet is not a failure)", got.Degraded)
		}
	}
}

func TestSummaryErpSourceErrorDegradesErpImport(t *testing.T) {
	installation, listings, linkage, orders, sync, _ := healthySources()
	erp := erpSourceStub{err: errors.New("erp unavailable")}
	service := newTestService(installation, listings, linkage, orders, sync, erp)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.LastImport != nil {
		t.Fatalf("LastImport = %v, want nil", got.LastImport)
	}
	found := false
	for _, source := range got.Degraded {
		if source == "erp_import" {
			found = true
		}
	}
	if !found {
		t.Fatalf("Degraded = %v, want erp_import present", got.Degraded)
	}
}

func TestSummaryAnunciosAtivosMapsFromListingsActiveAndNullsOnListingsDown(t *testing.T) {
	installation, listings, linkage, orders, sync, erp := healthySources()
	service := newTestService(installation, listings, linkage, orders, sync, erp)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.AnunciosAtivos == nil || *got.AnunciosAtivos != int64(listings.row.Active) {
		t.Fatalf("AnunciosAtivos = %v, want %v", got.AnunciosAtivos, listings.row.Active)
	}

	downListings := listingsSourceStub{err: errors.New("listings unavailable")}
	serviceDown := newTestService(installation, downListings, linkage, orders, sync, erp)
	gotDown, err := serviceDown.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if gotDown.AnunciosAtivos != nil {
		t.Fatalf("AnunciosAtivos = %v, want nil when listings degraded", gotDown.AnunciosAtivos)
	}
}

func TestSummaryPendingLinksCountsUnlinkedListingsAndNullsOnListingsDown(t *testing.T) {
	installation, listings, linkage, orders, sync, erp := healthySources()
	// The linkage source reports a DIFFERENT (smaller) number on purpose: the
	// dashboard must answer with the listings predicate, not the workflow rows.
	service := newTestService(installation, listings, linkage, orders, sync, erp)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.PendingLinks == nil || *got.PendingLinks != int64(listings.row.Unlinked) {
		t.Fatalf("PendingLinks = %v, want %v", got.PendingLinks, listings.row.Unlinked)
	}

	downListings := listingsSourceStub{err: errors.New("listings unavailable")}
	serviceDown := newTestService(installation, downListings, linkage, orders, sync, erp)
	gotDown, err := serviceDown.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if gotDown.PendingLinks != nil {
		t.Fatalf("PendingLinks = %v, want nil when listings degraded", gotDown.PendingLinks)
	}
}
