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
	for _, required := range []string{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO"} {
		if _, ok := columns[normalizeHeader(required)]; !ok {
			return nil, &ports.FileError{
				Code:   domain.CodeMissingRequiredColumn,
				Column: required,
				Detail: "required column is missing",
			}
		}
	}

	result := make([]domain.NormalizedRow, 0, len(rows)-1)
	for _, row := range rows[1:] {
		result = append(result, domain.NormalizedRow{
			Codprod:       cell(row, columns[normalizeHeader("CODPROD")]),
			Descrprod:     cell(row, columns[normalizeHeader("DESCRPROD")]),
			Custo:         domain.Decimal(cell(row, columns[normalizeHeader("CUSTO")])),
			StockPhysical: cell(row, columns[normalizeHeader("ESTOQUE_FISICO")]),
			StockReserved: optionalCell(row, columns, "ESTOQUE_RESERVADO"),
			EAN:           optionalCell(row, columns, "EAN"),
			Refforn:       optionalCell(row, columns, "REFFORN"),
			Marca:         optionalCell(row, columns, "MARCA"),
			NCM:           optionalCell(row, columns, "NCM"),
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
