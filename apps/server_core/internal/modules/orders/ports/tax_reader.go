package ports

import "context"

// TaxItem is one order item's inputs to the per-item ICMS/PIS-COFINS
// calculation (pricing/domain.TaxesForItem, D-41). InternalProductID nil
// means the item was never linked to a products_mirror row — the reader must
// treat that as a fully unknown ICMSCell for THIS item, which per the
// all-or-nothing aggregation contract takes the WHOLE ORDER's new tax
// components to nil (never a partial sum over only the linked items).
type TaxItem struct {
	InternalProductID *int
	UnitPrice         *float64
	Quantity          int
}

// OrderTaxes carries the tax components of an order's decomposição, derived
// per-item from the D-41 ICMS matrix (pricing/domain.TaxesForItem) and summed
// across the order's items. A nil field is honest-unknown (ADR-17); a 0 value
// is an explicit zero — e.g. ICMSSaida is legitimately 0 for an ST item, and
// that is a different fact from "we could not work it out".
//
// Imposto stays for the D-38 simulator surfaces that still read it; this
// reader no longer populates it (T5: the ProfileSource-based imposto/DIFAL
// calculation is replaced end to end by the per-item D-41 path).
type OrderTaxes struct {
	Imposto       *float64
	Difal         *float64
	ICMSSaida     *float64
	PisCofins     *float64
	RestituicaoST *float64
}

// TaxReader turns an order's items into their aggregate ICMS/DIFAL/PIS-COFINS
// amounts using the D-41 per-item matrix calculation (pricing/domain,
// Task 4). It is the orders-local consumer port for facts the pricing module
// owns; a nil reader leaves every component honest-unknown rather than
// panicking, matching the nil-port idiom of CostReader and ShipmentReader.
//
// destinoUF is the shipment's destination state, "" when unknown — with no
// destino there is no honest cell to resolve, so every component comes back
// nil instead of 0.
type TaxReader interface {
	TaxesForItems(ctx context.Context, items []TaxItem, destinoUF string) (OrderTaxes, error)
}
