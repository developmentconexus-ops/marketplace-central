package domain

import (
	"math/big"
	"reflect"
	"testing"
)

// TestDecomposeICMSGolden pins Decompose's D-41 path (ICMSCell present) end
// to end: Imposto is the named-absent nil (Task A6 — the legacy AliquotaPct
// path does not run at all when ICMSCell is present, so it is never parsed
// and never shown); ICMSSaida, Difal (now sourced from TaxesForItem instead
// of DifalEnabled/DestinoUF/EfetivoPct), PisCofins and RestituicaoST enter
// the sum in Imposto's former place. RestituicaoST is the one component
// that SUBTRACTS from the running sum (ADDS to margem) — Task 6 calls it "a
// única linha positiva da seção". Every expected value is hand-computed with
// exact big.Rat fractions (documented inline); the mandated four (interno,
// interestadual, ST, ambígua) plus a MAX(0) vector and a dedicated
// restituição-increases-margem pair.
func TestDecomposeICMSGolden(t *testing.T) {
	cases := []struct {
		name string
		in   DecomposeInput
		want Decomposition
	}{
		{
			// MG intra-UF: Task A2 rule (c)+(d) — MG lê a tabela (18%,
			// fcp_embutido = "0" explicito ⇒ fcp = 0) e ICMS usa a alíquota interna INTEIRA;
			// DIFAL é 0,00 explícito (destino=origem). ICMS+DIFAL total
			// (180,00) é o MESMO valor de antes (70+110=180) — a correção
			// só redistribui entre os dois campos, então Sum/Margem/Pct NÃO
			// mudam.
			// Comissao=1000x0.12=120.00. TaxaFixa=0 (>=79). Frete=30.00.
			// Imposto: ICMSCell present ⇒ legacy AliquotaPct=4 is never
			// parsed/shown — named-absent nil (Task A6).
			// ICMSSaida=180.00 Difal=0.00 PisCofins=75.85 RestituicaoST=0.00 (UF=MG).
			// Sum = 120+0+30+180+0+75.85+0(tarifa)+200(custo) - 0(restituicao) = 605.85.
			// Margem = 1000-605.85 = 394.15. Pct = 394.15/1000*100 = 39.415 -> 39.42 (tie, half-up).
			name: "ICMS-1: interno MG — ICMS pela alíquota interna inteira, DIFAL 0 explícito",
			in: DecomposeInput{
				Preco: "1000.00", ComissaoPct: "12", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("200.00"),
				FreteProduto: money("30.00"),
				ICMSCell: &ICMSCell{
					UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
					Origprod: intp(0), AliquotaInterna: strptr("18"),
					// Task A4 (D2): par completo -- "0" explicito, nao nil
					// (ver icms.go:FcpEmbutido / icms_test.go caso A).
					FcpEmbutido:     strptr("0"),
					RestituicaoUnit: strptr("50.00"),
				},
			},
			want: Decomposition{
				Preco: "1000.00", Comissao: "120.00", TaxaFixa: "0.00",
				Frete: strptr("30.00"), Imposto: nil,
				Difal: strptr("0.00"), TarifaFull: strptr("0.00"), Custo: strptr("200.00"),
				ICMSSaida: strptr("180.00"), PisCofins: strptr("75.85"), RestituicaoST: strptr("0.00"),
				MargemValor: strptr("394.15"), MargemPct: strptr("39.42"),
			},
		},
		{
			// Interestadual SP, diferença simples pela tabela legal (Task A2 rule (a):
			// a_int=18%, fcp_embutido = "0" explicito ⇒ fcp = 0, a_inter=12% pq SP no conjunto dos 12%).
			// ICMS_oper=500x0.12=60.00. ICMS_total=500x0.18=90.00. DIFAL=90-60=30.00.
			// Comissao=500x0.10=50.00. Frete=20.00. Imposto: legacy
			// AliquotaPct=4 never parsed/shown (ICMSCell present, Task A6).
			// ICMSSaida=60.00 Difal=30.00 PisCofins=33.30
			// (base=500x0.82-50=360.00, x0.0925=33.30) RestituicaoST=25.00 (UF!=MG).
			// Sum = 50+0+20+60+30.00+33.30+0+150 - 25 = 318.30.
			// Margem = 500-318.30 = 181.70. Pct = 181.70/500*100 = 36.34.
			name: "ICMS-2: interestadual SP — diferença simples da tabela legal, restituicao subtrai do sum",
			in: DecomposeInput{
				Preco: "500.00", ComissaoPct: "10", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("150.00"),
				FreteProduto: money("20.00"),
				ICMSCell: &ICMSCell{
					UFDestino: "SP", CodTrib: intp(0), Ambiguo: false,
					AliquotaInterna: strptr("18"), Origprod: intp(0),
					// Task A4 (D2): par completo -- ver comentario do ICMS-1.
					FcpEmbutido:     strptr("0"),
					StRetidoEntrada: strptr("50.00"), RestituicaoUnit: strptr("25.00"),
				},
			},
			want: Decomposition{
				Preco: "500.00", Comissao: "50.00", TaxaFixa: "0.00",
				Frete: strptr("20.00"), Imposto: nil,
				Difal: strptr("30.00"), TarifaFull: strptr("0.00"), Custo: strptr("150.00"),
				ICMSSaida: strptr("60.00"), PisCofins: strptr("33.30"), RestituicaoST: strptr("25.00"),
				MargemValor: strptr("181.70"), MargemPct: strptr("36.34"),
			},
		},
		{
			// ST no destino (RJ, codtrib=60): ICMSSaida/Difal = 0.00 explicito. PisCofins
			// deduz S do liquido (a=0): 800-40=760, x0.0925=70.30. RestituicaoST=15.00 (UF!=MG,
			// independente do ramo ST). Comissao=800x0.11=88.00. Frete=25.00.
			// Imposto: legacy AliquotaPct=4 never parsed/shown (ICMSCell
			// present, Task A6).
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
				Frete: strptr("25.00"), Imposto: nil,
				Difal: strptr("0.00"), TarifaFull: strptr("0.00"), Custo: strptr("300.00"),
				ICMSSaida: strptr("0.00"), PisCofins: strptr("70.30"), RestituicaoST: strptr("15.00"),
				MargemValor: strptr("331.70"), MargemPct: strptr("41.46"),
			},
		},
		{
			// Celula ambigua (PR): ICMSSaida/Difal/PisCofins ficam nil e nomeados —
			// RestituicaoST continua conhecida (10.00, independente da resolucao da celula).
			// Margem/Pct ficam nil (ADR-17: componente desconhecido nunca vira 0). Imposto
			// e nil (named-absent, ICMSCell present ⇒ legacy nao roda) — mas NAO entra em
			// ComponentesDesconhecidos: sua ausencia e estrutural (Task A6), nao um fato
			// fiscal que faltou (que ja sao icms_saida/difal/pis_cofins acima).
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
				Frete: strptr("10.00"), Imposto: nil,
				Difal: nil, TarifaFull: strptr("0.00"), Custo: strptr("200.00"),
				ICMSSaida: nil, PisCofins: nil, RestituicaoST: strptr("10.00"),
				MargemValor: nil, MargemPct: nil,
				ComponentesDesconhecidos: []string{"icms_saida", "difal", "pis_cofins"},
			},
		},
		{
			// MAX(0,...) nao e cosmetico: S=500 excede o liquido em MG intra-UF
			// (P=200, alíquota interna 18%, fcp_embutido = "0" explicito ⇒ fcp = 0). aBase=0.18.
			// P*(1-aBase)=164, 164-500=-336 -> clamp a 0 -> PisCofins=0.00 CONHECIDO
			// (zero legitimo por clamp, nao desconhecido). ICMSSaida=200x0.18=36.00
			// (intra-MG, alíquota inteira). Difal=0.00 explícito (intra-MG).
			// ICMS+DIFAL total (36,00) é o MESMO de antes (14+22=36) — a correção só
			// redistribui, Sum/Margem/Pct nao mudam. RestituicaoST=0.00 (UF=MG).
			// Comissao=200x0.10=20.00. Frete=10.00. Imposto: legacy AliquotaPct=4
			// never parsed/shown (ICMSCell present, Task A6).
			// Sum = 20+0+10+36+0+0+0+50 - 0 = 116.00. Margem=200-116=84.00. Pct=42.00 exato.
			name: "ICMS-5: MAX(0) clamp na base do PIS/COFINS — S excede o liquido",
			in: DecomposeInput{
				Preco: "200.00", ComissaoPct: "10", AliquotaPct: "4",
				Modalidade: ModalidadeClassico, Custo: money("50.00"),
				FreteProduto: money("10.00"),
				ICMSCell: &ICMSCell{
					UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
					Origprod: intp(0), AliquotaInterna: strptr("18"),
					// Task A4 (D2): par completo -- ver comentario do ICMS-1.
					FcpEmbutido:     strptr("0"),
					StRetidoEntrada: strptr("500.00"),
				},
			},
			want: Decomposition{
				Preco: "200.00", Comissao: "20.00", TaxaFixa: "0.00",
				Frete: strptr("10.00"), Imposto: nil,
				Difal: strptr("0.00"), TarifaFull: strptr("0.00"), Custo: strptr("50.00"),
				ICMSSaida: strptr("36.00"), PisCofins: strptr("0.00"), RestituicaoST: strptr("0.00"),
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
//
// Task A2 recompute (rule a, diferença simples — BA não está no conjunto dos
// 12%, a_inter=7%; fcp_embutido = "0" explicito ⇒ fcp = 0, aBase=aCusto=0.205):
// ICMS_oper=400x0.07=28.00. ICMS_total=400x0.205=82.00. DIFAL=82-28=54.00.
// BASE_PC=400x(1-0.205)=318.00. PIS/COFINS=0.0925x318.00=29.415 -> 29.42
// (meio exato, half-up sobe). Sum(sem restituição) = 40(comissao)+0(taxa)+
// 15(frete)+28.00(icms)+54.00(difal)+29.42(pis)+0(tarifa)+100(custo)-0 =
// 266.42. Margem = 400-266.42 = 133.58. Com restituição 40.00: Sum=226.42,
// Margem=173.58 (133.58+40.00).
func TestDecomposeICMSRestituicaoIncreasesMargem(t *testing.T) {
	base := func(restituicao string) DecomposeInput {
		return DecomposeInput{
			Preco: "400.00", ComissaoPct: "10", AliquotaPct: "4",
			Modalidade: ModalidadeClassico, Custo: money("100.00"),
			FreteProduto: money("15.00"),
			ICMSCell: &ICMSCell{
				UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
				AliquotaInterna: strptr("20.5"), Origprod: intp(0),
				// Task A4 (D2): par completo -- ver comentario do ICMS-1.
				FcpEmbutido:     strptr("0"),
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
	if *withoutR.MargemValor != "133.58" {
		t.Fatalf("withoutR.MargemValor = %s, want 133.58", *withoutR.MargemValor)
	}
	if *withR.MargemValor != "173.58" {
		t.Fatalf("withR.MargemValor = %s, want 173.58 (133.58 + 40.00 restituicao)", *withR.MargemValor)
	}

	diff := new(big.Rat).Sub(ratOf(t, *withR.MargemValor), ratOf(t, *withoutR.MargemValor))
	if FormatRatHalfUp(diff, 2) != "40.00" {
		t.Fatalf("margem delta = %s, want exactly 40.00 (restituicao adds dollar-for-dollar to margem)", FormatRatHalfUp(diff, 2))
	}

	assertSomaFechaICMS(t, withoutR)
	assertSomaFechaICMS(t, withR)
}
