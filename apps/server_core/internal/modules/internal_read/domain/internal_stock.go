package domain

import "time"

const (
	StockScopeCodeRevenda          = "revenda"
	StockScopeCodeShowroomExcluded = "showroom_excluded"
	StockScopeCodeCustomPolicy     = "custom_policy"
)

type StockScope struct {
	Companies []int
	Locations []int
	Formula   string
	ScopeCode string
}

func DefaultSellableStockScope() StockScope {
	return StockScope{
		Companies: []int{1, 2},
		Locations: []int{10101},
		Formula:   "SUM(ESTOQUE - RESERVADO)",
		ScopeCode: StockScopeCodeRevenda,
	}
}

type SellableStock struct {
	Codprod         int
	Quantity        *float64
	Scope           StockScope
	QualityFlags    []QualityFlag
	SourceFetchedAt time.Time
}
