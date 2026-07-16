//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	integrationspostgres "marketplace-central/apps/server_core/internal/modules/integrations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/integrations/domain"
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

func testRun(tenant, installation, id string) domain.OperationRun {
	now := time.Now().UTC()
	return domain.OperationRun{TenantID: tenant, InstallationID: installation, OperationRunID: id, OperationType: "listings_refresh", Status: domain.OperationRunStatusQueued, CreatedAt: now, UpdatedAt: now}
}
