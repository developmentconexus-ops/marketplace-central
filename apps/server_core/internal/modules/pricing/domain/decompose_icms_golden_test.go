package domain

import (
	"math/big"
	"reflect"
	"testing"
)

// TestDecomposeICMSGolden pins Decompose's D-41 path (ICMSCell present) end
// to end: Imposto stays computed/shown but stops being summed; ICMSSaida,
// Difal (now sourced from TaxesForItem instead of DifalEnabled/DestinoUF/
// EfetivoPct), PisCofins and RestituicaoST enter the sum in its place.
// RestituicaoST is the one component that SUBTRACTS from the running sum
// (ADDS to margem) — Task 6 calls it "a única linha positiva da seção".
// Every expected value is hand-computed with exact big.Rat fractions
// (documented inline); the mandated four (interno, interestadual, ST,
// ambígua) plus a MAX(0) vector and a dedicated restituição-increases-margem
// pair.
func TestDecomposeICMSGolden(t *testing.T) {
	cases := []struct {
		name string
		in   DecomposeInput
		want Decomposition
	}{
		{
			// MG interno: a fixo 0,18 (não consulta AliquotaInterna).
			// Comissao=1000x0.12=120.00. TaxaFixa=0 (>=79). Frete=30.00.
			// Imposto=1000x0.04=40.00 (mostrado, NAO somado).
			// ICMSSaida=70.00 Difal=110.00 PisCofins=75.85 RestituicaoST=0.00 (UF=MG).
			// Sum = 120+0+30+70+110+75.85+0(tarifa)+200(custo) - 0(restituicao) = 605.85.
			// Margem = 1000-605.85 = 394.15. Pct = 394.15/1000*100 = 39.415 -> 39.42 (tie, half-up).
			name: "ICMS-1: interno MG — a fixo, restituicao 0 explicito (UF=MG)",
			in: DecomposeInput{
				Preco: "1000.00", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("200.00"),
				FreteProduto: money("30.00"),
				ICMSCell: &ICMSCell{
					UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
					Origprod: intp(0), RestituicaoUnit: strptr("50.00"),
				},
			},
			want: Decomposition{
				Preco: "1000.00", Comissao: "120.00", TaxaFixa: "0.00",
				Frete: strptr("30.00"), Imposto: "40.00",
				Difal: strptr("110.00"), TarifaFull: strptr("0.00"), Custo: strptr("200.00"),
				ICMSSaida: strptr("70.00"), PisCofins: strptr("75.85"), RestituicaoST: strptr("0.00"),
				MargemValor: strptr("394.15"), MargemPct: strptr("39.42"),
			},
		},
		{
			// Interestadual SP, gross-up pela tabela legal (a_int=18%, a_inter=12% pq SP no
			// conjunto dos 12%). Comissao=500x0.10=50.00. Frete=20.00. Imposto=500x0.04=20.00
			// (mostrado, nao somado). ICMSSaida=60.00 Difal=36.59 PisCofins=32.69
			// RestituicaoST=25.00 (UF!=MG).
			// Sum = 50+0+20+60+36.59+32.69+0+150 - 25 = 324.28.
			// Margem = 500-324.28 = 175.72. Pct = 175.72/500*100 = 35.144 -> 35.14.
			name: "ICMS-2: interestadual SP — gross-up da tabela legal, restituicao subtrai do sum",
			in: DecomposeInput{
				Preco: "500.00", ComissaoPct: "10", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("150.00"),
				FreteProduto: money("20.00"),
				ICMSCell: &ICMSCell{
					UFDestino: "SP", CodTrib: intp(0), Ambiguo: false,
					AliquotaInterna: strptr("18"), Origprod: intp(0),
					StRetidoEntrada: strptr("50.00"), RestituicaoUnit: strptr("25.00"),
				},
			},
			want: Decomposition{
				Preco: "500.00", Comissao: "50.00", TaxaFixa: "0.00",
				Frete: strptr("20.00"), Imposto: "20.00",
				Difal: strptr("36.59"), TarifaFull: strptr("0.00"), Custo: strptr("150.00"),
				ICMSSaida: strptr("60.00"), PisCofins: strptr("32.69"), RestituicaoST: strptr("25.00"),
				MargemValor: strptr("175.72"), MargemPct: strptr("35.14"),
			},
		},
		{
			// ST no destino (RJ, codtrib=60): ICMSSaida/Difal = 0.00 explicito. PisCofins
			// deduz S do liquido (a=0): 800-40=760, x0.0925=70.30. RestituicaoST=15.00 (UF!=MG,
			// independente do ramo ST). Comissao=800x0.11=88.00. Frete=25.00.
			// Imposto=800x0.04=32.00 (mostrado, nao somado).
			// Sum = 88+0+25+0+0+70.30+0+300 - 15 = 468.30.
			// Margem = 800-468.30 = 331.70. Pct = 331.70/800*100 = 41.4625 -> 41.46.
			name: "ICMS-3: ST no destino — ICMS/DIFAL zero explicito, PIS/COFINS deduz S",
			in: DecomposeInput{
				Preco: "800.00", ComissaoPct: "11", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("300.00"),
				FreteProduto: money("25.00"),
				ICMSCell: &ICMSCell{
					UFDestino: "RJ", CodTrib: intp(60), Ambiguo: false,
					StRetidoEntrada: strptr("40.00"), RestituicaoUnit: strptr("15.00"),
				},
			},
			want: Decomposition{
				Preco: "800.00", Comissao: "88.00", TaxaFixa: "0.00",
				Frete: strptr("25.00"), Imposto: "32.00",
				Difal: strptr("0.00"), TarifaFull: strptr("0.00"), Custo: strptr("300.00"),
				ICMSSaida: strptr("0.00"), PisCofins: strptr("70.30"), RestituicaoST: strptr("15.00"),
				MargemValor: strptr("331.70"), MargemPct: strptr("41.46"),
			},
		},
		{
			// Celula ambigua (PR): ICMSSaida/Difal/PisCofins ficam nil e nomeados —
			// RestituicaoST continua conhecida (10.00, independente da resolucao da celula).
			// Margem/Pct ficam nil (ADR-17: componente desconhecido nunca vira 0). Imposto
			// continua mostrado (600x0.04=24.00) mesmo com margem indisponivel.
			name: "ICMS-4: ambiguo=true — icms/difal/pis_cofins unknown, restituicao continua conhecida",
			in: DecomposeInput{
				Preco: "600.00", ComissaoPct: "10", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("200.00"),
				FreteProduto: money("10.00"),
				ICMSCell: &ICMSCell{
					UFDestino: "PR", CodTrib: intp(0), Ambiguo: true,
					AliquotaInterna: strptr("19.5"), Origprod: intp(5),
					RestituicaoUnit: strptr("10.00"),
				},
			},
			want: Decomposition{
				Preco: "600.00", Comissao: "60.00", TaxaFixa: "0.00",
				Frete: strptr("10.00"), Imposto: "24.00",
				Difal: nil, TarifaFull: strptr("0.00"), Custo: strptr("200.00"),
				ICMSSaida: nil, PisCofins: nil, RestituicaoST: strptr("10.00"),
				MargemValor: nil, MargemPct: nil,
				ComponentesDesconhecidos: []string{"icms_saida", "difal", "pis_cofins"},
			},
		},
		{
			// MAX(0,...) nao e cosmetico: S=500 excede o liquido em MG (P=200, a=0.18 fixo).
			// P*(1-a)=164, 164-500=-336 -> clamp a 0 -> PisCofins=0.00 CONHECIDO (zero
			// legitimo por clamp, nao desconhecido). ICMSSaida=200x0.07=14.00 (a_inter, MG fora
			// do conjunto dos 12%). Difal=200x0.18-14=22.00. RestituicaoST=0.00 (UF=MG).
			// Comissao=200x0.10=20.00. Frete=10.00. Imposto=200x0.04=8.00 (mostrado, nao somado).
			// Sum = 20+0+10+14+22+0+0+50 - 0 = 116.00. Margem=200-116=84.00. Pct=42.00 exato.
			name: "ICMS-5: MAX(0) clamp na base do PIS/COFINS — S excede o liquido",
			in: DecomposeInput{
				Preco: "200.00", ComissaoPct: "10", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("50.00"),
				FreteProduto: money("10.00"),
				ICMSCell: &ICMSCell{
					UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
					Origprod: intp(0), StRetidoEntrada: strptr("500.00"),
				},
			},
			want: Decomposition{
				Preco: "200.00", Comissao: "20.00", TaxaFixa: "0.00",
				Frete: strptr("10.00"), Imposto: "8.00",
				Difal: strptr("22.00"), TarifaFull: strptr("0.00"), Custo: strptr("50.00"),
				ICMSSaida: strptr("14.00"), PisCofins: strptr("0.00"), RestituicaoST: strptr("0.00"),
				MargemValor: strptr("84.00"), MargemPct: strptr("42.00"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Decompose(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("Decompose mismatch\n got: %+v\nwant: %+v", got, tc.want)
			}
			if got.MargemValor != nil {
				assertSomaFechaICMS(t, got)
			}
		})
	}
}

// assertSomaFechaICMS is assertSomaFecha's ICMS-cell-aware sibling: Imposto
// is EXCLUDED (ICMSSaida already carries that tax) and RestituicaoST
// SUBTRACTS (it's a credit, the one component that grows margem).
func assertSomaFechaICMS(t *testing.T, d Decomposition) {
	t.Helper()
	sum := ratOf(t, d.Comissao)
	sum.Add(sum, ratOf(t, d.TaxaFixa))
	sum.Add(sum, ratOf(t, *d.Frete))
	sum.Add(sum, ratOf(t, *d.ICMSSaida))
	sum.Add(sum, ratOf(t, *d.Difal))
	sum.Add(sum, ratOf(t, *d.PisCofins))
	sum.Add(sum, ratOf(t, *d.TarifaFull))
	sum.Add(sum, ratOf(t, *d.Custo))
	sum.Sub(sum, ratOf(t, *d.RestituicaoST))
	sum.Add(sum, ratOf(t, *d.MargemValor))
	if FormatRatHalfUp(sum, 2) != d.Preco {
		t.Fatalf("soma nao fecha (excluindo Imposto, restituicao subtraida): Sum+margem = %s, preco = %s",
			FormatRatHalfUp(sum, 2), d.Preco)
	}
}

// TestDecomposeICMSRestituicaoIncreasesMargem is the T4-mandated dedicated
// vector: same item, same everything, only RestituicaoUnit differs (0.00 vs
// 40.00) — proves RestituicaoST is subtracted from the sum (added to
// margem) dollar-for-dollar, not the other way around. If the sign were
// flipped, margem_R40 would be LOWER than margem_R0, not higher.
func TestDecomposeICMSRestituicaoIncreasesMargem(t *testing.T) {
	base := func(restituicao string) DecomposeInput {
		return DecomposeInput{
			Preco: "400.00", ComissaoPct: "10", AliquotaPct: "4",
			Modalidade: ModalidadeClassico, Custo: money("100.00"),
			FreteProduto: money("15.00"),
			ICMSCell: &ICMSCell{
				UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
				AliquotaInterna: strptr("20.5"), Origprod: intp(0),
				RestituicaoUnit: strptr(restituicao),
			},
		}
	}

	withoutR := Decompose(base("0.00"))
	withR := Decompose(base("40.00"))

	if withoutR.RestituicaoST == nil || *withoutR.RestituicaoST != "0.00" {
		t.Fatalf("withoutR.RestituicaoST = %v, want 0.00", withoutR.RestituicaoST)
	}
	if withR.RestituicaoST == nil || *withR.RestituicaoST != "40.00" {
		t.Fatalf("withR.RestituicaoST = %v, want 40.00", withR.RestituicaoST)
	}
	if withoutR.MargemValor == nil || withR.MargemValor == nil {
		t.Fatalf("margem should be known in both cases: withoutR=%v withR=%v", withoutR.MargemValor, withR.MargemValor)
	}
	if *withoutR.MargemValor != "120.95" {
		t.Fatalf("withoutR.MargemValor = %s, want 120.95", *withoutR.MargemValor)
	}
	if *withR.MargemValor != "160.95" {
		t.Fatalf("withR.MargemValor = %s, want 160.95 (120.95 + 40.00 restituicao)", *withR.MargemValor)
	}

	diff := new(big.Rat).Sub(ratOf(t, *withR.MargemValor), ratOf(t, *withoutR.MargemValor))
	if FormatRatHalfUp(diff, 2) != "40.00" {
		t.Fatalf("margem delta = %s, want exactly 40.00 (restituicao adds dollar-for-dollar to margem)", FormatRatHalfUp(diff, 2))
	}

	assertSomaFechaICMS(t, withoutR)
	assertSomaFechaICMS(t, withR)
}
