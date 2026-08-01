package application

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

var ErrDuplicateCanonicalKey = errors.New("duplicate canonical listing key")

// maxIngestionPages is an honest ceiling on a single completed pull. A connector
// that never signals a short final page (bug or malice) would otherwise loop and
// grow memory without bound; past this many full pages we fail loudly instead of
// truncating silently. A real catalog that legitimately exceeds it must page in
// smaller runs — a wrong count is never treated as complete.
const maxIngestionPages = 10_000

type Ingestion struct {
	source   ports.PageSource
	store    ports.CompletedPullStore
	pageSize int
	now      func() time.Time
}

func NewIngestion(source ports.PageSource, store ports.CompletedPullStore, pageSize int, now func() time.Time) *Ingestion {
	return &Ingestion{source: source, store: store, pageSize: pageSize, now: now}
}

func (i *Ingestion) Pull(ctx context.Context, account ports.InstallationAccount) error {
	if i.source == nil || i.store == nil || i.pageSize <= 0 || i.now == nil {
		return errors.New("listing ingestion is not configured")
	}

	// Captured at the START of the pull, not at the end: it is this run's
	// reference time (ports.CompletedPullStore), used to bound "not seen
	// since" in the keep-absent step. Pull today only calls
	// ApplyCompletedPull once, after its full page loop drains — i.e. every
	// call this method makes already corresponds to a COMPLETE run (F-03
	// will own the resumable/partial case via a different caller).
	runStarted := i.now().UTC()

	var rows []domain.Listing
	seen := make(map[domain.ListingKey]struct{})
	offset := 0
	cursor := ""
	for pages := 0; ; pages++ {
		if pages >= maxIngestionPages {
			return fmt.Errorf("listing ingestion exceeded %d pages without a final short page", maxIngestionPages)
		}
		page, err := i.source.ReadPage(ctx, account, cursor, i.pageSize)
		if err != nil {
			return err
		}
		if page.ProviderItemCount < 0 || page.ProviderItemCount > i.pageSize {
			return fmt.Errorf("invalid provider page item count %d", page.ProviderItemCount)
		}
		for _, row := range page.Rows {
			if _, duplicate := seen[row.Key]; duplicate {
				return fmt.Errorf("%w: %s/%s", ErrDuplicateCanonicalKey, row.Key.ProviderListingID, row.Key.VariationID)
			}
			seen[row.Key] = struct{}{}
			rows = append(rows, row)
		}
		offset += page.ProviderItemCount
		if page.ProviderItemCount < i.pageSize {
			break
		}
		// Provider-issued cursor wins (ML scan scroll_id — offset paging is
		// hard-capped at 1000 items by the provider); numeric offset is the
		// fallback for offset-paged sources.
		if page.NextCursor != "" {
			cursor = page.NextCursor
		} else {
			cursor = strconv.Itoa(offset)
		}
	}

	return i.store.ApplyCompletedPull(ctx, account.InstallationID, rows, runStarted, true)
}
