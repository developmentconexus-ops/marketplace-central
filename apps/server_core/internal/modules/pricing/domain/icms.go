package domain

import "math/big"

// aInterUF12 is the set of UF destino whose D-41 a_inter is 12% instead of
// the 7% default (plan docs/superpowers/plans/2026-08-02-p2b-imposto-ex-ante-plan.md
// line 426). Deliberately NOT named interestadual12 and NOT shared with
// difal.go's own interestadual12 (which includes MG and serves a different,
// older DIFAL concept, IC-04) — the two sets differ (this one excludes MG:
// MG's a_inter falls to the 7% default, only its `a` itself is fixed at
// 0,18) and conflating them would silently mis-rate MG-destino items.
var aInterUF12 = map[string]bool{
	"SP": true, "RJ": true, "PR": true, "SC": true, "RS": true,
}

// origprodImportado is the ORIGPROD set that forces a_inter=4% regardless of
// UF destino. D-46 debt: FCI / conteúdo de importação por operação não é
// consultado — correto para o catálogo medido (metal sanitário/papeleiro),
// falso se entrar produto com conteúdo de importação resolvido caso a caso.
var origprodImportado = map[int]bool{1: true, 2: true, 3: true, 8: true}

// codTribST is the icms_matrix_mirror.codtrib value meaning "ICMS-ST neste
// destino" — o ICMS próprio já ficou embutido no custo (CUSSEMICM, ramo ST).
const codTribST = 60

// ufOrigemMG is the company's fixed origin UF (D-17). Two independent rules
// key off it: (1) intra-MG é operação INTERNA — ICMS de saída usa a alíquota
// interna cheia (aCusto, lida de icms_aliquota_interna como qualquer UF desde
// a Task A2 — o 0,18 hardcoded saiu) e DIFAL é 0 explícito, ao invés da
// quebra ICMS_oper=P×a_inter/DIFAL=P×aCusto−ICMS_oper que só vale
// interestadual; (2) restituição de ICMS-ST só sai de operação
// interestadual — gatilho literal "UF <> 13" da fonte Oracle
// (PAN_GET_CUSVAR_MNOBRE); intra-MG ela é 0 explícito.
const ufOrigemMG = "MG"

// pisCofinsRate is the STF Tema 69 (exclusão do ICMS da base) + STJ Tema 1125
// (exclusão do ICMS-ST retido) crédito: 9,25% é hardcoded no CUSSEMICM do
// ERP — D-39, espelhamos o número, não a regra.
var pisCofinsRate = big.NewRat(925, 10000)

// ICMSCell is the already-resolved input for TaxesForItem: the matrix cell
// (codtrib/ambiguo) plus the product's fiscal facts, all read by the port
// before this pure function runs — same shape convention as DecomposeInput's
// pre-resolved comissão/frete/custo.
type ICMSCell struct {
	// UFDestino is the order's destination UF. Only meaningful when the
	// *ICMSCell itself is non-nil (a nil *ICMSCell means the whole tax
	// picture is unknown, not that UFDestino is "").
	UFDestino string

	// CodTrib is icms_matrix_mirror.codtrib for (uf_origem=MG, UFDestino,
	// GrupoICMS). nil ⇒ célula não resolvida (ausente OU nunca lida) —
	// ICMSSaida/Difal/PisCofins ficam UNKNOWN. codTribST (60) ⇒ produto ST
	// neste destino: zero explícito, não desconhecido.
	CodTrib *int
	// Ambiguo ⇒ mais de uma linha TGFICM sobreviveu à lista branca sem
	// precedência estabelecida (ResolveCell) — mesmo tratamento de CodTrib
	// nil: nunca escolhe uma candidata às cegas.
	Ambiguo bool

	// Origprod is products_mirror.origprod (TGFPRO.ORIGPROD). Um valor fora
	// de {1,2,3,8} é nacional CONHECIDO e a_inter cai para o ramo por UF.
	// nil NÃO é a mesma coisa: origem desconhecida é desconhecida, e desde o
	// gate D5 (Task A4, linhas 248-250) ela É gatilho de desconhecido — mas
	// só no ramo interestadual, o único que consulta aInterFor. O ramo
	// intra-UF nunca lê este campo e não se importa se ele é nil.
	//
	// Este comentário dizia o contrário ("origprod nunca é ele próprio o
	// gatilho de desconhecido") e ficou falso no MESMO diff que introduziu o
	// gate; corrigido em vez de mantido, porque uma auditoria de ADR-17
	// lendo a doc do campo receberia o oposto do que o código faz.
	Origprod *int
	// AliquotaInterna is icms_aliquota_interna.aliquota — a_custo, a carga
	// interna cheia do destino COM o FCP embutido onde ele é geral (D-43),
	// para UFDestino, como string decimal em percentual ("18" = 18%). nil ⇒
	// UF sem linha na tabela legal (D-37) — pendência explícita. A2: MG deixa
	// de ser literal (era 0,18 hardcoded) e consulta este campo como
	// qualquer outra UF (migrations/0094_icms_matrix.sql:98, ('MG', 18.0));
	// só o ramo ST (a=0) não toca este campo.
	AliquotaInterna *string
	// FcpEmbutido is icms_aliquota_interna.fcp_embutido para UFDestino — a
	// fatia do FCP que já está embutida em AliquotaInterna. Existe porque,
	// medido, o ERP abate ICMS+DIFAL (que usam a alíquota cheia) da base do
	// PIS/COFINS mas DEIXA o FCP dentro dela: 6 itens do RJ com resíduo
	// exatamente igual a −FCP (roundV.txt V2) — uma fonte (a mesma linha
	// vigente de AliquotaInterna), duas alíquotas derivadas, nunca dois
	// cálculos independentes. nil é tratado como 0 (zero legítimo: a maioria
	// das UFs não tem FCP geral e a coluna vem 0, não ausente — só é
	// consultado quando AliquotaInterna também é, mesma linha vigente).
	FcpEmbutido *string
	// Task A4 (D2): nil aqui so e legitimo quando AliquotaInterna tambem e
	// nil (D-37, "UF sem linha" -- os dois viajam juntos, ver
	// pricing/ports/tax_matrix.go:18-22 e
	// orders/adapters/pricingtax/reader.go:201-208: AliquotaInternaFor
	// devolve o par da MESMA linha vigente; FcpEmbutidoPct e uma STRING
	// nao-ponteiro dentro do struct nao-nil -- "0" para UF sem FCP geral,
	// nunca ausente). AliquotaInterna presente com FcpEmbutido nil e um
	// estado que a fonte real nunca produz; so existe porque este campo e
	// *string e nao acompanha AliquotaInterna num tipo unico (mudanca de
	// forma adiada -- exigiria mexer em reader.go:201, fora do escopo desta
	// task). TaxesForItem trata essa combinacao como desconhecido nomeado
	// (fcp e pis_cofins), nunca como zero.

	// StRetidoEntrada is products_mirror.st_retido_entrada ("S"), string
	// decimal em dinheiro. nil é tratado como 0 — a própria fonte Oracle não
	// distingue "genuinamente zero" de "nunca registrado" para este fato.
	StRetidoEntrada *string
	// RestituicaoUnit is products_mirror.restituicao_unit ("R"), string
	// decimal em dinheiro. nil é tratado como 0, mesmo racional de
	// StRetidoEntrada.
	RestituicaoUnit *string

	// FiscalDtRef is products_mirror.fiscal_dt_ref, passado adiante sem
	// alteração para o chamador (Task 6/7) poder sinalizar bloco de ST
	// envelhecido (D-44). TaxesForItem não usa a idade para decidir nada.
	FiscalDtRef *string
}

// TaxComponents is TaxesForItem's pure output: ICMS de saída e DIFAL — ambos
// derivados da mesma alíquota efetiva colapsada `a` (D-41, "uma fonte, duas
// linhas, nunca dois cálculos") — mais PIS/COFINS débito e restituição de
// ICMS-ST. Qualquer campo nil se nomeia em Unknown; nunca um 0 fabricado
// (ADR-17).
type TaxComponents struct {
	ICMSSaida *string
	Difal     *string
	// FCP is P × fcp_embutido — a fatia do ICMS devido que corresponde ao
	// FCP, informativa (já está DENTRO de ICMSSaida+Difal, que somam
	// a_custo; somar FCP de novo no cálculo de margem dobraria a conta).
	// Task A4 review finding: até esta revisão FCP ficava nil e SEM NOME em
	// sete caminhos de retorno antecipado (violação do próprio ADR-17 que
	// este campo deveria seguir) — corrigido: toda vez que o fato-fonte
	// falta, "fcp" entra em Unknown junto de icms_saida/difal/pis_cofins;
	// no ramo ST o ICMS de saída (e portanto sua fatia FCP) é 0 explícito,
	// não desconhecido, mesma lógica de ICMSSaida/Difal nesse ramo. FCP só
	// é computado (nem nil nem 0 fabricado) no ramo MG-interno/interestadual
	// resolvido, quando FcpEmbutido está presente e dentro do invariante
	// 0 ≤ fcp ≤ aCusto (D2/D3).
	FCP           *string
	PisCofins     *string
	RestituicaoST *string
	FiscalDtRef   *string
	Unknown       []string
}

// TaxesForItem applies the D-41 formula to one item's preço, given its
// already-resolved ICMSCell. Pure, no I/O — mirrors decomposeWithLimiar's
// separation between pre-resolved inputs and formula.
func TaxesForItem(preco string, cell *ICMSCell) TaxComponents {
	if cell == nil {
		return TaxComponents{Unknown: []string{"icms_saida", "difal", "fcp", "pis_cofins", "restituicao_st"}}
	}

	// D4 (Task A4): UF vazia não é UF — não existe operação sem destino
	// conhecido. Antes deste gate, "" caía direto no default interestadual
	// de 7% (aInterFor não trata "" como especial, só como "fora do
	// conjunto dos 12%"), fabricando ICMS/DIFAL/PIS-COFINS a partir de um
	// destino que não existe. A comparação UFDestino==ufOrigemMG que decide
	// a restituição também não é um fato válido sem UF — mesma forma do
	// cell==nil (4 componentes, incluindo restituicao_st), preservando
	// FiscalDtRef (o cell em si existe, só o destino não).
	if cell.UFDestino == "" {
		return TaxComponents{
			FiscalDtRef: cell.FiscalDtRef,
			Unknown:     []string{"icms_saida", "difal", "fcp", "pis_cofins", "restituicao_st"},
		}
	}

	p := mustRat(preco)
	out := TaxComponents{FiscalDtRef: cell.FiscalDtRef}

	// restituição: independente da resolução da célula (codtrib/ambiguo) —
	// só depende da UF de destino (já garantida não-vazia acima) e do valor
	// unitário. R_aplicável = R se UF ≠ MG, senão 0.
	if cell.UFDestino == ufOrigemMG {
		s, _ := round2(new(big.Rat))
		out.RestituicaoST = &s
	} else {
		s, _ := round2(zeroIfNil(cell.RestituicaoUnit))
		out.RestituicaoST = &s
	}

	// D1 (Task A4): o portão de célula-não-resolvida (ausente ou ambígua)
	// tem que rodar ANTES de qualquer decisão de ramo, inclusive o ramo ST
	// — antes deste conserto, `isST` era decidido primeiro (CodTrib==60
	// bate mesmo com Ambiguo=true) e uma célula AMBÍGUA com CST 60 devolvia
	// ICMS/DIFAL "0.00" explícito com Unknown vazio: indistinguível do zero
	// LEGÍTIMO do ramo ST resolvido. Uma célula ambígua é uma célula que o
	// resolvedor NÃO conseguiu resolver — nunca escolhe às cegas.
	if cell.CodTrib == nil || cell.Ambiguo {
		out.Unknown = append(out.Unknown, "icms_saida", "difal", "fcp", "pis_cofins")
		return out
	}

	isST := *cell.CodTrib == codTribST

	if isST {
		// a = 0: ICMS de saída e DIFAL são 0 EXPLÍCITO — zero legítimo (o
		// próprio já está no custo), não desconhecido. PIS/COFINS ainda
		// roda com a=0: BASE_PC = MAX(0, P − S). FCP é uma FATIA do ICMS de
		// saída (P × fcp_embutido, ver comentário do campo FCP acima) — no
		// ramo ST o ICMS de saída em si é 0 explícito, então a fatia dele
		// também é 0 explícito, mesma justificativa que já vale para
		// ICMSSaida/Difal neste ramo, nunca desconhecido (Task A4 review
		// finding: FCP ficava nil e sem nome aqui, nem zero nem Unknown).
		zeroICMS, _ := round2(new(big.Rat))
		zeroDifal, _ := round2(new(big.Rat))
		zeroFCP, _ := round2(new(big.Rat))
		out.ICMSSaida = &zeroICMS
		out.Difal = &zeroDifal
		out.FCP = &zeroFCP
		base := new(big.Rat).Sub(p, zeroIfNil(cell.StRetidoEntrada))
		out.PisCofins = pisCofinsFrom(base)
		return out
	}

	// D-37: UF sem linha na tabela legal vira pendência explícita — nunca
	// cai na alíquota interna como aproximação. A2 (d): isso agora vale
	// uniformemente, inclusive para MG — o 18/100 hardcoded saiu; MG lê
	// icms_aliquota_interna como qualquer outra UF (constante de negócio em
	// código é gate de auto-revisão do projeto).
	if cell.AliquotaInterna == nil {
		out.Unknown = append(out.Unknown, "icms_saida", "difal", "fcp", "pis_cofins")
		return out
	}

	// aCusto é a carga interna do destino (valor devido, FCP já embutido
	// onde geral — D-43). O ERP calcula o DIFAL por DIFERENÇA SIMPLES:
	// TGFDIN.ALIQPARADIFAL = ALIQINTDEST − ALIQUOTA, com TIPCALCDIFAL = 0,
	// BASEDIFAL = BASE e PERCPARTDIFAL = 100 em 783/783 itens de notas TOP
	// 305/306 desde 2025. D-41 ("DIFAL pela lei") era sobre o DADO, não sobre
	// o método. Ver .mnfs/MIS-008/evidence/fiscal-2026-08-04/roundR.txt (R5)
	// e o plano docs/superpowers/plans/2026-08-04-fiscal-verdade-tgfaid-plan.md
	// §1.1.
	aCusto := new(big.Rat).Quo(mustRat(*cell.AliquotaInterna), cem)

	var icmsOper, difal *big.Rat
	if cell.UFDestino == ufOrigemMG {
		// Destino = origem é operação INTERNA: o ICMS de saída é a alíquota
		// interna cheia e o DIFAL é 0 EXPLÍCITO (zero legítimo, não
		// desconhecido). A derivação interestadual (ICMS = P×a_inter, DIFAL
		// = P×aCusto − ICMS) só vale quando há duas UFs; aplicada a MG→MG
		// ela devolvia ICMS 7% + DIFAL 11%, soma certa e quebra errada — e é
		// a quebra que aparece na tela. Este ramo NUNCA consulta Origprod
		// (D5): a_inter não entra na conta intra-UF.
		icmsOper = new(big.Rat).Mul(p, aCusto)
		difal = new(big.Rat)
	} else {
		// D5 (Task A4): a_inter (aInterFor) DEPENDE de Origprod para
		// decidir nacional×4 eixos importados versus a quebra por UF —
		// antes deste gate, Origprod nil era tratado como "fora do
		// conjunto importado", ou seja, nacional CONHECIDO por omissão.
		// Origem desconhecida é desconhecida; só bloqueia AQUI, no ramo que
		// efetivamente consulta o fato — o ramo intra-MG acima nunca chama
		// aInterFor e fica intocado.
		if cell.Origprod == nil {
			out.Unknown = append(out.Unknown, "icms_saida", "difal", "fcp", "pis_cofins")
			return out
		}
		// Exibição: ICMS_oper = P × a_inter, DIFAL = P × aCusto − ICMS_oper.
		// Uma fonte (aCusto), duas linhas — nunca dois cálculos
		// independentes.
		aInter := aInterFor(cell)
		icmsOper = new(big.Rat).Mul(p, aInter)
		icmsTotal := new(big.Rat).Mul(p, aCusto)
		difal = new(big.Rat).Sub(icmsTotal, icmsOper)
	}

	icmsOperStr, _ := round2(icmsOper)
	difalStr, _ := round2(difal)
	out.ICMSSaida = &icmsOperStr
	out.Difal = &difalStr

	// D2 (Task A4): AliquotaInterna presente com FcpEmbutido nil é um
	// estado que a fonte real nunca produz — a porta devolve os dois da
	// MESMA linha vigente, sempre juntos (pricing/ports/tax_matrix.go,
	// orders/adapters/pricingtax/reader.go:201-208: FcpEmbutidoPct é uma
	// STRING não-ponteiro dentro do struct não-nil, "0" para UF sem FCP
	// geral, nunca ausente). Antes deste gate, zeroIfNil tratava nil como 0
	// legítimo e fabricava FCP="0.00"/PisCofins a partir de um fato que não
	// foi lido. Escopo estreito: só PIS/COFINS (via aBase=aCusto−fcp) fica
	// desconhecido — ICMSSaida/Difal já foram calculados acima e não
	// dependem de FcpEmbutido, continuam conhecidos.
	if cell.FcpEmbutido == nil {
		out.Unknown = append(out.Unknown, "fcp", "pis_cofins")
		return out
	}

	fcp := new(big.Rat).Quo(mustRat(*cell.FcpEmbutido), cem)

	// D3 (Task A4): 0 ≤ fcp ≤ aCusto é invariante da própria fonte (FCP é
	// uma FATIA da alíquota headline, nunca pode superá-la nem ser
	// negativa). Violação é dado inválido, não entrada exótica — antes
	// deste gate, fcp > aCusto produzia aBase negativo e PIS/COFINS sobre
	// uma base MAIOR que o preço, sem bloquear. Mesmo escopo estreito do
	// D2: só PIS/COFINS fica desconhecido.
	if fcp.Sign() < 0 || fcp.Cmp(aCusto) > 0 {
		out.Unknown = append(out.Unknown, "fcp", "pis_cofins")
		return out
	}

	// aCusto inclui o FCP (a tabela legal embute o FCP na alíquota headline
	// onde ele é geral — D-43).
	// aBase EXCLUI o FCP porque, medido, o ERP abate ICMS + DIFAL da base do
	// PIS/COFINS e DEIXA O FCP DENTRO: 6 itens do RJ com resíduo exatamente
	// igual a −FCP (roundV.txt V2). Uma fonte, duas alíquotas, nunca dois
	// cálculos independentes.
	aBase := new(big.Rat).Sub(aCusto, fcp)

	fcpMoney := new(big.Rat).Mul(p, fcp)
	fcpStr, _ := round2(fcpMoney)
	out.FCP = &fcpStr

	// BASE_PC = MAX(0, P×(1−aBase) − S). Usa o `aBase` exato (não o
	// icms_oper+difal já arredondados) — o mesmo racional exato que já
	// valia para o cálculo, nunca um valor re-arredondado alimentando outro
	// cálculo. aBase exclui o FCP (rule b): o FCP fica DENTRO da base do
	// PIS/COFINS mesmo estando embutido na alíquota usada para ICMS/DIFAL.
	base := new(big.Rat).Sub(big.NewRat(1, 1), aBase)
	base.Mul(base, p)
	base.Sub(base, zeroIfNil(cell.StRetidoEntrada))
	out.PisCofins = pisCofinsFrom(base)

	return out
}

// aInterFor derives a_inter (D-41): origprod ∈ {1,2,3,8} vence a UF; senão UF
// ∈ {SP,RJ,PR,SC,RS} → 12%; demais UFs → 7%. Nunca em si um gatilho de
// desconhecido — sempre derivável dos fatos já presentes na célula.
func aInterFor(cell *ICMSCell) *big.Rat {
	if cell.Origprod != nil && origprodImportado[*cell.Origprod] {
		return big.NewRat(4, 100)
	}
	if aInterUF12[cell.UFDestino] {
		return big.NewRat(12, 100)
	}
	return big.NewRat(7, 100)
}

// pisCofinsFrom applies MAX(0, base) then 9,25% — o MAX não é cosmético: S
// retido pode exceder o líquido, e sem o clamp sai base negativa e crédito
// fantasma.
func pisCofinsFrom(base *big.Rat) *string {
	if base.Sign() < 0 {
		base = big.NewRat(0, 1)
	}
	pc := new(big.Rat).Mul(pisCofinsRate, base)
	s, _ := round2(pc)
	return &s
}

// zeroIfNil treats a nil money string as 0 — S e R são fatos cuja fonte
// Oracle não distingue "genuinamente zero" de "nunca registrado" (diferente
// de CodTrib/AliquotaInterna, cuja ausência É o gatilho de desconhecido).
func zeroIfNil(s *string) *big.Rat {
	if s == nil {
		return big.NewRat(0, 1)
	}
	return mustRat(*s)
}
