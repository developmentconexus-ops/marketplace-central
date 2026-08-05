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

	// ICMSCell is DecomposeInput's field of the same name (decompose.go) — forwarded to Decompose via margemDecompose (Task A3) so the
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
	// target at or above the asymptote can never be reached — LEGACY PATH
	// ONLY. Task A5 (C1): this short-circuit assumed every summed component
	// is a cost, so the asymptote is approached from BELOW and is a true
	// upper bound. That holds for the legacy path (fixed = taxa_fixa/frete +
	// tarifa_full + custo, all ≥0, no credit term — Fnet = −fixed ≤ 0
	// always). It does NOT hold once ICMSCell is set: RestituicaoST is
	// SUBTRACTED from the summed tax load (decompose.go, "a única
	// linha positiva da seção"), so a credit large enough (RestituicaoUnit,
	// or StRetidoEntrada once it clamps PisCofins to 0 — see
	// searchSegmentICMSCellBracket) can push margem_pct ABOVE this asymptote
	// at a finite preço, approaching it from ABOVE instead of below. For
	// ICMSCell this check is skipped and the real search below (which
	// reasons per sub-interval about the sign of its own net fixed
	// coefficient) decides reachability instead.
	if in.ICMSCell == nil && cmpPct(in.TargetMargemPct, ceiling) >= 0 {
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
	//
	// Coordinator review (post-A5, Achado 1 — REFUTED): an earlier version of
	// this gate probed reachability by substituting a synthetic frete=0
	// (icmsCellPossiblyReachableAboveCeiling, now deleted). That was wrong on
	// two counts, both measured: (1) it fabricated a concrete value for an
	// unknown fact — ADR-17 backwards — and (2) when the probe reported
	// "unreachable", it returned CeilingPct, the very asymptote this task's
	// own C1 fix proved is NOT a true ceiling once a credit is present; the
	// code asserted a limit it had itself just disproven. Measured
	// counter-example: target at the cell's real ceiling with a credit
	// present reported CeilingPct even though a real (larger) frete reaches
	// it.
	//
	// Fix: decide from the credit terms directly, which are known WITHOUT
	// frete (icmsCellCreditsAndFixed reads R/S from the cell, not from
	// FreteProduto). fnet is monotonically DECREASING in frete — frete only
	// ADDS to fixed0 (Money ≥ 0) — so evaluating fnet at frete=0 gives its
	// SUPREMUM over the whole domain of the unknown. That supremum, not a
	// guessed value, is what decides whether the asymptote can be trusted:
	// icmsCellCeilingHoldsForEveryFrete.
	//
	// DUAL GATE (Fatia A close, Opus + GPT-5.6 Sol independently, blocking):
	// the R≤0 ∧ S≤0 branch above proves only that the asymptote is an HONEST
	// ceiling — it does NOT prove the target sits above it. The legacy path
	// gets that comparison for free in the ceilingPct short-circuit above;
	// ICMSCell skips that short-circuit (C1),
	// so nothing made it, and returning CeilingPct with !Reached is mapped by
	// transport (calc_handler.go) to UNREACHABLE_TARGET: "no such preço
	// exists". Measured on the credit-free BA fixture, ceiling 60.15:
	//
	//	alvo   frete desconhecido        frete real 15.00
	//	40.00  UNREACHABLE, teto 60.15   alcançado a 124.03
	//	50.00  UNREACHABLE, teto 60.15   alcançado a 246.18
	//	55.00  UNREACHABLE, teto 60.15   alcançado a 485.17
	//
	// Three targets strictly BELOW the ceiling were told the price does not
	// exist while it does — a fabricated verdict of non-existence built on an
	// unread fact, the very class the deleted frete=0 probe was removed for.
	// Below the ceiling the high segment decides, and the high segment needs
	// the frete: the honest answer there is FreteDesconhecido.
	// Pinned by TestSolveTargetPriceICMSCellNoCreditFreteUnknownBelowCeiling-
	// IsNotUnreachable (both halves of the gate, side by side).
	//
	// RE-GATE (GPT-5.6 Sol, blocking): the repaired predicate was still
	// R ≤ 0 ∧ S ≤ 0, which is SUFFICIENT for "the asymptote holds for every
	// frete" but not NECESSARY — a credit dominated by a KNOWN custo or
	// tarifa_full falls in the gap and got FreteDesconhecido even though the
	// verdict does not depend on the frete at all. Measured, BA fixture with
	// RestituicaoUnit=1.00 and custo=100.00 (ceiling 60.15), alvo 70.00:
	// frete 0.00 / 0.01 / 15.00 / 1000.00 all return UNREACHABLE teto 60.15,
	// while frete unknown returned FreteDesconhecido. Invariant across the
	// unknown's entire domain ⇒ the frete was never needed. A guard wider
	// than the fact is the same defect as one narrower, pointing the other
	// way (HARNESS-PROFILE, ratified 2026-07-28), and no must-fail arm can
	// catch over-strictness — only a MUST-PASS, which the new test carries
	// (alvo 55.00, where the price genuinely moves with the frete:
	// 1921.77 at 0.00 vs 21334.61 at 1000.00, so FreteDesconhecido stays).
	// Pinned by TestSolveTargetPriceICMSCellDominatedCreditFreteUnknown-
	// StillCeilings.
	if in.FreteProduto == nil {
		if in.ICMSCell != nil {
			if in.icmsCellCeilingHoldsForEveryFrete() && cmpPct(in.TargetMargemPct, ceiling) >= 0 {
				return SolveResult{CeilingPct: ceiling}
			}
		}
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
//     sets icms_saida=P×aCusto, difal=0 (icms.go); interestadual sets
//     icms_saida=P×aInter, difal=P×aCusto−P×aInter, so the sum is P×aCusto
//     either way (icms.go) — a_inter cancels out, so UFDestino never
//     needs consulting here, only aCusto.
//   - pis_cofins's base has a flat subtrahend (StRetidoEntrada) that washes
//     out as preço→∞, leaving pisCofinsRate×(1−aBase) — the same aBase
//     (aCusto minus FCP, D-43) TaxesForItem derives from the cell.
//   - isST (codTribST) never consults AliquotaInterna/FcpEmbutido at all
//     (icms.go: icms_saida=difal=0, base=P−S) — folding aCusto=aBase=0
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
// icms.go. TaxesForItem (icms.go) is the single source of truth this
// function tracks: whenever aCusto, aBase/FCP folding, isST's zeroing, or the
// PIS/COFINS RATES OR FORMULA change there, re-derive this function's four
// bullets above against the new formula before merging — a drift here
// reproduces exactly the Finding-1 class of bug (asymptote silently diverges
// from the real D-41 formula).
//
// That trigger used to name only the constant `pisCofinsRate`, and Task A8
// walked straight through the hole: it changed the pis_cofins FORMULA (one
// combined round2 became two separate ones) without touching that constant,
// so nothing here said to re-check. The asymptote survived it — rounding
// washes out as preço→∞, leaving (pisRate+cofinsRate)×(1−aBase) — but only
// because pisRate+cofinsRate is EXACTLY pisCofinsRate. That identity is now
// asserted by TestPisCofinsRatesSumToTheSolverAsymptoteRate, because editing
// one of the two split rates would break this asymptote while leaving the
// named trigger untouched, and prose cannot catch that.
func (in SolveInput) icmsCellAsymptoticRatePct() *big.Rat {
	aCusto, aBase := in.icmsCellRates()
	oneMinusABase := new(big.Rat).Sub(big.NewRat(1, 1), aBase)
	rate := new(big.Rat).Mul(pisCofinsRate, oneMinusABase)
	rate.Add(rate, aCusto)
	rate.Mul(rate, cem)
	return rate
}

// icmsCellRates extracts the D-41 rate pair (aCusto, aBase) from the cell —
// shared by icmsCellAsymptoticRatePct (the asymptote, above) and
// searchSegmentICMSCellBracket (Task A5: the clamped/unclamped regime
// coefficients need the SAME rates, not a re-derivation, to stay coupled to
// TaxesForItem exactly as icmsCellAsymptoticRatePct's doc comment requires).
func (in SolveInput) icmsCellRates() (aCusto, aBase *big.Rat) {
	cell := in.ICMSCell
	isST := cell.CodTrib != nil && *cell.CodTrib == codTribST
	aCusto = new(big.Rat)
	aBase = new(big.Rat)
	if !isST {
		aCusto.Quo(mustRat(*cell.AliquotaInterna), cem)
		fcp := new(big.Rat)
		if cell.FcpEmbutido != nil {
			fcp.Quo(mustRat(*cell.FcpEmbutido), cem)
		}
		aBase.Sub(aCusto, fcp)
	}
	return aCusto, aBase
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
		// Task A5: ICMSCell's credit term (RestituicaoST) can flip the sign
		// of the net fixed coefficient within the segment (C1/C2), so the
		// legacy single-bracket-per-segment approach below is unsound for
		// it — it assumes one monotone direction across the WHOLE segment.
		// Dispatch to a dedicated per-regime bracket instead (C3: it also
		// derives its own window from the ICMSCell term count, not this
		// branch's legacy 150).
		if in.ICMSCell != nil {
			return in.searchSegmentICMSCellBracket(loCents, hiCents)
		}
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
		winCents := ceilRatToCents(new(big.Rat).Quo(big.NewRat(windowNumeratorCents(roundedTermsLegacy), 1), gStar)) + windowMarginCents
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

// firstCentSatisfying returns the smallest cents in [lo,hi] for which
// predicate is true, given predicate is false-then-true (monotone) over the
// range, or hi+1 if it is never true. Generic bisection core shared by the
// legacy ascending bracket (firstCentExactAtLeast) and the ICMSCell regime
// bracket (searchMonotoneRegime, Task A5 C2/C3) — the latter's predicate
// direction flips per sub-interval depending on the sign of that regime's
// net fixed coefficient (fnet), so the DIRECTION is chosen by the caller,
// never baked into the bisection itself.
func firstCentSatisfying(lo, hi int64, predicate func(int64) bool) int64 {
	if !predicate(hi) {
		return hi + 1
	}
	for lo < hi {
		mid := lo + (hi-lo)/2
		if !predicate(mid) {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// firstCentExactAtLeast returns the smallest cents in [lo,hi] whose EXACT
// (unrounded, strictly-increasing) margem_pct is ≥ bound, or hi+1 if none.
func (in SolveInput) firstCentExactAtLeast(lo, hi, limiar int64, bound *big.Rat) int64 {
	return firstCentSatisfying(lo, hi, func(c int64) bool {
		return in.exactMargemPctRat(c, limiar).Cmp(bound) >= 0
	})
}

// exactMargemPctRat is the UNROUNDED margem_pct at preço=cents for the LEGACY
// segment only: 100·k − 100·F/preço, where k = 1 − comissão/100 −
// AliquotaPct/100 − applied difal/100, and F is the fixed-cost sum (taxa_fixa
// below limiar / produto frete at-or-above, plus tarifa_full when full and
// custo). Strictly increasing in preço because F ≥ 0 for the legacy path (no
// credit term). Used ONLY to bracket the legacy high-segment search window.
//
// Task A5: this used to also branch on ICMSCell, sharing one bracket with the
// legacy path. That coupling was the root of C1/C2 — RestituicaoST (and,
// past the PIS/COFINS clamp, StRetidoEntrada's rate) can flip the sign of the
// net fixed coefficient within the ICMSCell domain, so a single
// ascending-only bracket over the WHOLE segment is unsound there. The
// ICMSCell path now gets its own dedicated per-regime bracket
// (searchSegmentICMSCellBracket / searchMonotoneRegime) that reasons about
// that sign per sub-interval instead. searchSegment dispatches ICMSCell
// inputs there before ever reaching this function — it is never called with
// in.ICMSCell != nil.
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

// exactCeilingRat is the UNROUNDED legacy-path margem_pct asymptote (100 −
// comissão − AliquotaPct − applied difal) as preço→∞. The window derivation
// needs the unrounded value: the 2dp ceiling gate can admit a target only
// ~0.005 below the true asymptote when a >2dp comissão is supplied, so a
// rounded gap would understate the window.
//
// Task A5: see exactMargemPctRat's comment — ICMSCell no longer flows through
// this function. Its asymptote/regime coefficients live in
// icmsCellAsymptoticRatePct and searchSegmentICMSCellBracket's own
// kClamped/kUnclamped instead.
func (in SolveInput) exactCeilingRat() *big.Rat {
	ceil := new(big.Rat).Set(cem)
	ceil.Sub(ceil, mustRat(in.ComissaoPct))
	ceil.Sub(ceil, mustRat(in.AliquotaPct))
	if in.difalApplied() {
		ceil.Sub(ceil, mustRat(in.EfetivoPct))
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

// roundedTermsClamped/Unclamped/Legacy count the components that are
// INDEPENDENTLY rounded to 2dp from a computed (not merely copied) value
// before Decompose sums them — Task A5 (C3): the window numerator is DERIVED
// from this count (windowNumeratorCents), never a bare literal, because the
// bound comes from the same "≤50 cents of rounding error per term, divided
// by the slope gStar" argument searchSegment's legacy comment spells out
// (150 = 3×50). A term only counts if summing its ROUNDED value can differ
// from summing its EXACT value by up to half a cent — copying an
// already-2dp money field through round2 (a no-op) does not qualify.
//
// Coordinator review (post-A5, Achado 2): the original derivation here
// enumerated the wrong three/four terms (folded ICMSSaida+Difal into one
// "ICMS (aCusto)" slot, and treated PIS/COFINS as a placeholder in the
// clamped regime instead of recognizing it is exactly, not approximately,
// absent there) and asserted a false claim about how Decompose rounds the
// D-41 tax load. Re-derived directly against decompose.go, which
// sums four ALREADY-ROUNDED-INSIDE-TaxesForItem components separately —
// mustRat(*tax.ICMSSaida), mustRat(*tax.Difal), mustRat(*tax.PisCofins) (each
// added), mustRat(*tax.RestituicaoST) (subtracted) — never "the exact tax
// load summed once then rounded":
//   - ICMSSaida = round2(P×aInter or P×aCusto) (icms.go, zero-exact in
//     the isST branch, icms.go) — genuinely rounded from a
//     price-proportional computation.
//   - Difal = round2(icmsTotal−icmsOper) (icms.go, zero-exact in the
//     isST branch) — rounded SEPARATELY from ICMSSaida even though the two
//     telescope to P×aCusto exactly in the unrounded/exact model this
//     file's kPct uses (icmsCellAsymptoticRatePct's doc comment) — each
//     independent round2 call is its own up-to-half-cent error source, so
//     ICMSSaida and Difal are TWO terms, not one "aCusto" term.
//   - PisCofins = round2(pisRate×base) + round2(cofinsRate×base) (icms.go,
//     pisCofinsFrom) — Task A8: PIS and COFINS are TWO independently-rounded
//     terms, not one. The base is clamped to MAX(0, base) once, before either
//     rate (icms.go's "base.Sign()<0 ⇒ base=0"); in the clamped sub-region
//     that forces base to exactly 0, so BOTH round2(0×pisRate) and
//     round2(0×cofinsRate) are "0.00" with ZERO uncertainty — neither is a
//     placeholder slot standing in for future rounding, both genuinely
//     contribute no slack in that sub-region and are correctly excluded from
//     roundedTermsClamped. In the unclamped sub-region the base is positive
//     and BOTH round2 calls are real, independent up-to-half-cent error
//     sources — two terms, not one.
//   - Comissão = round2(comissãoPct×preço) (decompose.go) — rounded in
//     BOTH regimes, present in every count below.
//   - RestituicaoST = round2(cell.RestituicaoUnit) (icms.go) is
//     EXCLUDED from every count: RestituicaoUnit is already a decimal-money
//     string from its source (products_mirror.restituicao_unit, ICMSCell.RestituicaoUnit), not derived from a price-proportional formula, so round2 here
//     is a no-op with no real information loss — a flat credit, never a
//     rounded proportional term (decompose.go, the credit is
//     SUBTRACTED from sum, the one non-cost line in the section).
//
// So per regime:
//   - roundedTermsLegacy = 3: comissão + imposto (AliquotaPct) + difal (the
//     OLD, non-ICMSCell path — decompose.go/165/196-201 — unaffected by
//     Task A8, unchanged).
//   - roundedTermsClamped = 3: comissão + ICMSSaida + Difal (PIS and COFINS
//     are exactly, not approximately, zero here — see above; the split into
//     two terms does not change this count because both are zero-uncertainty
//     in this sub-region).
//   - roundedTermsUnclamped = 5: comissão + ICMSSaida + Difal + PIS + COFINS
//     — Task A8: PIS and COFINS are now genuinely rounded SEPARATELY from a
//     positive base (pisCofinsFrom's two independent round2 calls), so what
//     was one PisCofins term is now two. Re-derived directly against
//     icms.go's pisCofinsFrom as it now reads (not bumped pre-emptively —
//     this is the re-count a prior revision of this comment said would be
//     needed once PIS/COFINS rounded separately, and that point is now).
const roundedTermsClamped = 3
const roundedTermsUnclamped = 5
const roundedTermsLegacy = 3

// windowNumeratorCents is the DERIVED window numerator (Task A5 C3): 50 cents
// of worst-case rounding slack per proportional/rounded term, summed across
// termCount terms. Replaces the legacy bracket's former bare `150` literal
// (now windowNumeratorCents(roundedTermsLegacy), same value, same math, just
// no longer hardcoded) and supplies the ICMSCell regimes' windows —
// windowNumeratorCents(roundedTermsClamped) and
// windowNumeratorCents(roundedTermsUnclamped). Written as symbols on
// purpose: this line read "150/200" until the dual gate caught it, A8 having
// moved roundedTermsUnclamped 4→5 (so the second number is now 250) without
// the prose following.
func windowNumeratorCents(termCount int) int64 {
	return 50 * int64(termCount)
}

// icmsCellCreditsAndFixed extracts the D-41 credit terms (R = RestituicaoST,
// zeroed for UFDestino=MG per icms.go's MG-interno branch; S =
// StRetidoEntrada) and the flat non-tax fixed-cost sum (frete + tarifa_full +
// custo — NOT taxa_fixa, which is segment-selected by the caller) that
// searchSegmentICMSCellBracket's two regimes both build their net fixed
// coefficient (fnet) from.
func (in SolveInput) icmsCellCreditsAndFixed() (r, s, fixed0 *big.Rat) {
	cell := in.ICMSCell
	if cell.UFDestino == ufOrigemMG {
		r = big.NewRat(0, 1)
	} else {
		r = zeroIfNil(cell.RestituicaoUnit)
	}
	s = zeroIfNil(cell.StRetidoEntrada)
	fixed0 = new(big.Rat)
	if in.FreteProduto != nil {
		fixed0.Add(fixed0, mustRat(in.FreteProduto.Amount))
	}
	if in.Modalidade == ModalidadeFull && in.TarifaFull != nil {
		fixed0.Add(fixed0, mustRat(in.TarifaFull.Amount))
	}
	if in.Custo != nil {
		fixed0.Add(fixed0, mustRat(in.Custo.Amount))
	}
	return r, s, fixed0
}

// icmsCellCeilingHoldsForEveryFrete answers, WITHOUT reading FreteProduto,
// whether the ICMSCell asymptote is a true upper bound for every frete the
// unknown could turn out to be.
//
// fnet = R − fixed0 in the clamped sub-region (PIS/COFINS pinned at 0) and
// R + pisCofinsRate·S − fixed0 in the unclamped one — the two values
// searchSegmentICMSCellBracket derives. fixed0 = frete + tarifa_full + custo,
// every term Money ≥ 0, so a real frete only moves fnet DOWN. Evaluating both
// at frete=0 (the caller reaches here only with FreteProduto == nil, so
// icmsCellCreditsAndFixed already omits it) therefore yields the SUPREMUM of
// fnet over the unknown's whole domain. Both suprema ≤ 0 ⇒ fnet ≤ 0 for every
// admissible frete ⇒ exact(P) approaches the asymptote from below in both
// sub-regions ⇒ the asymptote is a genuine ceiling and CeilingPct is honest.
//
// This is a BOUND on an unknown, not a substituted value: nothing downstream
// consumes a frete of 0, and the branch is only ever used to decide that the
// answer does not depend on the frete. That distinction is what separates it
// from the deleted icmsCellPossiblyReachableAboveCeiling probe, which computed
// a whole solve against a fabricated frete=0.
func (in SolveInput) icmsCellCeilingHoldsForEveryFrete() bool {
	r, s, fixed0 := in.icmsCellCreditsAndFixed()
	// Clamped sub-region: PIS/COFINS = 0, so S contributes nothing.
	if new(big.Rat).Sub(r, fixed0).Sign() > 0 {
		return false
	}
	// Unclamped sub-region.
	unclamped := new(big.Rat).Mul(pisCofinsRate, s)
	unclamped.Add(unclamped, r)
	unclamped.Sub(unclamped, fixed0)
	return unclamped.Sign() <= 0
}

// icmsCellClampBoundaryCents is the preço (in cents, ceil-rounded) at which
// the PIS/COFINS base P×(1−aBase) − S crosses 0 — Task A5 (C2)'s second
// discontinuity, independent of and additional to the existing taxa_fixa/
// frete limiar step. Below the boundary PIS/COFINS is clamped to 0 (MAX(0,·)
// in icms.go); at/above it, PIS/COFINS = pisCofinsRate×(P×(1−aBase) − S).
// s≤0 ⇒ no clamp region ever exists (boundary=0, always in the unclamped
// regime). aBase≥1 (pathological: FCP-adjusted rate ≥100%) ⇒ the base can
// never turn positive, so every preço is clamped (alwaysClamped=true).
func icmsCellClampBoundaryCents(aBase, s *big.Rat) (boundaryCents int64, alwaysClamped bool) {
	if s.Sign() <= 0 {
		return 0, false
	}
	denom := new(big.Rat).Sub(big.NewRat(1, 1), aBase)
	if denom.Sign() <= 0 {
		return 0, true
	}
	pclamp := new(big.Rat).Quo(s, denom)
	pclamp.Mul(pclamp, cem)
	return ceilRatToCents(pclamp), false
}

// searchMonotoneRegime brackets and scans ONE monotone sub-interval of the
// ICMSCell domain: exact(P) = kPct + 100·fnet/P. Task A5 (C1/C2): unlike the
// legacy bracket, fnet is not assumed ≤0 — its SIGN determines the search
// direction. fnet≤0 ⇒ exact(P) rises toward kPct from BELOW as P grows (same
// shape as the legacy asymptote, kPct is a true ceiling in this
// sub-interval). fnet>0 ⇒ exact(P) FALLS toward kPct from ABOVE as P grows
// (the C1 credit-driven mechanism: kPct is a floor here, and targets ABOVE
// kPct are reachable at small-enough P within this sub-interval — exactly
// the case the old unconditional ceiling short-circuit wrongly rejected).
// Both directions reduce to the same false-then-true bisection
// (firstCentSatisfying) by choosing the comparison operator that keeps the
// predicate monotone as cents increase.
func (in SolveInput) searchMonotoneRegime(lo, hi int64, kPct, fnet *big.Rat, termCount int) (string, bool) {
	if lo > hi {
		return "", false
	}
	half := big.NewRat(5, 1000) // 0.005 = round2 half-band, same as legacy
	exact := func(c int64) *big.Rat {
		preco := big.NewRat(c, 100)
		term := new(big.Rat).Quo(fnet, preco)
		term.Mul(term, cem)
		return new(big.Rat).Add(kPct, term)
	}
	var bound, gStar *big.Rat
	var predicate func(int64) bool
	if fnet.Sign() <= 0 {
		// rising toward kPct from below: first cent whose exact value is
		// ≥ (target−half) is the crossing, same shape as the legacy bracket.
		bound = new(big.Rat).Sub(mustRat(in.TargetMargemPct), half)
		gStar = new(big.Rat).Sub(kPct, bound)
		predicate = func(c int64) bool { return exact(c).Cmp(bound) >= 0 }
	} else {
		// falling toward kPct from above: first cent whose exact value has
		// dropped to ≤ (target+half) is the crossing — the mirror image.
		bound = new(big.Rat).Add(mustRat(in.TargetMargemPct), half)
		gStar = new(big.Rat).Sub(bound, kPct)
		predicate = func(c int64) bool { return exact(c).Cmp(bound) <= 0 }
	}
	if gStar.Sign() <= 0 {
		// target at/past kPct in this sub-interval's direction — no crossing
		// exists HERE (a different sub-interval, or CeilingPct, may still
		// reach it; that is the caller's decision, not this function's).
		return "", false
	}
	cStar := firstCentSatisfying(lo, hi, predicate)
	if cStar > hi {
		return "", false
	}
	win := ceilRatToCents(new(big.Rat).Quo(big.NewRat(windowNumeratorCents(termCount), 1), gStar)) + windowMarginCents
	return in.scanWindow(lo, hi, cStar, win)
}

// scanWindow runs the real Decompose exhaustively over [center-win,
// center+win] ∩ [lo,hi] — the same exact-match contract as searchSegment's
// legacy scan, factored out so both the legacy and ICMSCell paths return a
// price only when the ACTUAL rounded Decompose margem_pct matches, never the
// analytic approximation.
func (in SolveInput) scanWindow(lo, hi, center, win int64) (string, bool) {
	scanLo := center - win
	if scanLo < lo {
		scanLo = lo
	}
	scanHi := center + win
	if scanHi > hi {
		scanHi = hi
	}
	for c := scanLo; c <= scanHi; c++ {
		if m := in.margemAt(c); m != "" && cmpPct(m, in.TargetMargemPct) == 0 {
			return centsToStr(c), true
		}
	}
	return "", false
}

// searchSegmentICMSCellBracket is the ICMSCell replacement for the legacy
// single-bracket searchSegment body (Task A5, C1+C2+C3). The D-41 margem_pct
// is piecewise over TWO independent discontinuities: the existing taxa_fixa/
// frete limiar (selected by the caller via loCents/hiCents, unchanged) and
// the NEW PIS/COFINS clamp boundary (icmsCellClampBoundaryCents) — so within
// one taxa_fixa/frete segment there can still be a clamped sub-region
// (PIS/COFINS≡0) below the boundary and an unclamped sub-region at/above it.
// Each sub-region gets its own kPct/fnet pair and is searched independently
// via searchMonotoneRegime; the clamped sub-region is tried first (cheaper
// prices).
//
// It takes no limiar: the taxa_fixa/frete step is already resolved by the
// caller into loCents/hiCents, and each sub-region's fnet is built from
// icmsCellCreditsAndFixed (which folds frete in), so nothing here needs to
// re-test which side of the step a price falls on.
func (in SolveInput) searchSegmentICMSCellBracket(loCents, hiCents int64) (string, bool) {
	r, s, fixed0 := in.icmsCellCreditsAndFixed()
	aCusto, aBase := in.icmsCellRates()
	boundary, alwaysClamped := icmsCellClampBoundaryCents(aBase, s)

	// Clamped sub-region: PIS/COFINS ≡ 0 (base ≤ 0), so the fiscal load is
	// just aCusto (ICMS) — HIGHER margem than the unclamped asymptote since
	// no PIS/COFINS is charged. fnet = R − fixed0 (only the credit and the
	// non-tax fixed costs; S does not appear because PIS/COFINS is zero
	// here, C2's "before the clamp opens, do not use the asymptotic rate").
	kClamped := new(big.Rat).Sub(cem, mustRat(in.ComissaoPct))
	kClamped.Sub(kClamped, new(big.Rat).Mul(aCusto, cem))
	fnetClamped := new(big.Rat).Sub(r, fixed0)

	clampedHi := hiCents
	if !alwaysClamped && boundary-1 < clampedHi {
		clampedHi = boundary - 1
	}
	if clampedHi >= loCents {
		if p, ok := in.searchMonotoneRegime(loCents, clampedHi, kClamped, fnetClamped, roundedTermsClamped); ok {
			return p, true
		}
	}
	if alwaysClamped {
		return "", false
	}

	// Unclamped sub-region: PIS/COFINS = pisCofinsRate×(P×(1−aBase) − S),
	// contributing pisCofinsRate·S/preço to fnet (C2) once P is large enough
	// for the base to be positive. kPct here is icmsCellAsymptoticRatePct's
	// exact ceiling (the true D-41 asymptote for this cell).
	unclampedLo := loCents
	if boundary > unclampedLo {
		unclampedLo = boundary
	}
	if unclampedLo > hiCents {
		return "", false
	}
	kUnclamped := new(big.Rat).Sub(cem, mustRat(in.ComissaoPct))
	kUnclamped.Sub(kUnclamped, in.icmsCellAsymptoticRatePct())
	fnetUnclamped := new(big.Rat).Mul(pisCofinsRate, s)
	fnetUnclamped.Add(fnetUnclamped, r)
	fnetUnclamped.Sub(fnetUnclamped, fixed0)

	return in.searchMonotoneRegime(unclampedLo, hiCents, kUnclamped, fnetUnclamped, roundedTermsUnclamped)
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
