package domain

import "time"

type SalesHistoryWindow struct {
	Start time.Time
	End   time.Time
}

type SalesHistoryEntry struct {
	ProductID       int
	ProductGroupID  *int
	SaleDate        time.Time
	Quantity        float64
	NetAmount       float64
	IsConfirmed     bool
	PriceTableID    *int
	OperationTypeID *int
}

type SalesHistory struct {
	Entries      []SalesHistoryEntry
	Source       SourceMetadata
	QualityFlags []QualityFlag
}
