package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

var _ ports.ReferenceStore = (*Repository)(nil)

func (r *Repository) AppendReferences(ctx context.Context, references []domain.MarketReference) error {
	if len(references) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, reference := range references {
		catalogStats, err := json.Marshal(reference.CatalogStats)
		if err != nil {
			return err
		}
		if reference.CatalogStats == nil {
			catalogStats = nil
		}
		batch.Queue(`
			INSERT INTO market_references (
				tenant_id, product_id, catalog_product_id, match_state,
				match_method, catalog_stats, evidence_state, captured_at, source
			) VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8, $9
			)
		`, r.tenantID, reference.ProductID, reference.CatalogProductID, string(reference.MatchState),
			reference.MatchMethod, catalogStats, string(reference.EvidenceState), reference.CapturedAt, source(reference.Source))
	}

	results := r.pool.SendBatch(ctx, batch)
	return results.Close()
}

func (r *Repository) LatestReferences(ctx context.Context, productIDs []string) ([]domain.MarketReference, error) {
	if len(productIDs) == 0 {
		return []domain.MarketReference{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (product_id)
			product_id, catalog_product_id, match_state, match_method,
			catalog_stats, evidence_state, captured_at, source
		FROM market_references
		WHERE tenant_id = $1 AND product_id = ANY($2)
		ORDER BY product_id, captured_at DESC
	`, r.tenantID, productIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byProductID := make(map[string]domain.MarketReference, len(productIDs))
	for rows.Next() {
		reference, err := scanReference(rows)
		if err != nil {
			return nil, err
		}
		byProductID[reference.ProductID] = reference
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]domain.MarketReference, 0, len(byProductID))
	for _, productID := range productIDs {
		if reference, ok := byProductID[productID]; ok {
			result = append(result, reference)
		}
	}
	return result, nil
}

func scanReference(row interface{ Scan(dest ...any) error }) (domain.MarketReference, error) {
	var reference domain.MarketReference
	var catalogStats []byte
	var matchState, evidenceState, sourceValue string
	var capturedAt time.Time
	if err := row.Scan(
		&reference.ProductID, &reference.CatalogProductID, &matchState, &reference.MatchMethod,
		&catalogStats, &evidenceState, &capturedAt, &sourceValue,
	); err != nil {
		return domain.MarketReference{}, err
	}

	reference.MatchState = domain.MatchState(matchState)
	if len(catalogStats) > 0 {
		var value domain.CatalogStats
		if err := json.Unmarshal(catalogStats, &value); err != nil {
			return domain.MarketReference{}, err
		}
		reference.CatalogStats = &value
	}
	reference.EvidenceState = domain.EvidenceState(evidenceState)
	capturedAt = capturedAt.UTC()
	reference.CapturedAt = &capturedAt
	sourceTyped := domain.Source(sourceValue)
	reference.Source = &sourceTyped
	return reference, nil
}
