package domain

// TgficmRow is one raw METALPRD.TGFICM row for the company's fixed origin
// (UFORIG=13/MG, D-17 territory — not a per-tenant choice), already joined to
// its destination UF's text sigla by the Oracle adapter (icms_matrix.go) so
// ResolveCell never has to touch Sankhya's internal numeric UF codes.
//
// The two restriction pairs are NOT two interchangeable slots. Measured
// against all 4209 rows of the live matrix for UFORIG=13:
//
//	pair 1 (TipRestricao/CodRestricao) is the OPERATION axis
//	  "N" 3815 rows, code always 0                  -> no operation restriction
//	  "O"  240 rows, 238/240 match TGFTOP           -> one CODTIPOPER only
//	  "I"  125 rows, 123/125 match TGFPRO.GRUPOICMS -> one grupo only
//	  "S"    6 rows, code only -1 or 0              -> axis unused
//	  "L"/"X"/"T" 23 rows                           -> not decoded, discarded
//
//	pair 2 (TipRestricao2/CodRestricao2) is the PRODUCT axis
//	  "I" 4054 rows, matches GRUPOICMS              -> one grupo only
//	  "S"  137 rows, code only -1 or 0              -> axis unused
//	  "P"   14 rows, 14/14 match TGFPRO.CODPROD     -> one produto only
//	  "H"/"K"/"G" 4 rows (NCM and others)           -> not expressible at this grain
//
// CodRestricao/CodRestricao2 are nil when Oracle's value was NULL (ADR-17,
// honest-unknown).
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
	// Found is false when no row in the input survives. The caller must write
	// NO row at all for this cell — never a zero-value or placeholder row
	// (ADR-17: unknown never becomes 0).
	Found bool
	// Ambiguo is true when more than one row survives at the same precedence
	// tier. CodTrib incorrect is the difference between 0% and 24% ICMS, so an
	// ambiguous cell leaves CodTrib/Aliquota/RedBase/PercFcp all nil rather
	// than picking one candidate.
	Ambiguo bool
	// LinhasCandidatas is the number of rows that survived at the winning
	// precedence tier (1 when Ambiguo is false, >1 when Ambiguo is true).
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

// axisVerdict is how one restriction pair of one TGFICM row relates to the
// cell being resolved.
type axisVerdict int

const (
	// axisDiscard: the pair restricts to something this cell is not. The whole
	// row is out — a row is a candidate only when BOTH pairs admit it.
	axisDiscard axisVerdict = iota
	// axisNeutral: the pair is switched off (sentinel code), so it restricts
	// nothing and the row stays in as a generic rule.
	axisNeutral
	// axisSpecific: the pair names exactly this cell's grupo (or the company's
	// own operation), so the row is a specific rule and outranks generic ones.
	axisSpecific
)

// ResolveCell is pure: it decides which TGFICM rows for ONE uf_destino apply to
// one (grupo, operation) and resolves the cell from them.
//
// Both restriction pairs are evaluated — a row survives only when NEITHER pair
// excludes it. That is the correction for the original whitelist predicate,
// which read one pair at a time and let through rows the other pair was
// restricting: a row like O/103 + I/122 (a rule belonging to CODTIPOPER 103)
// was accepted for grupo 122 under every operation, and a row like I/504 + S/-1
// (a rule belonging to grupo 504) was accepted for every grupo because "S" was
// read as a wildcard. Measured, "S" never carries a real code — only -1 or 0 —
// so it is the sentinel for "this axis is unused", not a wildcard.
//
// topsOperacao is the set of CODTIPOPER the caller's own operation covers; a
// row restricted to any other operation is discarded. An empty set discards
// every operation-restricted row.
//
// Precedence: rows that name this cell specifically win over rows that restrict
// nothing. Ambiguo is reported only among rows at the winning tier, so a
// specific rule coexisting with the matrix-wide default is not ambiguity.
func ResolveCell(rows []TgficmRow, grupo int, topsOperacao map[int]bool) ICMSMatrixCellResolution {
	var especificas, genericas []TgficmRow
	for _, row := range rows {
		operacao := classifyOperacao(row, grupo, topsOperacao)
		if operacao == axisDiscard {
			continue
		}
		produto := classifyProduto(row, grupo)
		if produto == axisDiscard {
			continue
		}
		if operacao == axisSpecific || produto == axisSpecific {
			especificas = append(especificas, row)
			continue
		}
		genericas = append(genericas, row)
	}

	candidatas := especificas
	if len(candidatas) == 0 {
		candidatas = genericas
	}

	if len(candidatas) == 0 {
		return ICMSMatrixCellResolution{Found: false}
	}
	merged, compativel := mergeCompativel(candidatas)
	if !compativel {
		return ICMSMatrixCellResolution{
			Found:            true,
			Ambiguo:          true,
			LinhasCandidatas: len(candidatas),
		}
	}
	return ICMSMatrixCellResolution{
		Found:            true,
		Ambiguo:          false,
		LinhasCandidatas: len(candidatas),
		CodTrib:          merged.CodTrib,
		Aliquota:         merged.Aliquota,
		RedBase:          merged.RedBase,
		PercFcp:          merged.PercFcp,
	}
}

// mergeCompativel collapses candidate rows that never contradict each other
// into a single row, and reports false as soon as two of them state DIFFERENT
// non-nil values for the same field.
//
// This exists because the live matrix duplicates its per-destination default
// rule: DF, GO and SP each carry both "N/0 + S/-1" and "S/-1 + S/-1" with the
// same aliquota, red_base and fcp, one of the two copies simply leaving CODTRIB
// unset. Reporting those as ambiguous would withhold a rate the ERP states
// unambiguously. Merging invents nothing — the result carries only values the
// rows themselves assert, and a row that is silent on a field never overrides
// one that is not (a nil is missing information, not a competing claim).
//
// Genuine disagreement — two rules naming the same cell with different
// aliquotas or different CODTRIB — still resolves to Ambiguo, which is the
// whole point: codtrib wrong is 0% vs 24% ICMS.
func mergeCompativel(rows []TgficmRow) (TgficmRow, bool) {
	var merged TgficmRow
	for _, row := range rows {
		if !mergeInt(&merged.CodTrib, row.CodTrib) {
			return TgficmRow{}, false
		}
		if !mergeFloat(&merged.Aliquota, row.Aliquota) {
			return TgficmRow{}, false
		}
		if !mergeFloat(&merged.RedBase, row.RedBase) {
			return TgficmRow{}, false
		}
		if !mergeFloat(&merged.PercFcp, row.PercFcp) {
			return TgficmRow{}, false
		}
	}
	return merged, true
}

func mergeInt(dst **int, src *int) bool {
	if src == nil {
		return true
	}
	if *dst == nil {
		*dst = src
		return true
	}
	return **dst == *src
}

func mergeFloat(dst **float64, src *float64) bool {
	if src == nil {
		return true
	}
	if *dst == nil {
		*dst = src
		return true
	}
	return **dst == *src
}

// classifyOperacao reads pair 1, the operation axis.
func classifyOperacao(row TgficmRow, grupo int, topsOperacao map[int]bool) axisVerdict {
	switch row.TipRestricao {
	case "N":
		// Measured: "N" carries code 0 in all 3815 rows. Any other shape is one
		// this decoding has never seen, so it is not assumed harmless.
		if isSentinelCode(row.CodRestricao) {
			return axisNeutral
		}
		return axisDiscard
	case "S":
		if isSentinelCode(row.CodRestricao) {
			return axisNeutral
		}
		return axisDiscard
	case "O":
		if row.CodRestricao != nil && topsOperacao[*row.CodRestricao] {
			return axisSpecific
		}
		return axisDiscard
	case "I":
		if row.CodRestricao != nil && *row.CodRestricao == grupo {
			return axisSpecific
		}
		return axisDiscard
	}
	// "L", "X", "T" and anything else: the axis restricts on a dimension this
	// decoding cannot evaluate, so the row cannot be claimed to apply.
	return axisDiscard
}

// classifyProduto reads pair 2, the product axis.
func classifyProduto(row TgficmRow, grupo int) axisVerdict {
	switch row.TipRestricao2 {
	case "S":
		if isSentinelCode(row.CodRestricao2) {
			return axisNeutral
		}
		return axisDiscard
	case "I":
		if row.CodRestricao2 != nil && *row.CodRestricao2 == grupo {
			return axisSpecific
		}
		return axisDiscard
	}
	// "P" (one CODPROD), "H" (one NCM), "K", "G": real restrictions the
	// (uf_destino, grupo) grain cannot express. Discarded rather than applied
	// to every product of the grupo — the cell stays unknown instead of
	// becoming plausibly wrong.
	return axisDiscard
}

// isSentinelCode reports whether a restriction code means "this axis is
// switched off". Measured across the live matrix the off values are NULL, 0
// and -1; every real restriction carries a positive code.
func isSentinelCode(code *int) bool {
	return code == nil || *code <= 0
}
