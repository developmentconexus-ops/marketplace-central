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
// key off it: (1) a intra-MG usa a=0,18 FIXO, sem consultar
// icms_aliquota_interna; (2) restituição de ICMS-ST só sai de operação
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

	// Origprod is products_mirror.origprod (TGFPRO.ORIGPROD). nil is treated
	// the same as any value outside {1,2,3,8} — origprod nunca é ele próprio
	// o gatilho de desconhecido; a_inter só cai para o ramo por UF.
	Origprod *int
	// AliquotaInterna is icms_aliquota_interna.aliquota (FCP já embutido onde
	// geral, D-43) para UFDestino, como string decimal em percentual ("18" =
	// 18%). nil ⇒ UF sem linha na tabela legal (D-37) — só é consultado fora
	// de MG e fora do ramo ST; MG usa 0,18 fixo e ST usa a=0, nenhum dos dois
	// toca este campo.
	AliquotaInterna *string

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
	ICMSSaida     *string
	Difal         *string
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
		return TaxComponents{Unknown: []string{"icms_saida", "difal", "pis_cofins", "restituicao_st"}}
	}

	p := mustRat(preco)
	out := TaxComponents{FiscalDtRef: cell.FiscalDtRef}

	// restituição: independente da resolução da célula (codtrib/ambiguo) —
	// só depende da UF de destino e do valor unitário. R_aplicável = R se
	// UF ≠ MG, senão 0.
	if cell.UFDestino == ufOrigemMG {
		s, _ := round2(new(big.Rat))
		out.RestituicaoST = &s
	} else {
		s, _ := round2(zeroIfNil(cell.RestituicaoUnit))
		out.RestituicaoST = &s
	}

	isST := cell.CodTrib != nil && *cell.CodTrib == codTribST

	if isST {
		// a = 0: ICMS de saída e DIFAL são 0 EXPLÍCITO — zero legítimo (o
		// próprio já está no custo), não desconhecido. PIS/COFINS ainda
		// roda com a=0: BASE_PC = MAX(0, P − S).
		zeroICMS, _ := round2(new(big.Rat))
		zeroDifal, _ := round2(new(big.Rat))
		out.ICMSSaida = &zeroICMS
		out.Difal = &zeroDifal
		base := new(big.Rat).Sub(p, zeroIfNil(cell.StRetidoEntrada))
		out.PisCofins = pisCofinsFrom(base)
		return out
	}

	// célula não resolvida (ausente ou ambígua): nunca escolhe às cegas.
	if cell.CodTrib == nil || cell.Ambiguo {
		out.Unknown = append(out.Unknown, "icms_saida", "difal", "pis_cofins")
		return out
	}

	aInter := aInterFor(cell)

	var a *big.Rat
	if cell.UFDestino == ufOrigemMG {
		// intra-MG: a é FIXO, nunca consulta icms_aliquota_interna.
		a = big.NewRat(18, 100)
	} else {
		if cell.AliquotaInterna == nil {
			// D-37: UF sem linha na tabela legal vira pendência explícita —
			// nunca cai na alíquota interna como aproximação.
			out.Unknown = append(out.Unknown, "icms_saida", "difal", "pis_cofins")
			return out
		}
		aInt := new(big.Rat).Quo(mustRat(*cell.AliquotaInterna), cem)
		one := big.NewRat(1, 1)
		num := new(big.Rat).Mul(aInt, new(big.Rat).Sub(one, aInter))
		den := new(big.Rat).Sub(one, aInt)
		a = new(big.Rat).Quo(num, den)
	}

	// Exibição: ICMS_oper = P × a_inter, DIFAL = P × a − ICMS_oper. Uma
	// fonte (a), duas linhas — nunca dois cálculos independentes.
	icmsOper := new(big.Rat).Mul(p, aInter)
	icmsTotal := new(big.Rat).Mul(p, a)
	difal := new(big.Rat).Sub(icmsTotal, icmsOper)

	icmsOperStr, _ := round2(icmsOper)
	difalStr, _ := round2(difal)
	out.ICMSSaida = &icmsOperStr
	out.Difal = &difalStr

	// BASE_PC = MAX(0, P×(1−a) − S). Usa o `a` exato (não o icms_oper+difal
	// já arredondados) — o mesmo racional exato que já vale para o cálculo,
	// nunca um valor re-arredondado alimentando outro cálculo.
	base := new(big.Rat).Sub(big.NewRat(1, 1), a)
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
