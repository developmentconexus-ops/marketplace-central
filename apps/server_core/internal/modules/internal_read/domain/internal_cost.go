package domain

type CostAsOf struct {
	Codprod      int
	Codemp       int
	SaleDate     string
	CUSSEMICM    *float64
	QualityFlags []QualityFlag
}
