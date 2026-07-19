package domain

import "math/big"

// solveMaxCents caps the search at R$ 1.000.000,00. margem_pct approaches its
// ceiling asymptotically, so a finite cap is needed; a target unreached within
// it is treated as UNREACHABLE alongside the analytic ceiling.
const solveMaxCents int64 = 100_000_000

// SolveInput is a DecomposeInput with the preço replaced by a target
// margem_pct — the bidirectional (margem → preço) direction. Every other
// component is resolved exactly as for Decompose.
type SolveInput struct {
	TargetMargemPct string

	ComissaoPct  string
	AliquotaPct  string
	Modalidade   Modalidade
	TarifaFull   *Money
	FreteProduto *Money
	Custo        *Money

	DifalEnabled bool
	DestinoUF    string
	EfetivoPct   string

	// TaxaFixaLimiarCents overrides the taxa_fixa/frete step (in cents) for
	// this solve. 0 ⇒ default defaultTaxaFixaLimiarCents (79,00).
	TaxaFixaLimiarCents int64
}

// SolveResult is the outcome of SolveTargetPrice. Exactly one condition holds:
//   - Reached with a non-nil Preco: the target is attainable.
//   - !Reached with CeilingPct: UNREACHABLE_TARGET — the target is at or above
//     the analytic ceiling (or unattainable within the cap); CeilingPct is the
//     attainable margem_pct ceiling to surface (mapped to HTTP 200 in S7).
//   - Desconhecidos non-empty: BLOCKING — inputs leave margem unknown (ADR-17);
//     the solver reports the missing components (custo_erp / tarifa_full /
//     difal — segment-independent), never a preço or a ceiling.
//   - FreteDesconhecido true: the target is only attainable in the ≥limiar
//     segment, where produto frete is required but unknown (ADR-17). Distinct
//     from Desconhecidos (structural, segment-independent) and from
//     CeilingPct (target unreachable regardless of frete).
type SolveResult struct {
	Preco             *string
	Reached           bool
	CeilingPct        string
	Desconhecidos     []string
	FreteDesconhecido bool
}

// limiarCents resolves the effective taxa_fixa/frete step for this solve.
func (in SolveInput) limiarCents() int64 {
	if in.TaxaFixaLimiarCents <= 0 {
		return defaultTaxaFixaLimiarCents
	}
	return in.TaxaFixaLimiarCents
}

// SolveTargetPrice finds the cheapest 2dp preço whose Decompose margem_pct
// equals TargetMargemPct exactly. retorno(preço) is piecewise over the preço=79
// step, so the search runs per monotone segment (low first for the cheapest
// price, then high) rather than a single bisection that the discontinuity would
// break. The ceiling (100 − comissão − aliquota − difal) is the margem_pct
// asymptote as preço→∞.
func SolveTargetPrice(in SolveInput) SolveResult {
	limiar := in.limiarCents()
	if unk := in.structuralUnknowns(); len(unk) > 0 {
		return SolveResult{Desconhecidos: unk}
	}

	ceiling := in.ceilingPct()
	// target at or above the asymptote can never be reached.
	if cmpPct(in.TargetMargemPct, ceiling) >= 0 {
		return SolveResult{CeilingPct: ceiling}
	}

	// low segment (preço < limiar): frete is 0 here, so always evaluable even
	// when produto frete is unknown — this is the cheapest price when it reaches.
	if cmpPct(in.margemAt(limiar-1), in.TargetMargemPct) >= 0 {
		if p, ok := in.searchSegment(1, limiar-1); ok {
			return SolveResult{Preco: &p, Reached: true}
		}
	}
	// target needs the high segment (preço ≥ limiar), where produto frete is
	// consulted. If it is unknown, block ONLY this segment (ADR-17) — do not
	// fabricate a frete of 0.
	if in.FreteProduto == nil {
		return SolveResult{FreteDesconhecido: true}
	}
	if p, ok := in.searchSegment(limiar, solveMaxCents); ok {
		return SolveResult{Preco: &p, Reached: true}
	}
	return SolveResult{CeilingPct: ceiling}
}

// ceilingPct is the margem_pct asymptote: fixed costs (taxa_fixa/frete/
// tarifa_full/custo) wash out as preço→∞, leaving 100 − comissão − aliquota −
// difal (difal only when applied).
func (in SolveInput) ceilingPct() string {
	ceil := big.NewRat(100, 1)
	ceil.Sub(ceil, mustRat(in.ComissaoPct))
	ceil.Sub(ceil, mustRat(in.AliquotaPct))
	if in.difalApplied() {
		ceil.Sub(ceil, mustRat(in.EfetivoPct))
	}
	return FormatRatHalfUp(ceil, 2)
}

func (in SolveInput) difalApplied() bool {
	return in.DifalEnabled && in.DestinoUF != "" && in.EfetivoPct != ""
}

// structuralUnknowns returns the components that make margem unknown regardless
// of preço (custo nil, full+tarifa_full nil, difal enabled+destino unknown).
// frete_produto nil is NO LONGER caught here — it is now a segment-conditional
// state handled by SolveTargetPrice. Evaluated at limiar-1, a low-segment
// preço where frete is 0.
func (in SolveInput) structuralUnknowns() []string {
	return in.margemDecompose(in.limiarCents() - 1).ComponentesDesconhecidos
}

func (in SolveInput) searchSegment(loCents, hiCents int64) (string, bool) {
	// margem_pct is increasing within a segment; bail if even the top is short.
	if cmpPct(in.margemAt(hiCents), in.TargetMargemPct) < 0 {
		return "", false
	}
	lo, hi := loCents, hiCents
	for lo < hi {
		mid := lo + (hi-lo)/2
		if cmpPct(in.margemAt(mid), in.TargetMargemPct) < 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	// lo is the smallest cents whose margem_pct ≥ target; it matches exactly
	// unless a one-cent step skipped over the target band (only at tiny preços,
	// outside the money range these inputs target).
	if cmpPct(in.margemAt(lo), in.TargetMargemPct) == 0 {
		return centsToStr(lo), true
	}
	return "", false
}

// margemAt returns the margem_pct string at a candidate preço (in cents). The
// caller guarantees no structural unknowns, so within a searched segment this
// is always non-nil; "" is returned only defensively.
func (in SolveInput) margemAt(cents int64) string {
	d := in.margemDecompose(cents)
	if d.MargemPct == nil {
		return ""
	}
	return *d.MargemPct
}

func (in SolveInput) margemDecompose(cents int64) Decomposition {
	return decomposeWithLimiar(DecomposeInput{
		Preco: centsToStr(cents), ComissaoPct: in.ComissaoPct, AliquotaPct: in.AliquotaPct,
		Modalidade: in.Modalidade, TarifaFull: in.TarifaFull,
		FreteProduto: in.FreteProduto, Custo: in.Custo,
		DifalEnabled: in.DifalEnabled, DestinoUF: in.DestinoUF, EfetivoPct: in.EfetivoPct,
	}, big.NewRat(in.limiarCents(), 100))
}

// centsToStr renders integer cents as a 2dp decimal string via the single
// rounding routine (exact — cents are already 2dp).
func centsToStr(cents int64) string {
	return FormatRatHalfUp(big.NewRat(cents, 100), 2)
}

// cmpPct compares two 2dp percent strings numerically (-1/0/+1).
func cmpPct(a, b string) int {
	return mustRat(a).Cmp(mustRat(b))
}
