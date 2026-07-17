//go:build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	mutationspostgres "marketplace-central/apps/server_core/internal/modules/mutations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestCancelProtocolPersistsCancelledStateAndTimestampRoundTrip(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mutation_cancel")
	tenant := "mutation-cancel-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id=$1`, tenant)
	})
	repo := mutationspostgres.NewRepository(pool, tenant)
	created, err := repo.CreateProtocol(ctx, ports.CreateProtocolInput{InstallationID: "inst-1", Type: domain.ProtocolTypePriceUpdate, Actor: "operator_supplied_unverified", Intent: json.RawMessage(`{"new_price":{"amount":"49.90"}}`), Selection: json.RawMessage(`{"mode":"explicit"}`), CreatedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	cancelledAt := time.Date(2026, 7, 16, 12, 0, 0, 123456000, time.FixedZone("BRT", -3*60*60))
	if err := repo.CancelProtocol(ctx, created.ProtocolID, cancelledAt); err != nil {
		t.Fatal(err)
	}
	got, found, err := repo.GetProtocol(ctx, created.ProtocolID)
	if err != nil || !found || got.State != domain.ProtocolStateCancelled || got.FinishedAt == nil || !got.FinishedAt.Equal(cancelledAt) {
		t.Fatalf("found=%v protocol=%+v err=%v", found, got, err)
	}
	if err := repo.CancelProtocol(ctx, created.ProtocolID, cancelledAt.Add(time.Minute)); !errors.Is(err, domain.ErrInvalidStateTransition) {
		t.Fatalf("second cancel error=%v", err)
	}
}
