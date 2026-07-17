package domain

import "testing"

func TestValidateRowRejectionCodes(t *testing.T) {
	tests := []struct {
		name string
		row  NormalizedRow
		seen map[string]struct{}
		code string
	}{
		{
			name: "empty codprod",
			row:  validRow("   ", "Produto", "10.00", "2"),
			code: "EMPTY_CODPROD",
		},
		{
			name: "duplicate codprod",
			row:  validRow("P-1", "Produto", "10.00", "2"),
			seen: map[string]struct{}{"P-1": {}},
			code: "DUPLICATE_CODPROD",
		},
		{
			name: "empty descrprod",
			row:  validRow("P-1", "  ", "10.00", "2"),
			code: "EMPTY_DESCRPROD",
		},
		{
			name: "invalid custo",
			row:  validRow("P-1", "Produto", "0", "2"),
			code: "INVALID_CUSTO",
		},
		{
			name: "invalid estoque",
			row:  validRow("P-1", "Produto", "10.00", "-1"),
			code: "INVALID_ESTOQUE",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidateRow(test.row, test.seen)
			if result.Row != nil {
				t.Fatal("rejected row was returned as accepted")
			}
			if len(result.Issues) != 1 {
				t.Fatalf("got %d issues, want 1: %#v", len(result.Issues), result.Issues)
			}
			if result.Issues[0].Kind != Rejection {
				t.Fatalf("got kind %q, want %q", result.Issues[0].Kind, Rejection)
			}
			if string(result.Issues[0].Code) != test.code {
				t.Fatalf("got code %q, want %q", result.Issues[0].Code, test.code)
			}
		})
	}
}

func TestValidateRowWarningsAcceptAndNullInvalidOptionalValues(t *testing.T) {
	invalidEAN := "1234567890123"
	invalidNCM := "1234567"
	result := ValidateRow(NormalizedRow{
		Codprod:       "P-1",
		Descrprod:     "Produto",
		Custo:         "10,50",
		StockPhysical: "2",
		EAN:           &invalidEAN,
		NCM:           &invalidNCM,
	}, nil)

	if result.Row == nil {
		t.Fatal("warning-only row was rejected")
	}
	if len(result.Issues) != 2 {
		t.Fatalf("got %d issues, want 2: %#v", len(result.Issues), result.Issues)
	}
	if result.Issues[0].Kind != Warning || string(result.Issues[0].Code) != "INVALID_EAN" {
		t.Fatalf("unexpected EAN issue: %#v", result.Issues[0])
	}
	if result.Issues[1].Kind != Warning || string(result.Issues[1].Code) != "INVALID_NCM" {
		t.Fatalf("unexpected NCM issue: %#v", result.Issues[1])
	}
	if result.Row.EAN != nil || result.Row.NCM != nil {
		t.Fatalf("invalid optional values were not nulled: %#v", result.Row)
	}
	if result.Row.Custo != "10.50" {
		t.Fatalf("got custo %q, want exact canonical decimal 10.50", result.Row.Custo)
	}
}

func TestValidateRowCanonicalizesBareDecimalCusto(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{".5", "0.5"},
		{",50", "0.50"},
		{"+.75", "0.75"},
	} {
		result := ValidateRow(validRow("P-1", "Produto", tc.in, "1"), nil)
		if result.Row == nil {
			t.Fatalf("%q rejected, want accepted", tc.in)
		}
		if string(result.Row.Custo) != tc.want {
			t.Fatalf("custo %q normalized to %q, want canonical %q", tc.in, result.Row.Custo, tc.want)
		}
	}
}

func TestValidateRowsRejectsDuplicateAndKeepsCleanRows(t *testing.T) {
	accepted, issues, rejected := ValidateRows([]NormalizedRow{
		validRow("P-1", "Produto 1", "10.00", "0"),
		validRow("P-1", "Produto 1 duplicate", "11.00", "1"),
		validRow("P-2", "Produto 2", "12.00", "2"),
	})

	if rejected != 1 {
		t.Fatalf("got %d rejected rows, want 1", rejected)
	}
	if len(accepted) != 2 {
		t.Fatalf("got %d accepted rows, want 2", len(accepted))
	}
	if len(issues) != 1 || string(issues[0].Code) != "DUPLICATE_CODPROD" || issues[0].Row != 2 {
		t.Fatalf("unexpected duplicate issue: %#v", issues)
	}
}

func TestValidateRowCleanProducesZeroIssues(t *testing.T) {
	result := ValidateRow(validRow("P-1", "Produto", "10.25", "3"), nil)
	if result.Row == nil {
		t.Fatal("clean row was rejected")
	}
	if len(result.Issues) != 0 {
		t.Fatalf("got issues for clean row: %#v", result.Issues)
	}
	if result.Row.StockPhysical != "3" {
		t.Fatalf("got stock physical %q, want 3", result.Row.StockPhysical)
	}
}

func validRow(codprod, descrprod, custo, stockPhysical string) NormalizedRow {
	return NormalizedRow{
		Codprod:       codprod,
		Descrprod:     descrprod,
		Custo:         Decimal(custo),
		StockPhysical: stockPhysical,
	}
}
