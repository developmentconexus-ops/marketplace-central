package application

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	dashboarddomain "marketplace-central/apps/server_core/internal/modules/dashboard/domain"
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

func healthySources() (installationSourceStub, listingsSourceStub, linkageSourceStub, ordersSourceStub, syncSourceStub) {
	margin := 7
	lastSync := summaryReferenceTime.Add(-time.Hour)
	return installationSourceStub{found: true},
		listingsSourceStub{row: listingsports.ListingSummaryRow{SyncError: 2, BelowMarginWorstCase: &margin}},
		linkageSourceStub{summary: linkageports.LinkageSummary{PendingLinks: 3, MissingGTIN: 4}},
		ordersSourceStub{summary: ordersports.OrderSummary{Today: 5, SevenDays: 11}},
		syncSourceStub{runs: []integrationports.LatestRunByModule{{OperationType: "listings_refresh", LatestAttemptedAt: &lastSync}}}
}

func newTestService(
	installation installationSourceStub,
	listings listingsSourceStub,
	linkage linkageSourceStub,
	orders ordersSourceStub,
	sync syncSourceStub,
) Service {
	return NewService(installation, listings, linkage, orders, sync, func() time.Time { return summaryReferenceTime })
}

func TestSummaryOneFailedLinkageSourceReturnsPendingLinksNull(t *testing.T) {
	installation, listings, _, orders, sync := healthySources()
	service := newTestService(installation, listings, linkageSourceStub{err: errors.New("linkage unavailable")}, orders, sync)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.PendingLinks != nil {
		t.Fatalf("PendingLinks = %v, want nil", *got.PendingLinks)
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
	)

	got, err := service.Summary(context.Background(), "installation-1")
	if err != nil {
		t.Fatalf("Summary() error = %v", err)
	}
	if got.SyncErrors != nil || got.PendingLinks != nil || got.BelowMargin != nil || got.MissingGTIN != nil || got.OrdersToday != nil || got.Orders7d != nil {
		t.Fatalf("summary counters = %+v, want all nil", got)
	}
	if got.LastSyncAt != nil {
		t.Fatalf("LastSyncAt = %v, want nil", got.LastSyncAt)
	}
	if !reflect.DeepEqual(got.Degraded, []string{"listings", "linkage", "orders", "sync"}) {
		t.Fatalf("Degraded = %v, want full fixed order", got.Degraded)
	}
}

func TestSummaryUnknownInstallationReturnsInstallationNotFound(t *testing.T) {
	service := newTestService(installationSourceStub{}, listingsSourceStub{}, linkageSourceStub{}, ordersSourceStub{}, syncSourceStub{})

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
	installation, listings, linkage, orders, sync := healthySources()
	listings.row.BelowMarginWorstCase = nil
	service := newTestService(installation, listings, linkage, orders, sync)

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
	)

	for i := 0; i < 3; i++ {
		got, err := service.Summary(context.Background(), "installation-1")
		if err != nil {
			t.Fatalf("Summary() error = %v", err)
		}
		if !reflect.DeepEqual(got.Degraded, []string{"listings", "linkage", "orders", "sync"}) {
			t.Fatalf("run %d Degraded = %v, want fixed deduplicated order", i, got.Degraded)
		}
	}
}
