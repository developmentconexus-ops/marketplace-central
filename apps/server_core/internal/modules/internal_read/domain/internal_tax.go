package domain

import "time"

type TaxPolicy struct {
	EffectiveAt   time.Time
	IncidenceCode int
}

func DefaultTaxPolicy(now time.Time) TaxPolicy {
	return TaxPolicy{
		EffectiveAt:   now,
		IncidenceCode: 0,
	}
}

type TaxInputs struct {
	ProductID     int
	EffectiveAt   time.Time
	IncidenceCode int
	ICMSAmount    *float64
	IPIAmount     *float64
	PISAmount     *float64
	COFINSAmount  *float64
	Source        SourceMetadata
	QualityFlags  []QualityFlag
}
