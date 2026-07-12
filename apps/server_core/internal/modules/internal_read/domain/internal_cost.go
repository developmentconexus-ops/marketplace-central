package domain

import "time"

type CostBasis string
type CostAmountScope string

const (
	CostBasisCUSSEMICM     CostBasis       = "cussemicm"
	CostAmountScopePerUnit CostAmountScope = "per_unit"
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
	AmountScope  CostAmountScope
	Source       SourceMetadata
	QualityFlags []QualityFlag
}
