package domain

import "time"

type CurrentPrice struct {
	Codprod         int
	Codtab          int
	Codlocal        int
	Price           *float64
	QualityFlags    []QualityFlag
	SourceFetchedAt time.Time
}
