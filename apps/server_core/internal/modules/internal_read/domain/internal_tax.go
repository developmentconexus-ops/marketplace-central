package domain

import "time"

type TaxPolicy struct {
	EffectiveAt   time.Time
	IncidenceCode int
	Source        TaxSourceIdentity
}

// TaxSourceIdentity identifies exactly one Oracle sale line. Both values must
// be positive before tax amounts can be attributed to a marketplace item.
type TaxSourceIdentity struct {
	DocumentID int
	LineNumber int
}

func (s TaxSourceIdentity) Verified() bool {
	return s.DocumentID > 0 && s.LineNumber > 0
}

func DefaultTaxPolicy(now time.Time) TaxPolicy {
	return TaxPolicy{
		EffectiveAt:   now,
		IncidenceCode: 0,
	}
}

type TaxInputs struct {
	ProductID      int
	EffectiveAt    time.Time
	IncidenceCode  int
	SourceIdentity TaxSourceIdentity
	NCM            *string
	ICMSAmount     *float64
	IPIAmount      *float64
	PISAmount      *float64
	COFINSAmount   *float64
	Source         SourceMetadata
	QualityFlags   []QualityFlag
}
