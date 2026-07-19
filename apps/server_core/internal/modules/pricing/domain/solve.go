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
	// when produto frete is unknown, and being the cheapest segment it is tried
	// first. searchSegment scans it EXHAUSTIVELY — Decompose margem_pct is not
	// monotone in preço (2dp rounding of comissão/imposto/difal makes it wiggle
	// downward by sub-cent amounts), so a rounded-value bisection skips valid
	// prices and mis-reports reachable targets (FINDING-P6-SOLVER).
	if p, ok := in.searchSegment(1, limiar-1); ok {
		return SolveResult{Preco: &p, Reached: true}
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
	// probe a valid low-segment preço; floor at 1 cent so a degenerate limiar of
	// 1 (limiar-1 == 0) cannot divide by preço 0 in Decompose.
	probe := in.limiarCents() - 1
	if probe < 1 {
		probe = 1
	}
	return in.margemDecompose(probe).ComponentesDesconhecidos
}

// lowSegmentSpanCents is the widest span searchSegment scans cent-by-cent. The
// low segment (below the taxa_fixa limiar, ≤ 7899 cents) is always narrower, so
// it is scanned exhaustively; the unbounded high segment is bracketed first.
const lowSegmentSpanCents int64 = 200_000

// highScanCapCents caps the high-segment linear scan after bracketing so a
// pathological near-ceiling target cannot make it unbounded. A match beyond the
// bracketed+capped window is reported UNREACHABLE (honest) — never mis-solved.
const highScanCapCents int64 = 4_000_000

// searchSegment returns the cheapest 2dp preço (in cents) within [loCents,
// hiCents] whose Decompose margem_pct equals TargetMargemPct EXACTLY, and
// whether such a preço exists.
//
// Decompose rounds comissão/imposto/difal to 2dp before subtracting, so its
// margem_pct is only piecewise-increasing with sub-cent downward wiggles; a
// binary search over the rounded value is unsound (FINDING-P6-SOLVER: it skips
// valid prices and mis-reports reachable targets as SEM_FRETE/UNREACHABLE). The
// EXACT margem_pct (100·k − 100·F/preço) IS strictly increasing, so it is used
// only to BRACKET the candidate window; the exact-match test still runs the
// real Decompose so the returned price honors the frozen contract exactly.
func (in SolveInput) searchSegment(loCents, hiCents int64) (string, bool) {
	if loCents > hiCents {
		return "", false
	}
	limiar := in.limiarCents()
	scanLo, scanHi := loCents, hiCents
	if hiCents-loCents > lowSegmentSpanCents {
		// |Decompose.pct − exact.pct| < 2/preço + 0.02 pp (2dp rounding of ≤3
		// pct components, then of the ratio, with margin); preço ≥ loCents/100
		// bounds it. Any exact-target match lies where the monotone exact value
		// is within that wiggle of the target — bracket that band by bisection.
		wig := new(big.Rat).Quo(big.NewRat(200, 100), big.NewRat(loCents, 100))
		wig.Add(wig, big.NewRat(2, 100))
		tgt := mustRat(in.TargetMargemPct)
		scanLo = in.firstCentExactAtLeast(loCents, hiCents, limiar, new(big.Rat).Sub(tgt, wig))
		if scanLo > hiCents {
			return "", false // even the cap doesn't reach the target band
		}
		hiExceed := in.firstCentExactAtLeast(loCents, hiCents, limiar, new(big.Rat).Add(tgt, wig))
		scanHi = hiExceed - 1
		if scanHi > hiCents {
			scanHi = hiCents
		}
		if scanHi-scanLo > highScanCapCents {
			scanHi = scanLo + highScanCapCents
		}
	}
	for c := scanLo; c <= scanHi; c++ {
		// numeric compare: target may arrive unnormalized ("15" vs "15.00").
		if m := in.margemAt(c); m != "" && cmpPct(m, in.TargetMargemPct) == 0 {
			return centsToStr(c), true
		}
	}
	return "", false
}

// firstCentExactAtLeast returns the smallest cents in [lo,hi] whose EXACT
// (unrounded, strictly-increasing) margem_pct is ≥ bound, or hi+1 if none.
func (in SolveInput) firstCentExactAtLeast(lo, hi, limiar int64, bound *big.Rat) int64 {
	if in.exactMargemPctRat(hi, limiar).Cmp(bound) < 0 {
		return hi + 1
	}
	for lo < hi {
		mid := lo + (hi-lo)/2
		if in.exactMargemPctRat(mid, limiar).Cmp(bound) < 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// exactMargemPctRat is the UNROUNDED margem_pct at preço=cents for the segment
// selected by limiar: 100·k − 100·F/preço, where k = 1 − Σpct/100 (comissão +
// aliquota + applied difal) and F is the fixed-cost sum (taxa_fixa below limiar
// / produto frete at-or-above, plus tarifa_full when full and custo). Strictly
// increasing in preço because F ≥ 0 — unlike Decompose's rounded margem_pct, so
// it is safe to bisect. Used ONLY to bracket the search window.
func (in SolveInput) exactMargemPctRat(cents, limiar int64) *big.Rat {
	preco := big.NewRat(cents, 100)
	pct := new(big.Rat).Set(cem)
	pct.Sub(pct, mustRat(in.ComissaoPct))
	pct.Sub(pct, mustRat(in.AliquotaPct))
	if in.difalApplied() {
		pct.Sub(pct, mustRat(in.EfetivoPct))
	}
	fixed := new(big.Rat)
	if cents < limiar {
		fixed.Add(fixed, taxaFixaValor)
	} else if in.FreteProduto != nil {
		fixed.Add(fixed, mustRat(in.FreteProduto.Amount))
	}
	if in.Modalidade == ModalidadeFull && in.TarifaFull != nil {
		fixed.Add(fixed, mustRat(in.TarifaFull.Amount))
	}
	if in.Custo != nil {
		fixed.Add(fixed, mustRat(in.Custo.Amount))
	}
	fOverP := new(big.Rat).Quo(fixed, preco)
	fOverP.Mul(fOverP, cem)
	return pct.Sub(pct, fOverP)
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
