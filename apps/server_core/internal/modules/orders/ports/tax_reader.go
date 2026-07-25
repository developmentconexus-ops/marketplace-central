package ports

import "context"

// OrderTaxes carries the tax components of an order's decomposição, derived
// from the tenant's pricing profile applied to the order's own revenue. A nil
// field is honest-unknown (ADR-17); a 0 value is an explicit zero — DIFAL is
// legitimately 0 when the tenant has the feature switched off, and that is a
// different fact from "we could not work it out".
type OrderTaxes struct {
	Imposto *float64
	Difal   *float64
}

// TaxReader turns an order's revenue into its imposto/DIFAL amounts using the
// tenant's pricing configuration (regime aliquota, DIFAL toggle + UF table).
// It is the orders-local consumer port for facts the pricing module owns; a
// nil reader leaves both components honest-unknown rather than panicking,
// matching the nil-port idiom of CostReader and ShipmentReader.
//
// destinoUF is the shipment's destination state, "" when unknown. With DIFAL
// enabled and no destino there is no honest rate to apply, so Difal comes back
// nil instead of 0.
type TaxReader interface {
	TaxesForOrder(ctx context.Context, total float64, destinoUF string) (OrderTaxes, error)
}
