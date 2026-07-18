package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestMarketReadShapes(t *testing.T) {
	t.Parallel()

	fetchedAt := time.Now().UTC().Truncate(time.Microsecond)
	amount := "89.90"
	regularAmount := "120.00"
	pricing := OwnItemPricing{
		ItemID:        "MLB-1",
		Amount:        &amount,
		Currency:      "BRL",
		RegularAmount: &regularAmount,
		FetchedAt:     fetchedAt,
	}
	if pricing.ItemID != "MLB-1" || *pricing.Amount != "89.90" || *pricing.RegularAmount != "120.00" || pricing.Currency != "BRL" || !pricing.FetchedAt.Equal(fetchedAt) {
		t.Fatalf("pricing = %#v", pricing)
	}

	current := &Money{Amount: "120.00", Currency: "BRL"}
	target := &Money{Amount: "99.50", Currency: "BRL"}
	position := &MarketPosition{Rank: 2, Total: 5}
	priceToWin := PriceToWin{
		ItemID:       "MLB-1",
		Status:       "winning",
		CurrentPrice: current,
		TargetPrice:  target,
		Position:     position,
		FetchedAt:    fetchedAt,
	}
	if priceToWin.TargetPrice == nil || priceToWin.TargetPrice.Amount == priceToWin.CurrentPrice.Amount || !reflect.DeepEqual(priceToWin.Position, position) {
		t.Fatalf("price to win = %#v", priceToWin)
	}
}
