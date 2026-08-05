package domain

import "testing"

// TestProbeGateFreteUnknownBelowCeiling measures the dual gate's blocking claim:
// on the ICMSCell path with frete unknown and R<=0 && S<=0, solve.go
// returns CeilingPct WITHOUT comparing the target to it, because the ceilingPct short-circuit's
// comparison is skipped for ICMSCell. Decisive pair: same input, frete unknown
// vs frete known.
func TestProbeGateFreteUnknownBelowCeiling(t *testing.T) {
	for _, target := range []string{"30.00", "40.00", "50.00", "55.00", "70.00"} {
		in := baseSolveICMSCell()
		in.TargetMargemPct = target
		res := SolveTargetPrice(in)
		preco := "<nil>"
		if res.Preco != nil {
			preco = *res.Preco
		}
		t.Logf("frete=nil  target=%s -> Reached=%v Preco=%s Ceiling=%q FreteDesconhecido=%v",
			target, res.Reached, preco, res.CeilingPct, res.FreteDesconhecido)
	}

	// icmsCellRates reads *cell.AliquotaInterna unguarded. Measure whether a
	// nil AliquotaInterna can reach it through SolveTargetPrice, or whether
	// structuralUnknowns fires first.
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("PANIC with AliquotaInterna=nil: %v", r)
			}
		}()
		in := baseSolveICMSCell()
		in.ICMSCell.AliquotaInterna = nil
		in.TargetMargemPct = "20.00"
		res := SolveTargetPrice(in)
		t.Logf("AliquotaInterna=nil -> Desconhecidos=%v Ceiling=%q Reached=%v",
			res.Desconhecidos, res.CeilingPct, res.Reached)
	}()

	for _, target := range []string{"30.00", "40.00", "50.00", "55.00", "70.00"} {
		in := baseSolveICMSCell()
		in.FreteProduto = money("15.00")
		in.TargetMargemPct = target
		res := SolveTargetPrice(in)
		preco := "<nil>"
		if res.Preco != nil {
			preco = *res.Preco
		}
		t.Logf("frete=15.00 target=%s -> Reached=%v Preco=%s Ceiling=%q FreteDesconhecido=%v",
			target, res.Reached, preco, res.CeilingPct, res.FreteDesconhecido)
	}
}
