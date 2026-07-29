package domain

import "time"

// Decimal is an exact base-10 value represented without floating-point loss.
type Decimal string

// NormalizedRow is the column-normalized ERP product row before validation.
// Custo and stock values remain textual until validation accepts them.
//
// This carries the 10 E2 business fields (interface-contracts-mis006.md §E2.1)
// that feed products_mirror:
//
//	codigo_produto → Codprod    descricao      → Descrprod
//	custo          → Custo      preco_venda    → PrecoVenda
//	estoque_total  → StockPhysical (+ StockReserved)
//	ean            → EAN        referencia     → Refforn
//	marca          → Marca      ncm            → NCM
//	grupo_codigo   → Grupo      grupo_descricao → DescrGrupo
//	local          → Local (per stock location, feeds products_mirror_stock_locations)
//
// Missing values stay nil/empty (honest-unknown, ADR-17) — never coerced to 0.
type NormalizedRow struct {
	Codprod   string  `json:"codprod"`
	Descrprod string  `json:"descrprod"`
	Custo     Decimal `json:"custo"`
	// PrecoVenda is the ERP sale price (E2.1). Textual until validation accepts
	// it; empty means absent, kept honest-unknown (never 0).
	PrecoVenda    Decimal `json:"preco_venda,omitempty"`
	StockPhysical string  `json:"stock_physical"`
	StockReserved *string `json:"stock_reserved,omitempty"`
	EAN           *string `json:"ean,omitempty"`
	Refforn       *string `json:"refforn,omitempty"`
	Marca         *string `json:"marca,omitempty"`
	NCM           *string `json:"ncm,omitempty"`
	// Grupo is the ERP product group code (grupo_codigo, E2.1).
	Grupo *string `json:"grupo,omitempty"`
	// DescrGrupo is the ERP product group description (grupo_descricao, E2.1).
	DescrGrupo *string `json:"descrgrupo,omitempty"`
	// Local is the stock-location code the row's quantity belongs to (E2.1).
	// Feeds products_mirror_stock_locations; nil when the source has no
	// per-location breakdown.
	Local       *string `json:"local,omitempty"`
	Usoprod     *string `json:"usoprod,omitempty"`
	ADEcommerce *string `json:"ad_ecommerce,omitempty"`
}

type IssueKind string

const (
	Rejection IssueKind = "REJECTION"
	Warning   IssueKind = "WARNING"
)

type IssueCode string

const (
	CodeEmptyCodprod          IssueCode = "EMPTY_CODPROD"
	CodeDuplicateCodprod      IssueCode = "DUPLICATE_CODPROD"
	CodeEmptyDescrprod        IssueCode = "EMPTY_DESCRPROD"
	CodeInvalidCusto          IssueCode = "INVALID_CUSTO"
	CodeInvalidEstoque        IssueCode = "INVALID_ESTOQUE"
	CodeInvalidEAN            IssueCode = "INVALID_EAN"
	CodeInvalidNCM            IssueCode = "INVALID_NCM"
	CodeMissingRequiredColumn IssueCode = "MISSING_REQUIRED_COLUMN"
	CodeInvalidFile           IssueCode = "INVALID_FILE"
	CodeImportInProgress      IssueCode = "IMPORT_IN_PROGRESS"
	CodeDuplicateFile         IssueCode = "DUPLICATE_FILE"
	CodeMissingCusto          IssueCode = "MISSING_CUSTO"
	CodeMissingEstoque        IssueCode = "MISSING_ESTOQUE"
)

type Issue struct {
	Row            int       `json:"row"`
	Column         *string   `json:"column,omitempty"`
	Kind           IssueKind `json:"kind"`
	Code           IssueCode `json:"code"`
	Detail         string    `json:"detail,omitempty"`
	OffendingValue *string   `json:"offending_value,omitempty"`
}

type ImportID string

// IsValid reports whether the id is syntactically the UUID that
// erp_import_protocols.id holds. Without this guard a malformed path value
// binds straight into a `uuid` column, Postgres raises 22P02, and — because
// that is not pgx.ErrNoRows — the caller receives internal_error and cannot
// tell "you sent junk" from "we broke".
func (id ImportID) IsValid() bool {
	value := string(id)
	if len(value) != 36 {
		return false
	}
	for index, r := range value {
		if index == 8 || index == 13 || index == 18 || index == 23 {
			if r != '-' {
				return false
			}
			continue
		}
		if !isHexDigit(r) {
			return false
		}
	}
	return true
}

func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

type Protocol string
type FileSHA256 string

type ImportSource string

const (
	SourceXLSX            ImportSource = "xlsx"
	SourceCatalogoCliente ImportSource = "catalogo_cliente"
)

type ImportStatus string

const (
	ImportStatusInProgress ImportStatus = "IN_PROGRESS"
	ImportStatusCompleted  ImportStatus = "COMPLETED"
	ImportStatusRejected   ImportStatus = "REJECTED"
)

type ImportSnapshot struct {
	ID           ImportID
	Protocol     Protocol
	FileSHA256   FileSHA256
	Source       ImportSource
	ImportedAt   time.Time
	Status       ImportStatus
	AcceptedRows []NormalizedRow
	Issues       []Issue
}

type ImportReport struct {
	ID            ImportID
	Protocol      Protocol
	FileSHA256    FileSHA256
	Source        ImportSource
	ImportedAt    time.Time
	Status        ImportStatus
	AcceptedCount int
	RejectedCount int
	WarningCount  int
	Issues        []Issue
}

// ImportChain answers what became of an import, read from the live tables
// rather than from the counters frozen at import time.
//
// Enfileirados is the CURRENT market queue, not a historical total: it falls as
// the queue is consumed. No consumer is registered today — the scheduler runs a
// products job only — so the count only grows until one is wired, and a later
// drop is drainage rather than data loss. QueueReadAt is the database clock of the same
// statement that produced the counts, so a reader can tell when that snapshot
// was taken instead of inferring a trend from two unrelated instants.
type ImportChain struct {
	Protocol     Protocol
	Importados   int
	Vinculados   int
	Enfileirados int
	QueueReadAt  time.Time
}
