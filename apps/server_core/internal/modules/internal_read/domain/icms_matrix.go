package domain

// TgficmRow is one raw METALPRD.TGFICM row for the company's fixed origin
// (UFORIG=13/MG, D-17 territory — not a per-tenant choice), already joined to
// its destination UF's text sigla by the Oracle adapter (icms_matrix.go) so
// ResolveCell never has to touch Sankhya's internal numeric UF codes. Fields
// carry only what the P2.b Task 3 whitelist predicate and the reconciliation
// snapshot need.
//
// TipRestricao/TipRestricao2 are trimmed Sankhya restriction-type codes — "N",
// "I", "S" are whitelist members; "O", "P", "H", "K", "G", "L", "X", "T" are
// all real TGFICM values and none of them are. CodRestricao/CodRestricao2 are
// nil when Oracle's value was NULL (ADR-17, honest-unknown): a nil pointer can
// never equal a grupo int, so an unset code simply never satisfies the "I"
// branch below — it does not need special-casing.
type TgficmRow struct {
	TipRestricao  string
	CodRestricao  *int
	TipRestricao2 string
	CodRestricao2 *int
	CodTrib       *int
	Aliquota      *float64
	RedBase       *float64
	PercFcp       *float64
}

// ICMSMatrixCellResolution is what ResolveCell decides for one (uf_destino,
// grupo) pair, given the candidate TGFICM rows already scoped to that
// destination.
type ICMSMatrixCellResolution struct {
	// Found is false when no row in the input survives the whitelist. The
	// caller must write NO row at all for this cell — never a zero-value or
	// placeholder row (ADR-17: unknown never becomes 0).
	Found bool
	// Ambiguo is true when more than one row survives the whitelist. The
	// precedence between competing TGFICM rows has never been established —
	// only one dry-run observation exists, insufficient to generalize into a
	// hierarchy — so when Ambiguo is true, CodTrib/Aliquota/RedBase/PercFcp are
	// all left nil rather than guessed. CodTrib incorrect is the difference
	// between 0% and 24% ICMS: an ambiguous cell must never silently pick one
	// candidate's CodTrib.
	Ambiguo bool
	// LinhasCandidatas is the number of rows that survived the whitelist (1
	// when Ambiguo is false, >1 when Ambiguo is true).
	LinhasCandidatas int
	CodTrib          *int
	Aliquota         *float64
	RedBase          *float64
	PercFcp          *float64
}

// ICMSMatrixCell is one fully-identified, resolved cell ready for the
// versioned writer: ResolveCell's output plus the (uf_origem, uf_destino,
// grupo_icms) key the caller already knows when it invokes ResolveCell.
type ICMSMatrixCell struct {
	UFOrigem         string
	UFDestino        string
	GrupoICMS        int
	CodTrib          *int
	Aliquota         *float64
	RedBase          *float64
	PercFcp          *float64
	LinhasCandidatas int
	Ambiguo          bool
}

// ResolveCell is pure: it applies the ICMS-ST whitelist ratified in the P2.b
// Task 3 dispatch ("Resolução por lista branca") to the raw TGFICM candidate
// rows for ONE uf_destino, deciding which rows are acceptable for grupo. A row
// is acceptable if:
//
//   - TipRestricao2 == "I" and CodRestricao2 == grupo, OR
//   - TipRestricao  == "I" and CodRestricao  == grupo, OR
//   - TipRestricao == "N" AND TipRestricao2 == "N" (unrestricted — this branch
//     requires BOTH fields to be "N", not either alone), OR
//   - TipRestricao == "S" OR TipRestricao2 == "S" (curinga/wildcard)
//
// Anything else — "O", "P", "H", "K", "G", "L", "X", "T", or any row that
// matches none of the above — is discarded before candidates are even
// counted. Precedence between competing acceptable rows is deliberately NOT
// decided here: that is exactly what Ambiguo communicates to the caller.
func ResolveCell(rows []TgficmRow, grupo int) ICMSMatrixCellResolution {
	var candidates []TgficmRow
	for _, row := range rows {
		if isAcceptableCandidate(row, grupo) {
			candidates = append(candidates, row)
		}
	}

	switch len(candidates) {
	case 0:
		return ICMSMatrixCellResolution{Found: false}
	case 1:
		c := candidates[0]
		return ICMSMatrixCellResolution{
			Found:            true,
			Ambiguo:          false,
			LinhasCandidatas: 1,
			CodTrib:          c.CodTrib,
			Aliquota:         c.Aliquota,
			RedBase:          c.RedBase,
			PercFcp:          c.PercFcp,
		}
	default:
		return ICMSMatrixCellResolution{
			Found:            true,
			Ambiguo:          true,
			LinhasCandidatas: len(candidates),
		}
	}
}

func isAcceptableCandidate(row TgficmRow, grupo int) bool {
	if row.TipRestricao2 == "I" && row.CodRestricao2 != nil && *row.CodRestricao2 == grupo {
		return true
	}
	if row.TipRestricao == "I" && row.CodRestricao != nil && *row.CodRestricao == grupo {
		return true
	}
	if row.TipRestricao == "N" && row.TipRestricao2 == "N" {
		return true
	}
	if row.TipRestricao == "S" || row.TipRestricao2 == "S" {
		return true
	}
	return false
}
