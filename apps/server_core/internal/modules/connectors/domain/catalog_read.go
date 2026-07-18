package domain

import (
	"errors"
	"time"
)

type CatalogOffer struct {
	SellerID     string `json:"seller_id"`
	Price        *Money `json:"price"`
	Condition    string `json:"condition"`
	ShippingMode string `json:"shipping_mode"`
}

var ErrCatalogOffersUnavailable = errors.New("catalog offers unavailable")

type CatalogSearchResult struct {
	Products  []CatalogCandidate `json:"products"`
	FetchedAt time.Time          `json:"fetched_at"`
}

type CatalogCandidate struct {
	CatalogProductID string             `json:"catalog_product_id"`
	Title            string             `json:"title"`
	Attrs            []CatalogAttribute `json:"attrs"`
}

type CatalogProduct struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	BuyBoxWinner *BuyBoxWinner      `json:"buy_box_winner"`
	Attrs        []CatalogAttribute `json:"attrs"`
	FetchedAt    time.Time          `json:"fetched_at"`
}

type BuyBoxWinner struct {
	Price    *Money `json:"price"`
	SellerID string `json:"seller_id"`
}

type CatalogAttribute struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ValueName *string `json:"value_name"`
}
