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

func (p *Parser) Parse(ctx context.Context, source io.Reader) ([]domain.NormalizedRow, error) {
	return p.parseWithRequired(ctx, source, []string{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO"})
}

// ParseLenient parses a workbook that may omit CUSTO / ESTOQUE_FISICO (the
// client-catalog use case). Missing columns yield honest-empty values rather
// than a structural rejection; ValidateRowsLenient turns that into warnings
// instead of fabricating zeros (ADR-17).
func (p *Parser) ParseLenient(ctx context.Context, source io.Reader) ([]domain.NormalizedRow, error) {
	return p.parseWithRequired(ctx, source, []string{"CODPROD", "DESCRPROD"})
}

func (p *Parser) parseWithRequired(ctx context.Context, source io.Reader, required []string) ([]domain.NormalizedRow, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if source == nil {
		return nil, invalidFileError()
	}

	workbook, err := excelize.OpenReader(source)
	if err != nil {
		return nil, invalidFileError()
	}
	defer func() {
		_ = workbook.Close()
	}()

	sheets := workbook.GetSheetList()
	if len(sheets) == 0 {
		return nil, invalidFileError()
	}

	rows, err := workbook.GetRows(sheets[0])
	if err != nil {
		return nil, invalidFileError()
	}
	if len(rows) < 2 {
		return nil, invalidFileError()
	}

	columns := make(map[string]int, len(rows[0]))
	for index, header := range rows[0] {
		key := normalizeHeader(header)
		if key != "" {
			columns[key] = index
		}
	}
	for _, requiredColumn := range required {
		if _, ok := columns[normalizeHeader(requiredColumn)]; !ok {
			return nil, &ports.FileError{
				Code:   domain.CodeMissingRequiredColumn,
				Column: requiredColumn,
				Detail: "required column is missing",
			}
		}
	}

	result := make([]domain.NormalizedRow, 0, len(rows)-1)
	for _, row := range rows[1:] {
		result = append(result, domain.NormalizedRow{
			Codprod:       columnCell(row, columns, "CODPROD"),
			Descrprod:     columnCell(row, columns, "DESCRPROD"),
			Custo:         domain.Decimal(columnCell(row, columns, "CUSTO")),
			StockPhysical: columnCell(row, columns, "ESTOQUE_FISICO"),
			StockReserved: optionalCell(row, columns, "ESTOQUE_RESERVADO"),
			EAN:           optionalCell(row, columns, "EAN"),
			Refforn:       optionalCell(row, columns, "REFFORN"),
			Marca:         optionalCell(row, columns, "MARCA"),
			NCM:           optionalCell(row, columns, "NCM"),
			Grupo:         optionalCell(row, columns, "GRUPO"),
			DescrGrupo:    optionalCell(row, columns, "DESCRGRUPO"),
		})
	}
	return result, nil
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

// columnCell reads a column that may be entirely absent from the header row
// (the lenient path). When present it behaves exactly like the direct
// columns[...] lookup the strict path used to perform inline, since required
// columns are guaranteed present by the caller's required-column check above.
func columnCell(row []string, columns map[string]int, name string) string {
	index, ok := columns[normalizeHeader(name)]
	if !ok {
		return ""
	}
	return cell(row, index)
}

func optionalCell(row []string, columns map[string]int, name string) *string {
	index, ok := columns[normalizeHeader(name)]
	if !ok {
		return nil
	}
	value := cell(row, index)
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return &value
}

func invalidFileError() *ports.FileError {
	return &ports.FileError{
		Code:   domain.CodeInvalidFile,
		Detail: "invalid xlsx file",
	}
}
