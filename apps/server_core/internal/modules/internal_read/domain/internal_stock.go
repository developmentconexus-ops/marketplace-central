package domain

const (
	SellableStockScopeResale StockScope = "resale"
	SellableStockScopeCustom StockScope = "custom"
)

type StockScope string

type SellableStockPolicy struct {
	CompanyIDs          []int
	LocationIDs         []int
	ExcludedLocationIDs []int
	Formula             string
	Scope               StockScope
}

func DefaultSellableStockPolicy() SellableStockPolicy {
	return SellableStockPolicy{
		CompanyIDs:          []int{1, 2},
		LocationIDs:         []int{10101},
		ExcludedLocationIDs: []int{10108},
		Formula:             "SUM(ESTOQUE - RESERVADO)",
		Scope:               SellableStockScopeResale,
	}
}

type SellableStock struct {
	ProductID    int
	Quantity     *float64
	Policy       SellableStockPolicy
	Source       SourceMetadata
	QualityFlags []QualityFlag
}
