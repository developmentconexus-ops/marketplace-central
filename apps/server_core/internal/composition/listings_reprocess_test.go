package composition

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"marketplace-central/apps/server_core/internal/contexts/listings/contracts"
	"marketplace-central/apps/server_core/internal/kernel/channel"
	"marketplace-central/apps/server_core/internal/kernel/fact"
	"marketplace-central/apps/server_core/internal/kernel/provenance"
	"marketplace-central/apps/server_core/internal/kernel/tenant"
)

func reprocessKey(t *testing.T, listingID string) contracts.SourceListingKey {
	t.Helper()
	tid, err := tenant.Parse("tenant_default")
	if err != nil {
		t.Fatalf("tenant: %v", err)
	}
	code, err := channel.ParseCode("mercado_livre")
	if err != nil {
		t.Fatalf("channel: %v", err)
	}
	account, err := channel.NewAccountRef(code, "179571326")
	if err != nil {
		t.Fatalf("account: %v", err)
	}
	key, err := contracts.NewSourceListingKey(tid, account, listingID)
	if err != nil {
		t.Fatalf("key: %v", err)
	}
	return key
}

// fakeListings serves stored pages in the order given and records what was
// folded back in.
type fakeListings struct {
	pages     []contracts.StoredPage
	served    int
	ingested  []string
	result    map[string]contracts.Disposition
	ingestErr error
}

func (f *fakeListings) StoredObservations(_ context.Context, _ tenant.ID, _ contracts.StoredCursor, _ int) (contracts.StoredPage, error) {
	if f.served >= len(f.pages) {
		return contracts.StoredPage{Done: true}, nil
	}
	page := f.pages[f.served]
	f.served++
	return page, nil
}

func (f *fakeListings) IngestListing(_ context.Context, o contracts.ListingObservation) (contracts.IngestResult, error) {
	if f.ingestErr != nil {
		return contracts.IngestResult{}, f.ingestErr
	}
	f.ingested = append(f.ingested, o.Key.ListingID())
	d, ok := f.result[o.Key.ListingID()]
	if !ok {
		d = contracts.DispositionIdempotent
	}
	return contracts.IngestResult{Disposition: d, Version: 2}, nil
}

// fakeMapper turns a stored observation into an observation with one Known
// title, which is enough for the run to fold; mapErr fails a named listing.
type fakeMapper struct {
	mapErrFor string
}

func (m fakeMapper) MapStored(so contracts.StoredObservation) (contracts.ListingObservation, error) {
	if m.mapErrFor != "" && so.Key.ListingID() == m.mapErrFor {
		return contracts.ListingObservation{}, errors.New("payload does not hash to its stored hash")
	}
	evidence, err := provenance.NewEvidence("mercado_livre", "item", so.Key.ListingID(), so.ObservedAt, so.PayloadHash)
	if err != nil {
		return contracts.ListingObservation{}, err
	}
	title, err := fact.NewKnown("Produto "+so.Key.ListingID(), evidence)
	if err != nil {
		return contracts.ListingObservation{}, err
	}
	return contracts.ListingObservation{
		Key:        so.Key,
		State:      contracts.ListingState{Title: title},
		RawPayload: so.RawPayload,
		Evidence:   evidence,
	}, nil
}

func storedPage(t *testing.T, done bool, ids ...string) contracts.StoredPage {
	t.Helper()
	page := contracts.StoredPage{Done: done}
	for _, id := range ids {
		key := reprocessKey(t, id)
		page.Observations = append(page.Observations, contracts.StoredObservation{
			Key:         key,
			RawPayload:  []byte(`{"id":"` + id + `"}`),
			PayloadHash: "hash-" + id,
			ObservedAt:  time.Date(2026, 8, 10, 12, 0, 0, 0, time.UTC),
		})
		page.Next = contracts.StoredCursorAfter(key)
	}
	return page
}

// TestRunListingsReprocessFoldsEveryStoredPayload: the run reads what is
// stored, maps it through the channel adapter, folds it back, and reports each
// disposition separately — "nothing changed" and "nothing was read" are
// different outcomes, and only the counts distinguish them.
func TestRunListingsReprocessFoldsEveryStoredPayload(t *testing.T) {
	module := &fakeListings{
		pages: []contracts.StoredPage{
			storedPage(t, false, "MLB1", "MLB2"),
			storedPage(t, true, "MLB3"),
		},
		result: map[string]contracts.Disposition{
			"MLB1": contracts.DispositionChanged,
			"MLB2": contracts.DispositionChanged,
			"MLB3": contracts.DispositionIdempotent,
		},
	}
	tid, _ := tenant.Parse("tenant_default")

	report, err := RunListingsReprocess(context.Background(), module, fakeMapper{}, tid, 2)
	if err != nil {
		t.Fatalf("reprocess: %v", err)
	}
	if got := strings.Join(module.ingested, ","); got != "MLB1,MLB2,MLB3" {
		t.Fatalf("ingested %q, want every stored listing exactly once in order", got)
	}
	if report.Pages != 2 || report.Read != 3 || report.Changed != 2 || report.Idempotent != 1 || report.Created != 0 {
		t.Fatalf("report = %+v, want pages=2 read=3 changed=2 idempotent=1 created=0", report)
	}
}

// TestRunListingsReprocessFailsLoudOnAMappingError: a payload that cannot be
// mapped is a data-integrity failure, and skipping it would report a clean run
// over rows nobody reprocessed.
func TestRunListingsReprocessFailsLoudOnAMappingError(t *testing.T) {
	module := &fakeListings{pages: []contracts.StoredPage{storedPage(t, true, "MLB1", "MLB2")}}
	tid, _ := tenant.Parse("tenant_default")

	_, err := RunListingsReprocess(context.Background(), module, fakeMapper{mapErrFor: "MLB2"}, tid, 50)
	if err == nil || !strings.Contains(err.Error(), "MLB2") {
		t.Fatalf("a mapping failure must fail the run naming the listing, got: %v", err)
	}
}

func TestRunListingsReprocessFailsLoudOnAnIngestError(t *testing.T) {
	module := &fakeListings{
		pages:     []contracts.StoredPage{storedPage(t, true, "MLB1")},
		ingestErr: errors.New("save v2: connection refused"),
	}
	tid, _ := tenant.Parse("tenant_default")

	_, err := RunListingsReprocess(context.Background(), module, fakeMapper{}, tid, 50)
	if err == nil || !strings.Contains(err.Error(), "MLB1") {
		t.Fatalf("an ingest failure must fail the run naming the listing, got: %v", err)
	}
}

func TestRunListingsReprocessRejectsANonPositivePageSize(t *testing.T) {
	tid, _ := tenant.Parse("tenant_default")
	if _, err := RunListingsReprocess(context.Background(), &fakeListings{}, fakeMapper{}, tid, 0); err == nil {
		t.Fatal("page size 0 was accepted")
	}
}

// TestRunListingsReprocessRefusesAPageThatDoesNotAdvance: a page that says
// "not done" while handing back the start cursor would restart the walk from
// row one forever, reprocessing the same listings and never reaching the rest.
func TestRunListingsReprocessRefusesAPageThatDoesNotAdvance(t *testing.T) {
	page := storedPage(t, false, "MLB1")
	page.Next = contracts.StoredCursor{}
	module := &fakeListings{pages: []contracts.StoredPage{page}}
	tid, _ := tenant.Parse("tenant_default")

	_, err := RunListingsReprocess(context.Background(), module, fakeMapper{}, tid, 1)
	if err == nil || !strings.Contains(err.Error(), "cursor") {
		t.Fatalf("a non-advancing page must fail naming the cursor, got: %v", err)
	}
}
