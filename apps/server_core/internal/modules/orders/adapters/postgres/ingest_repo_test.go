//go:build integration

package postgres

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

// newIngestRepositoryForTest mirrors newOrderRepositoryForTest (order_repo_test.go) but also
// cleans up order_shipments, which that helper's fixture set never touches.
func newIngestRepositoryForTest(t *testing.T, ctx context.Context) (*OrderRepository, string, string, func()) {
	t.Helper()
	pool, _ := testpostgres.OpenPool(t, "tenant_harness_orders")
	tenantID := fmt.Sprintf("f02-ingest-%d", time.Now().UnixNano())
	repo := NewOrderRepository(pool, tenantID)
	installationID := tenantID + "-installation"
	return repo, tenantID, installationID, func() {
		if _, err := pool.Exec(ctx, `DELETE FROM order_shipments WHERE tenant_id = $1`, tenantID); err != nil {
			t.Errorf("cleanup order_shipments: %v", err)
		}
		if _, err := pool.Exec(ctx, `DELETE FROM orders_marketplace_order_items WHERE tenant_id = $1`, tenantID); err != nil {
			t.Errorf("cleanup order items: %v", err)
		}
		if _, err := pool.Exec(ctx, `DELETE FROM orders_marketplace_order_payments WHERE tenant_id = $1`, tenantID); err != nil {
			t.Errorf("cleanup order payments: %v", err)
		}
		if _, err := pool.Exec(ctx, `DELETE FROM orders_marketplace_orders WHERE tenant_id = $1`, tenantID); err != nil {
			t.Errorf("cleanup orders: %v", err)
		}
	}
}

func ingestFixtureOrder(installationID, orderID string, bucket ordersdomain.OrderBucket, providerShipmentID *string) ordersdomain.MarketplaceOrder {
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	name := "Comprador " + orderID
	return ordersdomain.MarketplaceOrder{
		InstallationID:     installationID,
		ProviderCode:       "mercado_livre",
		ProviderOrderID:    orderID,
		ProviderStatus:     "paid",
		ProviderUpdatedAt:  &now,
		FetchedAt:          now,
		RawProviderRef:     "/orders/" + orderID,
		CreatedAt:          now,
		UpdatedAt:          now,
		Bucket:             bucket,
		ProviderShipmentID: providerShipmentID,
		BuyerName:          &name,
		Items: []ordersdomain.MarketplaceOrderItem{{
			ProviderItemID: orderID + "-item",
			Quantity:       1,
			LinkQuality:    ordersdomain.LinkQualityMissing,
		}},
		Payments: []ordersdomain.MarketplaceOrderPayment{{
			ProviderPaymentID: orderID + "-payment",
			ProviderStatus:    "approved",
		}},
	}
}

func ingestFixtureShipment(providerOrderID, shipmentID string, fetchedAt time.Time, currency string) ordersdomain.OrderShipment {
	cost := 42.5
	return ordersdomain.OrderShipment{
		ProviderCode:       "mercado_livre",
		ProviderShipmentID: shipmentID,
		ProviderOrderID:    providerOrderID,
		Status:             "shipped",
		CostGross:          &cost,
		Currency:           &currency,
		FetchedAt:          fetchedAt,
		CreatedAt:          fetchedAt,
		UpdatedAt:          fetchedAt,
	}
}

type persistedOrderRow struct {
	bucket             *string
	providerShipmentID *string
	buyerName          *string
	status             string
}

func readOrderRowForTest(t *testing.T, ctx context.Context, repo *OrderRepository, installationID, orderID string) persistedOrderRow {
	t.Helper()
	var bucket, shipID, buyerName pgtype.Text
	var status string
	err := repo.pool.QueryRow(ctx, `
		SELECT provider_status, bucket, provider_shipment_id, buyer_name
		FROM orders_marketplace_orders
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
	`, repo.tenantID, installationID, orderID).Scan(&status, &bucket, &shipID, &buyerName)
	if err != nil {
		t.Fatalf("read order row %q: %v", orderID, err)
	}
	row := persistedOrderRow{status: status}
	if bucket.Valid {
		row.bucket = &bucket.String
	}
	if shipID.Valid {
		row.providerShipmentID = &shipID.String
	}
	if buyerName.Valid {
		row.buyerName = &buyerName.String
	}
	return row
}

func countShipmentRows(t *testing.T, ctx context.Context, repo *OrderRepository, shipmentID string) int {
	t.Helper()
	var count int
	if err := repo.pool.QueryRow(ctx, `
		SELECT count(*) FROM order_shipments WHERE tenant_id = $1 AND provider = $2 AND provider_shipment_id = $3
	`, repo.tenantID, "mercado_livre", shipmentID).Scan(&count); err != nil {
		t.Fatalf("count order_shipments rows: %v", err)
	}
	return count
}

func countRows(t *testing.T, ctx context.Context, repo *OrderRepository, table, installationID, orderID string) int {
	t.Helper()
	var count int
	query := fmt.Sprintf(`SELECT count(*) FROM %s WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3`, table)
	if err := repo.pool.QueryRow(ctx, query, repo.tenantID, installationID, orderID).Scan(&count); err != nil {
		t.Fatalf("count %s rows: %v", table, err)
	}
	return count
}

// TestOrderRepositoryIngestOrderPersistsHeaderItemsPaymentsAndShipment is scenario (a): the
// end-to-end fixture producing orders + order_shipments + the caller-supplied bucket, for a
// shipped/delivered-style arm (shipment present, BucketEnviado) — bucket DERIVATION itself is
// unit-tested DB-independently in orders/application/ingest_service_test.go; this test only
// proves the value round-trips through Postgres untouched.
func TestOrderRepositoryIngestOrderPersistsHeaderItemsPaymentsAndShipment(t *testing.T) {
	ctx := context.Background()
	repo, _, installationID, cleanup := newIngestRepositoryForTest(t, ctx)
	defer cleanup()

	shipID := "ship-" + installationID
	orderID := installationID + "-order-1"
	order := ingestFixtureOrder(installationID, orderID, ordersdomain.BucketEnviado, &shipID)
	shipment := ingestFixtureShipment(orderID, shipID, order.FetchedAt, "BRL")

	if err := repo.IngestOrder(ctx, order, &shipment); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}

	row := readOrderRowForTest(t, ctx, repo, installationID, orderID)
	if row.bucket == nil || ordersdomain.OrderBucket(*row.bucket) != ordersdomain.BucketEnviado {
		t.Fatalf("persisted bucket = %v, want %q", row.bucket, ordersdomain.BucketEnviado)
	}
	if row.providerShipmentID == nil || *row.providerShipmentID != shipID {
		t.Fatalf("persisted provider_shipment_id = %v, want %q", row.providerShipmentID, shipID)
	}
	if row.buyerName == nil || *row.buyerName != *order.BuyerName {
		t.Fatalf("persisted buyer_name = %v, want %q", row.buyerName, *order.BuyerName)
	}
	if got := countRows(t, ctx, repo, "orders_marketplace_order_items", installationID, orderID); got != 1 {
		t.Fatalf("item rows = %d, want 1", got)
	}
	if got := countRows(t, ctx, repo, "orders_marketplace_order_payments", installationID, orderID); got != 1 {
		t.Fatalf("payment rows = %d, want 1", got)
	}
	if got := countShipmentRows(t, ctx, repo, shipID); got != 1 {
		t.Fatalf("order_shipments rows = %d, want 1", got)
	}
}

// TestOrderRepositoryIngestOrderPersistsOrderWithoutShipmentRow is scenario (d): a nil shipment
// (F-01's honest-absence on 404/410) must still persist the order header — with
// provider_shipment_id set for a future retry — and write zero order_shipments rows.
func TestOrderRepositoryIngestOrderPersistsOrderWithoutShipmentRow(t *testing.T) {
	ctx := context.Background()
	repo, _, installationID, cleanup := newIngestRepositoryForTest(t, ctx)
	defer cleanup()

	shipID := "ship-absent-" + installationID
	orderID := installationID + "-order-2"
	order := ingestFixtureOrder(installationID, orderID, ordersdomain.BucketFaturar, &shipID)

	if err := repo.IngestOrder(ctx, order, nil); err != nil {
		t.Fatalf("IngestOrder() error = %v", err)
	}

	row := readOrderRowForTest(t, ctx, repo, installationID, orderID)
	if row.providerShipmentID == nil || *row.providerShipmentID != shipID {
		t.Fatalf("persisted provider_shipment_id = %v, want %q retained for future retry", row.providerShipmentID, shipID)
	}
	if got := countShipmentRows(t, ctx, repo, shipID); got != 0 {
		t.Fatalf("order_shipments rows = %d, want 0 on honest-absence", got)
	}
}

// TestOrderRepositoryIngestOrderIsIdempotentOnReplay is scenario (b): re-ingesting the identical
// order+shipment must not change persisted row counts.
func TestOrderRepositoryIngestOrderIsIdempotentOnReplay(t *testing.T) {
	ctx := context.Background()
	repo, _, installationID, cleanup := newIngestRepositoryForTest(t, ctx)
	defer cleanup()

	shipID := "ship-replay-" + installationID
	orderID := installationID + "-order-3"
	order := ingestFixtureOrder(installationID, orderID, ordersdomain.BucketEnviado, &shipID)
	shipment := ingestFixtureShipment(orderID, shipID, order.FetchedAt, "BRL")

	if err := repo.IngestOrder(ctx, order, &shipment); err != nil {
		t.Fatalf("IngestOrder() first call error = %v", err)
	}
	firstOrders := countRows(t, ctx, repo, "orders_marketplace_orders", installationID, orderID)
	firstItems := countRows(t, ctx, repo, "orders_marketplace_order_items", installationID, orderID)
	firstPayments := countRows(t, ctx, repo, "orders_marketplace_order_payments", installationID, orderID)
	firstShipments := countShipmentRows(t, ctx, repo, shipID)

	// Replay: same snapshot, freshness watermarks unchanged (same fetched/updated timestamps).
	if err := repo.IngestOrder(ctx, order, &shipment); err != nil {
		t.Fatalf("IngestOrder() replay error = %v", err)
	}

	if got := countRows(t, ctx, repo, "orders_marketplace_orders", installationID, orderID); got != firstOrders {
		t.Fatalf("orders row count after replay = %d, want unchanged %d", got, firstOrders)
	}
	if got := countRows(t, ctx, repo, "orders_marketplace_order_items", installationID, orderID); got != firstItems {
		t.Fatalf("items row count after replay = %d, want unchanged %d", got, firstItems)
	}
	if got := countRows(t, ctx, repo, "orders_marketplace_order_payments", installationID, orderID); got != firstPayments {
		t.Fatalf("payments row count after replay = %d, want unchanged %d", got, firstPayments)
	}
	if got := countShipmentRows(t, ctx, repo, shipID); got != firstShipments {
		t.Fatalf("order_shipments row count after replay = %d, want unchanged %d", got, firstShipments)
	}
}

// TestOrderRepositoryIngestOrderRollsBackWholeTransactionWhenShipmentWriteFails is the must-fail
// atomicity scenario: IngestOrder must commit header+items+payments+shipment as ONE transaction.
// currency is char(3) (0088) — a 4-char value forces upsertOrderShipment to fail AFTER
// upsertOrder has already run inside the same tx. If IngestOrder ever split into two
// transactions (header committed independently of the shipment write), this test starts
// failing: the header would show the SECOND call's "shipped" status even though the shipment
// write failed, instead of staying on the first call's "paid" status.
func TestOrderRepositoryIngestOrderRollsBackWholeTransactionWhenShipmentWriteFails(t *testing.T) {
	ctx := context.Background()
	repo, _, installationID, cleanup := newIngestRepositoryForTest(t, ctx)
	defer cleanup()

	shipID := "ship-atomic-" + installationID
	orderID := installationID + "-order-4"
	baseline := ingestFixtureOrder(installationID, orderID, ordersdomain.BucketFaturar, nil)
	if err := repo.IngestOrder(ctx, baseline, nil); err != nil {
		t.Fatalf("IngestOrder() baseline error = %v", err)
	}

	mutated := ingestFixtureOrder(installationID, orderID, ordersdomain.BucketEnviado, &shipID)
	mutated.ProviderStatus = "shipped"
	newer := baseline.FetchedAt.Add(time.Hour)
	mutated.ProviderUpdatedAt = &newer
	mutated.FetchedAt = newer
	mutated.UpdatedAt = newer
	badShipment := ingestFixtureShipment(orderID, shipID, newer, "USDX") // char(3) violation

	err := repo.IngestOrder(ctx, mutated, &badShipment)
	if err == nil {
		t.Fatal("IngestOrder() error = nil, want the char(3) currency violation to fail the call")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "value too long") {
		t.Fatalf("IngestOrder() error = %v, want a data-length violation naming the currency column", err)
	}

	row := readOrderRowForTest(t, ctx, repo, installationID, orderID)
	if row.status != baseline.ProviderStatus {
		t.Fatalf("persisted status after failed ingest = %q, want unchanged baseline %q (header write must have rolled back with the shipment failure)", row.status, baseline.ProviderStatus)
	}
	if row.bucket == nil || ordersdomain.OrderBucket(*row.bucket) != ordersdomain.BucketFaturar {
		t.Fatalf("persisted bucket after failed ingest = %v, want unchanged baseline %q", row.bucket, ordersdomain.BucketFaturar)
	}
	if got := countShipmentRows(t, ctx, repo, shipID); got != 0 {
		t.Fatalf("order_shipments rows after failed ingest = %d, want 0 (nothing committed)", got)
	}
}

// TestReadModelExposesCurrencyAndFulfillment is Task 5's red/green test (P2 slice):
// currency and fulfillment were contract fields on ordersdomain.OrderReadModel with zero
// producer — null on every order, not because storage was empty but because nothing ever
// wrote them. currency now has a real column (migration 0096) fed by the order's own
// provider currency_id. fulfillment has NO column of its own on orders_marketplace_orders —
// it is derived read-side from order_shipments.logistic_type (Task 4's producer), so the
// same fact is never stored twice. Both read-model projections (ListOrders and GetOrder)
// must agree, since scanReadModel is shared between them — this test exercises both.
func TestReadModelExposesCurrencyAndFulfillment(t *testing.T) {
	ctx := context.Background()
	repo, tenantID, installationID, cleanup := newIngestRepositoryForTest(t, ctx)
	defer cleanup()
	readRepo := NewOrderReadRepository(repo.pool, tenantID)

	brl := "BRL"
	fulfillmentType := "fulfillment"
	dropOffType := "drop_off"
	// ListOrders (unlike GetOrder) filters on provider_created_at IS NOT NULL —
	// ingestFixtureOrder leaves it nil, so both fixtures set it explicitly here to stay
	// visible to the ListOrders projection this test also exercises.
	providerCreatedAt := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)

	fullShipID := "ship-full-" + installationID
	fullOrderID := installationID + "-full"
	full := ingestFixtureOrder(installationID, fullOrderID, ordersdomain.BucketEnviado, &fullShipID)
	full.Currency = &brl
	full.ProviderCreatedAt = &providerCreatedAt
	fullShipment := ingestFixtureShipment(fullOrderID, fullShipID, full.FetchedAt, "BRL")
	fullShipment.LogisticType = &fulfillmentType
	if err := repo.IngestOrder(ctx, full, &fullShipment); err != nil {
		t.Fatalf("IngestOrder(full) error = %v", err)
	}

	dropShipID := "ship-drop-" + installationID
	dropOrderID := installationID + "-drop"
	drop := ingestFixtureOrder(installationID, dropOrderID, ordersdomain.BucketEnviado, &dropShipID)
	drop.Currency = &brl
	drop.ProviderCreatedAt = &providerCreatedAt
	dropShipment := ingestFixtureShipment(dropOrderID, dropShipID, drop.FetchedAt, "BRL")
	dropShipment.LogisticType = &dropOffType
	if err := repo.IngestOrder(ctx, drop, &dropShipment); err != nil {
		t.Fatalf("IngestOrder(drop) error = %v", err)
	}

	// GetOrder projection.
	gotFull, err := readRepo.GetOrder(ctx, installationID, fullOrderID)
	if err != nil {
		t.Fatalf("GetOrder(full) error = %v", err)
	}
	if gotFull.Currency == nil || *gotFull.Currency != "BRL" {
		t.Fatalf("GetOrder Currency = %v, want BRL — the field continues without a producer", gotFull.Currency)
	}
	if gotFull.Fulfillment == nil || *gotFull.Fulfillment != "fulfillment" {
		t.Fatalf("GetOrder Fulfillment = %v, want fulfillment", gotFull.Fulfillment)
	}

	// A non-fulfillment shipment modality must not become "fulfillment" (lazy always-nil-
	// unless-fulfillment implementation) and must not become nil either: drop_off is a known
	// fact, not an absence of one.
	gotDrop, err := readRepo.GetOrder(ctx, installationID, dropOrderID)
	if err != nil {
		t.Fatalf("GetOrder(drop) error = %v", err)
	}
	if gotDrop.Currency == nil || *gotDrop.Currency != "BRL" {
		t.Fatalf("GetOrder(drop) Currency = %v, want BRL", gotDrop.Currency)
	}
	if gotDrop.Fulfillment == nil || *gotDrop.Fulfillment != "drop_off" {
		t.Fatalf("GetOrder Fulfillment of drop_off order = %v, want drop_off", gotDrop.Fulfillment)
	}

	// ListOrders projection — same two facts, different query, same scanReadModel.
	page, err := readRepo.ListOrders(ctx, ports.OrderListQuery{InstallationID: installationID, Limit: 10})
	if err != nil {
		t.Fatalf("ListOrders error = %v", err)
	}
	byOrderID := make(map[string]ordersdomain.OrderReadModel, len(page.Items))
	for _, item := range page.Items {
		byOrderID[item.ProviderOrderID] = item
	}
	listFull, ok := byOrderID[fullOrderID]
	if !ok {
		t.Fatalf("ListOrders did not return %q", fullOrderID)
	}
	if listFull.Currency == nil || *listFull.Currency != "BRL" {
		t.Fatalf("ListOrders Currency = %v, want BRL", listFull.Currency)
	}
	if listFull.Fulfillment == nil || *listFull.Fulfillment != "fulfillment" {
		t.Fatalf("ListOrders Fulfillment = %v, want fulfillment", listFull.Fulfillment)
	}
	listDrop, ok := byOrderID[dropOrderID]
	if !ok {
		t.Fatalf("ListOrders did not return %q", dropOrderID)
	}
	if listDrop.Fulfillment == nil || *listDrop.Fulfillment != "drop_off" {
		t.Fatalf("ListOrders Fulfillment of drop_off order = %v, want drop_off", listDrop.Fulfillment)
	}
}
