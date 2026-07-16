//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	orderspostgres "marketplace-central/apps/server_core/internal/modules/orders/adapters/postgres"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
	testpostgres "marketplace-central/apps/server_core/internal/testsupport/postgres"
)

func TestListOrderPageUsesTenantInstallationAndKeyset(t *testing.T) {
	ctx, repo, pool, tenant, other, installation := orderReadFixture(t, "scope")
	at := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	seedReadOrder(t, pool, tenant, installation, "a", "paid", &at)
	seedReadOrder(t, pool, other, installation, "leak", "paid", &at)
	seedReadOrder(t, pool, tenant, installation+"-other", "wrong-install", "paid", &at)

	first, err := repo.ListOrders(ctx, ports.OrderListQuery{InstallationID: installation, Limit: 1})
	if err != nil || len(first.Items) != 1 || first.Items[0].ProviderOrderID != "a" || first.NextCursor != nil {
		t.Fatalf("page=%+v err=%v", first, err)
	}
	got := first.Items[0]
	if got.BuyerNickname != nil || got.Currency != nil || got.Fulfillment != nil || got.NFState != nil || got.Total != nil {
		t.Fatalf("unsupported or absent facts defaulted: %+v", got)
	}
}

func TestListOrderPageNewestFirstWithoutDuplicates(t *testing.T) {
	ctx, repo, pool, tenant, _, installation := orderReadFixture(t, "keyset")
	newest := time.Date(2026, 7, 16, 13, 0, 0, 0, time.UTC)
	older := newest.Add(-time.Hour)
	for _, row := range []struct {
		id string
		at *time.Time
	}{{"c", &newest}, {"b", &newest}, {"a", &older}, {"null", nil}} {
		seedReadOrder(t, pool, tenant, installation, row.id, "paid", row.at)
	}
	query := ports.OrderListQuery{InstallationID: installation, Limit: 2}
	first, err := repo.ListOrders(ctx, query)
	if err != nil || len(first.Items) != 2 || first.Items[0].ProviderOrderID != "c" || first.Items[1].ProviderOrderID != "b" || first.NextCursor == nil {
		t.Fatalf("first=%+v err=%v", first, err)
	}
	query.Cursor = *first.NextCursor
	second, err := repo.ListOrders(ctx, query)
	if err != nil || len(second.Items) != 1 || second.Items[0].ProviderOrderID != "a" || second.NextCursor != nil {
		t.Fatalf("second=%+v err=%v", second, err)
	}
}

func TestOrderFiltersStatusDateAndQuery(t *testing.T) {
	ctx, repo, pool, tenant, _, installation := orderReadFixture(t, "filters")
	at := time.Date(2026, 7, 16, 12, 0, 0, 0, time.UTC)
	seedReadOrder(t, pool, tenant, installation, "prefix-1", "paid", &at)
	seedReadOrder(t, pool, tenant, installation, "other", "cancelled", &at)
	_, err := pool.Exec(ctx, `INSERT INTO orders_marketplace_order_items(tenant_id,installation_id,provider_order_id,line_no,provider_item_id,seller_sku,title,quantity,link_quality,created_at,updated_at,mpc_line_id,reconciliation_state) VALUES($1,$2,'other',1,'item','SKU-Needle','Title',1,'resolved',now(),now(),'mpl_11111111111111111111111111111111','stable')`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	from, to := at.Add(-time.Minute), at.Add(time.Minute)
	for _, query := range []ports.OrderListQuery{
		{InstallationID: installation, Limit: 10, Status: "paid"},
		{InstallationID: installation, Limit: 10, DateFrom: &from, DateTo: &to},
		{InstallationID: installation, Limit: 10, Q: "prefix"},
		{InstallationID: installation, Limit: 10, Q: "needle"},
	} {
		page, e := repo.ListOrders(ctx, query)
		if e != nil || len(page.Items) == 0 {
			t.Fatalf("query=%+v page=%+v err=%v", query, page, e)
		}
	}
}

func TestGetOrderCannotCrossTenant(t *testing.T) {
	ctx, repo, pool, tenant, other, installation := orderReadFixture(t, "cross")
	at := time.Now().UTC()
	seedReadOrder(t, pool, tenant, installation, "visible", "paid", &at)
	seedReadOrder(t, pool, other, installation, "secret", "paid", &at)
	_, err := pool.Exec(ctx, `INSERT INTO orders_marketplace_order_items(tenant_id,installation_id,provider_order_id,line_no,provider_item_id,seller_sku,title,quantity,link_quality,created_at,updated_at,mpc_line_id,reconciliation_state) VALUES($1,$2,'visible',1,'item','SKU','Title',2,'resolved',now(),now(),'mpl_22222222222222222222222222222222','stable')`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `INSERT INTO orders_marketplace_order_payments(tenant_id,installation_id,provider_order_id,line_no,provider_payment_id,transaction_amount,total_paid_amount,created_at,updated_at) VALUES($1,$2,'visible',1,'payment',40,42,now(),now())`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec(ctx, `INSERT INTO orders_sankhya_linkage_events(tenant_id,installation_id,event_id,event_type,provider_order_id,external_order_key,internal_document_id,actor_type,actor_id,reason,idempotency_key,source_at,recorded_at,configuration_revision,evidence_state) VALUES($1,$2,'evt','confirmed','visible','ml:v1:'||$2||':visible',1,'operator','test','test','idem',now(),now(),'rev','exact')`, tenant, installation)
	if err != nil {
		t.Fatal(err)
	}
	visible, err := repo.GetOrder(ctx, installation, "visible")
	if err != nil || len(visible.Items) != 1 || len(visible.Payments) != 1 || visible.NFState == nil || *visible.NFState != "linked" || visible.Total == nil || *visible.Total != 42 {
		t.Fatalf("visible=%+v err=%v", visible, err)
	}
	if visible.BuyerNickname != nil || visible.Currency != nil || visible.Fulfillment != nil {
		t.Fatalf("unsupported facts=%+v", visible)
	}
	_, err = repo.GetOrder(ctx, installation, "secret")
	if !errors.Is(err, ports.ErrOrderNotFound) {
		t.Fatalf("err=%v", err)
	}
}

func TestOrderNFStateRequiresExactEvidence(t *testing.T) {
	ctx, repo, pool, tenant, _, installation := orderReadFixture(t, "nf-state")
	at := time.Now().UTC()
	seedReadOrder(t, pool, tenant, installation, "unknown", "paid", &at)
	seedReadOrder(t, pool, tenant, installation, "exact", "paid", &at)
	seedLinkageEvent := func(orderID, eventID, evidenceState string) {
		t.Helper()
		_, err := pool.Exec(ctx, `INSERT INTO orders_sankhya_linkage_events(tenant_id,installation_id,event_id,event_type,provider_order_id,external_order_key,internal_document_id,actor_type,actor_id,reason,idempotency_key,source_at,recorded_at,configuration_revision,evidence_state) VALUES($1,$2,$4,'confirmed',$3,'ml:v1:'||$2||':'||$3,1,'operator','test','test',$4,now(),now(),'rev',$5)`, tenant, installation, orderID, eventID, evidenceState)
		if err != nil {
			t.Fatal(err)
		}
	}
	seedLinkageEvent("unknown", "evt-unknown", "unknown")
	seedLinkageEvent("exact", "evt-exact", "exact")

	page, err := repo.ListOrders(ctx, ports.OrderListQuery{InstallationID: installation, Limit: 10})
	if err != nil || len(page.Items) != 2 {
		t.Fatalf("page=%+v err=%v", page, err)
	}
	for _, item := range page.Items {
		switch item.ProviderOrderID {
		case "unknown":
			if item.NFState != nil {
				t.Fatalf("unknown evidence reported nf_state=%q", *item.NFState)
			}
		case "exact":
			if item.NFState == nil || *item.NFState != "linked" {
				t.Fatalf("exact evidence reported nf_state=%v", item.NFState)
			}
		}
	}

	unknownOrder, err := repo.GetOrder(ctx, installation, "unknown")
	if err != nil || unknownOrder.NFState != nil {
		t.Fatalf("unknown order=%+v err=%v", unknownOrder, err)
	}
	exactOrder, err := repo.GetOrder(ctx, installation, "exact")
	if err != nil || exactOrder.NFState == nil || *exactOrder.NFState != "linked" {
		t.Fatalf("exact order=%+v err=%v", exactOrder, err)
	}
}

func TestGetUnknownOrderReturnsNotFound(t *testing.T) {
	ctx, repo, _, _, _, installation := orderReadFixture(t, "missing")
	_, err := repo.GetOrder(ctx, installation, "missing")
	var typed *ports.OrderNotFoundError
	if !errors.As(err, &typed) || !errors.Is(err, ports.ErrOrderNotFound) {
		t.Fatalf("err=%#v", err)
	}
}

func orderReadFixture(t *testing.T, suffix string) (context.Context, *orderspostgres.OrderReadRepository, *pgxpool.Pool, string, string, string) {
	t.Helper()
	ctx := context.Background()
	pool, _ := testpostgres.OpenPool(t, "orders_read_"+suffix)
	token := fmt.Sprint(time.Now().UnixNano())
	tenant, other, installation := "tenant-a-"+token, "tenant-b-"+token, "inst-"+token
	t.Cleanup(func() {
		for _, table := range []string{"orders_sankhya_linkage_events", "orders_marketplace_order_items", "orders_marketplace_order_payments", "orders_marketplace_orders"} {
			_, _ = pool.Exec(context.Background(), "DELETE FROM "+table+" WHERE tenant_id = ANY($1)", []string{tenant, other})
		}
	})
	return ctx, orderspostgres.NewOrderReadRepository(pool, tenant), pool, tenant, other, installation
}

func seedReadOrder(t *testing.T, pool *pgxpool.Pool, tenant, installation, id, status string, at *time.Time) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `INSERT INTO orders_marketplace_orders(tenant_id,installation_id,provider_code,provider_order_id,provider_status,provider_created_at,fetched_at,created_at,updated_at) VALUES($1,$2,'mercado_livre',$3,$4,$5,now(),now(),now())`, tenant, installation, id, status, at)
	if err != nil {
		t.Fatal(err)
	}
}
