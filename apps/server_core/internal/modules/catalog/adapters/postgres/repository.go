package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"marketplace-central/apps/server_core/internal/modules/catalog/domain"
	"marketplace-central/apps/server_core/internal/modules/catalog/ports"
)

var _ ports.EnrichmentStore = (*EnrichmentRepository)(nil)

type EnrichmentRepository struct {
	pool     *pgxpool.Pool
	tenantID string
}

func NewEnrichmentRepository(pool *pgxpool.Pool, tenantID string) *EnrichmentRepository {
	return &EnrichmentRepository{pool: pool, tenantID: tenantID}
}

func (r *EnrichmentRepository) GetEnrichment(ctx context.Context, productID string) (domain.ProductEnrichment, error) {
	canonicalID, err := positiveProductID(productID)
	if err != nil {
		return domain.ProductEnrichment{}, err
	}
	var e domain.ProductEnrichment
	err = r.pool.QueryRow(ctx, `
		SELECT internal_product_id::text, tenant_id, height_cm, width_cm, length_cm, weight_g, suggested_price_amount
		FROM product_enrichments
		WHERE tenant_id = $1 AND internal_product_id = $2 AND resolution_status = 'mapped'
	`, r.tenantID, canonicalID).Scan(
		&e.ProductID, &e.TenantID, &e.HeightCM, &e.WidthCM, &e.LengthCM, &e.WeightG, &e.SuggestedPriceAmount,
	)
	if err == pgx.ErrNoRows {
		return domain.ProductEnrichment{ProductID: productID, TenantID: r.tenantID}, nil
	}
	if err != nil {
		return domain.ProductEnrichment{}, fmt.Errorf("get enrichment: %w", err)
	}
	return e, nil
}

// UpsertEnrichment uses partial-update semantics: only non-nil fields overwrite existing values.
func (r *EnrichmentRepository) UpsertEnrichment(ctx context.Context, e domain.ProductEnrichment) error {
	canonicalID, err := positiveProductID(e.ProductID)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO product_enrichments (product_id, tenant_id, internal_product_id, resolution_status, height_cm, width_cm, length_cm, weight_g, suggested_price_amount, updated_at)
		VALUES ($1, $2, $3, 'mapped', $4, $5, $6, $7, $8, now())
		ON CONFLICT (tenant_id, product_id) DO UPDATE SET
			internal_product_id = EXCLUDED.internal_product_id,
			resolution_status = 'mapped',
			height_cm = COALESCE(EXCLUDED.height_cm, product_enrichments.height_cm),
			width_cm = COALESCE(EXCLUDED.width_cm, product_enrichments.width_cm),
			length_cm = COALESCE(EXCLUDED.length_cm, product_enrichments.length_cm),
			weight_g = COALESCE(EXCLUDED.weight_g, product_enrichments.weight_g),
			suggested_price_amount = COALESCE(EXCLUDED.suggested_price_amount, product_enrichments.suggested_price_amount),
			updated_at = now()
	`, strconv.Itoa(canonicalID), r.tenantID, canonicalID, e.HeightCM, e.WidthCM, e.LengthCM, e.WeightG, e.SuggestedPriceAmount)
	if err != nil {
		return fmt.Errorf("upsert enrichment: %w", err)
	}
	return nil
}

func (r *EnrichmentRepository) ListEnrichments(ctx context.Context, productIDs []string) (map[string]domain.ProductEnrichment, error) {
	if len(productIDs) == 0 {
		return make(map[string]domain.ProductEnrichment), nil
	}
	canonicalIDs := make([]int, 0, len(productIDs))
	for _, productID := range productIDs {
		canonicalID, err := positiveProductID(productID)
		if err != nil {
			return nil, err
		}
		canonicalIDs = append(canonicalIDs, canonicalID)
	}
	rows, err := r.pool.Query(ctx, `
		SELECT internal_product_id::text, tenant_id, height_cm, width_cm, length_cm, weight_g, suggested_price_amount
		FROM product_enrichments
		WHERE tenant_id = $1 AND internal_product_id = ANY($2) AND resolution_status = 'mapped'
	`, r.tenantID, canonicalIDs)
	if err != nil {
		return nil, fmt.Errorf("list enrichments: %w", err)
	}
	defer rows.Close()

	result := make(map[string]domain.ProductEnrichment)
	for rows.Next() {
		var e domain.ProductEnrichment
		if err := rows.Scan(&e.ProductID, &e.TenantID, &e.HeightCM, &e.WidthCM, &e.LengthCM, &e.WeightG, &e.SuggestedPriceAmount); err != nil {
			return nil, fmt.Errorf("scan enrichment: %w", err)
		}
		result[e.ProductID] = e
	}
	return result, rows.Err()
}

func positiveProductID(raw string) (int, error) {
	trimmed := strings.TrimSpace(raw)
	id, err := strconv.Atoi(trimmed)
	if err != nil || id <= 0 || strconv.Itoa(id) != trimmed {
		return 0, fmt.Errorf("CATALOG_INVALID_PRODUCT_ID")
	}
	return id, nil
}
