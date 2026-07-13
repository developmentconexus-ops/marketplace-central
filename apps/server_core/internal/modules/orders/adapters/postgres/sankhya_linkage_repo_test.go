//go:build integration

package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestSankhyaLinkageRepositoryAppendsAtomicallyAndRetriesSemantically(t *testing.T) {
	ctx := context.Background()
	repo, pool, items := newLinkageRepositoryForTest(t, ctx, "atomic", 2)
	linkage := linkageForItems(repo.tenantID, "installation-atomic", "order-atomic", items, "event-atomic", "key-atomic", 313001)

	stored, err := repo.AppendConfirmation(ctx, linkage)
	if err != nil {
		t.Fatalf("append confirmation: %v", err)
	}
	if !stored.SemanticallyEqual(linkage) {
		t.Fatalf("stored linkage = %#v, want semantic match", stored)
	}

	retry := linkage
	retry.Audit.EventID = "event-retry"
	retry.Audit.RecordedAt = retry.Audit.RecordedAt.Add(time.Minute)
	retry.Lines = []ordersdomain.SankhyaLineMapping{linkage.Lines[1], linkage.Lines[0]}
	retried, err := repo.AppendConfirmation(ctx, retry)
	if err != nil {
		t.Fatalf("semantic retry: %v", err)
	}
	if retried.Audit.EventID != linkage.Audit.EventID {
		t.Fatalf("retry event = %q, want persisted event %q", retried.Audit.EventID, linkage.Audit.EventID)
	}

	var events, lines int
	if err := pool.QueryRow(ctx, `
		SELECT
		  (SELECT count(*) FROM orders_sankhya_linkage_events WHERE tenant_id=$1 AND installation_id=$2 AND provider_order_id=$3),
		  (SELECT count(*) FROM orders_sankhya_linkage_lines WHERE tenant_id=$1 AND installation_id=$2 AND provider_order_id=$3)
	`, repo.tenantID, linkage.Scope.InstallationID, linkage.Scope.ProviderOrderID).Scan(&events, &lines); err != nil {
		t.Fatalf("count ledger rows: %v", err)
	}
	if events != 1 || lines != 2 {
		t.Fatalf("ledger counts = (%d events, %d lines), want (1, 2)", events, lines)
	}
}

func TestSankhyaLinkageRepositoryConcurrentRetriesCreateOneMapping(t *testing.T) {
	ctx := context.Background()
	repo, pool, items := newLinkageRepositoryForTest(t, ctx, "concurrent", 1)
	base := linkageForItems(repo.tenantID, "installation-concurrent", "order-concurrent", items, "event-concurrent", "key-concurrent", 313002)

	const workers = 6
	start := make(chan struct{})
	errorsByWorker := make(chan error, workers)
	var wait sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wait.Add(1)
		go func(worker int) {
			defer wait.Done()
			<-start
			attempt := base
			attempt.Audit.EventID = fmt.Sprintf("event-concurrent-%d", worker)
			attempt.Audit.RecordedAt = attempt.Audit.RecordedAt.Add(time.Duration(worker) * time.Second)
			_, err := repo.AppendConfirmation(ctx, attempt)
			errorsByWorker <- err
		}(worker)
	}
	close(start)
	wait.Wait()
	close(errorsByWorker)
	for err := range errorsByWorker {
		if err != nil {
			t.Fatalf("concurrent retry: %v", err)
		}
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM orders_sankhya_linkage_events WHERE tenant_id=$1 AND idempotency_key=$2`, repo.tenantID, base.Audit.IdempotencyKey).Scan(&count); err != nil {
		t.Fatalf("count concurrent events: %v", err)
	}
	if count != 1 {
		t.Fatalf("concurrent event count = %d, want 1", count)
	}
}

func TestOrderRepositoryRefreshesConfirmedLinkedOrderInPlace(t *testing.T) {
	ctx := context.Background()
	linkageRepo, pool, items := newLinkageRepositoryForTest(t, ctx, "linked-refresh", 2)
	linkage := linkageForItems(linkageRepo.tenantID, "installation-linked-refresh", "order-linked-refresh", items, "event-linked-refresh", "key-linked-refresh", 313010)
	if _, err := linkageRepo.AppendConfirmation(ctx, linkage); err != nil {
		t.Fatalf("append linked refresh fixture: %v", err)
	}

	var originalCreatedAt time.Time
	if err := pool.QueryRow(ctx, `
		SELECT created_at FROM orders_marketplace_order_items
		WHERE tenant_id=$1 AND installation_id=$2 AND provider_order_id=$3 AND mpc_line_id=$4
	`, linkageRepo.tenantID, linkage.Scope.InstallationID, linkage.Scope.ProviderOrderID, items[0].MPCLineID).Scan(&originalCreatedAt); err != nil {
		t.Fatalf("read original item creation time: %v", err)
	}
	orderRepo := NewOrderRepository(pool, linkageRepo.tenantID)
	orders, err := orderRepo.ListOrders(ctx, linkage.Scope.InstallationID, 10)
	if err != nil || len(orders) != 1 {
		t.Fatalf("load linked order: count=%d err=%v", len(orders), err)
	}
	originalIDs := lineIDsByProviderIdentity(orders[0].Items)
	refreshAt := time.Date(2026, 7, 13, 14, 0, 0, 0, time.UTC)
	refresh := orders[0]
	refresh.ProviderUpdatedAt = &refreshAt
	refresh.FetchedAt = refreshAt
	refresh.CreatedAt = refreshAt
	refresh.UpdatedAt = refreshAt
	refresh.Items = []ordersdomain.MarketplaceOrderItem{orders[0].Items[1], orders[0].Items[0]}
	refresh.Items[0].Quantity = 7
	refresh.Items[1].Quantity = 9
	for index := range refresh.Items {
		refresh.Items[index].MPCLineID = ""
		refresh.Items[index].ReconciliationState = ""
	}
	imported, skipped, err := orderRepo.UpsertOrders(ctx, []ordersdomain.MarketplaceOrder{refresh})
	if err != nil || imported != 1 || skipped != 0 {
		t.Fatalf("refresh linked order: imported=%d skipped=%d err=%v", imported, skipped, err)
	}
	refreshed, err := orderRepo.ListOrders(ctx, linkage.Scope.InstallationID, 10)
	if err != nil || len(refreshed) != 1 {
		t.Fatalf("load refreshed order: count=%d err=%v", len(refreshed), err)
	}
	if got := lineIDsByProviderIdentity(refreshed[0].Items); !sameLineIDMap(got, originalIDs) {
		t.Fatalf("linked refresh IDs = %#v, want %#v", got, originalIDs)
	}
	storedLinkage, found, err := linkageRepo.LoadCurrent(ctx, linkage.Scope)
	if err != nil || !found || !storedLinkage.SemanticallyEqual(linkage) {
		t.Fatalf("load linkage after refresh: found=%v linkage=%#v err=%v", found, storedLinkage, err)
	}
	var refreshedCreatedAt time.Time
	if err := pool.QueryRow(ctx, `
		SELECT created_at FROM orders_marketplace_order_items
		WHERE tenant_id=$1 AND installation_id=$2 AND provider_order_id=$3 AND mpc_line_id=$4
	`, linkageRepo.tenantID, linkage.Scope.InstallationID, linkage.Scope.ProviderOrderID, items[0].MPCLineID).Scan(&refreshedCreatedAt); err != nil {
		t.Fatalf("read refreshed item creation time: %v", err)
	}
	if !refreshedCreatedAt.Equal(originalCreatedAt) {
		t.Fatalf("linked row was recreated: created_at=%v want %v", refreshedCreatedAt, originalCreatedAt)
	}
}

func TestOrderRepositoryRemovedLinkedLineFailsWithIdentityConflict(t *testing.T) {
	ctx := context.Background()
	linkageRepo, pool, items := newLinkageRepositoryForTest(t, ctx, "linked-removal", 2)
	linkage := linkageForItems(linkageRepo.tenantID, "installation-linked-removal", "order-linked-removal", items[:1], "event-linked-removal", "key-linked-removal", 313011)
	if _, err := linkageRepo.AppendConfirmation(ctx, linkage); err != nil {
		t.Fatalf("append linked removal fixture: %v", err)
	}
	orderRepo := NewOrderRepository(pool, linkageRepo.tenantID)
	orders, err := orderRepo.ListOrders(ctx, linkage.Scope.InstallationID, 10)
	if err != nil || len(orders) != 1 {
		t.Fatalf("load linked removal order: count=%d err=%v", len(orders), err)
	}
	refreshAt := time.Date(2026, 7, 13, 14, 0, 0, 0, time.UTC)
	refresh := orders[0]
	refresh.ProviderStatus = "changed"
	refresh.ProviderUpdatedAt = &refreshAt
	refresh.FetchedAt = refreshAt
	refresh.UpdatedAt = refreshAt
	refresh.Items = []ordersdomain.MarketplaceOrderItem{orders[0].Items[1]}
	refresh.Items[0].MPCLineID = ""
	refresh.Items[0].ReconciliationState = ""
	if _, _, err := orderRepo.UpsertOrders(ctx, []ordersdomain.MarketplaceOrder{refresh}); !errors.Is(err, ordersdomain.ErrOrderLineIdentityConflict) {
		t.Fatalf("removed linked line error = %v, want identity conflict", err)
	}
	stored, err := orderRepo.ListOrders(ctx, linkage.Scope.InstallationID, 10)
	if err != nil || len(stored) != 1 {
		t.Fatalf("load order after rejected removal: count=%d err=%v", len(stored), err)
	}
	if stored[0].ProviderStatus == "changed" || len(stored[0].Items) != 2 {
		t.Fatalf("rejected refresh partially persisted: status=%q items=%d", stored[0].ProviderStatus, len(stored[0].Items))
	}
}

func TestSankhyaLinkageRepositoryConflictsFailClosedWithoutOverwrite(t *testing.T) {
	ctx := context.Background()
	repo, pool, items := newLinkageRepositoryForTest(t, ctx, "conflict", 1)
	base := linkageForItems(repo.tenantID, "installation-conflict", "order-conflict", items, "event-conflict", "key-conflict", 313003)
	if _, err := repo.AppendConfirmation(ctx, base); err != nil {
		t.Fatalf("append base mapping: %v", err)
	}

	conflict := base
	conflict.Lines = append([]ordersdomain.SankhyaLineMapping(nil), base.Lines...)
	conflict.Audit.EventID = "event-conflicting-retry"
	conflict.Header.DocumentID = 313099
	conflict.Lines[0].Origin.DocumentID = 313099
	conflict.Lines[0].Origin.LineNumber = 9
	if _, err := repo.AppendConfirmation(ctx, conflict); !errors.Is(err, ordersdomain.ErrSankhyaLinkageConflict) {
		t.Fatalf("conflicting retry error = %v, want conflict", err)
	}
	stored, found, err := repo.LoadCurrent(ctx, base.Scope)
	if err != nil || !found {
		t.Fatalf("load base mapping: found=%v err=%v", found, err)
	}
	if stored.Header != base.Header || stored.Lines[0] != base.Lines[0] {
		t.Fatalf("base mapping was overwritten: %#v", stored)
	}
	otherItems := persistOrderItems(t, ctx, pool, repo.tenantID, "installation-conflict", "order-key-conflict", 1)
	keyConflict := linkageForItems(repo.tenantID, "installation-conflict", "order-key-conflict", otherItems, "event-key-conflict", base.Audit.IdempotencyKey, 313100)
	if _, err := repo.AppendConfirmation(ctx, keyConflict); !errors.Is(err, ordersdomain.ErrSankhyaLinkageConflict) {
		t.Fatalf("external-key/idempotency conflict error = %v, want conflict", err)
	}
	assertNoLinkageEvent(t, ctx, pool, repo.tenantID, keyConflict.Audit.EventID)

	if _, err := pool.Exec(ctx, `UPDATE orders_sankhya_linkage_events SET reason='changed' WHERE tenant_id=$1 AND event_id=$2`, repo.tenantID, base.Audit.EventID); err == nil {
		t.Fatal("append-only event update unexpectedly succeeded")
	}
	if _, err := pool.Exec(ctx, `DELETE FROM orders_sankhya_linkage_lines WHERE tenant_id=$1 AND event_id=$2`, repo.tenantID, base.Audit.EventID); err == nil {
		t.Fatal("append-only line delete unexpectedly succeeded")
	}
}

func TestSankhyaLinkageRepositoryRejectsNonLinkableOrWrongOrderLinesAtomically(t *testing.T) {
	ctx := context.Background()
	repo, pool, items := newLinkageRepositoryForTest(t, ctx, "nonlinkable", 2)
	if _, err := pool.Exec(ctx, `
		UPDATE orders_marketplace_order_items SET reconciliation_state='ambiguous'
		WHERE tenant_id=$1 AND installation_id=$2 AND provider_order_id=$3 AND mpc_line_id=$4
	`, repo.tenantID, "installation-nonlinkable", "order-nonlinkable", items[0].MPCLineID); err != nil {
		t.Fatalf("mark ambiguous fixture: %v", err)
	}
	linkage := linkageForItems(repo.tenantID, "installation-nonlinkable", "order-nonlinkable", items[:1], "event-nonlinkable", "key-nonlinkable", 313004)
	if _, err := repo.AppendConfirmation(ctx, linkage); !errors.Is(err, ordersdomain.ErrSankhyaLinkageConflict) {
		t.Fatalf("ambiguous line error = %v, want conflict", err)
	}
	assertNoLinkageEvent(t, ctx, pool, repo.tenantID, linkage.Audit.EventID)

	otherOrder := persistOrderItems(t, ctx, pool, repo.tenantID, "installation-nonlinkable", "order-other", 1)
	wrongOrder := linkageForItems(repo.tenantID, "installation-nonlinkable", "order-nonlinkable", otherOrder, "event-wrong-order", "key-wrong-order", 313005)
	if _, err := repo.AppendConfirmation(ctx, wrongOrder); !errors.Is(err, ordersdomain.ErrSankhyaLinkageNotFound) {
		t.Fatalf("wrong-order line error = %v, want not found", err)
	}
	assertNoLinkageEvent(t, ctx, pool, repo.tenantID, wrongOrder.Audit.EventID)
}

func TestSankhyaLinkageRepositoryScopesOriginsByTenantAndInstallation(t *testing.T) {
	ctx := context.Background()
	repo, pool, items := newLinkageRepositoryForTest(t, ctx, "origin-one", 1)
	first := linkageForItems(repo.tenantID, "installation-origin-one", "order-origin-one", items, "event-origin-one", "key-origin-one", 313006)
	if _, err := repo.AppendConfirmation(ctx, first); err != nil {
		t.Fatalf("append first origin: %v", err)
	}

	secondItems := persistOrderItems(t, ctx, pool, repo.tenantID, "installation-origin-one", "order-origin-two", 1)
	second := linkageForItems(repo.tenantID, "installation-origin-one", "order-origin-two", secondItems, "event-origin-two", "key-origin-two", 313006)
	if _, err := repo.AppendConfirmation(ctx, second); !errors.Is(err, ordersdomain.ErrSankhyaLinkageConflict) {
		t.Fatalf("reused origin error = %v, want conflict", err)
	}
	assertNoLinkageEvent(t, ctx, pool, repo.tenantID, second.Audit.EventID)

	otherInstallationItems := persistOrderItems(t, ctx, pool, repo.tenantID, "installation-origin-other", "order-origin-three", 1)
	otherInstallation := linkageForItems(repo.tenantID, "installation-origin-other", "order-origin-three", otherInstallationItems, "event-origin-three", "key-origin-three", 313006)
	if _, err := repo.AppendConfirmation(ctx, otherInstallation); err != nil {
		t.Fatalf("same origin in another installation should be scoped independently: %v", err)
	}

	wrongTenant := first
	wrongTenant.Scope.TenantID = "other-tenant"
	if _, err := repo.AppendConfirmation(ctx, wrongTenant); !errors.Is(err, ordersdomain.ErrInvalidSankhyaLinkage) {
		t.Fatalf("wrong tenant error = %v, want invalid linkage", err)
	}
}

func newLinkageRepositoryForTest(t *testing.T, ctx context.Context, suffix string, itemCount int) (*SankhyaLinkageRepository, *pgxpool.Pool, []ordersdomain.MarketplaceOrderItem) {
	t.Helper()
	tenantID := fmt.Sprintf("f09-%s-%d", suffix, time.Now().UnixNano())
	pool, _ := testpostgres.OpenPool(t, tenantID)
	installationID := "installation-" + suffix
	orderID := "order-" + suffix
	items := persistOrderItems(t, ctx, pool, tenantID, installationID, orderID, itemCount)
	return NewSankhyaLinkageRepository(pool, tenantID), pool, items
}

func persistOrderItems(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenantID, installationID, orderID string, count int) []ordersdomain.MarketplaceOrderItem {
	t.Helper()
	now := time.Date(2026, 7, 13, 12, 0, 0, 0, time.UTC)
	items := make([]ordersdomain.MarketplaceOrderItem, 0, count)
	for index := 0; index < count; index++ {
		items = append(items, ordersdomain.MarketplaceOrderItem{
			ProviderItemID:      fmt.Sprintf("MLB-%s-%d", orderID, index),
			ProviderVariationID: fmt.Sprintf("variation-%d", index),
			SellerSKU:           fmt.Sprintf("sku-%d", index),
			Quantity:            1,
			LinkQuality:         ordersdomain.LinkQualityMissing,
		})
	}
	orderRepo := NewOrderRepository(pool, tenantID)
	imported, skipped, err := orderRepo.UpsertOrders(ctx, []ordersdomain.MarketplaceOrder{{
		InstallationID:    installationID,
		ProviderCode:      "mercado_livre",
		ProviderOrderID:   orderID,
		ProviderUpdatedAt: &now,
		FetchedAt:         now,
		CreatedAt:         now,
		UpdatedAt:         now,
		Items:             items,
	}})
	if err != nil || imported != 1 || skipped != 0 {
		t.Fatalf("persist order fixture: imported=%d skipped=%d err=%v", imported, skipped, err)
	}
	stored, err := orderRepo.ListOrders(ctx, installationID, 100)
	if err != nil {
		t.Fatalf("load order fixture: %v", err)
	}
	for _, order := range stored {
		if order.ProviderOrderID == orderID {
			return order.Items
		}
	}
	t.Fatalf("persisted order fixture %q not found", orderID)
	return nil
}

func linkageForItems(tenantID, installationID, orderID string, items []ordersdomain.MarketplaceOrderItem, eventID, idempotencyKey string, documentID int64) ordersdomain.SankhyaLinkage {
	now := time.Date(2026, 7, 13, 13, 0, 0, 0, time.UTC)
	lines := make([]ordersdomain.SankhyaLineMapping, 0, len(items))
	for index, item := range items {
		lines = append(lines, ordersdomain.SankhyaLineMapping{
			MPCLineID: item.MPCLineID,
			Origin:    ordersdomain.InternalDocumentLineIdentity{DocumentID: documentID, LineNumber: index + 1},
		})
	}
	return ordersdomain.SankhyaLinkage{
		Scope:            ordersdomain.LinkageScope{TenantID: tenantID, InstallationID: installationID, ProviderOrderID: orderID},
		ExternalOrderKey: ordersdomain.ExternalOrderKeyFor(ordersdomain.LinkageScope{InstallationID: installationID, ProviderOrderID: orderID}),
		Header:           ordersdomain.InternalDocumentIdentity{DocumentID: documentID},
		Lines:            lines,
		Audit: ordersdomain.LinkageAudit{
			EventID:               eventID,
			EventType:             ordersdomain.LinkageEventConfirmed,
			ActorType:             "operator",
			ActorID:               "actor-1",
			Reason:                "operator confirmed exact internal lines",
			IdempotencyKey:        idempotencyKey,
			SourceAt:              now,
			RecordedAt:            now,
			ConfigurationRevision: "config-1",
			EvidenceState:         ordersdomain.LinkageEvidenceExact,
			EvidenceReference:     "safe-reference",
		},
	}
}

func assertNoLinkageEvent(t *testing.T, ctx context.Context, pool *pgxpool.Pool, tenantID, eventID string) {
	t.Helper()
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM orders_sankhya_linkage_events WHERE tenant_id=$1 AND event_id=$2`, tenantID, eventID).Scan(&count); err != nil {
		t.Fatalf("count rolled-back event: %v", err)
	}
	if count != 0 {
		t.Fatalf("rolled-back event count = %d, want 0", count)
	}
}
