package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	catalogdomain "marketplace-central/apps/server_core/internal/modules/catalog/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/adapters/mirror"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
	"marketplace-central/apps/server_core/internal/modules/sourcekind"
)

// sankhyaCompany is CODEMP=1: the client is single-empresa METALPRD (ratified
// mapping, docs/design/evidence/sankhya-mapping-M04.md §"Empresa"). Revisit only
// for a multi-empresa tenant, which the foundation is not.
const sankhyaCompany = 1

// SankhyaAdapter is the live-read-through ProductSourceAdapter for a single
// tenant's Sankhya (Oracle) ERP. It embeds the existing read-side *Reader verbatim
// (no read signature changes — M04-C3/C9) and adds the write/sync side: a
// full-snapshot Sync that feeds products_mirror through the mirror.Writer.
//
// tenantID is fixed per adapter because one Oracle connection == one tenant's ERP;
// the Sankhya tables carry no tenant column, so tenant scoping is applied where the
// mirror is written (AC-01), not in the Oracle read.
type SankhyaAdapter struct {
	*Reader
	mirror   mirror.Writer
	tenantID string
	now      func() time.Time
}

// NewSankhyaAdapter builds the adapter over an Oracle handle, the mirror writer,
// and the tenant this ERP connection belongs to. Composition (root.go, hub-owned
// F1 pre-activation) supplies the concrete handles.
func NewSankhyaAdapter(oracleDB queryer, mirrorWriter mirror.Writer, tenantID string) *SankhyaAdapter {
	return &SankhyaAdapter{
		Reader:   NewReader(oracleDB),
		mirror:   mirrorWriter,
		tenantID: tenantID,
		now:      time.Now,
	}
}

var _ ports.ProductSourceAdapter = (*SankhyaAdapter)(nil)

// Kind reports this adapter as protocol-less live-read-through (F-02).
func (a *SankhyaAdapter) Kind() sourcekind.SourceKind {
	return sourcekind.LiveReadThrough
}

// Sync reads the current Sankhya snapshot (as-of today) and applies it to
// products_mirror via keep-absent merge (ADR-04) with honest-NULL fields
// (ADR-17). Full-snapshot v1: keep-absent needs the complete set and TGFEST has no
// change timestamp, so every run reads the whole active catalogue (ratified sync
// strategy, mapping doc §"Estratégia de sync").
func (a *SankhyaAdapter) Sync(ctx context.Context) (ports.SyncResult, error) {
	if err := a.ensureAvailable(ctx); err != nil {
		return ports.SyncResult{}, err
	}
	dataRef := a.now().UTC()

	rows, err := a.readBase(ctx)
	if err != nil {
		return ports.SyncResult{}, err
	}
	if err := a.applyCost(ctx, dataRef, rows); err != nil {
		return ports.SyncResult{}, err
	}
	if err := a.applyPrice(ctx, dataRef, rows); err != nil {
		return ports.SyncResult{}, err
	}
	locs, err := a.applyStock(ctx, rows)
	if err != nil {
		return ports.SyncResult{}, err
	}

	mirrorRows := make([]mirror.Row, 0, len(rows))
	for _, r := range rows {
		mirrorRows = append(mirrorRows, r.row)
	}

	processed, err := a.mirror.ApplySnapshot(ctx, a.tenantID, mirrorRows, locs)
	if err != nil {
		return ports.SyncResult{}, fmt.Errorf("sankhya sync: apply snapshot: %w", err)
	}
	return ports.SyncResult{Processed: processed, Errors: 0}, nil
}

// snapshotRow carries a mirror.Row plus its numeric CODPROD so the per-product
// cost/price/stock passes can join without re-parsing the string key.
type snapshotRow struct {
	codprod int
	row     mirror.Row
}

// Q1 — products_mirror BASE (identidade + grupo + marca + ncm + ean), ATIVO='S'.
// Ratified mapping doc §"Pack de queries" Q1. grupo FOLHA via TGFPRO.CODGRUPOPROD
// (already a leaf); marca via canonical CODMARCA join, not the denormalized text.
const sankhyaBaseSQL = `
SELECT p.CODPROD, p.DESCRPROD, p.NCM, p.REFERENCIA, p.REFFORN,
       p.CODGRUPOPROD, g.DESCRGRUPOPROD, m.DESCRICAO,
       p.USOPROD, p.AD_ECOMMERCE
FROM METALPRD.TGFPRO p
LEFT JOIN METALPRD.TGFGRU g ON g.CODGRUPOPROD = p.CODGRUPOPROD
LEFT JOIN METALPRD.TGFMAR m ON m.CODIGO       = p.CODMARCA
WHERE p.ATIVO = 'S'`

func (a *SankhyaAdapter) readBase(ctx context.Context) (map[int]*snapshotRow, error) {
	dbrows, err := a.db.QueryContext(ctx, sankhyaBaseSQL)
	if err != nil {
		return nil, wrapOracleError("sankhya base", err)
	}
	defer dbrows.Close()

	out := make(map[int]*snapshotRow)
	for dbrows.Next() {
		var (
			codprod   int
			descr     sql.NullString
			ncm       sql.NullString
			reference sql.NullString // REFERENCIA carries the EAN (name collision)
			refforn   sql.NullString // supplier reference, NOT the EAN
			grupoCod  sql.NullInt64
			grupoDesc sql.NullString
			marca     sql.NullString
			usoprod   sql.NullString
			adEcomm   sql.NullString
		)
		if err := dbrows.Scan(&codprod, &descr, &ncm, &reference, &refforn, &grupoCod, &grupoDesc, &marca, &usoprod, &adEcomm); err != nil {
			return nil, wrapOracleError("scan sankhya base", err)
		}

		// EAN lives in REFERENCIA but only when EAN-shaped; ~51% are not → NULL
		// (honest-unknown). Trim first: Oracle CHAR/space-padding would otherwise fail
		// the digit check and drop a genuine EAN to NULL. Validate AND store the
		// trimmed value so the mirror never persists a padded identifier. Shape rule =
		// catalogdomain.IsValidGTIN, deliberately identical to read-side
		// FindProductsForLinking (reader.go) so a product's EAN is the same fact whether
		// served live-read or from the mirror snapshot (cross-source consistency).
		var ean *string
		if reference.Valid {
			if trimmed := strings.TrimSpace(reference.String); catalogdomain.IsValidGTIN(trimmed) {
				ean = &trimmed
			}
		}

		out[codprod] = &snapshotRow{
			codprod: codprod,
			row: mirror.Row{
				CodigoProduto:  strconv.Itoa(codprod),
				Descricao:      nullStr(descr),
				Referencia:     nullStr(refforn),
				EAN:            ean,
				Marca:          nullStr(marca),
				GrupoCodigo:    nullIntStr(grupoCod),
				GrupoDescricao: nullStr(grupoDesc),
				NCM:            nullStr(ncm),
				Usoprod:        nullStr(usoprod),
				ADEcommerce:    nullStr(adEcomm),
			},
		}
	}
	if err := dbrows.Err(); err != nil {
		return nil, wrapOracleError("iterate sankhya base", err)
	}
	return out, nil
}

// Q2 — CUSTO as-of (CUSSEMICM, líquido de ICMS), CODEMP=1. Set-based equivalent of
// the ratified per-product query: ROW_NUMBER partitioned by CODPROD over the same
// filter/order (DTATUAL DESC, CODLOCAL DESC) picking rn=1 is exactly the per-product
// "FETCH FIRST 1 ROW ONLY". No row → custo stays NULL.
const sankhyaCostSQL = `
SELECT CODPROD, CUSSEMICM FROM (
	SELECT c.CODPROD, c.CUSSEMICM,
		ROW_NUMBER() OVER (PARTITION BY c.CODPROD ORDER BY c.DTATUAL DESC, c.CODLOCAL DESC) rn
	FROM METALPRD.TGFCUS c
	WHERE c.CODEMP = :1 AND c.DTATUAL <= :2
) WHERE rn = 1`

func (a *SankhyaAdapter) applyCost(ctx context.Context, dataRef time.Time, rows map[int]*snapshotRow) error {
	dbrows, err := a.db.QueryContext(ctx, sankhyaCostSQL, sankhyaCompany, dataRef)
	if err != nil {
		return wrapOracleError("sankhya cost", err)
	}
	defer dbrows.Close()
	for dbrows.Next() {
		var (
			codprod int
			custo   sql.NullFloat64
		)
		if err := dbrows.Scan(&codprod, &custo); err != nil {
			return wrapOracleError("scan sankhya cost", err)
		}
		if r, ok := rows[codprod]; ok && custo.Valid {
			v := custo.Float64
			r.row.Custo = &v
		}
	}
	if err := dbrows.Err(); err != nil {
		return wrapOracleError("iterate sankhya cost", err)
	}
	return nil
}

// Q3 — PREÇO DE LISTA as-of (tabela de venda CODTAB=0). Set-based equivalent of the
// per-product query; ORDER BY DTVIGOR DESC matches the ratified query. NUTAB DESC is
// an extra deterministic tiebreak (a higher NUTAB is the later-created version, so
// the newest wins) — the ratified query, being per-product FETCH-FIRST-1, leaves a
// same-DTVIGOR tie unspecified; in real data CODTAB=0 has ~one version/day (mapping
// doc FATO 1) so the tie is practically absent, but pinning it keeps the snapshot
// reproducible run-to-run. No row → preco_venda stays NULL.
const sankhyaPriceSQL = `
SELECT CODPROD, VLRVENDA FROM (
	SELECT e.CODPROD, e.VLRVENDA,
		ROW_NUMBER() OVER (PARTITION BY e.CODPROD ORDER BY t.DTVIGOR DESC, e.NUTAB DESC) rn
	FROM METALPRD.TGFEXC e
	JOIN METALPRD.TGFTAB t ON t.NUTAB = e.NUTAB
	WHERE t.CODTAB = 0 AND t.DTVIGOR <= :1
) WHERE rn = 1`

func (a *SankhyaAdapter) applyPrice(ctx context.Context, dataRef time.Time, rows map[int]*snapshotRow) error {
	dbrows, err := a.db.QueryContext(ctx, sankhyaPriceSQL, dataRef)
	if err != nil {
		return wrapOracleError("sankhya price", err)
	}
	defer dbrows.Close()
	for dbrows.Next() {
		var (
			codprod int
			preco   sql.NullFloat64
		)
		if err := dbrows.Scan(&codprod, &preco); err != nil {
			return wrapOracleError("scan sankhya price", err)
		}
		if r, ok := rows[codprod]; ok && preco.Valid {
			v := preco.Float64
			r.row.PrecoVenda = &v
		}
	}
	if err := dbrows.Err(); err != nil {
		return wrapOracleError("iterate sankhya price", err)
	}
	return nil
}

// Q4 — available sellable stock per whitelisted location: company-owned stock
// net of reservations in companies 1/2 and locations 1_REVENDA/2_OUTLET. Every
// Q1 product absent from this pinned result is known-zero; products absent from
// Q1 are untouched by the snapshot merge. Before the sellable pin, F-01 treated
// absence from the all-TGFEST query as NULL.
const sankhyaSellableLocationCodes = "10101, 10102" // 10101=1_REVENDA; 10102=2_OUTLET

const sankhyaStockSQL = `
SELECT CODPROD, CODLOCAL, SUM(NVL(ESTOQUE, 0) - NVL(RESERVADO, 0)) AS DISPONIVEL
FROM METALPRD.TGFEST
WHERE CODPARC = 0
AND CODEMP IN (1, 2)
AND CODLOCAL IN (` + sankhyaSellableLocationCodes + `)
GROUP BY CODPROD, CODLOCAL`

func (a *SankhyaAdapter) applyStock(ctx context.Context, rows map[int]*snapshotRow) ([]mirror.StockLocation, error) {
	dbrows, err := a.db.QueryContext(ctx, sankhyaStockSQL)
	if err != nil {
		return nil, wrapOracleError("sankhya stock", err)
	}
	defer dbrows.Close()

	totals := make(map[int]float64)
	var locs []mirror.StockLocation
	for dbrows.Next() {
		var (
			codprod    int
			codlocal   sql.NullInt64
			disponivel sql.NullFloat64
		)
		if err := dbrows.Scan(&codprod, &codlocal, &disponivel); err != nil {
			return nil, wrapOracleError("scan sankhya stock", err)
		}
		r, ok := rows[codprod]
		if !ok || !disponivel.Valid {
			continue
		}
		qty := disponivel.Float64
		totals[codprod] += qty
		locs = append(locs, mirror.StockLocation{
			CodigoProduto: r.row.CodigoProduto,
			LocalCodigo:   nullInt64ToStr(codlocal),
			Quantidade:    &qty,
		})
	}
	if err := dbrows.Err(); err != nil {
		return nil, wrapOracleError("iterate sankhya stock", err)
	}
	for codprod, snapshot := range rows {
		total := totals[codprod]
		snapshot.row.EstoqueTotal = &total
	}
	return locs, nil
}

// nullStr maps an Oracle text column to *string with honest-NULL discipline: an
// absent value, or one that is empty/whitespace-only after trimming (Oracle
// CHAR-padding, blank cells), is unknown → nil, never a "" that would read as a
// known-empty fact. The trimmed value is stored so no padding leaks into the mirror.
func nullStr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := strings.TrimSpace(v.String)
	if s == "" {
		return nil
	}
	return &s
}

func nullIntStr(v sql.NullInt64) *string {
	if !v.Valid {
		return nil
	}
	s := strconv.FormatInt(v.Int64, 10)
	return &s
}

func nullInt64ToStr(v sql.NullInt64) string {
	if !v.Valid {
		return ""
	}
	return strconv.FormatInt(v.Int64, 10)
}
