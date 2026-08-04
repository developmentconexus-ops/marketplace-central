package domain

import "testing"

func intp(v int) *int { return &v }

// topsOperacao mirrors what the Oracle adapter passes in production
// (icmsMatrixTOPs): 313 pedido + 306 NFE ENTREGA E-COMMERCE.
var topsOperacao = map[int]bool{306: true, 313: true}

// TestResolveCellDiscardsRowRestrictedToAnotherOperation uses the real DF rows
// measured in METALPRD.TGFICM. "O/103 + I/122" is a rule belonging to
// CODTIPOPER 103 — an operation with zero documents in 2026 — and it carries
// REDBASE=100, CODTRIB=60. The original predicate read only the second pair,
// saw I/122 == grupo, and accepted it, so it competed with the operation's own
// N/0 + I/122 rule and made the cell ambiguous. Only the N row may survive.
func TestResolveCellDiscardsRowRestrictedToAnotherOperation(t *testing.T) {
	aliquota, redBaseOutraOp, redBase := 7.0, 100.0, 0.0
	rows := []TgficmRow{
		{ // rule for CODTIPOPER 103 — not ours
			TipRestricao: "O", CodRestricao: intp(103),
			TipRestricao2: "I", CodRestricao2: intp(122),
			CodTrib: intp(60), Aliquota: &aliquota, RedBase: &redBaseOutraOp,
		},
		{ // our operation's rule: no operation restriction, grupo 122
			TipRestricao: "N", CodRestricao: intp(0),
			TipRestricao2: "I", CodRestricao2: intp(122),
			CodTrib: intp(0), Aliquota: &aliquota, RedBase: &redBase,
		},
	}

	got := ResolveCell(rows, 122, topsOperacao)

	if !got.Found {
		t.Fatalf("Found = false, want true (the N/0 + I/122 rule applies)")
	}
	if got.Ambiguo {
		t.Fatalf("Ambiguo = true (%d candidatas), want false — the O/103 row belongs to another operation", got.LinhasCandidatas)
	}
	if got.RedBase == nil || *got.RedBase != 0 {
		t.Fatalf("RedBase = %v, want 0 — 100 comes from the other operation's row", got.RedBase)
	}
	if got.CodTrib == nil || *got.CodTrib != 0 {
		t.Fatalf("CodTrib = %v, want 0 — 60 comes from the other operation's row", got.CodTrib)
	}
}

// TestResolveCellSIsSentinelNotWildcard uses the real SP row "I/504 + S/-1".
// It is a rule for grupo 504. The original predicate accepted any row with "S"
// in either pair as a curinga, so this single row became a candidate for all
// 142 grupos in the catalog — the dominant source of the 881 ambiguous cells.
// Measured, "S" never carries a code other than -1 or 0.
func TestResolveCellSIsSentinelNotWildcard(t *testing.T) {
	aliquota := 12.0
	row504 := TgficmRow{
		TipRestricao: "I", CodRestricao: intp(504),
		TipRestricao2: "S", CodRestricao2: intp(-1),
		CodTrib: intp(10), Aliquota: &aliquota,
	}

	if got := ResolveCell([]TgficmRow{row504}, 122, topsOperacao); got.Found {
		t.Fatalf("grupo 122 resolved from a grupo-504 rule: %+v", got)
	}
	if got := ResolveCell([]TgficmRow{row504}, 42, topsOperacao); got.Found {
		t.Fatalf("grupo 42 resolved from a grupo-504 rule: %+v", got)
	}

	got := ResolveCell([]TgficmRow{row504}, 504, topsOperacao)
	if !got.Found || got.Ambiguo {
		t.Fatalf("got = %+v, want the rule to resolve its OWN grupo 504", got)
	}
	if got.Aliquota == nil || *got.Aliquota != 12.0 {
		t.Fatalf("Aliquota = %v, want 12.0", got.Aliquota)
	}
}

// TestResolveCellSpecificRuleOutranksMatrixDefault proves precedence. Every
// destination carries a matrix-wide default "S/-1 + S/-1" alongside the
// per-grupo rules; without precedence the two compete and every cell that has
// its own rule would report ambiguity.
func TestResolveCellSpecificRuleOutranksMatrixDefault(t *testing.T) {
	aliqEspecifica, aliqPadrao := 4.0, 7.0
	rows := []TgficmRow{
		{ // matrix-wide default for this destination
			TipRestricao: "S", CodRestricao: intp(-1),
			TipRestricao2: "S", CodRestricao2: intp(-1),
			CodTrib: intp(0), Aliquota: &aliqPadrao,
		},
		{ // rule for grupo 122
			TipRestricao: "N", CodRestricao: intp(0),
			TipRestricao2: "I", CodRestricao2: intp(122),
			CodTrib: intp(10), Aliquota: &aliqEspecifica,
		},
	}

	got := ResolveCell(rows, 122, topsOperacao)
	if !got.Found || got.Ambiguo {
		t.Fatalf("got = %+v, want Found=true Ambiguo=false (specific outranks default)", got)
	}
	if got.Aliquota == nil || *got.Aliquota != 4.0 {
		t.Fatalf("Aliquota = %v, want 4.0 from the specific rule", got.Aliquota)
	}

	// A grupo with no rule of its own legitimately falls back to the default.
	fallback := ResolveCell(rows, 999, topsOperacao)
	if !fallback.Found || fallback.Ambiguo {
		t.Fatalf("fallback = %+v, want the default rule to resolve grupo 999", fallback)
	}
	if fallback.Aliquota == nil || *fallback.Aliquota != 7.0 {
		t.Fatalf("fallback Aliquota = %v, want 7.0 from the default rule", fallback.Aliquota)
	}
}

// TestResolveCellAcceptsRowForOwnOperation is the positive control for the
// operation axis: the same shape that is discarded for CODTIPOPER 103 must be
// accepted when it names 306, the operation that actually issues our notas.
func TestResolveCellAcceptsRowForOwnOperation(t *testing.T) {
	aliquota := 7.0
	rows := []TgficmRow{
		{
			TipRestricao: "O", CodRestricao: intp(306),
			TipRestricao2: "I", CodRestricao2: intp(122),
			CodTrib: intp(10), Aliquota: &aliquota,
		},
	}

	got := ResolveCell(rows, 122, topsOperacao)
	if !got.Found || got.Ambiguo {
		t.Fatalf("got = %+v, want the rule for our own CODTIPOPER 306 to resolve", got)
	}
	if got.CodTrib == nil || *got.CodTrib != 10 {
		t.Fatalf("CodTrib = %v, want 10", got.CodTrib)
	}
}

// TestResolveCellDiscardsProductAndNCMRules covers pair 2 values the
// (uf_destino, grupo) grain cannot express: "P" names one CODPROD and "H"
// names one NCM. Applying either to the whole grupo would produce a plausible
// wrong number, so the row is discarded and the cell stays unknown (ADR-17).
func TestResolveCellDiscardsProductAndNCMRules(t *testing.T) {
	rows := []TgficmRow{
		{ // real MG shape: TOP 59 + produto 18516
			TipRestricao: "O", CodRestricao: intp(59),
			TipRestricao2: "P", CodRestricao2: intp(18516),
			CodTrib: intp(0),
		},
		{ // real ES shape: TOP 115 + NCM 69072300
			TipRestricao: "O", CodRestricao: intp(115),
			TipRestricao2: "H", CodRestricao2: intp(69072300),
			CodTrib: intp(60),
		},
	}

	if got := ResolveCell(rows, 122, topsOperacao); got.Found {
		t.Fatalf("Found = true (%+v), want false — product/NCM rules are not expressible at grupo grain", got)
	}
}

// TestResolveCellAmbiguityOnlyAmongEqualPrecedence proves ambiguity survives
// where it is real: two rules that both name grupo 122 for the same
// destination. Nothing is guessed — codtrib wrong is 0% vs 24% ICMS.
func TestResolveCellAmbiguityOnlyAmongEqualPrecedence(t *testing.T) {
	a, b := 4.0, 12.0
	rows := []TgficmRow{
		{TipRestricao: "N", CodRestricao: intp(0), TipRestricao2: "I", CodRestricao2: intp(122), CodTrib: intp(0), Aliquota: &a},
		{TipRestricao: "I", CodRestricao: intp(122), TipRestricao2: "S", CodRestricao2: intp(-1), CodTrib: intp(10), Aliquota: &b},
	}

	got := ResolveCell(rows, 122, topsOperacao)
	if !got.Found || !got.Ambiguo {
		t.Fatalf("got = %+v, want Found=true Ambiguo=true", got)
	}
	if got.LinhasCandidatas != 2 {
		t.Fatalf("LinhasCandidatas = %d, want 2", got.LinhasCandidatas)
	}
	if got.CodTrib != nil || got.Aliquota != nil || got.RedBase != nil || got.PercFcp != nil {
		t.Fatalf("ambiguous cell carries a chosen fiscal value = %+v, want all nil", got)
	}
}

// TestResolveCellCollapsesDuplicatedDefaultRule uses the real GO pair: the
// per-destination default is stored twice, "N/0 + S/-1" and "S/-1 + S/-1",
// agreeing on aliquota and red_base with only one of the copies carrying
// CODTRIB. Reporting that as ambiguity would withhold a rate the ERP states
// unambiguously — it made 51 DF cells, 36 GO cells and 42 SP cells
// non-calculable on the first live run of the corrected predicate.
func TestResolveCellCollapsesDuplicatedDefaultRule(t *testing.T) {
	aliqA, aliqB, redA, redB := 7.0, 7.0, 0.0, 0.0
	rows := []TgficmRow{
		{TipRestricao: "N", CodRestricao: intp(0), TipRestricao2: "S", CodRestricao2: intp(-1),
			CodTrib: nil, Aliquota: &aliqA, RedBase: &redA},
		{TipRestricao: "S", CodRestricao: intp(-1), TipRestricao2: "S", CodRestricao2: intp(-1),
			CodTrib: intp(0), Aliquota: &aliqB, RedBase: &redB},
	}

	got := ResolveCell(rows, 507, topsOperacao)
	if !got.Found || got.Ambiguo {
		t.Fatalf("got = %+v, want Found=true Ambiguo=false — the two rows never contradict each other", got)
	}
	if got.LinhasCandidatas != 2 {
		t.Errorf("LinhasCandidatas = %d, want 2 — the reconciliation snapshot must still show both rows", got.LinhasCandidatas)
	}
	if got.CodTrib == nil || *got.CodTrib != 0 {
		t.Errorf("CodTrib = %v, want 0 — the only value either row states", got.CodTrib)
	}
	if got.Aliquota == nil || *got.Aliquota != 7.0 {
		t.Errorf("Aliquota = %v, want 7.0", got.Aliquota)
	}
}

// TestResolveCellDoesNotCollapseContradictingRules is the negative control for
// the merge: same shape, but the two rows state DIFFERENT CODTRIB. That is
// real disagreement and must stay ambiguous — codtrib wrong is 0% vs 24% ICMS.
func TestResolveCellDoesNotCollapseContradictingRules(t *testing.T) {
	aliq := 7.0
	rows := []TgficmRow{
		{TipRestricao: "N", CodRestricao: intp(0), TipRestricao2: "S", CodRestricao2: intp(-1),
			CodTrib: intp(60), Aliquota: &aliq},
		{TipRestricao: "S", CodRestricao: intp(-1), TipRestricao2: "S", CodRestricao2: intp(-1),
			CodTrib: intp(0), Aliquota: &aliq},
	}

	got := ResolveCell(rows, 507, topsOperacao)
	if !got.Found || !got.Ambiguo {
		t.Fatalf("got = %+v, want Ambiguo=true — CODTRIB 60 and 0 are different claims", got)
	}
	if got.CodTrib != nil || got.Aliquota != nil {
		t.Errorf("ambiguous cell carries a chosen value = %+v, want all nil", got)
	}
}

// TestResolveCellNoSurvivorWritesNoCell proves the "nenhuma sobrevive" branch:
// the caller must write no row at all — not a zero/placeholder row (ADR-17).
// "L" and "X" restrict on a dimension this decoding has not identified, so
// they can never be claimed to apply.
func TestResolveCellNoSurvivorWritesNoCell(t *testing.T) {
	rows := []TgficmRow{
		{TipRestricao: "L", CodRestricao: intp(-1), TipRestricao2: "I", CodRestricao2: intp(206)},
		{TipRestricao: "X", CodRestricao: intp(101), TipRestricao2: "I", CodRestricao2: intp(122)},
		{TipRestricao: "N", CodRestricao: intp(0), TipRestricao2: "I", CodRestricao2: intp(3)},
	}

	got := ResolveCell(rows, 7, topsOperacao)
	if got.Found {
		t.Fatalf("Found = true (%+v), want false — no row applies to grupo 7", got)
	}
}

// TestResolveCellEmptyTOPSetDiscardsOperationRules guards the adapter contract:
// passing no operation codes must not silently widen the predicate into
// accepting every operation's rules.
func TestResolveCellEmptyTOPSetDiscardsOperationRules(t *testing.T) {
	rows := []TgficmRow{
		{TipRestricao: "O", CodRestricao: intp(306), TipRestricao2: "I", CodRestricao2: intp(122), CodTrib: intp(10)},
	}

	if got := ResolveCell(rows, 122, nil); got.Found {
		t.Fatalf("Found = true (%+v), want false when no operation is declared", got)
	}
}
