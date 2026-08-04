package oracle

import (
	"context"
	"database/sql/driver"
	"strings"
	"testing"
)

func TestICMSMatrixCandidatesSQLJoinsTSIUFSWithCodPais55AndFixedOrigin(t *testing.T) {
	upper := strings.ToUpper(icmsMatrixCandidatesSQL)
	for _, want := range []string{
		"METALPRD.TGFICM", "METALPRD.TSIUFS", "CODPAIS = 55", "UFORIG = :1",
		"TIPRESTRICAO", "CODRESTRICAO", "TIPRESTRICAO2", "CODRESTRICAO2",
		"CODTRIB", "ALIQUOTA", "REDBASE", "PERCICMSFCP",
	} {
		if !strings.Contains(upper, want) {
			t.Errorf("icmsMatrixCandidatesSQL missing %q\n%s", want, icmsMatrixCandidatesSQL)
		}
	}
	// ALIQUFDEST is the ST-block ceiling rate (icms_ceiling.go, a different,
	// older feature) — ~2.5x the operation's own rate. Selecting it here would
	// silently poison icms_matrix_mirror.aliquota with a plausible wrong number.
	if strings.Contains(upper, "ALIQUFDEST") {
		t.Errorf("icmsMatrixCandidatesSQL must not select ALIQUFDEST (wrong rate, ST block, not the operation)\n%s", icmsMatrixCandidatesSQL)
	}
	// No grupo/uf_destino filter in SQL — architecture decision (a): the
	// predicate is applied per-cell in Go via domain.ResolveCell.
	if strings.Contains(upper, "GRUPOICMS") {
		t.Errorf("icmsMatrixCandidatesSQL must not filter by grupo in SQL (decision (a))\n%s", icmsMatrixCandidatesSQL)
	}
}

func TestICMSKnownGruposSQLFiltersNotNull(t *testing.T) {
	upper := strings.ToUpper(icmsKnownGruposSQL)
	for _, want := range []string{"DISTINCT GRUPOICMS", "METALPRD.TGFPRO", "GRUPOICMS IS NOT NULL"} {
		if !strings.Contains(upper, want) {
			t.Errorf("icmsKnownGruposSQL missing %q\n%s", want, icmsKnownGruposSQL)
		}
	}
}

// TestICMSMatrixTOPsAreTheMarketplaceOperation pins the operation the mirror
// resolves for: 313 (pedido) and 306 (NFE ENTREGA E-COMMERCE, the nota it
// becomes). A row restricted to any other CODTIPOPER belongs to a different
// operation and must never resolve this operation's cells.
func TestICMSMatrixTOPsAreTheMarketplaceOperation(t *testing.T) {
	if len(icmsMatrixTOPs) != 2 || !icmsMatrixTOPs[306] || !icmsMatrixTOPs[313] {
		t.Fatalf("icmsMatrixTOPs = %v, want exactly {306, 313}", icmsMatrixTOPs)
	}
}

// TestICMSMatrixReaderResolveCellsAppliesPredicatePerUFAndGrupo drives
// ResolveCells over a canned snapshot built from row shapes actually measured
// in METALPRD.TGFICM, and proves the wiring end to end:
//
//   - a rule restricted to another operation (O/103 + I/122) does not resolve
//     grupo 122 — the defect that turned 881 of 2700 live cells ambiguous;
//   - a rule for grupo 504 carrying the "S" sentinel does not leak into other
//     grupos — "S" is "axis unused", not a wildcard;
//   - the matrix-wide default (S/-1 + S/-1) resolves grupos that have no rule
//     of their own, and loses to a specific rule where one exists;
//   - a destination whose only row belongs to another operation produces no
//     cell at any grupo (ADR-17: no placeholder row);
//   - a row whose TSIUFS join produced no sigla is dropped rather than panicking.
func TestICMSMatrixReaderResolveCellsAppliesPredicatePerUFAndGrupo(t *testing.T) {
	q := &dispatchQueryer{results: map[string]fakeResult{
		"TGFICM": {cols: 9, rows: [][]driver.Value{
			// SP: our operation's rule for grupo 122, codtrib=0.
			{"SP", "N", int64(0), "I", int64(122), int64(0), 12.0, 0.0, 0.0},
			// SP: same grupo 122, but restricted to CODTIPOPER 103 — a different
			// operation, with a reduced base that must not reach our cell.
			{"SP", "O", int64(103), "I", int64(122), int64(60), 7.0, 100.0, 0.0},
			// SP: rule for grupo 504, product axis switched off.
			{"SP", "I", int64(504), "S", int64(-1), int64(10), 18.0, 0.0, 0.0},
			// SP: the matrix-wide default for this destination.
			{"SP", "S", int64(-1), "S", int64(-1), int64(20), 4.0, 0.0, 0.0},
			// RJ: only a rule for another operation — no cell may be written.
			{"RJ", "O", int64(103), "I", int64(122), int64(99), 20.0, 0.0, 0.0},
			// UFDEST whose TSIUFS join produced no sigla — dropped, not a panic.
			{nil, "N", int64(0), "I", int64(122), int64(1), 10.0, 0.0, 0.0},
		}},
		"TGFPRO": {cols: 1, rows: [][]driver.Value{
			{int64(3)},
			{int64(122)},
			{int64(504)},
		}},
	}}

	cells, err := NewICMSMatrixReader(q).ResolveCells(context.Background())
	if err != nil {
		t.Fatalf("ResolveCells: %v", err)
	}

	byKey := make(map[[2]any]icmsMatrixCellTestView, len(cells))
	for _, c := range cells {
		byKey[[2]any{c.UFDestino, c.GrupoICMS}] = icmsMatrixCellTestView{
			UFOrigem: c.UFOrigem, Ambiguo: c.Ambiguo, Linhas: c.LinhasCandidatas,
			CodTrib: c.CodTrib, RedBase: c.RedBase,
		}
	}

	// SP/122: only our operation's N/0 + I/122 rule applies. The O/103 row is
	// another operation's and must not make the cell ambiguous nor contribute
	// its RedBase=100.
	spCentoEVinteDois, ok := byKey[[2]any{"SP", 122}]
	if !ok {
		t.Fatalf("missing cell SP/122, cells=%+v", cells)
	}
	if spCentoEVinteDois.UFOrigem != "MG" || spCentoEVinteDois.Ambiguo || spCentoEVinteDois.Linhas != 1 {
		t.Errorf("SP/122 = %+v, want UFOrigem=MG Ambiguo=false Linhas=1", spCentoEVinteDois)
	}
	if spCentoEVinteDois.CodTrib == nil || *spCentoEVinteDois.CodTrib != 0 {
		t.Errorf("SP/122 CodTrib = %v, want 0 — 60 belongs to CODTIPOPER 103", spCentoEVinteDois.CodTrib)
	}
	if spCentoEVinteDois.RedBase == nil || *spCentoEVinteDois.RedBase != 0 {
		t.Errorf("SP/122 RedBase = %v, want 0 — 100 belongs to CODTIPOPER 103", spCentoEVinteDois.RedBase)
	}

	// SP/504: resolved by its own rule, not by the default.
	spQuinhentosEQuatro, ok := byKey[[2]any{"SP", 504}]
	if !ok {
		t.Fatalf("missing cell SP/504, cells=%+v", cells)
	}
	if spQuinhentosEQuatro.Ambiguo || spQuinhentosEQuatro.CodTrib == nil || *spQuinhentosEQuatro.CodTrib != 10 {
		t.Errorf("SP/504 = %+v, want Ambiguo=false CodTrib=10", spQuinhentosEQuatro)
	}

	// SP/3 has no rule of its own — the matrix-wide default resolves it, and
	// the grupo-504 rule must not have leaked in via its "S" sentinel.
	spTres, ok := byKey[[2]any{"SP", 3}]
	if !ok {
		t.Fatalf("missing cell SP/3, cells=%+v", cells)
	}
	if spTres.Ambiguo || spTres.Linhas != 1 || spTres.CodTrib == nil || *spTres.CodTrib != 20 {
		t.Errorf("SP/3 = %+v, want Ambiguo=false Linhas=1 CodTrib=20 (the default rule)", spTres)
	}

	// RJ: the only row belongs to CODTIPOPER 103 — no cell at any grupo.
	for _, grupo := range []int{3, 122, 504} {
		if _, ok := byKey[[2]any{"RJ", grupo}]; ok {
			t.Errorf("RJ/%d unexpectedly resolved from another operation's rule, cells=%+v", grupo, cells)
		}
	}

	// The UF-NULL row must not have manufactured a destination.
	for _, c := range cells {
		if c.UFDestino == "" {
			t.Errorf("cell with empty UFDestino leaked through: %+v", c)
		}
	}
}

type icmsMatrixCellTestView struct {
	UFOrigem string
	Ambiguo  bool
	Linhas   int
	CodTrib  *int
	RedBase  *float64
}
