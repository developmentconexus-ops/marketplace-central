package xlsx

import (
	"context"
	"io"
	"strings"
	"unicode"

	"github.com/xuri/excelize/v2"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

var _ ports.Parser = (*Parser)(nil)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(ctx context.Context, source io.Reader) ([]domain.NormalizedRow, []domain.Issue, error) {
	return p.parseWithRequired(ctx, source, []string{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO"}, false)
}

// ParseLenient parses a workbook that may omit CUSTO / ESTOQUE_FISICO (the
// client-catalog use case). Missing columns yield honest-empty values rather
// than a structural rejection; ValidateRowsLenient turns that into warnings
// instead of fabricating zeros (ADR-17). It also injects the sheet name as the
// product category (DescrGrupo) when the sheet carries no explicit group
// column — the raw Sankhya "Produto" export groups products across FACAS /
// FERRAMENTAS / MAQUINAS sheets with no in-sheet group code.
func (p *Parser) ParseLenient(ctx context.Context, source io.Reader) ([]domain.NormalizedRow, []domain.Issue, error) {
	return p.parseWithRequired(ctx, source, []string{"CODPROD", "DESCRPROD"}, true)
}

// parseWithRequired reads every worksheet in the workbook, detecting each
// sheet's header row past any preamble rows and UNIONing the data rows across
// sheets (GetSheetList order). Rows above a sheet's header are skipped; a sheet
// with no detectable product header is skipped with a WARNING file-issue rather
// than failing the whole import (ADR-17: a stray sheet never silently drops
// products, and never hard-fails the sheets that DO carry data). The strict
// required-column rejection still fires: a detected header missing a required
// column returns a MISSING_REQUIRED_COLUMN FileError for the whole file.
//
// Displayed issue row numbers are the position of the row in the concatenated
// cross-sheet stream (monotonic 1..N over accepted-order), not the original
// spreadsheet row — the NormalizedRow carries no origin coordinate and adding
// one is a schema change out of this seam's scope. File-level (sheet-skip)
// warnings carry row 0.
func (p *Parser) parseWithRequired(ctx context.Context, source io.Reader, required []string, sheetNameAsCategory bool) ([]domain.NormalizedRow, []domain.Issue, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if source == nil {
		return nil, nil, invalidFileError()
	}

	workbook, err := excelize.OpenReader(source)
	if err != nil {
		return nil, nil, invalidFileError()
	}
	defer func() {
		_ = workbook.Close()
	}()

	sheets := workbook.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, invalidFileError()
	}

	var (
		result     []domain.NormalizedRow
		fileIssues []domain.Issue
	)
	for _, sheetName := range sheets {
		rows, err := workbook.GetRows(sheetName)
		if err != nil {
			fileIssues = append(fileIssues, sheetSkippedIssue(sheetName, "sheet could not be read; skipped"))
			continue
		}

		headerIndex, columns := detectHeader(rows)
		if headerIndex < 0 {
			fileIssues = append(fileIssues, sheetSkippedIssue(sheetName, "no product header row (CODPROD/DESCRPROD) found; sheet skipped"))
			continue
		}

		// Strict required-column rejection: a detected header that is missing a
		// required column fails the whole file (genuinely-wrong file guard).
		for _, requiredColumn := range required {
			if _, ok := columns[canonicalHeaderKey(requiredColumn)]; !ok {
				return nil, nil, &ports.FileError{
					Code:   domain.CodeMissingRequiredColumn,
					Column: requiredColumn,
					Detail: "required column is missing",
				}
			}
		}

		_, hasDescrGrupo := columns[canonicalHeaderKey("DESCRGRUPO")]
		for _, row := range rows[headerIndex+1:] {
			normalized := domain.NormalizedRow{
				Codprod:       columnCell(row, columns, "CODPROD"),
				Descrprod:     columnCell(row, columns, "DESCRPROD"),
				Custo:         domain.Decimal(columnCell(row, columns, "CUSTO")),
				StockPhysical: columnCell(row, columns, "ESTOQUE_FISICO"),
				StockReserved: optionalCell(row, columns, "ESTOQUE_RESERVADO"),
				EAN:           optionalCell(row, columns, "EAN"),
				Refforn:       optionalCell(row, columns, "REFFORN"),
				Marca:         brandCell(row, columns),
				NCM:           optionalCell(row, columns, "NCM"),
				Grupo:         optionalCell(row, columns, "GRUPO"),
				DescrGrupo:    optionalCell(row, columns, "DESCRGRUPO"),
			}
			// Category from sheet name (lenient path only): only when the sheet
			// has no explicit group column, so files that DO carry group columns
			// keep their honest per-row values (empty cell stays nil, ADR-17).
			if sheetNameAsCategory && !hasDescrGrupo {
				category := sheetName
				normalized.DescrGrupo = &category
			}
			result = append(result, normalized)
		}
	}

	if len(result) == 0 {
		return nil, nil, invalidFileError()
	}
	return result, fileIssues, nil
}

// detectHeader returns the index of the first row that carries both identity
// columns (CODPROD and DESCRPROD, English or Portuguese Sankhya aliases) along
// with that row's canonical column map. Rows above it (title / emission
// preamble) are ignored. Returns -1 when no such row exists.
func detectHeader(rows [][]string) (int, map[string][]int) {
	for index, row := range rows {
		columns := buildColumns(row)
		_, hasCodprod := columns[canonicalHeaderKey("CODPROD")]
		_, hasDescrprod := columns[canonicalHeaderKey("DESCRPROD")]
		if hasCodprod && hasDescrprod {
			return index, columns
		}
	}
	return -1, nil
}

// buildColumns maps each canonical column key to the header indices carrying it.
// A key can map to several indices (the raw export ships two "Marca" columns);
// readers take the first non-empty. Unrecognized headers are ignored.
func buildColumns(header []string) map[string][]int {
	columns := make(map[string][]int, len(header))
	for index, headerCell := range header {
		key := canonicalHeaderKey(headerCell)
		if key == "" {
			continue
		}
		columns[key] = append(columns[key], index)
	}
	return columns
}

// canonicalHeaderKey folds a raw header label to a canonical column key,
// resolving both the English strict-path headers and the Portuguese Sankhya
// raw-export aliases to the same key. Unknown headers fold to "" (ignored).
func canonicalHeaderKey(raw string) string {
	switch normalizeHeader(raw) {
	case "codprod", "codigo":
		return "codprod"
	case "descrprod", "descricao":
		return "descrprod"
	case "custo":
		return "custo"
	case "estoque_fisico":
		return "estoque_fisico"
	case "estoque_reservado":
		return "estoque_reservado"
	case "ean", "codigo de barra", "codigo de barras":
		return "ean"
	case "refforn", "referencia do fornecedor":
		return "refforn"
	case "marca":
		return "marca"
	case "ncm":
		return "ncm"
	case "grupo":
		return "grupo"
	case "descrgrupo":
		return "descrgrupo"
	default:
		return ""
	}
}

// sheetSkippedIssue records a skipped sheet as a WARNING. It reuses the
// MISSING_REQUIRED_COLUMN code (truthful: the sheet lacks the required identity
// header) with column set to the sheet name, so persistence stays within the
// existing erp_import_issues code CHECK — a dedicated SHEET_SKIPPED code would
// need a migration (hub-owned) and the skipped-sheet path must never hard-fail
// the import, not even at the DB constraint.
func sheetSkippedIssue(sheetName, detail string) domain.Issue {
	column := sheetName
	value := sheetName
	return domain.Issue{
		Row:            0,
		Column:         &column,
		Kind:           domain.Warning,
		Code:           domain.CodeMissingRequiredColumn,
		Detail:         detail,
		OffendingValue: &value,
	}
}

func normalizeHeader(value string) string {
	return strings.Map(foldAccent, strings.ToLower(strings.TrimSpace(value)))
}

func foldAccent(r rune) rune {
	if unicode.Is(unicode.Mn, r) {
		return -1
	}
	switch r {
	case 'á', 'à', 'â', 'ã', 'ä', 'å':
		return 'a'
	case 'ç':
		return 'c'
	case 'ď', 'đ':
		return 'd'
	case 'é', 'è', 'ê', 'ë':
		return 'e'
	case 'í', 'ì', 'î', 'ï':
		return 'i'
	case 'ñ':
		return 'n'
	case 'ó', 'ò', 'ô', 'õ', 'ö', 'ø':
		return 'o'
	case 'ř':
		return 'r'
	case 'š', 'ș', 'ş':
		return 's'
	case 'ť':
		return 't'
	case 'ú', 'ù', 'û', 'ü':
		return 'u'
	case 'ý', 'ÿ':
		return 'y'
	case 'ž':
		return 'z'
	default:
		return r
	}
}

func cell(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}
	return row[index]
}

// columnCell returns the first non-empty cell among the header indices mapped to
// a column (raw value preserved untrimmed — validation owns trimming). Required
// columns are guaranteed present by the caller's required-column check; an
// absent or all-empty column yields "".
func columnCell(row []string, columns map[string][]int, name string) string {
	for _, index := range columns[canonicalHeaderKey(name)] {
		if strings.TrimSpace(cell(row, index)) != "" {
			return cell(row, index)
		}
	}
	return ""
}

// optionalCell returns a pointer to the first non-empty cell among a column's
// header indices ("first non-empty wins" across the raw export's twin Marca
// columns), or nil when the column is absent or every mapped cell is empty.
func optionalCell(row []string, columns map[string][]int, name string) *string {
	for _, index := range columns[canonicalHeaderKey(name)] {
		value := cell(row, index)
		if strings.TrimSpace(value) != "" {
			return &value
		}
	}
	return nil
}

// brandCell resolves MARCA across the raw export's twin "Marca" columns, where
// the first carries the Sankhya brand CODE and the second the brand NAME. A
// purely numeric value is the code, not a name — presenting it as the brand
// contradicts marketplace catalogs (EAN matches were rejected on marca), so
// the first non-empty NON-numeric value wins and an all-numeric/empty set is
// an honest unknown (nil), never a fabricated name.
func brandCell(row []string, columns map[string][]int) *string {
	for _, index := range columns[canonicalHeaderKey("MARCA")] {
		value := strings.TrimSpace(cell(row, index))
		if value == "" || isAllDigits(value) {
			continue
		}
		raw := cell(row, index)
		return &raw
	}
	return nil
}

func isAllDigits(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func invalidFileError() *ports.FileError {
	return &ports.FileError{
		Code:   domain.CodeInvalidFile,
		Detail: "invalid xlsx file",
	}
}
