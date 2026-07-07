package domain

type ProductCandidate struct {
	Codprod      int
	Produto      string
	EAN          string
	Reference    string
	QualityFlags []QualityFlag
}
