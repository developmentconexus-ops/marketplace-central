//go:build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	mutationspostgres "marketplace-central/apps/server_core/internal/modules/mutations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestRetryCloneAndReadPagesAreTenantScopedAndStable(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mutation_retry_read")
	token := time.Now().UTC().Format("150405.000000000")
	tenant, other, installation := "mutation-read-a-"+token, "mutation-read-b-"+token, "inst-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id = ANY($1)`, []string{tenant, other})
	})
	created1 := time.Date(2026, 7, 16, 10, 0, 0, 0, time.UTC)
	created2 := created1.Add(time.Minute)
	seedProtocol := func(owner, id string, state domain.ProtocolState, created time.Time) {
		t.Helper()
		_, err := pool.Exec(ctx, `INSERT INTO mutation_protocols(tenant_id,protocol_id,installation_id,type,state,actor,intent,selection,totals,created_at,finished_at) VALUES($1,$2,$3,'price_update',$4,'operator',$5,$6,'{"items":3}', $7,$7)`, owner, id, installation, state, json.RawMessage(`{ "new_price": "10.00" }`), json.RawMessage(`{"mode":"filter","q":"camisa"}`), created)
		if err != nil {
			t.Fatal(err)
		}
	}
	seedProtocol(tenant, "MP-000001", domain.ProtocolStatePartiallyFailed, created1)
	seedProtocol(tenant, "MP-000002", domain.ProtocolStateApplied, created2)
	seedProtocol(other, "MP-000001", domain.ProtocolStatePartiallyFailed, created2.Add(time.Minute))
	for seq, failure := range []string{
		`{"code":"provider_validation","retryable":true}`,
		`{"code":"provider_unavailable","retryable":false}`,
		`{"code":"internal","retryable":true}`,
	} {
		_, err := pool.Exec(ctx, `INSERT INTO mutation_items(tenant_id,protocol_id,seq,item_id,listing_id,idempotency_key,before,after,state,failure) VALUES($1,'MP-000001',$2,$3,$4,$5,$6,$7,'failed',$8)`, tenant, seq+1, "old-item", "listing-"+string(rune('a'+seq)), "MP-000001:key-"+string(rune('a'+seq)), json.RawMessage(`{"old":1}`), json.RawMessage(`{"new":2}`), json.RawMessage(failure))
		if err != nil {
			t.Fatal(err)
		}
	}
	var before string
	if err := pool.QueryRow(ctx, `SELECT jsonb_agg(to_jsonb(i) ORDER BY seq)::text FROM mutation_items i WHERE tenant_id=$1 AND protocol_id='MP-000001'`, tenant).Scan(&before); err != nil {
		t.Fatal(err)
	}
	repo := mutationspostgres.NewRepository(pool, tenant)
	clone, eligible, err := repo.CloneRetry(ctx, "MP-000001", created2.Add(2*time.Minute))
	if err != nil || !eligible {
		t.Fatalf("clone=%+v eligible=%v err=%v", clone, eligible, err)
	}
	if clone.ProtocolID != "MP-000003" || clone.RetriedFrom == nil || *clone.RetriedFrom != "MP-000001" || clone.State != domain.ProtocolStateDraft {
		t.Fatalf("clone=%+v", clone)
	}
	var cloneItems int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2`, tenant, clone.ProtocolID).Scan(&cloneItems); err != nil || cloneItems != 0 {
		t.Fatalf("clone items=%d err=%v", cloneItems, err)
	}
	var after string
	if err := pool.QueryRow(ctx, `SELECT jsonb_agg(to_jsonb(i) ORDER BY seq)::text FROM mutation_items i WHERE tenant_id=$1 AND protocol_id='MP-000001'`, tenant).Scan(&after); err != nil || after != before {
		t.Fatalf("original changed before=%s after=%s err=%v", before, after, err)
	}

	page, err := repo.ListProtocols(ctx, application.ProtocolQuery{InstallationID: installation, Limit: 2})
	if err != nil || len(page.Items) != 2 || page.Items[0].ProtocolID != clone.ProtocolID || page.NextCursor == nil {
		t.Fatalf("protocol page=%+v err=%v", page, err)
	}
	next, err := repo.ListProtocols(ctx, application.ProtocolQuery{InstallationID: installation, Limit: 2, Cursor: *page.NextCursor})
	if err != nil || len(next.Items) != 1 || next.Items[0].ProtocolID != "MP-000001" {
		t.Fatalf("next page=%+v err=%v", next, err)
	}
	items, err := repo.ListItems(ctx, application.ItemQuery{ProtocolID: "MP-000001", Limit: 2})
	if err != nil || len(items.Items) != 2 || items.Items[0].Seq != 1 || items.Items[1].Seq != 2 || items.NextCursor == nil {
		t.Fatalf("items=%+v err=%v", items, err)
	}
	items2, err := repo.ListItems(ctx, application.ItemQuery{ProtocolID: "MP-000001", Limit: 2, Cursor: *items.NextCursor})
	if err != nil || len(items2.Items) != 1 || items2.Items[0].Seq != 3 {
		t.Fatalf("items2=%+v err=%v", items2, err)
	}
	otherRepo := mutationspostgres.NewRepository(pool, other)
	foreign, err := otherRepo.ListProtocols(ctx, application.ProtocolQuery{InstallationID: installation, Limit: 10})
	if err != nil || len(foreign.Items) != 1 {
		t.Fatalf("foreign=%+v err=%v", foreign, err)
	}
}
