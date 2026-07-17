//go:build integration

package application_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	mutationspostgres "marketplace-central/apps/server_core/internal/modules/mutations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/mutations/adapters/stub"
	"marketplace-central/apps/server_core/internal/modules/mutations/application"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestPollerPassResumesPostgresClaimWithoutResendingAppliedItem(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mutations_poller_resume")
	tenant := "mutation-poller-resume-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id=$1`, tenant)
	})
	repo := mutationspostgres.NewRepository(pool, tenant)
	created, err := repo.CreateProtocol(ctx, ports.CreateProtocolInput{InstallationID: "installation-a", Type: domain.ProtocolTypePriceUpdate, Actor: "operator_supplied_unverified", Intent: json.RawMessage(`{"new_price":{"amount":"49.90","currency":"BRL"}}`), Selection: json.RawMessage(`{"mode":"explicit"}`), CreatedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	sourceAsOf := time.Now().UTC()
	items, err := repo.ReplaceItems(ctx, created.ProtocolID, []ports.ReplaceItemInput{
		{ListingID: "MLB-001", After: json.RawMessage(`{"price":{"amount":"49.90"}}`)},
		{ListingID: "MLB-002", After: json.RawMessage(`{"price":{"amount":"59.90"}}`)},
		{ListingID: "MLB-003", After: json.RawMessage(`{"price":{"amount":"69.90"}}`)},
	}, &sourceAsOf)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.ApproveItems(ctx, created.ProtocolID, time.Now().UTC()); err != nil {
		t.Fatal(err)
	}

	writer := stub.NewWriter(nil)
	if worked, err := application.NewPoller(&crashAfterOutcomeRepository{ProtocolRepository: repo}, writer, time.Now).Pass(ctx, "installation-a"); !worked || !errors.Is(err, errSimulatedCrash) {
		t.Fatalf("first Pass() worked=%v err=%v", worked, err)
	}
	afterCrash, found, err := repo.GetProtocol(ctx, created.ProtocolID)
	if err != nil || !found || afterCrash.State != domain.ProtocolStateApproved {
		t.Fatalf("after crash protocol=%+v found=%v err=%v", afterCrash, found, err)
	}
	var firstState domain.ItemState
	if err := pool.QueryRow(ctx, `SELECT state FROM mutation_items WHERE tenant_id=$1 AND item_id=$2`, tenant, items[0].ItemID).Scan(&firstState); err != nil || firstState != domain.ItemStateApplied {
		t.Fatalf("first item state=%q err=%v", firstState, err)
	}

	if worked, err := application.NewPoller(repo, writer, time.Now).Pass(ctx, "installation-a"); err != nil || !worked {
		t.Fatalf("resumed Pass() worked=%v err=%v", worked, err)
	}
	if got := writer.Keys(); len(got) != 3 || got[0] != items[0].IdempotencyKey || got[1] != items[1].IdempotencyKey || got[2] != items[2].IdempotencyKey {
		t.Fatalf("writer keys=%v; applied key was resent", got)
	}
	finished, found, err := repo.GetProtocol(ctx, created.ProtocolID)
	if err != nil || !found || finished.State != domain.ProtocolStateApplied || finished.FinishedAt == nil {
		t.Fatalf("finished protocol=%+v found=%v err=%v", finished, found, err)
	}
}

var errSimulatedCrash = errors.New("simulated crash after durable outcome")

type crashAfterOutcomeRepository struct{ ports.ProtocolRepository }

func (r *crashAfterOutcomeRepository) ClaimProtocol(ctx context.Context, installationID string) (ports.ProtocolClaim, bool, error) {
	claim, found, err := r.ProtocolRepository.ClaimProtocol(ctx, installationID)
	if err != nil || !found {
		return claim, found, err
	}
	return &crashAfterOutcomeClaim{ProtocolClaim: claim}, true, nil
}

type crashAfterOutcomeClaim struct {
	ports.ProtocolClaim
	written bool
}

func (c *crashAfterOutcomeClaim) WriteItemOutcome(ctx context.Context, itemID string, outcome ports.ItemOutcome) error {
	if err := c.ProtocolClaim.WriteItemOutcome(ctx, itemID, outcome); err != nil {
		return err
	}
	if !c.written {
		c.written = true
		return errSimulatedCrash
	}
	return nil
}
