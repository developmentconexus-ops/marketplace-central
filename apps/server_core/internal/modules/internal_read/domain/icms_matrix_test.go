package domain

import "testing"

func intp(v int) *int { return &v }

// TestResolveCellDiscardsOutsideWhitelistAndFlagsAmbiguity is the P2.b Task 3
// must-pass scenario from the dispatch: three candidate TGFICM rows for the
// same destination — N/0 (unrestricted), I/7 (restricted to grupo 7), and O/…
// (not a whitelist member at all). For grupo=7, O must be discarded outright
// and the remaining two (N/0 and I/7) both survive the whitelist, so the cell
// is ambiguous and CodTrib is never chosen — codtrib wrong is the difference
// between 0% and 24% ICMS, so an ambiguous cell must never guess.
func TestResolveCellDiscardsOutsideWhitelistAndFlagsAmbiguity(t *testing.T) {
	rows := []TgficmRow{
		{ // N/0 — unrestricted, both fields "N".
			TipRestricao: "N", TipRestricao2: "N",
			CodTrib: intp(0),
		},
		{ // I/7 — restricted to grupo 7.
			TipRestricao: "I", CodRestricao: intp(7),
			CodTrib: intp(10),
		},
		{ // O/… — not a whitelist member (also not I/N/S on the second field).
			TipRestricao: "O", TipRestricao2: "X",
			CodTrib: intp(99),
		},
	}

	got := ResolveCell(rows, 7)

	if !got.Found {
		t.Fatalf("Found = false, want true (two candidates survive)")
	}
	if !got.Ambiguo {
		t.Fatalf("Ambiguo = false, want true (N/0 and I/7 both survive)")
	}
	if got.LinhasCandidatas != 2 {
		t.Fatalf("LinhasCandidatas = %d, want 2 (O must be discarded)", got.LinhasCandidatas)
	}
	if got.CodTrib != nil {
		t.Fatalf("CodTrib = %v, want nil — an ambiguous cell must never choose a CodTrib", *got.CodTrib)
	}
	if got.Aliquota != nil || got.RedBase != nil || got.PercFcp != nil {
		t.Fatalf("ambiguous cell carries a chosen fiscal value = %+v, want all nil", got)
	}
}

// TestResolveCellSingleSurvivorWritesTheCell proves the non-ambiguous path:
// exactly one whitelist-acceptable row survives, so the cell is written with
// its CodTrib and reconciliation fields, Ambiguo=false.
func TestResolveCellSingleSurvivorWritesTheCell(t *testing.T) {
	aliquota, redBase, percFcp := 12.0, 0.0, 2.0
	rows := []TgficmRow{
		{TipRestricao: "I", CodRestricao: intp(7), CodTrib: intp(10), Aliquota: &aliquota, RedBase: &redBase, PercFcp: &percFcp},
		{TipRestricao: "I", CodRestricao: intp(3), CodTrib: intp(20)}, // different grupo, must not match
		{TipRestricao: "O"}, // not whitelisted
	}

	got := ResolveCell(rows, 7)

	if !got.Found || got.Ambiguo {
		t.Fatalf("got = %+v, want Found=true Ambiguo=false", got)
	}
	if got.LinhasCandidatas != 1 {
		t.Fatalf("LinhasCandidatas = %d, want 1", got.LinhasCandidatas)
	}
	if got.CodTrib == nil || *got.CodTrib != 10 {
		t.Fatalf("CodTrib = %v, want 10", got.CodTrib)
	}
	if got.Aliquota == nil || *got.Aliquota != 12.0 {
		t.Fatalf("Aliquota = %v, want 12.0", got.Aliquota)
	}
}

// TestResolveCellNoSurvivorWritesNoCell proves the "nenhuma sobrevive" branch:
// the caller must write no row at all — not a zero/placeholder row.
func TestResolveCellNoSurvivorWritesNoCell(t *testing.T) {
	rows := []TgficmRow{
		{TipRestricao: "O"},
		{TipRestricao: "P", TipRestricao2: "H"},
		{TipRestricao: "I", CodRestricao: intp(3)}, // different grupo
	}

	got := ResolveCell(rows, 7)

	if got.Found {
		t.Fatalf("Found = true, want false (no candidate survives the whitelist for grupo 7)")
	}
}

// TestResolveCellWildcardSurvivesRegardlessOfCode proves the "S" branch
// ignores CodRestricao/CodRestricao2 entirely — it is a curinga.
func TestResolveCellWildcardSurvivesRegardlessOfCode(t *testing.T) {
	rows := []TgficmRow{
		{TipRestricao: "S", CodRestricao: intp(-1), CodTrib: intp(0)},
	}

	got := ResolveCell(rows, 42)

	if !got.Found || got.Ambiguo || got.CodTrib == nil || *got.CodTrib != 0 {
		t.Fatalf("got = %+v, want single wildcard survivor with CodTrib=0", got)
	}
}
