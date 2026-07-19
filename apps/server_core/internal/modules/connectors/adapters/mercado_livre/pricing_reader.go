package mercadolivre

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/connectors/domain"
)

type mlOwnItemPricingResponse struct {
	Amount        *json.Number `json:"amount"`
	Currency      string       `json:"currency_id"`
	RegularAmount *json.Number `json:"regular_amount"`
}

type mlPriceToWinResponse struct {
	Status       string            `json:"status"`
	CurrentPrice *json.Number      `json:"current_price"`
	TargetPrice  *json.Number      `json:"price_to_win"`
	Currency     string            `json:"currency_id"`
	Position     *mlMarketPosition `json:"position"`
}

type mlMarketPosition struct {
	Rank  int `json:"rank"`
	Total int `json:"total"`
}

func (a *CapabilityAdapter) getOwnItemPricing(ctx context.Context, accountRef domain.ProviderAccountRef, token, itemID string) (domain.OwnItemPricing, error) {
	var response mlOwnItemPricingResponse
	// ML /items/{id}/sale_price 400s without context; channel_marketplace selects
	// standard marketplace pricing (FINDING-M02-LIVE-2, D-86; official docs
	// precio-por-cantidad/kits-virtuales).
	query := url.Values{}
	query.Set("context", "channel_marketplace")
	path := "/items/" + url.PathEscape(itemID) + "/sale_price?" + query.Encode()
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &response); err != nil {
		return domain.OwnItemPricing{}, mapPricingReaderError(err)
	}

	amount := decimalString(response.Amount)
	regularAmount := decimalString(response.RegularAmount)
	if err := validateSiblingMoney("amount", amount, response.Currency); err != nil {
		return domain.OwnItemPricing{}, err
	}
	if err := validateSiblingMoney("regular_amount", regularAmount, response.Currency); err != nil {
		return domain.OwnItemPricing{}, err
	}
	return domain.OwnItemPricing{
		ItemID:        itemID,
		Amount:        amount,
		Currency:      strings.TrimSpace(response.Currency),
		RegularAmount: regularAmount,
		FetchedAt:     a.now(),
	}, nil
}

func (a *CapabilityAdapter) getPriceToWin(ctx context.Context, accountRef domain.ProviderAccountRef, token, itemID string) (domain.PriceToWin, error) {
	var response mlPriceToWinResponse
	path := "/items/" + url.PathEscape(itemID) + "/price_to_win?version=v2"
	if err := a.doJSON(ctx, accountRef, token, http.MethodGet, path, nil, &response); err != nil {
		return domain.PriceToWin{}, mapPricingReaderError(err)
	}

	currentPrice, err := providerMoney("current_price", response.CurrentPrice, response.Currency)
	if err != nil {
		return domain.PriceToWin{}, err
	}
	targetPrice, err := providerMoney("target_price", response.TargetPrice, response.Currency)
	if err != nil {
		return domain.PriceToWin{}, err
	}
	var position *domain.MarketPosition
	if response.Position != nil {
		position = &domain.MarketPosition{Rank: response.Position.Rank, Total: response.Position.Total}
	}
	return domain.PriceToWin{
		ItemID:       itemID,
		Status:       strings.TrimSpace(response.Status),
		CurrentPrice: currentPrice,
		TargetPrice:  targetPrice,
		Position:     position,
		FetchedAt:    a.now(),
	}, nil
}

func decimalString(value *json.Number) *string {
	if value == nil {
		return nil
	}
	result := value.String()
	return &result
}

func providerMoney(field string, value *json.Number, currency string) (*domain.Money, error) {
	amount := decimalString(value)
	if amount == nil {
		return nil, nil
	}
	money := &domain.Money{Amount: *amount, Currency: strings.TrimSpace(currency)}
	if err := domain.ValidateMoney(field, money); err != nil {
		return nil, err
	}
	return money, nil
}

func validateSiblingMoney(field string, amount *string, currency string) error {
	if amount == nil {
		return nil
	}
	return domain.ValidateMoney(field, &domain.Money{Amount: *amount, Currency: strings.TrimSpace(currency)})
}

func mapPricingReaderError(err error) error {
	if err == nil {
		return nil
	}
	switch domain.ErrorCodeOf(err) {
	case domain.ErrCodeProviderAuth:
		return fmt.Errorf("%w: %w", domain.ErrUnauthorized, err)
	case domain.ErrCodeProviderInvalidReference:
		return fmt.Errorf("%w: %w", domain.ErrNotFound, err)
	case domain.ErrCodeProviderRateLimited:
		return fmt.Errorf("%w: %w", domain.ErrRateLimited, err)
	case domain.ErrCodeProviderTransient:
		return fmt.Errorf("%w: %w", domain.ErrProviderUnavailable, err)
	default:
		return err
	}
}
