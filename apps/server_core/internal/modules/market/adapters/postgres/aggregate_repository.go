package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

var _ ports.MarketAggregateStore = (*Repository)(nil)

func (r *Repository) AppendMarketAggregates(ctx context.Context, aggregates []domain.MarketAggregate) error {
	if len(aggregates) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, aggregate := range aggregates {
		batch.Queue(`
			INSERT INTO market_aggregates (
				tenant_id, product_id, median_amount, median_currency,
				min_valid_amount, min_valid_currency, max_valid_amount, max_valid_currency,
				n_offers, n_sellers, source, fetched_at, computed_at, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			ON CONFLICT (tenant_id, product_id, source, computed_at) DO NOTHING
		`, r.tenantID, aggregate.ProductID, moneyAmount(aggregate.Median), moneyCurrency(aggregate.Median),
			moneyAmount(aggregate.MinValid), moneyCurrency(aggregate.MinValid),
			moneyAmount(aggregate.MaxValid), moneyCurrency(aggregate.MaxValid), aggregate.NOffers,
			aggregate.NSellers, aggregate.Source, aggregate.FetchedAt, aggregate.ComputedAt, aggregate.Status)
	}
	results := r.pool.SendBatch(ctx, batch)
	return results.Close()
}

func (r *Repository) LatestMarketAggregates(ctx context.Context, productIDs []string) ([]domain.MarketAggregate, error) {
	if len(productIDs) == 0 {
		return []domain.MarketAggregate{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (product_id)
			product_id, median_amount, median_currency, min_valid_amount,
			min_valid_currency, max_valid_amount, max_valid_currency, n_offers, n_sellers, source, fetched_at,
			computed_at, status
		FROM market_aggregates
		WHERE tenant_id = $1 AND product_id = ANY($2)
		ORDER BY product_id, computed_at DESC, source
	`, r.tenantID, productIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byProductID := make(map[string]domain.MarketAggregate, len(productIDs))
	for rows.Next() {
		var aggregate domain.MarketAggregate
		var medianAmount, medianCurrency, minValidAmount, minValidCurrency, maxValidAmount, maxValidCurrency *string
		if err := rows.Scan(&aggregate.ProductID, &medianAmount, &medianCurrency, &minValidAmount,
			&minValidCurrency, &maxValidAmount, &maxValidCurrency, &aggregate.NOffers, &aggregate.NSellers, &aggregate.Source,
			&aggregate.FetchedAt, &aggregate.ComputedAt, &aggregate.Status); err != nil {
			return nil, err
		}
		aggregate.Median = scanMoney(medianAmount, medianCurrency)
		aggregate.MinValid = scanMoney(minValidAmount, minValidCurrency)
		aggregate.MaxValid = scanMoney(maxValidAmount, maxValidCurrency)
		byProductID[aggregate.ProductID] = aggregate
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.MarketAggregate, 0, len(byProductID))
	for _, productID := range productIDs {
		if aggregate, ok := byProductID[productID]; ok {
			result = append(result, aggregate)
		}
	}
	return result, nil
}

// StaleAggregateProductIDs enumerates which products' latest market_aggregates
// row has fallen behind olderThan, oldest first, capped at limit. market_aggregates
// is append-only (PK tenant_id, product_id, source, computed_at; migrations/0053:23-42),
// so DISTINCT ON picks the latest row per product exactly like LatestMarketAggregates
// above, and the age filter is applied AFTER that pick, on computed_at (when this
// aggregate was computed) rather than fetched_at (when the underlying evidence was
// first observed). Filtering age before picking latest-per-product would wrongly
// exclude a product that has a recent row buried under an older one for the same
// product/source pair.
func (r *Repository) StaleAggregateProductIDs(ctx context.Context, olderThan time.Duration, limit int) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT product_id
		  FROM (
		       SELECT DISTINCT ON (product_id) product_id, computed_at
		         FROM market_aggregates
		        WHERE tenant_id = $1
		        ORDER BY product_id, computed_at DESC, source
		       ) latest
		 WHERE latest.computed_at < now() - $2::interval
		 ORDER BY latest.computed_at ASC
		 LIMIT $3
	`, r.tenantID, olderThan.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0, limit)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *Repository) LatestMarketAggregatesBySource(ctx context.Context, productIDs []string, source domain.MarketPriceSource) ([]domain.MarketAggregate, error) {
	if len(productIDs) == 0 {
		return []domain.MarketAggregate{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (product_id)
			product_id, median_amount, median_currency, min_valid_amount,
			min_valid_currency, max_valid_amount, max_valid_currency, n_offers, n_sellers, source, fetched_at,
			computed_at, status
		FROM market_aggregates
		WHERE tenant_id = $1 AND product_id = ANY($2) AND source = $3
		ORDER BY product_id, computed_at DESC
	`, r.tenantID, productIDs, source)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byProductID := make(map[string]domain.MarketAggregate, len(productIDs))
	for rows.Next() {
		var aggregate domain.MarketAggregate
		var medianAmount, medianCurrency, minValidAmount, minValidCurrency, maxValidAmount, maxValidCurrency *string
		if err := rows.Scan(&aggregate.ProductID, &medianAmount, &medianCurrency, &minValidAmount,
			&minValidCurrency, &maxValidAmount, &maxValidCurrency, &aggregate.NOffers, &aggregate.NSellers, &aggregate.Source,
			&aggregate.FetchedAt, &aggregate.ComputedAt, &aggregate.Status); err != nil {
			return nil, err
		}
		aggregate.Median = scanMoney(medianAmount, medianCurrency)
		aggregate.MinValid = scanMoney(minValidAmount, minValidCurrency)
		aggregate.MaxValid = scanMoney(maxValidAmount, maxValidCurrency)
		byProductID[aggregate.ProductID] = aggregate
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.MarketAggregate, 0, len(byProductID))
	for _, productID := range productIDs {
		if aggregate, ok := byProductID[productID]; ok {
			result = append(result, aggregate)
		}
	}
	return result, nil
}
