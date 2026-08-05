package domain

import (
	"reflect"
	"testing"
)

// intp is the *int test helper for ICMSCell.CodTrib/Origprod, mirroring strptr's
// role for *string. Shared across this file and decompose_icms_golden_test.go
// (same package, same test binary).
func intp(i int) *int { return &i }

// TestTaxesForItemPureFormula pins the D-41 formula (docs/superpowers/plans/
// 2026-08-02-p2b-imposto-ex-ante-plan.md, Task 4) directly against TaxesForItem,
// independent of Decompose's comissão/frete/custo plumbing. Every expected value
// below is hand-computed with exact fractions (documented in the PR/commit
// description) — this is the vermelho-first pin for the pure math.
func TestTaxesForItemPureFormula(t *testing.T) {
	cases := []struct {
		name  string
		preco string
		cell  *ICMSCell
		want  TaxComponents
	}{
		{
			// P=1000, MG intra-UF (destino=origem): Task A2 rule (c) — operação
			// interna usa a alíquota interna INTEIRA para o ICMS, e o DIFAL é
			// 0,00 EXPLÍCITO (não a subtração interestadual aplicada por
			// engano). Rule (d): MG não tem mais bypass — lê
			// icms_aliquota_interna como qualquer UF (aqui 18%, fcp_embutido
			// nil ⇒ 0).
			// aCusto = 18/100 = 0,18. ICMS_oper = 1000×0,18 = 180,00.
			// DIFAL = 0,00 (UFDestino==UFOrigem=MG, ramo explícito).
			// FCP = 1000×0 = 0,00 (fcp_embutido nil ⇒ zeroIfNil=0).
			// aBase = aCusto−fcp = 0,18−0 = 0,18. BASE_PC = 1000×0,82−0 = 820,00.
			// PIS/COFINS = 0,0925×820,00 = 75,85 (inalterado — a_base não mudou).
			// UF=MG ⇒ restituição 0,00 explícito mesmo com restituicao_unit=
			// 50,00 (gatilho é a UF, não o valor).
			name:  "A: MG interno — ICMS pela alíquota interna inteira, DIFAL 0 explícito",
			preco: "1000.00",
			cell: &ICMSCell{
				UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
				Origprod: intp(0), AliquotaInterna: strptr("18"),
				// Task A4 (D2): par completo -- "0" explicito, nao nil.
				// AliquotaInterna presente com FcpEmbutido nil e um estado
				// que a fonte real nunca produz; sem este campo o gate novo
				// nomearia pis_cofins como desconhecido e este caso deixaria
				// de exercitar a formula que ele pretende pinar.
				FcpEmbutido:     strptr("0"),
				RestituicaoUnit: strptr("50.00"),
				FiscalDtRef:     strptr("2026-07-01"),
			},
			want: TaxComponents{
				ICMSSaida: strptr("180.00"), Difal: strptr("0.00"),
				FCP:       strptr("0.00"),
				PisCofins: strptr("75.85"), RestituicaoST: strptr("0.00"),
				FiscalDtRef: strptr("2026-07-01"),
			},
		},
		{
			// P=500, SP (não-MG, não-ST): origprod=0 não cai em {1,2,3,8}, e SP
			// está no conjunto dos 12% ⇒ a_inter=0,12. a_int=18% (tabela legal
			// SP, fcp_embutido nil ⇒ 0). Task A2 rule (a): DIFAL é a diferença
			// simples ICMS_total−ICMS_oper (o jeito que o ERP calcula,
			// TGFDIN.BASERED), não mais gross-up.
			// ICMS_oper = 500×0,12 = 60,00 exato.
			// ICMS_total = 500×0,18 = 90,00 exato.
			// DIFAL = 90,00−60,00 = 30,00.
			// FCP = 500×0 = 0,00 (fcp_embutido nil ⇒ zeroIfNil=0).
			// aBase = 0,18−0 = 0,18. BASE_PC = 500×0,82−50 = 360,00.
			// PIS/COFINS = 0,0925×360,00 = 33,30.
			// UF≠MG ⇒ restituição = restituicao_unit = 25,00.
			name:  "B: interestadual SP — diferença simples, a_int da tabela legal",
			preco: "500.00",
			cell: &ICMSCell{
				UFDestino: "SP", CodTrib: intp(0), Ambiguo: false,
				AliquotaInterna: strptr("18"), Origprod: intp(0),
				// Task A4 (D2): par completo -- ver comentario do caso A.
				FcpEmbutido:     strptr("0"),
				StRetidoEntrada: strptr("50.00"), RestituicaoUnit: strptr("25.00"),
				FiscalDtRef: strptr("2026-08-01"),
			},
			want: TaxComponents{
				ICMSSaida: strptr("60.00"), Difal: strptr("30.00"),
				FCP:       strptr("0.00"),
				PisCofins: strptr("33.30"), RestituicaoST: strptr("25.00"),
				FiscalDtRef: strptr("2026-08-01"),
			},
		},
		{
			// P=800, ST no destino (codtrib=60): a=0 força ICMS_oper e DIFAL a
			// "0.00" explícito — SEM consultar origprod/alíquota interna (ambos
			// nil de propósito, prova que o ramo ST não precisa deles).
			// PIS/COFINS ainda roda: BASE_PC = 800×(1−0) − 40(S retido) = 760,00.
			// PIS/COFINS = 0,0925×760,00 = 70,30 exato.
			// UF=RJ≠MG ⇒ restituição = 15,00 (prova que restituição não depende
			// do ramo ST).
			name:  "C: ST no destino — ICMS/DIFAL/FCP 0.00 explicito, PIS/COFINS deduz S",
			preco: "800.00",
			cell: &ICMSCell{
				UFDestino: "RJ", CodTrib: intp(60), Ambiguo: false,
				Origprod: nil, AliquotaInterna: nil,
				StRetidoEntrada: strptr("40.00"), RestituicaoUnit: strptr("15.00"),
			},
			want: TaxComponents{
				// Task A4 review finding: FCP costumava ficar nil e SEM NOME
				// neste ramo (nem zero fabricado nem desconhecido nomeado —
				// um terceiro estado que ADR-17 não admite). FCP é uma fatia
				// do ICMS de saída, e o ICMS de saída em si já é 0 explícito
				// no ramo ST (o próprio já está no custo) — a fatia dele
				// também é 0 explícito, mesma justificativa.
				ICMSSaida: strptr("0.00"), Difal: strptr("0.00"),
				FCP:       strptr("0.00"),
				PisCofins: strptr("70.30"), RestituicaoST: strptr("15.00"),
			},
		},
		{
			// célula ambígua: codtrib incerto ⇒ ICMS/DIFAL/PIS-COFINS ficam nil
			// e nomeados em Unknown — mas restituição NÃO depende da célula
			// resolvida (só de UF+restituicao_unit), então continua conhecida.
			name:  "D: ambiguo=true — icms/difal/pis_cofins unknown, restituicao independente",
			preco: "600.00",
			cell: &ICMSCell{
				UFDestino: "PR", CodTrib: intp(0), Ambiguo: true,
				AliquotaInterna: strptr("19.5"), Origprod: intp(5),
				RestituicaoUnit: strptr("10.00"),
			},
			want: TaxComponents{
				RestituicaoST: strptr("10.00"),
				// Task A4 review finding: "fcp" entra na lista — FCP também
				// fica nil aqui e antes desta correção ficava sem nome.
				Unknown: []string{"icms_saida", "difal", "fcp", "pis_cofins"},
			},
		},
		{
			// codtrib nil (classificação tributária nunca lida) é o MESMO
			// tratamento de ambiguo=true — célula não resolvida, nunca 0
			// fabricado.
			name:  "E: codtrib nil — mesmo tratamento de ambiguo",
			preco: "300.00",
			cell: &ICMSCell{
				UFDestino: "BA", CodTrib: nil, Ambiguo: false,
				AliquotaInterna: strptr("20.5"), Origprod: intp(0),
				RestituicaoUnit: strptr("5.00"),
			},
			want: TaxComponents{
				RestituicaoST: strptr("5.00"),
				// Task A4 review finding: "fcp" entra na lista — mesmo motivo
				// do caso D.
				Unknown: []string{"icms_saida", "difal", "fcp", "pis_cofins"},
			},
		},
		{
			// origprod ∈ {1,2,3,8} vence a UF mesmo quando a UF de destino está
			// no conjunto dos 12% (RJ) — a_inter=0,04, não 0,12.
			// a_int=22% (RJ, tabela legal), fcp_embutido=2% (D-43 — aCusto e
			// aBase derivam da MESMA linha). Task A2 rule (a): diferença
			// simples, não gross-up.
			// aCusto = 22/100 = 0,22. fcp = 2/100 = 0,02.
			// ICMS_oper = 1000×0,04 = 40,00.
			// ICMS_total = 1000×0,22 = 220,00.
			// DIFAL = 220,00−40,00 = 180,00.
			// FCP = 1000×0,02 = 20,00.
			// aBase = 0,22−0,02 = 0,20. BASE_PC = 1000×(1−0,20) = 800,00.
			// PIS/COFINS = 0,0925×800,00 = 74,00.
			// restituicao_unit nil ⇒ tratado como 0 (produto sem fato de
			// restituição registrado, não "desconhecido" — Oracle não distingue
			// as duas coisas para este fato, T4 escopo).
			name:  "F: origprod prioriza sobre UF — 4% mesmo em RJ (conjunto dos 12%)",
			preco: "1000.00",
			cell: &ICMSCell{
				UFDestino: "RJ", CodTrib: intp(0), Ambiguo: false,
				AliquotaInterna: strptr("22"), FcpEmbutido: strptr("2"),
				Origprod: intp(1),
			},
			want: TaxComponents{
				ICMSSaida: strptr("40.00"), Difal: strptr("180.00"),
				FCP:       strptr("20.00"),
				PisCofins: strptr("74.00"), RestituicaoST: strptr("0.00"),
			},
		},
		{
			name:  "G: ICMSCell inteiramente ausente — tudo unknown, inclusive restituicao",
			preco: "100.00",
			cell:  nil,
			want: TaxComponents{
				// Task A4 review finding: "fcp" entra na lista — cell nil
				// significa "o quadro fiscal inteiro é desconhecido",
				// inclusive a fatia de FCP.
				Unknown: []string{"icms_saida", "difal", "fcp", "pis_cofins", "restituicao_st"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := TaxesForItem(tc.preco, tc.cell)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("TaxesForItem mismatch\n got: %+v\nwant: %+v", dumpTax(got), dumpTax(tc.want))
			}
		})
	}
}

// dumpTax renders *string fields for readable test failure output.
func dumpTax(c TaxComponents) string {
	deref := func(s *string) string {
		if s == nil {
			return "<nil>"
		}
		return *s
	}
	return "ICMSSaida=" + deref(c.ICMSSaida) + " Difal=" + deref(c.Difal) +
		" PisCofins=" + deref(c.PisCofins) + " RestituicaoST=" + deref(c.RestituicaoST) +
		" FiscalDtRef=" + deref(c.FiscalDtRef) + " Unknown=" + strJoin(c.Unknown)
}

func strJoin(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ","
		}
		out += s
	}
	return out
}

// TestTaxesForItemMaxZeroClampsBase proves MAX(0,…) is not cosmetic: S retido
// maior que o líquido não pode virar base negativa / crédito fantasma.
// P=200, MG intra-UF (rule d: MG lê a tabela, 18%, fcp_embutido nil ⇒ 0).
// ICMS_oper = 200×0,18 = 36,00. DIFAL = 0,00 explícito (intra-MG, rule c).
// aBase = 0,18−0 = 0,18. P×(1−aBase) = 200×0,82 = 164,00. S=500,00 (muito
// maior) ⇒ BASE_PC = MAX(0,−336) = 0.
// PIS/COFINS = 0,0925×0 = 0,00 — CONHECIDO (zero legítimo por clamp), não nil.
func TestTaxesForItemMaxZeroClampsBase(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("18"),
		// Task A4 (D2): par completo -- ver comentario do caso A em
		// TestTaxesForItemPureFormula.
		FcpEmbutido:     strptr("0"),
		StRetidoEntrada: strptr("500.00"),
	}
	got := TaxesForItem("200.00", cell)
	if got.PisCofins == nil || *got.PisCofins != "0.00" {
		t.Fatalf("PisCofins = %v, want \"0.00\" (MAX(0,...) clamp) — sem o MAX sairia base negativa", dumpTax(got))
	}
	if got.ICMSSaida == nil || *got.ICMSSaida != "36.00" {
		t.Fatalf("ICMSSaida = %v, want 36.00 (200x0.18, intra-MG usa a alíquota inteira)", dumpTax(got))
	}
	if got.Difal == nil || *got.Difal != "0.00" {
		t.Fatalf("Difal = %v, want 0.00 (intra-MG explícito, rule c)", dumpTax(got))
	}
}

// TestTaxesForItemMutationControlNegativo is the mandated controle negativo: it
// simulates the mutation "ST devolve 0.00 silencioso" vs "ambigua devolve 0.00
// em vez de nil" by asserting the TWO branches are OBSERVABLY different in a way
// a broken implementation (that collapses them) would fail. See the T4 report
// for the actual red-run transcript where icms.go was mutated and this test
// (or its golden sibling) was watched fail.
func TestTaxesForItemSTZeroIsNotSameAsAmbiguoNil(t *testing.T) {
	st := TaxesForItem("100.00", &ICMSCell{UFDestino: "SP", CodTrib: intp(60), Ambiguo: false, Origprod: intp(0), AliquotaInterna: strptr("18")})
	if st.ICMSSaida == nil || *st.ICMSSaida != "0.00" {
		t.Fatalf("ST ICMSSaida = %v, want explicit \"0.00\" (zero legitimo)", dumpTax(st))
	}
	if len(st.Unknown) != 0 {
		t.Fatalf("ST Unknown = %v, want empty — ST zero e resposta, nao lacuna", st.Unknown)
	}

	ambiguo := TaxesForItem("100.00", &ICMSCell{UFDestino: "SP", CodTrib: intp(0), Ambiguo: true, Origprod: intp(0), AliquotaInterna: strptr("18")})
	if ambiguo.ICMSSaida != nil {
		t.Fatalf("ambiguo ICMSSaida = %v, want nil (unknown, nao zero fabricado)", *ambiguo.ICMSSaida)
	}
	if len(ambiguo.Unknown) == 0 {
		t.Fatalf("ambiguo Unknown vazio — celula incerta tem que nomear a lacuna")
	}
}
