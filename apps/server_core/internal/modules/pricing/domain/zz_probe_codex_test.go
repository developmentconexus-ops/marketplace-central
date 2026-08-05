package domain

import (
	"fmt"
	"math/big"
	"testing"
)

// zz_probe_codex_test.go holds the measurements that Fatia A's goldens cite as
// their provenance ("verbatim verified via zz_probe_codex_test.go", the 188.19
// and 454.02 numbers in solve_golden_test.go, and the legacy-ceiling
// counterfactual). It is committed for exactly that reason: the dual gate
// (Opus side) found the citations pointing at a file that existed in neither
// the tree nor the range, which makes a stated measurement unreachable and
// therefore a number taken on trust. The orders slice already did this right
// (zz_probe_1005_test.go, zz_probe_linetotal_test.go sit next to the comments
// citing them); this is the same rule applied here.
//
// Excluded from the normal lane by -skip 'TestProbe'. These probes log and
// assert measured behavior; they are not contract goldens.
//
// ============================================================================
// Claim A — solver rejects a REACHABLE target because the asymptote is not
// an upper ceiling.
// ============================================================================

// probeClaimAInput builds the BA-cell SolveInput used by both A instances,
// mirroring baseSolveICMSCell() (solve_golden_test.go) but with a caller
// supplied ICMSCell mutation and FreteProduto.
func probeClaimAInput(mutate func(c *ICMSCell)) SolveInput {
	in := baseSolveICMSCell()
	mutate(in.ICMSCell)
	in.FreteProduto = money("15.00")
	return in
}

func probeClaimADecomposeInput(in SolveInput, preco string) DecomposeInput {
	return DecomposeInput{
		Preco: preco, ComissaoPct: in.ComissaoPct, AliquotaPct: in.AliquotaPct,
		Modalidade: in.Modalidade, TarifaFull: in.TarifaFull,
		FreteProduto: in.FreteProduto, Custo: in.Custo,
		DifalEnabled: in.DifalEnabled, DestinoUF: in.DestinoUF, EfetivoPct: in.EfetivoPct,
		ICMSCell: in.ICMSCell,
	}
}

func TestProbeClaimA_Restituicao100_Price18819(t *testing.T) {
	in := probeClaimAInput(func(c *ICMSCell) { c.RestituicaoUnit = strptr("100") })

	d := Decompose(probeClaimADecomposeInput(in, "188.19"))
	t.Logf("Decompose(188.19) full struct = %+v", d)
	if d.MargemPct != nil {
		t.Logf("Decompose(188.19).MargemPct = %q", *d.MargemPct)
	} else {
		t.Logf("Decompose(188.19).MargemPct = nil, ComponentesDesconhecidos=%v", d.ComponentesDesconhecidos)
	}

	in.TargetMargemPct = "100.00"
	res := SolveTargetPrice(in)
	t.Logf("SolveTargetPrice(target=100.00) = %+v", res)

	// brute-force scan 1 cent .. R$20000 (2,000,000 cents) for ANY price
	// whose Decompose margem_pct == 100.00.
	foundCents := int64(-1)
	for c := int64(1); c <= 2_000_000; c++ {
		dd := Decompose(probeClaimADecomposeInput(in, centsToStr(c)))
		if dd.MargemPct != nil && *dd.MargemPct == "100.00" {
			foundCents = c
			break
		}
	}
	if foundCents >= 0 {
		t.Logf("brute-force: cheapest price yielding margem 100.00 = %s (cents=%d)", centsToStr(foundCents), foundCents)
	} else {
		t.Logf("brute-force: none found for margem 100.00 in [0.01, 20000.00]")
	}
}

func TestProbeClaimA_StRetido500_Price45402(t *testing.T) {
	in := probeClaimAInput(func(c *ICMSCell) { c.StRetidoEntrada = strptr("500") })

	d := Decompose(probeClaimADecomposeInput(in, "454.02"))
	t.Logf("Decompose(454.02) full struct = %+v", d)
	if d.PisCofins != nil {
		t.Logf("Decompose(454.02).PisCofins = %q", *d.PisCofins)
	}
	if d.MargemPct != nil {
		t.Logf("Decompose(454.02).MargemPct = %q", *d.MargemPct)
	} else {
		t.Logf("Decompose(454.02).MargemPct = nil, ComponentesDesconhecidos=%v", d.ComponentesDesconhecidos)
	}

	in.TargetMargemPct = "62.00"
	res := SolveTargetPrice(in)
	t.Logf("SolveTargetPrice(target=62.00) = %+v", res)

	foundCents := int64(-1)
	for c := int64(1); c <= 2_000_000; c++ {
		dd := Decompose(probeClaimADecomposeInput(in, centsToStr(c)))
		if dd.MargemPct != nil && *dd.MargemPct == "62.00" {
			foundCents = c
			break
		}
	}
	if foundCents >= 0 {
		t.Logf("brute-force: cheapest price yielding margem 62.00 = %s (cents=%d)", centsToStr(foundCents), foundCents)
	} else {
		t.Logf("brute-force: none found for margem 62.00 in [0.01, 20000.00]")
	}
}

// ============================================================================
// Claim B — an ambiguous CST 60 cell silently becomes explicit zeros.
// ============================================================================

func TestProbeClaimB_AmbiguousST60(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "SP", CodTrib: intp(60), Ambiguo: true,
	}
	out := TaxesForItem("100", cell)
	t.Logf("TaxesForItem(P=100, UFDestino=SP, CodTrib=60, Ambiguo=true) full struct = %+v", out)

	printOrNil := func(name string, s *string) {
		if s == nil {
			t.Logf("  %s = nil", name)
		} else {
			t.Logf("  %s = %q", name, *s)
		}
	}
	printOrNil("ICMSSaida", out.ICMSSaida)
	printOrNil("Difal", out.Difal)
	printOrNil("FCP", out.FCP)
	printOrNil("PisCofins", out.PisCofins)
	printOrNil("RestituicaoST", out.RestituicaoST)
	printOrNil("FiscalDtRef", out.FiscalDtRef)
	t.Logf("  Unknown = %#v (len=%d)", out.Unknown, len(out.Unknown))
}

// ============================================================================
// Claim C — the fabricated flat rate is still reachable and observable
// through Decompose (C1), and an otherwise-complete ICMSCell with
// AliquotaPct="" panics (C2).
// ============================================================================

func TestProbeClaimC1_FabricatedFlatRateStillExposed(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("20.5"),
	}
	in := DecomposeInput{
		Preco: "100", ComissaoPct: "0", AliquotaPct: "4",
		Modalidade: ModalidadeClassico,
		Custo:      money("0.00"),
		ICMSCell:   cell,
	}
	d := Decompose(in)
	t.Logf("Decompose(P=100, BA aCusto=20.5, legacy AliquotaPct=4) full struct = %+v", d)
	// %v not %q: Task A6 changed Imposto to *string (named-absent nil when
	// ICMSCell is present) — go vet's printf check rejects %q on a pointer,
	// which turns into a build failure under `go test`. Formatting-only
	// touch, the probe's own logic/assertions are untouched.
	t.Logf("  Imposto = %v", d.Imposto)
	if d.ICMSSaida != nil {
		t.Logf("  ICMSSaida = %q", *d.ICMSSaida)
	} else {
		t.Logf("  ICMSSaida = nil")
	}
	if d.Difal != nil {
		t.Logf("  Difal = %q", *d.Difal)
	} else {
		t.Logf("  Difal = nil")
	}
	if d.MargemPct != nil {
		t.Logf("  MargemPct = %q", *d.MargemPct)
	} else {
		t.Logf("  MargemPct = nil, ComponentesDesconhecidos=%v", d.ComponentesDesconhecidos)
	}
}

func TestProbeClaimC2_EmptyAliquotaPctPanics(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("20.5"),
	}
	in := DecomposeInput{
		Preco: "100", ComissaoPct: "12", AliquotaPct: "", // <-- legacy rate left empty
		Modalidade: ModalidadeClassico,
		Custo:      money("10.00"),
		ICMSCell:   cell,
	}

	var recovered any
	func() {
		defer func() {
			recovered = recover()
		}()
		d := Decompose(in)
		t.Logf("Decompose did NOT panic; result = %+v", d)
	}()

	if recovered != nil {
		t.Logf("Decompose(AliquotaPct=\"\", otherwise-complete ICMSCell) PANICKED: %v", fmt.Sprint(recovered))
	} else {
		t.Logf("Decompose(AliquotaPct=\"\", otherwise-complete ICMSCell) did NOT panic")
	}
}

// ============================================================================
// Claim D — the order tax money path crosses float64.
//
// orders/adapters/pricingtax/reader.go's lineTotal (reader.go) is
// exactly:
//
//	v := new(big.Rat).SetFloat64(unitPrice)
//	if v == nil { v = new(big.Rat) }
//	v.Mul(v, big.NewRat(int64(quantity), 1))
//	return pricingdomain.FormatRatHalfUp(v, 2)
//
// FormatRatHalfUp is exported by THIS package (decimal.go) and is the
// exact same routine reader.go calls — so this probe replicates reader.go's
// money path byte-for-byte using only exported domain-package symbols,
// without touching or importing the orders package.
// ============================================================================

func TestProbeClaimD_UnitPriceFloat64Crossing(t *testing.T) {
	const unitPrice = 1.005
	const quantity = 1

	// exact literal Go source of reader.go (lineTotal), substituting
	// unitPrice/quantity for item.UnitPrice/item.Quantity.
	v := new(big.Rat).SetFloat64(unitPrice)
	if v == nil {
		t.Fatal("SetFloat64(1.005) returned nil — unexpected for a finite literal")
	}
	t.Logf("new(big.Rat).SetFloat64(1.005) exact binary value = %s (%s)", v.RatString(), v.FloatString(20))
	v.Mul(v, big.NewRat(int64(quantity), 1))
	got := FormatRatHalfUp(v, 2)
	t.Logf("lineTotal(unitPrice=1.005, quantity=1) replica = FormatRatHalfUp(SetFloat64(1.005), 2) = %q", got)

	// reference: the decimal-exact value ParseRat("1.005") would carry (what
	// a string-sourced money path, never crossing float64, would round).
	exact := mustRat("1.005")
	exactFmt := FormatRatHalfUp(exact, 2)
	t.Logf("reference: FormatRatHalfUp(ParseRat(\"1.005\"), 2) [decimal-exact, no float64] = %q", exactFmt)

	if got != exactFmt {
		t.Logf("DIVERGENCE CONFIRMED: float64-sourced path yields %q, decimal-exact path yields %q", got, exactFmt)
	} else {
		t.Logf("no divergence at this input: both paths yield %q", got)
	}
}

// ============================================================================
// Claim A — legacy-path counterfactual (coordinator follow-up): is the
// ceiling-violation defect NEW to the D-41/ICMSCell solver path added in
// ecb13aa3..db3f1d90 (Task A3), or already present on the legacy flat-rate
// path?
//
// STRUCTURAL FINDING (measured by reading, confirmed by this probe compiling
// at all): RestituicaoUnit and StRetidoEntrada are fields of *ICMSCell
// (icms.go) ONLY. Neither DecomposeInput (decompose.go) nor
// SolveInput (solve.go) carries a restituição/st_retido field outside
// ICMSCell — so when ICMSCell is nil there is NO field to attach "the same
// credits" to; the request's instance-1 (RestituicaoUnit="100") and
// instance-2 (StRetidoEntrada="500") variants are therefore structurally
// IDENTICAL on the legacy path (both collapse to the same input below) —
// reported as one measurement, not duplicated to manufacture two different
// numbers.
//
// This is itself the load-bearing fact for the NEW-vs-PRE-EXISTING question:
// the legacy branch of decomposeWithLimiar (decompose.go, the
// `in.ICMSCell == nil` arm) has NO credit term at all — comissão, taxa_fixa,
// imposto, frete, difal, tarifa_full, custo are ALL *added* to `sum`
// (decompose.go,167,173,182,188,232,246); nothing is ever *subtracted*
// from `sum` the way RestituicaoST is on the ICMSCell branch
// (decompose.go, `sum.Sub(sum, mustRat(*tax.RestituicaoST))`). Since
// margem_pct = 100·(preço−sum)/preço and every summed legacy term is a cost
// (Sign() >= 0) added once, margem_pct = 100 − comissãoPct − aliquotaPct −
// (difalPct if applied) − 100·fixed/preço ≤ ceilingPct for EVERY preço, by
// construction — there is no mechanism on the legacy path (unlike the
// RestituicaoST credit on the ICMSCell path, Claim A) that could push margem
// above the asymptote at any finite preço. This probe brute-forces that
// prediction directly rather than trusting the derivation.
// ============================================================================

func probeClaimALegacyInput() SolveInput {
	return SolveInput{
		ComissaoPct: "12", AliquotaPct: "4", Modalidade: ModalidadeClassico,
		Custo:        money("10.00"),
		FreteProduto: money("15.00"),
		// ICMSCell intentionally nil (legacy path). RestituicaoUnit="100" /
		// StRetidoEntrada="500" from Claim A have no field to land in here —
		// see the structural finding above. Omitted, not zeroed-and-hidden.
	}
}

func TestProbeClaimA_LegacyPathCeilingCounterfactual(t *testing.T) {
	in := probeClaimALegacyInput()

	ceiling := in.ceilingPct()
	t.Logf("legacy SolveInput.ceilingPct() = %q  (ICMSCell=nil, ComissaoPct=%q, AliquotaPct=%q, DifalEnabled=%v)",
		ceiling, in.ComissaoPct, in.AliquotaPct, in.DifalEnabled)

	// brute-force scan 1 cent .. R$20000 (2,000,000 cents) for ANY price
	// whose legacy Decompose margem_pct is STRICTLY ABOVE the ceiling
	// SolveTargetPrice would report for this exact input.
	foundCents := int64(-1)
	foundMargem := ""
	for c := int64(1); c <= 2_000_000; c++ {
		d := decomposeWithLimiar(DecomposeInput{
			Preco: centsToStr(c), ComissaoPct: in.ComissaoPct, AliquotaPct: in.AliquotaPct,
			Modalidade: in.Modalidade, TarifaFull: in.TarifaFull,
			FreteProduto: in.FreteProduto, Custo: in.Custo,
			DifalEnabled: in.DifalEnabled, DestinoUF: in.DestinoUF, EfetivoPct: in.EfetivoPct,
		}, taxaFixaLimiar)
		if d.MargemPct == nil {
			continue
		}
		if cmpPct(*d.MargemPct, ceiling) > 0 {
			foundCents = c
			foundMargem = *d.MargemPct
			break
		}
	}

	if foundCents < 0 {
		t.Logf("brute-force [R$0.01, R$20000.00]: NO price found with legacy margem_pct strictly above ceiling %q", ceiling)
		t.Logf("legacy path does NOT reproduce Claim A's defect in this range: the ceiling holds as a true upper bound")
		return
	}

	t.Logf("brute-force: price %s (cents=%d) achieves legacy margem_pct=%q, ABOVE ceiling %q",
		centsToStr(foundCents), foundCents, foundMargem, ceiling)
	in.TargetMargemPct = foundMargem
	res := SolveTargetPrice(in)
	t.Logf("SolveTargetPrice(target=%s) on legacy path = %+v", foundMargem, res)
}
