package domain

import (
	"math/big"
	"reflect"
	"testing"
)

// strptr is a test helper for the nullable (*string) decomposition fields.
func strptr(s string) *string { return &s }

func money(amount string) *Money { return &Money{Amount: amount, Currency: "BRL"} }

// TestDecomposeGolden pins the SINGLE IC-04 decomposition formula (line 84)
// against hand-computed exact-decimal vectors. Every component is 2dp;
// margem_valor = preço − Σ(componentes) holds EXACTLY (soma fecha); any
// unknown component ⇒ margem_valor/margem_pct nil + componentes_desconhecidos
// naming the missing (NEVER 0). comissao_pct/aliquota/efetivo are percentages.
func TestDecomposeGolden(t *testing.T) {
	// Fixed resolved inputs the ports (S6) supply; chosen per-vector below.
	// comissão pct by modalidade (test values, not contract constants):
	//   classico 12 · premium 14 · full 16. efetivo: SP 6.00 · BA 13.50.
	cases := []struct {
		name string
		in   DecomposeInput
		want Decomposition
	}{
		{
			name: "V1 SIMPLES 150 classico SP custo-known — margin positive",
			in: DecomposeInput{
				Preco: "150.00", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("40.00"),
				FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "150.00", Comissao: "18.00", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("6.00"), Difal: strptr("9.00"),
				TarifaFull: strptr("0.00"), Custo: strptr("40.00"),
				MargemValor: strptr("62.00"), MargemPct: strptr("41.33"),
				ComponentesDesconhecidos: nil,
			},
		},
		{
			name: "V2 SIMPLES 78.90 classico SP custo-known — under 79: taxa_fixa 6.50, frete 0",
			in: DecomposeInput{
				Preco: "78.90", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("40.00"),
				FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "78.90", Comissao: "9.47", TaxaFixa: "6.50",
				Frete: strptr("0.00"), Imposto: strptr("3.16"), Difal: strptr("4.73"),
				TarifaFull: strptr("0.00"), Custo: strptr("40.00"),
				MargemValor: strptr("15.04"), MargemPct: strptr("19.06"),
			},
		},
		{
			name: "V3 PRESUMIDO 79.00 premium BA custo-known — negative margin",
			in: DecomposeInput{
				Preco: "79.00", ComissaoPct: "14", AliquotaPct: "9.25",
				Modalidade: ModalidadePremium, Custo: money("40.00"),
				FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "BA", EfetivoPct: "13.50",
			},
			want: Decomposition{
				Preco: "79.00", Comissao: "11.06", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("7.31"), Difal: strptr("10.67"),
				TarifaFull: strptr("0.00"), Custo: strptr("40.00"),
				MargemValor: strptr("-5.04"), MargemPct: strptr("-6.38"),
			},
		},
		{
			name: "V4 SIMPLES 150 full tarifa_full SET SP custo-known",
			in: DecomposeInput{
				Preco: "150.00", ComissaoPct: "16", AliquotaPct: "4",
				Modalidade: ModalidadeFull, TarifaFull: money("8.00"),
				Custo: money("40.00"), FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "150.00", Comissao: "24.00", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("6.00"), Difal: strptr("9.00"),
				TarifaFull: strptr("8.00"), Custo: strptr("40.00"),
				MargemValor: strptr("48.00"), MargemPct: strptr("32.00"),
			},
		},
		{
			name: "V5 full tarifa_full NULL ⇒ tarifa_full UNKNOWN (never 0)",
			in: DecomposeInput{
				Preco: "150.00", ComissaoPct: "16", AliquotaPct: "4",
				Modalidade: ModalidadeFull, TarifaFull: nil,
				Custo: money("40.00"), FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "150.00", Comissao: "24.00", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("6.00"), Difal: strptr("9.00"),
				TarifaFull: nil, Custo: strptr("40.00"),
				MargemValor: nil, MargemPct: nil,
				ComponentesDesconhecidos: []string{"tarifa_full"},
			},
		},
		{
			name: "V6 custo NIL ⇒ custo_erp UNKNOWN (never 0)",
			in: DecomposeInput{
				Preco: "150.00", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: nil,
				FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "150.00", Comissao: "18.00", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("6.00"), Difal: strptr("9.00"),
				TarifaFull: strptr("0.00"), Custo: nil,
				MargemValor: nil, MargemPct: nil,
				ComponentesDesconhecidos: []string{"custo_erp"},
			},
		},
		{
			name: "V7 78.90 frete_produto NIL but under 79 ⇒ frete 0 KNOWN (nil safe)",
			in: DecomposeInput{
				Preco: "78.90", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("40.00"),
				FreteProduto: nil,
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "78.90", Comissao: "9.47", TaxaFixa: "6.50",
				Frete: strptr("0.00"), Imposto: strptr("3.16"), Difal: strptr("4.73"),
				TarifaFull: strptr("0.00"), Custo: strptr("40.00"),
				MargemValor: strptr("15.04"), MargemPct: strptr("19.06"),
			},
		},
		{
			name: "V8 150 frete_produto NIL and >=79 ⇒ frete UNKNOWN",
			in: DecomposeInput{
				Preco: "150.00", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("40.00"),
				FreteProduto: nil,
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "150.00", Comissao: "18.00", TaxaFixa: "0.00",
				Frete: nil, Imposto: strptr("6.00"), Difal: strptr("9.00"),
				TarifaFull: strptr("0.00"), Custo: strptr("40.00"),
				MargemValor: nil, MargemPct: nil,
				ComponentesDesconhecidos: []string{"frete"},
			},
		},
		{
			name: "V9 difal_enabled FALSE ⇒ difal 0.00 explicit (legit zero, not unknown)",
			in: DecomposeInput{
				Preco: "150.00", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("40.00"),
				FreteProduto: money("15.00"),
				DifalEnabled: false, DestinoUF: "", EfetivoPct: "",
			},
			want: Decomposition{
				Preco: "150.00", Comissao: "18.00", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("6.00"), Difal: strptr("0.00"),
				TarifaFull: strptr("0.00"), Custo: strptr("40.00"),
				MargemValor: strptr("71.00"), MargemPct: strptr("47.33"),
			},
		},
		{
			name: "V10 difal_enabled TRUE but destino UNKNOWN ⇒ difal UNKNOWN",
			in: DecomposeInput{
				Preco: "150.00", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("40.00"),
				FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "", EfetivoPct: "",
			},
			want: Decomposition{
				Preco: "150.00", Comissao: "18.00", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("6.00"), Difal: nil,
				TarifaFull: strptr("0.00"), Custo: strptr("40.00"),
				MargemValor: nil, MargemPct: nil,
				ComponentesDesconhecidos: []string{"difal"},
			},
		},
		{
			name: "V11 full+null AND custo nil ⇒ desconhecidos ordered [tarifa_full, custo_erp]",
			in: DecomposeInput{
				Preco: "150.00", ComissaoPct: "16", AliquotaPct: "4",
				Modalidade: ModalidadeFull, TarifaFull: nil, Custo: nil,
				FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "150.00", Comissao: "24.00", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("6.00"), Difal: strptr("9.00"),
				TarifaFull: nil, Custo: nil,
				MargemValor: nil, MargemPct: nil,
				ComponentesDesconhecidos: []string{"tarifa_full", "custo_erp"},
			},
		},
		{
			name: "V12 PRESUMIDO 79.00 classico SP custo-known — boundary",
			in: DecomposeInput{
				Preco: "79.00", ComissaoPct: "12", AliquotaPct: "9.25",
				Modalidade: ModalidadeClassico, Custo: money("40.00"),
				FreteProduto: money("15.00"),
				DifalEnabled: true, DestinoUF: "SP", EfetivoPct: "6.00",
			},
			want: Decomposition{
				Preco: "79.00", Comissao: "9.48", TaxaFixa: "0.00",
				Frete: strptr("15.00"), Imposto: strptr("7.31"), Difal: strptr("4.74"),
				TarifaFull: strptr("0.00"), Custo: strptr("40.00"),
				MargemValor: strptr("2.47"), MargemPct: strptr("3.13"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Decompose(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Decompose mismatch\n got: %+v\nwant: %+v", got, tc.want)
			}
			// soma fecha: preço = Σ componentes + margem_valor (only when known).
			if got.MargemValor != nil {
				assertSomaFecha(t, got)
			}
		})
	}
}

// assertSomaFecha verifies preço = comissao+taxa_fixa+frete+imposto+difal+
// tarifa_full+custo + margem_valor EXACTLY at 2dp (all components known).
func assertSomaFecha(t *testing.T, d Decomposition) {
	t.Helper()
	sum := ratOf(t, d.Comissao)
	sum.Add(sum, ratOf(t, d.TaxaFixa))
	sum.Add(sum, ratOf(t, *d.Frete))
	sum.Add(sum, ratOf(t, *d.Imposto))
	sum.Add(sum, ratOf(t, *d.Difal))
	sum.Add(sum, ratOf(t, *d.TarifaFull))
	sum.Add(sum, ratOf(t, *d.Custo))
	sum.Add(sum, ratOf(t, *d.MargemValor))
	if FormatRatHalfUp(sum, 2) != d.Preco {
		t.Fatalf("soma não fecha: Σ+margem = %s, preço = %s", FormatRatHalfUp(sum, 2), d.Preco)
	}
}

func ratOf(t *testing.T, s string) *big.Rat {
	t.Helper()
	r, err := ParseRat(s)
	if err != nil {
		t.Fatalf("ParseRat(%q): %v", s, err)
	}
	return r
}
