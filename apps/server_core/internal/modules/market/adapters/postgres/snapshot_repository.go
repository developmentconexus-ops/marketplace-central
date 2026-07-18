package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

var _ ports.SnapshotStore = (*Repository)(nil)

func (r *Repository) AppendSnapshot(ctx context.Context, snapshot domain.MarketPriceSnapshot) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO market_price_snapshots (
			tenant_id, scope, ref_id, source, price_amount, price_currency,
			observed_at, fetched_at, expires_at, status, failure_reason, request_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (tenant_id, scope, ref_id, source, fetched_at) DO NOTHING
	`, r.tenantID, snapshot.Scope, snapshot.RefID, snapshot.Source,
		moneyAmount(snapshot.Price), moneyCurrency(snapshot.Price), snapshot.ObservedAt,
		snapshot.FetchedAt, snapshot.ExpiresAt, snapshot.Status, snapshot.FailureReason, snapshot.RequestID)
	return err
}

func (r *Repository) LatestValidSnapshot(ctx context.Context, scope domain.MarketPriceScope, refID string, source domain.MarketPriceSource, asOf time.Time) (*domain.MarketPriceSnapshot, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT scope, ref_id, source, price_amount, price_currency,
			observed_at, fetched_at, expires_at, status, failure_reason, request_id
		FROM market_price_snapshots
		WHERE tenant_id = $1 AND scope = $2 AND ref_id = $3 AND source = $4 AND status = 'VALID'
		ORDER BY fetched_at DESC
		LIMIT 1
	`, r.tenantID, scope, refID, source)
	snapshot, err := scanSnapshot(row, asOf)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}

func (r *Repository) SnapshotSeries(ctx context.Context, scope domain.MarketPriceScope, refID string, source domain.MarketPriceSource, asOf time.Time) ([]domain.MarketPriceSnapshot, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT scope, ref_id, source, price_amount, price_currency,
			observed_at, fetched_at, expires_at, status, failure_reason, request_id
		FROM market_price_snapshots
		WHERE tenant_id = $1 AND scope = $2 AND ref_id = $3 AND source = $4
		ORDER BY fetched_at DESC
	`, r.tenantID, scope, refID, source)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snapshots := []domain.MarketPriceSnapshot{}
	for rows.Next() {
		snapshot, err := scanSnapshot(rows, asOf)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return snapshots, nil
}

func scanSnapshot(row interface{ Scan(...any) error }, asOf time.Time) (domain.MarketPriceSnapshot, error) {
	var snapshot domain.MarketPriceSnapshot
	var amount, currency *string
	if err := row.Scan(&snapshot.Scope, &snapshot.RefID, &snapshot.Source, &amount, &currency,
		&snapshot.ObservedAt, &snapshot.FetchedAt, &snapshot.ExpiresAt, &snapshot.Status,
		&snapshot.FailureReason, &snapshot.RequestID); err != nil {
		return domain.MarketPriceSnapshot{}, err
	}
	snapshot.Price = scanMoney(amount, currency)
	if snapshot.Status == domain.MarketPriceSnapshotStatusValid && !snapshot.ExpiresAt.After(asOf) {
		snapshot.Status = domain.MarketPriceSnapshotStatusExpired
	}
	return snapshot, nil
}
