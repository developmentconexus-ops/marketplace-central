package postgres

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	profitabilitydomain "marketplace-central/apps/server_core/internal/modules/profitability/domain"
	"marketplace-central/apps/server_core/internal/modules/profitability/ports"
)

var _ ports.MarginInputStore = (*Store)(nil)
var _ ports.ManualAdjustmentStore = (*Store)(nil)
var _ ports.ProfitSnapshotStore = (*Store)(nil)

type Store struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewStore(pool *pgxpool.Pool, tenantID string) *Store {
	return &Store{pool: pool, tenantID: tenantID}
}

func (s *Store) ReplaceInputs(ctx context.Context, installationID string, orders []string, inputs []profitabilitydomain.MarginInput) error {
	if s.pool == nil {
		return nil
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, orderID := range orders {
		if _, err := tx.Exec(ctx, `
			DELETE FROM profitability_margin_inputs
			WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
		`, s.tenantID, installationID, strings.TrimSpace(orderID)); err != nil {
			return err
		}
	}
	for _, input := range inputs {
		_, err := tx.Exec(ctx, `
			INSERT INTO profitability_margin_inputs (
				tenant_id, installation_id, provider_order_id, provider_item_id, provider_variation_id,
				scope, kind, amount, currency, source_system, source_reference, observed_at,
				quality, quality_reason, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10, $11, $12,
				$13, $14, $15, $16
			)
		`, s.tenantID, input.InstallationID, input.ProviderOrderID, input.ProviderItemID, input.ProviderVariationID,
			input.Scope, input.Kind, nullableFloat8(input.Amount), input.Currency, input.SourceSystem, input.SourceReference, nullableTime(input.ObservedAt),
			input.Quality, input.QualityReason, input.CreatedAt, input.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *Store) ListInputs(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.MarginInput, error) {
	if s.pool == nil {
		return nil, nil
	}
	rows, err := s.pool.Query(ctx, `
		SELECT installation_id, provider_order_id, provider_item_id, provider_variation_id,
		       scope, kind, amount, currency, source_system, source_reference, observed_at,
		       quality, quality_reason, created_at, updated_at
		FROM profitability_margin_inputs
		WHERE tenant_id = $1 AND installation_id = $2
		ORDER BY updated_at DESC, provider_order_id DESC, provider_item_id DESC, kind DESC
		LIMIT $3
	`, s.tenantID, strings.TrimSpace(installationID), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]profitabilitydomain.MarginInput, 0)
	for rows.Next() {
		var item profitabilitydomain.MarginInput
		var amount pgtype.Float8
		var observedAt pgtype.Timestamptz
		var createdAt pgtype.Timestamptz
		var updatedAt pgtype.Timestamptz
		if err := rows.Scan(
			&item.InstallationID, &item.ProviderOrderID, &item.ProviderItemID, &item.ProviderVariationID,
			&item.Scope, &item.Kind, &amount, &item.Currency, &item.SourceSystem, &item.SourceReference, &observedAt,
			&item.Quality, &item.QualityReason, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		item.Amount = scanFloat8(amount)
		item.ObservedAt = scanTime(observedAt)
		item.CreatedAt = createdAt.Time.UTC()
		item.UpdatedAt = updatedAt.Time.UTC()
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) CreateOrGetAdjustment(ctx context.Context, adjustment profitabilitydomain.ManualAdjustment) (profitabilitydomain.ManualAdjustment, error) {
	if s.pool == nil {
		return adjustment, nil
	}
	commandTag, err := s.pool.Exec(ctx, `
		INSERT INTO profitability_manual_adjustments (
			tenant_id, adjustment_id, installation_id, provider_order_id, provider_item_id, provider_variation_id,
			scope, category, amount, currency, reason, idempotency_key,
			actor_type, actor_id, actor_name, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16
		) ON CONFLICT (tenant_id, installation_id, idempotency_key) WHERE idempotency_key <> '' DO NOTHING
	`, s.tenantID, adjustment.AdjustmentID, adjustment.InstallationID, adjustment.ProviderOrderID, adjustment.ProviderItemID, adjustment.ProviderVariationID,
		adjustment.Scope, adjustment.Category, adjustment.Amount, adjustment.Currency, adjustment.Reason, adjustment.IdempotencyKey,
		adjustment.Actor.ActorType, adjustment.Actor.ActorID, adjustment.Actor.ActorName, adjustment.CreatedAt)
	if err != nil {
		return profitabilitydomain.ManualAdjustment{}, err
	}
	if commandTag.RowsAffected() == 1 {
		return adjustment, nil
	}
	var stored profitabilitydomain.ManualAdjustment
	err = s.pool.QueryRow(ctx, `
		SELECT adjustment_id, idempotency_key, installation_id, provider_order_id, provider_item_id, provider_variation_id,
		       scope, category, amount, currency, reason, actor_type, actor_id, actor_name, created_at
		FROM profitability_manual_adjustments
		WHERE tenant_id = $1 AND installation_id = $2 AND idempotency_key = $3
	`, s.tenantID, adjustment.InstallationID, adjustment.IdempotencyKey).Scan(
		&stored.AdjustmentID, &stored.IdempotencyKey, &stored.InstallationID, &stored.ProviderOrderID, &stored.ProviderItemID, &stored.ProviderVariationID,
		&stored.Scope, &stored.Category, &stored.Amount, &stored.Currency, &stored.Reason,
		&stored.Actor.ActorType, &stored.Actor.ActorID, &stored.Actor.ActorName, &stored.CreatedAt,
	)
	if err == nil {
		stored.CreatedAt = stored.CreatedAt.UTC()
	}
	return stored, err
}

func (s *Store) ListAdjustments(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.ManualAdjustment, error) {
	if s.pool == nil {
		return nil, nil
	}
	rows, err := s.pool.Query(ctx, `
		SELECT adjustment_id, idempotency_key, installation_id, provider_order_id, provider_item_id, provider_variation_id,
		       scope, category, amount, currency, reason, actor_type, actor_id, actor_name, created_at
		FROM profitability_manual_adjustments
		WHERE tenant_id = $1 AND installation_id = $2
		ORDER BY created_at DESC, adjustment_id DESC
		LIMIT $3
	`, s.tenantID, strings.TrimSpace(installationID), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]profitabilitydomain.ManualAdjustment, 0)
	for rows.Next() {
		var item profitabilitydomain.ManualAdjustment
		if err := rows.Scan(
			&item.AdjustmentID, &item.IdempotencyKey, &item.InstallationID, &item.ProviderOrderID, &item.ProviderItemID, &item.ProviderVariationID,
			&item.Scope, &item.Category, &item.Amount, &item.Currency, &item.Reason,
			&item.Actor.ActorType, &item.Actor.ActorID, &item.Actor.ActorName, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) ReplaceSnapshots(ctx context.Context, installationID string, orders []string, snapshots []profitabilitydomain.ProfitSnapshot) error {
	if s.pool == nil {
		return nil
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for _, orderID := range orders {
		if _, err := tx.Exec(ctx, `
			DELETE FROM profitability_profit_snapshots
			WHERE tenant_id = $1 AND installation_id = $2 AND provider_order_id = $3
		`, s.tenantID, installationID, strings.TrimSpace(orderID)); err != nil {
			return err
		}
	}
	for _, snapshot := range snapshots {
		flagsJSON, err := json.Marshal(snapshot.Flags)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO profitability_profit_snapshots (
				tenant_id, installation_id, provider_order_id, provider_item_id, provider_variation_id,
				scope, currency, realization_state, revenue_amount, sale_fee_amount, cost_amount, tax_amount,
				freight_amount, commission_amount, adjustment_amount, contribution_amount, margin_percent,
				quality, flags_json, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10, $11, $12,
				$13, $14, $15, $16, $17,
				$18, $19, $20, $21
			)
		`, s.tenantID, snapshot.InstallationID, snapshot.ProviderOrderID, snapshot.ProviderItemID, snapshot.ProviderVariationID,
			snapshot.Scope, snapshot.Currency, snapshot.RealizationState, nullableFloat8(snapshot.RevenueAmount), nullableFloat8(snapshot.SaleFeeAmount), nullableFloat8(snapshot.CostAmount), nullableFloat8(snapshot.TaxAmount),
			nullableFloat8(snapshot.FreightAmount), nullableFloat8(snapshot.CommissionAmount), nullableFloat8(snapshot.AdjustmentAmount), nullableFloat8(snapshot.ContributionAmount), nullableFloat8(snapshot.MarginPercent),
			snapshot.Quality, flagsJSON, snapshot.CreatedAt, snapshot.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *Store) ListSnapshots(ctx context.Context, installationID string, limit int) ([]profitabilitydomain.ProfitSnapshot, error) {
	if s.pool == nil {
		return nil, nil
	}
	rows, err := s.pool.Query(ctx, `
		SELECT installation_id, provider_order_id, provider_item_id, provider_variation_id,
		       scope, currency, realization_state, revenue_amount, sale_fee_amount, cost_amount, tax_amount,
		       freight_amount, commission_amount, adjustment_amount, contribution_amount, margin_percent,
		       quality, flags_json, created_at, updated_at
		FROM profitability_profit_snapshots
		WHERE tenant_id = $1 AND installation_id = $2
		ORDER BY updated_at DESC, provider_order_id DESC, provider_item_id DESC
		LIMIT $3
	`, s.tenantID, strings.TrimSpace(installationID), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]profitabilitydomain.ProfitSnapshot, 0)
	for rows.Next() {
		var item profitabilitydomain.ProfitSnapshot
		var revenue pgtype.Float8
		var saleFee pgtype.Float8
		var cost pgtype.Float8
		var tax pgtype.Float8
		var freight pgtype.Float8
		var commission pgtype.Float8
		var adjustment pgtype.Float8
		var contribution pgtype.Float8
		var margin pgtype.Float8
		var createdAt pgtype.Timestamptz
		var updatedAt pgtype.Timestamptz
		var flagsJSON []byte
		if err := rows.Scan(
			&item.InstallationID, &item.ProviderOrderID, &item.ProviderItemID, &item.ProviderVariationID,
			&item.Scope, &item.Currency, &item.RealizationState, &revenue, &saleFee, &cost, &tax,
			&freight, &commission, &adjustment, &contribution, &margin,
			&item.Quality, &flagsJSON, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		item.RevenueAmount = scanFloat8(revenue)
		item.SaleFeeAmount = scanFloat8(saleFee)
		item.CostAmount = scanFloat8(cost)
		item.TaxAmount = scanFloat8(tax)
		item.FreightAmount = scanFloat8(freight)
		item.CommissionAmount = scanFloat8(commission)
		item.AdjustmentAmount = scanFloat8(adjustment)
		item.ContributionAmount = scanFloat8(contribution)
		item.MarginPercent = scanFloat8(margin)
		item.CreatedAt = createdAt.Time.UTC()
		item.UpdatedAt = updatedAt.Time.UTC()
		if len(flagsJSON) > 0 {
			_ = json.Unmarshal(flagsJSON, &item.Flags)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func nullableFloat8(value *float64) pgtype.Float8 {
	if value == nil {
		return pgtype.Float8{}
	}
	return pgtype.Float8{Float64: *value, Valid: true}
}

func nullableTime(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: value.UTC(), Valid: true}
}

func scanFloat8(value pgtype.Float8) *float64 {
	if !value.Valid {
		return nil
	}
	v := value.Float64
	return &v
}

func scanTime(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time.UTC()
	return &t
}
