package xlsx

import (
	"bytes"
	"context"
	"testing"

	"github.com/xuri/excelize/v2"
	"marketplace-central/apps/server_core/internal/modules/erp_import/domain"
)

// buildPriceWorkbook authors a single-sheet workbook whose header row is given
// verbatim, so each case can exercise a different sale-price header label.
func buildPriceWorkbook(t *testing.T, header []string, rows [][]string) []byte {
	t.Helper()
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()

	for col, value := range header {
		setRawCell(t, f, "Sheet1", col+1, 1, value)
	}
	for i, row := range rows {
		for col, value := range row {
			if value == "" {
				continue
			}
			setRawCell(t, f, "Sheet1", col+1, 2+i, value)
		}
	}

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		t.Fatalf("write workbook: %v", err)
	}
	return buf.Bytes()
}

func TestParseReadsSalePriceColumnAliases(t *testing.T) {
	cases := []struct {
		label string
		want  domain.Decimal
	}{
		{label: "PRECO_VENDA", want: "199.90"},
		{label: "Preço de Venda", want: "199.90"},
		{label: "VLRVENDA", want: "199.90"},
	}

	for _, testCase := range cases {
		t.Run(testCase.label, func(t *testing.T) {
			header := []string{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO", testCase.label}
			file := buildPriceWorkbook(t, header, [][]string{
				{"10010", "CHAVE BIELA 25B-8", "381.97", "96", "199,90"},
			})

			rows, _, err := (&Parser{}).Parse(context.Background(), bytes.NewReader(file))
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if len(rows) != 1 {
				t.Fatalf("rows=%d want 1", len(rows))
			}
			if rows[0].PrecoVenda != "199,90" {
				t.Fatalf("parsed preco_venda=%q want raw %q", rows[0].PrecoVenda, "199,90")
			}

			accepted, _, rejected := domain.ValidateRows(rows)
			if rejected != 0 {
				t.Fatalf("rejected=%d want 0", rejected)
			}
			if accepted[0].PrecoVenda != testCase.want {
				t.Fatalf("validated preco_venda=%q want %q", accepted[0].PrecoVenda, testCase.want)
			}
		})
	}
}

// A workbook with no sale-price column stays honest-unknown: empty value, no
// rejection, no fabricated zero (ADR-17).
func TestParseWithoutSalePriceColumnKeepsItEmpty(t *testing.T) {
	file := buildPriceWorkbook(t,
		[]string{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO"},
		[][]string{{"10010", "CHAVE BIELA 25B-8", "381.97", "96"}})

	rows, _, err := (&Parser{}).Parse(context.Background(), bytes.NewReader(file))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	accepted, _, rejected := domain.ValidateRows(rows)
	if rejected != 0 {
		t.Fatalf("rejected=%d want 0", rejected)
	}
	if accepted[0].PrecoVenda != "" {
		t.Fatalf("preco_venda=%q want empty", accepted[0].PrecoVenda)
	}
}

// An unparseable sale price is a WARNING and is cleared, never a rejection and
// never a fabricated number.
func TestUnparseableSalePriceWarnsAndClears(t *testing.T) {
	file := buildPriceWorkbook(t,
		[]string{"CODPROD", "DESCRPROD", "CUSTO", "ESTOQUE_FISICO", "PRECO_VENDA"},
		[][]string{{"10010", "CHAVE BIELA 25B-8", "381.97", "96", "sob consulta"}})

	rows, _, err := (&Parser{}).Parse(context.Background(), bytes.NewReader(file))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	accepted, issues, rejected := domain.ValidateRows(rows)
	if rejected != 0 {
		t.Fatalf("rejected=%d want 0", rejected)
	}
	if accepted[0].PrecoVenda != "" {
		t.Fatalf("preco_venda=%q want cleared", accepted[0].PrecoVenda)
	}
	var warned bool
	for _, current := range issues {
		if current.Column != nil && *current.Column == "PRECO_VENDA" && current.Kind == domain.Warning {
			warned = true
		}
	}
	if !warned {
		t.Fatal("expected a PRECO_VENDA warning issue")
	}
}
