//go:build integration

package postgres_test

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	mutationspostgres "marketplace-central/apps/server_core/internal/modules/mutations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestRepositoryCreatesTenantScopedMonotonicProtocolIDs(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mutations")
	token := time.Now().UTC().Format("150405.000000000")
	tenantA, tenantB := "mutation-a-"+token, "mutation-b-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id = ANY($1)`, []string{tenantA, tenantB})
	})

	createdAt := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	create := func(tenant string) ports.Protocol {
		repo := mutationspostgres.NewRepository(pool, tenant)
		protocol, err := repo.CreateProtocol(ctx, ports.CreateProtocolInput{
			InstallationID: "installation-1",
			Type:           domain.ProtocolTypePriceUpdate,
			Actor:          "operator_supplied_unverified",
			Intent:         json.RawMessage(`{"new_price":{"amount":"49.90"}}`),
			Selection:      json.RawMessage(`{"mode":"explicit"}`),
			CreatedAt:      createdAt,
		})
		if err != nil {
			t.Errorf("CreateProtocol(%s): %v", tenant, err)
		}
		return protocol
	}

	var wg sync.WaitGroup
	created := make(chan ports.Protocol, 2)
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			created <- create(tenantA)
		}()
	}
	wg.Wait()
	close(created)
	ids := map[string]bool{}
	for protocol := range created {
		ids[protocol.ProtocolID] = true
	}
	if !ids["MP-000001"] || !ids["MP-000002"] || len(ids) != 2 {
		t.Fatalf("concurrent tenant A IDs = %#v, want MP-000001 and MP-000002", ids)
	}

	if got := create(tenantB).ProtocolID; got != "MP-000001" {
		t.Fatalf("tenant B first ID = %q, want MP-000001", got)
	}
}

func TestRepositoryDraftRoundTripAndCrossTenantNotFound(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mutations")
	token := time.Now().UTC().Format("150405.000000000")
	tenantA, tenantB := "mutation-roundtrip-a-"+token, "mutation-roundtrip-b-"+token
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id = ANY($1)`, []string{tenantA, tenantB})
	})

	intent := json.RawMessage(`{"new_price":{"amount":"49.90"}}`)
	selection := json.RawMessage(`{"mode":"explicit"}`)
	createdAt := time.Date(2026, 7, 16, 12, 34, 56, 123000000, time.UTC)
	repoA := mutationspostgres.NewRepository(pool, tenantA)
	created, err := repoA.CreateProtocol(ctx, ports.CreateProtocolInput{
		InstallationID: "installation-1", Type: domain.ProtocolTypePriceUpdate,
		Actor: "operator_supplied_unverified", Intent: intent, Selection: selection, CreatedAt: createdAt,
	})
	if err != nil {
		t.Fatalf("CreateProtocol(): %v", err)
	}

	got, found, err := repoA.GetProtocol(ctx, created.ProtocolID)
	if err != nil || !found {
		t.Fatalf("GetProtocol(): found=%v err=%v", found, err)
	}
	if got.ProtocolID != "MP-000001" || got.InstallationID != "installation-1" || got.Type != domain.ProtocolTypePriceUpdate || got.State != domain.ProtocolStateDraft || got.Actor != "operator_supplied_unverified" || got.SourceAsOf != nil || !got.CreatedAt.Equal(createdAt) {
		t.Fatalf("round trip protocol = %#v", got)
	}
	assertEquivalentJSON(t, "intent", got.Intent, intent)
	assertEquivalentJSON(t, "selection", got.Selection, selection)

	if _, found, err := mutationspostgres.NewRepository(pool, tenantB).GetProtocol(ctx, created.ProtocolID); err != nil || found {
		t.Fatalf("cross-tenant GetProtocol(): found=%v err=%v, want not found", found, err)
	}
}

func assertEquivalentJSON(t *testing.T, field string, got, want json.RawMessage) {
	t.Helper()
	var gotCompact, wantCompact bytes.Buffer
	if err := json.Compact(&gotCompact, got); err != nil {
		t.Fatalf("compact returned %s: %v", field, err)
	}
	if err := json.Compact(&wantCompact, want); err != nil {
		t.Fatalf("compact expected %s: %v", field, err)
	}
	if !bytes.Equal(gotCompact.Bytes(), wantCompact.Bytes()) {
		t.Fatalf("%s = %s, want byte-equivalent %s", field, got, want)
	}
}
