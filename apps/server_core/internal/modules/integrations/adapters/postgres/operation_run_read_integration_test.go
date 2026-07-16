//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
	"marketplace-central/apps/server_core/internal/modules/integrations/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestListRunsTenantScopedNewestFirstWithCursor(t *testing.T) {
	ctx, repo, pool, tenant, other, installation := runReadHarness(t, "cursor")
	base := time.Now().UTC().Add(-time.Hour)
	seedReadRun(t, pool, tenant, installation, "run-1", "listing_read", domain.OperationRunStatusSucceeded, base, &base)
	seedReadRun(t, pool, tenant, installation, "run-2", "listing_read", domain.OperationRunStatusSucceeded, base.Add(time.Minute), &base)
	seedReadRun(t, pool, tenant, installation, "run-3", "listing_read", domain.OperationRunStatusSucceeded, base.Add(time.Minute), &base)
	seedReadRun(t, pool, other, installation, "leak", "listing_read", domain.OperationRunStatusSucceeded, base.Add(time.Hour), &base)

	first, err := repo.ListRuns(ctx, ports.RunListQuery{InstallationID: installation, Limit: 2})
	if err != nil || len(first.Items) != 2 || first.Items[0].OperationRunID != "run-3" || first.Items[1].OperationRunID != "run-2" || first.NextCursor == nil {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	last, err := repo.ListRuns(ctx, ports.RunListQuery{InstallationID: installation, Limit: 2, Cursor: *first.NextCursor})
	if err != nil || len(last.Items) != 1 || last.Items[0].OperationRunID != "run-1" || last.NextCursor != nil {
		t.Fatalf("last=%+v err=%v", last, err)
	}
}

func TestListRunsExcludesStartedBeforeNinetyDays(t *testing.T) {
	ctx, repo, pool, tenant, _, installation := runReadHarness(t, "window")
	now := time.Now().UTC()
	seedReadRun(t, pool, tenant, installation, "recent", "order_read", domain.OperationRunStatusSucceeded, now.Add(-89*24*time.Hour), &now)
	seedReadRun(t, pool, tenant, installation, "old", "order_read", domain.OperationRunStatusSucceeded, now.Add(-91*24*time.Hour), &now)
	page, err := repo.ListRuns(ctx, ports.RunListQuery{InstallationID: installation})
	if err != nil || len(page.Items) != 1 || page.Items[0].OperationRunID != "recent" {
		t.Fatalf("page=%+v err=%v", page, err)
	}
}

func TestListRunsFiltersStatusAndModule(t *testing.T) {
	ctx, repo, pool, tenant, _, installation := runReadHarness(t, "filters")
	now := time.Now().UTC()
	seedReadRun(t, pool, tenant, installation, "match", "pricing_fee_sync", domain.OperationRunStatusFailed, now, &now)
	seedReadRun(t, pool, tenant, installation, "wrong-module", "listing_read", domain.OperationRunStatusFailed, now, &now)
	seedReadRun(t, pool, tenant, installation, "wrong-status", "pricing_fee_sync", domain.OperationRunStatusSucceeded, now, &now)
	page, err := repo.ListRuns(ctx, ports.RunListQuery{InstallationID: installation, Module: "pricing_fee_sync", Status: domain.OperationRunStatusFailed})
	if err != nil || len(page.Items) != 1 || page.Items[0].OperationRunID != "match" || page.Items[0].Module != "pricing_fee_sync" {
		t.Fatalf("page=%+v err=%v", page, err)
	}
}

func TestRunningRunPreservesNullFinishedAt(t *testing.T) {
	ctx, repo, pool, tenant, _, installation := runReadHarness(t, "null-finished")
	now := time.Now().UTC()
	seedReadRun(t, pool, tenant, installation, "running", "stock_read", domain.OperationRunStatusRunning, now, nil)
	page, err := repo.ListRuns(ctx, ports.RunListQuery{InstallationID: installation})
	if err != nil || len(page.Items) != 1 || page.Items[0].FinishedAt != nil {
		t.Fatalf("page=%+v err=%v", page, err)
	}
}

func runReadHarness(t *testing.T, name string) (context.Context, *integrationspostgres.OperationRunRepository, *pgxpool.Pool, string, string, string) {
	t.Helper()
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "operation_run_read_"+name)
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	tenant, other, installation := "run-a-"+token, "run-b-"+token, "run-inst-"+token
	provider := "run-provider-" + token
	if _, err := pool.Exec(ctx, `INSERT INTO integration_provider_definitions(provider_code,family,display_name,auth_strategy,install_mode) VALUES($1,'marketplace',$1,'none','manual')`, provider); err != nil {
		t.Fatal(err)
	}
	for _, owner := range []string{tenant, other} {
		if _, err := pool.Exec(ctx, `INSERT INTO integration_installations(tenant_id,installation_id,provider_code,family,display_name,status) VALUES($1,$2,$3,'marketplace',$2,'connected')`, owner, installation, provider); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM integration_provider_definitions WHERE provider_code=$1`, provider)
	})
	return ctx, integrationspostgres.NewOperationRunRepository(pool, tenant), pool, tenant, other, installation
}

func seedReadRun(t *testing.T, pool *pgxpool.Pool, tenant, installation, id, module string, status domain.OperationRunStatus, startedAt time.Time, finishedAt *time.Time) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `INSERT INTO integration_operation_runs(tenant_id,operation_run_id,installation_id,operation_type,status,result_code,failure_code,attempt_count,duration_ms,started_at,completed_at) VALUES($1,$2,$3,$4,$5,'result','failure',2,123,$6,$7)`, tenant, id, installation, module, status, startedAt, finishedAt); err != nil {
		t.Fatal(err)
	}
}
