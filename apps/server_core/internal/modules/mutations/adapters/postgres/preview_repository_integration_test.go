//go:build integration

package postgres_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	mutationspostgres "marketplace-central/apps/server_core/internal/modules/mutations/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/mutations/domain"
	"marketplace-central/apps/server_core/internal/modules/mutations/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestReplacePreviewAtomicallyReplacesItemsAndSourceTime(t *testing.T) {
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_mutation_preview")
	tenant := "mutation-preview-" + time.Now().UTC().Format("150405.000000000")
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM mutation_protocols WHERE tenant_id=$1`, tenant)
	})
	repo := mutationspostgres.NewRepository(pool, tenant)
	created, err := repo.CreateProtocol(ctx, ports.CreateProtocolInput{InstallationID: "inst-1", Type: domain.ProtocolTypePriceUpdate, Actor: "operator_supplied_unverified", Intent: json.RawMessage(`{"new_price":{"amount":"49.90"}}`), Selection: json.RawMessage(`{"mode":"explicit"}`), CreatedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	firstSource := time.Date(2026, 7, 16, 10, 0, 0, 0, time.UTC)
	if _, err := repo.ReplacePreview(ctx, created.ProtocolID, []ports.ReplaceItemInput{{ListingID: "inst-1~MLB1~-", Before: json.RawMessage(`{"price":{"amount":"89.00"}}`), After: created.Intent}}, firstSource, firstSource.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	secondSource := firstSource.Add(2 * time.Minute)
	if _, err := repo.ReplacePreview(ctx, created.ProtocolID, []ports.ReplaceItemInput{{ListingID: "inst-1~MLB2~-", Before: json.RawMessage(`{"price":{"amount":"90.00"}}`), After: created.Intent}}, secondSource, secondSource.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}

	var count int
	var listingID string
	if err := pool.QueryRow(ctx, `SELECT count(*),min(listing_id) FROM mutation_items WHERE tenant_id=$1 AND protocol_id=$2`, tenant, created.ProtocolID).Scan(&count, &listingID); err != nil {
		t.Fatal(err)
	}
	got, found, err := repo.GetProtocol(ctx, created.ProtocolID)
	if err != nil || !found {
		t.Fatalf("GetProtocol found=%v err=%v", found, err)
	}
	if count != 1 || listingID != "inst-1~MLB2~-" || got.SourceAsOf == nil || !got.SourceAsOf.Equal(secondSource) {
		t.Fatalf("count=%d listing=%q source=%v", count, listingID, got.SourceAsOf)
	}
}
