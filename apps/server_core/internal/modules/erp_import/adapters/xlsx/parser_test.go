package xlsx

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"strings"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	"marketplace-central/apps/server_core/internal/modules/erp_import/ports"
)

type testSheet struct {
	name string
	rows [][]string
}

func TestParserValidAndInvalidRows(t *testing.T) {
	data := xlsxBytes(t, []testSheet{{
		name: "ERP",
		rows: [][]string{
			{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO", "ESTOQUE_RESERVADO", "EAN", "REFFORN", "MARCA", "NCM"},
			{"P1", "Produto 1", "12,34", "7", "2", "789", "REF-1", "Marca", "12345678"},
			{"P2", " Produto inválido ", "0", "-1", "", "", "", "", ""},
		},
	}})

	rows, _, err := NewParser().Parse(context.Background(), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("Parse() rows = %d, want 2", len(rows))
	}
	if rows[0].Custo != domain.Decimal("12,34") {
		t.Errorf("row 1 Custo = %q, want raw comma decimal", rows[0].Custo)
	}
	if rows[1].Custo != domain.Decimal("0") || rows[1].StockPhysical != "-1" {
		t.Errorf("invalid-value row was not preserved: %#v", rows[1])
	}
	if rows[1].Descrprod != " Produto inválido " {
		t.Errorf("Descrprod = %q, want raw cell content", rows[1].Descrprod)
	}
	if rows[0].StockReserved == nil || *rows[0].StockReserved != "2" {
		t.Errorf("StockReserved = %#v, want pointer to 2", rows[0].StockReserved)
	}
	if rows[1].StockReserved != nil || rows[1].EAN != nil || rows[1].Refforn != nil || rows[1].Marca != nil || rows[1].NCM != nil {
		t.Errorf("empty optional cells must be nil: %#v", rows[1])
	}
	// GRUPO/DESCRGRUPO columns absent from the header ⇒ honestly absent (nil),
	// never an empty string pretending a value (ADR-17). Covers legacy files
	// exported before the ERP grew the group columns.
	if rows[0].Grupo != nil || rows[0].DescrGrupo != nil || rows[1].Grupo != nil || rows[1].DescrGrupo != nil {
		t.Errorf("absent grupo columns must be nil: %#v / %#v", rows[0], rows[1])
	}
}

func TestParserSellableColumnsAreOptionalAndHonestUnknown(t *testing.T) {
	legacyData := xlsxBytes(t, []testSheet{{
		name: "ERP",
		rows: [][]string{
			{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO"},
			{"P1", "Produto 1", "12,34", "7"},
			{"P2", "Produto 2", "0", ""},
		},
	}})

	legacyRows, legacyIssues, err := NewParser().Parse(context.Background(), bytes.NewReader(legacyData))
	if err != nil {
		t.Fatalf("legacy Parse() error = %v", err)
	}
	if len(legacyRows) != 2 || len(legacyIssues) != 0 {
		t.Fatalf("legacy Parse() rows=%d issues=%d, want rows=2 issues=0", len(legacyRows), len(legacyIssues))
	}
	for index, row := range legacyRows {
		if row.Usoprod != nil || row.ADEcommerce != nil {
			t.Fatalf("legacy row %d sellable fields = usoprod=%#v ad_ecommerce=%#v, want nil", index, row.Usoprod, row.ADEcommerce)
		}
	}

	currentData := xlsxBytes(t, []testSheet{{
		name: "ERP",
		rows: [][]string{
			{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO", "USOPROD", "AD_ECOMMERCE"},
			{"P1", "Produto 1", "12,34", "7", "R", "S"},
			{"P2", "Produto 2", "0", "", "", ""},
		},
	}})

	currentRows, currentIssues, err := NewParser().Parse(context.Background(), bytes.NewReader(currentData))
	if err != nil {
		t.Fatalf("current Parse() error = %v", err)
	}
	if len(currentRows) != 2 || len(currentIssues) != 0 {
		t.Fatalf("current Parse() rows=%d issues=%d, want rows=2 issues=0", len(currentRows), len(currentIssues))
	}
	if currentRows[0].Usoprod == nil || *currentRows[0].Usoprod != "R" {
		t.Fatalf("current row 0 Usoprod = %#v, want R", currentRows[0].Usoprod)
	}
	if currentRows[0].ADEcommerce == nil || *currentRows[0].ADEcommerce != "S" {
		t.Fatalf("current row 0 ADEcommerce = %#v, want S", currentRows[0].ADEcommerce)
	}
	if currentRows[1].Usoprod != nil || currentRows[1].ADEcommerce != nil {
		t.Fatalf("blank current cells = usoprod=%#v ad_ecommerce=%#v, want nil", currentRows[1].Usoprod, currentRows[1].ADEcommerce)
	}
}

func TestParserCapturesGrupoColumns(t *testing.T) {
	data := xlsxBytes(t, []testSheet{{
		name: "ERP",
		rows: [][]string{
			{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO", "GRUPO", "DESCRGRUPO"},
			{"P1", "Produto 1", "12,34", "7", "07", "Ferramentas"},
			{"P2", "Produto 2", "5,00", "3", "", ""},
		},
	}})

	rows, _, err := NewParser().Parse(context.Background(), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if rows[0].Grupo == nil || *rows[0].Grupo != "07" {
		t.Errorf("row 1 Grupo = %#v, want pointer to 07", rows[0].Grupo)
	}
	if rows[0].DescrGrupo == nil || *rows[0].DescrGrupo != "Ferramentas" {
		t.Errorf("row 1 DescrGrupo = %#v, want pointer to Ferramentas", rows[0].DescrGrupo)
	}
	// Empty group cells on a row of a file that HAS the columns ⇒ still nil.
	if rows[1].Grupo != nil || rows[1].DescrGrupo != nil {
		t.Errorf("empty grupo cells must be nil: %#v", rows[1])
	}
}

func TestParserGrupoHeaderFoldsAccentAndCase(t *testing.T) {
	data := xlsxBytes(t, []testSheet{{
		name: "ERP",
		rows: [][]string{
			{"codprod", "descrprod", "custo", "estoque_fisico", " Grupo ", "DescrGrupo"},
			{"P1", "Produto 1", "9,90", "2", "12", "Solda"},
		},
	}})

	rows, _, err := NewParser().Parse(context.Background(), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if rows[0].Grupo == nil || *rows[0].Grupo != "12" {
		t.Errorf("folded Grupo header mapped incorrectly: %#v", rows[0].Grupo)
	}
	if rows[0].DescrGrupo == nil || *rows[0].DescrGrupo != "Solda" {
		t.Errorf("folded DescrGrupo header mapped incorrectly: %#v", rows[0].DescrGrupo)
	}
}

// TestParserHeaderFoldingAndUnionsAllSheets pins the multi-sheet UNION contract
// (M-01 defect fix: the parser previously read GetSheetList()[0] only and
// silently dropped every other sheet). Header folding/case-insensitivity is
// asserted on the first sheet; the second sheet's rows must also appear.
func TestParserHeaderFoldingAndUnionsAllSheets(t *testing.T) {
	data := xlsxBytes(t, []testSheet{
		{
			name: "First",
			rows: [][]string{
				{" Custo ", "ESTOQUE_FÍSICO", "codprod", " DESCRPROD "},
				{"10,50", "3", "P1", "Desc"},
			},
		},
		{
			name: "Second",
			rows: [][]string{
				{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO"},
				{"P2", "Segundo", "99", "99"},
			},
		},
	})

	rows, _, err := NewParser().Parse(context.Background(), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("Parse() rows = %d, want 2 (union of both sheets)", len(rows))
	}
	if rows[0].Codprod != "P1" || rows[0].Descrprod != "Desc" || rows[0].Custo != domain.Decimal("10,50") || rows[0].StockPhysical != "3" {
		t.Errorf("folded first-sheet headers mapped incorrectly: %#v", rows[0])
	}
	if rows[1].Codprod != "P2" || rows[1].Descrprod != "Segundo" {
		t.Errorf("second sheet row missing from union: %#v", rows[1])
	}
}

func TestParserMissingRequiredColumn(t *testing.T) {
	data := xlsxBytes(t, []testSheet{{
		name: "ERP",
		rows: [][]string{
			{"CODPROD", "DESCRPROD", "ESTOQUE_FISICO"},
			{"P1", "Desc", "3"},
		},
	}})

	rows, _, err := NewParser().Parse(context.Background(), bytes.NewReader(data))
	if rows != nil {
		t.Fatalf("Parse() rows = %#v, want nil on missing column", rows)
	}
	var fileErr *ports.FileError
	if !errors.As(err, &fileErr) {
		t.Fatalf("Parse() error = %T %v, want *ports.FileError", err, err)
	}
	if fileErr.Code != domain.CodeMissingRequiredColumn || fileErr.Column != "CUSTO" {
		t.Errorf("FileError = %#v, want missing CUSTO", fileErr)
	}
}

func TestParserInvalidFile(t *testing.T) {
	rows, _, err := NewParser().Parse(context.Background(), strings.NewReader("not an xlsx"))
	if rows != nil {
		t.Fatalf("Parse() rows = %#v, want nil for corrupt input", rows)
	}
	var fileErr *ports.FileError
	if !errors.As(err, &fileErr) {
		t.Fatalf("Parse() error = %T %v, want *ports.FileError", err, err)
	}
	if fileErr.Code != domain.CodeInvalidFile {
		t.Errorf("FileError.Code = %q, want %q", fileErr.Code, domain.CodeInvalidFile)
	}
}

func TestParserEmptyWorkbookOrFirstSheet(t *testing.T) {
	tests := []struct {
		name   string
		sheets []testSheet
	}{
		{name: "empty workbook"},
		{
			name: "empty first sheet",
			sheets: []testSheet{{
				name: "ERP",
				rows: [][]string{{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO"}},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := xlsxBytes(t, tt.sheets)
			rows, _, err := NewParser().Parse(context.Background(), bytes.NewReader(data))
			if rows != nil {
				t.Fatalf("Parse() rows = %#v, want nil", rows)
			}
			var fileErr *ports.FileError
			if !errors.As(err, &fileErr) {
				t.Fatalf("Parse() error = %T %v, want *ports.FileError", err, err)
			}
			if fileErr.Code != domain.CodeInvalidFile {
				t.Errorf("FileError.Code = %q, want %q", fileErr.Code, domain.CodeInvalidFile)
			}
		})
	}
}

func xlsxBytes(t *testing.T, sheets []testSheet) []byte {
	t.Helper()
	var output bytes.Buffer
	archive := zip.NewWriter(&output)
	writeZipFile(t, archive, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/><Override PartName="/xl/worksheets/sheet2.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>`)
	writeZipFile(t, archive, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`)
	writeZipFile(t, archive, "xl/_rels/workbook.xml.rels", workbookRelationships(len(sheets)))
	writeZipFile(t, archive, "xl/workbook.xml", workbookXML(sheets))
	for i, sheet := range sheets {
		writeZipFile(t, archive, "xl/worksheets/sheet"+itoa(i+1)+".xml", sheetXML(sheet.rows))
	}
	if err := archive.Close(); err != nil {
		t.Fatalf("close xlsx fixture: %v", err)
	}
	return output.Bytes()
}

func writeZipFile(t *testing.T, archive *zip.Writer, name, content string) {
	t.Helper()
	file, err := archive.Create(name)
	if err != nil {
		t.Fatalf("create xlsx fixture entry %q: %v", name, err)
	}
	if _, err := io.WriteString(file, content); err != nil {
		t.Fatalf("write xlsx fixture entry %q: %v", name, err)
	}
}

func workbookRelationships(count int) string {
	var relationships strings.Builder
	relationships.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	for i := 1; i <= count; i++ {
		relationships.WriteString(`<Relationship Id="rId` + itoa(i) + `" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet` + itoa(i) + `.xml"/>`)
	}
	relationships.WriteString(`</Relationships>`)
	return relationships.String()
}

func workbookXML(sheets []testSheet) string {
	var workbook strings.Builder
	workbook.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets>`)
	for i, sheet := range sheets {
		workbook.WriteString(`<sheet name="`)
		escape(&workbook, sheet.name)
		workbook.WriteString(`" sheetId="` + itoa(i+1) + `" r:id="rId` + itoa(i+1) + `"/>`)
	}
	workbook.WriteString(`</sheets></workbook>`)
	return workbook.String()
}

func sheetXML(rows [][]string) string {
	var sheet strings.Builder
	sheet.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for rowIndex, row := range rows {
		sheet.WriteString(`<row r="` + itoa(rowIndex+1) + `">`)
		for columnIndex, value := range row {
			if value == "" {
				continue
			}
			sheet.WriteString(`<c r="` + cellReference(columnIndex, rowIndex) + `" t="inlineStr"><is><t>`)
			escape(&sheet, value)
			sheet.WriteString(`</t></is></c>`)
		}
		sheet.WriteString(`</row>`)
	}
	sheet.WriteString(`</sheetData></worksheet>`)
	return sheet.String()
}

func cellReference(column, row int) string {
	return string(rune('A'+column)) + itoa(row+1)
}

func escape(builder *strings.Builder, value string) {
	_ = xml.EscapeText(builder, []byte(value))
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	i := len(digits)
	for value > 0 {
		i--
		digits[i] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[i:])
}
