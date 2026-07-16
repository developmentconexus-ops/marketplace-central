package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

var _ ports.ObservationStore = (*Repository)(nil)

func (r *Repository) AppendObservations(ctx context.Context, observations []domain.MarketObservation) error {
	if len(observations) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, observation := range observations {
		catalogStats, err := json.Marshal(observation.CatalogStats)
		if err != nil {
			return err
		}
		if observation.CatalogStats == nil {
			catalogStats = nil
		}
		batch.Queue(`
			INSERT INTO market_observations (
				tenant_id, listing_id,
				our_sale_price_amount, our_sale_price_currency, competitive_status,
				winner_price_amount, winner_price_currency,
				competitive_target_amount, competitive_target_currency,
				catalog_offer_price_amount, catalog_offer_price_currency,
				catalog_stats, evidence_state, captured_at, source
			) VALUES (
				$1, $2,
				$3, $4, $5,
				$6, $7,
				$8, $9,
				$10, $11,
				$12, $13, $14, $15
			)
		`, r.tenantID, observation.ListingID,
			moneyAmount(observation.OurSalePrice), moneyCurrency(observation.OurSalePrice), competitiveStatus(observation.CompetitiveStatus),
			moneyAmount(observation.WinnerPrice), moneyCurrency(observation.WinnerPrice),
			moneyAmount(observation.CompetitiveTarget), moneyCurrency(observation.CompetitiveTarget),
			moneyAmount(observation.CatalogOfferPrice), moneyCurrency(observation.CatalogOfferPrice),
			catalogStats, string(observation.EvidenceState), observation.CapturedAt, source(observation.Source))
	}

	results := r.pool.SendBatch(ctx, batch)
	return results.Close()
}

func (r *Repository) LatestObservations(ctx context.Context, listingIDs []string) ([]domain.MarketObservation, error) {
	if len(listingIDs) == 0 {
		return []domain.MarketObservation{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (listing_id)
			listing_id,
			our_sale_price_amount, our_sale_price_currency, competitive_status,
			winner_price_amount, winner_price_currency,
			competitive_target_amount, competitive_target_currency,
			catalog_offer_price_amount, catalog_offer_price_currency,
			catalog_stats, evidence_state, captured_at, source
		FROM market_observations
		WHERE tenant_id = $1 AND listing_id = ANY($2)
		ORDER BY listing_id, captured_at DESC
	`, r.tenantID, listingIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byListingID := make(map[string]domain.MarketObservation, len(listingIDs))
	for rows.Next() {
		observation, err := scanObservation(rows)
		if err != nil {
			return nil, err
		}
		byListingID[observation.ListingID] = observation
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]domain.MarketObservation, 0, len(byListingID))
	for _, listingID := range listingIDs {
		if observation, ok := byListingID[listingID]; ok {
			result = append(result, observation)
		}
	}
	return result, nil
}

func scanObservation(row interface{ Scan(dest ...any) error }) (domain.MarketObservation, error) {
	var observation domain.MarketObservation
	var ourAmount, ourCurrency, status *string
	var winnerAmount, winnerCurrency *string
	var targetAmount, targetCurrency *string
	var offerAmount, offerCurrency *string
	var catalogStats []byte
	var evidenceState, sourceValue string
	var capturedAt time.Time
	if err := row.Scan(
		&observation.ListingID,
		&ourAmount, &ourCurrency, &status,
		&winnerAmount, &winnerCurrency,
		&targetAmount, &targetCurrency,
		&offerAmount, &offerCurrency,
		&catalogStats, &evidenceState, &capturedAt, &sourceValue,
	); err != nil {
		return domain.MarketObservation{}, err
	}

	observation.OurSalePrice = scanMoney(ourAmount, ourCurrency)
	observation.WinnerPrice = scanMoney(winnerAmount, winnerCurrency)
	observation.CompetitiveTarget = scanMoney(targetAmount, targetCurrency)
	observation.CatalogOfferPrice = scanMoney(offerAmount, offerCurrency)
	if status != nil {
		value := domain.CompetitiveStatus(*status)
		observation.CompetitiveStatus = &value
	}
	if len(catalogStats) > 0 {
		var value domain.CatalogStats
		if err := json.Unmarshal(catalogStats, &value); err != nil {
			return domain.MarketObservation{}, err
		}
		observation.CatalogStats = &value
	}
	observation.EvidenceState = domain.EvidenceState(evidenceState)
	capturedAt = capturedAt.UTC()
	observation.CapturedAt = &capturedAt
	sourceTyped := domain.Source(sourceValue)
	observation.Source = &sourceTyped
	return observation, nil
}

func moneyAmount(value *domain.Money) any {
	if value == nil {
		return nil
	}
	return value.Amount
}

func moneyCurrency(value *domain.Money) any {
	if value == nil {
		return nil
	}
	return value.Currency
}

func competitiveStatus(value *domain.CompetitiveStatus) any {
	if value == nil {
		return nil
	}
	return string(*value)
}

func source(value *domain.Source) any {
	if value == nil {
		return nil
	}
	return string(*value)
}

func scanMoney(amount, currency *string) *domain.Money {
	if amount == nil && currency == nil {
		return nil
	}
	return &domain.Money{Amount: *amount, Currency: *currency}
}
