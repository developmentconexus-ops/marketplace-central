package domain

import "time"

// OrderDecomposition carries the per-order cost/fee breakdown behind
// retorno_liquido/margem_pct. Every amount is a *float64 pointer: nil means
// the component could not be honestly sourced (ADR-17: unknown != zero,
// never a fabricated 0). ComponentesDesconhecidos names every source
// component that is unknown, so a "—" in the UI is always explained.
type OrderDecomposition struct {
	Comissao *float64
	TaxaFixa *float64
	Frete    *float64
	// Imposto is the D-38 legacy regime-aliquota field, kept for the
	// simulator surfaces that still read it. T5 (P2.b) replaces it end to
	// end with the per-item D-41 path (ICMSSaida/Difal/PisCofins/
	// RestituicaoST below) — this field is never populated by
	// BuildProfitability anymore; it survives only as inert pass-through.
	Imposto *float64
	Difal   *float64
	// ICMSSaida is the D-41 per-item ICMS de saída, summed across the
	// order's items (orders/adapters/pricingtax, T5). nil is honest-unknown
	// (ADR-17), a real amount is a real amount — including an explicit 0.
	ICMSSaida *float64
	// PisCofins is the D-41 PIS/COFINS débito, summed across items. Same nil
	// rule as ICMSSaida.
	PisCofins  *float64
	TarifaFull *float64
	Custo      *float64
	// RestituicaoST is the ICMS-ST restituição credit, summed across items —
	// the one component that ADDS to margem instead of subtracting
	// (BuildProfitability below). nil is honest-unknown.
	RestituicaoST            *float64
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
	"icms_saida",
	"difal",
	"pis_cofins",
	"tarifa_full",
	"custo",
	"restituicao_st",
}

// UnknownOrderProfitability is the honest-empty value emitted whenever no
// Decomposer is wired (C1): every pointer is nil so the UI renders "—",
// explained by ComponentesDesconhecidos naming all 9 unknown source
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
	// Difal is the destination-state tax differential. A pointer to 0 is an
	// explicit zero (the tenant has DIFAL switched off), which is a different
	// fact from nil (no honest way to work it out) — only nil blocks the margin.
	Difal *float64
	// ICMSSaida, PisCofins and RestituicaoST are the D-41 per-item ICMS
	// components (T5), already summed across the order's items by the tax
	// reader. RestituicaoST is a CREDIT: it adds to margem instead of
	// subtracting, same sign convention as pricing/domain/decompose.go.
	ICMSSaida     *float64
	PisCofins     *float64
	RestituicaoST *float64
}

// BuildProfitability assembles an OrderProfitability from real facts in hand.
// Comissao/Custo/Frete/ICMSSaida/Difal/PisCofins/RestituicaoST are surfaced
// verbatim when known. TaxaFixa and TarifaFull are already inside the
// marketplace's own sale fee for a sold order, so breaking them back out
// needs a pricing engine that is not wired; they stay honest-unknown (nil)
// and named in ComponentesDesconhecidos. Margem (valor + pct) and
// RetornoLiquido derive ONLY when Total, Comissao, Frete, ICMSSaida, Difal,
// PisCofins, Custo and RestituicaoST are ALL known; any unknown input yields
// nil margins — never a partial fabricated margin (ADR-17).
func BuildProfitability(in ProfitabilityInputs) OrderProfitability {
	dec := OrderDecomposition{
		Comissao:      in.Comissao,
		Frete:         in.Frete,
		Difal:         in.Difal,
		ICMSSaida:     in.ICMSSaida,
		PisCofins:     in.PisCofins,
		Custo:         in.Custo,
		RestituicaoST: in.RestituicaoST,
	}
	var unknown []string
	if in.Comissao == nil {
		unknown = append(unknown, "comissao")
	}
	unknown = append(unknown, "taxa_fixa")
	if in.Frete == nil {
		unknown = append(unknown, "frete")
	}
	if in.ICMSSaida == nil {
		unknown = append(unknown, "icms_saida")
	}
	if in.Difal == nil {
		unknown = append(unknown, "difal")
	}
	if in.PisCofins == nil {
		unknown = append(unknown, "pis_cofins")
	}
	unknown = append(unknown, "tarifa_full")
	if in.Custo == nil {
		unknown = append(unknown, "custo")
	}
	if in.RestituicaoST == nil {
		unknown = append(unknown, "restituicao_st")
	}
	dec.ComponentesDesconhecidos = unknown

	if in.Total != nil && in.Comissao != nil && in.Frete != nil && in.ICMSSaida != nil && in.Difal != nil && in.PisCofins != nil && in.Custo != nil && in.RestituicaoST != nil {
		margem := *in.Total - *in.Comissao - *in.Frete - *in.ICMSSaida - *in.Difal - *in.PisCofins - *in.Custo + *in.RestituicaoST
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
