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

// ProfitabilityInputs carries the already-resolved real facts the honest
// decomposition is built from. A nil field means that component could not be
// honestly sourced (ADR-17). Total is the order revenue.
type ProfitabilityInputs struct {
	Total    *float64
	Comissao *float64
	Custo    *float64
	Frete    *float64
	Imposto  *float64
}

// BuildProfitability assembles an OrderProfitability from real facts in hand.
// Comissao/Custo/Frete/Imposto are surfaced verbatim when known. TaxaFixa,
// Difal and TarifaFull need a pricing engine not yet wired, so they are always
// honest-unknown (nil) and named in ComponentesDesconhecidos. Margem
// (valor + pct) and RetornoLiquido derive ONLY when Total, Comissao, Frete,
// Imposto and Custo are ALL known; any unknown input yields nil margins —
// never a partial fabricated margin (ADR-17).
func BuildProfitability(in ProfitabilityInputs) OrderProfitability {
	dec := OrderDecomposition{
		Comissao: in.Comissao,
		Frete:    in.Frete,
		Imposto:  in.Imposto,
		Custo:    in.Custo,
	}
	var unknown []string
	if in.Comissao == nil {
		unknown = append(unknown, "comissao")
	}
	unknown = append(unknown, "taxa_fixa")
	if in.Frete == nil {
		unknown = append(unknown, "frete")
	}
	if in.Imposto == nil {
		unknown = append(unknown, "imposto")
	}
	unknown = append(unknown, "difal", "tarifa_full")
	if in.Custo == nil {
		unknown = append(unknown, "custo")
	}
	dec.ComponentesDesconhecidos = unknown

	if in.Total != nil && in.Comissao != nil && in.Frete != nil && in.Imposto != nil && in.Custo != nil {
		margem := *in.Total - *in.Comissao - *in.Frete - *in.Imposto - *in.Custo
		dec.MargemValor = &margem
		var pctPtr *float64
		if *in.Total != 0 {
			pct := margem / *in.Total // FRACTION (0.18 = 18%); FE formatPercent multiplies by 100 — do NOT ×100 here
			pctPtr = &pct
			dec.MargemPct = &pct
		}
		return OrderProfitability{RetornoLiquido: &margem, MargemPct: pctPtr, Decomposition: dec}
	}
	return OrderProfitability{Decomposition: dec}
}
