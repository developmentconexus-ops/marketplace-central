package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"marketplace-central/apps/server_core/internal/modules/market/domain"
	"marketplace-central/apps/server_core/internal/modules/market/ports"
)

var _ ports.CompetitiveSignalStore = (*Repository)(nil)

func (r *Repository) AppendCompetitiveSignals(ctx context.Context, signals []domain.CompetitiveSignal) error {
	if len(signals) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, signal := range signals {
		var positionRank, positionTotal any
		if signal.Position != nil {
			positionRank = signal.Position.Rank
			positionTotal = signal.Position.Total
		}
		batch.Queue(`
			INSERT INTO market_competitive_signals (
				tenant_id, listing_id, our_price_amount, our_price_currency,
				winner_price_amount, winner_price_currency, target_price_amount,
				target_price_currency, position_rank, position_total, fetched_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (tenant_id, listing_id, fetched_at) DO NOTHING
		`, r.tenantID, signal.ListingID, moneyAmount(signal.OurPrice), moneyCurrency(signal.OurPrice),
			moneyAmount(signal.WinnerPrice), moneyCurrency(signal.WinnerPrice), moneyAmount(signal.TargetPrice),
			moneyCurrency(signal.TargetPrice), positionRank, positionTotal, signal.FetchedAt)
	}
	results := r.pool.SendBatch(ctx, batch)
	return results.Close()
}

func (r *Repository) LatestCompetitiveSignals(ctx context.Context, listingIDs []string) ([]domain.CompetitiveSignal, error) {
	if len(listingIDs) == 0 {
		return []domain.CompetitiveSignal{}, nil
	}
	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT ON (listing_id)
			listing_id, our_price_amount, our_price_currency,
			winner_price_amount, winner_price_currency, target_price_amount,
			target_price_currency, position_rank, position_total, fetched_at
		FROM market_competitive_signals
		WHERE tenant_id = $1 AND listing_id = ANY($2)
		ORDER BY listing_id, fetched_at DESC
	`, r.tenantID, listingIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byListingID := make(map[string]domain.CompetitiveSignal, len(listingIDs))
	for rows.Next() {
		var signal domain.CompetitiveSignal
		var ourAmount, ourCurrency, winnerAmount, winnerCurrency, targetAmount, targetCurrency *string
		var positionRank, positionTotal *int
		if err := rows.Scan(&signal.ListingID, &ourAmount, &ourCurrency, &winnerAmount, &winnerCurrency,
			&targetAmount, &targetCurrency, &positionRank, &positionTotal, &signal.FetchedAt); err != nil {
			return nil, err
		}
		signal.OurPrice = scanMoney(ourAmount, ourCurrency)
		signal.WinnerPrice = scanMoney(winnerAmount, winnerCurrency)
		signal.TargetPrice = scanMoney(targetAmount, targetCurrency)
		if positionRank != nil && positionTotal != nil {
			signal.Position = &domain.MarketPosition{Rank: *positionRank, Total: *positionTotal}
		}
		byListingID[signal.ListingID] = signal
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]domain.CompetitiveSignal, 0, len(byListingID))
	for _, listingID := range listingIDs {
		if signal, ok := byListingID[listingID]; ok {
			result = append(result, signal)
		}
	}
	return result, nil
}
