package domain

import "time"

type TaxInputs struct {
	Codprod         int
	SaleDate        string
	ICMSAmount      *float64
	IPIAmount       *float64
	PISAmount       *float64
	COFINSAmount    *float64
	QualityFlags    []QualityFlag
	SourceFetchedAt time.Time
}
