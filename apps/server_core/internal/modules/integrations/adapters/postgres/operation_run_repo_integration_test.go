//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestOperationRunRepositoryBeginExclusiveIsAtomic(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "operation_run_exclusive")
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	provider := "exclusive_provider_" + token
	if _, err := pool.Exec(ctx, `INSERT INTO integration_provider_definitions(provider_code,family,display_name,auth_strategy,install_mode) VALUES($1,'marketplace',$1,'none','manual')`, provider); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM integration_provider_definitions WHERE provider_code=$1`, provider)
	})
	seed := func(tenant, installation string) {
		if _, err := pool.Exec(ctx, `INSERT INTO integration_installations(tenant_id,installation_id,provider_code,family,display_name,status) VALUES($1,$2,$3,'marketplace',$2,'connected')`, tenant, installation, provider); err != nil {
			t.Fatal(err)
		}
	}
	seed("tenant-a-"+token, "inst-a-"+token)
	seed("tenant-a-"+token, "inst-b-"+token)
	seed("tenant-b-"+token, "inst-c-"+token)

	tenant, installation := "tenant-a-"+token, "inst-a-"+token
	repo := integrationspostgres.NewOperationRunRepository(pool, tenant)
	start := make(chan struct{})
	results := make(chan struct {
		run    domain.OperationRun
		active bool
		err    error
	}, 2)
	var wg sync.WaitGroup
	for _, id := range []string{"op-1-" + token, "op-2-" + token} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			<-start
			run, active, err := repo.BeginExclusive(ctx, testRun(tenant, installation, id))
			results <- struct {
				run    domain.OperationRun
				active bool
				err    error
			}{run, active, err}
		}(id)
	}
	close(start)
	wg.Wait()
	close(results)
	activeCount := 0
	ids := map[string]bool{}
	for result := range results {
		if result.err != nil {
			t.Fatal(result.err)
		}
		if result.active {
			activeCount++
		}
		ids[result.run.OperationRunID] = true
	}
	if activeCount != 1 || len(ids) != 1 {
		t.Fatalf("active losers=%d returned ids=%v", activeCount, ids)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM integration_operation_runs WHERE tenant_id=$1 AND installation_id=$2 AND operation_type='listings_refresh'`, tenant, installation).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}

	for _, tc := range []struct{ tenant, installation, id string }{{tenant, "inst-b-" + token, "op-other-installation" + token}, {"tenant-b-" + token, "inst-c-" + token, "op-other-tenant" + token}} {
		got, active, err := integrationspostgres.NewOperationRunRepository(pool, tc.tenant).BeginExclusive(ctx, testRun(tc.tenant, tc.installation, tc.id))
		if err != nil || active || got.OperationRunID != tc.id {
			t.Fatalf("independent start got=%+v active=%v err=%v", got, active, err)
		}
	}
}

func TestLatestRunsByModuleTenantScoped(t *testing.T) {
	ctx, repo, pool, tenant, other, installation, otherInstallation := runLatestRunsHarness(t, "tenant")
	now := time.Now().UTC().Truncate(time.Microsecond)
	successAt := now.Add(-time.Minute)
	seedLatestRun(t, pool, tenant, installation, "latest-success", "listing_read", domain.OperationRunStatusSucceeded, successAt, &successAt)
	seedLatestRun(t, pool, tenant, installation, "latest-attempt", "listing_read", domain.OperationRunStatusFailed, now, &now)
	seedLatestRun(t, pool, tenant, installation, "other-module", "order_read", domain.OperationRunStatusSucceeded, now, &now)
	seedLatestRun(t, pool, other, otherInstallation, "other-tenant", "pricing_fee_sync", domain.OperationRunStatusSucceeded, now.Add(time.Minute), &now)

	got, err := repo.LatestRunsByModule(ctx, installation)
	if err != nil {
		t.Fatalf("LatestRunsByModule() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("projection rows = %d, want 2: %+v", len(got), got)
	}
	byModule := make(map[string]ports.LatestRunByModule, len(got))
	for _, item := range got {
		byModule[item.OperationType] = item
	}
	if _, ok := byModule["pricing_fee_sync"]; ok {
		t.Fatal("projection leaked another tenant's operation module")
	}
	otherTenantRuns, err := repo.LatestRunsByModule(ctx, otherInstallation)
	if err != nil || len(otherTenantRuns) != 0 {
		t.Fatalf("other-tenant installation leaked: rows=%+v err=%v", otherTenantRuns, err)
	}
	if byModule["listing_read"].LatestAttemptedAt == nil || !byModule["listing_read"].LatestAttemptedAt.Equal(now) {
		t.Fatalf("listing_read timestamps = %+v", byModule["listing_read"])
	}
	if byModule["listing_read"].LatestSuccessfulAt == nil || !byModule["listing_read"].LatestSuccessfulAt.Equal(successAt) {
		t.Fatalf("listing_read successful timestamp = %v, want %v", byModule["listing_read"].LatestSuccessfulAt, successAt)
	}
}

func TestLatestRunsKeepsUnknownTimestampNull(t *testing.T) {
	ctx, repo, pool, tenant, _, installation, _ := runLatestRunsHarness(t, "null")
	now := time.Now().UTC().Truncate(time.Microsecond)
	seedLatestRun(t, pool, tenant, installation, "running", "listing_read", domain.OperationRunStatusRunning, now, nil)

	got, err := repo.LatestRunsByModule(ctx, installation)
	if err != nil {
		t.Fatalf("LatestRunsByModule() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("projection rows = %d, want 1: %+v", len(got), got)
	}
	if got[0].LatestSuccessfulAt != nil {
		t.Fatalf("latest successful timestamp = %v, want nil", got[0].LatestSuccessfulAt)
	}
	if got[0].LatestAttemptedAt == nil || !got[0].LatestAttemptedAt.Equal(now) {
		t.Fatalf("latest attempted timestamp = %v, want %v", got[0].LatestAttemptedAt, now)
	}
}

func runLatestRunsHarness(t *testing.T, name string) (context.Context, *integrationspostgres.OperationRunRepository, *pgxpool.Pool, string, string, string, string) {
	t.Helper()
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "operation_run_latest_"+name)
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant, other := "latest-a-"+token, "latest-b-"+token
	installation, otherInstallation := "latest-inst-a-"+token, "latest-inst-b-"+token
	provider := "latest-provider-" + token
	if _, err := pool.Exec(ctx, `INSERT INTO integration_provider_definitions(provider_code,family,display_name,auth_strategy,install_mode) VALUES($1,'marketplace',$1,'none','manual')`, provider); err != nil {
		t.Fatal(err)
	}
	for owner, install := range map[string]string{tenant: installation, other: otherInstallation} {
		if _, err := pool.Exec(ctx, `INSERT INTO integration_installations(tenant_id,installation_id,provider_code,family,display_name,status) VALUES($1,$2,$3,'marketplace',$2,'connected')`, owner, install, provider); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM integration_provider_definitions WHERE provider_code=$1`, provider)
	})
	return ctx, integrationspostgres.NewOperationRunRepository(pool, tenant), pool, tenant, other, installation, otherInstallation
}

func seedLatestRun(t *testing.T, pool *pgxpool.Pool, tenant, installation, id, module string, status domain.OperationRunStatus, startedAt time.Time, finishedAt *time.Time) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `INSERT INTO integration_operation_runs(tenant_id,operation_run_id,installation_id,operation_type,status,result_code,failure_code,attempt_count,duration_ms,started_at,completed_at) VALUES($1,$2,$3,$4,$5,'result','failure',2,123,$6,$7)`, tenant, id, installation, module, status, startedAt, finishedAt); err != nil {
		t.Fatal(err)
	}
}

func testRun(tenant, installation, id string) domain.OperationRun {
	now := time.Now().UTC()
	return domain.OperationRun{TenantID: tenant, InstallationID: installation, OperationRunID: id, OperationType: "listings_refresh", Status: domain.OperationRunStatusQueued, CreatedAt: now, UpdatedAt: now}
}

// A run left "running" by a process that died mid-refresh is what BeginExclusive
// keeps reporting as active, so the operator's Atualizar button stays refused
// forever. The reaper closes those out at boot; these tests hold it to doing only
// that, and only for this tenant.
func TestFailInterruptedRunsUnblocksTheOperation(t *testing.T) {
	ctx, repo, pool, tenant, other, installation, otherInstallation := runLatestRunsHarness(t, "interrupted")
	now := time.Now().UTC().Truncate(time.Microsecond)
	seedLatestRun(t, pool, tenant, installation, "orphan", "listings_refresh", domain.OperationRunStatusRunning, now.Add(-time.Hour), nil)
	seedLatestRun(t, pool, tenant, installation, "finished", "listing_read", domain.OperationRunStatusSucceeded, now.Add(-2*time.Hour), &now)
	seedLatestRun(t, pool, other, otherInstallation, "other-tenant-orphan", "listings_refresh", domain.OperationRunStatusRunning, now.Add(-time.Hour), nil)

	reaped, err := repo.FailInterruptedRuns(ctx, now)
	if err != nil {
		t.Fatalf("FailInterruptedRuns() error = %v", err)
	}
	if reaped != 1 {
		t.Fatalf("reaped = %d, want 1 (only this tenant's unfinished run)", reaped)
	}

	var status, failureCode string
	var completedAt *time.Time
	if err := pool.QueryRow(ctx, `SELECT status, failure_code, completed_at FROM integration_operation_runs WHERE tenant_id=$1 AND operation_run_id='orphan'`, tenant).Scan(&status, &failureCode, &completedAt); err != nil {
		t.Fatal(err)
	}
	if status != string(domain.OperationRunStatusFailed) {
		t.Fatalf("status = %q, want failed", status)
	}
	// The seed carries a real failure cause; the reaper must not overwrite it.
	if failureCode != "failure" {
		t.Fatalf("failure_code = %q, want the pre-existing cause preserved", failureCode)
	}
	// We never observed a completion, so we do not invent a completion time.
	if completedAt != nil {
		t.Fatalf("completed_at = %v, want NULL", completedAt)
	}

	if _, active, err := repo.BeginExclusive(ctx, testRun(tenant, installation, "after-reap")); err != nil || active {
		t.Fatalf("BeginExclusive after reap: active=%v err=%v, want a fresh run", active, err)
	}

	var otherStatus string
	if err := pool.QueryRow(ctx, `SELECT status FROM integration_operation_runs WHERE tenant_id=$1 AND operation_run_id='other-tenant-orphan'`, other).Scan(&otherStatus); err != nil {
		t.Fatal(err)
	}
	if otherStatus != string(domain.OperationRunStatusRunning) {
		t.Fatalf("other tenant's run status = %q, want untouched", otherStatus)
	}
}

// A run killed by a restart never recorded why it stopped, so the reaper supplies
// the only cause it can honestly name.
func TestFailInterruptedRunsNamesTheCauseWhenNoneWasRecorded(t *testing.T) {
	ctx, repo, pool, tenant, _, installation, _ := runLatestRunsHarness(t, "interrupted_cause")
	now := time.Now().UTC().Truncate(time.Microsecond)
	if _, err := pool.Exec(ctx, `INSERT INTO integration_operation_runs(tenant_id,operation_run_id,installation_id,operation_type,status,started_at) VALUES($1,'no-cause',$2,'listings_refresh','running',$3)`, tenant, installation, now.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}

	if _, err := repo.FailInterruptedRuns(ctx, now); err != nil {
		t.Fatalf("FailInterruptedRuns() error = %v", err)
	}

	var failureCode string
	if err := pool.QueryRow(ctx, `SELECT failure_code FROM integration_operation_runs WHERE tenant_id=$1 AND operation_run_id='no-cause'`, tenant).Scan(&failureCode); err != nil {
		t.Fatal(err)
	}
	if failureCode != "interrupted" {
		t.Fatalf("failure_code = %q, want interrupted", failureCode)
	}
}
