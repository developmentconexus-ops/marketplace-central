package domain

import "time"

// Decimal is an exact base-10 value represented without floating-point loss.
type Decimal string

// NormalizedRow is the column-normalized ERP product row before validation.
// Custo and stock values remain textual until validation accepts them.
type NormalizedRow struct {
	Codprod       string  `json:"codprod"`
	Descrprod     string  `json:"descrprod"`
	Custo         Decimal `json:"custo"`
	StockPhysical string  `json:"stock_physical"`
	StockReserved *string `json:"stock_reserved,omitempty"`
	EAN           *string `json:"ean,omitempty"`
	Refforn       *string `json:"refforn,omitempty"`
	Marca         *string `json:"marca,omitempty"`
	NCM           *string `json:"ncm,omitempty"`
	Grupo         *string `json:"grupo,omitempty"`
	DescrGrupo    *string `json:"descrgrupo,omitempty"`
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
type Protocol string
type FileSHA256 string

type ImportSource string

const SourceXLSX ImportSource = "xlsx"

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
