package domain

import "math/big"

// Modalidade is the Mercado Livre listing modalidade that selects the
// comissão rate and whether the full tarifa applies. Values are the wire
// tokens IC-04 uses.
type Modalidade string

const (
	ModalidadeClassico Modalidade = "classico"
	ModalidadePremium  Modalidade = "premium"
	ModalidadeFull     Modalidade = "full"
)

// taxaFixaLimiar is the preço threshold (exclusive) below which the fixed
// per-sale tarifa of R$ 6,50 applies; at or above it the fixed tarifa is 0
// and the produto frete is charged instead (IC-04 line 84).
var (
	taxaFixaLimiar = big.NewRat(79, 1)
	taxaFixaValor  = big.NewRat(650, 100) // 6.50
	cem            = big.NewRat(100, 1)
)

// defaultTaxaFixaLimiarCents is taxaFixaLimiar expressed in cents (79,00 ⇒
// 7900). This is the ONLY surviving form of the 79 literal for solve segment
// math (solve.go); it must equal taxaFixaLimiar × 100.
const defaultTaxaFixaLimiarCents int64 = 7900

// DecomposeInput carries the ALREADY-RESOLVED decimal inputs for one
// decomposition. The engine is pure (note G1): comissão pct, aliquota,
// efetivo, frete, custo and tarifa_full are resolved by the ports (S6)
// before this is called, so the golden test stays hermetic. Percentages are
// decimal strings expressed in percent (e.g. "4" = 4%, "6.00" = 6%).
type DecomposeInput struct {
	Preco       string     // > 0, transport-validated (INVALID_PRICE) upstream
	ComissaoPct string     // comissão rate for Modalidade, resolved from the fee schedule
	AliquotaPct string     // CalcProfile regime aliquota
	Modalidade  Modalidade // classico | premium | full

	// TarifaFull is CalcProfile.tarifa_full; consulted only when Modalidade
	// is full. nil there ⇒ tarifa_full UNKNOWN (propagates; never 0).
	TarifaFull *Money
	// FreteProduto is the resolved produto frete; consulted only when
	// preço ≥ 79. nil there ⇒ frete UNKNOWN. Below 79 frete is 0 and this
	// is not consulted.
	FreteProduto *Money
	// Custo is the resolved custo_erp; always applied. nil ⇒ custo UNKNOWN.
	Custo *Money

	// DifalEnabled false ⇒ difal is 0 explicit (feature off, legit zero).
	// true with a known DestinoUF + EfetivoPct ⇒ difal = efetivo × preço.
	// true with an unknown destino ⇒ difal UNKNOWN.
	//
	// Consulted ONLY when ICMSCell is nil. When ICMSCell is present, Difal
	// (and ICMSSaida/PisCofins/RestituicaoST) come exclusively from
	// TaxesForItem (D-41) — these three fields are not read at all in that
	// case.
	DifalEnabled bool
	DestinoUF    string // "" = destino unknown
	EfetivoPct   string // DifalForUF efetivo for DestinoUF; "" when not applied

	// ICMSCell is the P2.b Task 4 tax-matrix cell + product fiscal facts,
	// already resolved by the port (S6-style) before Decompose runs — nil
	// means the caller has no ICMS matrix data for this item (today's
	// behavior: Decompose falls back to the DifalEnabled/DestinoUF/EfetivoPct
	// path above and AliquotaPct stays summed as Imposto, byte-for-byte
	// identical to pre-P2.b output). Non-nil switches Decompose onto the D-41
	// formula: ICMSSaida/Difal/PisCofins/RestituicaoST are computed via
	// TaxesForItem and enter the sum in ICMSSaida/Difal/PisCofins/
	// RestituicaoST's place — Imposto is still computed and shown (D-38: the
	// field survives for the 7 simulator sites) but stops being summed, so it
	// is never double-counted alongside ICMSSaida.
	ICMSCell *ICMSCell
}

// Decomposition is the frozen IC-04 output shape (line 85). Always-known
// components are plain decimal strings at 2dp; components that can be UNKNOWN
// (frete, difal, tarifa_full, custo) are *string — nil means unknown, never
// "0.00". When ANY component is unknown MargemValor/MargemPct are nil and
// ComponentesDesconhecidos names the missing (ADR-17: unknown ≠ zero). When
// every component is known, preço = Σ(componentes) + margem_valor EXACTLY.
type Decomposition struct {
	Preco    string
	Comissao string
	TaxaFixa string
	Frete    *string
	// Imposto is the CalcProfile regime aliquota × preço (SIMPLES/PRESUMIDO
	// placeholder, D-38). Always computed and shown for the simulator sites
	// that still read it. Summed into MargemValor ONLY when ICMSCell is nil
	// (today's behavior); when ICMSCell is present, ICMSSaida replaces it in
	// the sum — Imposto keeps rendering but a silently-double-counted
	// component is exactly how the fabricated 4% survived this long.
	Imposto    string
	Difal      *string
	TarifaFull *string
	Custo      *string
	// ICMSSaida is P × a_inter (D-41) — nil only when ICMSCell is nil OR the
	// matrix cell/alíquota interna could not be resolved (ADR-17).
	ICMSSaida *string
	// PisCofins is 0,0925 × MAX(0, P×(1−a) − S) (STF Tema 69 + STJ Tema
	// 1125). Same nil rule as ICMSSaida — both come from the same resolved
	// `a`.
	PisCofins *string
	// RestituicaoST is the ICMS-ST restituição credit — the ONE component
	// that ADDS to margem instead of subtracting (Task 6: "a única linha
	// positiva da seção"). 0 explicit for UF=MG (intra-empresa, not
	// unknown); the resolved value otherwise. nil only when ICMSCell itself
	// is nil.
	RestituicaoST *string
	MargemValor   *string
	MargemPct     *string

	ComponentesDesconhecidos []string
}

// Decompose applies the single IC-04 formula. Inputs are transport-validated;
// a decimal parse failure is a programmer error upstream and panics (there is
// no error in the frozen port sig). The difal money component is computed from
// the resolved EfetivoPct string here — but that string itself is derived from
// exact interna−interestadual via big.Rat in DifalForUF, so no rounded value
// is ever re-rounded into a cost.
func Decompose(in DecomposeInput) Decomposition {
	return decomposeWithLimiar(in, taxaFixaLimiar)
}

// decomposeWithLimiar is the IC-04 formula parameterized on the taxa_fixa/
// frete step threshold. Decompose calls it with the default taxaFixaLimiar
// (79); solve.go's per-segment math calls it directly with an override so
// the frozen Decompose signature (DecomposeInput) Decomposition never
// changes shape.
func decomposeWithLimiar(in DecomposeInput, limiar *big.Rat) Decomposition {
	preco := mustRat(in.Preco)

	comissao, comissaoRat := pctOfPrice(in.ComissaoPct, preco)
	imposto, impostoRat := pctOfPrice(in.AliquotaPct, preco)

	var taxaFixa string
	var taxaFixaRat *big.Rat
	if preco.Cmp(limiar) < 0 {
		taxaFixa, taxaFixaRat = round2(new(big.Rat).Set(taxaFixaValor))
	} else {
		taxaFixa, taxaFixaRat = round2(new(big.Rat))
	}

	out := Decomposition{
		Preco:    FormatRatHalfUp(preco, 2),
		Comissao: comissao,
		TaxaFixa: taxaFixa,
		Imposto:  imposto,
	}

	// running sum of the KNOWN component rats (for the exact soma-fecha).
	// impostoRat joins the sum only when ICMSCell is nil (D-38/D-41): when
	// the matrix cell is present, ICMSSaida takes Imposto's place below —
	// summing both would double-count the same tax.
	sum := new(big.Rat).Add(comissaoRat, taxaFixaRat)
	if in.ICMSCell == nil {
		sum.Add(sum, impostoRat)
	}
	var unknown []string

	// frete: applied only when preço ≥ limiar.
	if preco.Cmp(limiar) < 0 {
		s, r := round2(new(big.Rat))
		out.Frete = &s
		sum.Add(sum, r)
	} else if in.FreteProduto == nil {
		unknown = append(unknown, "frete")
	} else {
		s, r := round2(mustRat(in.FreteProduto.Amount))
		out.Frete = &s
		sum.Add(sum, r)
	}

	if in.ICMSCell == nil {
		// difal (old path): off ⇒ 0 explicit; on+destino known ⇒ efetivo ×
		// preço; on+destino unknown ⇒ UNKNOWN.
		if !in.DifalEnabled {
			s, r := round2(new(big.Rat))
			out.Difal = &s
			sum.Add(sum, r)
		} else if in.DestinoUF == "" || in.EfetivoPct == "" {
			unknown = append(unknown, "difal")
		} else {
			s, r := pctOfPrice(in.EfetivoPct, preco)
			out.Difal = &s
			sum.Add(sum, r)
		}
	} else {
		// D-41: ICMSSaida/Difal/PisCofins/RestituicaoST all come from the
		// single collapsed rate `a` — "uma fonte, duas linhas" for
		// ICMSSaida/Difal, plus PisCofins and the RestituicaoST credit.
		// RestituicaoST is the one component that SUBTRACTS from sum
		// (ADDS to margem) — it is a credit back to the seller, not a cost.
		tax := TaxesForItem(in.Preco, in.ICMSCell)
		out.ICMSSaida = tax.ICMSSaida
		out.Difal = tax.Difal
		out.PisCofins = tax.PisCofins
		out.RestituicaoST = tax.RestituicaoST

		if tax.ICMSSaida == nil {
			unknown = append(unknown, "icms_saida")
		} else {
			sum.Add(sum, mustRat(*tax.ICMSSaida))
		}
		if tax.Difal == nil {
			unknown = append(unknown, "difal")
		} else {
			sum.Add(sum, mustRat(*tax.Difal))
		}
		if tax.PisCofins == nil {
			unknown = append(unknown, "pis_cofins")
		} else {
			sum.Add(sum, mustRat(*tax.PisCofins))
		}
		if tax.RestituicaoST == nil {
			unknown = append(unknown, "restituicao_st")
		} else {
			sum.Sub(sum, mustRat(*tax.RestituicaoST))
		}
	}

	// tarifa_full: full ⇒ CalcProfile.tarifa_full (nil ⇒ UNKNOWN); other
	// modalidades ⇒ 0 explicit.
	if in.Modalidade == ModalidadeFull {
		if in.TarifaFull == nil {
			unknown = append(unknown, "tarifa_full")
		} else {
			s, r := round2(mustRat(in.TarifaFull.Amount))
			out.TarifaFull = &s
			sum.Add(sum, r)
		}
	} else {
		s, r := round2(new(big.Rat))
		out.TarifaFull = &s
		sum.Add(sum, r)
	}

	// custo_erp: always applied. nil ⇒ UNKNOWN.
	if in.Custo == nil {
		unknown = append(unknown, "custo_erp")
	} else {
		s, r := round2(mustRat(in.Custo.Amount))
		out.Custo = &s
		sum.Add(sum, r)
	}

	if len(unknown) > 0 {
		out.ComponentesDesconhecidos = unknown
		return out // margem stays nil — never fabricate 0
	}

	// soma fecha exactly: every summed rat is 2dp and preço is 2dp, so
	// margem = preço − Σ is exactly 2dp (no rounding drift).
	margem := new(big.Rat).Sub(preco, sum)
	mv := FormatRatHalfUp(margem, 2)
	out.MargemValor = &mv

	// margem_pct = margem / preço × 100 (2dp). preço > 0 upstream.
	pct := new(big.Rat).Quo(margem, preco)
	pct.Mul(pct, cem)
	mp := FormatRatHalfUp(pct, 2)
	out.MargemPct = &mp
	return out
}

// pctOfPrice returns pct% of price as (2dp string, exact 2dp rat).
func pctOfPrice(pct string, price *big.Rat) (string, *big.Rat) {
	r := new(big.Rat).Quo(mustRat(pct), cem)
	r.Mul(r, price)
	return round2(r)
}

// round2 formats r at 2dp and returns the string plus the exact 2dp rat it
// parses back to, so summed components and the formatted output agree bit for
// bit. FormatRatHalfUp is the single rounding routine (G1).
func round2(r *big.Rat) (string, *big.Rat) {
	s := FormatRatHalfUp(r, 2)
	return s, mustRat(s)
}

func mustRat(s string) *big.Rat {
	r, err := ParseRat(s)
	if err != nil {
		panic(err)
	}
	return r
}
