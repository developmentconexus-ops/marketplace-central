package xlsx

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/xuri/excelize/v2"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

// rawSheet mirrors one worksheet of the customer's real Sankhya "Produto"
// export: two preamble rows (title + emission line) above a Portuguese header
// row, then the product rows.
type rawSheet struct {
	name     string
	total    int
	products [][]string // each row aligned to rawHeader
}

// rawHeader is the Portuguese Sankhya header shape from the ground-truth
// FINDING: two "Marca" columns, "Código de Barra", "Referência do Fornecedor",
// and "Estoque Mínimo" (which is NOT physical stock) — and crucially NO CUSTO
// and NO ESTOQUE_FISICO column.
var rawHeader = []string{
	"Código", "Descrição", "Marca", "Marca", "Código de Barra",
	"NCM", "Referência do Fornecedor", "Estoque Mínimo",
}

// buildRawSankhyaExport authors a workbook with excelize that mirrors the raw
// export shape exactly: for each sheet, row 1 = "Produto", row 2 = the emission
// line, row 3 = the PT header, rows 4+ = products.
func buildRawSankhyaExport(t *testing.T, sheets []rawSheet) []byte {
	t.Helper()
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	for _, sheet := range sheets {
		if _, err := f.NewSheet(sheet.name); err != nil {
			t.Fatalf("NewSheet(%q): %v", sheet.name, err)
		}
		setRawCell(t, f, sheet.name, 1, 1, "Produto")
		setRawCell(t, f, sheet.name, 1, 2, "Emissão: 20/07/2026 10:00")
		setRawCell(t, f, sheet.name, 2, 2, fmt.Sprintf("Total de registros: %d", sheet.total))
		setRawCell(t, f, sheet.name, 3, 2, "Usuário: DEMO")
		for col, header := range rawHeader {
			setRawCell(t, f, sheet.name, col+1, 3, header)
		}
		for i, product := range sheet.products {
			for col, value := range product {
				if value == "" {
					continue
				}
				setRawCell(t, f, sheet.name, col+1, 4+i, value)
			}
		}
	}
	f.DeleteSheet("Sheet1")

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		t.Fatalf("write workbook: %v", err)
	}
	return buf.Bytes()
}

func setRawCell(t *testing.T, f *excelize.File, sheet string, col, row int, value string) {
	t.Helper()
	name, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		t.Fatalf("cell name (%d,%d): %v", col, row, err)
	}
	if err := f.SetCellValue(sheet, name, value); err != nil {
		t.Fatalf("set cell %s!%s: %v", sheet, name, err)
	}
}

// product builds one raw-export product row aligned to rawHeader.
func product(codigo, descricao, marca1, marca2, ean, ncm, ref, estoqueMinimo string) []string {
	return []string{codigo, descricao, marca1, marca2, ean, ncm, ref, estoqueMinimo}
}

func syntheticSheet(name, prefix string, count int) rawSheet {
	products := make([][]string, 0, count)
	for i := 1; i <= count; i++ {
		code := fmt.Sprintf("%s%04d", prefix, i)
		products = append(products, product(code, name+" produto "+code, "", "", "", "", "REF-"+code, "0"))
	}
	return rawSheet{name: name, total: count, products: products}
}

// TestParserRawExportGolden2012 is the ground-truth golden: the raw 3-sheet
// Sankhya export (FACAS 138 + FERRAMENTAS 1837 + MAQUINAS 37 = 2012 products,
// with 2 preamble rows and PT headers per sheet, no CUSTO/ESTOQUE_FISICO
// columns) imports to exactly 2012 accepted rows on the lenient path with ZERO
// fabricated cost/stock (ADR-17) and the sheet name preserved as category.
func TestParserRawExportGolden2012(t *testing.T) {
	sheets := []rawSheet{
		syntheticSheet("FACAS", "F", 138),
		syntheticSheet("FERRAMENTAS", "G", 1837),
		syntheticSheet("MAQUINAS", "M", 37),
	}
	data := buildRawSankhyaExport(t, sheets)

	rows, fileIssues, err := NewParser().ParseLenient(context.Background(), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseLenient() error = %v", err)
	}
	if len(fileIssues) != 0 {
		t.Fatalf("ParseLenient() fileIssues = %+v, want none (every sheet has a header)", fileIssues)
	}
	if len(rows) != 2012 {
		t.Fatalf("ParseLenient() rows = %d, want 2012 (union of all 3 sheets)", len(rows))
	}

	accepted, issues, rejected := domain.ValidateRowsLenient(rows)
	if len(accepted) != 2012 {
		t.Fatalf("accepted = %d, want 2012", len(accepted))
	}
	if rejected != 0 {
		t.Fatalf("rejected = %d, want 0", rejected)
	}

	// ADR-17: absent CUSTO / ESTOQUE_FISICO are honest-unknown (empty), never a
	// fabricated 0. Not one accepted row may carry a synthesized cost or stock.
	category := map[string]int{}
	for _, row := range accepted {
		if row.Custo != "" {
			t.Fatalf("row %s has fabricated custo %q, want empty", row.Codprod, row.Custo)
		}
		if row.StockPhysical != "" {
			t.Fatalf("row %s has fabricated stock %q, want empty", row.Codprod, row.StockPhysical)
		}
		if row.DescrGrupo != nil {
			category[*row.DescrGrupo]++
		}
	}
	if category["FACAS"] != 138 || category["FERRAMENTAS"] != 1837 || category["MAQUINAS"] != 37 {
		t.Fatalf("category counts = %v, want FACAS 138 / FERRAMENTAS 1837 / MAQUINAS 37", category)
	}

	// Every accepted row must have emitted honest-unknown warnings for the two
	// absent columns (never silent) — 2012 * 2 = 4024 missing-value warnings.
	missingCusto, missingEstoque := 0, 0
	for _, issue := range issues {
		if issue.Kind != domain.Warning {
			continue
		}
		switch issue.Code {
		case domain.CodeMissingCusto:
			missingCusto++
		case domain.CodeMissingEstoque:
			missingEstoque++
		}
	}
	if missingCusto != 2012 || missingEstoque != 2012 {
		t.Fatalf("honest-unknown warnings = custo %d / estoque %d, want 2012 each", missingCusto, missingEstoque)
	}
}

// TestParserRawExportHeaderDetectionSkipsPreamble proves the header-row
// detection: the two preamble rows above the PT header are skipped, not parsed
// as data. Revert guard: reading rows[0] as the header would map "Produto" as
// the only cell and yield garbage / zero identity columns.
func TestParserRawExportHeaderDetectionSkipsPreamble(t *testing.T) {
	data := buildRawSankhyaExport(t, []rawSheet{{
		name:  "FACAS",
		total: 1,
		products: [][]string{
			product("F0001", "Faca do chef", "Tramontina", "", "7891234567895", "82119100", "REF-F0001", "3"),
		},
	}})

	rows, _, err := NewParser().ParseLenient(context.Background(), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseLenient() error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d, want 1 (preamble rows skipped, single product)", len(rows))
	}
	if rows[0].Codprod != "F0001" || rows[0].Descrprod != "Faca do chef" {
		t.Fatalf("preamble not skipped: row[0] = %#v", rows[0])
	}
}

// TestParserRawExportPTAliases proves every Portuguese Sankhya header alias maps
// to the canonical column, including "first non-empty wins" across the twin
// "Marca" columns.
func TestParserRawExportPTAliases(t *testing.T) {
	data := buildRawSankhyaExport(t, []rawSheet{{
		name:  "FERRAMENTAS",
		total: 2,
		products: [][]string{
			// Marca in the SECOND marca column only → first-non-empty picks it.
			product("G0001", "Furadeira", "", "Bosch", "7891234567895", "84672100", "REF-FORN-1", "5"),
			// Marca in the FIRST marca column → wins over empty second.
			product("G0002", "Parafusadeira", "Makita", "", "", "", "REF-FORN-2", "2"),
		},
	}})

	rows, _, err := NewParser().ParseLenient(context.Background(), bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseLenient() error = %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}

	first := rows[0]
	if first.Codprod != "G0001" || first.Descrprod != "Furadeira" {
		t.Fatalf("Código/Descrição alias failed: %#v", first)
	}
	if first.EAN == nil || *first.EAN != "7891234567895" {
		t.Fatalf("Código de Barra→EAN alias failed: %#v", first.EAN)
	}
	if first.NCM == nil || *first.NCM != "84672100" {
		t.Fatalf("NCM alias failed: %#v", first.NCM)
	}
	if first.Refforn == nil || *first.Refforn != "REF-FORN-1" {
		t.Fatalf("Referência do Fornecedor→Refforn alias failed: %#v", first.Refforn)
	}
	if first.Marca == nil || *first.Marca != "Bosch" {
		t.Fatalf("Marca first-non-empty (second column) failed: %#v", first.Marca)
	}
	if rows[1].Marca == nil || *rows[1].Marca != "Makita" {
		t.Fatalf("Marca first-non-empty (first column) failed: %#v", rows[1].Marca)
	}
}

// TestParserRawExportSkipsHeaderlessSheetWarnNotFail proves a sheet with no
// detectable product header is skipped with a WARNING file-issue, and the
// sheets that DO carry data still import — the whole file never hard-fails.
func TestParserRawExportSkipsHeaderlessSheetWarnNotFail(t *testing.T) {
	// Author a good sheet, then inject a second sheet with only junk rows (no
	// CODPROD/DESCRPROD header anywhere).
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	if _, err := f.NewSheet("FACAS"); err != nil {
		t.Fatal(err)
	}
	setRawCell(t, f, "FACAS", 1, 1, "Produto")
	setRawCell(t, f, "FACAS", 1, 2, "Emissão: x")
	for col, header := range rawHeader {
		setRawCell(t, f, "FACAS", col+1, 3, header)
	}
	setRawCell(t, f, "FACAS", 1, 4, "F0001")
	setRawCell(t, f, "FACAS", 2, 4, "Faca")
	if _, err := f.NewSheet("RESUMO"); err != nil {
		t.Fatal(err)
	}
	setRawCell(t, f, "RESUMO", 1, 1, "Relatório de resumo")
	setRawCell(t, f, "RESUMO", 1, 2, "Total geral: 1")
	f.DeleteSheet("Sheet1")

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		t.Fatal(err)
	}

	rows, fileIssues, err := NewParser().ParseLenient(context.Background(), bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("ParseLenient() error = %v, want the good sheet to import (never hard-fail)", err)
	}
	if len(rows) != 1 || rows[0].Codprod != "F0001" {
		t.Fatalf("rows = %#v, want the FACAS product to survive", rows)
	}
	if len(fileIssues) != 1 {
		t.Fatalf("fileIssues = %+v, want exactly one skipped-sheet warning", fileIssues)
	}
	skip := fileIssues[0]
	if skip.Kind != domain.Warning {
		t.Fatalf("skipped-sheet issue kind = %q, want WARNING (never a rejection/hard-fail)", skip.Kind)
	}
	if skip.Column == nil || *skip.Column != "RESUMO" {
		t.Fatalf("skipped-sheet issue column = %#v, want the sheet name RESUMO", skip.Column)
	}
}
