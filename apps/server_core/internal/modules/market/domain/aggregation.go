package domain

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"
)

type pricedOffer struct {
	offer ValidatedOffer
	value *big.Rat
}

func ComputeMarketAggregate(productID string, offers []ValidatedOffer, source MarketPriceSource, fetchedAt, computedAt time.Time) (MarketAggregate, error) {
	qualified := make([]pricedOffer, 0, len(offers))
	for _, offer := range offers {
		if offer.Price == nil || offer.Price.Currency != "BRL" || strings.TrimSpace(offer.Condition) != "new" || offer.MatchStatus != MatchStatusAccept {
			continue
		}
		value, ok := new(big.Rat).SetString(offer.Price.Amount)
		if !ok {
			return MarketAggregate{}, fmt.Errorf("offer price amount %q is not a decimal", offer.Price.Amount)
		}
		qualified = append(qualified, pricedOffer{offer: offer, value: value})
	}

	bySeller := make(map[string]pricedOffer, len(qualified))
	for _, candidate := range qualified {
		current, exists := bySeller[candidate.offer.SellerID]
		if !exists || candidate.value.Cmp(current.value) < 0 {
			bySeller[candidate.offer.SellerID] = candidate
		}
	}
	deduped := make([]pricedOffer, 0, len(bySeller))
	for _, offer := range bySeller {
		deduped = append(deduped, offer)
	}
	sort.SliceStable(deduped, func(i, j int) bool {
		cmp := deduped[i].value.Cmp(deduped[j].value)
		if cmp != 0 {
			return cmp < 0
		}
		return deduped[i].offer.SellerID < deduped[j].offer.SellerID
	})

	input := MarketAggregateInput{ProductID: productID, NOffers: len(qualified), NSellers: len(deduped), Source: source, FetchedAt: fetchedAt, ComputedAt: computedAt}
	if len(deduped) == 0 {
		input.Status = MarketAggregateStatusNoPriceEvidence
		return NewMarketAggregate(input)
	}
	input.MinValid = deduped[0].offer.Price
	input.Median = medianMoney(deduped)
	input.Status = MarketAggregateStatusOK
	if len(deduped) < 5 {
		input.Status = MarketAggregateStatusInsufficientMarket
	}
	return NewMarketAggregate(input)
}

func medianMoney(offers []pricedOffer) *Money {
	middle := len(offers) / 2
	if len(offers)%2 == 1 {
		return offers[middle].offer.Price
	}
	mean := new(big.Rat).Add(offers[middle-1].value, offers[middle].value)
	mean.Quo(mean, big.NewRat(2, 1))
	return &Money{Amount: formatRatHalfUp(mean, 2), Currency: "BRL"}
}

func formatRatHalfUp(value *big.Rat, places int) string {
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(places)), nil)
	scaledNumerator := new(big.Int).Mul(value.Num(), scale)
	quotient, remainder := new(big.Int).QuoRem(scaledNumerator, value.Denom(), new(big.Int))
	if new(big.Int).Mul(remainder, big.NewInt(2)).Cmp(value.Denom()) >= 0 {
		quotient.Add(quotient, big.NewInt(1))
	}
	digits := quotient.String()
	for len(digits) <= places {
		digits = "0" + digits
	}
	return digits[:len(digits)-places] + "." + digits[len(digits)-places:]
}
