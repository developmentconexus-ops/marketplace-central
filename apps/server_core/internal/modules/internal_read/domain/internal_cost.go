package domain

import "time"

type CostBasis string

const (
	CostBasisCUSSEMICM CostBasis = "cussemicm"
)

type CostAsOfPolicy struct {
	CompanyID   int
	EffectiveAt time.Time
	Basis       CostBasis
}

type CostAsOf struct {
	ProductID    int
	CompanyID    int
	Basis        CostBasis
	EffectiveAt  time.Time
	Amount       *float64
	Source       SourceMetadata
	QualityFlags []QualityFlag
}
