package domain

import "time"

// OrderDecomposition carries the per-order cost/fee breakdown behind
// retorno_liquido/margem_pct. Every amount is a *float64 pointer: nil means
// the component could not be honestly sourced (ADR-17: unknown != zero,
// never a fabricated 0). ComponentesDesconhecidos names every source
// component that is unknown, so a "—" in the UI is always explained.
type OrderDecomposition struct {
	Comissao                 *float64
	TaxaFixa                 *float64
	Frete                    *float64
	Imposto                  *float64
	Difal                    *float64
	TarifaFull               *float64
	Custo                    *float64
	MargemValor              *float64
	MargemPct                *float64
	ComponentesDesconhecidos []string
}

// OrderDifal carries the per-order DIFAL fact. All fields nil means unknown
// (ADR-17) — never a fabricated amount/route/date/paid flag.
type OrderDifal struct {
	Amount  *float64
	UFRoute *string
	DueDate *time.Time
	Paid    *bool
}

// OrderProfitability is the composite value a Decomposer returns and
// EnrichedOrder carries: the top-level retorno/margem plus the full
// decomposition and DIFAL breakdown.
type OrderProfitability struct {
	RetornoLiquido *float64
	MargemPct      *float64
	Decomposition  OrderDecomposition
	Difal          OrderDifal
}

// unknownDecompositionComponents names, in this exact order, every
// cost/fee source component the honest-empty path cannot source.
// margem_valor/margem_pct/retorno are derived FROM these components, so
// they are not themselves listed as source components.
var unknownDecompositionComponents = []string{
	"comissao",
	"taxa_fixa",
	"frete",
	"imposto",
	"difal",
	"tarifa_full",
	"custo",
}

// UnknownOrderProfitability is the honest-empty value emitted whenever no
// Decomposer is wired (C1): every pointer is nil so the UI renders "—",
// explained by ComponentesDesconhecidos naming all 7 unknown source
// components (ADR-17: unknown != zero, never fabricated).
func UnknownOrderProfitability() OrderProfitability {
	components := make([]string, len(unknownDecompositionComponents))
	copy(components, unknownDecompositionComponents)
	return OrderProfitability{
		Decomposition: OrderDecomposition{
			ComponentesDesconhecidos: components,
		},
	}
}
