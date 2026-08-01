package application_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	synccomposition "marketplace-central/apps/server_core/internal/modules/sync/composition"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
	"marketplace-central/apps/server_core/internal/modules/tenant_config"
)

// regressionStore is a minimal fake syncapp.StateStore local to this external
// test package. The internal fakeStore in scheduler_test.go is unexported and
// lives in package application — invisible from package application_test,
// which is a distinct package used specifically so this file can import
// sync/composition (the real job body) without an import cycle.
type regressionStore struct {
	mu    sync.Mutex
	state domain.SyncState

	successes []regressionSuccess
}

type regressionSuccess struct {
	cursor      json.RawMessage
	incremental bool
}

func (r *regressionStore) Read(context.Context, string, domain.Entity) (domain.SyncState, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state, true, nil
}

func (r *regressionStore) RecordSuccess(_ context.Context, _ string, _ domain.Entity, cursor json.RawMessage, _ time.Time, incremental bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.successes = append(r.successes, regressionSuccess{cursor: cursor, incremental: incremental})
	return nil
}

func (r *regressionStore) RecordFailure(context.Context, string, domain.Entity, domain.SyncError) error {
	return nil
}

type regressionLookup struct{ source tenant_config.ActiveSource }

func (l regressionLookup) Get(_ context.Context, tenantID string) (tenant_config.Config, error) {
	return tenant_config.Config{TenantID: tenantID, Source: l.source}, nil
}

type regressionAdapter struct {
	readports.Reader
	result readports.SyncResult
}

func (a *regressionAdapter) Sync(context.Context) (readports.SyncResult, error) {
	return a.result, nil
}

func (a *regressionAdapter) Kind() sourcekind.SourceKind { return sourcekind.UploadSnapshot }

// TestRealProductsJobThroughSchedulerRecordsFullSync is the binding
// regression for F-03: the REAL products job (synccomposition.NewProductsJob,
// untouched by this feature) produces a ProductsCursor with no "phase" key at
// all. Run through the real scheduler, RunOnce must still record
// incremental=false — identical to the behavior before this change — and the
// persisted cursor's fields must be exactly what the job produced, compared
// structurally (never a byte-exact JSONB round-trip, a known trap here).
func TestRealProductsJobThroughSchedulerRecordsFullSync(t *testing.T) {
	const tenantID = "11111111-1111-4111-8111-111111111111"
	previous := json.RawMessage(`{"source":"xlsx","processed":3,"completed_at":"2026-07-01T00:00:00Z"}`)
	store := &regressionStore{state: domain.SyncState{Cursor: previous}}
	at := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return at }

	scheduler := syncapp.NewScheduler(store, synccomposition.InstallationScopeERP, time.Minute, clock)

	job := synccomposition.NewProductsJob(
		tenantID,
		regressionLookup{source: tenant_config.SourceXLSX},
		map[tenant_config.ActiveSource]readports.ProductSourceAdapter{
			tenant_config.SourceXLSX: &regressionAdapter{result: readports.SyncResult{Processed: 1845}},
		},
		clock,
	)
	if err := scheduler.RegisterJob(domain.EntityProducts, job); err != nil {
		t.Fatalf("register: %v", err)
	}

	scheduler.RunOnce(context.Background())

	if len(store.successes) != 1 {
		t.Fatalf("successes = %d, want 1", len(store.successes))
	}
	got := store.successes[0]

	// (a) incremental must be false: the real job's cursor has no phase key.
	if got.incremental {
		t.Error("incremental = true, want false: the products job's ProductsCursor has no phase key")
	}

	// (b) the persisted cursor is exactly what the job produced.
	var cursor synccomposition.ProductsCursor
	if err := json.Unmarshal(got.cursor, &cursor); err != nil {
		t.Fatalf("persisted cursor did not unmarshal: %v", err)
	}
	if cursor.Source != string(tenant_config.SourceXLSX) {
		t.Errorf("cursor.Source = %q, want xlsx", cursor.Source)
	}
	if cursor.Processed != 1845 {
		t.Errorf("cursor.Processed = %d, want 1845", cursor.Processed)
	}
	if !cursor.CompletedAt.Equal(at) {
		t.Errorf("cursor.CompletedAt = %v, want %v", cursor.CompletedAt, at)
	}
}
