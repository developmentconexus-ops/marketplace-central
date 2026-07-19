package domain

import "testing"

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
// return the margem_pct the simulator would show.
func resim(in SolveInput, preco string) string {
	d := Decompose(DecomposeInput{
		Preco: preco, ComissaoPct: in.ComissaoPct, AliquotaPct: in.AliquotaPct,
		Modalidade: in.Modalidade, TarifaFull: in.TarifaFull,
		FreteProduto: in.FreteProduto, Custo: in.Custo,
		DifalEnabled: in.DifalEnabled, DestinoUF: in.DestinoUF, EfetivoPct: in.EfetivoPct,
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
		i := 0
		for target, wantCents := range cheapest {
			i++
			if i%50 != 0 { // sample to bound runtime; map order varies across runs
				continue
			}
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
