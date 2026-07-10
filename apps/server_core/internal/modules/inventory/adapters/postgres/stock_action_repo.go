package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/inventory/domain"
	"marketplace-central/apps/server_core/internal/modules/inventory/ports"
)

var _ ports.StockActionStore = (*StockActionRepository)(nil)

type StockActionRepository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewStockActionRepository(pool *pgxpool.Pool, tenantID string) *StockActionRepository {
	return &StockActionRepository{pool: pool, tenantID: tenantID}
}

func (r *StockActionRepository) FindByID(ctx context.Context, actionID string) (domain.StockAction, bool, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT
			stock_action_id, installation_id, provider_code, provider_account_id, provider_item_id, provider_variation_id,
			state, trigger, before_quantity, requested_quantity, recommended_quantity, policy_id,
			internal_observed_at, provider_observed_at,
			operator_actor_type, operator_actor_id, operator_actor_name,
			reason, idempotency_key,
			blocking_code, blocking_message, failure_code, failure_message,
			provider_result_status, provider_result_code, provider_result_message, provider_response_summary,
			audit_events_json, created_at, updated_at
		FROM inventory_stock_actions
		WHERE tenant_id = $1 AND stock_action_id = $2
	`, r.tenantID, actionID)

	action, err := scanStockAction(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.StockAction{}, false, nil
		}
		return domain.StockAction{}, false, err
	}
	return action, true, nil
}

func (r *StockActionRepository) Save(ctx context.Context, action domain.StockAction) error {
	auditJSON, err := json.Marshal(action.AuditEvents)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO inventory_stock_actions (
			tenant_id, stock_action_id, installation_id, provider_code, provider_account_id, provider_item_id, provider_variation_id,
			state, trigger, before_quantity, requested_quantity, recommended_quantity, policy_id,
			internal_observed_at, provider_observed_at,
			operator_actor_type, operator_actor_id, operator_actor_name,
			reason, idempotency_key,
			blocking_code, blocking_message, failure_code, failure_message,
			provider_result_status, provider_result_code, provider_result_message, provider_response_summary,
			audit_events_json, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13,
			$14, $15,
			$16, $17, $18,
			$19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28,
			$29, $30, $31
		)
		ON CONFLICT (tenant_id, stock_action_id) DO UPDATE SET
			installation_id = EXCLUDED.installation_id,
			provider_code = EXCLUDED.provider_code,
			provider_account_id = EXCLUDED.provider_account_id,
			provider_item_id = EXCLUDED.provider_item_id,
			provider_variation_id = EXCLUDED.provider_variation_id,
			state = EXCLUDED.state,
			trigger = EXCLUDED.trigger,
			before_quantity = EXCLUDED.before_quantity,
			requested_quantity = EXCLUDED.requested_quantity,
			recommended_quantity = EXCLUDED.recommended_quantity,
			policy_id = EXCLUDED.policy_id,
			internal_observed_at = EXCLUDED.internal_observed_at,
			provider_observed_at = EXCLUDED.provider_observed_at,
			operator_actor_type = EXCLUDED.operator_actor_type,
			operator_actor_id = EXCLUDED.operator_actor_id,
			operator_actor_name = EXCLUDED.operator_actor_name,
			reason = EXCLUDED.reason,
			idempotency_key = EXCLUDED.idempotency_key,
			blocking_code = EXCLUDED.blocking_code,
			blocking_message = EXCLUDED.blocking_message,
			failure_code = EXCLUDED.failure_code,
			failure_message = EXCLUDED.failure_message,
			provider_result_status = EXCLUDED.provider_result_status,
			provider_result_code = EXCLUDED.provider_result_code,
			provider_result_message = EXCLUDED.provider_result_message,
			provider_response_summary = EXCLUDED.provider_response_summary,
			audit_events_json = EXCLUDED.audit_events_json,
			updated_at = EXCLUDED.updated_at
	`, r.tenantID, action.ActionID, action.ProviderRef.InstallationID, action.ProviderRef.ProviderCode, action.ProviderRef.ProviderAccountID, action.ProviderRef.ProviderItemID, action.ProviderRef.ProviderVariationID,
		action.State, action.Trigger, nullableInt4(action.BeforeQuantity), action.RequestedQuantity, nullableInt4(action.RecommendedQuantity), action.PolicyID,
		nullableTime(action.InternalObservedAt), nullableTime(action.ProviderObservedAt),
		action.Operator.ActorType, action.Operator.ActorID, action.Operator.ActorName,
		action.Reason, action.IdempotencyKey,
		action.BlockingReason.Code, action.BlockingReason.Message, action.FailureReason.Code, action.FailureReason.Message,
		action.ProviderResult.Status, "", action.ProviderResult.Message, action.ProviderResult.ResponseSummary,
		auditJSON, action.CreatedAt, action.UpdatedAt)
	return err
}

func scanStockAction(scanner interface{ Scan(dest ...any) error }) (domain.StockAction, error) {
	var action domain.StockAction
	var beforeQty pgtype.Int4
	var recommendedQty pgtype.Int4
	var internalObserved pgtype.Timestamptz
	var providerObserved pgtype.Timestamptz
	var providerResultCode string
	var auditJSON []byte
	err := scanner.Scan(
		&action.ActionID,
		&action.ProviderRef.InstallationID,
		&action.ProviderRef.ProviderCode,
		&action.ProviderRef.ProviderAccountID,
		&action.ProviderRef.ProviderItemID,
		&action.ProviderRef.ProviderVariationID,
		&action.State,
		&action.Trigger,
		&beforeQty,
		&action.RequestedQuantity,
		&recommendedQty,
		&action.PolicyID,
		&internalObserved,
		&providerObserved,
		&action.Operator.ActorType,
		&action.Operator.ActorID,
		&action.Operator.ActorName,
		&action.Reason,
		&action.IdempotencyKey,
		&action.BlockingReason.Code,
		&action.BlockingReason.Message,
		&action.FailureReason.Code,
		&action.FailureReason.Message,
		&action.ProviderResult.Status,
		&providerResultCode,
		&action.ProviderResult.Message,
		&action.ProviderResult.ResponseSummary,
		&auditJSON,
		&action.CreatedAt,
		&action.UpdatedAt,
	)
	_ = providerResultCode
	if err != nil {
		return domain.StockAction{}, err
	}
	action.BeforeQuantity = scanInt4(beforeQty)
	action.RecommendedQuantity = scanInt4(recommendedQty)
	action.InternalObservedAt = scanTime(internalObserved)
	action.ProviderObservedAt = scanTime(providerObserved)
	if len(auditJSON) > 0 {
		_ = json.Unmarshal(auditJSON, &action.AuditEvents)
	}
	return action, nil
}

func nullableInt4(value *int) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: int32(*value), Valid: true}
}

func nullableTime(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: value.UTC(), Valid: true}
}

func scanInt4(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	v := int(value.Int32)
	return &v
}

func scanTime(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time.UTC()
	return &t
}
