//go:build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	mutationspostgres "marketplace-central/apps/server_core/internal/modules/mutations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestApplyRepositorySnapshotClaimChunkAndImmutableOutcomes(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mutations_apply")
	tenant := "mutation-apply-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id = $1`, tenant)
	})
	repo := mutationspostgres.NewRepository(pool, tenant)
	created, err := repo.CreateProtocol(ctx, ports.CreateProtocolInput{InstallationID: "installation-a", Type: domain.ProtocolTypePriceUpdate, Actor: "operator_supplied_unverified", Intent: json.RawMessage(`{"price":{"currency":"BRL","amount":"49.90"}}`), Selection: json.RawMessage(`{"mode":"explicit"}`), CreatedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	inputs := make([]ports.ReplaceItemInput, 23)
	for i := range inputs {
		inputs[i] = ports.ReplaceItemInput{ListingID: fmt.Sprintf("MLB-%03d", i+1), Before: json.RawMessage(`{"price":{"currency":"BRL","amount":"39.90"},"stock":{"available":7}}`), After: json.RawMessage(`{"price":{"amount":"49.90","currency":"BRL"},"stock":{"available":7}}`)}
	}
	items, err := repo.ReplaceItems(ctx, created.ProtocolID, inputs)
	if err != nil {
		t.Fatalf("ReplaceItems: %v", err)
	}
	if items[0].ItemID != created.ProtocolID+"-001" || items[0].IdempotencyKey != created.ProtocolID+":MLB-001" {
		t.Fatalf("generated item = %#v", items[0])
	}
	assertJSONEqual(t, items[0].After, inputs[0].After)
	approvedAt := time.Now().UTC()
	if err := repo.ApproveItems(ctx, created.ProtocolID, approvedAt); err != nil {
		t.Fatalf("ApproveItems: %v", err)
	}
	if _, err := repo.ReplaceItems(ctx, created.ProtocolID, inputs[:1]); err == nil {
		t.Fatal("ReplaceItems on approved protocol succeeded")
	}
	claim, found, err := repo.ClaimProtocol(ctx, "installation-a")
	if err != nil || !found {
		t.Fatalf("ClaimProtocol: found=%v err=%v", found, err)
	}
	defer claim.Rollback(ctx)
	if claim.Protocol().State != domain.ProtocolStateApplying || claim.Protocol().ApprovedAt == nil {
		t.Fatalf("claim protocol = %#v", claim.Protocol())
	}
	other, found, err := repo.ClaimProtocol(ctx, "installation-a")
	if err != nil || found || other != nil {
		t.Fatalf("second claim: claim=%v found=%v err=%v", other, found, err)
	}
	chunk, err := claim.FetchPendingItems(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunk) != 20 {
		t.Fatalf("chunk size=%d, want 20", len(chunk))
	}
	for i, item := range chunk {
		if item.Seq != i+1 {
			t.Fatalf("chunk[%d].Seq=%d", i, item.Seq)
		}
	}
	appliedAt := time.Now().UTC()
	if err := claim.WriteItemOutcome(ctx, chunk[0].ItemID, ports.ItemOutcome{State: domain.ItemStateApplied, AppliedAt: appliedAt}); err != nil {
		t.Fatal(err)
	}
	if err := claim.WriteItemOutcome(ctx, chunk[0].ItemID, ports.ItemOutcome{State: domain.ItemStateFailed, AppliedAt: appliedAt}); err == nil {
		t.Fatal("terminal item update succeeded")
	}
	keys, err := claim.AppliedIdempotencyKeys(ctx, []string{chunk[0].IdempotencyKey, chunk[1].IdempotencyKey})
	if err != nil || !reflect.DeepEqual(keys, []string{chunk[0].IdempotencyKey}) {
		t.Fatalf("applied keys=%v err=%v", keys, err)
	}
	if err := claim.WriteItemOutcome(ctx, chunk[1].ItemID, ports.ItemOutcome{State: domain.ItemStateFailed, Failure: json.RawMessage(`{"code":"provider_unavailable","message_pt":"indisponível","retryable":true}`), AppliedAt: appliedAt}); err != nil {
		t.Fatal(err)
	}
	if err := claim.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	var failedAppliedAt *time.Time
	if err := pool.QueryRow(ctx, `SELECT applied_at FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2 AND item_id=$3`, tenant, created.ProtocolID, chunk[1].ItemID).Scan(&failedAppliedAt); err != nil {
		t.Fatal(err)
	}
	if failedAppliedAt != nil {
		t.Fatalf("failed item applied_at = %v, want NULL", failedAppliedAt)
	}
}

func TestApplyRepositoryClaimsDifferentInstallationsConcurrently(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mutations_claim_installations")
	tenant := "mutation-claim-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id = $1`, tenant)
	})
	now := time.Now().UTC()
	for i, installation := range []string{"installation-a", "installation-b"} {
		_, err := pool.Exec(ctx, `INSERT INTO mutation_protocols (tenant_id,protocol_id,installation_id,type,state,actor,intent,selection,totals,source_as_of,created_at,previewed_at,approved_at) VALUES ($1,$2,$3,'price_update','approved','operator_supplied_unverified','{}','{}','{}',$4,$4,$4,$4)`, tenant, fmt.Sprintf("MP-%06d", i+1), installation, now.Add(time.Duration(i)*time.Second))
		if err != nil {
			t.Fatal(err)
		}
	}
	repo := mutationspostgres.NewRepository(pool, tenant)
	a, found, err := repo.ClaimProtocol(ctx, "installation-a")
	if err != nil || !found {
		t.Fatalf("claim a found=%v err=%v", found, err)
	}
	defer a.Rollback(ctx)
	b, found, err := repo.ClaimProtocol(ctx, "installation-b")
	if err != nil || !found {
		t.Fatalf("claim b found=%v err=%v", found, err)
	}
	defer b.Rollback(ctx)
}

func assertJSONEqual(t *testing.T, got, want json.RawMessage) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &w); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g, w) {
		t.Fatalf("JSON got=%s want=%s", got, want)
	}
}
