package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

var _ ports.ValidatedOfferStore = (*Repository)(nil)

func (r *Repository) AppendValidatedOffers(ctx context.Context, offers []domain.ValidatedOffer) error {
	if len(offers) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, offer := range offers {
		batch.Queue(`
			INSERT INTO market_validated_offers (
				tenant_id, catalog_product_id, seller_id, price_amount, price_currency,
				condition, match_status, fetched_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (tenant_id, catalog_product_id, seller_id, fetched_at) DO NOTHING
		`, r.tenantID, offer.CatalogProductID, offer.SellerID, moneyAmount(offer.Price),
			moneyCurrency(offer.Price), offer.Condition, offer.MatchStatus, offer.FetchedAt)
	}
	results := r.pool.SendBatch(ctx, batch)
	return results.Close()
}

func (r *Repository) ValidatedOffersByCatalogProduct(ctx context.Context, catalogProductID string) ([]domain.ValidatedOffer, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT catalog_product_id, seller_id, price_amount, price_currency, condition, match_status, fetched_at
		FROM market_validated_offers
		WHERE tenant_id = $1 AND catalog_product_id = $2
		ORDER BY fetched_at DESC, seller_id
	`, r.tenantID, catalogProductID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	offers := []domain.ValidatedOffer{}
	for rows.Next() {
		var offer domain.ValidatedOffer
		var amount, currency *string
		if err := rows.Scan(&offer.CatalogProductID, &offer.SellerID, &amount, &currency,
			&offer.Condition, &offer.MatchStatus, &offer.FetchedAt); err != nil {
			return nil, err
		}
		offer.Price = scanMoney(amount, currency)
		offers = append(offers, offer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return offers, nil
}
