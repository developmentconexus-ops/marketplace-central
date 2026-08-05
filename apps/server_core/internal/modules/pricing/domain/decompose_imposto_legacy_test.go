package domain

import "testing"

// TestDecomposeICMSCellPresentImpostoAbsent is Task A6's D1 vector: when
// ICMSCell is present, the legacy AliquotaPct is NOT summed (that half of D1
// was already fixed) AND is no longer even shown — Imposto comes back as the
// named-absent nil, never a fabricated number sitting next to the real
// ICMSSaida/Difal. Asserting the ABSENCE STATE (nil), not `!= "4.00"`: an
// inequality assertion would still pass if some OTHER fabricated number
// leaked out.
//
// P=100, BA (AliquotaInterna=aCusto=20.5%, not in aInterUF12 ⇒ a_inter=7%
// default, FcpEmbutido nil ⇒ 0 so aBase=aCusto=0.205). ICMS_oper (ICMSSaida)
// = 100×0.07 = 7.00. ICMS_total = 100×0.205 = 20.50. DIFAL = 20.50−7.00 =
// 13.50 — matches the brief's measured D1 numbers exactly. Comissao=100×0=
// 0.00 (isolate the ICMS/DIFAL pair). Custo=0.00 known so margem is
// observable too: PisCofins base = 100×(1−aBase)−S = 100×0.795−0 = 79.50,
// ×0.0925 = 7.35375 -> 7.35 (half-up). RestituicaoST=0.00 (not set,
// nil-as-zero, UF≠MG). Sum = 0(comissao)+0(taxa,>=79)+0(frete,>=79 but
// FreteProduto known 0)+7.00(icms)+13.50(difal)+7.35(pis)+0(tarifa)+
// 0(custo)-0(restituicao) = 27.85. Margem=100-27.85=72.15. Pct=72.15.
func TestDecomposeICMSCellPresentImpostoAbsent(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("20.5"),
		// Task A4 (D2): par completo -- "0" explicito, nao nil.
		// AliquotaInterna presente com FcpEmbutido nil e um estado que a
		// fonte real nunca produz; numericamente identico (fcp=0 nos dois
		// casos), so a REALISMO do fixture muda -- MargemValor abaixo nao
		// muda.
		FcpEmbutido: strptr("0"),
	}
	in := DecomposeInput{
		Preco: "100.00", ComissaoPct: "0", AliquotaPct: "4",
		Modalidade:   ModalidadeClassico,
		Custo:        money("0.00"),
		FreteProduto: money("0.00"),
		ICMSCell:     cell,
	}

	got := Decompose(in)

	if got.Imposto != nil {
		t.Fatalf("Imposto = %q, want nil (named-absent) — legacy path must not run when ICMSCell is present", *got.Imposto)
	}
	if got.ICMSSaida == nil || *got.ICMSSaida != "7.00" {
		t.Fatalf("ICMSSaida = %v, want 7.00 (real component, unaffected by Task A6)", derefOrNil(got.ICMSSaida))
	}
	if got.Difal == nil || *got.Difal != "13.50" {
		t.Fatalf("Difal = %v, want 13.50 (real component, unaffected by Task A6)", derefOrNil(got.Difal))
	}
	// Imposto's absence must not become an ADR-17 "unknown" that blocks
	// margem — it is a structural does-not-apply, not a missing fact.
	for _, c := range got.ComponentesDesconhecidos {
		if c == "imposto" {
			t.Fatalf("ComponentesDesconhecidos = %v; Imposto's absence is structural (ICMSCell path), never an unknown component", got.ComponentesDesconhecidos)
		}
	}
	if got.MargemValor == nil || *got.MargemValor != "72.15" {
		t.Fatalf("MargemValor = %v, want 72.15 (unchanged account: comissao+icms+difal+pis)", derefOrNil(got.MargemValor))
	}
}

// derefOrNil formats a *string for test failure messages without panicking
// on nil.
func derefOrNil(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}

// TestDecomposeICMSCellPresentEmptyAliquotaPctDoesNotPanic is Task A6's D2
// vector: an otherwise-complete ICMSCell with AliquotaPct="" must NOT panic
// (today it does: pctOfPrice("", preco) -> ParseRat("") -> panic). The pair
// (AliquotaPct="" vs AliquotaPct="4") is the control that proves the legacy
// field is genuinely OUT of the path, not just suppressed in the output —
// both inputs must produce byte-identical fiscal components.
//
// Dual gate (Opus side): the fixture used to omit FcpEmbutido, so A4's D2
// gate fired on BOTH sides and PisCofins came back nil in each — the
// comparison then held over two decompositions whose fiscal arithmetic never
// ran. A vacuous pass, and precisely on the component A8 changed. The pair is
// now D2-clean, so the equality is over computed numbers.
func TestDecomposeICMSCellPresentEmptyAliquotaPctDoesNotPanic(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("20.5"),
		FcpEmbutido: strptr("0"),
	}
	base := DecomposeInput{
		Preco: "100.00", ComissaoPct: "0",
		Modalidade:   ModalidadeClassico,
		Custo:        money("0.00"),
		FreteProduto: money("0.00"),
		ICMSCell:     cell,
	}

	withEmpty := base
	withEmpty.AliquotaPct = ""
	withFour := base
	withFour.AliquotaPct = "4"

	var gotEmpty, gotFour Decomposition
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Decompose(AliquotaPct=\"\", ICMSCell present) panicked: %v — the legacy field must not be parsed at all when the cell is resolved", r)
			}
		}()
		gotEmpty = Decompose(withEmpty)
	}()
	gotFour = Decompose(withFour)

	if gotEmpty.Imposto != nil || gotFour.Imposto != nil {
		t.Fatalf("Imposto must be nil regardless of legacy AliquotaPct when ICMSCell present: empty=%v four=%v", gotEmpty.Imposto, gotFour.Imposto)
	}
	if !decompositionEqualIgnoringImposto(gotEmpty, gotFour) {
		t.Fatalf("AliquotaPct=\"\" and AliquotaPct=\"4\" must yield identical fiscal components under ICMSCell:\n empty=%+v\n four =%+v", gotEmpty, gotFour)
	}
	// Guard against the vacuous version of this test coming back: if the
	// fixture ever trips a gate again, PisCofins goes nil on both sides and
	// the equality above holds without the arithmetic running.
	if gotEmpty.PisCofins == nil {
		t.Fatalf("PisCofins = nil — o fixture voltou a ser bloqueado por um gate, e a igualdade acima passa por vacuidade; desconhecidos=%v", gotEmpty.ComponentesDesconhecidos)
	}
}

// TestDecomposeICMSCellAbsentLegacyUnchanged is Task A6's mandatory NEGATIVE
// control: when ICMSCell is nil, the legacy path is completely untouched —
// Imposto is still summed and rendered exactly as before (now behind a
// pointer, but the same "4.00"). If this test breaks, the fix killed the
// entire legacy path instead of just the cell/legacy CROSSING.
func TestDecomposeICMSCellAbsentLegacyUnchanged(t *testing.T) {
	in := DecomposeInput{
		Preco: "100.00", ComissaoPct: "0", AliquotaPct: "4",
		Modalidade:   ModalidadeClassico,
		Custo:        money("0.00"),
		FreteProduto: money("0.00"),
		DifalEnabled: false,
		ICMSCell:     nil,
	}

	got := Decompose(in)

	if got.Imposto == nil || *got.Imposto != "4.00" {
		t.Fatalf("Imposto = %v, want \"4.00\" (ICMSCell absent ⇒ legacy path identical to today)", got.Imposto)
	}
}

// decompositionEqualIgnoringImposto compares every field except Imposto
// (already asserted nil==nil by the caller) field by field — avoids
// reflect.DeepEqual pulling in Imposto's pointer identity, which is
// irrelevant here (both are nil anyway; this guards the OTHER fields).
//
// ComponentesDesconhecidos is compared too (dual gate, Opus side): omitting
// it let the two sides disagree about WHICH facts are unknown while every
// value matched — the one divergence this pair exists to catch.
func decompositionEqualIgnoringImposto(a, b Decomposition) bool {
	eqStrPtr := func(x, y *string) bool {
		if x == nil || y == nil {
			return x == y
		}
		return *x == *y
	}
	if len(a.ComponentesDesconhecidos) != len(b.ComponentesDesconhecidos) {
		return false
	}
	for i := range a.ComponentesDesconhecidos {
		if a.ComponentesDesconhecidos[i] != b.ComponentesDesconhecidos[i] {
			return false
		}
	}
	return a.Preco == b.Preco &&
		a.Comissao == b.Comissao &&
		a.TaxaFixa == b.TaxaFixa &&
		eqStrPtr(a.Frete, b.Frete) &&
		eqStrPtr(a.Difal, b.Difal) &&
		eqStrPtr(a.TarifaFull, b.TarifaFull) &&
		eqStrPtr(a.Custo, b.Custo) &&
		eqStrPtr(a.ICMSSaida, b.ICMSSaida) &&
		eqStrPtr(a.PisCofins, b.PisCofins) &&
		eqStrPtr(a.RestituicaoST, b.RestituicaoST) &&
		eqStrPtr(a.MargemValor, b.MargemValor) &&
		eqStrPtr(a.MargemPct, b.MargemPct)
}
