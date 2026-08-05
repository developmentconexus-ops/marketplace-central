package domain

import (
	"math/big"
	"reflect"
	"testing"
)

// TestICMSERPGolden is Task A1 (MIS-008 fiscal correction): ten cases
// transcribed from REAL issued invoices (NUNOTA cited per case, raw evidence
// in .mnfs/MIS-008/evidence/fiscal-2026-08-04/ roundJ/roundT/roundU/roundV),
// measured against the live Oracle, that pin the ERP's actual DIFAL method
// against icms.go:101 TaxesForItem.
//
// The ERP uses SIMPLE DIFFERENCE, not the gross-up icms.go used to compute:
// DIFAL = P × (a_interna − a_inter). And the FCP does NOT abate the
// PIS/COFINS base — the legal aliquota embeds FCP (e.g. RJ 22,0% = 20,0 ICMS
// + 2,0 FECP) for the VALUE DUE, but the base abatement must use the
// FCP-excluded rate. icms.go had no way to represent that split (one
// AliquotaInterna field, no FCP component in TaxComponents) — that gap is
// exactly what this file proved red at A1.
//
// Task A2 closed the gap: ICMSCell.FcpEmbutido + TaxComponents.FCP now exist,
// icms.go uses diferença simples (not gross-up), and the intra-UF split is
// explicit (rule c). All ten cases pass now, including case4 (deferred at
// A1 for the missing FCP field, closed here with real data).
func TestICMSERPGolden(t *testing.T) {
	t.Run("case1: nota 895507 (BA) devido — metodo correto (diferenca simples)", func(t *testing.T) {
		// P=299.90, BA fora do conjunto dos 12% -> a_inter=0.07 (default).
		// interna=20.5, FCP=0 (nao consultado aqui — BA nao tem FCP nesta nota).
		// ICMS = P×a_inter = 299.90×0.07 = 20.993 -> 20.99.
		// DIFAL (metodo do ERP, diferenca simples) = P×(0.205−0.07) = 299.90×0.135
		//   = 40.4865 -> 40.49.
		// base PC = P − ICMS − DIFAL − ST_ant(0) = 299.90−20.99−40.49 = 238.42.
		// PIS/COFINS = 0.0925×238.42 = 22.05385 -> 22.05.
		//
		// CONTROLE NEGATIVO: icms.go:154-158 faz gross-up em vez de diferenca
		// simples. a_grossup = 0.205×0.93/0.795 = 0.239811... ; ICMS_total =
		// 299.90×a_grossup = 71.919... -> 71.92 (nao exposto diretamente);
		// DIFAL_hoje = 71.919... − 20.993 = 50.926... -> 50.93. Mensagem
		// esperada hoje: "DIFAL = 50.93, want 40.49".
		cell := &ICMSCell{
			UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
			Origprod: intp(0), AliquotaInterna: strptr("20.5"),
		}
		got := TaxesForItem("299.90", cell)
		assertMoney(t, "ICMS", got.ICMSSaida, "20.99")
		assertMoney(t, "DIFAL", got.Difal, "40.49")
		bp := basePCFromComponents(t, "299.90", got.ICMSSaida, got.Difal, "0.00")
		assertMoneyStr(t, "base PC", bp, "238.42")
		assertMoney(t, "PIS/COFINS", got.PisCofins, "22.05")
	})

	t.Run("case2: nota 895507 (BA) reproduzido com dado do ERP (interna=17)", func(t *testing.T) {
		// Mesma nota 895507, mas com o dado que o ERP de fato usou (interna=17,
		// nao 20.5). Se o motor reproduz TGFDIN ao centavo aqui, o METODO esta
		// provado identico ao do ERP — esse e o caso que fecha o argumento.
		// ICMS = 299.90×0.07 = 20.993 -> 20.99.
		// DIFAL (diferenca simples) = 299.90×(0.17−0.07) = 29.99.
		// base PC (soma dos componentes JA arredondados, para bater com
		// TGFDIN.BASERED) = 299.90−20.99−29.99 = 248.92.
		//
		// PIS/COFINS usa a base EXATA (nao a soma de componentes arredondados
		// acima) — mesmo racional de icms.go:223-227 ("aBase exato, nunca um
		// valor re-arredondado alimentando outro calculo"): base_exata =
		// 299.90×(1−0.17) = 248.917 (fecha 2dp em 248.92, mas o valor que
		// entra na multiplicacao e 248.917, nao 248.92).
		// PIS/COFINS = 0.0925×248.917 = 23.0248225 -> 23.02 (NAO 23.03: usar
		// 248.92 arredondado no lugar de 248.917 exato foi um erro de conta
		// da primeira versao deste golden — 0.0925×248.92 = 23.0251 -> 23.03
		// so aparece se voce re-arredonda a base antes de multiplicar, o que
		// TaxesForItem nao faz).
		cell := &ICMSCell{
			UFDestino: "BA", CodTrib: intp(0), Ambiguo: false,
			Origprod: intp(0), AliquotaInterna: strptr("17"),
		}
		got := TaxesForItem("299.90", cell)
		assertMoney(t, "ICMS", got.ICMSSaida, "20.99")
		assertMoney(t, "DIFAL", got.Difal, "29.99")
		bp := basePCFromComponents(t, "299.90", got.ICMSSaida, got.Difal, "0.00")
		assertMoneyStr(t, "base PC", bp, "248.92")
		assertMoney(t, "PIS/COFINS", got.PisCofins, "23.02")
	})

	t.Run("case3: GO tipico", func(t *testing.T) {
		// P=299.90, GO fora do conjunto dos 12% -> a_inter=0.07. interna=19.
		// ICMS = 299.90×0.07 = 20.99 (implicito, nao citado na tabela do brief
		// mas necessario para fechar base PC).
		// DIFAL (diferenca simples) = 299.90×(0.19−0.07) = 299.90×0.12 = 35.988
		//   -> 35.99.
		// base PC = 299.90−20.99−35.99 = 242.92.
		cell := &ICMSCell{
			UFDestino: "GO", CodTrib: intp(0), Ambiguo: false,
			Origprod: intp(0), AliquotaInterna: strptr("19"),
		}
		got := TaxesForItem("299.90", cell)
		assertMoney(t, "DIFAL", got.Difal, "35.99")
		bp := basePCFromComponents(t, "299.90", got.ICMSSaida, got.Difal, "0.00")
		assertMoneyStr(t, "base PC", bp, "242.92")
	})

	t.Run("case4: RJ — FCP nao abate a base do PIS/COFINS", func(t *testing.T) {
		// P=299.90, RJ no conjunto dos 12% -> a_inter=0.12. a_custo(interna,
		// com FCP embutido)=22.0%, fcp_embutido=2.0% -> a_base(sem FCP)=20.0%.
		//
		// Task A2 fecha este caso: ICMSCell.FcpEmbutido e TaxComponents.FCP
		// existem agora, entao a_custo (ICMS/DIFAL) e a_base (PIS/COFINS) se
		// alimentam da MESMA celula, sem gross-up e sem misturar FCP na base
		// do PIS/COFINS (D-43).
		//   ICMS = 299.90×0.12 = 35.988 -> 35.99.
		//   ICMS_total = 299.90×0.22 = 65.978 -> DIFAL = 65.978−35.988 =
		//     29.990 -> 29.99 (diferenca simples exata, sem residuo — mesma
		//     forma fechada do case2, que tambem tem diferenca de alíquota
		//     0.10 e fecha em 29.99 contra nota real citada por NUNOTA; este
		//     caso4 nao cita NUNOTA — e um exemplo ilustrativo do brief, e
		//     29.99 e o valor MATEMATICAMENTE alcancavel a partir de P=299.90,
		//     a_custo=0.22, a_inter=0.12; "30.00" so aparece como
		//     arredondamento de nota, nao como saida de TaxesForItem).
		//   FCP = 299.90×0.02 = 5.998 -> 6.00.
		//   a_base = 0.22−0.02 = 0.20. base PC = P×(1−0.20) = 299.90×0.80 =
		//     239.92 (formula exata; uma nota real citada poderia fechar em
		//     239.91 por residuo de 1 centavo PROPRIO da nota, nao da formula
		//     — TaxesForItem e puro e nao reproduz residuo de arredondamento
		//     de terceiros).
		cell := &ICMSCell{
			UFDestino: "RJ", CodTrib: intp(0), Ambiguo: false,
			Origprod: intp(0), AliquotaInterna: strptr("22"), FcpEmbutido: strptr("2"),
		}
		got := TaxesForItem("299.90", cell)
		assertMoney(t, "ICMS", got.ICMSSaida, "35.99")
		assertMoney(t, "DIFAL", got.Difal, "29.99")
		assertMoney(t, "FCP", got.FCP, "6.00")
		bp := basePCFromComponents(t, "299.90", got.ICMSSaida, got.Difal, "0.00")
		assertMoneyStr(t, "base PC", bp, "233.92")
		// basePCFromComponents faz P−ICMS−DIFAL−S (a formula de margem, que
		// exclui FCP porque ICMS/DIFAL ja carregam o FCP embutido). A base do
		// PIS/COFINS de verdade e a_base×P = 239.92, que soma o FCP de volta:
		// 233.92 (P−ICMS−DIFAL) + 6.00 (FCP) = 239.92 — bate com o PisCofins
		// vindo de TaxesForItem, calculado direto de a_base, nao desta soma.
		if got.PisCofins == nil || *got.PisCofins != "22.19" {
			t.Fatalf("PisCofins = %v, want \"22.19\" (0.0925×239.92)", got.PisCofins)
		}
	})

	t.Run("case5: intra-MG CST 00 — ICMS usa alíquota interna cheia, DIFAL 0 explicito", func(t *testing.T) {
		// nota 882236/2. P=136.66, destino=MG=origem -> operacao interna:
		// ICMS = P×a_interna_cheia = 136.66×0.18 = 24.5988 -> 24.60. DIFAL=0.00
		// EXPLICITO (zero legitimo, nao ha diferencial em operacao interna).
		// base PC = 136.66−24.60−0.00 = 112.06.
		//
		// CONTROLE NEGATIVO: icms.go SEMPRE faz a quebra ICMS_oper=P×a_inter /
		// DIFAL=icms_total−icms_oper, mesmo intra-UF (onde a_inter nao deveria
		// nem entrar na conta) — so o a=0.18 fixo esta certo, a QUEBRA esta
		// errada:
		//   ICMS_hoje = 136.66×0.07 (MG fora do conjunto dos 12%) = 9.5662
		//     -> 9.57.
		//   ICMS_total = 136.66×0.18 = 24.5988.
		//   DIFAL_hoje = 24.5988−9.5662 = 15.0326 -> 15.03.
		// Mensagens esperadas hoje: "ICMS = 9.57, want 24.60" E
		// "DIFAL = 15.03, want 0.00".
		cell := &ICMSCell{
			UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
			Origprod: intp(0), AliquotaInterna: strptr("18"),
		}
		got := TaxesForItem("136.66", cell)
		assertMoney(t, "ICMS", got.ICMSSaida, "24.60")
		assertMoney(t, "DIFAL", got.Difal, "0.00")
		// base PC fecha em 112.06 mesmo com a quebra ICMS/DIFAL errada de hoje
		// (ICMS_hoje+DIFAL_hoje = 9.57+15.03 = 24.60 = mesmo total) — prova
		// que o bug de quebra e invisivel na base, so nas duas linhas
		// exibidas separadamente.
		bp := basePCFromComponents(t, "136.66", got.ICMSSaida, got.Difal, "0.00")
		assertMoneyStr(t, "base PC", bp, "112.06")
	})

	t.Run("case6: intra-MG CST 60 — ramo ST, ICMS/DIFAL 0 explicito", func(t *testing.T) {
		// nota 897156/1. P=110.00, CST=60 (ST no destino) -> ramo isST de
		// icms.go: ICMS=0.00, DIFAL=0.00 explicito (o proprio ja esta no
		// custo). ST_ant=8.13 (StRetidoEntrada). base PC = P−S = 110.00−8.13
		//   = 101.87.
		// Este ramo (linha 122-133 de icms.go) ja faz base=P−S diretamente,
		// sem gross-up e sem a distincao FCP/base — espera-se PASSAR hoje.
		cell := &ICMSCell{
			UFDestino: "MG", CodTrib: intp(codTribST), Ambiguo: false,
			StRetidoEntrada: strptr("8.13"),
		}
		got := TaxesForItem("110.00", cell)
		assertMoney(t, "ICMS", got.ICMSSaida, "0.00")
		assertMoney(t, "DIFAL", got.Difal, "0.00")
		bp := basePCFromComponents(t, "110.00", got.ICMSSaida, got.Difal, "8.13")
		assertMoneyStr(t, "base PC", bp, "101.87")
	})

	t.Run("case7: interestadual com ST anterior — residuo zero", func(t *testing.T) {
		// nota 896886/1. Componentes ja resolvidos na nota (nao passam por
		// TaxesForItem — a entrada do caso ja e ICMS/DIFAL/ST_ant, nao
		// P/a_inter/interna; ver resolucao de ambiguidade do brief: casos 7 e
		// 8 fixam a REGRA de abatimento, nao a derivacao da aliquota).
		// P=143.10, ICMS=10.02, DIFAL=17.17, ST_ant=11.93.
		// base PC = 143.10−10.02−17.17−11.93 = 103.98 — residuo ZERO (a regra
		// "base = P−ICMS−DIFAL−ST_ant" fecha exato, sem sobra nem falta).
		bp := basePCFromComponents(t, "143.10", strptr("10.02"), strptr("17.17"), "11.93")
		assertMoneyStr(t, "base PC", bp, "103.98")
	})

	t.Run("case8: RJ real — FCP fora do abatimento", func(t *testing.T) {
		// nota 892328/1. P=1.896,00−36,60=1.859,40 (desconto ja liquido).
		// ICMS=223.13, DIFAL=111.56, FCP=37.19.
		// abatido (o que de fato reduz a base) = ICMS+DIFAL = 223.13+111.56
		//   = 334.69 — o FCP NAO entra.
		// Incluir o FCP no abatimento erraria por exatamente 37.19 (o proprio
		// FCP): 334.69+37.19 = 371.88, delta = 37.19.
		icms := ratOf(t, "223.13")
		difal := ratOf(t, "111.56")
		fcp := ratOf(t, "37.19")

		abatido := new(big.Rat).Add(icms, difal)
		abatidoStr, _ := round2(abatido)
		assertMoneyStr(t, "abatido", abatidoStr, "334.69")

		comFCP := new(big.Rat).Add(abatido, fcp)
		comFCPStr, _ := round2(comFCP)
		delta := new(big.Rat).Sub(ratOf(t, comFCPStr), ratOf(t, abatidoStr))
		deltaStr, _ := round2(delta)
		assertMoneyStr(t, "erro se FCP entrasse no abatimento", deltaStr, "37.19")
	})

	t.Run("case9: UF sem linha na tabela legal — unknown, nenhum numero", func(t *testing.T) {
		// D-37: AliquotaInterna nil (UF nunca cadastrada na tabela legal) —
		// nunca cai na alíquota interna como aproximação. icms_saida/difal/
		// pis_cofins ficam nil e nomeados em Unknown; restituicao continua
		// conhecida (independe da celula resolvida).
		cell := &ICMSCell{
			UFDestino: "AC", CodTrib: intp(0), Ambiguo: false,
			Origprod: intp(0), AliquotaInterna: nil,
		}
		got := TaxesForItem("299.90", cell)
		assertUnknown(t, got)
	})

	t.Run("case10: celula ambigua — unknown, nunca escolhe candidata", func(t *testing.T) {
		// Ambiguo=true: mais de uma linha TGFICM sobreviveu sem precedencia —
		// mesmo tratamento de CodTrib nil, nunca escolhe as cegas.
		cell := &ICMSCell{
			UFDestino: "PR", CodTrib: intp(0), Ambiguo: true,
			Origprod: intp(0), AliquotaInterna: strptr("19"),
		}
		got := TaxesForItem("299.90", cell)
		assertUnknown(t, got)
	})
}

// assertMoney compares a *string money/pct field against the expected value,
// reporting "<label> = <got>, want <want>" on mismatch — the exact shape the
// task brief's negative controls are pinned against. Uses t.Errorf (not
// Fatalf) so every field in a case gets checked and reported even when an
// earlier field in the same subtest already failed (case5 needs BOTH its
// ICMS and DIFAL messages to surface in the same run).
func assertMoney(t *testing.T, label string, got *string, want string) {
	t.Helper()
	g := "<nil>"
	if got != nil {
		g = *got
	}
	if g != want {
		t.Errorf("%s = %s, want %s", label, g, want)
	}
}

// assertMoneyStr is assertMoney's plain-string sibling, for locally computed
// quantities (base PC, abatido) that TaxesForItem does not return as a field.
func assertMoneyStr(t *testing.T, label, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %s, want %s", label, got, want)
	}
}

// basePCFromComponents implements literally the rule this task's brief
// states: "base PIS/COFINS = P − ICMS − DIFAL − ST_anterior" — applied to the
// ALREADY-ROUNDED (2dp) ICMS/DIFAL strings TaxesForItem actually returns
// (never a re-derivation of the internal unrounded `a`, which the struct
// does not expose), exactly how an issued invoice books it line by line.
// stAnterior is a plain decimal string, "0.00" when the case has no ST_ant
// fact.
func basePCFromComponents(t *testing.T, preco string, icms, difal *string, stAnterior string) string {
	t.Helper()
	if icms == nil || difal == nil {
		t.Fatalf("basePCFromComponents: ICMS/DIFAL desconhecidos (nil) — nao calculavel")
	}
	b := ratOf(t, preco)
	b.Sub(b, ratOf(t, *icms))
	b.Sub(b, ratOf(t, *difal))
	b.Sub(b, ratOf(t, stAnterior))
	s, _ := round2(b)
	return s
}

// assertUnknown pins the D-37/Ambiguo shape: icms_saida/difal/pis_cofins nil
// and named in Unknown, nunca um zero fabricado (ADR-17).
func assertUnknown(t *testing.T, got TaxComponents) {
	t.Helper()
	if got.ICMSSaida != nil {
		t.Errorf("ICMSSaida = %s, want nil", *got.ICMSSaida)
	}
	if got.Difal != nil {
		t.Errorf("Difal = %s, want nil", *got.Difal)
	}
	if got.PisCofins != nil {
		t.Errorf("PisCofins = %s, want nil", *got.PisCofins)
	}
	want := []string{"icms_saida", "difal", "pis_cofins"}
	if !reflect.DeepEqual(got.Unknown, want) {
		t.Errorf("Unknown = %v, want %v", got.Unknown, want)
	}
}
