package oracle

import (
	"context"
	"database/sql"
	"sort"
	"strings"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
)

// icmsMatrixOriginUF is the Sankhya-internal numeric origin state code the
// client's single empresa always ships from (13 = MG). Fixed the same way
// sankhyaCompany (sync.go) is fixed: D-17 territory, not a decision this task
// reopens.
const icmsMatrixOriginUF = 13

// icmsMatrixOriginUFSigla is the text sigla icms_matrix_mirror.uf_origem
// always carries, given icmsMatrixOriginUF.
const icmsMatrixOriginUFSigla = "MG"

// icmsMatrixTOPs are the Sankhya CODTIPOPER values the marketplace operation
// runs under: 313 is the pedido, 306 ("NFE ENTREGA E-COMMERCE") is the nota it
// becomes when invoiced. TGFICM rows whose operation axis (TIPRESTRICAO='O')
// names any other CODTIPOPER belong to a different operation and must not
// resolve this operation's cells — that is the correction measured in the
// P2.b Task 7 follow-up, where 240 rows restricted to unrelated TOPs (several
// of them with zero documents in 2026) were resolving as candidates and
// turning 881 of 2700 cells ambiguous.
//
// Same D-17 territory as icmsMatrixOriginUF: fixed for the client's single
// empresa, not a per-tenant choice.
var icmsMatrixTOPs = map[int]bool{306: true, 313: true}

// icmsMatrixCandidatesSQL reads every METALPRD.TGFICM row for the fixed
// origin (UFORIG=:1) with NO grupo/uf_destino filter — architecture decision
// (a) from the P2.b Task 3 dispatch: the whitelist predicate
// (domain.ResolveCell) is applied per (uf_destino, grupo) in Go, not folded
// into this SQL, so the extraction stays a single auditable read of the whole
// origin's matrix and never needs to know which grupos exist up front.
//
// UFDEST is translated to its text sigla via TSIUFS (TGFICM.UFDEST is a
// Sankhya-internal numeric code, NOT the IBGE code — MG=13 there, see the
// tgficm-ddl-traps evidence); CODPAIS=55 is required because TSIUFS is not
// scoped to Brazil alone.
//
// ALIQUFDEST is deliberately NOT selected here: it is the ST block's ceiling
// rate (icms_ceiling.go — a different, older feature) and is ~2.5x the
// operation's own rate. icms_matrix_mirror.aliquota must carry ALIQUOTA (the
// operation's rate); the wrong column would not fail any test, only produce a
// plausible, wrong number.
const icmsMatrixCandidatesSQL = `
SELECT u.UF, t.TIPRESTRICAO, t.CODRESTRICAO, t.TIPRESTRICAO2, t.CODRESTRICAO2,
       t.CODTRIB, t.ALIQUOTA, t.REDBASE, t.PERCICMSFCP
FROM METALPRD.TGFICM t
JOIN METALPRD.TSIUFS u ON u.CODUF = t.UFDEST AND u.CODPAIS = 55
WHERE t.UFORIG = :1`

// icmsKnownGruposSQL lists every fiscal grupo_icms a product actually carries.
// Read from TGFPRO directly (never a cross-database join against the
// Postgres mirror) so the Oracle extraction layer stays self-contained.
const icmsKnownGruposSQL = `SELECT DISTINCT GRUPOICMS FROM METALPRD.TGFPRO WHERE GRUPOICMS IS NOT NULL`

// ICMSMatrixReader is the read-only Oracle side of the ICMS matrix extraction
// (P2.b Task 3): it reads the raw TGFICM candidates and the set of known
// grupos, then resolves one cell per (uf_destino, grupo) via the pure
// domain.ResolveCell. It never writes to Oracle.
type ICMSMatrixReader struct {
	db queryer
}

// NewICMSMatrixReader builds a reader over an Oracle queryer handle.
func NewICMSMatrixReader(db queryer) *ICMSMatrixReader {
	return &ICMSMatrixReader{db: db}
}

// ResolveCells reads the current TGFICM snapshot for the fixed origin and the
// current set of known grupos, then resolves a domain.ICMSMatrixCell for
// every (uf_destino, grupo) combination that has at least one whitelist
// survivor. Combinations with zero survivors are omitted entirely (ADR-17:
// unknown never becomes a placeholder row) — the caller (the mirror writer)
// is responsible for closing any previously-open cell that no longer appears.
func (r *ICMSMatrixReader) ResolveCells(ctx context.Context) ([]domain.ICMSMatrixCell, error) {
	candidatesByUF, err := r.readCandidates(ctx)
	if err != nil {
		return nil, err
	}
	grupos, err := r.readKnownGrupos(ctx)
	if err != nil {
		return nil, err
	}

	ufs := make([]string, 0, len(candidatesByUF))
	for uf := range candidatesByUF {
		ufs = append(ufs, uf)
	}
	sort.Strings(ufs)

	var cells []domain.ICMSMatrixCell
	for _, uf := range ufs {
		rows := candidatesByUF[uf]
		for _, grupo := range grupos {
			resolved := domain.ResolveCell(rows, grupo, icmsMatrixTOPs)
			if !resolved.Found {
				continue
			}
			cells = append(cells, domain.ICMSMatrixCell{
				UFOrigem:         icmsMatrixOriginUFSigla,
				UFDestino:        uf,
				GrupoICMS:        grupo,
				CodTrib:          resolved.CodTrib,
				Aliquota:         resolved.Aliquota,
				RedBase:          resolved.RedBase,
				PercFcp:          resolved.PercFcp,
				LinhasCandidatas: resolved.LinhasCandidatas,
				Ambiguo:          resolved.Ambiguo,
			})
		}
	}
	return cells, nil
}

func (r *ICMSMatrixReader) readCandidates(ctx context.Context) (map[string][]domain.TgficmRow, error) {
	dbrows, err := r.db.QueryContext(ctx, icmsMatrixCandidatesSQL, icmsMatrixOriginUF)
	if err != nil {
		return nil, wrapOracleError("icms matrix candidates", err)
	}
	defer dbrows.Close()

	out := make(map[string][]domain.TgficmRow)
	for dbrows.Next() {
		var (
			uf            sql.NullString
			tipRestricao  sql.NullString
			codRestricao  sql.NullInt64
			tipRestricao2 sql.NullString
			codRestricao2 sql.NullInt64
			codTrib       sql.NullInt64
			aliquota      sql.NullFloat64
			redBase       sql.NullFloat64
			percFcp       sql.NullFloat64
		)
		if err := dbrows.Scan(&uf, &tipRestricao, &codRestricao, &tipRestricao2, &codRestricao2, &codTrib, &aliquota, &redBase, &percFcp); err != nil {
			return nil, wrapOracleError("scan icms matrix candidates", err)
		}
		if !uf.Valid {
			// TSIUFS join failed to resolve a sigla for this UFDEST — cannot place
			// the row against any destination, so it cannot become a candidate.
			continue
		}
		ufSigla := strings.TrimSpace(uf.String)
		if ufSigla == "" {
			continue
		}
		out[ufSigla] = append(out[ufSigla], domain.TgficmRow{
			TipRestricao:  strings.TrimSpace(tipRestricao.String),
			CodRestricao:  nullableInt(codRestricao),
			TipRestricao2: strings.TrimSpace(tipRestricao2.String),
			CodRestricao2: nullableInt(codRestricao2),
			CodTrib:       nullableInt(codTrib),
			Aliquota:      nullableFloat(aliquota),
			RedBase:       nullableFloat(redBase),
			PercFcp:       nullableFloat(percFcp),
		})
	}
	if err := dbrows.Err(); err != nil {
		return nil, wrapOracleError("iterate icms matrix candidates", err)
	}
	return out, nil
}

func (r *ICMSMatrixReader) readKnownGrupos(ctx context.Context) ([]int, error) {
	dbrows, err := r.db.QueryContext(ctx, icmsKnownGruposSQL)
	if err != nil {
		return nil, wrapOracleError("icms known grupos", err)
	}
	defer dbrows.Close()

	var out []int
	for dbrows.Next() {
		var grupo sql.NullInt64
		if err := dbrows.Scan(&grupo); err != nil {
			return nil, wrapOracleError("scan icms known grupos", err)
		}
		if grupo.Valid {
			out = append(out, int(grupo.Int64))
		}
	}
	if err := dbrows.Err(); err != nil {
		return nil, wrapOracleError("iterate icms known grupos", err)
	}
	sort.Ints(out)
	return out, nil
}
