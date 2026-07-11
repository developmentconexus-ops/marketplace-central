//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	profitabilitypostgres "marketplace-central/apps/server_core/internal/modules/profitability/adapters/postgres"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestProfitSnapshotRealizationPersistence(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_profit_snapshots")

	testID := integrationRandomToken(t)
	tenantID := "profit-snapshot-test-" + testID
	installationID := "installation-" + testID
	orderIDs := []string{"realized-" + testID, "cancelled-" + testID, "unknown-" + testID}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `
			DELETE FROM profitability_profit_snapshots
			WHERE tenant_id = $1 AND installation_id = $2
		`, tenantID, installationID)
	})

	contribution := 8.0
	margin := 40.0
	createdAt := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)
	snapshots := []profitabilitydomain.ProfitSnapshot{
		{
			InstallationID:     installationID,
			ProviderOrderID:    orderIDs[0],
			Scope:              profitabilitydomain.InputScopeOrder,
			Currency:           "BRL",
			RealizationState:   profitabilitydomain.OrderRealizationRealized,
			ContributionAmount: &contribution,
			MarginPercent:      &margin,
			Quality:            profitabilitydomain.ProfitSnapshotComplete,
			Flags:              []profitabilitydomain.ProfitSnapshotFlag{},
			CreatedAt:          createdAt,
			UpdatedAt:          createdAt,
		},
		{
			InstallationID:   installationID,
			ProviderOrderID:  orderIDs[1],
			Scope:            profitabilitydomain.InputScopeOrder,
			Currency:         "BRL",
			RealizationState: profitabilitydomain.OrderRealizationNotRealized,
			Quality:          profitabilitydomain.ProfitSnapshotNotRealized,
			Flags:            []profitabilitydomain.ProfitSnapshotFlag{profitabilitydomain.ProfitFlagOrderCancelled},
			CreatedAt:        createdAt.Add(time.Minute),
			UpdatedAt:        createdAt.Add(time.Minute),
		},
		{
			InstallationID:   installationID,
			ProviderOrderID:  orderIDs[2],
			Scope:            profitabilitydomain.InputScopeOrder,
			Currency:         "BRL",
			RealizationState: profitabilitydomain.OrderRealizationUnknown,
			Quality:          profitabilitydomain.ProfitSnapshotIncomplete,
			Flags:            []profitabilitydomain.ProfitSnapshotFlag{profitabilitydomain.ProfitFlagOrderStateUnknown},
			CreatedAt:        createdAt.Add(2 * time.Minute),
			UpdatedAt:        createdAt.Add(2 * time.Minute),
		},
	}

	store := profitabilitypostgres.NewStore(pool, tenantID)
	if err := store.ReplaceSnapshots(ctx, installationID, orderIDs, snapshots); err != nil {
		t.Fatalf("ReplaceSnapshots(): %v", err)
	}

	stored, err := store.ListSnapshots(ctx, installationID, 10)
	if err != nil {
		t.Fatalf("ListSnapshots(): %v", err)
	}
	if len(stored) != len(snapshots) {
		t.Fatalf("ListSnapshots() returned %d rows, want %d", len(stored), len(snapshots))
	}

	byOrderID := make(map[string]profitabilitydomain.ProfitSnapshot, len(stored))
	for _, snapshot := range stored {
		byOrderID[snapshot.ProviderOrderID] = snapshot
	}

	realized := byOrderID[orderIDs[0]]
	if realized.RealizationState != profitabilitydomain.OrderRealizationRealized || realized.Quality != profitabilitydomain.ProfitSnapshotComplete {
		t.Fatalf("realized snapshot contract changed after readback: %+v", realized)
	}
	if realized.ContributionAmount == nil || *realized.ContributionAmount != contribution || realized.MarginPercent == nil || *realized.MarginPercent != margin {
		t.Fatalf("realized monetary values changed after readback: %+v", realized)
	}

	cancelled := byOrderID[orderIDs[1]]
	assertSnapshotRealization(t, cancelled, profitabilitydomain.OrderRealizationNotRealized, profitabilitydomain.ProfitSnapshotNotRealized, profitabilitydomain.ProfitFlagOrderCancelled)
	unknown := byOrderID[orderIDs[2]]
	assertSnapshotRealization(t, unknown, profitabilitydomain.OrderRealizationUnknown, profitabilitydomain.ProfitSnapshotIncomplete, profitabilitydomain.ProfitFlagOrderStateUnknown)
}

func assertSnapshotRealization(t *testing.T, snapshot profitabilitydomain.ProfitSnapshot, realization profitabilitydomain.OrderRealizationState, quality profitabilitydomain.ProfitSnapshotQuality, flag profitabilitydomain.ProfitSnapshotFlag) {
	t.Helper()
	if snapshot.ProviderOrderID == "" || snapshot.RealizationState != realization || snapshot.Quality != quality {
		t.Fatalf("snapshot realization contract changed after readback: %+v", snapshot)
	}
	if len(snapshot.Flags) != 1 || snapshot.Flags[0] != flag {
		t.Fatalf("snapshot flags changed after readback: %+v", snapshot.Flags)
	}
	if snapshot.ContributionAmount != nil || snapshot.MarginPercent != nil {
		t.Fatalf("missing monetary values defaulted during readback: contribution=%v margin=%v", snapshot.ContributionAmount, snapshot.MarginPercent)
	}
}
