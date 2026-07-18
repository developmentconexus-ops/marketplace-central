package domain

import "errors"

type CatalogOffer struct {
	SellerID     string `json:"seller_id"`
	Price        *Money `json:"price"`
	Condition    string `json:"condition"`
	ShippingMode string `json:"shipping_mode"`
}

var ErrCatalogOffersUnavailable = errors.New("catalog offers unavailable")
