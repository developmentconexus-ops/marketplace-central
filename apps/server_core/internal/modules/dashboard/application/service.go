package application

import (
	"context"
	"strings"
	"time"

	"marketplace-central/apps/server_core/internal/modules/dashboard/domain"
	"marketplace-central/apps/server_core/internal/modules/dashboard/ports"
	erpimportdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

var ErrInstallationNotFound = domain.ErrInstallationNotFound

type InstallationNotFoundError = domain.InstallationNotFoundError

type Service struct {
	installation ports.InstallationSource
	listings     ports.ListingsSource
	linkage      ports.LinkageSource
	orders       ports.OrdersSource
	sync         ports.SyncSource
	erp          ports.ErpImportSource
	now          func() time.Time
}

func NewService(
	installation ports.InstallationSource,
	listings ports.ListingsSource,
	linkage ports.LinkageSource,
	orders ports.OrdersSource,
	sync ports.SyncSource,
	erp ports.ErpImportSource,
	now func() time.Time,
) Service {
	return Service{
		installation: installation,
		listings:     listings,
		linkage:      linkage,
		orders:       orders,
		sync:         sync,
		erp:          erp,
		now:          now,
	}
}

func (s Service) Summary(ctx context.Context, installationID string) (domain.Summary, error) {
	installationID = strings.TrimSpace(installationID)
	if s.installation == nil {
		return domain.Summary{}, &domain.InstallationNotFoundError{InstallationID: installationID}
	}
	_, found, err := s.installation.Get(ctx, installationID)
	if err != nil || !found {
		return domain.Summary{}, &domain.InstallationNotFoundError{InstallationID: installationID}
	}

	asOf := time.Time{}
	if s.now != nil {
		asOf = s.now().UTC()
	}
	result := domain.Summary{
		LastSyncAt: make(map[string]*time.Time),
		AsOf:       asOf,
	}
	failed := make(map[string]bool, len(sourceOrder))

	if s.listings == nil {
		failed["listings"] = true
	} else if row, sourceErr := s.listings.Summary(ctx, listingsports.SummaryQuery{InstallationID: installationID}); sourceErr != nil {
		failed["listings"] = true
	} else {
		result.SyncErrors = int64Pointer(int64(row.SyncError))
		result.BelowMargin = intPointer(row.BelowMarginWorstCase)
		result.AnunciosAtivos = int64Pointer(int64(row.Active))
	}

	if s.linkage == nil {
		failed["linkage"] = true
	} else if summary, sourceErr := s.linkage.Summary(ctx, installationID); sourceErr != nil {
		failed["linkage"] = true
	} else {
		result.PendingLinks = int64Pointer(summary.PendingLinks)
		result.MissingGTIN = int64Pointer(summary.MissingGTIN)
	}

	if s.orders == nil {
		failed["orders"] = true
	} else if summary, sourceErr := s.orders.Summary(ctx, installationID, asOf); sourceErr != nil {
		failed["orders"] = true
	} else {
		result.OrdersToday = int64Pointer(summary.Today)
		result.Orders7d = int64Pointer(summary.SevenDays)
	}

	if s.sync == nil {
		failed["sync"] = true
	} else if runs, sourceErr := s.sync.LatestRunsByModule(ctx, installationID); sourceErr != nil {
		failed["sync"] = true
	} else {
		for _, run := range runs {
			result.LastSyncAt[run.OperationType] = cloneTime(run.LatestAttemptedAt)
		}
	}

	if s.erp == nil {
		failed["erp_import"] = true
	} else if reports, sourceErr := s.erp.ListImports(ctx); sourceErr != nil {
		failed["erp_import"] = true
	} else {
		result.LastImport = newestCompletedImport(reports, asOf)
	}

	result.Degraded = degradedSources(failed)
	if failed["sync"] {
		result.LastSyncAt = nil
	}
	return result, nil
}

func newestCompletedImport(reports []erpimportdomain.ImportReport, asOf time.Time) *domain.LastImport {
	var newest *erpimportdomain.ImportReport
	for i := range reports {
		report := reports[i]
		if report.Status != erpimportdomain.ImportStatusCompleted {
			continue
		}
		if newest == nil || report.ImportedAt.After(newest.ImportedAt) {
			newest = &report
		}
	}
	if newest == nil {
		return nil
	}
	return &domain.LastImport{
		At:         newest.ImportedAt.UTC(),
		AgeSeconds: int64(asOf.Sub(newest.ImportedAt).Seconds()),
	}
}

var sourceOrder = [...]string{"listings", "linkage", "orders", "sync", "erp_import"}

func degradedSources(failed map[string]bool) []string {
	degraded := make([]string, 0, len(sourceOrder))
	seen := make(map[string]bool, len(sourceOrder))
	for _, source := range sourceOrder {
		if failed[source] && !seen[source] {
			degraded = append(degraded, source)
			seen[source] = true
		}
	}
	return degraded
}

func int64Pointer(value int64) *int64 { return &value }

func intPointer(value *int) *int64 {
	if value == nil {
		return nil
	}
	converted := int64(*value)
	return &converted
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := value.UTC()
	return &cloned
}
