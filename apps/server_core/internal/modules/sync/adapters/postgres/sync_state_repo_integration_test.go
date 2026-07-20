//go:build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	syncpostgres "marketplace-central/apps/server_core/internal/modules/sync/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// token isolates rows per run without Date.now() collisions across tests.
func newToken(t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// M01-C3: cursor read/write end-to-end.
func TestSyncStateCursorRoundTrip(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "sync_state_c3")
	tenant := "tenant-c3-" + newToken(t)
	inst := "erp"
	repo := syncpostgres.NewSyncStateRepository(pool, tenant)
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant) })

	// First run: never synced — honest not-found, no fabricated cursor.
	if _, found, err := repo.Read(ctx, inst, domain.EntityProducts); err != nil || found {
		t.Fatalf("first read found=%v err=%v, want found=false", found, err)
	}

	at := time.Date(2026, 7, 20, 9, 30, 0, 0, time.UTC)
	if err := repo.RecordSuccess(ctx, inst, domain.EntityProducts, json.RawMessage(`{"page":3}`), at, false); err != nil {
		t.Fatalf("record success: %v", err)
	}

	state, found, err := repo.Read(ctx, inst, domain.EntityProducts)
	if err != nil || !found {
		t.Fatalf("read after success found=%v err=%v", found, err)
	}
	if state.LastFullSyncAt == nil || !state.LastFullSyncAt.Equal(at) {
		t.Errorf("last_full_sync_at = %v, want %v", state.LastFullSyncAt, at)
	}
	if state.LastIncrementalAt != nil {
		t.Errorf("last_incremental_at must stay NULL after a full sync, got %v", state.LastIncrementalAt)
	}
	// Assert the cursor SEMANTICALLY, not byte-exact: JSONB is a normalized store,
	// so Postgres re-serializes `{"page":3}` as `{"page": 3}` (space after the
	// colon). A raw string equality on a JSONB round-trip is inherently wrong — it
	// tests Postgres' whitespace, not our fidelity. Unmarshal and compare the value.
	var gotCursor struct {
		Page int `json:"page"`
	}
	if err := json.Unmarshal(state.Cursor, &gotCursor); err != nil {
		t.Fatalf("cursor %q not valid JSON: %v", state.Cursor, err)
	}
	if gotCursor.Page != 3 {
		t.Errorf("cursor page = %d, want 3 (raw: %q)", gotCursor.Page, state.Cursor)
	}
	if state.ConsecutiveFailures != 0 {
		t.Errorf("consecutive_failures = %d, want 0", state.ConsecutiveFailures)
	}
	if state.LastError != nil {
		t.Errorf("last_error must be NULL after success, got %+v", state.LastError)
	}
}

// M01-C4: failure records last_error + increments consecutive_failures; a later
// success clears both.
func TestSyncStateFailureThenRecovery(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "sync_state_c4")
	tenant := "tenant-c4-" + newToken(t)
	inst := "erp"
	repo := syncpostgres.NewSyncStateRepository(pool, tenant)
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant) })

	for i := 1; i <= 2; i++ {
		if err := repo.RecordFailure(ctx, inst, domain.EntityProducts, domain.SyncError{
			Message: "job failed", At: time.Date(2026, 7, 20, 9, 0, i, 0, time.UTC),
		}); err != nil {
			t.Fatalf("record failure %d: %v", i, err)
		}
	}

	state, found, err := repo.Read(ctx, inst, domain.EntityProducts)
	if err != nil || !found {
		t.Fatalf("read found=%v err=%v", found, err)
	}
	if state.ConsecutiveFailures != 2 {
		t.Errorf("consecutive_failures = %d, want 2", state.ConsecutiveFailures)
	}
	if state.LastError == nil || state.LastError.Message != "job failed" {
		t.Errorf("last_error = %+v, want message=job failed", state.LastError)
	}

	// Recovery resets the counter and clears the error.
	if err := repo.RecordSuccess(ctx, inst, domain.EntityProducts, nil, time.Now().UTC(), false); err != nil {
		t.Fatalf("record success: %v", err)
	}
	state, _, err = repo.Read(ctx, inst, domain.EntityProducts)
	if err != nil {
		t.Fatal(err)
	}
	if state.ConsecutiveFailures != 0 || state.LastError != nil {
		t.Errorf("after recovery: failures=%d last_error=%+v, want 0/nil", state.ConsecutiveFailures, state.LastError)
	}
	// nil cursor on success must persist as SQL NULL, not a fabricated {}.
	if state.Cursor != nil {
		t.Errorf("cursor must be NULL, got %q", state.Cursor)
	}
}

// M01-C11: concurrent atomic cursor appends never lose an update.
func TestSyncStateConcurrentCursorAppend(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "sync_state_c11")
	tenant := "tenant-c11-" + newToken(t)
	inst := "erp"
	repo := syncpostgres.NewSyncStateRepository(pool, tenant)
	t.Cleanup(func() { _, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id=$1`, tenant) })

	const n = 50
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			if err := repo.AppendPendingCodigo(ctx, inst, domain.EntityMarket, fmt.Sprintf("cod-%03d", i)); err != nil {
				t.Errorf("append %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	state, found, err := repo.Read(ctx, inst, domain.EntityMarket)
	if err != nil || !found {
		t.Fatalf("read found=%v err=%v", found, err)
	}
	var cursor struct {
		Pending []string `json:"pending"`
	}
	if err := json.Unmarshal(state.Cursor, &cursor); err != nil {
		t.Fatalf("unmarshal cursor %q: %v", state.Cursor, err)
	}
	if len(cursor.Pending) != n {
		t.Fatalf("pending has %d codes, want %d — a concurrent append was lost", len(cursor.Pending), n)
	}
	sort.Strings(cursor.Pending)
	for i := 0; i < n; i++ {
		want := fmt.Sprintf("cod-%03d", i)
		if cursor.Pending[i] != want {
			t.Errorf("pending[%d] = %s, want %s", i, cursor.Pending[i], want)
		}
	}
}

// AC-01 / M01-C10: a repository scoped to one tenant cannot read another tenant's
// row for the same (installation, entity) key.
func TestSyncStateTenantIsolation(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "sync_state_tenant")
	tok := newToken(t)
	tenantA := "tenant-a-" + tok
	tenantB := "tenant-b-" + tok
	repoA := syncpostgres.NewSyncStateRepository(pool, tenantA)
	repoB := syncpostgres.NewSyncStateRepository(pool, tenantB)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM sync_state WHERE tenant_id IN ($1,$2)`, tenantA, tenantB)
	})

	if err := repoA.RecordSuccess(ctx, "erp", domain.EntityProducts, json.RawMessage(`{"a":1}`), time.Now().UTC(), false); err != nil {
		t.Fatal(err)
	}
	if _, found, err := repoB.Read(ctx, "erp", domain.EntityProducts); err != nil || found {
		t.Errorf("tenant B read found=%v err=%v, want found=false (isolation breach)", found, err)
	}
}
