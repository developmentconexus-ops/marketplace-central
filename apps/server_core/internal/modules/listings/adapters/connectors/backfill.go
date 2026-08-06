package connectors

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	mercadolivre "marketplace-central/apps/server_core/internal/modules/connectors/adapters/mercado_livre"
	connectorsapp "marketplace-central/apps/server_core/internal/modules/connectors/application"
	connectorsdomain "marketplace-central/apps/server_core/internal/modules/connectors/domain"
	listingsdomain "marketplace-central/apps/server_core/internal/modules/listings/domain"
	listingsports "marketplace-central/apps/server_core/internal/modules/listings/ports"
)

// idScanReader is the optional ids-only enumeration capability (Passo 1's
// ListListingsScanIDs). Not every ports.ListingReader implements it — same
// optional-capability pattern as this package's scanPageReader (source.go):
// a type assertion against the concrete reader returned by
// MarketplaceCapabilityService.ListingReader, never a new method on the
// shared ports.ListingReader interface.
type idScanReader interface {
	ListListingsScanIDs(ctx context.Context, input connectorsdomain.ListListingsInput) ([]string, string, error)
}

// BackfillSource is F-03's IDEnumerator: it produces ids only, never
// hydrates (IC-06).
type BackfillSource struct {
	capabilities *connectorsapp.MarketplaceCapabilityService
}

func NewBackfillSource(capabilities *connectorsapp.MarketplaceCapabilityService) BackfillSource {
	return BackfillSource{capabilities: capabilities}
}

func (s BackfillSource) ReadIDPage(ctx context.Context, account listingsports.InstallationAccount, cursor string, limit int) (listingsports.BatchIDPage, error) {
	reader, err := s.capabilities.ListingReader(account.ProviderCode)
	if err != nil {
		return listingsports.BatchIDPage{}, err
	}
	scanner, ok := reader.(idScanReader)
	if !ok {
		return listingsports.BatchIDPage{}, fmt.Errorf("provider %q does not support ids-only backfill enumeration", account.ProviderCode)
	}
	ref := connectorsdomain.ProviderAccountRef{TenantID: account.TenantID, InstallationID: account.InstallationID, ProviderCode: account.ProviderCode, ProviderAccountID: account.ProviderAccountID}
	ids, next, err := scanner.ListListingsScanIDs(ctx, connectorsdomain.ListListingsInput{AccountRef: ref, Cursor: cursor, Limit: limit})
	if err != nil {
		return listingsports.BatchIDPage{}, err
	}
	return listingsports.BatchIDPage{IDs: ids, NextCursor: next}, nil
}

// multigetReader is the optional batch-hydration capability (M-01/F-02's
// GetItemsMultiget) — same optional-capability, type-assertion pattern as
// idScanReader above.
type multigetReader interface {
	GetItemsMultiget(ctx context.Context, accountRef connectorsdomain.ProviderAccountRef, providerItemIDs []string) ([]mercadolivre.ItemMultigetDTO, error)
}

// MultigetHydrator is F-03's BatchHydrator: it hydrates a batch of ids via
// multiget, never enumerates (IC-06).
//
// ADR-019 non-regression: BEFORE mapping to canonical listing rows, every
// successfully-hydrated item is ALSO converted to a
// connectorsdomain.ListingSnapshot (multiget_mapper.go's
// MapMultigetItemToListingSnapshot) and handed to the SAME SnapshotObserver
// the existing single-item ReadListing path feeds via Source.observe
// (source.go:40-45) — so the new backfill/sweep path keeps the product-links
// matcher's EAN/SKU anchors fed exactly like the old path does, per IC-07
// ("hidratação nova CONTINUA alimentando AbsorbProviderSnapshots"). Observer
// failure fails the whole HydrateBatch call, mirroring Source.observe's own
// fail-closed contract (a pull that silently skipped the linker would report
// a complete catalog while the matcher stayed blind to part of it).
type MultigetHydrator struct {
	capabilities *connectorsapp.MarketplaceCapabilityService
	observer     SnapshotObserver
	now          func() time.Time
}

// NewMultigetHydrator builds the hydrator. observer may be nil (no
// product-links import service wired) — the plain hydrator, matching
// NewSource/NewSourceWithObserver's own nil-is-plain convention.
func NewMultigetHydrator(capabilities *connectorsapp.MarketplaceCapabilityService, observer SnapshotObserver, now func() time.Time) MultigetHydrator {
	if now == nil {
		now = time.Now
	}
	return MultigetHydrator{capabilities: capabilities, observer: observer, now: now}
}

func (h MultigetHydrator) HydrateBatch(ctx context.Context, account listingsports.InstallationAccount, ids []string) ([]listingsdomain.Listing, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	reader, err := h.capabilities.ListingReader(account.ProviderCode)
	if err != nil {
		return nil, err
	}
	hydrator, ok := reader.(multigetReader)
	if !ok {
		return nil, fmt.Errorf("provider %q does not support multiget hydration", account.ProviderCode)
	}
	ref := connectorsdomain.ProviderAccountRef{TenantID: account.TenantID, InstallationID: account.InstallationID, ProviderCode: account.ProviderCode, ProviderAccountID: account.ProviderAccountID}
	items, err := hydrator.GetItemsMultiget(ctx, ref, ids)
	if err != nil {
		// Batch-level failure (e.g. 429 series beyond M-01's retry budget):
		// fails this whole tick, never attributed to a single item.
		return nil, err
	}

	fetchedAt := h.now().UTC()

	if h.observer != nil {
		snapshots := make([]connectorsdomain.ListingSnapshot, 0, len(items))
		for _, item := range items {
			if item.Err != nil {
				continue
			}
			snapshots = append(snapshots, MapMultigetItemToListingSnapshot(account.ProviderCode, item, fetchedAt))
		}
		if len(snapshots) > 0 {
			if err := h.observer.AbsorbProviderSnapshots(ctx, account.InstallationID, snapshots); err != nil {
				return nil, fmt.Errorf("absorb provider snapshots: %w", err)
			}
		}
	}

	rows, itemErrs := MapMultigetItemsToListings(account.TenantID, account.InstallationID, account.ProviderCode, items, fetchedAt)
	for _, itemErr := range itemErrs {
		// IC-06: a per-item failure is recorded and skipped, never aborts
		// the batch.
		slog.WarnContext(ctx, "listings backfill: item hydration/mapping failed",
			"provider_item_id", itemErr.ProviderItemID, "error", itemErr.Err, "provider", account.ProviderCode)
	}
	return rows, nil
}
