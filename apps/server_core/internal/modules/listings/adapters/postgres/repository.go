package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/listings/domain"
	"marketplace-central/apps/server_core/internal/modules/listings/ports"
)

var _ ports.CompletedPullStore = (*Repository)(nil)

type Repository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewRepository(pool *pgxpool.Pool, tenantID string) *Repository {
	return &Repository{pool: pool, tenantID: tenantID}
}

func (r *Repository) ApplyCompletedPull(ctx context.Context, installationID string, rows []domain.Listing, completedAt time.Time) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin listing completed pull: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	closedRows, err := tx.Query(ctx, `
		UPDATE listings SET status = 'closed', updated_at = $3
		WHERE tenant_id = $1 AND installation_id = $2
		RETURNING provider_listing_id, variation_id
	`, r.tenantID, installationID, completedAt.UTC())
	if err != nil {
		return fmt.Errorf("close listings: %w", err)
	}
	previous := make([]domain.ListingKey, 0)
	for closedRows.Next() {
		var key domain.ListingKey
		key.TenantID, key.InstallationID = r.tenantID, installationID
		if err := closedRows.Scan(&key.ProviderListingID, &key.VariationID); err != nil {
			closedRows.Close()
			return fmt.Errorf("scan closed listing key: %w", err)
		}
		previous = append(previous, key)
	}
	if err := closedRows.Err(); err != nil {
		closedRows.Close()
		return fmt.Errorf("iterate closed listing keys: %w", err)
	}
	closedRows.Close()

	returned := make(map[domain.ListingKey]struct{}, len(rows))
	for _, row := range rows {
		if row.Key.TenantID != r.tenantID || row.Key.InstallationID != installationID {
			return fmt.Errorf("listing key is outside repository tenant or installation")
		}
		syncError, err := json.Marshal(row.SyncError)
		if err != nil {
			return fmt.Errorf("marshal listing sync error: %w", err)
		}
		if row.SyncError == nil {
			syncError = nil
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO listings (tenant_id, installation_id, provider, provider_listing_id, variation_id, title,
				listing_type_code, status, price_amount, price_currency, published_quantity, sync_state,
				sync_error, quality_score, sales_30d, fetched_at, created_at, updated_at)
			SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18
			WHERE $1 = $19 AND $2 = $20
			ON CONFLICT (tenant_id, installation_id, provider_listing_id, variation_id) DO UPDATE SET
				provider=EXCLUDED.provider, title=EXCLUDED.title, listing_type_code=EXCLUDED.listing_type_code,
				status=EXCLUDED.status, price_amount=EXCLUDED.price_amount, price_currency=EXCLUDED.price_currency,
				published_quantity=EXCLUDED.published_quantity, sync_state=EXCLUDED.sync_state,
				sync_error=EXCLUDED.sync_error, quality_score=EXCLUDED.quality_score, sales_30d=EXCLUDED.sales_30d,
				fetched_at=EXCLUDED.fetched_at, updated_at=EXCLUDED.updated_at
		`, r.tenantID, installationID, row.Provider, row.Key.ProviderListingID, row.Key.VariationID, row.Title,
			row.ListingTypeCode, row.Status, row.PriceAmount, row.PriceCurrency, row.PublishedQuantity, row.SyncState,
			syncError, row.QualityScore, row.Sales30D, row.FetchedAt, row.CreatedAt.UTC(), completedAt.UTC(), r.tenantID, installationID)
		if err != nil {
			return fmt.Errorf("upsert listing %s/%s: %w", row.Key.ProviderListingID, row.Key.VariationID, err)
		}
		returned[row.Key] = struct{}{}
		kind, message := "synced", "Sincronizado"
		if row.Status == domain.ListingStatusPaused {
			kind, message = "paused", "Anúncio pausado no provedor"
		}
		if err := insertEvent(ctx, tx, row.Key, completedAt, kind, message); err != nil {
			return err
		}
	}
	for _, key := range previous {
		if _, ok := returned[key]; ok {
			continue
		}
		if err := insertEvent(ctx, tx, key, completedAt, "closed", "Anúncio ausente no provedor — fechado"); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit listing completed pull: %w", err)
	}
	return nil
}

func insertEvent(ctx context.Context, tx pgx.Tx, key domain.ListingKey, at time.Time, kind, message string) error {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return fmt.Errorf("create listing sync event id: %w", err)
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO listing_sync_events (tenant_id, installation_id, provider_listing_id, variation_id, event_id, at, kind, message_pt)
		SELECT $1,$2,$3,$4,$5,$6,$7,$8 WHERE $1 = $9 AND $2 = $10
	`, key.TenantID, key.InstallationID, key.ProviderListingID, key.VariationID, hex.EncodeToString(idBytes), at.UTC(), kind, message, key.TenantID, key.InstallationID)
	if err != nil {
		return fmt.Errorf("insert listing sync event %s/%s: %w", key.ProviderListingID, key.VariationID, err)
	}
	return nil
}
