package postgres_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	profitabilitypostgres "marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres"
	profitabilityapp "marketplace-central/apps/server_core/internal/modules/profitability/application"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
)

func TestManualAdjustmentsAppendOnlyReadbackAndConstraints(t *testing.T) {
	databaseURL := os.Getenv("MC_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("MC_DATABASE_URL not set; real PostgreSQL append-only validation pending")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect PostgreSQL: %v", err)
	}
	defer pool.Close()

	testID := integrationRandomToken(t)
	tenantID := "profitability-test-" + testID
	installationID := "installation-" + testID
	orderID := "order-" + testID
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `
			DELETE FROM profitability_manual_adjustments
			WHERE tenant_id = $1 AND installation_id = $2
		`, tenantID, installationID)
	})

	createdAt := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	store := profitabilitypostgres.NewStore(pool, tenantID)
	service := profitabilityapp.NewService(profitabilityapp.ServiceConfig{
		Adjustments: store,
		Now:         func() time.Time { return createdAt },
	})

	inputs := []profitabilityapp.CreateManualAdjustmentInput{
		{
			InstallationID:  installationID,
			IdempotencyKey:  "freight-" + testID,
			ProviderOrderID: orderID,
			Scope:           profitabilitydomain.InputScopeOrder,
			Category:        profitabilitydomain.ManualAdjustmentFreight,
			Amount:          -12.50,
			Currency:        "BRL",
			Reason:          " freight correction ",
			Actor: profitabilitydomain.ActorMetadata{
				ActorType: " operator ",
				ActorID:   " operator-1 ",
				ActorName: " Leandro ",
			},
		},
		{
			InstallationID:  installationID,
			IdempotencyKey:  "commission-" + testID,
			ProviderOrderID: orderID,
			Scope:           profitabilitydomain.InputScopeOrder,
			Category:        profitabilitydomain.ManualAdjustmentCommission,
			Amount:          3.25,
			Currency:        "BRL",
			Reason:          " commission correction ",
			Actor: profitabilitydomain.ActorMetadata{
				ActorType: " operator ",
				ActorID:   " operator-1 ",
				ActorName: " Leandro ",
			},
		},
	}

	created := make([]profitabilitydomain.ManualAdjustment, 0, len(inputs))
	for _, input := range inputs {
		adjustment, err := service.CreateManualAdjustment(ctx, input)
		if err != nil {
			t.Fatalf("CreateManualAdjustment(%s): %v", input.Category, err)
		}
		created = append(created, adjustment)
	}
	retry := inputs[0]
	retry.Amount = 999.99
	retry.Reason = "changed retry payload must not rewrite audit"
	duplicate, err := service.CreateManualAdjustment(ctx, retry)
	if err != nil {
		t.Fatalf("CreateManualAdjustment(duplicate): %v", err)
	}
	if duplicate != created[0] {
		t.Fatalf("duplicate adjustment = %+v, want original %+v", duplicate, created[0])
	}
	if created[0].AdjustmentID == created[1].AdjustmentID {
		t.Fatalf("same-timestamp adjustments reused id %q", created[0].AdjustmentID)
	}

	stored, err := store.ListAdjustments(ctx, installationID, 10)
	if err != nil {
		t.Fatalf("ListAdjustments(): %v", err)
	}
	if len(stored) != 2 {
		t.Fatalf("ListAdjustments() returned %d rows, want 2", len(stored))
	}
	byCategory := make(map[profitabilitydomain.ManualAdjustmentCategory]profitabilitydomain.ManualAdjustment, len(stored))
	for _, adjustment := range stored {
		byCategory[adjustment.Category] = adjustment
		if adjustment.Actor.ActorType != "operator" || adjustment.Actor.ActorID != "operator-1" || adjustment.Actor.ActorName != "Leandro" {
			t.Fatalf("incomplete actor audit after readback: %+v", adjustment.Actor)
		}
		if adjustment.Reason == "" || !adjustment.CreatedAt.Equal(createdAt) {
			t.Fatalf("incomplete reason/time audit after readback: %+v", adjustment)
		}
	}
	if byCategory[profitabilitydomain.ManualAdjustmentFreight].Amount != -12.50 {
		t.Fatalf("freight signed amount changed: %+v", byCategory[profitabilitydomain.ManualAdjustmentFreight])
	}
	if byCategory[profitabilitydomain.ManualAdjustmentCommission].Amount != 3.25 {
		t.Fatalf("commission signed amount changed: %+v", byCategory[profitabilitydomain.ManualAdjustmentCommission])
	}

	constraintCases := []struct {
		name                string
		providerItemID      string
		providerVariationID string
		scope               string
		category            string
		currency            string
		reason              string
		actorType           string
		actorID             string
	}{
		{name: "empty actor type", scope: "order", category: "freight", currency: "BRL", reason: "reason", actorID: "operator-1"},
		{name: "empty actor id", scope: "order", category: "freight", currency: "BRL", reason: "reason", actorType: "operator"},
		{name: "empty reason", scope: "order", category: "freight", currency: "BRL", actorType: "operator", actorID: "operator-1"},
		{name: "empty currency", scope: "order", category: "freight", reason: "reason", actorType: "operator", actorID: "operator-1"},
		{name: "invalid scope", scope: "listing", category: "freight", currency: "BRL", reason: "reason", actorType: "operator", actorID: "operator-1"},
		{name: "invalid category", scope: "order", category: "rebate", currency: "BRL", reason: "reason", actorType: "operator", actorID: "operator-1"},
		{name: "order scope with provider item", providerItemID: "MLB-1", scope: "order", category: "freight", currency: "BRL", reason: "reason", actorType: "operator", actorID: "operator-1"},
		{name: "order scope with provider variation", providerVariationID: "variation-1", scope: "order", category: "freight", currency: "BRL", reason: "reason", actorType: "operator", actorID: "operator-1"},
		{name: "item scope without provider item", scope: "item", category: "freight", currency: "BRL", reason: "reason", actorType: "operator", actorID: "operator-1"},
	}
	for _, tt := range constraintCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pool.Exec(ctx, `
				INSERT INTO profitability_manual_adjustments (
					tenant_id, adjustment_id, installation_id, provider_order_id,
					provider_item_id, provider_variation_id,
					scope, category, amount, currency, reason,
					actor_type, actor_id, created_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1, $9, $10, $11, $12, $13)
			`, tenantID, "constraint-"+tt.name+"-"+testID, installationID, orderID,
				tt.providerItemID, tt.providerVariationID, tt.scope, tt.category,
				tt.currency, tt.reason, tt.actorType, tt.actorID, createdAt)
			if err == nil {
				t.Fatal("invalid adjustment bypassed PostgreSQL constraint")
			}
		})
	}
}

func integrationRandomToken(t *testing.T) string {
	t.Helper()
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("generate integration test token: %v", err)
	}
	return hex.EncodeToString(b)
}
