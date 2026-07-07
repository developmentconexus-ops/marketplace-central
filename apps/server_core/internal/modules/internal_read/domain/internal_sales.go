package domain

import "time"

type SalesHistoryEntry struct {
	Codprod   int
	SaleDate  string
	Quantity  float64
	NetAmount float64
}

type SalesHistory struct {
	Entries         []SalesHistoryEntry
	QualityFlags    []QualityFlag
	SourceFetchedAt time.Time
}
