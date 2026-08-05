package domain

import "testing"

// Task A4 — "desconhecido nunca vira zero em TaxesForItem". Cinco defeitos
// MEDIDOS (task-A4-brief.md), vermelho-primeiro: cada teste abaixo reproduz
// o número fabricado que icms.go devolvia ANTES do conserto e, depois do
// conserto, assere que o componente correspondente virou nil e o
// desconhecido NOMEADO exato aparece em Unknown — nunca só "não vazio".

// TestTaxesForItemD1_AmbiguousST60NeverZero is D1: uma célula AMBÍGUA com
// CST 60 caía direto no ramo ST (icms.go, antes do portão de
// desconhecido) e devolvia ICMS/DIFAL "0.00" explícito com Unknown vazio —
// indistinguível do zero LEGÍTIMO do ramo ST resolvido. Medido:
// TaxesForItem(P=100, UFDestino="SP", CodTrib=60, Ambiguo=true) devolvia
// ICMSSaida="0.00" Difal="0.00" PisCofins="9.25" Unknown=[] (len=0).
// Uma célula ambígua é uma célula que o resolvedor NÃO conseguiu resolver —
// o portão de desconhecido (CodTrib==nil || Ambiguo) tem que rodar ANTES de
// qualquer decisão de ramo, inclusive o ramo ST.
func TestTaxesForItemD1_AmbiguousST60NeverZero(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "SP", CodTrib: intp(codTribST), Ambiguo: true,
		Origprod: intp(0), AliquotaInterna: strptr("18"),
	}
	got := TaxesForItem("100.00", cell)

	if got.ICMSSaida != nil {
		t.Errorf("ICMSSaida = %s, want nil — célula ambígua nunca escolhe o ramo ST às cegas", *got.ICMSSaida)
	}
	if got.Difal != nil {
		t.Errorf("Difal = %s, want nil", *got.Difal)
	}
	// Coordinator review (post-A4): FCP ficava nil e SEM NOME aqui — nem
	// zero fabricado nem desconhecido nomeado. Célula não resolvida também
	// torna a fatia FCP desconhecida (não há sequer um ICMS de saída do qual
	// derivá-la).
	if got.FCP != nil {
		t.Errorf("FCP = %s, want nil — célula ambígua nunca escolhe o ramo ST às cegas", *got.FCP)
	}
	if got.PisCofins != nil {
		t.Errorf("PisCofins = %s, want nil", *got.PisCofins)
	}
	want := []string{"icms_saida", "difal", "fcp", "pis_cofins"}
	if !slicesEqual(got.Unknown, want) {
		t.Errorf("Unknown = %v, want %v", got.Unknown, want)
	}
	// Restituição continua conhecida: não depende do CodTrib/Ambiguo, só da
	// UF de destino (SP≠MG ⇒ RestituicaoUnit nil ⇒ 0 explícito).
	if got.RestituicaoST == nil || *got.RestituicaoST != "0.00" {
		t.Errorf("RestituicaoST = %v, want \"0.00\" (independente da célula resolvida)", derefOrNil(got.RestituicaoST))
	}
}

// TestTaxesForItemD2_FcpEmbutidoNilIsUnknownNotZero is D2: com
// AliquotaInterna presente e FcpEmbutido nil, icms.go tratava nil como 0
// (zeroIfNil) e devolvia FCP="0.00"/PisCofins fabricado. Medido: RJ, P=100,
// AliquotaInterna="22", FcpEmbutido nil ⇒ FCP="0.00" PisCofins="7.22" (com o
// FCP real "2" da mesma linha ⇒ FCP="2.00" PisCofins="7.40" — duas saídas
// diferentes, nenhuma com desconhecido). A porta devolve alíquota+FCP da
// MESMA linha vigente e os dois nil JUNTOS (pricing/ports/tax_matrix.go,
// orders/adapters/pricingtax/reader.go) — presente+nil é um estado
// que a fonte real nunca produz.
//
// Escopo do conserto: só o componente que DEPENDE do FCP (PIS/COFINS, via
// aBase=aCusto−fcp) fica desconhecido. ICMSSaida/Difal dependem só de
// aCusto/a_inter (nunca de FcpEmbutido) e continuam CONHECIDOS — RJ está no
// conjunto dos 12%: ICMS_oper=100×0,12=12,00, ICMS_total=100×0,22=22,00,
// DIFAL=22,00−12,00=10,00.
func TestTaxesForItemD2_FcpEmbutidoNilIsUnknownNotZero(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "RJ", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("22"),
		FcpEmbutido: nil,
	}
	got := TaxesForItem("100.00", cell)

	if got.FCP != nil {
		t.Errorf("FCP = %s, want nil — FcpEmbutido ausente com AliquotaInterna presente é estado inválido, não zero", *got.FCP)
	}
	if got.PisCofins != nil {
		t.Errorf("PisCofins = %s, want nil", *got.PisCofins)
	}
	// Coordinator review (post-A4): "fcp" entra na lista — a mesma ausência
	// (FcpEmbutido nil) que torna PisCofins desconhecido também torna a
	// própria fatia FCP desconhecida; antes desta correção FCP ficava nil
	// mas sem nome, ao lado de um pis_cofins nomeado pela MESMA causa.
	want := []string{"fcp", "pis_cofins"}
	if !slicesEqual(got.Unknown, want) {
		t.Errorf("Unknown = %v, want %v", got.Unknown, want)
	}
	if got.ICMSSaida == nil || *got.ICMSSaida != "12.00" {
		t.Errorf("ICMSSaida = %v, want 12.00 (não depende de FcpEmbutido)", derefOrNil(got.ICMSSaida))
	}
	if got.Difal == nil || *got.Difal != "10.00" {
		t.Errorf("Difal = %v, want 10.00 (não depende de FcpEmbutido)", derefOrNil(got.Difal))
	}
}

// TestTaxesForItemD3_FcpEmbutidoAboveAliquotaIsUnknown is D3: fcp > aCusto
// viola o invariante da fonte (0 ≤ fcp ≤ aCusto ≤ 100) e icms.go não
// bloqueava — calculava aBase negativo e PIS/COFINS sobre uma base maior que
// o preço. P=100, MG, AliquotaInterna="18" (aCusto=0,18), FcpEmbutido="25"
// (fcp=0,25 > aCusto): dado inválido, não entrada exótica. ICMS/Difal
// continuam conhecidos (dependem só de aCusto, MG intra-UF): ICMS=18,00,
// DIFAL=0,00 explícito.
func TestTaxesForItemD3_FcpEmbutidoAboveAliquotaIsUnknown(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("18"),
		FcpEmbutido: strptr("25"),
	}
	got := TaxesForItem("100.00", cell)

	if got.FCP != nil {
		t.Errorf("FCP = %s, want nil — fcp=0.25 > aCusto=0.18 viola o invariante 0≤fcp≤aCusto", *got.FCP)
	}
	if got.PisCofins != nil {
		t.Errorf("PisCofins = %s, want nil — base não pode nascer negativa por FCP inválido", *got.PisCofins)
	}
	// Coordinator review (post-A4): "fcp" entra na lista — mesmo motivo do
	// D2 (o dado inválido que bloqueia PisCofins também bloqueia FCP).
	want := []string{"fcp", "pis_cofins"}
	if !slicesEqual(got.Unknown, want) {
		t.Errorf("Unknown = %v, want %v", got.Unknown, want)
	}
	if got.ICMSSaida == nil || *got.ICMSSaida != "18.00" {
		t.Errorf("ICMSSaida = %v, want 18.00", derefOrNil(got.ICMSSaida))
	}
	if got.Difal == nil || *got.Difal != "0.00" {
		t.Errorf("Difal = %v, want 0.00 (intra-MG explícito)", derefOrNil(got.Difal))
	}
}

// TestTaxesForItemD4_EmptyUFDestinoIsUnknownNotInterestadualDefault is D4:
// UFDestino="" caía no default interestadual de 7% em vez de virar
// desconhecido. Medido: P=100, UFDestino="", AliquotaInterna="18" ⇒
// ICMS="7.00" DIFAL="11.00" PisCofins="7.59" — UF vazia NÃO é UF, não existe
// operação sem destino conhecido. O gate roda antes de qualquer outra
// decisão (inclusive antes da restituição, que também depende de comparar a
// UF contra MG — "" não é um fato de UF válido para essa comparação),
// mesmo shape do cell==nil (4 componentes, incluindo restituicao_st).
func TestTaxesForItemD4_EmptyUFDestinoIsUnknownNotInterestadualDefault(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("18"),
	}
	got := TaxesForItem("100.00", cell)

	if got.ICMSSaida != nil {
		t.Errorf("ICMSSaida = %s, want nil — não existe operação sem destino conhecido", *got.ICMSSaida)
	}
	if got.Difal != nil {
		t.Errorf("Difal = %s, want nil", *got.Difal)
	}
	if got.FCP != nil {
		t.Errorf("FCP = %s, want nil", *got.FCP)
	}
	if got.PisCofins != nil {
		t.Errorf("PisCofins = %s, want nil", *got.PisCofins)
	}
	if got.RestituicaoST != nil {
		t.Errorf("RestituicaoST = %s, want nil — comparação contra MG não é fato válido sem UF", *got.RestituicaoST)
	}
	// Coordinator review (post-A4): "fcp" entra na lista — UF vazia bloqueia
	// o quadro fiscal inteiro, inclusive a fatia FCP.
	want := []string{"icms_saida", "difal", "fcp", "pis_cofins", "restituicao_st"}
	if !slicesEqual(got.Unknown, want) {
		t.Errorf("Unknown = %v, want %v", got.Unknown, want)
	}
}

// TestTaxesForItemD5_OrigprodNilBlocksOnlyInterestadualBranch is D5:
// Origprod nil era tratado como "nacional" (fora de origprodImportado) e o
// a_inter interestadual usava a UF normalmente — produto IMPORTADO
// devolveria 4%/18%, mas Origprod nil devolvia a quebra de 12%/10% (RJ) como
// se fosse nacional CONHECIDO. Origem desconhecida é desconhecida; a
// alíquota interestadual DEPENDE dela.
//
// Escopo do conserto: só bloqueia no ramo que efetivamente CONSULTA
// Origprod — o interestadual (aInterFor). O ramo intra-MG nunca chama
// aInterFor e não deve ser afetado (ver o controle positivo MG→MG abaixo).
func TestTaxesForItemD5_OrigprodNilBlocksOnlyInterestadualBranch(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "RJ", CodTrib: intp(0), Ambiguo: false,
		Origprod: nil, AliquotaInterna: strptr("22"), FcpEmbutido: strptr("2"),
	}
	got := TaxesForItem("100.00", cell)

	if got.ICMSSaida != nil {
		t.Errorf("ICMSSaida = %s, want nil — a_inter depende de Origprod (nacional vs importado)", *got.ICMSSaida)
	}
	if got.Difal != nil {
		t.Errorf("Difal = %s, want nil", *got.Difal)
	}
	if got.FCP != nil {
		t.Errorf("FCP = %s, want nil — sem ICMS de saída conhecido não há fatia FCP para derivar", *got.FCP)
	}
	if got.PisCofins != nil {
		t.Errorf("PisCofins = %s, want nil", *got.PisCofins)
	}
	// Coordinator review (post-A4): "fcp" entra na lista — mesmo bloqueio de
	// Origprod nil que impede ICMSSaida/Difal impede a fatia FCP.
	want := []string{"icms_saida", "difal", "fcp", "pis_cofins"}
	if !slicesEqual(got.Unknown, want) {
		t.Errorf("Unknown = %v, want %v", got.Unknown, want)
	}
}

// TestTaxesForItemPositiveControlSTResolvedStaysExplicitZero is the D1
// mandated POSITIVE control: uma célula ST RESOLVIDA (CodTrib=60,
// Ambiguo=false, dados completos) continua ICMS/DIFAL 0 EXPLÍCITO, sem
// nenhum desconhecido — prova que o conserto de D1 não afrouxou o caso
// legítimo, só o ambíguo. MG, P=110,00, S=8,13 (mesmos números do golden
// case6 de icms_erp_golden_test.go, repetidos aqui como controle isolado do
// D1).
func TestTaxesForItemPositiveControlSTResolvedStaysExplicitZero(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "MG", CodTrib: intp(codTribST), Ambiguo: false,
		StRetidoEntrada: strptr("8.13"),
	}
	got := TaxesForItem("110.00", cell)

	if got.ICMSSaida == nil || *got.ICMSSaida != "0.00" {
		t.Errorf("ICMSSaida = %v, want explicit \"0.00\" (zero legítimo do ramo ST resolvido)", derefOrNil(got.ICMSSaida))
	}
	if got.Difal == nil || *got.Difal != "0.00" {
		t.Errorf("Difal = %v, want explicit \"0.00\"", derefOrNil(got.Difal))
	}
	// Coordinator review (post-A4): FCP é uma fatia do ICMS de saída — no
	// ramo ST resolvido o ICMS de saída em si é 0 explícito, então a fatia
	// dele também é 0 explícito, nunca nil/desconhecido. Antes desta
	// correção FCP ficava nil aqui SEM entrar em Unknown (nem zero, nem
	// nomeado — o próprio defeito que esta task existe para matar).
	if got.FCP == nil || *got.FCP != "0.00" {
		t.Errorf("FCP = %v, want explicit \"0.00\" (fatia do ICMS de saída, que já é 0 no ramo ST)", derefOrNil(got.FCP))
	}
	if len(got.Unknown) != 0 {
		t.Errorf("Unknown = %v, want empty — célula ST resolvida não tem lacuna", got.Unknown)
	}
}

// TestTaxesForItemPositiveControlMGIntraResolvedStaysExplicitZero is the
// mandated POSITIVE control for the intra-MG rule: uma operação MG→MG
// RESOLVIDA (AliquotaInterna presente, FcpEmbutido presente — par completo)
// continua DIFAL 0,00 EXPLÍCITO, sem desconhecido, mesmo depois dos gates de
// D2/D3/D5 — prova que os novos portões só bloqueiam dado FALTANTE/INVÁLIDO,
// nunca o caso resolvido. P=100, aCusto=18%, fcp=0% (par completo, "0"
// explícito, não nil): ICMS=100×0,18=18,00, DIFAL=0,00.
func TestTaxesForItemPositiveControlMGIntraResolvedStaysExplicitZero(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "MG", CodTrib: intp(0), Ambiguo: false,
		Origprod: intp(0), AliquotaInterna: strptr("18"), FcpEmbutido: strptr("0"),
	}
	got := TaxesForItem("100.00", cell)

	if got.ICMSSaida == nil || *got.ICMSSaida != "18.00" {
		t.Errorf("ICMSSaida = %v, want 18.00", derefOrNil(got.ICMSSaida))
	}
	if got.Difal == nil || *got.Difal != "0.00" {
		t.Errorf("Difal = %v, want explicit \"0.00\" (intra-MG, zero legítimo)", derefOrNil(got.Difal))
	}
	// Coordinator review (post-A4), question 2 (intra-MG half): o ramo
	// intra-MG NÃO tem atalho próprio para FCP — cai no mesmo cálculo
	// compartilhado dos gates D2/D3 que o ramo interestadual usa (icms.go,
	// após o split MG/interestadual). Com FcpEmbutido="0" presente (par
	// completo), FCP é CALCULADO normalmente (100×0=0,00), nunca nil.
	if got.FCP == nil || *got.FCP != "0.00" {
		t.Errorf("FCP = %v, want computed 0.00 (par completo, ramo intra-MG usa o mesmo cálculo compartilhado de D2/D3)", derefOrNil(got.FCP))
	}
	if got.PisCofins == nil {
		t.Errorf("PisCofins = nil, want a computed value — par AliquotaInterna+FcpEmbutido completo, nunca desconhecido")
	}
	if len(got.Unknown) != 0 {
		t.Errorf("Unknown = %v, want empty — célula MG→MG resolvida não tem lacuna", got.Unknown)
	}
}

// slicesEqual is a tiny []string comparator local to this file — avoids
// pulling reflect.DeepEqual's nil-vs-empty-slice foot-gun into the D1/D4/D5
// assertions (Unknown is built with append, so a "no unknown" result is a
// nil slice, and every non-empty want here is a literal — ordinary equality
// by index is what these tests actually mean).
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
