//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	listingspostgres "marketplace-central/apps/server_core/internal/modules/listings/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// TestRepositoryCompletedPullIsAtomicAndTenantScoped covers what survives
// from the pre-F-02 test of the same name: atomicity (a mid-batch CHECK
// violation rolls back every effect of that call) and tenant scoping (a
// sibling tenant's row is untouched). The MASS-CLOSURE assertions that used
// to live here (every omitted/absent row flips to status='closed', an empty
// pull closes the whole installation) asserted the exact defect F-02 kills
// (ADR-06, audit D-120 F1) and are gone — see
// TestApplyCompletedPull_AbortAfterPage1_NeverClosesUnseenRows and friends
// below for the behavior that replaced them.
func TestRepositoryCompletedPullIsAtomicAndTenantScoped(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings")
	token := time.Now().UTC().Format("150405.000000000")
	tenantA, tenantB, installation := "listing-a-"+token, "listing-b-"+token, "installation-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_sync_events WHERE tenant_id = ANY($1) AND installation_id = $2`, []string{tenantA, tenantB}, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id = ANY($1) AND installation_id = $2`, []string{tenantA, tenantB}, installation)
	})

	created := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	seedListing(t, pool, tenantA, installation, "same", "-", "old", "active", created, created)
	seedListing(t, pool, tenantA, installation, "omitted", "-", "old omitted", "active", created, created)
	seedListing(t, pool, tenantB, installation, "same", "-", "tenant b", "active", created, created)
	beforeB := readRawListing(t, pool, tenantB, installation, "same", "-")

	runStarted := created.Add(24 * time.Hour)
	fetched := runStarted.Add(-time.Minute)
	price := domain.PriceAmount("123.45")
	currency := domain.PriceCurrency("BRL")
	quantity := 7
	rows := []domain.Listing{
		listing(t, tenantA, installation, "same", "-", domain.ListingStatusActive, "updated", &price, &currency, &quantity, fetched, 80),
		listing(t, tenantA, installation, "same", "variation-2", domain.ListingStatusPaused, "variation", nil, nil, nil, fetched, 70),
	}
	repo := listingspostgres.NewRepository(pool, tenantA)
	if err := repo.UpsertPulledRows(ctx, installation, rows, runStarted); err != nil {
		t.Fatalf("UpsertPulledRows(): %v", err)
	}
	if err := repo.MarkRunComplete(ctx, installation, runStarted); err != nil {
		t.Fatalf("MarkRunComplete(): %v", err)
	}

	assertListing(t, pool, tenantA, installation, "same", "-", "active", "updated", created, runStarted, fetched, "123.45", "BRL", 7)
	assertListing(t, pool, tenantA, installation, "same", "variation-2", "paused", "variation", fetched, runStarted, fetched, "", "", -1)
	// "omitted" was never in rows and this run never advanced its
	// last_seen_at (it has none — seedListing never set one), so it is
	// untouched: still active, never closed, never even marked absent.
	assertStatus(t, pool, tenantA, installation, "omitted", "-", "active")
	afterB := readRawListing(t, pool, tenantB, installation, "same", "-")
	if !reflect.DeepEqual(beforeB, afterB) {
		t.Fatalf("tenant B changed: before=%#v after=%#v", beforeB, afterB)
	}
	assertEventKinds(t, pool, tenantA, installation, map[string]string{"same/-": "synced", "same/variation-2": "paused"})

	beforeRollback := readRawListing(t, pool, tenantA, installation, "same", "-")
	eventsBefore := eventCount(t, pool, tenantA, installation)
	invalid := listing(t, tenantA, installation, "failure", "-", domain.ListingStatusActive, "invalid", nil, nil, nil, fetched.Add(time.Hour), 101)
	if err := repo.UpsertPulledRows(ctx, installation, []domain.Listing{rows[0], invalid}, runStarted.Add(time.Hour)); err == nil {
		t.Fatal("UpsertPulledRows() error = nil, want check failure")
	}
	if got := readRawListing(t, pool, tenantA, installation, "same", "-"); !reflect.DeepEqual(beforeRollback, got) {
		t.Fatalf("listing effects did not roll back: before=%#v after=%#v", beforeRollback, got)
	}
	if got := eventCount(t, pool, tenantA, installation); got != eventsBefore {
		t.Fatalf("event count after rollback = %d, want %d", got, eventsBefore)
	}

	// An empty-rows tick followed by MarkRunComplete used to close the
	// entire installation (MASS-CLOSURE, F-02); now it is a no-op — the
	// run-scoped touched-check (F-03) finds no row with
	// last_seen_at = runStarted.Add(2h) (UpsertPulledRows was never called
	// with that timestamp), so MarkRunComplete skips the keep-absent UPDATE
	// entirely.
	if err := repo.UpsertPulledRows(ctx, installation, nil, runStarted.Add(2*time.Hour)); err != nil {
		t.Fatalf("empty UpsertPulledRows(): %v", err)
	}
	if err := repo.MarkRunComplete(ctx, installation, runStarted.Add(2*time.Hour)); err != nil {
		t.Fatalf("MarkRunComplete() after empty tick: %v", err)
	}
	assertStatus(t, pool, tenantA, installation, "same", "-", "active")
	assertStatus(t, pool, tenantA, installation, "same", "variation-2", "paused")
	assertStatus(t, pool, tenantA, installation, "omitted", "-", "active")
	if got := readRawListing(t, pool, tenantB, installation, "same", "-"); !reflect.DeepEqual(beforeB, got) {
		t.Fatalf("empty pull changed tenant B: before=%#v after=%#v", beforeB, got)
	}
}

// TestApplyCompletedPull_AbortAfterPage1_NeverClosesUnseenRows is the must-fail
// central to this feature (R-B / audit D-120 F1), now covering BOTH halves in
// one run: (phase A) a run that only ever saw page 1 of a paginated scan,
// then aborted (429/deadline/kill) before paging further, must never flip
// status nor mark absent_since on any row — absence requires MarkRunComplete,
// which an aborted run structurally never calls (F-03 split: no bool flag to
// forget, the call itself is the signal). (phase B) the SAME run resumed —
// same runStartedAt, remaining page ticked in via UpsertPulledRows, then
// MarkRunComplete — marks absent_since ONLY for rows genuinely never seen in
// ANY tick of the run; rows seen in either tick are exempt, and status is
// never touched by absence either way.
//
// Verified against the pre-F-02 code: the old ApplyCompletedPull opened with
// an unconditional `UPDATE listings SET status = 'closed' WHERE tenant_id =
// $1 AND installation_id = $2` (repository.go:390-394 before that change) —
// it has no `rows`/`complete` awareness at all, so it would have closed all
// 30 seeded rows before this test ever got to its assertions. This test would
// have failed on the old code; it is the red half of this must-fail's
// red/green pair.
func TestApplyCompletedPull_AbortAfterPage1_NeverClosesUnseenRows(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, installation := "listing-abort-"+token, "installation-abort-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_variations WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})

	priorSync := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	const total = 30
	const page1 = 10 // phase A: ids [0,10) — the only page an aborted run ever saw
	const page2 = 20 // phase B resume: ids [10,20) additionally seen; [20,30) genuinely vanished
	statusFor := func(i int) string {
		if i%2 == 0 {
			return "active"
		}
		return "paused"
	}
	for i := 0; i < total; i++ {
		item := fmt.Sprintf("item-%02d", i)
		seedListingWithLastSeen(t, pool, tenant, installation, item, "-", "title "+item, statusFor(i), priorSync, priorSync, &priorSync)
	}
	rowFor := func(runStarted time.Time, i int) domain.Listing {
		item := fmt.Sprintf("item-%02d", i)
		status := domain.ListingStatus(statusFor(i))
		return listing(t, tenant, installation, item, "-", status, "title "+item, nil, nil, nil, runStarted, 50)
	}
	countAbsent := func() int {
		n := 0
		for i := 0; i < total; i++ {
			item := fmt.Sprintf("item-%02d", i)
			if readAbsentSince(t, pool, tenant, installation, item, "-") != nil {
				n++
			}
		}
		return n
	}

	runStarted := priorSync.Add(time.Hour)
	repo := listingspostgres.NewRepository(pool, tenant)

	// --- Phase A: page 1 only, then the run aborts. No MarkRunComplete. ---
	rowsPage1 := make([]domain.Listing, 0, page1)
	for i := 0; i < page1; i++ {
		rowsPage1 = append(rowsPage1, rowFor(runStarted, i))
	}
	if err := repo.UpsertPulledRows(ctx, installation, rowsPage1, runStarted); err != nil {
		t.Fatalf("UpsertPulledRows() page 1: %v", err)
	}

	if got := countAbsent(); got != 0 {
		t.Fatalf("absent_since count after aborted phase A = %d, want 0 (incomplete run must never mark absent)", got)
	}
	for i := page1; i < total; i++ {
		item := fmt.Sprintf("item-%02d", i)
		got := readRawListing(t, pool, tenant, installation, item, "-")
		if got.Status != statusFor(i) {
			t.Fatalf("unseen row %s status = %q, want unchanged %q (mass-closure regression)", item, got.Status, statusFor(i))
		}
		if !got.Updated.Equal(priorSync) {
			t.Fatalf("unseen row %s updated_at = %v, want untouched %v", item, got.Updated, priorSync)
		}
	}
	for i := 0; i < page1; i++ {
		item := fmt.Sprintf("item-%02d", i)
		if got := readRawListing(t, pool, tenant, installation, item, "-").Status; got != statusFor(i) {
			t.Fatalf("seen row %s status = %q, want %q", item, got, statusFor(i))
		}
	}

	// --- Phase B: full retomada. Same run (same runStarted) resumes,
	// ticks in ids [10,20), then declares the run complete. ids [20,30)
	// are never seen in any tick of this run — genuinely vanished. ---
	rowsPage2 := make([]domain.Listing, 0, page2-page1)
	for i := page1; i < page2; i++ {
		rowsPage2 = append(rowsPage2, rowFor(runStarted, i))
	}
	if err := repo.UpsertPulledRows(ctx, installation, rowsPage2, runStarted); err != nil {
		t.Fatalf("UpsertPulledRows() page 2 (resume): %v", err)
	}
	if err := repo.MarkRunComplete(ctx, installation, runStarted); err != nil {
		t.Fatalf("MarkRunComplete() after full retomada: %v", err)
	}

	if got := countAbsent(); got != total-page2 {
		t.Fatalf("absent_since count after completed retomada = %d, want %d (only genuinely unseen ids)", got, total-page2)
	}
	for i := 0; i < page2; i++ {
		item := fmt.Sprintf("item-%02d", i)
		if absent := readAbsentSince(t, pool, tenant, installation, item, "-"); absent != nil {
			t.Fatalf("seen-in-run row %s absent_since = %v, want nil", item, *absent)
		}
	}
	for i := page2; i < total; i++ {
		item := fmt.Sprintf("item-%02d", i)
		absent := readAbsentSince(t, pool, tenant, installation, item, "-")
		if absent == nil {
			t.Fatalf("genuinely vanished row %s absent_since = nil, want set", item)
		}
		if !absent.Equal(runStarted) {
			t.Fatalf("vanished row %s absent_since = %v, want %v", item, *absent, runStarted)
		}
		if got := readRawListing(t, pool, tenant, installation, item, "-").Status; got != statusFor(i) {
			t.Fatalf("vanished row %s status = %q, want unchanged %q (absence must never flip status)", item, got, statusFor(i))
		}
	}
}

// TestApplyCompletedPull_CompleteRunMarksGenuinelyAbsentRowWithoutTouchingStatus
// is the IC-06 positive case: once a run is declared COMPLETE, a row this
// run legitimately never saw again gets absent_since stamped — but status
// stays exactly what it was; absence is never inferred as closed.
func TestApplyCompletedPull_CompleteRunMarksGenuinelyAbsentRowWithoutTouchingStatus(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, installation := "listing-absent-"+token, "installation-absent-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})

	priorSync := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	seedListingWithLastSeen(t, pool, tenant, installation, "stays", "-", "stays", "active", priorSync, priorSync, &priorSync)
	seedListingWithLastSeen(t, pool, tenant, installation, "vanished", "-", "vanished", "paused", priorSync, priorSync, &priorSync)

	runStarted := priorSync.Add(time.Hour)
	rows := []domain.Listing{
		listing(t, tenant, installation, "stays", "-", domain.ListingStatusActive, "stays", nil, nil, nil, runStarted, 50),
	}
	repo := listingspostgres.NewRepository(pool, tenant)
	if err := repo.UpsertPulledRows(ctx, installation, rows, runStarted); err != nil {
		t.Fatalf("UpsertPulledRows(): %v", err)
	}
	if err := repo.MarkRunComplete(ctx, installation, runStarted); err != nil {
		t.Fatalf("MarkRunComplete(): %v", err)
	}

	if got := readRawListing(t, pool, tenant, installation, "vanished", "-").Status; got != "paused" {
		t.Fatalf("vanished row status = %q, want unchanged %q (absence must never flip status)", got, "paused")
	}
	absent := readAbsentSince(t, pool, tenant, installation, "vanished", "-")
	if absent == nil {
		t.Fatal("vanished row absent_since = nil, want set")
	}
	if !absent.Equal(runStarted) {
		t.Fatalf("vanished row absent_since = %v, want %v", *absent, runStarted)
	}
	if got := readAbsentSince(t, pool, tenant, installation, "stays", "-"); got != nil {
		t.Fatalf("stays row absent_since = %v, want nil", *got)
	}
}

// TestApplyCompletedPull_ReappearanceClearsAbsentSince proves a row marked
// absent by an earlier complete run is un-marked the moment it is seen
// again — absence is not sticky once the provider re-reports the item.
func TestApplyCompletedPull_ReappearanceClearsAbsentSince(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, installation := "listing-reappear-"+token, "installation-reappear-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})

	priorSync := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	wentAbsent := priorSync.Add(time.Hour)
	seedListingWithAbsentSince(t, pool, tenant, installation, "back", "-", "back", "active", priorSync, priorSync, &priorSync, &wentAbsent)

	runStarted := wentAbsent.Add(time.Hour)
	rows := []domain.Listing{
		listing(t, tenant, installation, "back", "-", domain.ListingStatusActive, "back", nil, nil, nil, runStarted, 50),
	}
	repo := listingspostgres.NewRepository(pool, tenant)
	if err := repo.UpsertPulledRows(ctx, installation, rows, runStarted); err != nil {
		t.Fatalf("UpsertPulledRows(): %v", err)
	}
	if err := repo.MarkRunComplete(ctx, installation, runStarted); err != nil {
		t.Fatalf("MarkRunComplete(): %v", err)
	}

	if got := readAbsentSince(t, pool, tenant, installation, "back", "-"); got != nil {
		t.Fatalf("reappeared row absent_since = %v, want nil", *got)
	}
	var lastSeen time.Time
	if err := pool.QueryRow(ctx, `SELECT last_seen_at FROM listings WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id=$3 AND variation_id=$4`, tenant, installation, "back", "-").Scan(&lastSeen); err != nil {
		t.Fatal(err)
	}
	if !lastSeen.Equal(runStarted) {
		t.Fatalf("last_seen_at = %v, want advanced to %v", lastSeen, runStarted)
	}
}

// TestMarkRunComplete_NeverUpsertedIsNoOp is the run-scoped safety pin
// required test (4): a run that never called UpsertPulledRows at all (e.g. a
// broken/empty enumerator whose every tick came back empty) must have
// MarkRunComplete alone be a no-op — never read as "the whole catalog is
// gone". This replaces the old per-CALL pin
// (TestApplyCompletedPull_ZeroRowsNeverMarksAbsentEvenWhenCallerClaimsComplete):
// that pin only guarded a single call's `len(rows) == 0`, which would have
// wrongly disabled keep-absent for an entire multi-tick run just because its
// LAST tick happened to be empty (see
// TestApplyCompletedPull_AbortAfterPage1_NeverClosesUnseenRows phase B, where
// a real terminal tick legitimately has zero NEW rows for ids already seen
// in an earlier tick but must still complete the run). The new pin is
// scoped to "did this run ever touch anything", checked directly against
// last_seen_at = runStartedAt, not to any single call's row count.
func TestMarkRunComplete_NeverUpsertedIsNoOp(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, installation := "listing-zero-"+token, "installation-zero-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})

	priorSync := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	seedListingWithLastSeen(t, pool, tenant, installation, "solo", "-", "solo", "active", priorSync, priorSync, &priorSync)

	runStarted := priorSync.Add(time.Hour)
	repo := listingspostgres.NewRepository(pool, tenant)
	// UpsertPulledRows is never called for this run at all — MarkRunComplete
	// is the only call the run ever makes.
	if err := repo.MarkRunComplete(ctx, installation, runStarted); err != nil {
		t.Fatalf("MarkRunComplete(): %v", err)
	}

	if got := readAbsentSince(t, pool, tenant, installation, "solo", "-"); got != nil {
		t.Fatalf("solo row absent_since = %v, want nil (run-scoped safety pin: nothing was ever upserted this run)", *got)
	}
	if got := readRawListing(t, pool, tenant, installation, "solo", "-").Status; got != "active" {
		t.Fatalf("solo row status = %q, want unchanged %q", got, "active")
	}
}

// TestApplyCompletedPull_StatusPersistsVerbatimEvenForNovelProviderValues
// proves the 0090 CHECK-drop and this writer cooperate: a status value that
// never appeared in any CHECK constraint this repo ever wrote persists
// exactly as given, with no app-level normalization or rejection.
func TestApplyCompletedPull_StatusPersistsVerbatimEvenForNovelProviderValues(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, installation := "listing-novel-"+token, "installation-novel-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})

	runStarted := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	// Built directly (not via domain.NewListing/ListingInput), which is the
	// point: this is the provider's own vocabulary passing straight through
	// to the writer with zero app-level validation in its path, exactly as
	// IC-07 ("status: verbatim do provider ... sem CHECK") requires.
	key, err := domain.NewListingKey(tenant, installation, "novel-item", "-")
	if err != nil {
		t.Fatal(err)
	}
	novel := domain.Listing{
		Key:       key,
		Provider:  "mercado_livre",
		Title:     "novel status item",
		Status:    domain.ListingStatus("under_review_custom"),
		SyncState: domain.ListingSyncStateSynced,
		FetchedAt: &runStarted,
		CreatedAt: runStarted,
		UpdatedAt: runStarted,
	}
	// Second row built THROUGH domain.NewListing (the 0a fix path): proves a
	// novel status doesn't just avoid a DB CHECK, it also clears domain
	// construction and still comes out verbatim end to end (required test 7:
	// "status verbatim survives domain->repository via the 0a fix").
	throughDomain := listing(t, tenant, installation, "novel-item-domain", "-", domain.ListingStatus("suspended"), "novel via NewListing", nil, nil, nil, runStarted, 50)

	repo := listingspostgres.NewRepository(pool, tenant)
	if err := repo.UpsertPulledRows(ctx, installation, []domain.Listing{novel, throughDomain}, runStarted); err != nil {
		t.Fatalf("UpsertPulledRows(): %v", err)
	}
	if got := readRawListing(t, pool, tenant, installation, "novel-item", "-").Status; got != "under_review_custom" {
		t.Fatalf("status = %q, want verbatim %q", got, "under_review_custom")
	}
	if got := readRawListing(t, pool, tenant, installation, "novel-item-domain", "-").Status; got != "suspended" {
		t.Fatalf("status (via NewListing) = %q, want verbatim %q", got, "suspended")
	}
}

// TestApplyCompletedPull_VariationsUpsertIdempotently proves listing_variations
// rows land alongside the parent in the same transaction, that a second call
// updates in place (no duplicate rows for the same PK tuple), and — required
// test (5), kill-and-resume — that reprocessing the IDENTICAL page a third
// time (as a restarted run replaying the same scroll_id page would) is fully
// idempotent: no duplicate rows, no overcounted quantities.
func TestApplyCompletedPull_VariationsUpsertIdempotently(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_listings")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, installation := "listing-var-"+token, "installation-var-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_variations WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
		_, _ = pool.Exec(context.Background(), `DELETE FROM listings WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	})

	runStarted := time.Date(2026, 7, 1, 9, 0, 0, 0, time.UTC)
	priceV1 := domain.PriceAmount("10.00")
	priceV2 := domain.PriceAmount("20.00")
	sku1, sku2 := "SKU-1", "SKU-2"
	qty := 5
	row := listing(t, tenant, installation, "parent-item", "-", domain.ListingStatusActive, "parent", nil, nil, nil, runStarted, 50)
	row.Variations = []domain.ListingVariation{
		{VariationID: "var-1", Price: &priceV1, AvailableQuantity: &qty, SellerSKU: &sku1},
		{VariationID: "var-2", Price: &priceV2, AvailableQuantity: &qty, SellerSKU: &sku2},
	}
	repo := listingspostgres.NewRepository(pool, tenant)
	if err := repo.UpsertPulledRows(ctx, installation, []domain.Listing{row}, runStarted); err != nil {
		t.Fatalf("UpsertPulledRows(): %v", err)
	}

	got := readListingVariations(t, pool, tenant, installation, "parent-item")
	if len(got) != 2 {
		t.Fatalf("variations = %d, want 2: %#v", len(got), got)
	}
	if got["var-1"].price != "10.00" || got["var-1"].sellerSKU != sku1 {
		t.Fatalf("var-1 = %#v", got["var-1"])
	}
	if got["var-2"].price != "20.00" || got["var-2"].sellerSKU != sku2 {
		t.Fatalf("var-2 = %#v", got["var-2"])
	}

	// Second run: var-1's price changes, var-2 repeated unchanged. Must
	// update in place, never duplicate.
	priceV1b := domain.PriceAmount("15.00")
	row2 := listing(t, tenant, installation, "parent-item", "-", domain.ListingStatusActive, "parent", nil, nil, nil, runStarted.Add(time.Hour), 50)
	row2.Variations = []domain.ListingVariation{
		{VariationID: "var-1", Price: &priceV1b, AvailableQuantity: &qty, SellerSKU: &sku1},
		{VariationID: "var-2", Price: &priceV2, AvailableQuantity: &qty, SellerSKU: &sku2},
	}
	if err := repo.UpsertPulledRows(ctx, installation, []domain.Listing{row2}, runStarted.Add(time.Hour)); err != nil {
		t.Fatalf("UpsertPulledRows() second call: %v", err)
	}
	got = readListingVariations(t, pool, tenant, installation, "parent-item")
	if len(got) != 2 {
		t.Fatalf("variations after second call = %d, want 2 (no duplicate insert): %#v", len(got), got)
	}
	if got["var-1"].price != "15.00" {
		t.Fatalf("var-1 price after second call = %q, want 15.00", got["var-1"].price)
	}

	// Required test (5), kill-and-resume: reprocess the IDENTICAL page a
	// third time (same rows, same content) with a distinct tick timestamp —
	// the same scroll_id page delivered twice after a restart must be fully
	// idempotent: no duplicate variation rows, no overcounted quantities.
	if err := repo.UpsertPulledRows(ctx, installation, []domain.Listing{row2}, runStarted.Add(2*time.Hour)); err != nil {
		t.Fatalf("UpsertPulledRows() reprocessed identical page: %v", err)
	}
	got = readListingVariations(t, pool, tenant, installation, "parent-item")
	if len(got) != 2 {
		t.Fatalf("variations after reprocessing identical page = %d, want 2 (no duplicates on kill-and-resume replay): %#v", len(got), got)
	}
	if got["var-1"].price != "15.00" || *got["var-1"].availableQuantity != qty {
		t.Fatalf("var-1 after replay = %#v, want price 15.00 qty %d (no overcount)", got["var-1"], qty)
	}
	if got["var-2"].price != "20.00" || *got["var-2"].availableQuantity != qty {
		t.Fatalf("var-2 after replay = %#v, want price 20.00 qty %d (no overcount)", got["var-2"], qty)
	}
	if got["var-1"].price != "15.00" {
		t.Fatalf("var-1 price after update = %q, want 15.00", got["var-1"].price)
	}
}

func seedListing(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation, title, status string, created, updated time.Time) {
	t.Helper()
	seedListingWithLastSeen(t, pool, tenant, installation, item, variation, title, status, created, updated, nil)
}

func seedListingWithLastSeen(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation, title, status string, created, updated time.Time, lastSeen *time.Time) {
	t.Helper()
	seedListingWithAbsentSince(t, pool, tenant, installation, item, variation, title, status, created, updated, lastSeen, nil)
}

func seedListingWithAbsentSince(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation, title, status string, created, updated time.Time, lastSeen, absentSince *time.Time) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `INSERT INTO listings
		(tenant_id,installation_id,provider,provider_listing_id,variation_id,title,status,sync_state,created_at,updated_at,fetched_at,last_seen_at,absent_since)
		VALUES ($1,$2,'mercado_livre',$3,$4,$5,$6,'synced',$7,$8,$8,$9,$10)`, tenant, installation, item, variation, title, status, created, updated, lastSeen, absentSince)
	if err != nil {
		t.Fatalf("seed listing: %v", err)
	}
}

func listing(t *testing.T, tenant, installation, item, variation string, status domain.ListingStatus, title string, price *domain.PriceAmount, currency *domain.PriceCurrency, quantity *int, fetched time.Time, quality int) domain.Listing {
	t.Helper()
	row, err := domain.NewListing(domain.ListingInput{TenantID: tenant, InstallationID: installation, Provider: "mercado_livre", ProviderListingID: item, VariationID: variation, Title: title, Status: status, PriceAmount: price, PriceCurrency: currency, PublishedQuantity: quantity, SyncState: domain.ListingSyncStateSynced, QualityScore: &quality, FetchedAt: &fetched, CreatedAt: fetched, UpdatedAt: fetched})
	if err != nil {
		t.Fatal(err)
	}
	return row
}

type rawListing struct {
	Provider, Title, Status, Price, Currency string
	Quantity                                 *int
	Created, Updated, Fetched                time.Time
}

func readRawListing(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation string) rawListing {
	t.Helper()
	var got rawListing
	var price, currency *string
	err := pool.QueryRow(context.Background(), `SELECT provider,title,status,price_amount::text,price_currency,published_quantity,created_at,updated_at,fetched_at FROM listings WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id=$3 AND variation_id=$4`, tenant, installation, item, variation).Scan(&got.Provider, &got.Title, &got.Status, &price, &currency, &got.Quantity, &got.Created, &got.Updated, &got.Fetched)
	if err != nil {
		t.Fatalf("read listing %s/%s: %v", item, variation, err)
	}
	if price != nil {
		got.Price = *price
	}
	if currency != nil {
		got.Currency = *currency
	}
	return got
}

func readAbsentSince(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation string) *time.Time {
	t.Helper()
	var got *time.Time
	if err := pool.QueryRow(context.Background(), `SELECT absent_since FROM listings WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id=$3 AND variation_id=$4`, tenant, installation, item, variation).Scan(&got); err != nil {
		t.Fatalf("read absent_since %s/%s: %v", item, variation, err)
	}
	return got
}

type rawVariation struct {
	price             string
	availableQuantity *int
	sellerSKU         string
}

func readListingVariations(t *testing.T, pool *pgxpool.Pool, tenant, installation, item string) map[string]rawVariation {
	t.Helper()
	rows, err := pool.Query(context.Background(), `SELECT variation_id,price::text,available_quantity,seller_sku FROM listing_variations WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id=$3`, tenant, installation, item)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got := map[string]rawVariation{}
	for rows.Next() {
		var variationID string
		var v rawVariation
		var price, sku *string
		if err := rows.Scan(&variationID, &price, &v.availableQuantity, &sku); err != nil {
			t.Fatal(err)
		}
		if price != nil {
			v.price = *price
		}
		if sku != nil {
			v.sellerSKU = *sku
		}
		got[variationID] = v
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return got
}

func assertListing(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation, status, title string, created, updated, fetched time.Time, price, currency string, quantity int) {
	t.Helper()
	got := readRawListing(t, pool, tenant, installation, item, variation)
	if got.Status != status || got.Title != title || !got.Created.Equal(created) || !got.Updated.Equal(updated) || !got.Fetched.Equal(fetched) || got.Price != price || got.Currency != currency {
		t.Fatalf("listing %s/%s = %#v", item, variation, got)
	}
	if quantity < 0 {
		if got.Quantity != nil {
			t.Fatalf("quantity = %v, want nil", *got.Quantity)
		}
	} else if got.Quantity == nil || *got.Quantity != quantity {
		t.Fatalf("quantity = %v, want %d", got.Quantity, quantity)
	}
}
func assertStatus(t *testing.T, pool *pgxpool.Pool, tenant, installation, item, variation, want string) {
	t.Helper()
	var got string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM listings WHERE tenant_id=$1 AND installation_id=$2 AND provider_listing_id=$3 AND variation_id=$4`, tenant, installation, item, variation).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("status %s/%s = %s, want %s", item, variation, got, want)
	}
}
func assertEventKinds(t *testing.T, pool *pgxpool.Pool, tenant, installation string, want map[string]string) {
	t.Helper()
	rows, err := pool.Query(context.Background(), `SELECT provider_listing_id,variation_id,kind FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got := map[string]string{}
	for rows.Next() {
		var item, v, k string
		if err := rows.Scan(&item, &v, &k); err != nil {
			t.Fatal(err)
		}
		got[item+"/"+v] = k
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("event kinds = %#v, want %#v", got, want)
	}
}
func eventCount(t *testing.T, pool *pgxpool.Pool, tenant, installation string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM listing_sync_events WHERE tenant_id=$1 AND installation_id=$2`, tenant, installation).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}
