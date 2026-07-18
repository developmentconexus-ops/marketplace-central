package domain

import "time"

type OwnItemPricing struct {
	ItemID        string    `json:"item_id"`
	Amount        *string   `json:"amount"`
	Currency      string    `json:"currency"`
	RegularAmount *string   `json:"regular_amount"`
	FetchedAt     time.Time `json:"fetched_at"`
}

type PriceToWin struct {
	ItemID       string          `json:"item_id"`
	Status       string          `json:"status"`
	CurrentPrice *Money          `json:"current_price"`
	TargetPrice  *Money          `json:"target_price"`
	Position     *MarketPosition `json:"position"`
	FetchedAt    time.Time       `json:"fetched_at"`
}

type MarketPosition struct {
	Rank  int `json:"rank"`
	Total int `json:"total"`
}
