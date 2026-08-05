package domain

import (
	"math/big"
	"reflect"
	"testing"
	"time"
)

// baseSolve returns a fully-resolved SolveInput (classico, SIMPLES 4%, difal
// SP 6.00, frete 15.00, custo 40.00) with no target set — tests fill Target.
func baseSolve() SolveInput {
	return SolveInput{
		ComissaoPct: "12", AliquotaPct: "4", Modalidade: ModalidadeClassico,
		FreteProduto: money("15.00"), Custo: money("40.00"),
		DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
	}
}

// resim is the oracle: re-Decompose the solved preço under the SAME inputs and
// return the margem_pct the simulator would show. ICMSCell is forwarded
// (Task A3) so the oracle reflects the SAME D-41 path margemDecompose now
// feeds — a solved preço is only a true fixed point if re-decomposing it
// through the real tax cell reproduces the target exactly.
func resim(in SolveInput, preco string) string {
	d := Decompose(DecomposeInput{
		Preco: preco, ComissaoPct: in.ComissaoPct, AliquotaPct: in.AliquotaPct,
		Modalidade: in.Modalidade, TarifaFull: in.TarifaFull,
		FreteProduto: in.FreteProduto, Custo: in.Custo,
		DifalEnabled: in.DifalEnabled, DestinoUF: in.DestinoUF, EfetivoPct: in.EfetivoPct,
		ICMSCell: in.ICMSCell,
	})
	if d.MargemPct == nil {
		return "<unknown>"
	}
	return *d.MargemPct
}

// C03(a) — margem-alvo 15% ⇒ solved preço whose re-Decompose margem_pct is
// EXACTLY "15.00". Solution lands below 79 (cheapest price hitting target).
func TestSolveTargetPriceExactMargin(t *testing.T) {
	in := baseSolve()
	in.TargetMargemPct = "15.00"

	res := SolveTargetPrice(in)
	if !res.Reached || res.Preco == nil {
		t.Fatalf("target 15.00 must be reachable; got %+v", res)
	}
	if got := resim(in, *res.Preco); got != "15.00" {
		t.Fatalf("re-sim margem_pct = %q at preço %q, want 15.00 EXACT", got, *res.Preco)
	}
}

// C03(b) — a target only attainable ABOVE the preço=79 discontinuity. A naive
// bisection over [1, max] assuming continuity would break on the downward jump
// at 79; the solver must reject the low segment (its ceiling < target) and
// converge in the high segment. Solved preço ≥ 79, re-sim EXACT.
func TestSolveTargetPriceCrossesDiscontinuity(t *testing.T) {
	in := baseSolve()
	in.TargetMargemPct = "30.00" // above the below-79 segment's max (~19.13%)

	res := SolveTargetPrice(in)
	if !res.Reached || res.Preco == nil {
		t.Fatalf("target 30.00 must be reachable in the high segment; got %+v", res)
	}
	if got := resim(in, *res.Preco); got != "30.00" {
		t.Fatalf("re-sim margem_pct = %q at preço %q, want 30.00 EXACT", got, *res.Preco)
	}
	// must be on/above the discontinuity (proves segment selection).
	if p := ratOf(t, *res.Preco); p.Cmp(taxaFixaLimiar) < 0 {
		t.Fatalf("solved preço %q is below 79 — high-segment target mis-solved", *res.Preco)
	}
}

// C03(c) — margem-alvo 60% with comissão 16% (+ PRESUMIDO 9,25% + difal MA
// 16,00%) ⇒ UNREACHABLE_TARGET citing the attainable ceiling. Ceiling =
// 100 − 16 − 9.25 − 16 = 58.75, so 60% can never be reached.
func TestSolveTargetPriceUnreachableReturnsCeiling(t *testing.T) {
	in := SolveInput{
		TargetMargemPct: "60.00",
		ComissaoPct:     "16", AliquotaPct: "9.25", Modalidade: ModalidadeClassico,
		FreteProduto: money("15.00"), Custo: money("40.00"),
		DifalEnabled: true, DestinoUF: "MA", EfetivoPct: "16.00",
	}
	res := SolveTargetPrice(in)
	if res.Reached || res.Preco != nil {
		t.Fatalf("target 60.00 must be UNREACHABLE; got %+v", res)
	}
	if res.CeilingPct != "58.75" {
		t.Fatalf("attainable ceiling = %q, want 58.75", res.CeilingPct)
	}
}

// UNKNOWN inputs (custo nil) are a BLOCKING condition, distinct from
// UNREACHABLE: the solver reports componentes_desconhecidos, not a ceiling,
// and does not fabricate a preço.
func TestSolveTargetPriceUnknownInputsBlock(t *testing.T) {
	in := baseSolve()
	in.Custo = nil
	in.TargetMargemPct = "15.00"

	res := SolveTargetPrice(in)
	if res.Reached || res.Preco != nil || res.CeilingPct != "" {
		t.Fatalf("custo nil must block (no preço, no ceiling); got %+v", res)
	}
	if len(res.Desconhecidos) != 1 || res.Desconhecidos[0] != "custo_erp" {
		t.Fatalf("Desconhecidos = %v, want [custo_erp]", res.Desconhecidos)
	}
}

// FreteProduto nil but the target is attainable in the low segment (preço <
// limiar), where frete is 0 and never consulted — must still solve.
func TestSolveFreteNilLowSegmentSolves(t *testing.T) {
	in := baseSolve()
	in.FreteProduto = nil
	in.TargetMargemPct = "15.00"

	res := SolveTargetPrice(in)
	if !res.Reached || res.Preco == nil || res.FreteDesconhecido || len(res.Desconhecidos) != 0 {
		t.Fatalf("target 15.00 with frete nil must solve in the low segment; got %+v", res)
	}
	if p := ratOf(t, *res.Preco); p.Cmp(taxaFixaLimiar) >= 0 {
		t.Fatalf("solved preço %q must be below 79 (low segment)", *res.Preco)
	}
	if got := resim(in, *res.Preco); got != "15.00" {
		t.Fatalf("re-sim margem_pct = %q at preço %q, want 15.00 EXACT", got, *res.Preco)
	}
}

// FreteProduto nil and the target is only attainable in the high segment
// (preço ≥ limiar), where produto frete is required but unknown — must block
// with FreteDesconhecido, not fabricate a preço or a ceiling.
func TestSolveFreteNilHighSegmentBlocks(t *testing.T) {
	in := baseSolve()
	in.FreteProduto = nil
	in.TargetMargemPct = "30.00"

	res := SolveTargetPrice(in)
	if !res.FreteDesconhecido || res.Reached || res.Preco != nil || len(res.Desconhecidos) != 0 || res.CeilingPct != "" {
		t.Fatalf("target 30.00 with frete nil must block on FreteDesconhecido; got %+v", res)
	}
}

// custo nil is a structural (segment-independent) block that must win over
// the frete-nil segment-conditional block.
func TestSolveCustoNilBlocksBeforeFrete(t *testing.T) {
	in := baseSolve()
	in.Custo = nil
	in.FreteProduto = nil
	in.TargetMargemPct = "30.00"

	res := SolveTargetPrice(in)
	if len(res.Desconhecidos) != 1 || res.Desconhecidos[0] != "custo_erp" || res.FreteDesconhecido || res.Reached {
		t.Fatalf("custo nil must block structurally (not FreteDesconhecido); got %+v", res)
	}
}

// C03(d) — FINDING-P6-SOLVER regression. Decompose's 2dp component rounding
// makes margem_pct NON-monotone in preço (sub-cent downward wiggles), so the
// old rounded-value bisection skipped valid prices and mis-reported reachable
// targets. This brute-forces the cheapest exact-match preço for EVERY distinct
// low-segment margem_pct and asserts SolveTargetPrice returns exactly it —
// across ordinary comissão/aliquota grids that demonstrably contain dips.
func TestSolveMatchesBruteForceLowSegment(t *testing.T) {
	grids := []SolveInput{
		{ComissaoPct: "12", AliquotaPct: "4", Modalidade: ModalidadeClassico, FreteProduto: money("15.00"), Custo: money("10.00")},
		{ComissaoPct: "13", AliquotaPct: "4", Modalidade: ModalidadeClassico, FreteProduto: money("15.00"), Custo: money("40.00")},
		{ComissaoPct: "16", AliquotaPct: "12", Modalidade: ModalidadeClassico, FreteProduto: money("15.00"), Custo: money("10.00")},
	}
	for _, base := range grids {
		// cheapest[target] = smallest cents whose Decompose margem_pct == target.
		cheapest := map[string]int64{}
		for c := int64(1); c <= 7898; c++ {
			m := resim(base, centsToStr(c))
			if m == "<unknown>" {
				continue
			}
			if _, seen := cheapest[m]; !seen {
				cheapest[m] = c
			}
		}
		// Deterministic sample: stride over cents in order (map iteration order is
		// nondeterministic). Each sampled cent's target is looked up against the
		// brute cheapest map, so the assertion is still cheapest-exact-match.
		for c := int64(1); c <= 7898; c += 37 {
			target := resim(base, centsToStr(c))
			if target == "<unknown>" {
				continue
			}
			wantCents := cheapest[target]
			in := base
			in.TargetMargemPct = target
			res := SolveTargetPrice(in)
			if !res.Reached || res.Preco == nil {
				t.Fatalf("com=%s ali=%s target=%s must reach (brute @ %s); got %+v",
					base.ComissaoPct, base.AliquotaPct, target, centsToStr(wantCents), res)
			}
			if *res.Preco != centsToStr(wantCents) {
				t.Fatalf("com=%s ali=%s target=%s: solver=%s cheapest-brute=%s",
					base.ComissaoPct, base.AliquotaPct, target, *res.Preco, centsToStr(wantCents))
			}
		}
	}
}

// C03(e) — high-segment (preço ≥ limiar) exactness + cheapest: for a witness
// price the solver must return SOME preço whose re-Decompose margem_pct equals
// the target EXACTLY and is no costlier than the witness (proving the exact
// bracket does not skip the reachable band the way the old bisection did).
func TestSolveHighSegmentExactAndCheapest(t *testing.T) {
	base := SolveInput{
		ComissaoPct: "13", AliquotaPct: "4", Modalidade: ModalidadeClassico,
		FreteProduto: money("15.00"), Custo: money("40.00"),
	}
	for _, witness := range []int64{9000, 15000, 30000, 90000} {
		tgt := resim(base, centsToStr(witness))
		in := base
		in.TargetMargemPct = tgt
		res := SolveTargetPrice(in)
		if !res.Reached || res.Preco == nil {
			t.Fatalf("witness %s target %s must reach; got %+v", centsToStr(witness), tgt, res)
		}
		if got := resim(in, *res.Preco); got != tgt {
			t.Fatalf("solver preço %s re-sim=%s want %s EXACT", *res.Preco, got, tgt)
		}
		if ratOf(t, *res.Preco).Cmp(ratOf(t, centsToStr(witness))) > 0 {
			t.Fatalf("solver preço %s costlier than witness %s (not cheapest)", *res.Preco, centsToStr(witness))
		}
	}
}

// centsOf parses a 2dp preço string to integer cents (fails if not whole cents).
func centsOf(t *testing.T, price string) int64 {
	t.Helper()
	r := ratOf(t, price)
	r.Mul(r, big.NewRat(100, 1))
	if !r.IsInt() {
		t.Fatalf("preço %s is not whole cents", price)
	}
	return r.Num().Int64()
}

// C03(f) — FINDING-P6-SOLVER-2 regime: a >2dp comissão makes the EXACT ceiling
// carry a 3rd decimal, so the exact gap to the ceiling can be as small as 0.005
// for a reachable (2dp) target. The high-segment scan window must be DERIVED
// from that gap, not a fixed cap that could truncate the reachable band. Here
// exactCeiling = 100 − 12.005 = 87.995 (2dp ceiling gate = 88.00) and F = 0.50
// keeps the crossing price inside the R$1M cap as the gap shrinks. Assert exact
// match AND that no cheaper exact match exists just below the solved price (the
// window did not truncate the band). Targets are 2dp — a >2dp target is
// unreachable by construction since Decompose emits 2dp.
func TestSolveHighSegmentNearCeilingCheapest(t *testing.T) {
	base := SolveInput{
		ComissaoPct: "12.005", AliquotaPct: "0", Modalidade: ModalidadeClassico,
		FreteProduto: money("0.00"), Custo: money("0.50"),
	}
	for _, tgt := range []string{"87.90", "87.98", "87.99"} {
		in := base
		in.TargetMargemPct = tgt
		res := SolveTargetPrice(in)
		if !res.Reached || res.Preco == nil {
			t.Fatalf("target %s must reach (sub-0.01 ceiling-gap regime); got %+v", tgt, res)
		}
		if got := resim(in, *res.Preco); cmpPct(got, tgt) != 0 {
			t.Fatalf("target %s: solver preço %s re-sim=%s want EXACT", tgt, *res.Preco, got)
		}
		// no cheaper exact match below the solved price → window not truncated.
		solved := centsOf(t, *res.Preco)
		lo := solved - 200_000
		if lo < defaultTaxaFixaLimiarCents {
			lo = defaultTaxaFixaLimiarCents
		}
		for c := lo; c < solved; c++ {
			if cmpPct(resim(in, centsToStr(c)), tgt) == 0 {
				t.Fatalf("target %s: cheaper exact match at %s < solver %s (window truncated)",
					tgt, centsToStr(c), *res.Preco)
			}
		}
	}
}

// --- Task A3: solver enters the D-41 (ICMSCell) path ---------------------

// icmsCellBA is a RESOLVED, unambiguous cell for a BA destino (a_interna
// 20,5%, matching icms_erp_golden_test.go's case1/TestDecomposeICMSRestitui-
// caoIncreasesMargem fixture) — BA is outside the 12% a_inter set, so
// a_inter=7% default; fcp_embutido nil ⇒ 0 (BA has no general FCP in this
// fixture).
func icmsCellBA() *ICMSCell {
	return &ICMSCell{
		UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("20.5"),
		// Task A4 (D2): par completo -- "0" explicito, nao nil.
		// AliquotaInterna presente com FcpEmbutido nil e um estado que a
		// fonte real nunca produz; numericamente identico (fcp=0 nos dois
		// casos, ja que BA nesta fixture nao tem FCP geral), so a REALISMO
		// do fixture muda -- nenhum valor esperado nos goldens abaixo muda.
		FcpEmbutido: strptr("0"),
	}
}

// baseSolveICMSCell is a SolveInput exercising the D-41 path. AliquotaPct
// stays populated ("4") only so the fixture keeps a realistic legacy value
// around; it is NOT read on this path. Task A6 made that explicit —
// decomposeWithLimiar parses AliquotaPct and sets Imposto ONLY when
// ICMSCell == nil (decompose.go, the `if in.ICMSCell == nil` branch), so
// with a cell present Imposto stays nil (named-absent) and the fiscal load
// comes from ICMSSaida/Difal/PisCofins/RestituicaoST via TaxesForItem.
//
// This comment used to claim the opposite ("Decompose ALWAYS computes the
// Imposto field, decompose.go") — true before A6, false after it, in the
// same commit range, about the exact field the range changed. Caught by the
// dual gate (Opus side). Symbol names, not line numbers, from here on.
//
// DifalEnabled/DestinoUF/EfetivoPct are left at their zero value for the
// same reason — the legacy difal branch is not consulted once ICMSCell is
// non-nil.
func baseSolveICMSCell() SolveInput {
	return SolveInput{
		ComissaoPct: "12", AliquotaPct: "4", Modalidade: ModalidadeClassico,
		Custo:    money("10.00"),
		ICMSCell: icmsCellBA(),
	}
}

// TestSolveTargetPriceICMSCellBA is the task's mandated contract golden:
// "alvo de margem 20% para BA". Before Task A3, margemDecompose never
// forwarded SolveInput.ICMSCell into DecomposeInput, so the solver searched
// for the preço whose margem — computed with the LEGACY AliquotaPct="4"
// fabricated SIMPLES rate (D-38) as the only tax — hit 20%. That preço, when
// actually re-decomposed through the real tax cell (BA's alíquota interna
// 20,5%, D-41), does NOT sit at margem 20% — the same "decomposição certa,
// alvo com 4% fabricado, na mesma tela" bug the brief describes. After the
// fix, margemDecompose forwards ICMSCell and ceilingPct derives its
// asymptote from the SAME cell (icmsCellAsymptoticRatePct), so the solved
// preço's REAL margem_pct (re-simulated through the D-41 path, not the
// legacy one) is exactly the target.
func TestSolveTargetPriceICMSCellBA(t *testing.T) {
	in := baseSolveICMSCell()
	in.TargetMargemPct = "20.00"

	res := SolveTargetPrice(in)
	if !res.Reached || res.Preco == nil {
		t.Fatalf("target 20.00 (BA, alíquota interna 20.5) must be reachable; got %+v", res)
	}
	if got := resim(in, *res.Preco); got != "20.00" {
		t.Fatalf("re-sim (real ICMSCell decompose) margem_pct = %q at preço %q, want 20.00 EXACT — "+
			"the solver must use the SAME D-41 formula Decompose uses, never the legacy 4%% AliquotaPct",
			got, *res.Preco)
	}
	// This golden is scoped to the exhaustive low-segment scanner (span
	// 7898 cents, always run to completion — searchSegment never brackets
	// it), which always tests candidates against the REAL Decompose
	// (margemAt→margemDecompose) regardless of ICMSCell. The high-segment
	// bracket math (exactCeilingRat/exactMargemPctRat) still assumes the
	// legacy AliquotaPct/EfetivoPct source and is untouched by Task A3 (out
	// of scope — Fatia C removes the legacy path entirely), so this
	// fixture is deliberately built to land below the taxa_fixa/frete
	// limiar (79) rather than exercise that bracket.
	if p := ratOf(t, *res.Preco); p.Cmp(taxaFixaLimiar) >= 0 {
		t.Fatalf("solved preço %q is at/above 79 — expected low segment for this fixture", *res.Preco)
	}
}

// TestSolveTargetPriceICMSCellBAHighSegment is the A3 review's Finding-1
// reproduction. exactCeilingRat/exactMargemPctRat (the high-segment bracket
// math) were NOT updated to branch on ICMSCell when ceilingPct itself was —
// they kept deriving the asymptote from the legacy fabricated AliquotaPct=4
// (ceiling 100−12−4=84) instead of the real cell ceiling (60.15, see
// TestCeilingPctICMSCellBALiteral). Target 59.00 sits below the REAL ceiling
// (60.15) but only in the high segment (FreteProduto must be set — the low
// segment tops out at margem ≈19% for this fixture's small custo). Before the
// fix, gStar/cStar bracket around the WRONG (84) asymptote, so the bisection
// converges near preço≈79 instead of the true crossing near R$1302, the
// window scan never reaches it, and the solver wrongly reports
// UNREACHABLE_TARGET citing CeilingPct="60.15" — the right ceiling with the
// wrong reachability verdict underneath it.
func TestSolveTargetPriceICMSCellBAHighSegment(t *testing.T) {
	in := baseSolveICMSCell()
	in.FreteProduto = money("15.00")
	in.TargetMargemPct = "59.00" // < real ceiling 60.15, < legacy fabricated ceiling 84

	res := SolveTargetPrice(in)
	if !res.Reached || res.Preco == nil {
		t.Fatalf("target 59.00 (BA real ceiling 60.15) must be reachable in the high segment; got %+v", res)
	}
	// must be a REAL price in the high segment (≥79), not merely Reached=true.
	if p := ratOf(t, *res.Preco); p.Cmp(taxaFixaLimiar) < 0 {
		t.Fatalf("solved preço %q is below 79 — expected high segment for this fixture", *res.Preco)
	}
	if got := resim(in, *res.Preco); got != "59.00" {
		t.Fatalf("re-sim (real ICMSCell decompose) margem_pct = %q at preço %q, want 59.00 EXACT", got, *res.Preco)
	}
}

// TestCeilingPctICMSCellBALiteral (Finding 2, A3 review) pins ceilingPct()'s
// literal value for the BA fixture so a badly-wrong or zeroed
// icmsCellAsymptoticRatePct cannot pass the suite silently — the mandated
// golden's targets never approach this value. Verified by hand against D-41
// (icms.go, non-ST branch): aCusto = 20.5/100 = 0.205 (AliquotaInterna
// "20.5", D-43 FCP already embedded); fcp = 0 (FcpEmbutido nil ⇒ 0, D-43);
// aBase = aCusto − fcp = 0.205; rate = pisCofinsRate×(1−aBase) + aCusto =
// 0.0925×0.795 + 0.205 = 0.0735375 + 0.205 = 0.2785375, i.e. 27.85375 in
// percent units; ceiling = 100 − comissão(12) − 27.85375 = 60.14625, which
// FormatRatHalfUp rounds to 60.15 (2dp).
func TestCeilingPctICMSCellBALiteral(t *testing.T) {
	in := baseSolveICMSCell()
	if got := in.ceilingPct(); got != "60.15" {
		t.Fatalf("ceilingPct() for BA cell = %q, want 60.15", got)
	}
}

// TestSolveTargetPriceICMSCellBABetweenCeilings (Finding 2, A3 review) targets
// 70.00 — strictly between the REAL cell ceiling (60.15) and the LEGACY
// fabricated ceiling (100 − comissão(12) − AliquotaPct(4) = 84). If
// ceilingPct/icmsCellAsymptoticRatePct silently fell back to the legacy
// AliquotaPct path (or returned 0), this target would wrongly report
// Reached=true (70 < 84) instead of UNREACHABLE_TARGET citing 60.15 — the
// exact discrimination the mandated golden's targets (which never sit in this
// band) cannot make.
func TestSolveTargetPriceICMSCellBABetweenCeilings(t *testing.T) {
	in := baseSolveICMSCell()
	in.TargetMargemPct = "70.00"

	res := SolveTargetPrice(in)
	if res.Reached || res.Preco != nil {
		t.Fatalf("target 70.00 (between cell ceiling 60.15 and legacy fabricated 84) must be UNREACHABLE; got %+v", res)
	}
	if res.CeilingPct != "60.15" {
		t.Fatalf("CeilingPct = %q, want 60.15 (cell ceiling, never the legacy fabricated 84)", res.CeilingPct)
	}
}

// wantICMSCellUnknown is the Unknown list Decompose/TaxesForItem produce for
// an unresolved OR ambíguo cell (icms.go, decompose.go) — the
// same three components icms_erp_golden_test.go's assertUnknown pins for
// TaxesForItem directly.
var wantICMSCellUnknown = []string{"icms_saida", "difal", "pis_cofins"}

// TestSolveTargetPriceICMSCellAusenteBlocks is the brief's first blocking
// control: célula AUSENTE (CodTrib nil — "célula não resolvida (ausente OU
// nunca lida)", icms.go). The solver must block with the desconhecido
// nomeado Decompose/TaxesForItem already produce — never a preço, never a
// ceiling, and never falling back to the 4% fabricated aliquota.
func TestSolveTargetPriceICMSCellAusenteBlocks(t *testing.T) {
	in := baseSolveICMSCell()
	in.ICMSCell = &ICMSCell{UFDestino: "AC", CodTrib: nil, Ambiguo: false, Origprod: intp(0)}
	in.TargetMargemPct = "20.00"

	res := SolveTargetPrice(in)
	if res.Reached || res.Preco != nil || res.CeilingPct != "" || res.FreteDesconhecido {
		t.Fatalf("célula ausente (CodTrib nil) must block (no preço, no ceiling); got %+v", res)
	}
	if !reflect.DeepEqual(res.Desconhecidos, wantICMSCellUnknown) {
		t.Fatalf("Desconhecidos = %v, want %v", res.Desconhecidos, wantICMSCellUnknown)
	}
}

// TestSolveTargetPriceICMSCellAmbiguaBlocks is the brief's second blocking
// control: célula AMBÍGUA (Ambiguo=true), even with CodTrib/AliquotaInterna
// both present and resolved — icms.go treats Ambiguo the same as
// CodTrib nil, never choosing a candidate blindly. Same blocking shape as
// the ausente control.
func TestSolveTargetPriceICMSCellAmbiguaBlocks(t *testing.T) {
	in := baseSolveICMSCell()
	cell := icmsCellBA()
	cell.Ambiguo = true
	in.ICMSCell = cell
	in.TargetMargemPct = "20.00"

	res := SolveTargetPrice(in)
	if res.Reached || res.Preco != nil || res.CeilingPct != "" || res.FreteDesconhecido {
		t.Fatalf("célula ambígua must block (no preço, no ceiling); got %+v", res)
	}
	if !reflect.DeepEqual(res.Desconhecidos, wantICMSCellUnknown) {
		t.Fatalf("Desconhecidos = %v, want %v", res.Desconhecidos, wantICMSCellUnknown)
	}
}

// --- Task A5: solver stops rejecting reachable targets when credits exist ---
//
// C1 (solve.go): ceilingPct/exactCeilingRat treat "100 − carga proporcional"
// as an upper bound on margem_pct. That only holds while every summed term
// is a COST. With ICMSCell, RestituicaoST is the one component that
// SUBTRACTS from the sum (decompose.go) — a large enough credit pushes
// margem_pct ABOVE the asymptote at a finite preço; the curve approaches the
// asymptote from ABOVE instead of from below. C2: the old high-segment
// bracket (exactMargemPctRat) never accounted for RestituicaoST, nor for
// StRetidoEntrada's effect once the PIS/COFINS MAX(0,…) clamp opens, nor did
// it know the asymptotic PIS/COFINS rate does NOT apply BELOW the clamp.
//
// Both goldens below are the brief's MEASURED reproduction — verbatim
// verified via zz_probe_codex_test.go (TestProbeClaimA_Restituicao100_
// Price18819 / TestProbeClaimA_StRetido500_Price45402): before the Task A5
// fix, SolveTargetPrice(target) returned {Preco:<nil> Reached:false
// CeilingPct:60.15} for BOTH, even though Decompose at the cited preço
// genuinely yields exactly that margem_pct (confirmed by brute force over
// [R$0.01, R$20000.00]).

// TestSolveTargetPriceICMSCellRestituicaoReachesAboveCeiling is the first
// measured golden. RestituicaoUnit=100 against Fixed0=frete(15)+custo(10)=25
// flips the high segment's net fixed coefficient negative (Fnet=25−100=−75),
// so margem_pct approaches the 60.15 ceiling from ABOVE (decreasing in
// preço) and exceeds it for every preço below the (very large) crossing —
// target 100.00 is reached at preço 188.19, the brute-forced cheapest.
func TestSolveTargetPriceICMSCellRestituicaoReachesAboveCeiling(t *testing.T) {
	in := baseSolveICMSCell()
	in.FreteProduto = money("15.00")
	in.ICMSCell.RestituicaoUnit = strptr("100")
	in.TargetMargemPct = "100.00"

	res := SolveTargetPrice(in)
	if !res.Reached || res.Preco == nil {
		t.Fatalf("target 100.00 (RestituicaoUnit=100, real ceiling 60.15) must be reachable ABOVE the ceiling via the credit; got %+v", res)
	}
	if *res.Preco != "188.19" {
		t.Fatalf("solved preço = %s, want 188.19 (measured cheapest, zz_probe_codex_test.go)", *res.Preco)
	}
	if got := resim(in, *res.Preco); got != "100.00" {
		t.Fatalf("re-sim margem_pct = %q at preço %q, want 100.00 EXACT", got, *res.Preco)
	}
}

// TestSolveTargetPriceICMSCellStRetidoClampReachesAboveCeiling is the second
// measured golden. StRetidoEntrada=500 clamps PIS/COFINS to zero for a wide
// band of preço (base = P×0.795 − 500 < 0 below preço ≈ 629.56), so within
// that band the effective ceiling is 67.50 (no PIS/COFINS term at all, not
// 60.15) and the net fixed coefficient there is POSITIVE (Fixed0=25, no
// RestituicaoST credit) — an ordinary increasing crossing, but toward the
// WRONG (higher, un-clamped-blind) asymptote if the old bracket's single
// always-active PIS/COFINS rate is used. Target 62.00 is reached at preço
// 454.02, the brute-forced cheapest, well inside the clamped band.
func TestSolveTargetPriceICMSCellStRetidoClampReachesAboveCeiling(t *testing.T) {
	in := baseSolveICMSCell()
	in.FreteProduto = money("15.00")
	in.ICMSCell.StRetidoEntrada = strptr("500")
	in.TargetMargemPct = "62.00"

	res := SolveTargetPrice(in)
	if !res.Reached || res.Preco == nil {
		t.Fatalf("target 62.00 (StRetidoEntrada=500, real ceiling 60.15, clamped effective ceiling 67.50) must be reachable; got %+v", res)
	}
	if *res.Preco != "454.02" {
		t.Fatalf("solved preço = %s, want 454.02 (measured cheapest, zz_probe_codex_test.go)", *res.Preco)
	}
	if got := resim(in, *res.Preco); got != "62.00" {
		t.Fatalf("re-sim margem_pct = %q at preço %q, want 62.00 EXACT", got, *res.Preco)
	}
}

// TestSolveTargetPriceLegacyCeilingUnchanged is the Task A5 non-regression
// control (brief item 3). The legacy (ICMSCell=nil) branch of Decompose has
// NO credit term at all — comissão, taxa_fixa, imposto, frete, difal,
// tarifa_full, custo are ALL added to sum, never subtracted (decompose.go's
// `in.ICMSCell == nil` arm) — so ceilingPct is a PROVEN upper bound there
// unconditionally (verified by brute force over [R$0.01, R$20000.00] in
// zz_probe_codex_test.go's TestProbeClaimA_LegacyPathCeilingCounterfactual,
// using this exact fixture). C1's fix must not touch this path: 84.00 stays
// the ceiling and a target at/above it stays UNREACHABLE_TARGET via the
// SAME early short-circuit as before Task A5.
func TestSolveTargetPriceLegacyCeilingUnchanged(t *testing.T) {
	in := SolveInput{
		ComissaoPct: "12", AliquotaPct: "4", Modalidade: ModalidadeClassico,
		Custo:        money("10.00"),
		FreteProduto: money("15.00"),
	}
	if got := in.ceilingPct(); got != "84.00" {
		t.Fatalf("legacy ceilingPct() = %q, want 84.00 (fixture unchanged by Task A5)", got)
	}
	in.TargetMargemPct = "84.00"
	res := SolveTargetPrice(in)
	if res.Reached || res.Preco != nil {
		t.Fatalf("target 84.00 (AT the legacy ceiling) must be UNREACHABLE; got %+v", res)
	}
	if res.CeilingPct != "84.00" {
		t.Fatalf("CeilingPct = %q, want 84.00", res.CeilingPct)
	}
}

// icmsCellRegimeGrid is the Task A5 property test's four regimes (brief item
// 2): no credit (baseline — proves the fix does not regress the ordinary
// increasing-toward-ceiling case), RestituicaoUnit (C1, decreasing
// throughout — same fixture as the mandated golden above), StRetidoEntrada
// large enough to clamp PIS/COFINS to zero over part of the high segment
// (C1+C2, increasing-then-decreasing across the clamp transition), and the
// ST branch (isST — a_custo/a_base fold to 0 per icms.go, a
// completely different rate but the SAME two-regime shape). capCents is
// chosen per regime to comfortably straddle its PIS/COFINS clamp transition
// (Pclamp = S/(1−aBase) dollars) where one exists, so the sample below
// exercises BOTH monotone pieces, not just one:
//   - "sem crédito": no StRetidoEntrada ⇒ never clamped, single increasing
//     piece throughout.
//   - "RestituicaoUnit": StRetidoEntrada also unset ⇒ never clamped, single
//     DEcreasing piece (the C1 case).
//   - "StRetidoEntrada clampa PIS/COFINS": aBase=0.205 (BA) ⇒ Pclamp =
//     500/0.795 ≈ R$629.56 (cents 62956); capCents=100000 straddles it.
//   - "ramo ST": aBase=0 (isST) ⇒ Pclamp = S = R$200.00 exactly (cents
//     20000); capCents=60000 straddles it.
func icmsCellRegimeGrid() []struct {
	name     string
	in       SolveInput
	capCents int64
} {
	withFrete := func(cell *ICMSCell) SolveInput {
		in := baseSolveICMSCell()
		in.ICMSCell = cell
		in.FreteProduto = money("15.00")
		return in
	}
	return []struct {
		name     string
		in       SolveInput
		capCents int64
	}{
		{"sem credito", withFrete(icmsCellBA()), 50_000},
		{"RestituicaoUnit", withFrete(&ICMSCell{
			UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
			Origprod: intp(0), AliquotaInterna: strptr("20.5"), FcpEmbutido: strptr("0"),
			RestituicaoUnit: strptr("100"),
		}), 40_000},
		{"StRetidoEntrada clampa PIS/COFINS", withFrete(&ICMSCell{
			UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
			Origprod: intp(0), AliquotaInterna: strptr("20.5"), FcpEmbutido: strptr("0"),
			StRetidoEntrada: strptr("500"),
		}), 100_000},
		{"ramo ST", withFrete(&ICMSCell{
			UFDestino: "RJ", CodTrib: intp(60), Ambiguo: false,
			StRetidoEntrada: strptr("200.00"), RestituicaoUnit: strptr("15.00"),
		}), 60_000},
	}
}

// TestSolveMatchesBruteForceICMSCellFourRegimes is the Task A5 mandated
// strong acceptance test (brief item 2): for each regime, brute-force EVERY
// cent from 1 to capCents through the real Decompose, record the cheapest
// preço reaching each distinct margem_pct, then — for an evenly-strided
// sample of those cents (kept small so the lane runs in acceptable time,
// per the brief) — ask SolveTargetPrice for that exact target and require
// Reached=true at EXACTLY the brute-forced cheapest preço. Zero exceptions:
// every mismatch is reported via t.Errorf (not the first-only t.Fatalf), so
// a run names every failing target instead of stopping at one.
func TestSolveMatchesBruteForceICMSCellFourRegimes(t *testing.T) {
	const samplesPerRegime = 40 // brief: "grade pequena o suficiente para rodar em tempo aceitável"
	start := time.Now()
	totalSamples, totalMismatches := 0, 0
	for _, regime := range icmsCellRegimeGrid() {
		regime := regime
		t.Run(regime.name, func(t *testing.T) {
			cheapest := map[string]int64{}
			for c := int64(1); c <= regime.capCents; c++ {
				m := resim(regime.in, centsToStr(c))
				if m == "<unknown>" {
					continue
				}
				if _, seen := cheapest[m]; !seen {
					cheapest[m] = c
				}
			}
			stride := regime.capCents / samplesPerRegime
			if stride < 1 {
				stride = 1
			}
			samples, mismatches := 0, 0
			for c := int64(1); c <= regime.capCents; c += stride {
				target := resim(regime.in, centsToStr(c))
				if target == "<unknown>" {
					continue
				}
				wantCents := cheapest[target]
				in := regime.in
				in.TargetMargemPct = target
				res := SolveTargetPrice(in)
				samples++
				if !res.Reached || res.Preco == nil {
					t.Errorf("target=%s must reach (brute cheapest %s); got %+v", target, centsToStr(wantCents), res)
					mismatches++
					continue
				}
				if *res.Preco != centsToStr(wantCents) {
					t.Errorf("target=%s solver=%s brute-cheapest=%s", target, *res.Preco, centsToStr(wantCents))
					mismatches++
				}
			}
			t.Logf("%s: %d distinct margins in [1,%d], %d samples, %d mismatches", regime.name, len(cheapest), regime.capCents, samples, mismatches)
			totalSamples += samples
			totalMismatches += mismatches
		})
	}
	t.Logf("TestSolveMatchesBruteForceICMSCellFourRegimes TOTAL wall-clock = %s, samples=%d, mismatches=%d",
		time.Since(start), totalSamples, totalMismatches)
}

// --- Coordinator review (post-A5) — Achado 1: frete=0 fabrication ---------

// TestSolveTargetPriceICMSCellCreditFreteUnknownStaysAmbiguous is the red
// test for the coordinator's Achado 1 (REFUTED
// icmsCellPossiblyReachableAboveCeiling): when ICMSCell has a credit term
// (RestituicaoUnit here, R>0) and FreteProduto is nil, the OLD code probed
// reachability by substituting a FABRICATED frete=0 — which is ADR-17
// backwards (an unknown fact silently became a concrete value) and, when the
// probe came back "unreachable", asserted CeilingPct — the very asymptote
// the fixture's own C1 fix proved is NOT a true ceiling once a credit is
// present. The correct answer never depends on frete's actual (unknown)
// value at all: fnet's sign is decided entirely by the credit terms R/S,
// which are known without frete (icmsCellCreditsAndFixed does not need
// FreteProduto to read them) — R>0 (or S>0) means frete is genuinely needed
// to know which regime/price applies, so the honest answer is
// FreteDesconhecido, never a guessed ceiling.
//
// Target is pinned to "60.15" — the fixture's own real ceiling (see
// TestCeilingPctICMSCellBALiteral) — because at frete=0 the analytic probe's
// bracket sits exactly on its own gStar=0 boundary for this fixture,
// reporting "not reachable at frete=0" (a degenerate coincidence of the
// fabricated substitution, not a fact about the real, unknown frete) and so
// exercises exactly the wrong CeilingPct branch the old code took.
func TestSolveTargetPriceICMSCellCreditFreteUnknownStaysAmbiguous(t *testing.T) {
	in := baseSolveICMSCell()
	in.ICMSCell.RestituicaoUnit = strptr("100")
	in.TargetMargemPct = "60.15"
	// in.FreteProduto stays nil — the fact under test.

	res := SolveTargetPrice(in)
	if !res.FreteDesconhecido {
		t.Fatalf("célula COM crédito (RestituicaoUnit=100) + FreteProduto desconhecido + alvo na assíntota real: "+
			"frete é genuinamente necessário para decidir (fnet muda de sinal conforme o frete real) — "+
			"esperava FreteDesconhecido=true, got %+v", res)
	}
	if res.CeilingPct != "" {
		t.Fatalf("FreteDesconhecido=true não deve vir acompanhado de CeilingPct (o solver não pode afirmar um teto "+
			"que não investigou) — got CeilingPct=%q", res.CeilingPct)
	}
	if res.Reached || res.Preco != nil {
		t.Fatalf("não deveria ter Reached/Preco junto de FreteDesconhecido; got %+v", res)
	}
}

// TestSolveTargetPriceICMSCellNoCreditFreteUnknownStillCeilings is the
// non-regression control for Achado 1's fix: célula SEM crédito nenhum
// (icmsCellBA "as-is", nem RestituicaoUnit nem StRetidoEntrada — the SAME
// fixture TestSolveTargetPriceICMSCellBABetweenCeilings already uses) means
// R=0 ∧ S=0, so fnet = −fixed0 ≤ 0 for EVERY real frete≥0 (frete only adds
// to fixed0, never subtracts) — the asymptote 60.15 is a true ceiling
// regardless of the unknown frete's actual value, so CeilingPct (not
// FreteDesconhecido) remains the honest, provable answer. This duplicates
// TestSolveTargetPriceICMSCellBABetweenCeilings's assertion deliberately —
// it is the control half of Achado 1's red/green pair, kept next to the new
// credit test so the sign-based branch (R==0∧S==0 vs R>0∨S>0) is pinned by
// two tests that sit side by side, not one alone in a different file section.
func TestSolveTargetPriceICMSCellNoCreditFreteUnknownStillCeilings(t *testing.T) {
	in := baseSolveICMSCell()
	in.TargetMargemPct = "70.00" // acima da assíntota real 60.15

	res := SolveTargetPrice(in)
	if res.FreteDesconhecido {
		t.Fatalf("célula SEM crédito (R=0, S=0): frete não pode mudar o veredito (fnet=−fixed0≤0 para todo "+
			"frete≥0) — esperava CeilingPct honesto, got FreteDesconhecido=true, %+v", res)
	}
	if res.CeilingPct != "60.15" {
		t.Fatalf("CeilingPct = %q, want 60.15 (assíntota real da célula, sem crédito)", res.CeilingPct)
	}
	if res.Reached || res.Preco != nil {
		t.Fatalf("target 70.00 acima da assíntota real 60.15 não deveria ser Reached; got %+v", res)
	}
}

// TestSolveTargetPriceICMSCellNoCreditFreteUnknownBelowCeilingIsNotUnreachable
// é a OUTRA metade do gate acima, e ela estava descoberta — o dual gate (lado
// Opus) achou, e a medição confirmou com margem maior que a relatada.
//
// solve.go pula a comparação alvo-vs-teto quando ICMSCell != nil (o teto
// não é limite superior verdadeiro na presença de crédito — C1). O atalho de
// frete desconhecido em solve.go então devolvia CeilingPct sempre que
// R≤0 ∧ S≤0, SEM NUNCA comparar o alvo com esse teto. Ninguém rio acima faz
// essa comparação, e o transporte (calc_handler.go) mapeia
// `!Reached && CeilingPct != ""` para UNREACHABLE_TARGET.
//
// Resultado medido nesta MESMA fixture (sem crédito, frete desconhecido):
//
//	alvo   frete=nil                       frete=15.00 (real)
//	40.00  UNREACHABLE, teto 60.15         Reached, preço 124.03
//	50.00  UNREACHABLE, teto 60.15         Reached, preço 246.18
//	55.00  UNREACHABLE, teto 60.15         Reached, preço 485.17
//	70.00  UNREACHABLE, teto 60.15         UNREACHABLE, teto 60.15  <- correto
//
// Três alvos ESTRITAMENTE ABAIXO do teto recebiam "não existe preço" enquanto
// o preço existe — a tela diria "inalcançável, teto 60,15%" para um alvo de
// 50%. É a mesma classe que A5 dizia consertar (veredito de inexistência
// fabricado a partir de um fato não lido), só que mudou de sítio: saiu da
// sonda com frete=0 sintético e entrou no atalho que a substituiu.
//
// O que R≤0 ∧ S≤0 realmente prova é só que o teto é HONESTO — não que o alvo
// esteja acima dele. Quando o alvo fica abaixo, o segmento alto é quem
// decidiria, e ele precisa do frete: a resposta é FreteDesconhecido.
//
// O teste acima (alvo 70.00) só exercitava a metade que por acaso acertava.
func TestSolveTargetPriceICMSCellNoCreditFreteUnknownBelowCeilingIsNotUnreachable(t *testing.T) {
	for _, target := range []string{"40.00", "50.00", "55.00"} {
		in := baseSolveICMSCell()
		in.TargetMargemPct = target

		res := SolveTargetPrice(in)

		if !res.FreteDesconhecido {
			t.Errorf("alvo %s (abaixo do teto 60.15) com frete desconhecido: esperava FreteDesconhecido=true — "+
				"o segmento alto é quem decide e ele precisa do frete; got %+v", target, res)
		}
		if res.CeilingPct != "" {
			t.Errorf("alvo %s: CeilingPct=%q com o alvo ABAIXO dele — o transporte mapeia isso para "+
				"UNREACHABLE_TARGET (calc_handler.go) e afirma que o preço não existe quando existe",
				target, res.CeilingPct)
		}
	}

	// Controle positivo: o preço que o veredito nega. Com um frete real, os
	// mesmos três alvos são alcançáveis — então "inalcançável" era falso, não
	// conservador.
	for _, tc := range []struct{ target, preco string }{
		{"40.00", "124.03"}, {"50.00", "246.18"}, {"55.00", "485.17"},
	} {
		in := baseSolveICMSCell()
		in.FreteProduto = money("15.00")
		in.TargetMargemPct = tc.target

		res := SolveTargetPrice(in)

		if !res.Reached || res.Preco == nil || *res.Preco != tc.preco {
			t.Errorf("controle positivo, alvo %s com frete 15.00: esperava Reached preço %s; got %+v",
				tc.target, tc.preco, res)
		}
	}
}
