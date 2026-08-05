package domain

import (
	"math/big"
	"testing"
)

// TestPisCofinsRatesSumToTheSolverAsymptoteRate fecha um acoplamento que a
// Task A8 deixou vivo e sem instrumento.
//
// `solve.go`'s icmsCellAsymptoticRatePct é uma re-derivação À MÃO do limite de
// TaxesForItem quando preço→∞. No limite os arredondamentos lavam, então ela
// usa a alíquota COMBINADA (`pisCofinsRate`), enquanto o caminho de dinheiro
// real passou a usar `pisRate` e `cofinsRate` separados. Isso continua exato —
// mas SÓ enquanto pisRate + cofinsRate for exatamente pisCofinsRate.
//
// O contrato de acoplamento em solve.go:197-211 manda re-derivar "whenever
// aCusto, aBase/FCP folding, isST's zeroing, or pisCofinsRate change". Repare
// no buraco: mexer só em `pisRate` ou só em `cofinsRate` NÃO muda
// `pisCofinsRate` e portanto não dispara aquele gatilho — a assíntota
// divergiria calada da fórmula real, que é literalmente a "Finding-1 class"
// que aquele comentário diz querer evitar. Prosa não pega isso; este teste
// pega.
func TestPisCofinsRatesSumToTheSolverAsymptoteRate(t *testing.T) {
	sum := new(big.Rat).Add(pisRate, cofinsRate)
	if sum.Cmp(pisCofinsRate) != 0 {
		t.Fatalf("pisRate + cofinsRate = %s, mas pisCofinsRate = %s — a assíntota do solver (solve.go:215) usa a combinada e passa a divergir da fórmula real; ou corrija as alíquotas ou re-derive icmsCellAsymptoticRatePct",
			sum.RatString(), pisCofinsRate.RatString())
	}

	// Controle do próprio instrumento: uma asserção de igualdade que não
	// consegue DISTINGUIR não vale nada. Perturbando um lado em 1 ponto-base,
	// a comparação acima tem que reprovar — senão ela estaria passando por
	// vacuidade, não por acerto.
	perturbado := new(big.Rat).Add(pisRate, new(big.Rat).Add(cofinsRate, big.NewRat(1, 10000)))
	if perturbado.Cmp(pisCofinsRate) == 0 {
		t.Fatal("a comparação não distingue uma perturbação de 1 ponto-base — o instrumento é vácuo")
	}
}

// Task A8 — PIS e COFINS são DOIS tributos, arredondados SEPARADAMENTE.
//
// RED escrito ANTES da implementação e contra a árvore pristina, por quem não
// vai implementar (mandato do operador). Hoje `pisCofinsFrom` (icms.go:332-342)
// aplica uma alíquota combinada 9,25% e arredonda UMA vez. O ERP não faz isso:
// PIS 1,65% e COFINS 7,60% são apurados e arredondados cada um, e só depois
// somados.
//
// Matriz medida contra o ERP (base = 248,917):
//
//	ERP    separado  sobre base arredondada -> 4.11 + 18.92 = 23.03
//	       combinado sobre base arredondada -> 23.0251      -> 23.03
//	       separado  sobre base exata       -> 4.11 + 18.92 = 23.03
//	nosso  combinado sobre base exata       -> 23.0248225   -> 23.02
//
// Um centavo, e ele é sistemático — não é ruído de fixture.
//
// LIMITE HONESTO DESTE VETOR: as três primeiras linhas da matriz concordam em
// 23.03, então este teste distingue APENAS "separado vs combinado". Ele NÃO
// decide se a base entra arredondada ou exata na conta — nenhum vetor medido
// até agora separa essas duas hipóteses. Se um caso futuro separar, ele precisa
// de teste próprio; não leia este arquivo como se a questão da base estivesse
// fechada.
//
// O ramo ST é o sítio mais direto: lá a base é literalmente P − S
// (icms.go:204), sem alíquota efetiva no meio, então o número abaixo é
// atribuível só ao arredondamento de PIS/COFINS.
func TestTaxesForItemPisCofinsRoundedSeparately(t *testing.T) {
	cell := &ICMSCell{
		UFDestino: "BA",
		CodTrib:   intptr(codTribST),
	}

	got := TaxesForItem("248.917", cell)

	if got.PisCofins == nil {
		t.Fatalf("PisCofins = nil, want 23.03 — Unknown=%v", got.Unknown)
	}
	// PIS    = round2(248.917 x 0.0165) = round2(4.10713205)  = 4.11
	// COFINS = round2(248.917 x 0.0760) = round2(18.9176920)  = 18.92
	// soma   = 23.03      (combinado hoje: round2(23.0248225) = 23.02)
	if *got.PisCofins != "23.03" {
		t.Errorf("PisCofins = %s, want 23.03 (PIS 4.11 + COFINS 18.92, arredondados separadamente)", *got.PisCofins)
	}
}

// TestTaxesForItemPisCofinsSeparateRoundingAgreesWhereItShould é o controle que
// impede a mudança de virar um deslocamento cego: nos casos em que as duas
// formas concordam, o número NÃO pode mexer. Se um destes mudar, a alíquota
// composta saiu do lugar (1,65 + 7,60 tem que continuar valendo 9,25), e isso
// é um defeito diferente do que A8 conserta.
func TestTaxesForItemPisCofinsSeparateRoundingAgreesWhereItShould(t *testing.T) {
	for _, tc := range []struct{ preco, want string }{
		{"100.00", "9.25"},  // 1.65 + 7.60 exatos
		{"10.00", "0.93"},   // 0.17 + 0.76 ; combinado round2(0.925) = 0.93
		{"30.00", "2.78"},   // 0.50 + 2.28 ; combinado round2(2.775) = 2.78
		{"1000.00", "92.50"}, // 16.50 + 76.00 exatos
	} {
		cell := &ICMSCell{UFDestino: "BA", CodTrib: intptr(codTribST)}
		got := TaxesForItem(tc.preco, cell)
		if got.PisCofins == nil {
			t.Errorf("preço %s: PisCofins = nil, want %s — Unknown=%v", tc.preco, tc.want, got.Unknown)
			continue
		}
		if *got.PisCofins != tc.want {
			t.Errorf("preço %s: PisCofins = %s, want %s (separado e combinado concordam aqui — este valor não pode mexer)", tc.preco, *got.PisCofins, tc.want)
		}
	}
}

// TestTaxesForItemPisCofinsClampSurvivesSeparateRounding fixa o MAX(0, base)
// (icms.go:336-338): S retido pode exceder o líquido, e sem o clamp sai base
// negativa e crédito fantasma. A separação em dois tributos não pode
// ressuscitar isso por um dos lados — o clamp é sobre a BASE, uma vez, antes
// de qualquer alíquota, não duas vezes depois.
func TestTaxesForItemPisCofinsClampSurvivesSeparateRounding(t *testing.T) {
	cell := &ICMSCell{
		UFDestino:       "BA",
		CodTrib:         intptr(codTribST),
		StRetidoEntrada: strptr("500.00"), // S > P ⇒ base negativa antes do clamp
	}

	got := TaxesForItem("100.00", cell)

	if got.PisCofins == nil {
		t.Fatalf("PisCofins = nil, want 0.00 — Unknown=%v", got.Unknown)
	}
	if *got.PisCofins != "0.00" {
		t.Errorf("PisCofins = %s, want 0.00 — base P−S é negativa e o clamp é 0 explícito, nunca crédito fantasma", *got.PisCofins)
	}
}

func intptr(i int) *int { return &i }
