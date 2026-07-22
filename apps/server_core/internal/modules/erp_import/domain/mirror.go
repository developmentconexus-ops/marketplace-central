package domain

import "time"

type MirrorProduct struct {
	CodigoProduto  string
	Descricao      *string
	Referencia     *string
	EAN            *string
	Marca          *string
	GrupoCodigo    *string
	GrupoDescricao *string
	NCM            *string
	Custo          *string
	PrecoVenda     *string
	EstoqueTotal   *string
	UpdatedAt      time.Time
}
