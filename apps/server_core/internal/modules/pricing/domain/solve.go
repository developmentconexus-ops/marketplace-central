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

	// ICMSCell is DecomposeInput's field of the same name (decompose.go:63-
	// 75) — forwarded to Decompose via margemDecompose (Task A3) so the
	// solver's margem search runs the exact D-41 formula the decompose
	// direction already uses, never a second tax computation. nil ⇒ legacy
	// AliquotaPct/DifalEnabled/EfetivoPct path (unchanged; Fatia C removes
	// that path once every caller has migrated). Non-nil switches BOTH
	// margemDecompose's tax AND ceilingPct's asymptote onto the cell — never
	// a mix of the two sources within the same solve.
	ICMSCell *ICMSCell

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
// tarifa_full/custo/restituição) wash out as preço→∞, leaving 100 − comissão
// − carga fiscal. The fiscal load comes from ICMSCell when present (D-41
// asymptote, icmsCellAsymptoticRatePct — Task A3); otherwise from the legacy
// AliquotaPct + applied difal (Fatia C removes this branch once the legacy
// path itself is retired). Never mixes the two sources for the same solve.
func (in SolveInput) ceilingPct() string {
	ceil := big.NewRat(100, 1)
	ceil.Sub(ceil, mustRat(in.ComissaoPct))
	if in.ICMSCell != nil {
		ceil.Sub(ceil, in.icmsCellAsymptoticRatePct())
	} else {
		ceil.Sub(ceil, mustRat(in.AliquotaPct))
		if in.difalApplied() {
			ceil.Sub(ceil, mustRat(in.EfetivoPct))
		}
	}
	return FormatRatHalfUp(ceil, 2)
}

// icmsCellAsymptoticRatePct is the D-41 tax load's margem_pct asymptote, in
// percent units (0-100 scale, matching AliquotaPct's convention) — the exact
// analogue of the legacy branch's AliquotaPct+EfetivoPct sum, but derived
// from the SAME formula TaxesForItem (icms.go) uses:
//
//   - icms_saida + difal telescope to aCusto regardless of branch: MG-interno
//     sets icms_saida=P×aCusto, difal=0 (icms.go:196-204); interestadual sets
//     icms_saida=P×aInter, difal=P×aCusto−P×aInter, so the sum is P×aCusto
//     either way (icms.go:205-212) — a_inter cancels out, so UFDestino never
//     needs consulting here, only aCusto.
//   - pis_cofins's base has a flat subtrahend (StRetidoEntrada) that washes
//     out as preço→∞, leaving pisCofinsRate×(1−aBase) — the same aBase
//     (aCusto minus FCP, D-43) TaxesForItem derives from the cell.
//   - isST (codTribST) never consults AliquotaInterna/FcpEmbutido at all
//     (icms.go:145-156: icms_saida=difal=0, base=P−S) — folding aCusto=aBase=0
//     into the general formula reproduces that exactly (0 + pisCofinsRate×1).
//   - restituicao_st is a flat credit like custo/tarifa_full/taxa_fixa and is
//     never added here — it washes out the same way those do.
//
// Callers only reach this with a RESOLVED, unambiguous cell: structural-
// Unknowns already blocks the unresolved/ambíguo case before
// SolveTargetPrice calls ceilingPct (Unknown-ness is price-independent, so a
// probe at one preço decides it for every preço).
//
// COUPLING CONTRACT (Task A3 review Finding 3): this is a hand-re-derivation
// of TaxesForItem's limit as preço→∞, not a call into TaxesForItem itself —
// TaxesForItem only evaluates at a finite, already-2dp-rounded preço, and its
// MAX(0, …) clamp on the pis_cofins base and its round2 calls make any
// finite-sample slope inexact (off by cents), which the solve bracket cannot
// tolerate (it needs the EXACT asymptote, see exactCeilingRat/
// exactMargemPctRat above). A second finite evaluation of TaxesForItem was
// deliberately rejected as the derivation strategy for that reason; this
// closed-form re-encoding is the only exact option that does not touch
// icms.go. TaxesForItem (icms.go:124-234) is the single source of truth this
// function tracks: whenever aCusto, aBase/FCP folding, isST's zeroing, or
// pisCofinsRate change there, re-derive this function's four bullets above
// against the new formula before merging — a drift here reproduces exactly
// the Finding-1 class of bug (asymptote silently diverges from the real
// D-41 formula).
func (in SolveInput) icmsCellAsymptoticRatePct() *big.Rat {
	cell := in.ICMSCell
	isST := cell.CodTrib != nil && *cell.CodTrib == codTribST

	aCusto := new(big.Rat)
	aBase := new(big.Rat)
	if !isST {
		aCusto.Quo(mustRat(*cell.AliquotaInterna), cem)
		fcp := new(big.Rat)
		if cell.FcpEmbutido != nil {
			fcp.Quo(mustRat(*cell.FcpEmbutido), cem)
		}
		aBase.Sub(aCusto, fcp)
	}

	oneMinusABase := new(big.Rat).Sub(big.NewRat(1, 1), aBase)
	rate := new(big.Rat).Mul(pisCofinsRate, oneMinusABase)
	rate.Add(rate, aCusto)
	rate.Mul(rate, cem)
	return rate
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

// windowMarginCents pads the DERIVED high-segment scan half-window (150/gStar
// cents, see searchSegment) to absorb the integer ceil. The window itself is
// computed per-input from the exact ceiling gap, never a fixed cap, so it can
// never silently truncate a reachable target (FINDING-P6-SOLVER-2).
const windowMarginCents int64 = 8

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
		// High segment: bracket the round2(exact.pct)=target crossing by bisection
		// (exact margem_pct is strictly increasing), then scan a window around it
		// against the real Decompose. round2(exact) ≥ target ⇔ exact ≥ target−0.005,
		// so cStar = firstCentExactAtLeast(target−0.005) is that crossing (the
		// −0.005 is the final-round2 half-band, already folded in here).
		//
		// The cheapest Decompose match lies at most 150/gStar cents from cStar,
		// gStar = exactCeiling − (target−0.005): Decompose = round2(exact + ρ) with
		// |ρ| ≤ 1.5/preço pp (2dp rounding of ≤3 cost components), and the exact
		// slope is 10000·F/preço² pp/cent, so |Δpreço| ≤ (1.5/preço)/(10000F/preço²)
		// = 150·preço/(10000F) = 150/gStar cents at the crossing (preço = 10000F/
		// gStar). gStar ≥ 0.005 by construction (target < the 2dp ceiling ≤
		// exactCeiling+0.005), so the window is finite for EVERY input — including
		// sub-0.01 ceiling gaps a >2dp comissão admits, which a fixed cap would
		// truncate (FINDING-P6-SOLVER-2).
		half := big.NewRat(5, 1000) // 0.005 = round2 half-band
		bound := new(big.Rat).Sub(mustRat(in.TargetMargemPct), half)
		gStar := new(big.Rat).Sub(in.exactCeilingRat(), bound)
		if gStar.Sign() <= 0 {
			return "", false // target at/above the exact asymptote — unreachable
		}
		cStar := in.firstCentExactAtLeast(loCents, hiCents, limiar, bound)
		if cStar > hiCents {
			return "", false // crossing beyond the search cap — unreachable
		}
		winCents := ceilRatToCents(new(big.Rat).Quo(big.NewRat(150, 1), gStar)) + windowMarginCents
		scanLo = cStar - winCents
		if scanLo < loCents {
			scanLo = loCents
		}
		scanHi = cStar + winCents
		if scanHi > hiCents {
			scanHi = hiCents
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
// fiscal load) and F is the fixed-cost sum (taxa_fixa below limiar / produto
// frete at-or-above, plus tarifa_full when full and custo). Strictly
// increasing in preço because F ≥ 0 — unlike Decompose's rounded margem_pct, so
// it is safe to bisect. Used ONLY to bracket the search window.
//
// The fiscal load is the SAME source ceilingPct uses (Task A3 review Finding
// 1): icmsCellAsymptoticRatePct when ICMSCell is set, otherwise the legacy
// AliquotaPct + applied difal. Before this fix this always used the legacy
// path even with ICMSCell set, so the bracket was built around the fabricated
// asymptote — the high-segment scan converged on the wrong crossing and
// wrongly reported UNREACHABLE_TARGET for targets only reachable via the real
// cell ceiling (see TestSolveTargetPriceICMSCellBAHighSegment).
func (in SolveInput) exactMargemPctRat(cents, limiar int64) *big.Rat {
	preco := big.NewRat(cents, 100)
	pct := new(big.Rat).Set(cem)
	pct.Sub(pct, mustRat(in.ComissaoPct))
	if in.ICMSCell != nil {
		pct.Sub(pct, in.icmsCellAsymptoticRatePct())
	} else {
		pct.Sub(pct, mustRat(in.AliquotaPct))
		if in.difalApplied() {
			pct.Sub(pct, mustRat(in.EfetivoPct))
		}
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

// exactCeilingRat is the UNROUNDED margem_pct asymptote (100 − comissão −
// fiscal load) as preço→∞ — the exact analogue of ceilingPct, which rounds to
// 2dp. The window derivation needs the unrounded value: the 2dp ceiling gate
// can admit a target only ~0.005 below the true asymptote when a >2dp
// comissão is supplied, so a rounded gap would understate the window.
//
// The fiscal load branches on ICMSCell exactly like ceilingPct (Task A3
// review Finding 1) — icmsCellAsymptoticRatePct when set, otherwise the
// legacy AliquotaPct + applied difal. Never mixes the two sources.
func (in SolveInput) exactCeilingRat() *big.Rat {
	ceil := new(big.Rat).Set(cem)
	ceil.Sub(ceil, mustRat(in.ComissaoPct))
	if in.ICMSCell != nil {
		ceil.Sub(ceil, in.icmsCellAsymptoticRatePct())
	} else {
		ceil.Sub(ceil, mustRat(in.AliquotaPct))
		if in.difalApplied() {
			ceil.Sub(ceil, mustRat(in.EfetivoPct))
		}
	}
	return ceil
}

// ceilRatToCents returns ⌈x⌉ as int64 for a non-negative *big.Rat.
func ceilRatToCents(x *big.Rat) int64 {
	q := new(big.Int)
	m := new(big.Int)
	q.DivMod(x.Num(), x.Denom(), m) // Euclidean: m ≥ 0
	if m.Sign() != 0 {
		q.Add(q, big.NewInt(1))
	}
	return q.Int64()
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
		ICMSCell: in.ICMSCell,
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
