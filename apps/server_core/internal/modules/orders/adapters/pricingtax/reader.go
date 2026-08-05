// Package pricingtax adapts the pricing module's D-41 per-item ICMS matrix
// calculation (pricing/domain.TaxesForItem, Task 4) to the orders module's
// TaxReader port, so a sold order's ICMSSaida/Difal/PisCofins/RestituicaoST
// come from exactly the matrix cell + product fiscal facts the Simulador
// prices with — one tax truth, not a second copy that can drift.
package pricingtax

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"

	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// ufOrigemMG is the company's fixed origin UF (D-17 debt: CODEMP=1 fixo,
// origin is not resolved per tenant/company today). Deliberately NOT shared
// with pricing/domain's own unexported ufOrigemMG — same rationale as that
// package's comment: two different packages, duplicating a stable literal
// beats an import coupling two modules just to share one constant.
const ufOrigemMG = "MG"

// Reader implements ordersports.TaxReader over the D-41 per-item path: for
// each order item it resolves the product's fiscal facts (grupo_icms,
// origprod, st_retido_entrada, restituicao_unit) and the ICMS matrix cell for
// (ufOrigem, destinoUF, grupo_icms), builds a pricing/domain.ICMSCell, and
// calls TaxesForItem once per item — never once for the whole order total
// (D-41's non-linear MAX(0,...) and gross-up formulas do not commute with
// rateio-after-the-fact).
type Reader struct {
	matrix   pricingports.TaxMatrixReader
	products pricingports.ProductFiscalReader
	tenantID string
}

// NewReader builds a Reader over matrix (Task 4's TaxMatrixReader) and
// products (this task's ProductFiscalReader).
func NewReader(matrix pricingports.TaxMatrixReader, products pricingports.ProductFiscalReader, tenantID string) *Reader {
	return &Reader{matrix: matrix, products: products, tenantID: tenantID}
}

var _ ordersports.TaxReader = (*Reader)(nil)

// TaxesForItems resolves and sums ICMSSaida/Difal/PisCofins/RestituicaoST
// across items, per-component all-or-nothing (ADR-17): if ANY item's
// TaxesForItem call comes back unknown for a component, that component is
// nil for the WHOLE ORDER, never a partial sum over only the resolved items.
// An item with no product link, or whose product has no fiscal snapshot,
// resolves to a nil ICMSCell for that item — which pricing/domain.TaxesForItem
// already treats as "every component unknown", so the all-or-nothing rule
// above naturally makes the whole order's aggregate honest without a second,
// separate branch for the "unlinked item" case.
//
// destinoUF == "" means the order has no known destination — there is no
// honest cell to resolve for any item, so every component comes back nil.
func (r *Reader) TaxesForItems(ctx context.Context, items []ordersports.TaxItem, destinoUF string) (ordersports.OrderTaxes, error) {
	if destinoUF == "" {
		return ordersports.OrderTaxes{}, nil
	}

	// aliquota interna depends only on destinoUF, not on the item — resolve
	// it once for the whole order instead of once per item.
	aliquotaInterna, err := r.matrix.AliquotaInternaFor(ctx, destinoUF)
	if err != nil {
		return ordersports.OrderTaxes{}, err
	}

	icmsKnown, difalKnown, pisKnown, restKnown := true, true, true, true
	icmsSum := new(big.Rat)
	difalSum := new(big.Rat)
	pisSum := new(big.Rat)
	restSum := new(big.Rat)

	for _, item := range items {
		cell, preco, err := r.resolveItem(ctx, destinoUF, aliquotaInterna, item)
		if err != nil {
			return ordersports.OrderTaxes{}, err
		}

		tax := pricingdomain.TaxesForItem(preco, cell)

		if tax.ICMSSaida == nil {
			icmsKnown = false
		} else if icmsKnown {
			v, err := parseTaxComponent(*tax.ICMSSaida, "ICMS saída")
			if err != nil {
				return ordersports.OrderTaxes{}, err
			}
			icmsSum.Add(icmsSum, v)
		}
		if tax.Difal == nil {
			difalKnown = false
		} else if difalKnown {
			v, err := parseTaxComponent(*tax.Difal, "DIFAL")
			if err != nil {
				return ordersports.OrderTaxes{}, err
			}
			difalSum.Add(difalSum, v)
		}
		if tax.PisCofins == nil {
			pisKnown = false
		} else if pisKnown {
			v, err := parseTaxComponent(*tax.PisCofins, "PIS/COFINS")
			if err != nil {
				return ordersports.OrderTaxes{}, err
			}
			pisSum.Add(pisSum, v)
		}
		if tax.RestituicaoST == nil {
			restKnown = false
		} else if restKnown {
			v, err := parseTaxComponent(*tax.RestituicaoST, "restituição de ST")
			if err != nil {
				return ordersports.OrderTaxes{}, err
			}
			restSum.Add(restSum, v)
		}
	}

	var out ordersports.OrderTaxes
	if icmsKnown {
		v, err := ratToFloat(icmsSum)
		if err != nil {
			return ordersports.OrderTaxes{}, err
		}
		out.ICMSSaida = &v
	}
	if difalKnown {
		v, err := ratToFloat(difalSum)
		if err != nil {
			return ordersports.OrderTaxes{}, err
		}
		out.Difal = &v
	}
	if pisKnown {
		v, err := ratToFloat(pisSum)
		if err != nil {
			return ordersports.OrderTaxes{}, err
		}
		out.PisCofins = &v
	}
	if restKnown {
		v, err := ratToFloat(restSum)
		if err != nil {
			return ordersports.OrderTaxes{}, err
		}
		out.RestituicaoST = &v
	}
	return out, nil
}

// resolveItem builds the ICMSCell + preço (line total, UnitPrice × Quantity)
// for one order item. A nil *ICMSCell — no product link, no UnitPrice, no
// readable price, or no products_mirror fiscal snapshot for the linked product
// — is the single, uniform "this item's tax picture is unknown" signal that
// TaxesForItem already understands; the caller (TaxesForItems) never
// special-cases the unlinked case separately from a linked-but-factless one.
func (r *Reader) resolveItem(ctx context.Context, destinoUF string, aliquotaInterna *pricingports.AliquotaInterna, item ordersports.TaxItem) (*pricingdomain.ICMSCell, string, error) {
	if item.InternalProductID == nil || item.UnitPrice == nil {
		return nil, "0.00", nil
	}

	// A price we cannot read faithfully is not a price: it takes the same
	// named-unknown channel as an unlinked item, never a silently rounded
	// number. Checked before any read, so an unusable price costs nothing.
	preco, priceReadable := lineTotal(*item.UnitPrice, item.Quantity)
	if !priceReadable {
		return nil, "0.00", nil
	}

	codigoProduto := strconv.Itoa(*item.InternalProductID)
	facts, err := r.products.FactsFor(ctx, r.tenantID, codigoProduto)
	if err != nil {
		return nil, "", err
	}
	if !facts.Found {
		return nil, "0.00", nil
	}

	// S/R are per-unit product facts (products_mirror.st_retido_entrada /
	// restituicao_unit) — scale by Quantity BEFORE building the ICMSCell, so
	// TaxesForItem's formula (which expects the line's own S/R, not a
	// per-unit fact) sees the right magnitude.
	stScaled, err := scaleMoney(facts.StRetidoEntrada, item.Quantity)
	if err != nil {
		return nil, "", err
	}
	restScaled, err := scaleMoney(facts.RestituicaoUnit, item.Quantity)
	if err != nil {
		return nil, "", err
	}

	var codTrib *int
	var ambiguo bool
	if facts.GrupoICMS != nil {
		mc, err := r.matrix.CellFor(ctx, r.tenantID, ufOrigemMG, destinoUF, *facts.GrupoICMS)
		if err != nil {
			return nil, "", err
		}
		if mc.Found {
			codTrib = mc.CodTrib
			ambiguo = mc.Ambiguo
		}
	}

	// aliquotaInterna carries BOTH facts (headline + fcp_embutido) from the
	// same vigente row (Task A2) — nil means "UF sem linha" (D-37) and both
	// ICMSCell fields stay nil together, never one without the other.
	var aliquotaPct, fcpEmbutidoPct *string
	if aliquotaInterna != nil {
		aliquotaPct = &aliquotaInterna.AliquotaPct
		fcpEmbutidoPct = &aliquotaInterna.FcpEmbutidoPct
	}

	cell := &pricingdomain.ICMSCell{
		UFDestino:       destinoUF,
		CodTrib:         codTrib,
		Ambiguo:         ambiguo,
		Origprod:        facts.Origprod,
		AliquotaInterna: aliquotaPct,
		FcpEmbutido:     fcpEmbutidoPct,
		StRetidoEntrada: stScaled,
		RestituicaoUnit: restScaled,
		FiscalDtRef:     facts.FiscalDtRef,
	}

	return cell, preco, nil
}

// lineTotal converts an order item's UnitPrice (float64, the orders domain's
// money representation) × Quantity into the 2dp decimal string TaxesForItem
// expects — the one float64→decimal crossing in this adapter, same idiom the
// pre-T5 reader used for order.Total.
//
// It reports ok=false when unitPrice is not a readable amount, so the caller
// can route it to the named-unknown channel instead of a number. Before this
// guard the conversion rounded to 2dp in SILENCE: lineTotal(1.005, 1) produced
// "1.00" where the decimal-exact answer is "1.01", and NaN/±Inf produced
// "0.00" — a fabricated zero, the ADR-17 failure mode. That was only ever safe
// by accident of schema (unit_price is numeric(14,2), migrations/0027:31), a
// premise no code asserted; ingest_service.go:241 fills UnitPrice straight from
// the provider payload, in memory, before that column ever rounds it.
//
// The predicate is NOT "exactly representable as n/100" — that would reject
// almost every real price, since float64's nearest value to 0.10 is
// 0.1000000000000000055…, so 0.10 × 100 is not an integer. The question is the
// other direction: is this float the faithful IMAGE of some 2dp decimal? Round
// it, take the decimal back to float64, compare. 0.10 survives (it IS the
// nearest double to "0.10"); 1.005 does not (its nearest double rounds to
// "1.00", a different double). The line total is then derived from that
// DECIMAL, never from the binary value it came in as.
func lineTotal(unitPrice float64, quantity int) (string, bool) {
	if math.IsNaN(unitPrice) || math.IsInf(unitPrice, 0) {
		return "", false
	}
	v := new(big.Rat).SetFloat64(unitPrice)
	if v == nil {
		return "", false
	}

	unit := pricingdomain.FormatRatHalfUp(v, 2)
	back, err := strconv.ParseFloat(unit, 64)
	if err != nil || back != unitPrice {
		return "", false
	}

	total, ok := new(big.Rat).SetString(unit)
	if !ok {
		return "", false
	}
	total.Mul(total, big.NewRat(int64(quantity), 1))
	return pricingdomain.FormatRatHalfUp(total, 2), true
}

// scaleMoney multiplies a decimal money string by qty, or passes nil through
// unscaled — nil already means "not registered" (zeroIfNil's job downstream,
// pricing/domain), scaling nil would still be nil, no need to touch it.
func scaleMoney(s *string, qty int) (*string, error) {
	if s == nil {
		return nil, nil
	}
	v, err := pricingdomain.ParseRat(*s)
	if err != nil {
		return nil, err
	}
	v.Mul(v, big.NewRat(int64(qty), 1))
	out := pricingdomain.FormatRatHalfUp(v, 2)
	return &out, nil
}

// ratToFloat renders r at 2dp through the pricing module's single rounding
// routine (so orders and Simulador round identically) and parses it back to
// float64 for the OrderTaxes wire shape.
func ratToFloat(r *big.Rat) (float64, error) {
	return strconv.ParseFloat(pricingdomain.FormatRatHalfUp(r, 2), 64)
}

// parseTaxComponent lê uma saída de TaxesForItem de volta para big.Rat. Essas
// strings são sempre produto do próprio FormatRatHalfUp, então uma falha aqui é
// erro de programação em pricing/domain, não dado ruim do operador — mas ela
// vira erro, não panic: o que não calcula é UM pedido, e derrubar o processo do
// servidor por isso troca um número faltando por um serviço fora do ar. O nome
// do componente entra na mensagem porque um erro que diz apenas "parse falhou"
// não diz qual dos quatro impostos quebrou.
func parseTaxComponent(s, componente string) (*big.Rat, error) {
	r, err := pricingdomain.ParseRat(s)
	if err != nil {
		return nil, fmt.Errorf("orders pricingtax: %s inválido (%q): %w", componente, s, err)
	}
	return r, nil
}
