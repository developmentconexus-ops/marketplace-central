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

// C03(d) — monotonic-bracket invariant: within a single segment margem_pct is
// strictly increasing in preço, which is the precondition the per-segment
// binary search relies on. Sampled across the low segment (< 79) and the high
// segment (≥ 79) separately (the jump lives only at the segment boundary).
func TestSolveMargemMonotonicWithinSegment(t *testing.T) {
	in := baseSolve()
	assertIncreasing := func(t *testing.T, from, to, step int64) {
		t.Helper()
		prev := ""
		for c := from; c <= to; c += step {
			cur := resim(in, centsToStr(c))
			if cur == "<unknown>" {
				t.Fatalf("unexpected unknown margem at cents %d", c)
			}
			if prev != "" && cmpPct(cur, prev) < 0 {
				t.Fatalf("margem_pct not monotonic at cents %d: %q < previous %q", c, cur, prev)
			}
			prev = cur
		}
	}
	assertIncreasing(t, 5000, 7899, 137)    // low segment, preço < 79
	assertIncreasing(t, 7900, 500000, 4211) // high segment, preço ≥ 79
}
