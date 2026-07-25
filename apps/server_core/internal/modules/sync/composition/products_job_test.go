package composition_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
	synccomposition "marketplace-central/apps/server_core/internal/modules/sync/composition"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

const jobTenant = "11111111-1111-4111-8111-111111111111"

type stubLookup struct {
	source tenant_config.ActiveSource
	err    error
}

func (s stubLookup) Get(ctx context.Context, tenantID string) (tenant_config.Config, error) {
	if s.err != nil {
		return tenant_config.Config{}, s.err
	}
	return tenant_config.Config{TenantID: tenantID, Source: s.source}, nil
}

type stubAdapter struct {
	readports.Reader
	calls    int
	result   readports.SyncResult
	err      error
	lastCtx  context.Context
	kindWant sourcekind.SourceKind
}

func (a *stubAdapter) Sync(ctx context.Context) (readports.SyncResult, error) {
	a.calls++
	a.lastCtx = ctx
	return a.result, a.err
}

func (a *stubAdapter) Kind() sourcekind.SourceKind { return a.kindWant }

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

// The job runs the adapter of the CONFIGURED source, not a hardcoded one, and
// leaves every other source's adapter untouched.
func TestProductsJobRunsTheActiveSourceAdapter(t *testing.T) {
	active := &stubAdapter{result: readports.SyncResult{Processed: 1845}}
	other := &stubAdapter{}
	at := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)

	job := synccomposition.NewProductsJob(jobTenant,
		stubLookup{source: tenant_config.SourceXLSX},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{
			tenant_config.SourceXLSX:    active,
			tenant_config.SourceSankhya: other,
		}, fixedClock(at))

	next, err := job(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if active.calls != 1 {
		t.Fatalf("active adapter calls=%d want 1", active.calls)
	}
	if other.calls != 0 {
		t.Fatalf("inactive adapter calls=%d want 0", other.calls)
	}

	var cursor synccomposition.ProductsCursor
	if err := json.Unmarshal(next, &cursor); err != nil {
		t.Fatal(err)
	}
	if cursor.Source != string(tenant_config.SourceXLSX) {
		t.Fatalf("cursor source=%q want xlsx", cursor.Source)
	}
	if cursor.Processed != 1845 {
		t.Fatalf("cursor processed=%d want 1845", cursor.Processed)
	}
	if !cursor.CompletedAt.Equal(at) {
		t.Fatalf("cursor completed_at=%v want %v", cursor.CompletedAt, at)
	}
}

// Flipping the active source flips which adapter runs, with no re-wiring.
func TestProductsJobFollowsTheActiveSourceFlip(t *testing.T) {
	upload := &stubAdapter{}
	live := &stubAdapter{}
	adapters := map[tenant_config.ActiveSource]readports.ProductSourceAdapter{
		tenant_config.SourceXLSX:    upload,
		tenant_config.SourceSankhya: live,
	}

	uploadJob := synccomposition.NewProductsJob(jobTenant, stubLookup{source: tenant_config.SourceXLSX}, adapters, nil)
	if _, err := uploadJob(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	liveJob := synccomposition.NewProductsJob(jobTenant, stubLookup{source: tenant_config.SourceSankhya}, adapters, nil)
	if _, err := liveJob(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if upload.calls != 1 || live.calls != 1 {
		t.Fatalf("upload=%d live=%d want 1 each", upload.calls, live.calls)
	}
}

// No configured source is a failure, never a default: syncing an assumed source
// would silently overwrite the mirror from the wrong ERP.
func TestProductsJobFailsClosedWhenTheActiveSourceIsUnreadable(t *testing.T) {
	adapter := &stubAdapter{}
	job := synccomposition.NewProductsJob(jobTenant,
		stubLookup{err: errors.New("no active source configured")},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{tenant_config.SourceXLSX: adapter}, nil)

	previous := json.RawMessage(`{"source":"xlsx","processed":3}`)
	next, err := job(context.Background(), previous)
	if err == nil {
		t.Fatal("want error when the active source cannot be resolved")
	}
	if adapter.calls != 0 {
		t.Fatalf("adapter ran %d times despite an unresolved source", adapter.calls)
	}
	if string(next) != string(previous) {
		t.Fatalf("cursor=%s want the previous cursor preserved", next)
	}
}

// A source with no registered adapter (e.g. Sankhya with no ERP connection) is
// named in the error instead of being silently skipped as a success.
func TestProductsJobFailsClosedWhenNoAdapterIsRegistered(t *testing.T) {
	job := synccomposition.NewProductsJob(jobTenant,
		stubLookup{source: tenant_config.SourceSankhya},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{
			tenant_config.SourceXLSX: &stubAdapter{},
		}, nil)

	if _, err := job(context.Background(), nil); err == nil {
		t.Fatal("want error when the active source has no adapter")
	} else if !strings.Contains(err.Error(), "sankhya") {
		t.Fatalf("error %q does not name the unserved source", err)
	}
}

// A nil entry is as absent as a missing key — a half-wired map must not panic
// the scheduler loop.
func TestProductsJobFailsClosedOnANilAdapter(t *testing.T) {
	job := synccomposition.NewProductsJob(jobTenant,
		stubLookup{source: tenant_config.SourceXLSX},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{tenant_config.SourceXLSX: nil}, nil)

	if _, err := job(context.Background(), nil); err == nil {
		t.Fatal("want error for a nil adapter")
	}
}

type stubRefresher struct {
	calls int
	err   error
}

func (r *stubRefresher) RefreshLinkCandidates(ctx context.Context) error {
	r.calls++
	return r.err
}

// M05-C11 (sync half): the ERP sync that brought new product data in is what
// asks "which product is this anúncio?" again — the operator does not have to
// press a button for a codprod that arrived minutes ago to become a candidate.
func TestProductsJobRefreshesLinkCandidatesAfterASuccessfulSync(t *testing.T) {
	refresher := &stubRefresher{}
	job := synccomposition.NewProductsJob(jobTenant,
		stubLookup{source: tenant_config.SourceSankhya},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{
			tenant_config.SourceSankhya: &stubAdapter{result: readports.SyncResult{Processed: 10529}},
		}, nil, synccomposition.WithLinkCandidateRefresh(refresher))

	if _, err := job(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if refresher.calls != 1 {
		t.Fatalf("refresh calls=%d want 1", refresher.calls)
	}
}

// M05-C11: the products landed. A matcher that fails afterwards is logged and
// retried next cycle — reporting the cycle as failed would claim the product
// data did not arrive, which is false, and would erase the cursor with it.
func TestProductsJobSurvivesALinkCandidateRefreshFailure(t *testing.T) {
	refresher := &stubRefresher{err: errors.New("candidate generation exploded")}
	at := time.Date(2026, 7, 25, 9, 0, 0, 0, time.UTC)
	job := synccomposition.NewProductsJob(jobTenant,
		stubLookup{source: tenant_config.SourceSankhya},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{
			tenant_config.SourceSankhya: &stubAdapter{result: readports.SyncResult{Processed: 42}},
		}, fixedClock(at), synccomposition.WithLinkCandidateRefresh(refresher))

	next, err := job(context.Background(), nil)
	if err != nil {
		t.Fatalf("refresh failure must not fail the sync cycle: %v", err)
	}
	var cursor synccomposition.ProductsCursor
	if err := json.Unmarshal(next, &cursor); err != nil {
		t.Fatal(err)
	}
	if cursor.Processed != 42 {
		t.Fatalf("cursor processed=%d want the sync that actually landed (42)", cursor.Processed)
	}
}

// Nothing new arrived, so there is nothing new to match: a failed sync must not
// trigger a regeneration over stale data.
func TestProductsJobDoesNotRefreshWhenTheSyncFailed(t *testing.T) {
	refresher := &stubRefresher{}
	job := synccomposition.NewProductsJob(jobTenant,
		stubLookup{source: tenant_config.SourceSankhya},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{
			tenant_config.SourceSankhya: &stubAdapter{err: errors.New("oracle unreachable")},
		}, nil, synccomposition.WithLinkCandidateRefresh(refresher))

	if _, err := job(context.Background(), nil); err == nil {
		t.Fatal("want the adapter error surfaced")
	}
	if refresher.calls != 0 {
		t.Fatalf("refresh calls=%d want 0 after a failed sync", refresher.calls)
	}
}

// An adapter failure keeps the stored cursor: the scheduler records last_error,
// and the next cycle must resume from real progress, not from an erased cursor.
func TestProductsJobKeepsTheCursorWhenTheAdapterFails(t *testing.T) {
	job := synccomposition.NewProductsJob(jobTenant,
		stubLookup{source: tenant_config.SourceXLSX},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{
			tenant_config.SourceXLSX: &stubAdapter{err: errors.New("oracle unreachable")},
		}, nil)

	previous := json.RawMessage(`{"source":"xlsx","processed":7}`)
	next, err := job(context.Background(), previous)
	if err == nil {
		t.Fatal("want the adapter error surfaced")
	}
	if string(next) != string(previous) {
		t.Fatalf("cursor=%s want %s", next, previous)
	}
}
