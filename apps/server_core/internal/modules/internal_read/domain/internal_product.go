package domain

import "time"

type ProductCandidate struct {
	Codprod         int
	Produto         string
	EAN             *string
	Reference       *string
	CodgrupoProd    *int
	GrupoProduto    *string
	Codmarca        *int
	Marca           *string
	QualityFlags    []QualityFlag
	SourceFetchedAt time.Time
}
