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
	// whitelist is applied per-cell in Go via domain.ResolveCell.
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

// TestICMSMatrixReaderResolveCellsAppliesWhitelistPerUFAndGrupo drives
// ResolveCells over a canned Oracle snapshot spanning two destinations and
// proves: an unrestricted (N/0) row alone resolves cleanly; the SAME
// unrestricted row alongside a grupo-restricted (I) row for the same grupo is
// ambiguous (CodTrib withheld); a destination whose only row is outside the
// whitelist ("O") produces no cell for any grupo; a row whose TSIUFS join
// could not resolve a sigla (UF NULL) is dropped rather than panicking.
func TestICMSMatrixReaderResolveCellsAppliesWhitelistPerUFAndGrupo(t *testing.T) {
	q := &dispatchQueryer{results: map[string]fakeResult{
		"TGFICM": {cols: 9, rows: [][]driver.Value{
			// SP: N/0 unrestricted, codtrib=0.
			{"SP", "N", nil, "N", nil, int64(0), 12.0, 0.0, 0.0},
			// SP: I/7 restricted to grupo 7, codtrib=10 — collides with the N/0
			// row above for grupo=7 only.
			{"SP", "I", int64(7), nil, nil, int64(10), 18.0, 0.0, 0.0},
			// RJ: outside the whitelist entirely.
			{"RJ", "O", nil, "P", nil, int64(99), 20.0, 0.0, 0.0},
			// UFDEST whose TSIUFS join produced no sigla — must be dropped, not panic.
			{nil, "N", nil, "N", nil, int64(1), 10.0, 0.0, 0.0},
		}},
		"TGFPRO": {cols: 1, rows: [][]driver.Value{
			{int64(3)},
			{int64(7)},
		}},
	}}

	cells, err := NewICMSMatrixReader(q).ResolveCells(context.Background())
	if err != nil {
		t.Fatalf("ResolveCells: %v", err)
	}

	byKey := make(map[[2]any]icmsMatrixCellTestView, len(cells))
	for _, c := range cells {
		byKey[[2]any{c.UFDestino, c.GrupoICMS}] = icmsMatrixCellTestView{
			UFOrigem: c.UFOrigem, Ambiguo: c.Ambiguo, Linhas: c.LinhasCandidatas, CodTrib: c.CodTrib,
		}
	}

	// SP/grupo=3: only the unrestricted N/0 row applies (I/7 does not match
	// grupo 3) → single survivor, codtrib=0.
	spThree, ok := byKey[[2]any{"SP", 3}]
	if !ok {
		t.Fatalf("missing cell SP/3, cells=%+v", cells)
	}
	if spThree.UFOrigem != "MG" || spThree.Ambiguo || spThree.Linhas != 1 || spThree.CodTrib == nil || *spThree.CodTrib != 0 {
		t.Errorf("SP/3 = %+v, want UFOrigem=MG Ambiguo=false Linhas=1 CodTrib=0", spThree)
	}

	// SP/grupo=7: both N/0 and I/7 apply → ambiguous, CodTrib withheld.
	spSeven, ok := byKey[[2]any{"SP", 7}]
	if !ok {
		t.Fatalf("missing cell SP/7, cells=%+v", cells)
	}
	if !spSeven.Ambiguo || spSeven.Linhas != 2 || spSeven.CodTrib != nil {
		t.Errorf("SP/7 = %+v, want Ambiguo=true Linhas=2 CodTrib=nil", spSeven)
	}

	// RJ: the only row is "O" (not whitelisted) — no cell for RJ at any grupo.
	if _, ok := byKey[[2]any{"RJ", 3}]; ok {
		t.Errorf("RJ/3 unexpectedly resolved, cells=%+v", cells)
	}
	if _, ok := byKey[[2]any{"RJ", 7}]; ok {
		t.Errorf("RJ/7 unexpectedly resolved, cells=%+v", cells)
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
}
