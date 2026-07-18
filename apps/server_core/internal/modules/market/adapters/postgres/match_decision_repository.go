package postgres

import (
	"context"

	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

var _ ports.MatchDecisionStore = (*Repository)(nil)

func (r *Repository) AppendMatchDecision(ctx context.Context, decision domain.MatchDecision) error {
	// anchors/contradictions are TEXT[] NOT NULL; a decision with no anchors
	// (e.g. NO_CANDIDATE) carries nil slices — persist as empty arrays, not NULL.
	anchors, contradictions := decision.Anchors, decision.Contradictions
	if anchors == nil {
		anchors = []string{}
	}
	if contradictions == nil {
		contradictions = []string{}
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO market_match_decisions (
			tenant_id, product_id, catalog_product_id, match_status,
			anchors, contradictions, confidence_band, decided_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (tenant_id, product_id, decided_at) DO NOTHING
	`, r.tenantID, decision.ProductID, decision.CatalogProductID, decision.Status,
		anchors, contradictions, decision.ConfidenceBand, decision.DecidedAt)
	return err
}

func (r *Repository) LatestMatchDecisions(ctx context.Context, productIDs []string) ([]domain.MatchDecision, error) {
	if len(productIDs) == 0 {
		return []domain.MatchDecision{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (product_id)
			product_id, catalog_product_id, match_status, anchors,
			contradictions, confidence_band, decided_at
		FROM market_match_decisions
		WHERE tenant_id = $1 AND product_id = ANY($2)
		ORDER BY product_id, decided_at DESC
	`, r.tenantID, productIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byProductID := make(map[string]domain.MatchDecision, len(productIDs))
	for rows.Next() {
		var decision domain.MatchDecision
		if err := rows.Scan(&decision.ProductID, &decision.CatalogProductID, &decision.Status,
			&decision.Anchors, &decision.Contradictions, &decision.ConfidenceBand, &decision.DecidedAt); err != nil {
			return nil, err
		}
		byProductID[decision.ProductID] = decision
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.MatchDecision, 0, len(byProductID))
	for _, productID := range productIDs {
		if decision, ok := byProductID[productID]; ok {
			result = append(result, decision)
		}
	}
	return result, nil
}
