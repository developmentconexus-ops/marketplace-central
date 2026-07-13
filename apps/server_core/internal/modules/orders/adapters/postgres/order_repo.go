package postgres

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	ordersdomain "marketplace-central/apps/server_core/internal/modules/orders/domain"
	"marketplace-central/apps/server_core/internal/modules/orders/ports"
)

var _ ports.OrderStore = (*OrderRepository)(nil)

type OrderRepository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewOrderRepository(pool *pgxpool.Pool, tenantID string) *OrderRepository {
	return &OrderRepository{pool: pool, tenantID: tenantID}
}

func (r *OrderRepository) UpsertOrders(ctx context.Context, orders []ordersdomain.MarketplaceOrder) (int, int, error) {
	if r.pool == nil || len(orders) == 0 {
		return 0, 0, nil
	}
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	imported := 0
	skipped := 0
	for _, order := range orders {
		won, err := r.upsertOrder(ctx, tx, order)
		if err != nil {
			return 0, 0, err
		}
		if !won {
			skipped++
			continue
		}
		if err := r.replaceItems(ctx, tx, order); err != nil {
			return 0, 0, err
		}
		if err := r.replacePayments(ctx, tx, order); err != nil {
			return 0, 0, err
		}
		imported++
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, 0, err
	}
	return imported, skipped, nil
}

func (r *OrderRepository) ListOrders(ctx context.Context, installationID string, limit int) ([]ordersdomain.MarketplaceOrder, error) {
	if r.pool == nil {
		return nil, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT
			installation_id, provider_code, provider_order_id, provider_status, provider_status_detail,
			provider_created_at, provider_closed_at, provider_updated_at, fetched_at,
			shipping_id, cancellation_detail, tags_json, raw_provider_ref, created_at, updated_at
		FROM orders_marketplace_orders
		WHERE tenant_id = $1
		  AND installation_id = $2
		ORDER BY provider_updated_at DESC NULLS LAST, fetched_at DESC, provider_order_id DESC
		LIMIT $3
	`, r.tenantID, strings.TrimSpace(installationID), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]ordersdomain.MarketplaceOrder, 0)
	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		order.Items, err = r.listItems(ctx, order.InstallationID, order.ProviderOrderID)
		if err != nil {
			return nil, err
		}
		order.Payments, err = r.listPayments(ctx, order.InstallationID, order.ProviderOrderID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *OrderRepository) upsertOrder(ctx context.Context, tx pgx.Tx, order ordersdomain.MarketplaceOrder) (bool, error) {
	tagsJSON, err := json.Marshal(order.Tags)
	if err != nil {
		return false, err
	}
	var rawRefJSON []byte
	if order.RawProviderRef != "" {
		rawRefJSON, err = json.Marshal(order.RawProviderRef)
		if err != nil {
			return false, err
		}
	}
	tag, err := tx.Exec(ctx, `
		INSERT INTO orders_marketplace_orders (
			tenant_id, installation_id, provider_code, provider_order_id, provider_status, provider_status_detail,
			provider_created_at, provider_closed_at, provider_updated_at, fetched_at,
			shipping_id, cancellation_detail, tags_json, raw_provider_ref, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16
		)
		ON CONFLICT (tenant_id, installation_id, provider_order_id) DO UPDATE SET
			provider_code = EXCLUDED.provider_code,
			provider_status = EXCLUDED.provider_status,
			provider_status_detail = EXCLUDED.provider_status_detail,
			provider_created_at = EXCLUDED.provider_created_at,
			provider_closed_at = EXCLUDED.provider_closed_at,
			provider_updated_at = EXCLUDED.provider_updated_at,
			fetched_at = EXCLUDED.fetched_at,
			shipping_id = EXCLUDED.shipping_id,
			cancellation_detail = EXCLUDED.cancellation_detail,
			tags_json = EXCLUDED.tags_json,
			raw_provider_ref = EXCLUDED.raw_provider_ref,
			updated_at = EXCLUDED.updated_at
		WHERE orders_marketplace_orders.provider_updated_at IS NULL
		   OR (
			EXCLUDED.provider_updated_at IS NOT NULL
			AND EXCLUDED.provider_updated_at >= orders_marketplace_orders.provider_updated_at
		   )
	`, r.tenantID, order.InstallationID, order.ProviderCode, order.ProviderOrderID, order.ProviderStatus, order.ProviderStatusDetail,
		nullableTime(order.ProviderCreatedAt), nullableTime(order.ProviderClosedAt), nullableTime(order.ProviderUpdatedAt), order.FetchedAt,
		order.ShippingID, order.CancellationDetail, tagsJSON, rawRefJSON, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

func (r *OrderRepository) replaceItems(ctx context.Context, tx pgx.Tx, order ordersdomain.MarketplaceOrder) error {
	existing, err := r.loadItemsForReconciliation(ctx, tx, order.InstallationID, order.ProviderOrderID)
	if err != nil {
		return err
	}
	reconciled, err := ordersdomain.ReconcileOrderLineIdentities(existing, order.Items, ordersdomain.NewMPCLineID)
	if err != nil {
		return err
	}
	retained := make(map[ordersdomain.MPCLineID]struct{}, len(reconciled))
	for _, item := range reconciled {
		retained[item.MPCLineID] = struct{}{}
	}
	for _, item := range existing {
		if _, found := retained[item.MPCLineID]; found {
			continue
		}
		linked, err := r.isLineLinked(ctx, tx, order.InstallationID, order.ProviderOrderID, item.MPCLineID)
		if err != nil {
			return err
		}
		if linked {
			return ordersdomain.ErrOrderLineIdentityConflict
		}
	}

	// Move current positions out of the positive storage range so reordered
	// identities can be updated without transient primary-key collisions.
	if _, err := tx.Exec(ctx, `
		UPDATE orders_marketplace_order_items
		SET line_no = -line_no
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
	`, r.tenantID, order.InstallationID, order.ProviderOrderID); err != nil {
		return err
	}
	for index, item := range reconciled {
		_, err := tx.Exec(ctx, `
			INSERT INTO orders_marketplace_order_items (
				tenant_id, installation_id, provider_order_id, line_no, mpc_line_id, reconciliation_state,
				provider_item_id, provider_variation_id, seller_sku, title, quantity,
				unit_price, sale_fee_amount, link_quality, internal_product_id, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10, $11,
				$12, $13, $14, $15, $16, $17
			)
			ON CONFLICT (tenant_id, installation_id, provider_order_id, mpc_line_id) DO UPDATE SET
				line_no = EXCLUDED.line_no,
				reconciliation_state = EXCLUDED.reconciliation_state,
				provider_item_id = EXCLUDED.provider_item_id,
				provider_variation_id = EXCLUDED.provider_variation_id,
				seller_sku = EXCLUDED.seller_sku,
				title = EXCLUDED.title,
				quantity = EXCLUDED.quantity,
				unit_price = EXCLUDED.unit_price,
				sale_fee_amount = EXCLUDED.sale_fee_amount,
				link_quality = EXCLUDED.link_quality,
				internal_product_id = EXCLUDED.internal_product_id,
				updated_at = EXCLUDED.updated_at
		`, r.tenantID, order.InstallationID, order.ProviderOrderID, index+1,
			item.MPCLineID, item.ReconciliationState,
			item.ProviderItemID, item.ProviderVariationID, item.SellerSKU, item.Title, item.Quantity,
			nullableFloat8(item.UnitPrice), nullableFloat8(item.SaleFeeAmount), item.LinkQuality, nullableInt4(item.InternalProductID), order.CreatedAt, order.UpdatedAt)
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `
		DELETE FROM orders_marketplace_order_items
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
		  AND line_no < 0
	`, r.tenantID, order.InstallationID, order.ProviderOrderID)
	return err
}

func (r *OrderRepository) isLineLinked(ctx context.Context, tx pgx.Tx, installationID, providerOrderID string, mpcLineID ordersdomain.MPCLineID) (bool, error) {
	var linked bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM orders_sankhya_linkage_lines
			WHERE tenant_id = $1 AND installation_id = $2
			  AND provider_order_id = $3 AND mpc_line_id = $4
		)
	`, r.tenantID, installationID, providerOrderID, mpcLineID).Scan(&linked)
	return linked, err
}

func (r *OrderRepository) loadItemsForReconciliation(ctx context.Context, tx pgx.Tx, installationID, providerOrderID string) ([]ordersdomain.MarketplaceOrderItem, error) {
	rows, err := tx.Query(ctx, `
		SELECT mpc_line_id, reconciliation_state, provider_item_id, provider_variation_id, seller_sku
		FROM orders_marketplace_order_items
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
		ORDER BY line_no ASC
		FOR UPDATE
	`, r.tenantID, installationID, providerOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ordersdomain.MarketplaceOrderItem, 0)
	for rows.Next() {
		var item ordersdomain.MarketplaceOrderItem
		if err := rows.Scan(
			&item.MPCLineID,
			&item.ReconciliationState,
			&item.ProviderItemID,
			&item.ProviderVariationID,
			&item.SellerSKU,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *OrderRepository) replacePayments(ctx context.Context, tx pgx.Tx, order ordersdomain.MarketplaceOrder) error {
	if _, err := tx.Exec(ctx, `
		DELETE FROM orders_marketplace_order_payments
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
	`, r.tenantID, order.InstallationID, order.ProviderOrderID); err != nil {
		return err
	}
	for index, payment := range order.Payments {
		_, err := tx.Exec(ctx, `
			INSERT INTO orders_marketplace_order_payments (
				tenant_id, installation_id, provider_order_id, line_no,
				provider_payment_id, provider_status, transaction_amount, total_paid_amount,
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8,
				$9, $10
			)
		`, r.tenantID, order.InstallationID, order.ProviderOrderID, index+1,
			payment.ProviderPaymentID, payment.ProviderStatus, nullableFloat8(payment.TransactionAmount), nullableFloat8(payment.TotalPaidAmount),
			order.CreatedAt, order.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderRepository) listItems(ctx context.Context, installationID, providerOrderID string) ([]ordersdomain.MarketplaceOrderItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT mpc_line_id, reconciliation_state, provider_item_id, provider_variation_id, seller_sku, title, quantity, unit_price, sale_fee_amount, link_quality, internal_product_id
		FROM orders_marketplace_order_items
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
		ORDER BY line_no ASC
	`, r.tenantID, installationID, providerOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ordersdomain.MarketplaceOrderItem, 0)
	for rows.Next() {
		var item ordersdomain.MarketplaceOrderItem
		var unitPrice pgtype.Float8
		var saleFeeAmount pgtype.Float8
		var internalProductID pgtype.Int4
		if err := rows.Scan(
			&item.MPCLineID,
			&item.ReconciliationState,
			&item.ProviderItemID,
			&item.ProviderVariationID,
			&item.SellerSKU,
			&item.Title,
			&item.Quantity,
			&unitPrice,
			&saleFeeAmount,
			&item.LinkQuality,
			&internalProductID,
		); err != nil {
			return nil, err
		}
		item.UnitPrice = scanFloat8(unitPrice)
		item.SaleFeeAmount = scanFloat8(saleFeeAmount)
		item.InternalProductID = scanInt4(internalProductID)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *OrderRepository) listPayments(ctx context.Context, installationID, providerOrderID string) ([]ordersdomain.MarketplaceOrderPayment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT provider_payment_id, provider_status, transaction_amount, total_paid_amount
		FROM orders_marketplace_order_payments
		WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
		ORDER BY line_no ASC
	`, r.tenantID, installationID, providerOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]ordersdomain.MarketplaceOrderPayment, 0)
	for rows.Next() {
		var payment ordersdomain.MarketplaceOrderPayment
		var transactionAmount pgtype.Float8
		var totalPaidAmount pgtype.Float8
		if err := rows.Scan(
			&payment.ProviderPaymentID,
			&payment.ProviderStatus,
			&transactionAmount,
			&totalPaidAmount,
		); err != nil {
			return nil, err
		}
		payment.TransactionAmount = scanFloat8(transactionAmount)
		payment.TotalPaidAmount = scanFloat8(totalPaidAmount)
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}

func scanOrder(scanner interface{ Scan(dest ...any) error }) (ordersdomain.MarketplaceOrder, error) {
	var order ordersdomain.MarketplaceOrder
	var providerCreatedAt pgtype.Timestamptz
	var providerClosedAt pgtype.Timestamptz
	var providerUpdatedAt pgtype.Timestamptz
	var fetchedAt pgtype.Timestamptz
	var createdAt pgtype.Timestamptz
	var updatedAt pgtype.Timestamptz
	var tagsJSON []byte
	var rawRefJSON []byte
	if err := scanner.Scan(
		&order.InstallationID,
		&order.ProviderCode,
		&order.ProviderOrderID,
		&order.ProviderStatus,
		&order.ProviderStatusDetail,
		&providerCreatedAt,
		&providerClosedAt,
		&providerUpdatedAt,
		&fetchedAt,
		&order.ShippingID,
		&order.CancellationDetail,
		&tagsJSON,
		&rawRefJSON,
		&createdAt,
		&updatedAt,
	); err != nil {
		return ordersdomain.MarketplaceOrder{}, err
	}
	order.ProviderCreatedAt = scanTime(providerCreatedAt)
	order.ProviderClosedAt = scanTime(providerClosedAt)
	order.ProviderUpdatedAt = scanTime(providerUpdatedAt)
	order.FetchedAt = fetchedAt.Time.UTC()
	order.CreatedAt = createdAt.Time.UTC()
	order.UpdatedAt = updatedAt.Time.UTC()
	if len(tagsJSON) > 0 {
		_ = json.Unmarshal(tagsJSON, &order.Tags)
	}
	if len(rawRefJSON) > 0 {
		_ = json.Unmarshal(rawRefJSON, &order.RawProviderRef)
	}
	return order, nil
}

func nullableTime(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: value.UTC(), Valid: true}
}

func nullableFloat8(value *float64) pgtype.Float8 {
	if value == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *value, Valid: true}
}

func nullableInt4(value *int) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: int32(*value), Valid: true}
}

func scanTime(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time.UTC()
	return &t
}

func scanFloat8(value pgtype.Float8) *float64 {
	if !value.Valid {
		return nil
	}
	v := value.Float64
	return &v
}

func scanInt4(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	v := int(value.Int32)
	return &v
}
