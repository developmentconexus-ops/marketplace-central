package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/product_links/domain"
	"marketplace-central/apps/server_core/internal/modules/product_links/ports"
)

var _ ports.ListingSnapshotStore = (*ListingSnapshotRepository)(nil)
var _ ports.ListingSnapshotReader = (*ListingSnapshotRepository)(nil)

type ListingSnapshotRepository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewListingSnapshotRepository(pool *pgxpool.Pool, tenantID string) *ListingSnapshotRepository {
	return &ListingSnapshotRepository{pool: pool, tenantID: tenantID}
}

func (r *ListingSnapshotRepository) UpsertListingSnapshots(ctx context.Context, snapshots []domain.ListingSnapshot) error {
	if r.pool == nil || len(snapshots) == 0 {
		return nil
	}

	for _, snapshot := range snapshots {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO product_link_listing_snapshots (
				tenant_id, installation_id, provider_code, provider_item_id, provider_variation_id,
				provider_status, seller_sku, ean, title, available_quantity,
				source_updated_at, fetched_at, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10,
				$11, $12, $13, $14
			)
			ON CONFLICT (tenant_id, installation_id, provider_item_id, provider_variation_id) DO UPDATE SET
				provider_code = EXCLUDED.provider_code,
				provider_status = EXCLUDED.provider_status,
				seller_sku = EXCLUDED.seller_sku,
				ean = EXCLUDED.ean,
				title = EXCLUDED.title,
				available_quantity = EXCLUDED.available_quantity,
				source_updated_at = EXCLUDED.source_updated_at,
				fetched_at = EXCLUDED.fetched_at,
				updated_at = EXCLUDED.updated_at
		`, r.tenantID,
			snapshot.InstallationID,
			snapshot.ProviderCode,
			snapshot.ProviderItemID,
			strings.TrimSpace(snapshot.ProviderVariationID),
			snapshot.ProviderStatus,
			snapshot.SellerSKU,
			snapshot.EAN,
			snapshot.Title,
			snapshot.AvailableQuantity,
			timestamptzArg(snapshot.SourceUpdatedAt),
			snapshot.FetchedAt,
			snapshot.CreatedAt,
			snapshot.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ListingSnapshotRepository) ListListingSnapshots(ctx context.Context, installationID string, limit int) ([]domain.ListingSnapshot, error) {
	if r.pool == nil {
		return nil, nil
	}
	// limit <= 0 means EVERY snapshot of the installation, not a default page:
	// candidate generation has to see the whole catalog, and a silent cap there
	// leaves the uncapped remainder reported as "sem vínculo" forever.
	const selectSnapshots = `
		SELECT
			installation_id, provider_code, provider_item_id, provider_variation_id,
			provider_status, seller_sku, ean, title, available_quantity,
			source_updated_at, fetched_at, created_at, updated_at
		FROM product_link_listing_snapshots
		WHERE tenant_id = $1
		  AND installation_id = $2
		ORDER BY updated_at DESC, provider_item_id DESC, provider_variation_id DESC`
	query := selectSnapshots + "\n\t\tLIMIT $3"
	args := []any{r.tenantID, strings.TrimSpace(installationID), limit}
	if limit <= 0 {
		query = selectSnapshots
		args = args[:2]
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snapshots := make([]domain.ListingSnapshot, 0)
	for rows.Next() {
		var snapshot domain.ListingSnapshot
		var availableQuantity pgtype.Int4
		var sourceUpdatedAt pgtype.Timestamptz
		var fetchedAt pgtype.Timestamptz
		var createdAt pgtype.Timestamptz
		var updatedAt pgtype.Timestamptz

		if err := rows.Scan(
			&snapshot.InstallationID,
			&snapshot.ProviderCode,
			&snapshot.ProviderItemID,
			&snapshot.ProviderVariationID,
			&snapshot.ProviderStatus,
			&snapshot.SellerSKU,
			&snapshot.EAN,
			&snapshot.Title,
			&availableQuantity,
			&sourceUpdatedAt,
			&fetchedAt,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		snapshot.AvailableQuantity = scanInt4Ptr(availableQuantity)
		snapshot.SourceUpdatedAt = scanTimestamptz(sourceUpdatedAt)
		snapshot.FetchedAt = fetchedAt.Time.UTC()
		snapshot.CreatedAt = createdAt.Time.UTC()
		snapshot.UpdatedAt = updatedAt.Time.UTC()
		snapshots = append(snapshots, snapshot)
	}

	return snapshots, rows.Err()
}

func timestamptzArg(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: value.UTC(), Valid: true}
}

func scanTimestamptz(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time.UTC()
	return &t
}

func scanInt4Ptr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	v := int(value.Int32)
	return &v
}
