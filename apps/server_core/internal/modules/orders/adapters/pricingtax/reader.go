// Package pricingtax adapts the pricing module's tenant calculation profile to
// the orders module's TaxReader port, so a sold order's imposto and DIFAL come
// from exactly the same configuration the Simulador prices with — one tax
// truth, not a second copy that can drift.
package pricingtax

import (
	"context"
	"errors"
	"math/big"
	"strconv"

	ordersports "marketplace-central/apps/server_core/internal/modules/orders/ports"
	pricingdomain "marketplace-central/apps/server_core/internal/modules/pricing/domain"
	pricingports "marketplace-central/apps/server_core/internal/modules/pricing/ports"
)

// cem is the percent divisor. Kept as an exact rational so a 4% aliquota on a
// R$ 129,90 order is 5,196 → 5,20, never a binary-float 5,1959999.
var cem = big.NewRat(100, 1)

// ProfileSource is the slice of the pricing CalcRepository this adapter needs:
// the tenant profile (regime aliquota, DIFAL toggle) and the per-UF DIFAL rate.
type ProfileSource interface {
	GetProfile(ctx context.Context, tenantID string) (pricingdomain.CalcProfile, error)
	RateForUF(ctx context.Context, tenantID, uf string) (pricingdomain.DifalRate, error)
}

type Reader struct {
	source   ProfileSource
	tenantID string
}

func NewReader(source ProfileSource, tenantID string) *Reader {
	return &Reader{source: source, tenantID: tenantID}
}

var _ ordersports.TaxReader = (*Reader)(nil)

// TaxesForOrder applies the tenant profile to the order's revenue.
//
// Imposto is the regime aliquota on the order total — always resolvable, since
// an unset profile still answers with its explicit defaults (SIMPLES 4%,
// origem "default"), which is a real configured rate rather than a guess.
//
// DIFAL follows the same three-way rule the pricing engine uses: switched off
// is an explicit 0; switched on with a known destino state is efetivo × total;
// switched on with no destino (or no rate row for that state) is honest-unknown
// (nil), never 0 — an order we cannot place is not an order that owes nothing.
func (r *Reader) TaxesForOrder(ctx context.Context, total float64, destinoUF string) (ordersports.OrderTaxes, error) {
	profile, err := r.source.GetProfile(ctx, r.tenantID)
	if err != nil {
		return ordersports.OrderTaxes{}, err
	}

	revenue := new(big.Rat).SetFloat64(total)
	if revenue == nil { // total was NaN or ±Inf — no honest tax on a non-amount
		return ordersports.OrderTaxes{}, nil
	}

	var taxes ordersports.OrderTaxes
	imposto, err := pctOf(profile.AliquotaPct, revenue)
	if err != nil {
		return ordersports.OrderTaxes{}, err
	}
	taxes.Imposto = &imposto

	if !profile.DifalEnabled {
		zero := 0.0
		taxes.Difal = &zero
		return taxes, nil
	}

	uf := destinoUF
	if uf == "" && profile.DifalDestinoUF != nil {
		uf = *profile.DifalDestinoUF
	}
	if uf == "" {
		return taxes, nil
	}
	rate, err := r.source.RateForUF(ctx, r.tenantID, uf)
	if err != nil {
		if errors.Is(err, pricingports.ErrDifalUFNotFound) {
			return taxes, nil
		}
		return ordersports.OrderTaxes{}, err
	}
	difal, err := pctOf(pricingdomain.DifalForUF(rate).EfetivoPct, revenue)
	if err != nil {
		return ordersports.OrderTaxes{}, err
	}
	taxes.Difal = &difal
	return taxes, nil
}

// pctOf returns pct% of amount, rounded to centavos through the pricing
// module's single rounding routine so orders and Simulador round identically.
func pctOf(pct string, amount *big.Rat) (float64, error) {
	parsed, err := pricingdomain.ParseRat(pct)
	if err != nil {
		return 0, err
	}
	value := new(big.Rat).Quo(parsed, cem)
	value.Mul(value, amount)
	return strconv.ParseFloat(pricingdomain.FormatRatHalfUp(value, 2), 64)
}
